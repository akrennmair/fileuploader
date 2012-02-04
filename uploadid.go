package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
)

var secretKey = "jf7$SD!&/kI9IBjk<Lz8FV"

// this function generates a new upload ID
func GenerateUploadID() (string, error) {
	// first, we read 8 bytes of random data
	rnd := make([]byte, 8)
	if n, err := rand.Read(rnd); err != nil || n != len(rnd) {
		return "", err
	}

	// then we hash the random data with the secret key.
	hash := md5.New()
	hash.Write([]byte(secretKey))
	hash.Write(rnd)

	// and then we append MD5 hash and random data and convert it 
	// to lower-case hexadecimal representation.
	result := append(hash.Sum(nil), rnd...)
	return fmt.Sprintf("%x", result), nil
}

// this function verifies the validity of an upload ID
func VerifyUploadID(id string) bool {
	// first, we parse the upload ID byte by byte from hexadecimal
	// representation to the actual bytes.
	upid := []byte{}
	for i := 0; (i + 1) < len(id); i += 2 {
		var b byte
		if n, err := fmt.Sscanf(string(id[i:i+2]), "%x", &b); err != nil || n != 1 {
			return false
		}
		upid = append(upid, b)
	}

	if len(upid) != 24 {
		return false
	}

	md5check := upid[0:16]
	rnd := upid[16:]

	// then we hash the secret key and the random data that we extracted from
	// the upload ID.
	hash := md5.New()
	hash.Write([]byte(secretKey))
	hash.Write(rnd)

	// if the resulting hash is equal to the MD5 checksum from the upload ID
	// then the provided upload ID is valid.
	return bytes.Equal(hash.Sum(nil), md5check)
}

// this helper function parses an upload ID from a URI and verifies it. It
// assumes the all data after the last / of the URI is an upload ID. If
// no upload ID could be found or the upload ID is invalid, an error is returned.
func GetUploadID(uri string) (upload_id string, err error) {
	slash_pos := strings.LastIndex(uri, "/")
	if slash_pos < 0 {
		err = errors.New("no Upload ID")
		return
	}
	upload_id = uri[slash_pos+1:]
	if !VerifyUploadID(upload_id) {
		err = errors.New("invalid Upload ID")
	}
	return
}
