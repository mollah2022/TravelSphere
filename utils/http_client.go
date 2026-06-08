package utils

import (
	"net/http"
	"time"
)

// NewHTTPClient creates a new HTTP client with a custom timeout
// This helps prevent requests from hanging indefinitely
func NewHTTPClient(timeoutSeconds int) *http.Client {
	return &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}
}
