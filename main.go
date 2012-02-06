package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// setup logging to stdout; show exact time in log, including lines in sourcecode
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	// setup URI multiplexing
	servemux := http.NewServeMux()

	// the / prefix is special... it handles all the request that didn't match the other prefixes
	// that's why we create the custom StartPageHandler.
	servemux.Handle("/", &StartPageHandler{Other: http.FileServer(http.Dir("htdocs"))})
	servemux.Handle("/progress/", &ProgressHandler{})
	servemux.Handle("/upload/", &PostUploadHandler{})
	servemux.Handle("/show/", &ShowHandler{})
	servemux.Handle("/savedesc/", &SaveDescHandler{})
	servemux.Handle("/files/", &DeliverFileHandler{})
	servemux.Handle("/requpid", &RequestIDHandler{})

	// create HTTP server and run it on port 8000.
	httpsrv := &http.Server{Handler: servemux, Addr: "0.0.0.0:8000"}
	httpsrv.ListenAndServe()
}

// this function creates wrapper functions for http.ServeMux handler functions
// that log all requests that are directed to the handler function specified
// as argument.
func logReq(f func(rw http.ResponseWriter, r *http.Request)) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		LogHTTPRequest(r)
		f(rw, r)
	}
}

func MakeResponseNonCachable(rw *http.ResponseWriter) {
	(*rw).Header()["Expires"] = []string{"Sat, 1 Jan 2005 00:00:00 GMT"}
	(*rw).Header()["Last-Modified"] = []string{time.Now().Format(time.RFC1123)}
	(*rw).Header()["Cache-Control"] = []string{"no-cache, must-revalidate, no-store"}
	(*rw).Header()["Pragma"] = []string{"no-cache"}
}

// helper function to log an HTTP request.
func LogHTTPRequest(r *http.Request) {
	ua := r.Header["User-Agent"]
	if len(ua) == 0 {
		ua = []string{"-"}
	}
	log.Printf("HTTP Request: %s %s %s %s", r.RemoteAddr, r.Method, r.URL.Path, ua[0])
}
