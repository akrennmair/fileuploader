package main

import (
	"testing"
	"net/http"
)

func TestRequestIDHandler(t *testing.T) {
	rh := &RequestIDHandler{}

	req, _ := http.NewRequest("GET", "http://localhost:80/requpid", nil)

	resp := NewMockResponseWriter()

	rh.ServeHTTP(resp, req)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("RequestIDHandler doesn't return HTTP status code 200 (%d instead)", resp.StatusCode)
	}
	if !VerifyUploadID(resp.Buffer.String()) {
		t.Errorf("RequestIDHandler doesn't return a valid upload ID")
	}
}
