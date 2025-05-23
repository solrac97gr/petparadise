package integration

import (
	"fmt"
	"net/http"

	"github.com/cucumber/godog"
)

// UserSteps contains user-related test steps
type UserSteps struct {
	client *APIClient
}

// RegisterUserSteps registers step definitions for user testing
func RegisterUserSteps(ctx *godog.ScenarioContext, client *APIClient) {
	steps := &UserSteps{client: client}

	// Given steps
	ctx.Step(`^I am authenticated as an? "([^"]*)"$`, steps.iAmAuthenticatedAs)
	ctx.Step(`^there is a user with email "([^"]*)"$`, steps.thereIsAUserWithEmail)
	ctx.Step(`^there is no user with email "([^"]*)"$`, steps.thereIsNoUserWithEmail)
	ctx.Step(`^the user "([^"]*)" has role "([^"]*)"$`, steps.theUserHasRole)
	ctx.Step(`^the user "([^"]*)" has status "([^"]*)"$`, steps.theUserHasStatus)

	// When steps
	ctx.Step(`^I create a new user with name "([^"]*)", email "([^"]*)", password "([^"]*)"$`, steps.iCreateANewUser)
	ctx.Step(`^I get user details for user with ID "([^"]*)"$`, steps.iGetUserDetailsForID)
	ctx.Step(`^I get user details for email "([^"]*)"$`, steps.iGetUserDetailsForEmail)
	ctx.Step(`^I update user "([^"]*)" with name "([^"]*)"$`, steps.iUpdateUserWithName)
	ctx.Step(`^I update user "([^"]*)" role to "([^"]*)"$`, steps.iUpdateUserRoleTo)
	ctx.Step(`^I update user "([^"]*)" status to "([^"]*)"$`, steps.iUpdateUserStatusTo)
	ctx.Step(`^I change password for user "([^"]*)" from "([^"]*)" to "([^"]*)"$`, steps.iChangePasswordForUser)
	ctx.Step(`^I delete user "([^"]*)"$`, steps.iDeleteUser)

	// Then steps
	ctx.Step(`^I should see user details with email "([^"]*)"$`, steps.iShouldSeeUserDetailsWithEmail)
	ctx.Step(`^I should see user role is "([^"]*)"$`, steps.iShouldSeeUserRoleIs)
	ctx.Step(`^I should see user status is "([^"]*)"$`, steps.iShouldSeeUserStatusIs)
	ctx.Step(`^the user should be created successfully$`, steps.theUserShouldBeCreatedSuccessfully)
	ctx.Step(`^the user should be updated successfully$`, steps.theUserShouldBeUpdatedSuccessfully)
	ctx.Step(`^the user should be deleted successfully$`, steps.theUserShouldBeDeletedSuccessfully)
	ctx.Step(`^I should see a list of users$`, steps.iShouldSeeAListOfUsers)
}

// Given step implementations

func (s *UserSteps) iAmAuthenticatedAs(role string) error {
	// Set up authentication for the specified role
	// This would typically involve logging in with a user of the specified role
	loginData := map[string]string{
		"email":    fmt.Sprintf("%s@example.com", role),
		"password": "password123",
	}

	err := s.client.Post("/login", loginData)
	if err != nil {
		return err
	}

	if s.client.GetResponseStatusCode() != http.StatusOK {
		return fmt.Errorf("login failed, got status %d", s.client.GetResponseStatusCode())
	}

	// Extract and set the token
	respBody := s.client.GetResponseBodyAsMap()
	if tokensObj, ok := respBody["tokens"].(map[string]interface{}); ok {
		if accessToken, ok := tokensObj["access_token"].(string); ok {
			s.client.SetAuthToken(accessToken)
		} else {
			return fmt.Errorf("access_token not found in response")
		}
	} else {
		return fmt.Errorf("tokens object not found in response")
	}

	return nil
}

func (s *UserSteps) thereIsAUserWithEmail(email string) error {
	// This step assumes the user exists in the test database
	// In a real implementation, you would verify the user exists or create it if not
	return nil
}

func (s *UserSteps) thereIsNoUserWithEmail(email string) error {
	// This step assumes the user does not exist in the test database
	// In a real implementation, you would verify the user doesn't exist or delete it if it does
	return nil
}

func (s *UserSteps) theUserHasRole(email, role string) error {
	// This step assumes the user exists in the test database with the specified role
	// In a real implementation, you would verify or set the user's role
	return nil
}

func (s *UserSteps) theUserHasStatus(email, status string) error {
	// This step assumes the user exists in the test database with the specified status
	// In a real implementation, you would verify or set the user's status
	return nil
}

// When step implementations

func (s *UserSteps) iCreateANewUser(name, email, password string) error {
	userData := map[string]interface{}{
		"name":     name,
		"email":    email,
		"password": password,
		"role":     "user", // Default role
	}

	return s.client.Post("/register", userData)
}

func (s *UserSteps) iGetUserDetailsForID(id string) error {
	if id == "" {
		return s.client.Get("/")
	}
	return s.client.Get(fmt.Sprintf("/%s", id))
}

func (s *UserSteps) iGetUserDetailsForEmail(email string) error {
	return s.client.Get(fmt.Sprintf("/email?email=%s", email))
}

func (s *UserSteps) iUpdateUserWithName(id, name string) error {
	userData := map[string]interface{}{
		"name": name,
	}

	return s.client.Put(fmt.Sprintf("/%s", id), userData)
}

func (s *UserSteps) iUpdateUserRoleTo(id, role string) error {
	roleData := map[string]interface{}{
		"role": role,
	}

	return s.client.Patch(fmt.Sprintf("/%s/role", id), roleData)
}

func (s *UserSteps) iUpdateUserStatusTo(id, status string) error {
	statusData := map[string]interface{}{
		"status": status,
	}

	return s.client.Patch(fmt.Sprintf("/%s/status", id), statusData)
}

func (s *UserSteps) iChangePasswordForUser(id, oldPassword, newPassword string) error {
	passwordData := map[string]interface{}{
		"old_password": oldPassword,
		"new_password": newPassword,
	}

	return s.client.Post(fmt.Sprintf("/%s/password", id), passwordData)
}

func (s *UserSteps) iDeleteUser(id string) error {
	return s.client.Delete(fmt.Sprintf("/%s", id))
}

// Then step implementations

func (s *UserSteps) iShouldSeeUserDetailsWithEmail(email string) error {
	respBody := s.client.GetResponseBodyAsMap()
	if respBody == nil {
		return fmt.Errorf("response body is empty or not valid JSON")
	}

	var userEmail string
	if user, ok := respBody["user"].(map[string]interface{}); ok {
		if email, ok := user["email"].(string); ok {
			userEmail = email
		}
	} else if email, ok := respBody["email"].(string); ok {
		userEmail = email
	}

	if userEmail != email {
		return fmt.Errorf("expected email %s, got %s", email, userEmail)
	}

	return nil
}

func (s *UserSteps) iShouldSeeUserRoleIs(role string) error {
	respBody := s.client.GetResponseBodyAsMap()
	if respBody == nil {
		return fmt.Errorf("response body is empty or not valid JSON")
	}

	var userRole string
	if user, ok := respBody["user"].(map[string]interface{}); ok {
		if role, ok := user["role"].(string); ok {
			userRole = role
		}
	} else if role, ok := respBody["role"].(string); ok {
		userRole = role
	}

	if userRole != role {
		return fmt.Errorf("expected role %s, got %s", role, userRole)
	}

	return nil
}

func (s *UserSteps) iShouldSeeUserStatusIs(status string) error {
	respBody := s.client.GetResponseBodyAsMap()
	if respBody == nil {
		return fmt.Errorf("response body is empty or not valid JSON")
	}

	var userStatus string
	if user, ok := respBody["user"].(map[string]interface{}); ok {
		if status, ok := user["status"].(string); ok {
			userStatus = status
		}
	} else if status, ok := respBody["status"].(string); ok {
		userStatus = status
	}

	if userStatus != status {
		return fmt.Errorf("expected status %s, got %s", status, userStatus)
	}

	return nil
}

func (s *UserSteps) theUserShouldBeCreatedSuccessfully() error {
	if s.client.GetResponseStatusCode() != http.StatusCreated {
		return fmt.Errorf("expected status code 201, got %d", s.client.GetResponseStatusCode())
	}

	respBody := s.client.GetResponseBodyAsMap()
	if respBody == nil {
		return fmt.Errorf("response body is empty or not valid JSON")
	}

	if _, exists := respBody["id"]; !exists {
		return fmt.Errorf("user ID not found in response")
	}

	return nil
}

func (s *UserSteps) theUserShouldBeUpdatedSuccessfully() error {
	if s.client.GetResponseStatusCode() != http.StatusOK {
		return fmt.Errorf("expected status code 200, got %d", s.client.GetResponseStatusCode())
	}
	return nil
}

func (s *UserSteps) theUserShouldBeDeletedSuccessfully() error {
	if s.client.GetResponseStatusCode() != http.StatusNoContent {
		return fmt.Errorf("expected status code 204, got %d", s.client.GetResponseStatusCode())
	}
	return nil
}

func (s *UserSteps) iShouldSeeAListOfUsers() error {
	if s.client.GetResponseStatusCode() != http.StatusOK {
		return fmt.Errorf("expected status code 200, got %d", s.client.GetResponseStatusCode())
	}

	respBody := s.client.GetResponseBody()
	if respBody == nil || len(respBody) == 0 {
		return fmt.Errorf("response body is empty")
	}

	// Check if the response is an array
	if respBody[0] != '[' {
		return fmt.Errorf("response is not a JSON array")
	}

	return nil
}
