package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}

	if hash == "" {
		t.Error("Expected hash to not be empty")
	}

	if hash == "password" {
		t.Error("Expected hash to be different from password")
	}
}

func TestComparePassword(t *testing.T) {
	hash, err := HashPassword("password")

	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}

	if !ComparePasswords(hash, []byte("password")) {
		t.Errorf("Expected password to match hash")
	}

	if ComparePasswords(hash, []byte("notpassword")) {
		t.Errorf("Expected password to not match hash")
	}
}
