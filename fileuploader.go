package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
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
		fmt.Printf("%s %s %s %s\n", r.RemoteAddr, r.Method, r.URL.Path, ua[0])
		f(rw, r)
	}
}

// renders the upload dialog
func uploadDialog(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)

	if upid, err := GenerateUploadID(); err != nil {
		rw.Write(ErrorPage(err.Error()))
	} else {
		rw.Write(UploadPage(upid))
	}
}

// returns the progress for the current upload
func progress(rw http.ResponseWriter, r *http.Request) {
	upload_id, err := GetUploadID(r.URL.Path)
	if err != nil {
		rw.WriteHeader(http.StatusOK)
		rw.Write(ErrorPage(err.Error()))
		return
	}
	percent, err := ReadProgress(upload_id)
	if err != nil {
		percent = -1
	}
	rw.Header()["Content-Type"] = []string{"text/plain"}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(fmt.Sprintf("%d", percent)))
}

// accepts the POST with the uploaded file
func postUpload(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write(ErrorPage("POST expected"))
		return
	}
	upload_id, err := GetUploadID(r.URL.Path)
	if err != nil {
		rw.Write(ErrorPage(err.Error()))
		return
	}

	content_length, _ := strconv.ParseUint(r.Header["Content-Length"][0], 10, 64) // TODO: check error

	r.Body = NewProgressReadCloser(r.Body, content_length, upload_id)

	fmt.Printf("before creating MultipartReader\n")

	mpr, err := r.MultipartReader() // TODO: check error
	if err != nil {
		rw.WriteHeader(http.StatusOK)
		rw.Write(ErrorPage(err.Error()))
		return
	}

	for {
		fmt.Printf("handling next part\n")
		if part, err := mpr.NextPart(); err == io.EOF {
			break
		} else {
			if f, err := os.OpenFile("files/" + upload_id, os.O_CREATE | os.O_WRONLY, 0444); err == nil {
				io.Copy(f, part)
				f.Close()
			} else {
				rw.WriteHeader(http.StatusOK)
				rw.Write([]byte("Opening file for " + upload_id + " failed"))
				return
			}
		}
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("postUpload: upload_id = " + upload_id))
}

// show link to file + description
func show(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("show: "))
	rw.Write([]byte(r.URL.Path))
}

// deliver file from file system by Upload ID
func deliverFile(rw http.ResponseWriter, r *http.Request) {
	upload_id, err := GetUploadID(r.URL.Path)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write(ErrorPage(err.Error()))
		return
	}

	if f, err := os.Open("files/" + upload_id); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write(ErrorPage(err.Error()))
		return
	} else {
		rw.Header()["Content-Type"] = []string{"application/octet-stream"}
		rw.WriteHeader(http.StatusOK)
		io.Copy(rw, f)
		f.Close()
	}
}
