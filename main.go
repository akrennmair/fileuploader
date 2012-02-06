package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// setup logging to stdout; show exact time in log, including lines in sourcecode
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	pm := NewFilePersistenceManager()

	// setup URI multiplexing
	servemux := http.NewServeMux()

	// the / prefix is special... it handles all the request that didn't match the other prefixes
	// that's why we create the custom StartPageHandler.
	servemux.Handle("/", NewLoggingHandler(&StartPageHandler{Other: http.FileServer(http.Dir("htdocs"))}))
	servemux.Handle("/progress/", NewLoggingHandler(&ProgressHandler{Persistence: pm}))
	servemux.Handle("/upload/", NewLoggingHandler(&PostUploadHandler{Persistence: pm}))
	servemux.Handle("/show/", NewLoggingHandler(&ShowHandler{Persistence: pm}))
	servemux.Handle("/savedesc/", NewLoggingHandler(&SaveDescHandler{Persistence: pm}))
	servemux.Handle("/files/", NewLoggingHandler(&DeliverFileHandler{Persistence: pm}))
	servemux.Handle("/requpid", NewLoggingHandler(&RequestIDHandler{}))

	// create HTTP server and run it on port 8000.
	httpsrv := &http.Server{Handler: servemux, Addr: "0.0.0.0:8000"}
	httpsrv.ListenAndServe()
}

type LoggingHandler struct {
	handler http.Handler
}

// this function creates a wrapper handler that logs all HTTP requests for the specified
// http.Handler
func NewLoggingHandler(h http.Handler) http.Handler {
	return &LoggingHandler{handler: h}
}

func (h *LoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	LogHTTPRequest(r)
	h.handler.ServeHTTP(w, r)
}

// helper function to log an HTTP request.
func LogHTTPRequest(r *http.Request) {
	ua := r.Header["User-Agent"]
	if len(ua) == 0 {
		ua = []string{"-"}
	}
	log.Printf("HTTP Request: %s %s %s %s", r.RemoteAddr, r.Method, r.URL.Path, ua[0])
}
