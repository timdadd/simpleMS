package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"lib/utils"
	"log"
	"net/http"
	"os"
)

const (
	defaultPort        = "8080"
	version     string = "2.0.0"
)

//type frontendServer struct {
//	systemSvcAddr string
//	systemSvcConn *grpc.ClientConn
//}

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
		utils.MustMapEnv(srvPort, "PORT")
	}
	//host,err := os.Hostname()
	//if err != nil {
	//	panic(fmt.Sprintf("No host name: %v\n", err))
	//}
	//addrs, err := net.LookupHost(host)
	//if err != nil {
	//	panic(fmt.Sprintf("No IP addresses for host name %s: %v\n", host, err))
	//}
	addr := utils.GetLocalIP()
	utils.MustMapEnv(&addr, "LISTEN_ADDR")
	log.Printf("http://%s:%s", addr, *srvPort)
	//svc := new(frontendServer)

	// Now connect to the backend stuff
	utils.MustMapEnv(systemService, "SYSTEM_SERVICE_ADDR")
	//mustConnGRPC(ctx, &svc.systemSvcConn, svc.systemSvcAddr)

	frontendMode(*srvPort, *systemService)
}

func frontendMode(port string, backendURL string) {
	log.Printf("Starting frontend on port %s", port)
	// This path works at the command line, but in GoLand it starts at route
	tpl := template.Must(template.ParseFiles("./templates/serverStatus.html"))

	transport := http.Transport{DisableKeepAlives: false}
	client := &http.Client{Transport: &transport}
	req, _ := http.NewRequest(
		"GET",
		backendURL,
		nil,
	)
	req.Close = false

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		i := &utils.ServerInstance{}
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
