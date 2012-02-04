package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"strings"
	"errors"
)

var sharedSecret = "jf7$SD!&/kI9IBjk<Lz8FV"

func GenerateUploadID() (string, error) {
	rnd := make([]byte, 8)
	if n, err := rand.Read(rnd); err != nil || n != len(rnd) {
		return "", err
	}

	hash := md5.New()
	hash.Write([]byte(sharedSecret))
	hash.Write(rnd)

	result := append(hash.Sum(nil), rnd...)

	return fmt.Sprintf("%x", result), nil
}

func VerifyUploadID(id string) bool {
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

	hash := md5.New()
	hash.Write([]byte(sharedSecret))
	hash.Write(rnd)

	return bytes.Equal(hash.Sum(nil), md5check)
}

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
