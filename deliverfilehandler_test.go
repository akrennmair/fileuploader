package main

import (
	"testing"
	"net/http"
)

func TestDeliverFileHandler(t *testing.T) {
	pm := NewMockPersistenceManager()

	dfh := &DeliverFileHandler{Persistence: pm}

	upload_id, _ := GenerateUploadID()
	req, _ := http.NewRequest("GET", "http://localhost:80/progress/" + upload_id, nil)

	resp := NewMockResponseWriter()

	pm.FileContent = NewBufferReaderCloser("test data")
	pm.Progress = 42

	dfh.ServeHTTP(resp, req)

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("DeliverHandler didn't deliver correct 403 status code for incomplete upload")
	}

	pm.FileContent = NewBufferReaderCloser("test data")
	pm.Progress = 100
	resp.Buffer.Reset()

	dfh.ServeHTTP(resp, req)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("DeliverFileHandler didn't set correct HTTP 200 status code (%d instead)", resp.StatusCode)
	}
	if resp.Buffer.String() != "test data" {
		t.Errorf("DeliverFileHandler didn't deliver 'test data' (%s)", resp.Buffer.String())
	}
	if len(resp.Header()["Content-Type"]) < 1 || resp.Header()["Content-Type"][0] != "application/octet-stream" {
		t.Errorf("DeliverFileHandler didn't set correct Content-Type")
	}
	if len(resp.Header()["Content-Length"]) < 1 || resp.Header()["Content-Length"][0] != "9" {
		t.Errorf("DeliverFileHandler didn't set correct Content-Length")
	}
}
