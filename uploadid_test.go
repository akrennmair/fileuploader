package main

import (
	"testing"
)

func TestUploadID(t *testing.T) {
	id, err := GenerateUploadID()
	if err != nil {
		t.Fatalf("GenerateUploadID doesn't generate a valid ID: %s", err.Error())
	}
	if len(id) != 48 {
		t.Errorf("ID %s is not 48 characters long", id)
	}

	if !VerifyUploadID(id) {
		t.Errorf("Previously generated ID %s doesn't verify", id)
	}

	if VerifyUploadID("") {
		t.Errorf("empty string verifies")
	}

	if VerifyUploadID("x") {
		t.Errorf("x verifies")
	}

	if VerifyUploadID("random text") {
		t.Errorf("random text verifies")
	}

	if VerifyUploadID("ffffffffffffffffffffffffffffffffffffffffffffffff") {
		t.Errorf("48 characters of f verify")
	}
}
