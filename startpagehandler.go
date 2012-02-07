package main

import (
	"net/http"
	"fmt"
)

type StartPageHandler struct {
	Other     http.Handler
}

// The / prefix is special, so I wrote a handler to handles /, and hands everything else to
// a http.Fileserver handler that delivers files from a directory. That means that if we want
// to make static files available, we can simply place them there.

func (h *StartPageHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		page := UploadPage()
		rw.Header()["Content-Length"] = []string{fmt.Sprintf("%d", len(page))}
		rw.WriteHeader(http.StatusOK)
		rw.Write(page)
	} else {
		h.Other.ServeHTTP(rw, r)
	}
}

