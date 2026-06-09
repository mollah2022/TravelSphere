package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient interface — এটা দিয়ে real ও mock client আলাদা করা যায়
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// DefaultHTTPClient production এ ব্যবহার হয়
type DefaultHTTPClient struct {
	client *http.Client
}

// NewHTTPClient নতুন HTTP client তৈরি করে timeout সহ
func NewHTTPClient(timeoutSeconds int) *DefaultHTTPClient {
	return &DefaultHTTPClient{
		client: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
	}
}

// Get HTTP GET request করে
func (c *DefaultHTTPClient) Get(url string) (*http.Response, error) {
	return c.client.Get(url)
}

// FetchJSON যেকোনো URL থেকে JSON fetch করে struct এ decode করে
// target must be a pointer to a struct
func FetchJSON(client HTTPClient, url string, target interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d for URL: %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return nil
}
