package main

import (
	"testing"
	"errors"
	"bytes"
	"net/http"
	"fmt"
	"strings"
	"io"
)


func TestProgressHandler(t *testing.T) {
	pm := &MockPersistenceManager{}

	ph := &ProgressHandler{Persistence: pm}

	upload_id, _ := GenerateUploadID()

	req, _ := http.NewRequest("GET", "http://localhost:80/progress/" + upload_id, nil)

	resp := NewMockResponseWriter()

	pm.Progress = 23

	ph.ServeHTTP(resp, req)

	if resp.Buffer.String() != "23" {
		t.Errorf("for progress 23, ProgressHandler didn't deliver string '10'")
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

type MockResponseWriter struct {
	StatusCode int
	Buffer      *bytes.Buffer
	header      http.Header
}

func NewMockResponseWriter() *MockResponseWriter {
	return &MockResponseWriter { Buffer: new(bytes.Buffer), header: make(http.Header) };
}

func (w *MockResponseWriter) Header() http.Header {
	return w.header
}

func (w *MockResponseWriter) WriteHeader(c int) {
	w.StatusCode = c
}

func (w *MockResponseWriter) Write(b []byte) (int, error) {
	if w.StatusCode == 0 {
		w.StatusCode = http.StatusOK
	}
	return w.Buffer.Write(b)
}

type MockPersistenceManager struct {
	Progress int
}

func (h *MockPersistenceManager) WriteUploadProgress(upload_id string, percent int) error {
	h.Progress = percent
	return nil
}

func (h *MockPersistenceManager) GetUploadProgress(upload_id string) (percent int, err error) {
	return h.Progress, nil
}

func (h *MockPersistenceManager) OpenUpload(upload_id string) (io.ReadCloser, error) {
	return nil, errors.New("mock persistence manager")
}

func (h *MockPersistenceManager) OpenUploadWritable(upload_id string) (io.WriteCloser, error) {
	return nil, errors.New("mock persistence manager")
}

func (h *MockPersistenceManager) UploadExists(upload_id string) bool {
	return true
}

func (h *MockPersistenceManager) SaveUploadText(upload_id, text string) error {
	return nil
}

func (h *MockPersistenceManager) GetUploadText(upload_id string) (string, error) { 
	return "", nil
}

func (h *MockPersistenceManager) SaveUploadFilename(upload_id, filename string) error { 
	return nil
}

func (h *MockPersistenceManager) GetUploadFilename(upload_id string) (string, error) { 
	return "", nil
}

func (h *MockPersistenceManager) GetUploadSize(upload_id string) (size int64, err error) { 
	return 23, nil
}
