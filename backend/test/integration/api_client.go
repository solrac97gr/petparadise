package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// APIClient wraps HTTP client for API testing
type APIClient struct {
	BaseURL      string
	HTTPClient   *http.Client
	AuthToken    string
	Headers      map[string]string
	LastResponse *APIResponse
}

// APIResponse captures response details for assertions
type APIResponse struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
	ParsedBody map[string]interface{}
}

// NewAPIClient creates a new API client for testing
func NewAPIClient() *APIClient {
	baseURL := os.Getenv("TEST_API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3001/api/users" // Default API URL for testing
	}

	return &APIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
		Headers: make(map[string]string),
	}
}

// SetAuthToken sets the authorization token for requests
func (c *APIClient) SetAuthToken(token string) {
	c.AuthToken = token
}

// AddHeader adds a custom header for requests
func (c *APIClient) AddHeader(key, value string) {
	c.Headers[key] = value
}

// Get performs a GET request to the specified endpoint
func (c *APIClient) Get(endpoint string) error {
	return c.doRequest("GET", endpoint, nil)
}

// Post performs a POST request with JSON body
func (c *APIClient) Post(endpoint string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.doRequest("POST", endpoint, jsonData)
}

// Put performs a PUT request with JSON body
func (c *APIClient) Put(endpoint string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.doRequest("PUT", endpoint, jsonData)
}

// Patch performs a PATCH request with JSON body
func (c *APIClient) Patch(endpoint string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.doRequest("PATCH", endpoint, jsonData)
}

// Delete performs a DELETE request
func (c *APIClient) Delete(endpoint string) error {
	return c.doRequest("DELETE", endpoint, nil)
}

// doRequest handles the HTTP request execution
func (c *APIClient) doRequest(method, endpoint string, body []byte) error {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return err
	}

	// Set content type for requests with body
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set auth token if present
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	// Add custom headers
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Create response object
	c.LastResponse = &APIResponse{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
	}

	// Try to parse JSON response
	if len(respBody) > 0 {
		var result map[string]interface{}
		if err := json.Unmarshal(respBody, &result); err == nil {
			c.LastResponse.ParsedBody = result
		}
	}

	return nil
}

// GetResponseStatusCode returns the status code of the last response
func (c *APIClient) GetResponseStatusCode() int {
	if c.LastResponse != nil {
		return c.LastResponse.StatusCode
	}
	return 0
}

// GetResponseBody returns the body of the last response
func (c *APIClient) GetResponseBody() []byte {
	if c.LastResponse != nil {
		return c.LastResponse.Body
	}
	return nil
}

// GetResponseBodyAsMap returns the parsed JSON body as a map
func (c *APIClient) GetResponseBodyAsMap() map[string]interface{} {
	if c.LastResponse != nil {
		return c.LastResponse.ParsedBody
	}
	return nil
}

// GetValueFromResponse extracts a value from parsed JSON response using dot notation
func (c *APIClient) GetValueFromResponse(path string) (interface{}, bool) {
	if c.LastResponse == nil || c.LastResponse.ParsedBody == nil {
		return nil, false
	}

	// For simple top-level keys
	if value, exists := c.LastResponse.ParsedBody[path]; exists {
		return value, true
	}

	// TODO: Implement support for nested paths with dot notation

	return nil, false
}
