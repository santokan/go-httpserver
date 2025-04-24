package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "password123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an unexpected error: %v", err)
	}

	if hash == "" {
		t.Fatal("HashPassword returned an empty hash")
	}

	// Verify the hash matches the original password
	err = CheckPasswordHash(hash, password)
	if err != nil {
		t.Errorf("CheckPasswordHash failed to verify the generated hash: %v", err)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "correcthorsebatterystaple"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Setup failed: Could not hash password: %v", err)
	}

	// Test case 1: Correct password
	err = CheckPasswordHash(hash, password)
	if err != nil {
		t.Errorf("Test Case 1 Failed: CheckPasswordHash returned error for correct password: %v", err)
	}

	// Test case 2: Incorrect password
	incorrectPassword := "wrongpassword"
	err = CheckPasswordHash(hash, incorrectPassword)
	if err == nil {
		t.Error("Test Case 2 Failed: CheckPasswordHash did not return error for incorrect password")
	}

	// Test case 3: Invalid hash format (bcrypt hashes start with $)
	invalidHash := "notavalidbcryptHash"
	err = CheckPasswordHash(invalidHash, password)
	if err == nil {
		t.Error("Test Case 3 Failed: CheckPasswordHash did not return error for invalid hash format")
	}
}
