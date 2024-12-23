package atlassian

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"markcli/internal/logging"
	"net/http"
)

// Client represents an Atlassian API client
type Client struct {
	baseURL    string
	email      string
	token      string
	httpClient *http.Client
}

// NewClient creates a new Atlassian API client
func NewClient(baseURL, email, token string) *Client {
	return &Client{
		baseURL:    baseURL,
		email:      email,
		token:      token,
		httpClient: &http.Client{},
	}
}

// newRequest creates a new HTTP request with authentication and common headers
func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	// Build full URL
	reqURL := fmt.Sprintf("%s%s", c.baseURL, path)

	// Create request body if provided
	var buf io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		buf = bytes.NewBuffer(jsonBody)
		logging.LogDebug("Request Body: %s", string(jsonBody))
	}

	// Create request
	req, err := http.NewRequest(method, reqURL, buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set common headers
	req.SetBasicAuth(c.email, c.token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}
