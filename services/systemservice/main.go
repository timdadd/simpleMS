package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"lib/utils"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

const (
	defaultPort        = "8081"
	version     string = "1.0.0"
)

func main() {
	// Command line stuff
	showversion := flag.Bool("version", false, "display version")
	srvPort := flag.String("port", defaultPort, "port to bind")
	flag.Parse()

	if *showversion {
		fmt.Printf("Version %s\n", version)
		return
	}

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s\n", version)
	})

	addr := utils.GetLocalIP()
	utils.MustMapEnv(&addr, "LISTEN_ADDR")
	log.Printf("http://%s:%s", addr, *srvPort)

	if os.Getenv("PORT") != "" {
		*srvPort = os.Getenv("PORT")
	}
	systemService(*srvPort)
}

func systemService(port string) {
	log.Printf("Starting system service on port %s", port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		i := utils.NewServerInstance(version)
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
