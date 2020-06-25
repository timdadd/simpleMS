package main

import (
	"cloud.google.com/go/storage"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"google.golang.org/grpc"
	"lib/common"
	"net/http"
	//"time"
)

const (
	serviceName string = "frontend" // Make sure same cfg name in Dockerfile
	svcBook     string = "book"     // Name used for book service
	svcSystem   string = "string"   // Name used for system service
	version     string = "1.0.1"

	cookieMaxAge = 60 * 60 * 48

	cookiePrefix    = "simplems_"
	cookieSessionID = cookiePrefix + "session-id"
)

type ctxKeySessionID struct{}

type frontendServer struct {
	bookSvcConn *grpc.ClientConn

	StorageBucket     *storage.BucketHandle
	StorageBucketName string

	log *logrus.Logger
}

func main() {
	//ctx := context.Background()
	// Command line stuff
	showversion := flag.Bool("version", false, "display version")
	flag.Parse()
	if *showversion {
		fmt.Printf("Version %s\n", version)
		return
	}

	// Load configuration & Service details
	var c, err = common.LoadConfig(serviceName, "")
	if err != nil {
		fmt.Printf("Cannot load the configuration: %s", err)
	}

	// Create connections to the RPC services
	c.ConnGRPC(svcBook)
	svc := new(frontendServer)
	svc.bookSvcConn = c.SvcConn[svcBook]
	svc.log = c.Log
	svc.registerHandlers(c)
	svc.log.Debug("Connected to book service")
}

func (fe frontendServer) registerHandlers(c common.AppConfig) {
	// Use gorilla/mux for rich routing.
	// See https://www.gorillatoolkit.org/pkg/mux.
	r := mux.NewRouter()
	r.StrictSlash(true) // Odd that this is true when now the strictness is false
	// Homepage redirects to /books
	r.Handle("/", http.RedirectHandler("/books", http.StatusFound))

	// GET books, HEAD is like GET but without the body returned
	r.HandleFunc("/books", fe.listBook).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/books/add", fe.addBook).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/books/templates", fe.bookTemplates).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/books/{id:[0-9a-zA-Z_\\-]+}", fe.bookDetail).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/books/{id:[0-9a-zA-Z_\\-]+}/edit", fe.editBook).Methods(http.MethodGet, http.MethodHead)

	// POST/PUT books
	r.HandleFunc("/books/add", fe.createBook).Methods(http.MethodPost)
	r.HandleFunc("/books/{id:[0-9a-zA-Z_\\-]+}", fe.updateBook).Methods(http.MethodPost, http.MethodPut)
	r.HandleFunc("/books/{id:[0-9a-zA-Z_\\-]+}:delete", fe.deleteBook).Methods(http.MethodPost)

	// Admin stuff
	r.HandleFunc("/version", fe.version).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/logout", fe.logoutHandler).Methods(http.MethodGet)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "User-agent: *\nDisallow: /") })
	r.HandleFunc("/_healthz", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "ok") })

	var handler http.Handler = r
	handler = &logHandler{log: c.Log, next: handler} // add logging
	handler = ensureSessionID(handler)               // add session ID
	handler = &ochttp.Handler{                       // add opencensus instrumentation
		Handler:     handler,
		Propagation: &b3.HTTPFormat{}}

	c.KeyPrefix("frontend")
	addr := c.ListenAddress()
	port := c.Port()
	c.Log.Infof("starting server on %s:%v", addr, port)

	c.KeyPrefix("frontend")
	c.Log.Printf("Listening on http://%s:%v", common.GetLocalIP(), c.Port())
	c.Log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", c.Port()), handler))

	//c.Log.Fatal(http.ListenAndServe(addr+":"+port, handler))

	//r.Methods("GET").Path("/logs").Handler(appHandler(fe.sendLog))
	//r.Methods("GET").Path("/errors").Handler(appHandler(fe.sendError))

	// Delegate all of the HTTP routing and serving to the gorilla/mux router.
	// Log all requests using the standard Apache format.
	//http.Handle("/", handlers.CombinedLoggingHandler(fe.logWriter, r))
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

//c.KeyPrefix("system")
//log.Printf("System Service URL:%s", c.ServiceAddress())
//
//// This path works at the command line, but in GoLand it starts at route
//tpl := template.Must(template.ParseFiles("./templates/serverStatus.gohtml"))
//transport := http.Transport{DisableKeepAlives: false}
//client := &http.Client{Transport: &transport}
//req, _ := http.NewRequest(
//	"GET",
//	c.ServiceAddress(),
//	nil,
//)
//req.Close = false
//
//http.HandleFunc("/system", func(w http.ResponseWriter, r *http.Request) {
//	i := &common.ServerInstance{}
//	resp, err := client.Do(req)
//	if err != nil {
//		w.WriteHeader(http.StatusServiceUnavailable)
//		fmt.Fprintf(w, "Error: %s\n", err.Error())
//		return
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		fmt.Fprintf(w, "Error: %s\n", err.Error())
//		return
//	}
//	err = json.Unmarshal([]byte(body), i)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		fmt.Fprintf(w, "Error: %s\n", err.Error())
//		return
//	}
//	tpl.Execute(w, i)
//})
//
//http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
//	resp, err := client.Do(req)
//	if err != nil {
//		w.WriteHeader(http.StatusServiceUnavailable)
//		fmt.Fprintf(w, "Backend could not be connected to: %s", err.Error())
//		return
//	}
//	defer resp.Body.Close()
//	ioutil.ReadAll(resp.Body)
//	w.WriteHeader(http.StatusOK)
//})
//
//c.KeyPrefix("frontend")
//c.Log.Printf("Listening on http://%s:%v", common.GetLocalIP(), c.Port())
//c.Log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", c.Port()), nil))
