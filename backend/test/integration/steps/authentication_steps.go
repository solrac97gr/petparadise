package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cucumber/godog"
)

// RegisterAuthenticationSteps registers step definitions for authentication scenarios
func RegisterAuthenticationSteps(ctx *godog.ScenarioContext, client *APIClient) {
	// Given steps
	ctx.Step(`^I am not authenticated$`, client.iAmNotAuthenticated)
	ctx.Step(`^I have valid user credentials$`, client.iHaveValidUserCredentials)
	ctx.Step(`^I have invalid user credentials$`, client.iHaveInvalidUserCredentials)
	ctx.Step(`^I have a valid refresh token$`, client.iHaveValidRefreshToken)
	ctx.Step(`^I have an expired refresh token$`, client.iHaveExpiredRefreshToken)
	ctx.Step(`^I have a valid access token$`, client.iHaveValidAccessToken)
	ctx.Step(`^I have an expired access token$`, client.iHaveExpiredAccessToken)

	// When steps
	ctx.Step(`^I login with my credentials$`, client.iLoginWithMyCredentials)
	ctx.Step(`^I request to refresh my tokens$`, client.iRequestToRefreshMyTokens)
	ctx.Step(`^I logout$`, client.iLogout)
	ctx.Step(`^I use my token to access a protected resource$`, client.iUseMyTokenToAccessProtectedResource)
	ctx.Step(`^I try to access a protected resource without authentication$`, client.iTryToAccessProtectedResourceWithoutAuth)
	ctx.Step(`^I revoke all my user tokens$`, client.iRevokeAllMyUserTokens)

	// Then steps
	ctx.Step(`^I should receive a valid token pair$`, client.iShouldReceiveValidTokenPair)
	ctx.Step(`^I should receive new valid tokens$`, client.iShouldReceiveNewValidTokens)
	ctx.Step(`^I should receive an authentication error$`, client.iShouldReceiveAuthenticationError)
	ctx.Step(`^I should receive a (\d+) status code$`, client.iShouldReceiveStatusCode)
	ctx.Step(`^the response should contain "([^"]*)"$`, client.theResponseShouldContain)
	ctx.Step(`^my tokens should be invalidated$`, client.myTokensShouldBeInvalidated)
}

// Step definition implementations

// Given steps
func (c *APIClient) iAmNotAuthenticated() error {
	c.AuthToken = ""
	return nil
}

func (c *APIClient) iHaveValidUserCredentials() error {
	// Set up valid test credentials
	// These would typically be test user credentials in the test database
	return nil
}

func (c *APIClient) iHaveInvalidUserCredentials() error {
	// Set up invalid test credentials
	return nil
}

func (c *APIClient) iHaveValidRefreshToken() error {
	// Set up valid refresh token for testing
	// This could be fetched from the test database or generated on the fly
	return nil
}

func (c *APIClient) iHaveExpiredRefreshToken() error {
	// Set up expired refresh token for testing
	return nil
}

func (c *APIClient) iHaveValidAccessToken() error {
	// Set up valid access token for testing
	return nil
}

func (c *APIClient) iHaveExpiredAccessToken() error {
	// Set up expired access token for testing
	return nil
}

// When steps
func (c *APIClient) iLoginWithMyCredentials() error {
	// Prepare login request data
	loginData := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}

	// Send login request
	err := c.Post("/users/login", loginData)
	if err != nil {
		return err
	}

	return nil
}

func (c *APIClient) iRequestToRefreshMyTokens() error {
	// Extract refresh token from previous response or use the one set directly
	refreshToken := "test-refresh-token"
	if value, exists := c.GetValueFromResponse("tokens.refresh_token"); exists {
		refreshToken = value.(string)
	}

	// Prepare refresh request
	refreshData := map[string]string{
		"refresh_token": refreshToken,
	}

	// Send refresh request
	err := c.Post("/users/refresh", refreshData)
	if err != nil {
		return err
	}

	return nil
}

func (c *APIClient) iLogout() error {
	// Send logout request
	return c.Post("/users/logout", nil)
}

func (c *APIClient) iUseMyTokenToAccessProtectedResource() error {
	// Use the token to access a protected endpoint (e.g., user profile)
	return c.Get("/users/profile")
}

func (c *APIClient) iTryToAccessProtectedResourceWithoutAuth() error {
	// Clear any authentication token
	c.AuthToken = ""
	
	// Try to access a protected endpoint
	return c.Get("/users/profile")
}

func (c *APIClient) iRevokeAllMyUserTokens() error {
	// Send request to revoke all tokens for the current user
	userId := "test-user-id" // This would come from test context or previous responses
	return c.Post(fmt.Sprintf("/users/%s/revoke-tokens", userId), nil)
}

// Then steps
func (c *APIClient) iShouldReceiveValidTokenPair() error {
	// Check if response contains access_token and refresh_token
	respBody := c.GetResponseBodyAsMap()
	if respBody == nil {
		return fmt.Errorf("response body is empty or not valid JSON")
	}

	// Check for tokens object
	tokensObj, ok := respBody["tokens"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("tokens object not found in response")
	}

	// Check for access_token
	if _, exists := tokensObj["access_token"]; !exists {
		return fmt.Errorf("access_token not found in response")
	}

	// Check for refresh_token
	if _, exists := tokensObj["refresh_token"]; !exists {
		return fmt.Errorf("refresh_token not found in response")
	}

	// Check for expires_in
	if _, exists := tokensObj["expires_in"]; !exists {
		return fmt.Errorf("expires_in not found in response")
	}

	return nil
}

func (c *APIClient) iShouldReceiveNewValidTokens() error {
	// Similar to iShouldReceiveValidTokenPair, but could add additional checks
	// for token freshness if needed
	return c.iShouldReceiveValidTokenPair()
}

func (c *APIClient) iShouldReceiveAuthenticationError() error {
	// Check for 401 Unauthorized status code
	if c.GetResponseStatusCode() != http.StatusUnauthorized {
		return fmt.Errorf("expected status code 401, got %d", c.GetResponseStatusCode())
	}

	// Check for error message in response
	respBody := c.GetResponseBodyAsMap()
	if respBody == nil {
		return fmt.Errorf("response body is empty or not valid JSON")
	}

	if _, exists := respBody["error"]; !exists {
		return fmt.Errorf("error message not found in response")
	}

	return nil
}

func (c *APIClient) iShouldReceiveStatusCode(code int) error {
	if c.GetResponseStatusCode() != code {
		return fmt.Errorf("expected status code %d, got %d", code, c.GetResponseStatusCode())
	}
	return nil
}

func (c *APIClient) theResponseShouldContain(text string) error {
	body := c.GetResponseBody()
	if body == nil {
		return fmt.Errorf("response body is empty")
	}

	if !bytes.Contains(body, []byte(text)) {
		return fmt.Errorf("response does not contain %q", text)
	}

	return nil
}

func (c *APIClient) myTokensShouldBeInvalidated() error {
	// Try to use the token to access a protected resource
	// This should now fail with 401 Unauthorized
	err := c.Get("/users/profile")
	if err != nil {
		return err
	}

	if c.GetResponseStatusCode() != http.StatusUnauthorized {
		return fmt.Errorf("token was not properly invalidated, got status code %d", c.GetResponseStatusCode())
	}

	return nil
}
