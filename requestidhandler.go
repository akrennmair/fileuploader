package main

import (
	"net/http"
)

type RequestIDHandler struct { }

func (h *RequestIDHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if upload_id, err := GenerateUploadID(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	} else {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(upload_id))
	}
}
