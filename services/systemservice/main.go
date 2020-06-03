package main

import (
	"cloud.google.com/go/compute/metadata"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
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
	defaultPort        = "8081"
	version     string = "1.0.0"
)

//type frontendServer struct {
//	systemSvcAddr string
//	systemSvcConn *grpc.ClientConn
//}

func main() {
	// Command line stuff
	showversion := flag.Bool("version", false, "display version")
	port := flag.String("port", defaultPort, "port to bind")
	flag.Parse()

	if *showversion {
		fmt.Printf("Version %s\n", version)
		return
	}

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s\n", version)
	})

	srvPort := port
	if os.Getenv("PORT") != "" {
		*srvPort = os.Getenv("PORT")
	}
	systemService(*srvPort)
}

func systemService(port string) {
	log.Printf("Starting system service on port %s", port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		i := newInstance()
		raw, _ := httputil.DumpRequest(r, true)
		i.LBRequest = string(raw)
		resp, _ := json.Marshal(i)
		fmt.Fprintf(w, "%s", resp)
	})
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))

}

func mustMapEnv(target *string, envKey string) {
	v := os.Getenv(envKey)
	if v == "" {
		panic(fmt.Sprintf("environment variable %q not set", envKey))
	}
	*target = v
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

//
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
