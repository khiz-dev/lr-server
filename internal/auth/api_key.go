package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	APIKeyPrefix = "ApiKey "
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKeyWithPrefix := headers.Get(AuthHeader)
	fmt.Println("API key with prefix:", apiKeyWithPrefix)
	apiKeyWithPrefix = strings.TrimSpace(apiKeyWithPrefix)
	if apiKeyWithPrefix == "" {
		return "", errors.New("API key is required")
	}
	if !strings.HasPrefix(apiKeyWithPrefix, APIKeyPrefix) {
		return "", errors.New("Invalid API key format")
	}
	apiKey := apiKeyWithPrefix[len(APIKeyPrefix):]
	fmt.Println("Extracted API key:", apiKey)
	if apiKey == "" {
		return "", errors.New("API key is required")
	}
	return strings.TrimSpace(apiKey), nil
}