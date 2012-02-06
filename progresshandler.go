package main

import (
	"log"
	"net/http"
	"fmt"
)

type ProgressHandler struct { 
	Persistence PersistenceManager
}

// this handler returns the current upload progress for an upload ID
func (h *ProgressHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// first, the upload ID is parsed and validated. Requests with invalid
	// upload IDs are rejected.
	upload_id, err := GetUploadID(r.URL.Path)
	if err != nil {
		log.Printf("Progress: UploadID doesn't validate: %s", err.Error())
		rw.WriteHeader(http.StatusOK)
		rw.Write(ErrorPage(err.Error()))
		return
	}

	// then the upload progress is fetched and sent to the client.
	percent, err := h.Persistence.GetUploadProgress(upload_id)
	if err != nil {
		log.Printf("Progress: there was an error fetching the upload progress: %s", err.Error())
		percent = -1
	}
	percent_str := fmt.Sprintf("%d", percent)
	rw.Header()["Content-Type"] = []string{"text/plain"}
	rw.Header()["Content-Length"] = []string{fmt.Sprintf("%d", len(percent_str))}
	MakeResponseNonCachable(&rw)

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(percent_str))
}
