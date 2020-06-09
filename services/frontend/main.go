package main

import (
	//"context"
	"encoding/json"
	"flag"
	"fmt"
	//"google.golang.org/grpc"
	"html/template"
	"io/ioutil"
	"lib/common"
	"log"
	"net/http"
	//"time"
)

const (
	serviceName string = "frontend" // Make sure same cfg name in Dockerfile
	version     string = "1.0.1"
)

func main() {
	//ctx := context.Background()
	// Command line stuff
	showversion := flag.Bool("version", false, "display version")
	flag.Parse()
	if *showversion {
		fmt.Printf("Version %s\n", version)
		return
	}

	// Service details
	var c, err = common.LoadConfig(serviceName, "")
	if err != nil {
		fmt.Printf("Cannot load the configuration: %s", err)
	}

	// Now we've parsed all the connection paramaters we can connect to the services
	//mustConnGRPC(ctx, &servers["route-guide"].Conn, *servers["route-guide"].ServiceParams[common.ServiceAddress])

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s\n", version)
	})

	//host,err := os.Hostname()
	//if err != nil {
	//	panic(fmt.Sprintf("No host name: %v\n", err))
	//}
	//addrs, err := net.LookupHost(host)
	//if err != nil {
	//	panic(fmt.Sprintf("No IP addresses for host name %s: %v\n", host, err))
	//}

	// Now connect to the backend stuff
	//common.MustMapEnv(&servers.routeGuideSvcAddr, "ROUTE_GUIDE_SERVICE_ADDR")
	//mustConnGRPC(ctx, &servers.systemSvcConn, servers.systemSvcAddr)

	c.KeyPrefix("system")
	log.Printf("System Service URL:%s", c.ServiceAddress())

	// This path works at the command line, but in GoLand it starts at route
	tpl := template.Must(template.ParseFiles("./templates/serverStatus.html"))
	transport := http.Transport{DisableKeepAlives: false}
	client := &http.Client{Transport: &transport}
	req, _ := http.NewRequest(
		"GET",
		c.ServiceAddress(),
		nil,
	)
	req.Close = false

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		i := &common.ServerInstance{}
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

	c.KeyPrefix("frontend")
	c.Log.Printf("Listening on http://%s:%v", common.GetLocalIP(), c.Port())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", c.Port()), nil))
}

//func mustConnGRPC(ctx context.Context, conn **grpc.ClientConn, addr string) {
//	var err error
//	*conn, err = grpc.DialContext(ctx, addr,
//		grpc.WithInsecure(),
//		grpc.WithTimeout(time.Second*3),
//		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}))
//	if err != nil {
//		panic(fmt.Errorf("grpc: failed to connect %s : %w", addr,err))
//	}
//}
