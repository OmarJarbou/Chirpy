package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", errors.New("Error while signing jwt string: " + err.Error())
	}

	return jwtString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.UUID{}, errors.New("Error while parsing token string: " + err.Error())
	}

	expiry, err2 := token.Claims.GetExpirationTime()
	if err2 != nil {
		return uuid.UUID{}, errors.New("Error while fetching expiration time claim from the token: " + err2.Error())
	}

	if expiry != nil && expiry.Time.Before(time.Now()) {
		return uuid.UUID{}, fmt.Errorf("token expired at: %v", expiry.Time)
	}

	userid_string, err2 := token.Claims.GetSubject()
	if err2 != nil {
		return uuid.UUID{}, errors.New("Error while fetching subject claim from the token: " + err2.Error())
	}

	user_id, err3 := uuid.Parse(userid_string)
	if err3 != nil {
		return uuid.UUID{}, errors.New("Error while parsing userid string to uuid: " + err3.Error())
	}

	return user_id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	bearer_token := headers.Get("Authorization")
	if bearer_token == "" {
		return "", errors.New("Authorization(Bearer Token) header not found")
	}

	if strings.HasPrefix(bearer_token, "Bearer ") {
		token := strings.TrimSpace(bearer_token[len("Bearer "):])
		return token, nil
	} else {
		return "", errors.New("Header does not start with 'Bearer '")
	}
}
