package main

import (
	"testing"
	"net/http"
	"bytes"
	"fmt"
)

func TestPostUploadHandler(t *testing.T) {
	pm := NewMockPersistenceManager()

	puh := &PostUploadHandler{Persistence: pm}

	upload_id, _ := GenerateUploadID()

	req, _ := http.NewRequest("POST", "http://localhost/upload/" + upload_id, bytes.NewBufferString(multipart_formdata_str))
	req.Header["Content-Length"] = []string{fmt.Sprintf("%d", len(multipart_formdata_str))}
	req.Header["Content-Type"] = []string{"multipart/form-data; boundary=AaB03x"}

	resp := NewMockResponseWriter()

	puh.ServeHTTP(resp, req)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("PostUploadHandler didn't answer with 200 (%d instead)", resp.StatusCode)
	}
	if resp.Buffer.String() != "" {
		t.Errorf("PostUploadHandler printed out a (error?) message: %s", resp.Buffer.String())
	}
	if pm.FileContent.String() != "example file upload" {
		t.Errorf("PostUploadHandler didn't store the correct file content ('%s' instead)", pm.FileContent.String())
	}
	fn, _ := pm.GetUploadFilename(upload_id)
	if fn != "testfile.txt" {
		t.Errorf("PostUploadHandler didn't store the correct original filename (%s instead)", fn)
	}
}

var multipart_formdata_str = `--AaB03x
Content-Disposition: file; filename="testfile.txt"
Content-Type: text/plain

example file upload
--AaB03x--
`
