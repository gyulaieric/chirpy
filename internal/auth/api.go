package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authorizationHeader := headers.Get("Authorization")
	if authorizationHeader == "" {
		return "", fmt.Errorf(`invalid "Authorization" Header`)
	}
	return strings.Replace(authorizationHeader, "ApiKey ", "", 1), nil
}
