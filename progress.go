package main

import (
	"fmt"
	"io"
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
	n, err = prc.r.Read(p)
	if err == nil {
		prc.received += uint64(n)
		percent := int((100 * prc.received) / prc.content_length)
		if prc.prev_percent < 0 || percent > prc.prev_percent {
			fmt.Printf("upload progress for %s: %d (%d of %d)\n", prc.upload_id, percent, prc.received, prc.content_length)
			if err = WriteUploadProgress(prc.upload_id, percent); err != nil {
				fmt.Printf("error writing progress: %s\n", err.Error())
			}
		}
		prc.prev_percent = percent
	}
	return
}

func (prc *ProgressReadCloser) Close() error {
	return prc.r.Close()
}

