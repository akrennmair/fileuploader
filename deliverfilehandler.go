package main

import (
	"net/http"
	"strconv"
	"io"
)

type DeliverFileHandler struct { 
	Persistence PersistenceManager
}

// this handler delivers an uploaded file, identified by its
// upload ID to the client, with content-type application/octet-stream.
func (h *DeliverFileHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// parse and verify upload iD
	upload_id, err := GetUploadID(r.URL.Path)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write(ErrorPage(err.Error()))
		return
	}

	// reject requests if the upload isn't complete yet
	percent, err := h.Persistence.GetUploadProgress(upload_id)
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
	if f, err := h.Persistence.OpenUpload(upload_id); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write(ErrorPage(err.Error()))
	} else {
		// if an original filename was saved, set it as Content-Disposition header
		if filename, err := h.Persistence.GetUploadFilename(upload_id); err == nil {
			rw.Header()["Content-Disposition"] = []string{"attachment; filename=" + filename}
		}
		// if the filesize can be determined, set the Content-Length header
		if size, err := h.Persistence.GetUploadSize(upload_id); err == nil {
			rw.Header()["Content-Length"] = []string{strconv.FormatInt(size, 10)}
		}
		// deliver uploaded file as application/octet-stream
		rw.Header()["Content-Type"] = []string{"application/octet-stream"}
		rw.WriteHeader(http.StatusOK)
		io.Copy(rw, f)
		f.Close()
	}
}
