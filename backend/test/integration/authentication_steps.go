package integration

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// AuthSteps contains authentication test steps
type AuthSteps struct {
	client               *APIClient
	db                   *sqlx.DB
	credentialsValid     bool
	testEmail            string
	testPassword         string
	testToken            string
	testUserID           string // Store the user ID from authentication
	explicitRefreshToken string // Added to store an explicitly set refresh token
}

// RegisterAuthenticationSteps registers step definitions for authentication scenarios
func RegisterAuthenticationSteps(ctx *godog.ScenarioContext, client *APIClient, db *sqlx.DB) {
	steps := &AuthSteps{client: client, db: db}

	// Given steps
	ctx.Step(`^I am authenticated as an? "([^"]*)"$`, steps.iAmAuthenticatedAs)
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
func (s *AuthSteps) iAmAuthenticatedAs(role string) error {
	randomUUID := uuid.New().String()
	s.testEmail = "test" + randomUUID + role + "@example.com"
	s.testPassword = "password123"

	loginData := map[string]string{
		"email":    s.testEmail,
		"password": s.testPassword,
	}

	// For login we need to create a new user in the database matching the test credentials
	// the password should be hashed in a real scenario
	if err := s.client.Post("/users/register", map[string]string{
		"name":     "Test User",
		"email":    s.testEmail,
		"password": s.testPassword,
	}); err != nil {
		return fmt.Errorf("failed to create test user: %v", err)
	}
	if s.client.GetResponseStatusCode() != http.StatusCreated {
		return fmt.Errorf("failed to create test user, got status %d, and body %v", s.client.GetResponseStatusCode(), string(s.client.GetResponseBody()))
	}

	// we convert the user to the specified role for testing purposes
	res, err := s.db.Exec("UPDATE users SET role = $1 WHERE email = $2", role, s.testEmail)
	if err != nil {
		return fmt.Errorf("failed to update user role: %v", err)
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("no rows affected when updating user role")
	}
	// Send login request
	if err := s.client.Post("/users/login", loginData); err != nil {
		return err
	}

	if s.client.GetResponseStatusCode() != http.StatusOK {
		fmt.Println("Login failed:", string(s.client.GetResponseBody()))
		return fmt.Errorf("login failed, got status %d", s.client.GetResponseStatusCode())
	}

	// Extract and set the token
	respBody := s.client.GetResponseBodyAsMap()
	if tokensObj, ok := respBody["tokens"].(map[string]interface{}); ok {
		if accessToken, ok := tokensObj["access_token"].(string); ok {
			s.client.SetAuthToken(accessToken)
			s.testToken = accessToken
		} else {
			return fmt.Errorf("access_token not found in response")
		}
	} else {
		return fmt.Errorf("tokens object not found in response")
	}

	// Extract and store the user ID
	if userObj, ok := respBody["user"].(map[string]interface{}); ok {
		if userID, ok := userObj["id"].(string); ok {
			s.testUserID = userID
		} else {
			return fmt.Errorf("user id not found or not a string")
		}
	} else {
		return fmt.Errorf("user object not found in response")
	}

	return nil
}

// Given steps
func (s *AuthSteps) iAmNotAuthenticated() error {
	s.client.AuthToken = ""
	return nil
}

func (s *AuthSteps) iHaveValidUserCredentials() error {
	// Set up valid test credentials from our test data
	s.credentialsValid = true
	s.testEmail = "user@example.com"
	s.testPassword = "password123"
	return nil
}

func (s *AuthSteps) iHaveInvalidUserCredentials() error {
	// Set up invalid test credentials
	s.credentialsValid = false
	s.testEmail = "invalid@test.com"
	s.testPassword = "wrongpassword"
	return nil
}

func (s *AuthSteps) iHaveValidRefreshToken() error {
	// The refresh token should already be available from the login step
	// Let's verify we have one from the current response
	respBody := s.client.GetResponseBodyAsMap()
	if respBody == nil {
		return fmt.Errorf("no response body available")
	}

	// Check for tokens object
	tokensObj, ok := respBody["tokens"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("tokens object not found in login response")
	}

	// Check for refresh_token
	if _, exists := tokensObj["refresh_token"]; !exists {
		return fmt.Errorf("refresh_token not found in login response")
	}

	return nil
}

func (s *AuthSteps) iHaveExpiredRefreshToken() error {
	// Set an expired/invalid token to be used by a subsequent step
	s.explicitRefreshToken = "expired.or.invalid.token"
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

	// Use the credentials set in the previous steps
	email := s.testEmail
	password := s.testPassword

	// Default to test user if no specific credentials were set
	if email == "" {
		email = "user@example.com"
		password = "password123"
	}

	// Prepare login request data
	loginData := map[string]string{
		"email":    email,
		"password": password,
	}

	// For login we need to create a new user in the database matching the test credentials
	// the password should be hashed in a real scenario
	if err := s.client.Post("/users/register", map[string]string{
		"name":     "Test User",
		"email":    email,
		"password": password,
	}); err != nil {
		return fmt.Errorf("failed to create test user: %v", err)
	}
	if s.client.GetResponseStatusCode() != http.StatusCreated {
		return fmt.Errorf("failed to create test user, got status %d, and body %v", s.client.GetResponseStatusCode(), string(s.client.GetResponseBody()))
	}

	// we convert the user to an admin for testing purposes
	res, err := s.db.Exec("UPDATE users SET role = 'admin' WHERE email = $1", email)
	if err != nil {
		return fmt.Errorf("failed to update user role: %v", err)
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("no rows affected when updating user role")
	}
	if !s.credentialsValid {
		loginData["email"] = "invalid@email.com"
	}

	// Send login request
	if err := s.client.Post("/users/login", loginData); err != nil {
		return err
	}

	return nil
}

func (s *AuthSteps) iRequestToRefreshMyTokens() error {
	var tokenForRefresh string

	if s.explicitRefreshToken != "" {
		tokenForRefresh = s.explicitRefreshToken
		s.explicitRefreshToken = "" // Consume the explicitly set token
	} else {
		// Fallback to getting token from previous response (e.g., after a login)
		respBody := s.client.GetResponseBodyAsMap()
		if respBody == nil {
			return fmt.Errorf("no response body available to extract refresh token for refresh request")
		}
		tokensObj, ok := respBody["tokens"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("tokens object not found in previous response to extract refresh token")
		}
		refreshTokenFromResponse, ok := tokensObj["refresh_token"].(string)
		if !ok {
			return fmt.Errorf("refresh_token not found or not a string in previous response")
		}
		tokenForRefresh = refreshTokenFromResponse
	}

	if tokenForRefresh == "" {
		return fmt.Errorf("no refresh token available to make the refresh request")
	}

	// Prepare refresh request
	refreshData := map[string]string{
		"refresh_token": tokenForRefresh,
	}

	// Send refresh request
	// The client.Post method should handle storing the response (status code, body)
	// so subsequent "Then" steps can assert on it.
	return s.client.Post("/users/refresh", refreshData)
}

func (s *AuthSteps) iLogout() error {
	// Send logout request
	return s.client.Post("/users/logout", nil)
}

func (s *AuthSteps) iUseMyTokenToAccessProtectedResource() error {
	// Use the token to access a protected endpoint (e.g., get all users)
	if err := s.client.Get("/users/"); err != nil {

		return fmt.Errorf("failed to access protected resource: %v", err)
	}
	return nil
}

func (s *AuthSteps) iTryToAccessProtectedResourceWithoutAuth() error {
	// Clear any authentication token
	s.client.AuthToken = ""

	// Try to access a protected endpoint
	return s.client.Get("/users/")
}

func (s *AuthSteps) iRevokeAllMyUserTokens() error {
	// Use the stored user ID from authentication
	if s.testUserID == "" {
		return fmt.Errorf("no user ID available - user may not be authenticated")
	}
	// Send request to revoke all tokens for the current user
	if err := s.client.Post(fmt.Sprintf("/users/%s/revoke-tokens", s.testUserID), nil); err != nil {
		return fmt.Errorf("failed to revoke tokens: %v", err)
	}
	return nil
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
		return fmt.Errorf("tokens object not found in response body: %v", respBody)
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
		return fmt.Errorf("expected status code %d, got %d, with body: %s", code, s.client.GetResponseStatusCode(), string(s.client.GetResponseBody()))
	}
	return nil
}

func (s *AuthSteps) theResponseShouldContain(text string) error {
	body := s.client.GetResponseBody()
	if body == nil {
		return fmt.Errorf("response body is empty")
	}

	if !bytes.Contains(body, []byte(text)) {
		return fmt.Errorf("response does not contain %q. Actual response: %s", text, string(body))
	}

	return nil
}

func (s *AuthSteps) myTokensShouldBeInvalidated() error {
	// Try to use the token to access a protected resource
	// This should now fail with 401 Unauthorized
	err := s.client.Get("/users/")
	if err != nil {
		return err
	}

	if s.client.GetResponseStatusCode() != http.StatusUnauthorized {
		return fmt.Errorf("token was not properly invalidated, got status code %d", s.client.GetResponseStatusCode())
	}

	// Clear the auth token from the client after confirming it's invalid
	// This prevents the invalid token from affecting subsequent test scenarios
	s.client.SetAuthToken("")

	return nil
}
