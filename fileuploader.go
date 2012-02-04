package main

import (
	"fmt"
	"net/http"
	"strconv"
	"io"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	servemux := http.NewServeMux()
	servemux.HandleFunc("/", logReq(UploadDialog))
	servemux.HandleFunc("/progress/", logReq(Progress))
	servemux.HandleFunc("/upload/", logReq(PostUpload))
	servemux.HandleFunc("/show/", logReq(Show))
	servemux.HandleFunc("/savedesc/", logReq(SaveDesc))
	servemux.HandleFunc("/files/", logReq(DeliverFile))

	httpsrv := &http.Server{Handler: servemux, Addr: "0.0.0.0:8000"}
	httpsrv.ListenAndServe()
}

// generate wrapper functions for request logging
func logReq(f http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ua := r.Header["User-Agent"]
		if len(ua) == 0 {
			ua = []string{"-"}
		}
		log.Printf("HTTP Request: %s %s %s %s", r.RemoteAddr, r.Method, r.URL.Path, ua[0])
		f(rw, r)
	}
}

// renders the upload dialog
func UploadDialog(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)

	if upload_id, err := GenerateUploadID(); err != nil {
		rw.Write(ErrorPage(err.Error()))
	} else {
		rw.Write(UploadPage(upload_id))
	}
}

// returns the progress for the current upload
func Progress(rw http.ResponseWriter, r *http.Request) {
	upload_id, err := GetUploadID(r.URL.Path)
	if err != nil {
		rw.WriteHeader(http.StatusOK)
		rw.Write(ErrorPage(err.Error()))
		return
	}
	percent, err := GetUploadProgress(upload_id)
	if err != nil {
		percent = -1
	}
	rw.Header()["Content-Type"] = []string{"text/plain"}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(fmt.Sprintf("%d", percent)))
}

// accepts the POST with the uploaded file
func PostUpload(rw http.ResponseWriter, r *http.Request) {
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

	mpr, err := r.MultipartReader() // TODO: check error
	if err != nil {
		rw.WriteHeader(http.StatusOK)
		rw.Write(ErrorPage(err.Error()))
		return
	}

	part_count := 0
	for {
		if part, err := mpr.NextPart(); err == io.EOF {
			break
		} else {
			part_count++
			if part_count > 1 {
				continue
			}
			if f, err := OpenUploadWritable(upload_id); err == nil {
				io.Copy(f, part)
				f.Close()
			} else {
				rw.WriteHeader(http.StatusOK)
				rw.Write([]byte("Opening file for " + upload_id + " failed"))
				return
			}
			if err := SaveUploadFilename(upload_id, part.FileName()); err != nil {
				log.Printf("couldn't save upload filename for %s", upload_id)
			}
		}
	}

	if part_count > 1 {
		log.Printf("found %d parts, saved only first one.", part_count)
	}

	rw.WriteHeader(http.StatusOK)
}

func SaveDesc(rw http.ResponseWriter, r *http.Request) {
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
	if !UploadExists(upload_id) {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte("Forbidden"))
		return
	}

	r.ParseForm()

	text := r.Form.Get("input_desc")
	if err := SaveUploadText(upload_id, text); err != nil {
		log.Printf("saving description failed: %s", err.Error())
		rw.WriteHeader(http.StatusOK)
		rw.Write(ErrorPage("saving description failed: " + err.Error()))
	} else {
		rw.Header()["Location"] = []string{"/show/" + upload_id}
		rw.WriteHeader(http.StatusFound)
	}
}

// show link to file + description
func Show(rw http.ResponseWriter, r *http.Request) {
	upload_id, err := GetUploadID(r.URL.Path)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write(ErrorPage(err.Error()))
		return
	}

	percent, err := GetUploadProgress(upload_id)
	if err != nil {
		rw.Write(ErrorPage(err.Error()))
		return
	}
	if percent < 100 {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte("Forbidden"))
		return
	}

	desc, err := GetUploadText(upload_id)
	if err != nil {
		rw.Write(ErrorPage(err.Error()))
		return
	}

	filename, err := GetUploadFilename(upload_id)
	if err != nil {
		filename = ""
		log.Printf("couldn't retrieve upload filename for %s", upload_id)
	}

	rw.Write(InformationPage(upload_id, desc, filename))
}

// deliver file from file system by Upload ID
func DeliverFile(rw http.ResponseWriter, r *http.Request) {
	upload_id, err := GetUploadID(r.URL.Path)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write(ErrorPage(err.Error()))
		return
	}

	percent, err := GetUploadProgress(upload_id)
	if err != nil {
		rw.Write(ErrorPage(err.Error()))
		return
	}
	if percent < 100 {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte("Forbidden"))
		return
	}

	if f, err := OpenUpload(upload_id); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write(ErrorPage(err.Error()))
	} else {
		if filename, err := GetUploadFilename(upload_id); err == nil {
			rw.Header()["Content-Disposition"] = []string{"attachment; filename=" + filename}
		}
		if size, err := GetUploadSize(upload_id); err == nil {
			rw.Header()["Content-Length"] = []string{strconv.FormatInt(size,10)}
		}
		rw.Header()["Content-Type"] = []string{"application/octet-stream"}
		rw.WriteHeader(http.StatusOK)
		io.Copy(rw, f)
		f.Close()
	}
}
