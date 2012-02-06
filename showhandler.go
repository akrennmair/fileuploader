package main

import (
	"log"
	"net/http"
)

type ShowHandler struct  {
	Persistence PersistenceManager
}

// this handler renders the description page, which contains
// a link to the uploaded file plus the description text that
// was saved with it.
func (h *ShowHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// parse and verify upload ID
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

	// fetch description text and original filename and
	// render the description page.
	desc, err := h.Persistence.GetUploadText(upload_id)
	if err != nil {
		rw.Write(ErrorPage(err.Error()))
		return
	}

	filename, err := h.Persistence.GetUploadFilename(upload_id)
	if err != nil {
		filename = ""
		log.Printf("couldn't retrieve upload filename for %s", upload_id)
	}

	rw.Write(InformationPage(upload_id, desc, filename))
}
