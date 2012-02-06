package main

import (
	"log"
	"net/http"
	"io"
	"strconv"
)

type PostUploadHandler struct { }

// this handler receives the uploaded file, parses the multipart/form-data
// and stores the file.
func (h *PostUploadHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
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
