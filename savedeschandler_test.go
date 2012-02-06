package main

import (
	"testing"
	"net/http"
	"bytes"
)

func TestSaveDescHandler(t *testing.T) {
	pm := NewMockPersistenceManager()

	sdh := &SaveDescHandler{Persistence: pm}

	upload_id, _ := GenerateUploadID()

	form_str := "input_desc=Hello%20World"

	req, _ := http.NewRequest("POST", "http://localhost/savedesc/" + upload_id, bytes.NewBufferString(form_str))
	req.ContentLength = int64(len(form_str))
	req.Header["Content-Type"] = []string{"application/x-www-form-urlencoded"}

	resp := NewMockResponseWriter()

	pm.Progress = 100
	pm.Exists = true

	sdh.ServeHTTP(resp, req)

	if (resp.StatusCode != http.StatusFound) {
		t.Errorf("SaveDescHandler didn't answer with 302 (%d instead)", resp.StatusCode)
	}
	if len(resp.Header()["Location"]) < 1 || resp.Header()["Location"][0] != "/show/" + upload_id {
		t.Errorf("SaveDescHandler didn't set proper Location header")
	}

	text, _ := pm.GetUploadText(upload_id)
	if text != "Hello World" {
		t.Errorf("SaveDescHandler didn't properly save the uploaded text to the persistence manager: %s", text)
	}
}
