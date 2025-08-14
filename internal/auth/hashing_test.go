package auth

import (
	"testing"
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
