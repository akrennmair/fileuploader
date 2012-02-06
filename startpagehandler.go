package main

import (
	"net/http"
)

type StartPageHandler struct {
	Other     http.Handler
}

// The / prefix is special, so I wrote a handler to handles /, and hands everything else to
// a http.Fileserver handler that delivers files from a directory. That means that if we want
// to make static files available, we can simply place them there.

func (h *StartPageHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		rw.WriteHeader(http.StatusOK)
		rw.Write(UploadPage())
	} else {
		h.Other.ServeHTTP(rw, r)
	}
}

