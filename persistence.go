package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// the functions in this file implement the "persistence layer" of the uploaded files and the
// corresponding metadata, and implement all the details how this data is stored and retrieved.

// this function stores the current upload progress for an upload ID.
func WriteUploadProgress(upload_id string, percent int) error {
	// the Unix way of atomically overwriting a file (i.e. so that race conditions of possibly
	// parallel reading and writing to the same files are avoided) is to write the new data to
	// a temporary file and then to rename the temporary file to the file that needs to be overwritten.
	if f, err := os.OpenFile("files/"+upload_id+".prog.tmp", os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644); err == nil {
		fmt.Fprintf(f, "%d", percent)
		f.Close()
		return os.Rename("files/"+upload_id+".prog.tmp", "files/"+upload_id+".prog")
	} else {
		return err
	}
	return nil
}

// this function fetches the current upload progress for an upload ID.
func GetUploadProgress(upload_id string) (percent int, err error) {
	if f, err := os.Open("files/" + upload_id + ".prog"); err == nil {
		_, err = fmt.Fscanf(f, "%d", &percent)
	}
	return
}

// this function opens an existing uploaded file and returns an io.ReadCloser through
// which the uploaded file can be read.
func OpenUpload(upload_id string) (io.ReadCloser, error) {
	return os.Open("files/" + upload_id)
}

// this function opens a non-existing uploaded file and returns an io.WriteCloser through
// which the file can be filled. If a file with this upload ID has already been written,
// the function returns an error.
func OpenUploadWritable(upload_id string) (io.WriteCloser, error) {
	return os.OpenFile("files/"+upload_id, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
}

// this function returns whether a file with that upload ID already exists.
func UploadExists(upload_id string) bool {
	_, err := os.Stat("files/" + upload_id)
	return err == nil
}

// this function saves the description text for an upload ID. If the description text
// has already been saved, the function returns an error.
func SaveUploadText(upload_id, text string) error {
	if f, err := os.OpenFile("files/"+upload_id+".desc", os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644); err != nil {
		return err
	} else {
		f.WriteString(text)
		f.Close()
	}
	return nil
}

// this function fetches the saved descripton for an upload ID.
func GetUploadText(upload_id string) (string, error) {
	return slurp("files/" + upload_id + ".desc")
}

// this function saves the original filename for an upload ID. If the original filename
// has already been saved, the function returns an error.
func SaveUploadFilename(upload_id, filename string) error {
	if f, err := os.OpenFile("files/"+upload_id+".fn", os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644); err != nil {
		return err
	} else {
		f.WriteString(filename)
		f.Close()
	}
	return nil
}

// this function fetches the saved original filename for an upload ID.
func GetUploadFilename(upload_id string) (string, error) {
	return slurp("files/" + upload_id + ".fn")
}

// helper function that attemps to open a file and return its complete
// content as string.
func slurp(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	io.Copy(buf, f)
	f.Close()
	return buf.String(), nil
}

// this function returns the size of a previously uploaded file, identified
// by its upload ID.
func GetUploadSize(upload_id string) (size int64, err error) {
	if fi, err := os.Stat("files/" + upload_id); err != nil {
		return -1, err
	} else {
		size = fi.Size()
	}
	return
}
