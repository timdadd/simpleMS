package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"lib/common"
	"log"
	"net/http"
	"net/http/httputil"
)

const (
	version string = "1.0.0"
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
	c, err := common.LoadConfig("frontend.yaml", "")
	if err != nil {
		fmt.Printf("Cannot load the configuration: %s", err)
	}

	// Now we've parsed all the connection paramaters we can connect to the services
	//mustConnGRPC(ctx, &servers["route-guide"].Conn, *servers["route-guide"].ServiceParams[common.ServiceAddress])

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s\n", version)
	})

	c.KeyPrefix("system")
	c.Log.Debug("http://", common.GetLocalIP(), ":", c.GetIntKey(c.Key.Port))

	c.Log.Debug("Starting system service on port ", c.GetIntKey(c.Key.Port))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		i := common.NewServerInstance(version)
		raw, _ := httputil.DumpRequest(r, true)
		i.LBRequest = string(raw)
		resp, _ := json.Marshal(i)
		fmt.Fprintf(w, "%s", resp)
	})
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", c.GetIntKey(c.Key.Port)), nil))
}
