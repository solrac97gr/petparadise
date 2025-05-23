package integration

import (
	"fmt"
	"net/http"

	"github.com/cucumber/godog"
)

// AdoptionSteps contains adoption-related test steps
type AdoptionSteps struct {
	client *APIClient
}

// RegisterAdoptionSteps registers step definitions for adoption testing
func RegisterAdoptionSteps(ctx *godog.ScenarioContext, client *APIClient) {
	steps := &AdoptionSteps{client: client}

	// Given steps
	ctx.Step(`^I have submitted an adoption request$`, steps.iHaveSubmittedAdoptionRequest)
	ctx.Step(`^I have submitted an adoption request with ID "([^"]*)"$`, steps.iHaveSubmittedAdoptionRequestWithID)
	ctx.Step(`^there is a pending adoption request with ID "([^"]*)"$`, steps.thereIsPendingAdoptionRequest)

	// When steps
	ctx.Step(`^I create an adoption request for pet "([^"]*)"$`, steps.iCreateAdoptionRequest)
	ctx.Step(`^I request my adoption requests$`, steps.iRequestMyAdoptionRequests)
	ctx.Step(`^I cancel the adoption request with ID "([^"]*)"$`, steps.iCancelAdoptionRequest)
	ctx.Step(`^I approve the adoption request with ID "([^"]*)"$`, steps.iApproveAdoptionRequest)
	ctx.Step(`^I reject the adoption request with ID "([^"]*)"$`, steps.iRejectAdoptionRequest)

	// Then steps
	ctx.Step(`^the adoption request should be created successfully$`, steps.adoptionRequestShouldBeCreatedSuccessfully)
	ctx.Step(`^I should see a list of my adoption requests$`, steps.iShouldSeeListOfMyAdoptionRequests)
	ctx.Step(`^the adoption request status should be "([^"]*)"$`, steps.adoptionRequestStatusShouldBe)
}

// Given step implementations

func (s *AdoptionSteps) iHaveSubmittedAdoptionRequest() error {
	// This step assumes the user has already submitted an adoption request
	// In a real implementation, you would create an adoption request or verify one exists
	return nil
}

func (s *AdoptionSteps) iHaveSubmittedAdoptionRequestWithID(id string) error {
	// This step assumes the user has already submitted an adoption request with the given ID
	// In a real implementation, you would create an adoption request or verify it exists
	return nil
}

func (s *AdoptionSteps) thereIsPendingAdoptionRequest(id string) error {
	// This step assumes there is a pending adoption request with the given ID
	// In a real implementation, you would create or verify such a request exists
	return nil
}

// When step implementations

func (s *AdoptionSteps) iCreateAdoptionRequest(petID string) error {
	adoptionData := map[string]interface{}{
		"pet_id":     petID,
		"notes":      "I would love to adopt this pet!",
		"home_type":  "apartment",
		"has_yard":   false,
		"has_kids":   false,
		"has_pets":   true,
		"work_hours": "9-5",
	}

	return s.client.Post("/adoptions", adoptionData)
}

func (s *AdoptionSteps) iRequestMyAdoptionRequests() error {
	return s.client.Get("/adoptions/my")
}

func (s *AdoptionSteps) iCancelAdoptionRequest(id string) error {
	statusData := map[string]interface{}{
		"status": "cancelled",
	}

	return s.client.Patch(fmt.Sprintf("/adoptions/%s/status", id), statusData)
}

func (s *AdoptionSteps) iApproveAdoptionRequest(id string) error {
	statusData := map[string]interface{}{
		"status": "approved",
	}

	return s.client.Patch(fmt.Sprintf("/adoptions/%s/status", id), statusData)
}

func (s *AdoptionSteps) iRejectAdoptionRequest(id string) error {
	statusData := map[string]interface{}{
		"status": "rejected",
	}

	return s.client.Patch(fmt.Sprintf("/adoptions/%s/status", id), statusData)
}

// Then step implementations

func (s *AdoptionSteps) adoptionRequestShouldBeCreatedSuccessfully() error {
	if s.client.GetResponseStatusCode() != http.StatusCreated {
		return fmt.Errorf("expected status code 201, got %d", s.client.GetResponseStatusCode())
	}

	respBody := s.client.GetResponseBodyAsMap()
	if respBody == nil {
		return fmt.Errorf("response body is empty or not valid JSON")
	}

	if _, exists := respBody["id"]; !exists {
		return fmt.Errorf("adoption request ID not found in response")
	}

	return nil
}

func (s *AdoptionSteps) iShouldSeeListOfMyAdoptionRequests() error {
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

func (s *AdoptionSteps) adoptionRequestStatusShouldBe(status string) error {
	respBody := s.client.GetResponseBodyAsMap()
	if respBody == nil {
		return fmt.Errorf("response body is empty or not valid JSON")
	}

	adoptionStatus, ok := respBody["status"].(string)
	if !ok {
		return fmt.Errorf("adoption status not found in response")
	}

	if adoptionStatus != status {
		return fmt.Errorf("expected adoption status %s, got %s", status, adoptionStatus)
	}

	return nil
}
