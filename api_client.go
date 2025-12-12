package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// APIClient handles HTTP requests to the CyberArk API
// In Go, we often create wrapper types to encapsulate functionality
type APIClient struct {
	config     *Config
	httpClient *http.Client
}

// NewAPIClient creates a new API client
// This is a constructor function (Go doesn't have constructors like Python)
// By convention, constructor functions are named New[TypeName]
func NewAPIClient(config *Config) *APIClient {
	// Set default timeout if not specified
	timeout := 30
	if config.Timeout > 0 {
		timeout = config.Timeout
	}

	// Create HTTP client with timeout
	// Unlike Python's requests library, Go's http.Client is built-in
	httpClient := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	return &APIClient{
		config:     config,
		httpClient: httpClient,
	}
}

// Get performs a GET request to the API
// This method demonstrates how to make HTTP calls in Go
func (c *APIClient) Get(endpoint string) ([]byte, error) {
	// Build full URL
	// Go doesn't have string interpolation like f-strings in Python
	url := fmt.Sprintf("%s/%s", c.config.BaseURL, endpoint)

	// Create new HTTP request
	// http.NewRequest returns a request and potentially an error
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	// All CyberArk API calls typically need authentication
	req.Header.Set("Authorization", c.config.APISecret)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	// defer ensures this runs when the function returns
	// Similar to Python's context managers or finally blocks
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check HTTP status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Post performs a POST request to the API
// payload is an interface{} which means it can be any type
// interface{} is similar to Python's Any type
func (c *APIClient) Post(endpoint string, payload interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.config.BaseURL, endpoint)

	// Marshal payload to JSON
	// json.Marshal converts a Go value to JSON bytes
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create request with body
	// bytes.NewBuffer creates an io.Reader from a byte slice
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", c.config.APISecret)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Put performs a PUT request (for updates)
func (c *APIClient) Put(endpoint string, payload interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.config.BaseURL, endpoint)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.config.APISecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Delete performs a DELETE request
func (c *APIClient) Delete(endpoint string) error {
	url := fmt.Sprintf("%s/%s", c.config.BaseURL, endpoint)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.config.APISecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// Example usage in a workflow:
//
// func (w *MyWorkflow) Execute(config *Config, args []string) error {
//     client := NewAPIClient(config)
//
//     // Make a GET request
//     data, err := client.Get("accounts")
//     if err != nil {
//         return err
//     }
//
//     // Parse JSON response
//     var accounts []Account
//     if err := json.Unmarshal(data, &accounts); err != nil {
//         return err
//     }
//
//     // Process accounts...
//     return nil
// }

// Key Differences from Python's requests library:
//
// 1. No automatic JSON encoding - must use json.Marshal/Unmarshal
// 2. Must manually close response bodies (use defer)
// 3. Must check errors explicitly at every step
// 4. Timeouts are set on the client, not per-request
// 5. Headers are set on the request object, not as a dict parameter
