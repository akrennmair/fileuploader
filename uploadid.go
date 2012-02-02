package main

import (
	"crypto/rand"
	"crypto/md5"
	"bytes"
	"os"
	"fmt"
)

var sharedSecret = "jf7$SD!&/kI9IBjk<Lz8FV"

func generateUploadID() (string, os.Error) {
	rnd := make([]byte, 8)
	if n, err := rand.Read(rnd); err != nil || n != len(rnd) {
		return "", err
	}

	hash := md5.New()
	hash.Write([]byte(sharedSecret))
	hash.Write(rnd)

	result := append(hash.Sum(), rnd...)

	return fmt.Sprintf("%x", result), nil
}

func verifyUploadID(id string) bool {
	upid := []byte{}
	for i:=0;(i+1)<len(id);i+=2 {
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

	return bytes.Equal(hash.Sum(), md5check)
}
