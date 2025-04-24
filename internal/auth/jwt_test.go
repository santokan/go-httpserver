package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"
	validDuration := time.Hour
	expiredDuration := -time.Hour // A duration in the past

	// Test Case 1: Successful creation and validation
	tokenString, err := MakeJWT(userID, secret, validDuration)
	if err != nil {
		t.Fatalf("Test Case 1 Failed: MakeJWT returned an error: %v", err)
	}
	if tokenString == "" {
		t.Fatal("Test Case 1 Failed: MakeJWT returned an empty token string")
	}

	validatedUserID, err := ValidateJWT(tokenString, secret)
	if err != nil {
		t.Errorf("Test Case 1 Failed: ValidateJWT returned an error for a valid token: %v", err)
	}
	if validatedUserID != userID {
		t.Errorf("Test Case 1 Failed: ValidateJWT returned incorrect userID. Got %v, want %v", validatedUserID, userID)
	}

	// Test Case 2: Validation with incorrect secret
	_, err = ValidateJWT(tokenString, "wrong-secret")
	if err == nil {
		t.Error("Test Case 2 Failed: ValidateJWT did not return an error for an incorrect secret")
	}

	// Test Case 3: Validation with invalid token string
	invalidTokenString := "this.is.not.a.valid.jwt"
	_, err = ValidateJWT(invalidTokenString, secret)
	if err == nil {
		t.Error("Test Case 3 Failed: ValidateJWT did not return an error for an invalid token string")
	}

	// Test Case 4: Validation with expired token
	expiredTokenString, err := MakeJWT(userID, secret, expiredDuration)
	if err != nil {
		t.Fatalf("Test Case 4 Failed: MakeJWT returned an error creating expired token: %v", err)
	}

	_, err = ValidateJWT(expiredTokenString, secret)
	if err == nil {
		t.Error("Test Case 4 Failed: ValidateJWT did not return an error for an expired token")
	} else {
		// Check if the error is specifically about token expiration (or invalidity)
		// The jwt library might return different specific errors, but it should be non-nil
		t.Logf("Test Case 4 Passed: Received expected error for expired token: %v", err)
	}

	// Test Case 5: Validation with token missing parts
	malformedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ" // Missing signature
	_, err = ValidateJWT(malformedToken, secret)
	if err == nil {
		t.Error("Test Case 5 Failed: ValidateJWT did not return an error for a malformed token (missing signature)")
	}
}
