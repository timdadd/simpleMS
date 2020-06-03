package main

import (
	"cloud.google.com/go/compute/metadata"
	"encoding/json"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

type Instance struct {
	Id         string
	Name       string
	Version    string
	Hostname   string
	Zone       string
	Project    string
	InternalIP string
	ExternalIP string
	LBRequest  string
	ClientIP   string
	Error      string
}

const (
	defaultPort        = "8080"
	version     string = "2.0.0"
)

type frontendServer struct {
	systemSvcAddr string
	systemSvcConn *grpc.ClientConn
}

func main() {

	// Command line stuff
	showversion := flag.Bool("version", false, "display version")
	//frontend := flag.Bool("frontend", false, "run in frontend mode")
	srvPort := flag.String("port", defaultPort, "port to bind")
	systemService := flag.String("system-service", "http://127.0.0.1:8081", "hostname of backend server")
	flag.Parse()

	if *showversion {
		fmt.Printf("Version %s\n", version)
		return
	}

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s\n", version)
	})

	if os.Getenv("PORT") != "" {
		mustMapEnv(srvPort,"PORT")
	}
	//host,err := os.Hostname()
	//if err != nil {
	//	panic(fmt.Sprintf("No host name: %v\n", err))
	//}
	//addrs, err := net.LookupHost(host)
	//if err != nil {
	//	panic(fmt.Sprintf("No IP addresses for host name %s: %v\n", host, err))
	//}
	addr := GetLocalIP()
	mustMapEnv(&addr, "LISTEN_ADDR")
	log.Printf("http://%s:%s",addr,*srvPort)
	//svc := new(frontendServer)

	// Now connect to the backend stuff
	mustMapEnv(systemService, "SYSTEM_SERVICE_ADDR")
	//mustConnGRPC(ctx, &svc.systemSvcConn, svc.systemSvcAddr)

	frontendMode(*srvPort, *systemService)
}

func frontendMode(port string, backendURL string) {
	log.Printf("Starting frontend on port %s", port)
	tpl := template.Must(template.New("out").Parse(html))

	transport := http.Transport{DisableKeepAlives: false}
	client := &http.Client{Transport: &transport}
	req, _ := http.NewRequest(
		"GET",
		backendURL,
		nil,
	)
	req.Close = false

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		i := &Instance{}
		resp, err := client.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "Error: %s\n", err.Error())
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %s\n", err.Error())
			return
		}
		err = json.Unmarshal([]byte(body), i)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %s\n", err.Error())
			return
		}
		tpl.Execute(w, i)
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		resp, err := client.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "Backend could not be connected to: %s", err.Error())
			return
		}
		defer resp.Body.Close()
		ioutil.ReadAll(resp.Body)
		w.WriteHeader(http.StatusOK)
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

type assigner struct {
	err error
}

func (a *assigner) assign(getVal func() (string, error)) string {
	if a.err != nil {
		return ""
	}
	s, err := getVal()
	if err != nil {
		a.err = err
	}
	return s
}

func newInstance() *Instance {
	var i = new(Instance)
	if !metadata.OnGCE() {
		i.Error = "Not running on GCE"
		return i
	}

	a := &assigner{}
	i.Id = a.assign(metadata.InstanceID)
	i.Zone = a.assign(metadata.Zone)
	i.Name = a.assign(metadata.InstanceName)
	i.Hostname = a.assign(metadata.Hostname)
	i.Project = a.assign(metadata.ProjectID)
	i.InternalIP = a.assign(metadata.InternalIP)
	i.ExternalIP = a.assign(metadata.ExternalIP)
	i.Version = version

	if a.err != nil {
		i.Error = a.err.Error()
	}
	return i
}

func mustMapEnv(target *string, envKey string) {
	v := os.Getenv(envKey)
	if v == "" {
		v = *target
	}
	if v == "" {
		panic(fmt.Sprintf("environment variable %q not set", envKey))
	}
	*target = v
	log.Printf("%q=%q", envKey, v)
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}