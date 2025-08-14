package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	passwords := []string{
		"Omar@123456",
		"Ahmad",
		"",
	}
	for i, password := range passwords {
		hashed, err := HashPassword(password)
		if err != nil {
			t.Errorf("%v", err)
		}
		t.Logf("Hashing %d succeeded: %s", i, hashed)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	passwords := map[string]string{
		"Omar@123456": "$2a$10$g2MD5zJMhX71H5matSzG3OSYpyyBAVYlZMB10miN5STiE.97TINVm",
		"Ahmad":       "$2a$10$bRlYjIKxKW65XrkLqLC9I.mpOB/4Wo0Cr6JqwyXDULIJR7X8GlqbK",
		"":            "$2a$10$KLKGDhyqKoaCcS5dLTdd6.8jb6H0h9kysOIZGoe2ikppzDFb95xYy",
	}
	for key, value := range passwords {
		err := CheckPasswordHash(key, value)
		if err != nil {
			t.Errorf("%v", err)
		}
		t.Logf("Checking Password Hash for %s succeeded", key)
	}
}

type makeJWTTestData struct {
	UserID        uuid.UUID
	TokenSecret   string
	ExpiresIn     time.Duration
	ErrorExpected bool
}

func TestJWTFlow(t *testing.T) {
	tests := []makeJWTTestData{
		{
			UserID:        uuid.MustParse("b3a29e2e-54e4-4b84-a991-07b5f63c2a6a"),
			TokenSecret:   "super-secret-key-123!@#",
			ExpiresIn:     15 * time.Minute,
			ErrorExpected: false,
		},
		{
			UserID:        uuid.MustParse("9f36d4c2-7d3b-4e12-8b65-3f8936e4babc"),
			TokenSecret:   "another-secret-key-456$%^",
			ExpiresIn:     -10 * time.Minute,
			ErrorExpected: true,
		},
	}

	tokens := []string{}

	for i, test := range tests {
		token_string, err := MakeJWT(test.UserID, test.TokenSecret, test.ExpiresIn)
		if err != nil {
			t.Errorf("%v", err)
		}
		tokens = append(tokens, token_string)
		t.Logf("Making token %d succeeded: %s", i, token_string)
	}

	for i, token := range tokens {
		user_id, err := ValidateJWT(token, tests[i].TokenSecret)
		if tests[i].ErrorExpected {
			if err == nil {
				t.Errorf("Expected an error for token %d, but got none", i)
			} else {
				t.Logf("Expected error occurred for token %d: %v", i, err)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for token %d: %v", i, err)
			} else {
				t.Logf("Validating token %d for user %s succeeded", i, user_id.String())
			}
		}
	}

}
