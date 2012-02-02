package main

import (
	"http"
	"fmt"
	"strings"
	"io"
	"os"
)

func main() {

	servemux := http.NewServeMux()
	servemux.HandleFunc("/", logReq(uploadDialog))
	servemux.HandleFunc("/progress/", logReq(progress))
	servemux.HandleFunc("/upload/", logReq(postUpload))
	servemux.HandleFunc("/show/", logReq(show))
	servemux.HandleFunc("/files/", logReq(deliverFile))

	httpsrv := &http.Server{Handler: servemux, Addr: "0.0.0.0:8000"}
	httpsrv.ListenAndServe()
}

// generate wrapper functions for request logging
func logReq(f func(rw http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		ua := r.Header["User-Agent"]
		if len(ua) == 0 {
			ua = []string{"-"}
		}
		fmt.Printf("%s %s %s %s\n", r.RemoteAddr, r.Method, r.RawURL, ua[0])
		f(rw, r)
	}
}

// renders the upload dialog
func uploadDialog(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)

	if upid, err := generateUploadID(); err != nil {
		rw.Write(ErrorPage(err.String()))
	} else {
		rw.Write(UploadPage(upid))
	}
}

// returns the progress for the current upload
func progress(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	slash_pos := strings.LastIndex(r.RawURL, "/")
	if slash_pos < 0 {
		rw.Write(ErrorPage("no Upload ID"))
		return
	}
	upload_id := r.RawURL[slash_pos+1:]
	if !verifyUploadID(upload_id) {
		rw.Write(ErrorPage("invalid Upload ID"))
		return
	}
	rw.Write([]byte("progress: upload_id = " + upload_id))
}

// accepts the POST with the uploaded file
func postUpload(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write(ErrorPage("POST expected"))
		return
	}
	slash_pos := strings.LastIndex(r.RawURL, "/")
	if slash_pos < 0 {
		rw.WriteHeader(http.StatusOK)
		rw.Write(ErrorPage("no Upload ID"))
		return
	}
	upload_id := r.RawURL[slash_pos+1:]
	if !verifyUploadID(upload_id) {
		rw.WriteHeader(http.StatusOK)
		rw.Write(ErrorPage("invalid Upload ID"))
		return
	}
	io.Copy(os.Stdout, r.Body)

	rw.Write([]byte("postUpload: upload_id = " + upload_id))
}

// show link to file + description
func show(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("show: "))
	rw.Write([]byte(r.RawURL))
}

// deliver file from file system by Upload ID
func deliverFile(rw http.ResponseWriter, r *http.Request) {
	// TODO: implement
}
