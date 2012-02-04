package main

import (
	"os"
	"fmt"
	"io"
)

func WriteUploadProgress(upload_id string, percent int) error {
	if f, err := os.OpenFile("files/" + upload_id + ".prog.tmp", os.O_CREATE | os.O_WRONLY, 0644); err == nil {
		fmt.Fprintf(f, "%d", percent)
		f.Close()
		return os.Rename("files/" + upload_id + ".prog.tmp", "files/" + upload_id + ".prog")
	} else {
		return err
	}
	return nil
}

func ReadUploadProgress(upload_id string) (percent int, err error) {
	if f, err := os.Open("files/" + upload_id + ".prog"); err == nil {
		_, err = fmt.Fscanf(f, "%d", &percent)
	}
	return
}

func OpenUpload(upload_id string) (io.ReadCloser, error) {
	return os.Open("files/" + upload_id)
}

func OpenUploadWritable(upload_id string) (io.WriteCloser, error) {
	return os.OpenFile("files/" + upload_id, os.O_CREATE | os.O_WRONLY, 0644)
}
