package integration

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// UserSteps contains user management test steps
type UserSteps struct {
	client               *APIClient
	db                   *sqlx.DB
	testUser             map[string]interface{}
	anotherTestUser      map[string]interface{}
	registrationData     map[string]interface{}
	registrationValid    bool
	testUserID           string
	anotherTestUserID    string
	updateData           map[string]interface{}
	originalPassword     string
	newPassword          string
	wrongCurrentPassword string
}

// RegisterUserSteps registers step definitions for user management scenarios
func RegisterUserSteps(ctx *godog.ScenarioContext, client *APIClient, db *sqlx.DB) {
	steps := &UserSteps{client: client, db: db}

	// Given steps
	ctx.Step(`^the system is initialized$`, steps.theSystemIsInitialized)
	ctx.Step(`^I have valid registration data$`, steps.iHaveValidRegistrationData)
	ctx.Step(`^I have invalid registration data$`, steps.iHaveInvalidRegistrationData)
	ctx.Step(`^I have registration data with an existing email$`, steps.iHaveRegistrationDataWithExistingEmail)
	ctx.Step(`^a user exists in the system$`, steps.aUserExistsInTheSystem)
	ctx.Step(`^another user exists in the system$`, steps.anotherUserExistsInTheSystem)
	ctx.Step(`^users with different statuses exist$`, steps.usersWithDifferentStatusesExist)

	// When steps
	ctx.Step(`^I register a new user$`, steps.iRegisterANewUser)
	ctx.Step(`^I request all users$`, steps.iRequestAllUsers)
	ctx.Step(`^I request the user by ID$`, steps.iRequestTheUserByID)
	ctx.Step(`^I request the other user by ID$`, steps.iRequestTheOtherUserByID)
	ctx.Step(`^I request the user by email$`, steps.iRequestTheUserByEmail)
	ctx.Step(`^I request users by status "([^"]*)"$`, steps.iRequestUsersByStatus)
	ctx.Step(`^I update the user information$`, steps.iUpdateTheUserInformation)
	ctx.Step(`^I try to update the other user information$`, steps.iTryToUpdateTheOtherUserInformation)
	ctx.Step(`^I update the user role to "([^"]*)"$`, steps.iUpdateTheUserRoleTo)
	ctx.Step(`^I try to update the user role to "([^"]*)"$`, steps.iTryToUpdateTheUserRoleTo)
	ctx.Step(`^I update the user status to "([^"]*)"$`, steps.iUpdateTheUserStatusTo)
	ctx.Step(`^I try to update the user status to "([^"]*)"$`, steps.iTryToUpdateTheUserStatusTo)
	ctx.Step(`^I delete the user$`, steps.iDeleteTheUser)
	ctx.Step(`^I try to delete the other user$`, steps.iTryToDeleteTheOtherUser)
	ctx.Step(`^I try to access user endpoints without authentication$`, steps.iTryToAccessUserEndpointsWithoutAuthentication)

	// Then steps
	ctx.Step(`^the response should contain user details$`, steps.theResponseShouldContainUserDetails)
	ctx.Step(`^the user should be created in the database$`, steps.theUserShouldBeCreatedInTheDatabase)
	ctx.Step(`^the response should contain an error message$`, steps.theResponseShouldContainAnErrorMessage)
	ctx.Step(`^the response should contain a list of users$`, steps.theResponseShouldContainAListOfUsers)
	ctx.Step(`^the response should contain the user details$`, steps.theResponseShouldContainTheUserDetails)
	ctx.Step(`^the response should contain only active users$`, steps.theResponseShouldContainOnlyActiveUsers)
	ctx.Step(`^the user information should be updated$`, steps.theUserInformationShouldBeUpdated)
	ctx.Step(`^the user role should be updated$`, steps.theUserRoleShouldBeUpdated)
	ctx.Step(`^the user status should be updated$`, steps.theUserStatusShouldBeUpdated)
	ctx.Step(`^the user should be removed from the system$`, steps.theUserShouldBeRemovedFromTheSystem)
}

// Given step implementations
func (s *UserSteps) theSystemIsInitialized() error {
	// Clean up any existing test data
	s.db.Exec("DELETE FROM users WHERE email LIKE '%test%' OR email LIKE '%example%'")
	return nil
}

func (s *UserSteps) iHaveValidRegistrationData() error {
	randomUUID := uuid.New().String()
	s.registrationData = map[string]interface{}{
		"name":      "Test User",
		"email":     "testuser" + randomUUID + "@example.com",
		"password":  "password123",
		"role":      "user",
		"address":   "123 Test Street",
		"phone":     "+1234567890",
		"documents": []string{"doc1.pdf", "doc2.pdf"},
	}
	s.registrationValid = true
	return nil
}

func (s *UserSteps) iHaveInvalidRegistrationData() error {
	s.registrationData = map[string]interface{}{
		"name":     "", // Invalid: empty name
		"email":    "invalid-email",
		"password": "",
	}
	s.registrationValid = false
	return nil
}

func (s *UserSteps) iHaveRegistrationDataWithExistingEmail() error {
	// First, create a user
	randomUUID := uuid.New().String()
	existingEmail := "existing" + randomUUID + "@example.com"

	err := s.client.Post("/users/register", map[string]interface{}{
		"name":     "Existing User",
		"email":    existingEmail,
		"password": "password123",
		"role":     "user",
		"address":  "123 Existing Street",
		"phone":    "+1234567890",
	})
	if err != nil {
		return fmt.Errorf("failed to create existing user: %v", err)
	}

	// Now set registration data with the same email
	s.registrationData = map[string]interface{}{
		"name":     "Test User",
		"email":    existingEmail,
		"password": "password123",
		"role":     "user",
		"address":  "123 Test Street",
		"phone":    "+1234567890",
	}
	s.registrationValid = false
	return nil
}

func (s *UserSteps) aUserExistsInTheSystem() error {
	randomUUID := uuid.New().String()
	s.testUser = map[string]interface{}{
		"name":     "Test User",
		"email":    "testuser" + randomUUID + "@example.com",
		"password": "password123",
		"role":     "user",
		"address":  "123 Test Street",
		"phone":    "+1234567890",
	}

	err := s.client.Post("/users/register", s.testUser)
	if err != nil {
		return fmt.Errorf("failed to create test user: %v", err)
	}

	if s.client.GetResponseStatusCode() != http.StatusCreated {
		return fmt.Errorf("failed to create test user, got status %d", s.client.GetResponseStatusCode())
	}

	// Extract user ID from response (user object is returned directly)
	respBody := s.client.GetResponseBodyAsMap()
	if userID, ok := respBody["id"].(string); ok {
		s.testUserID = userID
	}

	return nil
}

func (s *UserSteps) anotherUserExistsInTheSystem() error {
	randomUUID := uuid.New().String()
	s.anotherTestUser = map[string]interface{}{
		"name":     "Another Test User",
		"email":    "anothertestuser" + randomUUID + "@example.com",
		"password": "password123",
		"role":     "user",
		"address":  "456 Another Street",
		"phone":    "+0987654321",
	}

	err := s.client.Post("/users/register", s.anotherTestUser)
	if err != nil {
		return fmt.Errorf("failed to create another test user: %v", err)
	}

	if s.client.GetResponseStatusCode() != http.StatusCreated {
		return fmt.Errorf("failed to create another test user, got status %d", s.client.GetResponseStatusCode())
	}

	// Extract user ID from response (user object is returned directly)
	respBody := s.client.GetResponseBodyAsMap()
	if userID, ok := respBody["id"].(string); ok {
		s.anotherTestUserID = userID
	}

	return nil
}

func (s *UserSteps) usersWithDifferentStatusesExist() error {
	// Create active user
	randomUUID1 := uuid.New().String()
	activeUser := map[string]interface{}{
		"name":     "Active User",
		"email":    "activeuser" + randomUUID1 + "@example.com",
		"password": "password123",
		"role":     "user",
		"address":  "123 Active Street",
		"phone":    "+1111111111",
	}

	err := s.client.Post("/users/register", activeUser)
	if err != nil {
		return fmt.Errorf("failed to create active user: %v", err)
	}

	// Create another user and make it inactive
	randomUUID2 := uuid.New().String()
	inactiveUserEmail := "inactiveuser" + randomUUID2 + "@example.com"
	inactiveUser := map[string]interface{}{
		"name":     "Inactive User",
		"email":    inactiveUserEmail,
		"password": "password123",
		"role":     "user",
		"address":  "456 Inactive Street",
		"phone":    "+2222222222",
	}

	err = s.client.Post("/users/register", inactiveUser)
	if err != nil {
		return fmt.Errorf("failed to create inactive user: %v", err)
	}

	// Update user status to inactive in database
	_, err = s.db.Exec("UPDATE users SET status = 'inactive' WHERE email = $1", inactiveUserEmail)
	if err != nil {
		return fmt.Errorf("failed to update user status: %v", err)
	}

	return nil
}

// When step implementations
func (s *UserSteps) iRegisterANewUser() error {
	return s.client.Post("/users/register", s.registrationData)
}

func (s *UserSteps) iRequestAllUsers() error {
	s.client.SetAuthToken(s.client.AuthToken)
	return s.client.Get("/users/")
}

func (s *UserSteps) iRequestTheUserByID() error {
	if s.testUserID == "" {
		return fmt.Errorf("no test user ID available")
	}
	return s.client.Get("/users/" + s.testUserID)
}

func (s *UserSteps) iRequestMyOwnUserDetailsByID() error {
	// Get current user ID from authentication context
	// We need to extract the user ID from the login response or JWT token
	// For now, let's assume we can get it from the last authentication response
	respBody := s.client.GetResponseBodyAsMap()
	var userID string

	if respBody != nil {
		if userObj, ok := respBody["user"].(map[string]interface{}); ok {
			if id, ok := userObj["id"].(string); ok {
				userID = id
			}
		}
	}

	// If we don't have a user ID from response, we might need to get it differently
	// For testing purposes, let's use the authenticated user's info
	if userID == "" {
		// Try to get current user info first
		err := s.client.Get("/users/")
		if err != nil {
			return err
		}
		// This will get all users, we'd need a different endpoint for current user
		// For now, we'll just use the first available user ID
		return fmt.Errorf("unable to get current user ID - this needs proper implementation")
	}

	return s.client.Get("/users/" + userID)
}

func (s *UserSteps) iRequestTheOtherUserByID() error {
	if s.anotherTestUserID == "" {
		return fmt.Errorf("no other test user ID available")
	}
	return s.client.Get("/users/" + s.anotherTestUserID)
}

func (s *UserSteps) iRequestTheUserByEmail() error {
	if s.testUser == nil {
		return fmt.Errorf("no test user available")
	}
	email := s.testUser["email"].(string)
	return s.client.Get("/users/email?email=" + email)
}

func (s *UserSteps) iRequestUsersByStatus(status string) error {
	return s.client.Get("/users/status?status=" + status)
}

func (s *UserSteps) iUpdateTheUserInformation() error {
	if s.testUserID == "" {
		return fmt.Errorf("no test user ID available")
	}

	s.updateData = map[string]interface{}{
		"name":    "Updated Test User",
		"address": "456 Updated Street",
		"phone":   "+9999999999",
	}

	return s.client.Put("/users/"+s.testUserID, s.updateData)
}

func (s *UserSteps) iUpdateMyOwnUserInformation() error {
	// For this scenario, we need to get the current user's ID from the authentication context
	// Since the API requires user ID in the path, we need to extract it
	// This would typically be extracted from the JWT token or stored during login

	s.updateData = map[string]interface{}{
		"name":    "Updated My User",
		"address": "789 My Updated Street",
		"phone":   "+8888888888",
	}

	// We need to get the authenticated user's ID
	// For now, let's assume we can get it from the current auth context
	// In a real scenario, you'd get the user ID from the JWT token
	// For testing, we'll need to extract it from the login response

	// This is a placeholder - in reality you'd extract from JWT or have a /users/me endpoint
	return fmt.Errorf("this scenario needs a /users/me endpoint or JWT parsing implementation")
}

func (s *UserSteps) iTryToUpdateTheOtherUserInformation() error {
	if s.anotherTestUserID == "" {
		return fmt.Errorf("no other test user ID available")
	}

	updateData := map[string]interface{}{
		"name":    "Unauthorized Update",
		"address": "Unauthorized Street",
	}

	return s.client.Put("/users/"+s.anotherTestUserID, updateData)
}

func (s *UserSteps) iUpdateTheUserRoleTo(role string) error {
	if s.testUserID == "" {
		return fmt.Errorf("no test user ID available")
	}

	roleData := map[string]interface{}{
		"role": role,
	}

	return s.client.Patch("/users/"+s.testUserID+"/role", roleData)
}

func (s *UserSteps) iTryToUpdateTheUserRoleTo(role string) error {
	if s.anotherTestUserID == "" {
		return fmt.Errorf("no other test user ID available")
	}

	roleData := map[string]interface{}{
		"role": role,
	}

	return s.client.Patch("/users/"+s.anotherTestUserID+"/role", roleData)
}

func (s *UserSteps) iUpdateTheUserStatusTo(status string) error {
	if s.testUserID == "" {
		return fmt.Errorf("no test user ID available")
	}

	statusData := map[string]interface{}{
		"status": status,
	}

	return s.client.Patch("/users/"+s.testUserID+"/status", statusData)
}

func (s *UserSteps) iTryToUpdateTheUserStatusTo(status string) error {
	if s.anotherTestUserID == "" {
		return fmt.Errorf("no other test user ID available")
	}

	statusData := map[string]interface{}{
		"status": status,
	}

	return s.client.Patch("/users/"+s.anotherTestUserID+"/status", statusData)
}

func (s *UserSteps) iChangeMyPassword() error {
	s.originalPassword = "password123"
	s.newPassword = "newpassword456"

	// The API expects /:id/password, so we need to get the current user ID
	// For testing, we'll need to extract this from authentication context
	return fmt.Errorf("this scenario needs user ID from authentication context")
}

func (s *UserSteps) iTryToChangeMyPasswordWithWrongCurrentPassword() error {
	s.originalPassword = "password123"
	s.wrongCurrentPassword = "wrongpassword"
	s.newPassword = "newpassword456"

	// The API expects /:id/password, so we need to get the current user ID
	return fmt.Errorf("this scenario needs user ID from authentication context")
}

func (s *UserSteps) iDeleteTheUser() error {
	if s.testUserID == "" {
		return fmt.Errorf("no test user ID available")
	}
	return s.client.Delete("/users/" + s.testUserID)
}

func (s *UserSteps) iTryToDeleteTheOtherUser() error {
	if s.anotherTestUserID == "" {
		return fmt.Errorf("no other test user ID available")
	}
	return s.client.Delete("/users/" + s.anotherTestUserID)
}

func (s *UserSteps) iTryToAccessUserEndpointsWithoutAuthentication() error {
	// Clear authentication
	s.client.AuthToken = ""
	return s.client.Get("/users/")
}

func (s *UserSteps) iRevokeAllMyTokens() error {
	// The API expects /:id/revoke-tokens, so we need to get the current user ID
	// For testing, we'll need to extract this from authentication context
	return fmt.Errorf("this scenario needs user ID from authentication context")
}

// Then step implementations
func (s *UserSteps) theResponseShouldContainUserDetails() error {
	respBody := s.client.GetResponseBodyAsMap()
	if respBody == nil {
		return fmt.Errorf("response body is empty or not valid JSON")
	}

	// Check that essential fields are present in the user object (returned directly)
	if _, exists := respBody["id"]; !exists {
		return fmt.Errorf("user ID not found in response")
	}
	if _, exists := respBody["email"]; !exists {
		return fmt.Errorf("user email not found in response")
	}
	if _, exists := respBody["name"]; !exists {
		return fmt.Errorf("user name not found in response")
	}
	// Ensure password is not returned
	if _, exists := respBody["password"]; exists {
		return fmt.Errorf("password should not be returned in response")
	}

	return nil
}

func (s *UserSteps) theUserShouldBeCreatedInTheDatabase() error {
	if s.registrationData == nil {
		return fmt.Errorf("no registration data available")
	}

	email := s.registrationData["email"].(string)
	var count int
	err := s.db.Get(&count, "SELECT COUNT(*) FROM users WHERE email = $1", email)
	if err != nil {
		return fmt.Errorf("failed to check user in database: %v", err)
	}

	if count != 1 {
		return fmt.Errorf("user was not created in database")
	}

	return nil
}

func (s *UserSteps) theResponseShouldContainAnErrorMessage() error {
	respBody := s.client.GetResponseBodyAsMap()
	if respBody == nil {
		return fmt.Errorf("response body is empty or not valid JSON")
	}

	if _, exists := respBody["error"]; !exists {
		return fmt.Errorf("error message not found in response")
	}

	return nil
}

func (s *UserSteps) theResponseShouldContainAListOfUsers() error {
	respBody := s.client.GetResponseBody()
	if respBody == nil {
		return fmt.Errorf("response body is empty")
	}

	// Check if response is an array or contains an array
	var users []interface{}
	if err := json.Unmarshal(respBody, &users); err != nil {
		// Try to get users from a wrapper object
		var wrapper map[string]interface{}
		if err := json.Unmarshal(respBody, &wrapper); err != nil {
			return fmt.Errorf("response is not valid JSON")
		}

		if usersArray, ok := wrapper["users"].([]interface{}); ok {
			users = usersArray
		} else {
			return fmt.Errorf("users array not found in response")
		}
	}

	if len(users) == 0 {
		return fmt.Errorf("response does not contain any users")
	}

	return nil
}

func (s *UserSteps) theResponseShouldContainTheUserDetails() error {
	return s.theResponseShouldContainUserDetails()
}

func (s *UserSteps) theResponseShouldContainMyUserDetails() error {
	return s.theResponseShouldContainUserDetails()
}

func (s *UserSteps) theResponseShouldContainOnlyActiveUsers() error {
	respBody := s.client.GetResponseBody()
	if respBody == nil {
		return fmt.Errorf("response body is empty")
	}

	var users []map[string]interface{}
	if err := json.Unmarshal(respBody, &users); err != nil {
		return fmt.Errorf("response is not valid JSON: %v", err)
	}

	for _, user := range users {
		if status, ok := user["status"].(string); ok {
			if status != "active" {
				return fmt.Errorf("found non-active user in response: status = %s", status)
			}
		} else {
			return fmt.Errorf("user status not found or not a string")
		}
	}

	return nil
}

func (s *UserSteps) theUserInformationShouldBeUpdated() error {
	if s.testUserID == "" || s.updateData == nil {
		return fmt.Errorf("no test user ID or update data available")
	}

	// Fetch the updated user
	err := s.client.Get("/users/" + s.testUserID)
	if err != nil {
		return fmt.Errorf("failed to fetch updated user: %v", err)
	}

	respBody := s.client.GetResponseBodyAsMap()
	if respBody == nil {
		return fmt.Errorf("response body is empty")
	}

	// Check if updated fields match (user object is returned directly)
	for key, expectedValue := range s.updateData {
		if actualValue, exists := respBody[key]; !exists || actualValue != expectedValue {
			return fmt.Errorf("field %s was not updated correctly: expected %v, got %v", key, expectedValue, actualValue)
		}
	}

	return nil
}

func (s *UserSteps) myUserInformationShouldBeUpdated() error {
	return s.theUserInformationShouldBeUpdated()
}

func (s *UserSteps) theUserRoleShouldBeUpdated() error {
	if s.testUserID == "" {
		return fmt.Errorf("no test user ID available")
	}

	// Check in database
	var role string
	err := s.db.Get(&role, "SELECT role FROM users WHERE id = $1", s.testUserID)
	if err != nil {
		return fmt.Errorf("failed to get user role from database: %v", err)
	}

	if role != "admin" {
		return fmt.Errorf("user role was not updated in database: expected admin, got %s", role)
	}

	return nil
}

func (s *UserSteps) theUserStatusShouldBeUpdated() error {
	if s.testUserID == "" {
		return fmt.Errorf("no test user ID available")
	}

	// Check in database
	var status string
	err := s.db.Get(&status, "SELECT status FROM users WHERE id = $1", s.testUserID)
	if err != nil {
		return fmt.Errorf("failed to get user status from database: %v", err)
	}

	if status != "inactive" {
		return fmt.Errorf("user status was not updated in database: expected inactive, got %s", status)
	}

	return nil
}

func (s *UserSteps) iShouldBeAbleToLoginWithTheNewPassword() error {
	if s.testUser == nil || s.newPassword == "" {
		return fmt.Errorf("no test user or new password available")
	}

	loginData := map[string]string{
		"email":    s.testUser["email"].(string),
		"password": s.newPassword,
	}

	err := s.client.Post("/users/login", loginData)
	if err != nil {
		return fmt.Errorf("failed to login with new password: %v", err)
	}

	if s.client.GetResponseStatusCode() != http.StatusOK {
		return fmt.Errorf("login with new password failed: got status %d", s.client.GetResponseStatusCode())
	}

	return nil
}

func (s *UserSteps) theUserShouldBeRemovedFromTheSystem() error {
	if s.testUserID == "" {
		return fmt.Errorf("no test user ID available")
	}

	// Check that user no longer exists in database
	var count int
	err := s.db.Get(&count, "SELECT COUNT(*) FROM users WHERE id = $1", s.testUserID)
	if err != nil {
		return fmt.Errorf("failed to check user in database: %v", err)
	}

	if count != 0 {
		return fmt.Errorf("user was not removed from database")
	}

	return nil
}

func (s *UserSteps) allMyTokensShouldBeInvalidated() error {
	// Try to access a protected resource
	err := s.client.Get("/users/")
	if err != nil {
		return err
	}

	if s.client.GetResponseStatusCode() != http.StatusUnauthorized {
		return fmt.Errorf("tokens were not properly invalidated, got status code %d", s.client.GetResponseStatusCode())
	}

	// Clear the auth token from the client after confirming it's invalid
	// This prevents the invalid token from affecting subsequent test scenarios
	s.client.SetAuthToken("")

	return nil
}
