package main

import (
	"bytes"
	"net/http"
	"io"
)

// mockup for http.ResponseWriter

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


// mockup for io.ReaderCloser resp. io.WriterCloser

type BufferReaderCloser struct {
	b *bytes.Buffer
}

func (f *BufferReaderCloser) Close() error {
	return nil
}

func NewBufferReaderCloser(data string) *BufferReaderCloser {
	return &BufferReaderCloser{b: bytes.NewBufferString(data)}
}

func (f *BufferReaderCloser) Read(b []byte) (int, error) {
	return f.b.Read(b)
}

func (f *BufferReaderCloser) Write(b []byte) (int, error) {
	return f.b.Write(b)
}

func (f *BufferReaderCloser) Len() int {
	return f.b.Len()
}

type MockPersistenceManager struct {
	Progress    int
	FileContent *BufferReaderCloser
	Exists      bool
	description string
}

func NewMockPersistenceManager() *MockPersistenceManager {
	return &MockPersistenceManager{FileContent: &BufferReaderCloser{}}
}

func (h *MockPersistenceManager) WriteUploadProgress(upload_id string, percent int) error {
	h.Progress = percent
	return nil
}

func (h *MockPersistenceManager) GetUploadProgress(upload_id string) (percent int, err error) {
	return h.Progress, nil
}

func (h *MockPersistenceManager) OpenUpload(upload_id string) (io.ReadCloser, error) {
	return h.FileContent, nil
}

func (h *MockPersistenceManager) OpenUploadWritable(upload_id string) (io.WriteCloser, error) {
	return h.FileContent, nil
}

func (h *MockPersistenceManager) UploadExists(upload_id string) bool {
	return h.Exists
}

func (h *MockPersistenceManager) SaveUploadText(upload_id, text string) error {
	h.description = text
	return nil
}

func (h *MockPersistenceManager) GetUploadText(upload_id string) (string, error) { 
	return h.description, nil
}

func (h *MockPersistenceManager) SaveUploadFilename(upload_id, filename string) error { 
	return nil
}

func (h *MockPersistenceManager) GetUploadFilename(upload_id string) (string, error) { 
	return "", nil
}

func (h *MockPersistenceManager) GetUploadSize(upload_id string) (size int64, err error) { 
	return int64(h.FileContent.Len()), nil
}
