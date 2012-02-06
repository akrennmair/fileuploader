package main

import (
	"testing"
	"net/http"
	"strings"
)

func TestShowHandler(t *testing.T) {
	pm := NewMockPersistenceManager()

	sh := &ShowHandler{Persistence: pm}

	upload_id, _ := GenerateUploadID()

	req, _ := http.NewRequest("GET", "http://localhost/show/" + upload_id, nil)

	resp := NewMockResponseWriter()

	pm.Progress = 100
	pm.SaveUploadText(upload_id, "Test Description")
	pm.SaveUploadFilename(upload_id, "testfilename.zip")

	sh.ServeHTTP(resp, req)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("ShowHandler didn't answer with 200 (%d instead)", resp.StatusCode)
	}

	if !strings.Contains(resp.Buffer.String(), "Test Description") {
		t.Error("ShowHandler doesn't display description text")
	}
	if !strings.Contains(resp.Buffer.String(), "testfilename.zip") {
		t.Error("ShowHandler doesn't display original filename")
	}
}
