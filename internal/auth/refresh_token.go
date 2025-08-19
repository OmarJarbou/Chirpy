package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

func MakeRefreshToken() (string, error) {
	random_32_byte := make([]byte, 32)
	n, err := rand.Read(random_32_byte)
	if err != nil {
		return "", errors.New("Error while generating random 32 byte using rand.Read: " + err.Error())
	}
	if n != 32 {
		return "", errors.New("the generating random slice of bytes using rand.Read isn't exactly 32 byte")
	}

	refresh_token_hex_string := hex.EncodeToString(random_32_byte)
	return refresh_token_hex_string, nil
}
