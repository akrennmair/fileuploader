package main

import (
	"log"
	"net/http"
)

type SaveDescHandler struct { 
	Persistence PersistenceManager
}

// this handler saves the description text.
func(h *SaveDescHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
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
	if !h.Persistence.UploadExists(upload_id) {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte("Forbidden"))
		return
	}

	r.ParseForm()

	text := r.Form.Get("input_desc")

	// saved the description text, and if that went fine, redirect
	// to the description page.
	if err := h.Persistence.SaveUploadText(upload_id, text); err != nil {
		log.Printf("saving description failed: %s", err.Error())
		rw.WriteHeader(http.StatusOK)
		rw.Write(ErrorPage("saving description failed: " + err.Error()))
	} else {
		rw.Header()["Location"] = []string{"/show/" + upload_id}
		rw.WriteHeader(http.StatusFound)
	}

}
