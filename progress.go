package main

import (
	"io"
	"log"
)

type ProgressReadCloser struct {
	r              io.ReadCloser
	content_length uint64
	received       uint64
	upload_id      string
	prev_percent   int
}

func NewProgressReadCloser(r io.ReadCloser, content_length uint64, upload_id string) io.ReadCloser {
	return &ProgressReadCloser{r: r, content_length: content_length, upload_id: upload_id, prev_percent: -1}
}

func (prc *ProgressReadCloser) Read(p []byte) (n int, err error) {
	// first, call the Read() method on our wrapped io.ReadCloser
	n, err = prc.r.Read(p)
	if err == nil {
		// if the read operation went fine, we add the number of received bytes to the count
		// of what we received so far. Then we compute how many percent of the uploaded have
		// already been completed, and if it's bigger than the previous value, we persist this
		// new upload progress value. The check whether the new percentage is bigger than the
		// older one is to limit the I/O load that is otherwise created by consistently updating
		// the pgoress.
		prc.received += uint64(n)
		percent := int((100 * prc.received) / prc.content_length)
		if prc.prev_percent < 0 || percent > prc.prev_percent {
			if err = WriteUploadProgress(prc.upload_id, percent); err != nil {
				log.Printf("error writing progress: %s", err.Error())
			}
		}
		prc.prev_percent = percent
	}
	return
}

func (prc *ProgressReadCloser) Close() error {
	return prc.r.Close()
}
