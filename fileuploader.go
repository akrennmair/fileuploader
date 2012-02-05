package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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
	servemux.Handle("/", &StartPageHandler{StartPage: logReq(UploadDialog), Other: http.FileServer(http.Dir("htdocs"))})
	servemux.HandleFunc("/progress/", logReq(Progress))
	servemux.HandleFunc("/upload/", logReq(PostUpload))
	servemux.HandleFunc("/show/", logReq(Show))
	servemux.HandleFunc("/savedesc/", logReq(SaveDesc))
	servemux.HandleFunc("/files/", logReq(DeliverFile))
	servemux.HandleFunc("/requpid", logReq(RequestUploadID))

	// create HTTP server and run it on port 8000.
	httpsrv := &http.Server{Handler: servemux, Addr: "0.0.0.0:8000"}
	httpsrv.ListenAndServe()
}

// custom Handler... as mentioned before, the / prefix is special, so instead of serving
// "everything else" from the UploadDialog handler, we create our own handler that handles
// exactly "/" using the UploadDialog handler, and everything else with the default
// http.Fileserver from the htdocs subdirectory. That means that if we want to make
// static files available, we can simple place them there.
type StartPageHandler struct {
	StartPage http.HandlerFunc
	Other     http.Handler
}

func (h *StartPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		h.StartPage(w, r)
	} else {
		LogHTTPRequest(r)
		h.Other.ServeHTTP(w, r)
	}
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

// this function render the upload dialog and sends it to the client. 
// really nothing fancy.
func UploadDialog(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write(UploadPage())
}

// this functions returns the current upload progress for an upload ID
func Progress(rw http.ResponseWriter, r *http.Request) {
	// first, the upload ID is parsed and validated. Requests with invalid
	// upload IDs are rejected.
	log.Printf("Progress called")
	upload_id, err := GetUploadID(r.URL.Path)
	if err != nil {
		log.Printf("Progress: UploadID doesn't validate: %s", err.Error())
		rw.WriteHeader(http.StatusOK)
		rw.Write(ErrorPage(err.Error()))
		return
	}

	// then the upload progress is fetched and sent to the client.
	percent, err := GetUploadProgress(upload_id)
	if err != nil {
		log.Printf("Progress: there was an error fetching the upload progress: %s", err.Error())
		percent = -1
	}
	percent_str := fmt.Sprintf("%d", percent)
	rw.Header()["Content-Type"] = []string{"text/plain"}
	rw.Header()["Content-Length"] = []string{fmt.Sprintf("%d", len(percent_str))}
	MakeResponseNonCachable(&rw)
	log.Printf("Progress: percent = %d %s", percent, percent_str)

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(percent_str))
	log.Printf("Progress: finished sending data")
}

func MakeResponseNonCachable(rw *http.ResponseWriter) {
	(*rw).Header()["Expires"] = []string{"Sat, 1 Jan 2005 00:00:00 GMT"}
	(*rw).Header()["Last-Modified"] = []string{time.Now().Format(time.RFC1123)}
	(*rw).Header()["Cache-Control"] = []string{"no-cache, must-revalidate, no-store"}
	(*rw).Header()["Pragma"] = []string{"no-cache"}
}

// this function receives the uploaded file, parses the multipart/form-data
// and stores the file.
func PostUpload(rw http.ResponseWriter, r *http.Request) {
	// all methods other than POST are rejected. They don't make any sense, anyway.
	if r.Method != "POST" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write(ErrorPage("POST expected"))
		return
	}
	// then, the upload ID is validated. Invalid ones are rejected.
	upload_id, err := GetUploadID(r.URL.Path)
	if err != nil {
		rw.Write(ErrorPage(err.Error()))
		return
	}

	// in order to track the upload progress, we need to know how much is
	// getting uploaded. The Content-Length header is a rather reliable source
	// for that.
	content_length, _ := strconv.ParseUint(r.Header["Content-Length"][0], 10, 64) // TODO: check error

	// we then wrap the r.Body io.ReadClose with our own custom ProgressReadCloser
	// that will count how much data was received and will continuously update
	// the upload progress.
	r.Body = NewProgressReadCloser(r.Body, content_length, upload_id)

	// and then we start parsing the multipart/form-data request body.
	mpr, err := r.MultipartReader()
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
			// only the first part will be saved, all others are discarded.
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
			// when upload is finished, we also store the original filename.
			if err := SaveUploadFilename(upload_id, part.FileName()); err != nil {
				log.Printf("couldn't save upload filename for %s", upload_id)
			}
		}
	}

	// more than one part in the multipart-formdata body is odd, so we log it.
	// this might be an attacker or a buggy client.
	if part_count > 1 {
		log.Printf("found %d parts, saved only first one.", part_count)
	}

	// respond with HTTP 200, but no other information (nothing is shown on the
	// client's side, anyway, this all goes into a hidden iframe).
	rw.WriteHeader(http.StatusOK)
}

// this function saves the description text.
func SaveDesc(rw http.ResponseWriter, r *http.Request) {
	// reject anything but POST requests.
	if r.Method != "POST" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write(ErrorPage("POST expected"))
		return
	}
	// parse and verify upload ID
	upload_id, err := GetUploadID(r.URL.Path)
	if err != nil {
		rw.Write(ErrorPage(err.Error()))
		return
	}

	// according to the JavaScript client logic, the file must be
	// uploaded before the description is being saved, so we explicitly
	// disallow saving a description text for a file that hasn't been
	// uploaded yet as this is bogus and somebody might be emulating
	// the actual client's requests to the server without following
	// the proper logic.
	if !UploadExists(upload_id) {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte("Forbidden"))
		return
	}

	r.ParseForm()

	text := r.Form.Get("input_desc")

	// saved the description text, and if that went fine, redirect
	// to the description page.
	if err := SaveUploadText(upload_id, text); err != nil {
		log.Printf("saving description failed: %s", err.Error())
		rw.WriteHeader(http.StatusOK)
		rw.Write(ErrorPage("saving description failed: " + err.Error()))
	} else {
		rw.Header()["Location"] = []string{"/show/" + upload_id}
		rw.WriteHeader(http.StatusFound)
	}
}

// this function renders the description page, which contains
// a link to the uploaded file plus the description text that
// was saved with it.
func Show(rw http.ResponseWriter, r *http.Request) {
	// parse and verify upload ID
	upload_id, err := GetUploadID(r.URL.Path)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write(ErrorPage(err.Error()))
		return
	}

	// reject requests if the upload isn't complete yet
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

	// fetch description text and original filename and
	// render the description page.
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

// this function delivers an uploaded file, identified by its
// upload ID to the client, with content-type application/octet-stream.
func DeliverFile(rw http.ResponseWriter, r *http.Request) {
	// parse and verify upload iD
	upload_id, err := GetUploadID(r.URL.Path)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write(ErrorPage(err.Error()))
		return
	}

	// reject requests if the upload isn't complete yet
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

	// open uploaded file.
	if f, err := OpenUpload(upload_id); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write(ErrorPage(err.Error()))
	} else {
		// if an original filename was saved, set it as Content-Disposition header
		if filename, err := GetUploadFilename(upload_id); err == nil {
			rw.Header()["Content-Disposition"] = []string{"attachment; filename=" + filename}
		}
		// if the filesize can be determined, set the Content-Length header
		if size, err := GetUploadSize(upload_id); err == nil {
			rw.Header()["Content-Length"] = []string{strconv.FormatInt(size, 10)}
		}
		// deliver uploaded file as application/octet-stream
		rw.Header()["Content-Type"] = []string{"application/octet-stream"}
		rw.WriteHeader(http.StatusOK)
		io.Copy(rw, f)
		f.Close()
	}
}

func RequestUploadID(rw http.ResponseWriter, r *http.Request) {
	if upload_id, err := GenerateUploadID(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	} else {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(upload_id))
	}
}

// helper function to log an HTTP request.
func LogHTTPRequest(r *http.Request) {
	ua := r.Header["User-Agent"]
	if len(ua) == 0 {
		ua = []string{"-"}
	}
	log.Printf("HTTP Request: %s %s %s %s", r.RemoteAddr, r.Method, r.URL.Path, ua[0])
}
