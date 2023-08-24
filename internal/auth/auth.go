package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Example
// Authorization: ApiKey { value here }
func GetApiKey(headers http.Header) (string, error) {
	header := headers.Get("Authorization")

	if header == "" {
		return "", errors.New("no auth info")
	}

	apiKeyValues := strings.Split(header, " ")
	if len(apiKeyValues) != 2 || apiKeyValues[0] != "ApiKey" {
		return "", errors.New("malformed auth info format")
	}

	return apiKeyValues[1], nil
}
