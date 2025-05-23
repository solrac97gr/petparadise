package integration

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/cucumber/godog"
)

// AuthSteps contains authentication test steps
type AuthSteps struct {
	client *APIClient
}

// RegisterAuthenticationSteps registers step definitions for authentication scenarios
func RegisterAuthenticationSteps(ctx *godog.ScenarioContext, client *APIClient) {
	steps := &AuthSteps{client: client}

	// Given steps
	ctx.Step(`^I am not authenticated$`, steps.iAmNotAuthenticated)
	ctx.Step(`^I have valid user credentials$`, steps.iHaveValidUserCredentials)
	ctx.Step(`^I have invalid user credentials$`, steps.iHaveInvalidUserCredentials)
	ctx.Step(`^I have a valid refresh token$`, steps.iHaveValidRefreshToken)
	ctx.Step(`^I have an expired refresh token$`, steps.iHaveExpiredRefreshToken)
	ctx.Step(`^I have a valid access token$`, steps.iHaveValidAccessToken)
	ctx.Step(`^I have an expired access token$`, steps.iHaveExpiredAccessToken)

	// When steps
	ctx.Step(`^I login with my credentials$`, steps.iLoginWithMyCredentials)
	ctx.Step(`^I request to refresh my tokens$`, steps.iRequestToRefreshMyTokens)
	ctx.Step(`^I logout$`, steps.iLogout)
	ctx.Step(`^I use my token to access a protected resource$`, steps.iUseMyTokenToAccessProtectedResource)
	ctx.Step(`^I try to access a protected resource without authentication$`, steps.iTryToAccessProtectedResourceWithoutAuth)
	ctx.Step(`^I revoke all my user tokens$`, steps.iRevokeAllMyUserTokens)

	// Then steps
	ctx.Step(`^I should receive a valid token pair$`, steps.iShouldReceiveValidTokenPair)
	ctx.Step(`^I should receive new valid tokens$`, steps.iShouldReceiveNewValidTokens)
	ctx.Step(`^I should receive an authentication error$`, steps.iShouldReceiveAuthenticationError)
	ctx.Step(`^I should receive a (\d+) status code$`, steps.iShouldReceiveStatusCode)
	ctx.Step(`^the response should contain "([^"]*)"$`, steps.theResponseShouldContain)
	ctx.Step(`^my tokens should be invalidated$`, steps.myTokensShouldBeInvalidated)
}

// Step definition implementations

// Given steps
func (s *AuthSteps) iAmNotAuthenticated() error {
	s.client.AuthToken = ""
	return nil
}

func (s *AuthSteps) iHaveValidUserCredentials() error {
	// Set up valid test credentials
	// These would typically be test user credentials in the test database
	return nil
}

func (s *AuthSteps) iHaveInvalidUserCredentials() error {
	// Set up invalid test credentials
	return nil
}

func (s *AuthSteps) iHaveValidRefreshToken() error {
	// Set up valid refresh token for testing
	// This could be fetched from the test database or generated on the fly
	return nil
}

func (s *AuthSteps) iHaveExpiredRefreshToken() error {
	// Set up expired refresh token for testing
	return nil
}

func (s *AuthSteps) iHaveValidAccessToken() error {
	// Set up valid access token for testing
	return nil
}

func (s *AuthSteps) iHaveExpiredAccessToken() error {
	// Set up expired access token for testing
	return nil
}

// When steps
func (s *AuthSteps) iLoginWithMyCredentials() error {
	// Prepare login request data
	loginData := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}

	// Send login request
	err := s.client.Post("/users/login", loginData)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthSteps) iRequestToRefreshMyTokens() error {
	// Extract refresh token from previous response or use the one set directly
	refreshToken := "test-refresh-token"
	if value, exists := s.client.GetValueFromResponse("tokens.refresh_token"); exists {
		refreshToken = value.(string)
	}

	// Prepare refresh request
	refreshData := map[string]string{
		"refresh_token": refreshToken,
	}

	// Send refresh request
	err := s.client.Post("/users/refresh", refreshData)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthSteps) iLogout() error {
	// Send logout request
	return s.client.Post("/users/logout", nil)
}

func (s *AuthSteps) iUseMyTokenToAccessProtectedResource() error {
	// Use the token to access a protected endpoint (e.g., user profile)
	return s.client.Get("/users/profile")
}

func (s *AuthSteps) iTryToAccessProtectedResourceWithoutAuth() error {
	// Clear any authentication token
	s.client.AuthToken = ""

	// Try to access a protected endpoint
	return s.client.Get("/users/profile")
}

func (s *AuthSteps) iRevokeAllMyUserTokens() error {
	// Send request to revoke all tokens for the current user
	userId := "test-user-id" // This would come from test context or previous responses
	return s.client.Post(fmt.Sprintf("/users/%s/revoke-tokens", userId), nil)
}

// Then steps
func (s *AuthSteps) iShouldReceiveValidTokenPair() error {
	// Check if response contains access_token and refresh_token
	respBody := s.client.GetResponseBodyAsMap()
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

func (s *AuthSteps) iShouldReceiveNewValidTokens() error {
	// Similar to iShouldReceiveValidTokenPair, but could add additional checks
	// for token freshness if needed
	return s.iShouldReceiveValidTokenPair()
}

func (s *AuthSteps) iShouldReceiveAuthenticationError() error {
	// Check for 401 Unauthorized status code
	if s.client.GetResponseStatusCode() != http.StatusUnauthorized {
		return fmt.Errorf("expected status code 401, got %d", s.client.GetResponseStatusCode())
	}

	// Check for error message in response
	respBody := s.client.GetResponseBodyAsMap()
	if respBody == nil {
		return fmt.Errorf("response body is empty or not valid JSON")
	}

	if _, exists := respBody["error"]; !exists {
		return fmt.Errorf("error message not found in response")
	}

	return nil
}

func (s *AuthSteps) iShouldReceiveStatusCode(code int) error {
	if s.client.GetResponseStatusCode() != code {
		return fmt.Errorf("expected status code %d, got %d", code, s.client.GetResponseStatusCode())
	}
	return nil
}

func (s *AuthSteps) theResponseShouldContain(text string) error {
	body := s.client.GetResponseBody()
	if body == nil {
		return fmt.Errorf("response body is empty")
	}

	if !bytes.Contains(body, []byte(text)) {
		return fmt.Errorf("response does not contain %q", text)
	}

	return nil
}

func (s *AuthSteps) myTokensShouldBeInvalidated() error {
	// Try to use the token to access a protected resource
	// This should now fail with 401 Unauthorized
	err := s.client.Get("/users/profile")
	if err != nil {
		return err
	}

	if s.client.GetResponseStatusCode() != http.StatusUnauthorized {
		return fmt.Errorf("token was not properly invalidated, got status code %d", s.client.GetResponseStatusCode())
	}

	return nil
}
