package main

import (
	"testing"
)

func TestUploadID(t *testing.T) {
	id, err := generateUploadID()
	if err != nil {
		t.Fatalf("generateUploadID doesn't generate a valid ID: %s", err.String())
	}
	if len(id) != 48 {
		t.Errorf("ID %s is not 48 characters long", id)
	}

	if !verifyUploadID(id) {
		t.Errorf("Previously generated ID %s doesn't verify", id)
	}

	if verifyUploadID("") {
		t.Errorf("empty string verifies")
	}

	if verifyUploadID("x") {
		t.Errorf("x verifies")
	}

	if verifyUploadID("random text") {
		t.Errorf("random text verifies")
	}

	if verifyUploadID("ffffffffffffffffffffffffffffffffffffffffffffffff") {
		t.Errorf("48 characters of f verify")
	}
}
