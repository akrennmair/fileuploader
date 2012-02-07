package main

import (
	"testing"
	"net/http"
	"fmt"
	"strings"
)


func TestProgressHandler(t *testing.T) {
	pm := NewMockPersistenceManager()

	ph := &ProgressHandler{Persistence: pm}

	upload_id, _ := GenerateUploadID()

	req, _ := http.NewRequest("GET", "http://localhost:80/progress/" + upload_id, nil)

	resp := NewMockResponseWriter()

	pm.Progress = 23

	ph.ServeHTTP(resp, req)

	if resp.Buffer.String() != "23" {
		t.Errorf("for progress 23, ProgressHandler didn't deliver string '23'")
	}
	if resp.Header()["Content-Length"][0] != "2" {
		t.Errorf("for progress 23, ProgressHandler didn't set the correct Content-Length")
	}
	if resp.Header()["Content-Type"][0] != "text/plain" {
		t.Errorf("ProgressHandler didn't set correct Content-Type: text/plain")
	}

	resp.Buffer.Truncate(0)
	pm.Progress = 100

	ph.ServeHTTP(resp, req)
	if resp.Buffer.String() != "100" {
		t.Errorf("for progress 10, ProgressHandler didn't deliver string '100'")
	}
	if resp.Header()["Content-Length"][0] != fmt.Sprintf("%d", len(resp.Buffer.String())) {
		t.Errorf("ProgressHandler didn't set correct Content-Length")
	}

	resp.Buffer.Truncate(0)
	req, _ = http.NewRequest("GET", "http://localhost:80/progress/INVALID", nil)

	ph.ServeHTTP(resp, req)
	if !strings.Contains(resp.Buffer.String(),"invalid Upload ID") {
		t.Errorf("ProgressHandler didn't detect an invalid upload ID as such")
	}
}

