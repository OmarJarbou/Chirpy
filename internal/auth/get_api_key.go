package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	api_key := headers.Get("Authorization")
	if api_key == "" {
		return "", errors.New("Authorization(Api Key) header not found")
	}

	if strings.HasPrefix(api_key, "ApiKey ") {
		apikey := strings.TrimSpace(api_key[len("ApiKey "):])
		return apikey, nil
	} else {
		return "", errors.New("Header does not start with 'ApiKey '")
	}
}
