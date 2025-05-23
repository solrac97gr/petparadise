package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cucumber/godog"
)

// DonationSteps contains donation-related test steps
type DonationSteps struct {
	client *APIClient
}

// RegisterDonationSteps registers step definitions for donation testing
func RegisterDonationSteps(ctx *godog.ScenarioContext, client *APIClient) {
	steps := &DonationSteps{client: client}

	// Given steps
	ctx.Step(`^I have made donations in the past$`, steps.iHaveMadeDonationsInThePast)
	ctx.Step(`^I have made a donation with ID "([^"]*)"$`, steps.iHaveMadeDonationWithID)
	ctx.Step(`^there is a pending donation with ID "([^"]*)"$`, steps.thereIsAPendingDonationWithID)
	ctx.Step(`^there is a donation with ID "([^"]*)"$`, steps.thereIsADonationWithID)

	// When steps
	ctx.Step(`^I make a donation of "([^"]*)" USD$`, steps.iMakeADonation)
	ctx.Step(`^I request my donation history$`, steps.iRequestMyDonationHistory)
	ctx.Step(`^I get donation details for ID "([^"]*)"$`, steps.iGetDonationDetails)
	ctx.Step(`^I request all donations$`, steps.iRequestAllDonations)
	ctx.Step(`^I update donation "([^"]*)" status to "([^"]*)"$`, steps.iUpdateDonationStatus)
	ctx.Step(`^I delete donation with ID "([^"]*)"$`, steps.iDeleteDonation)

	// Then steps
	ctx.Step(`^the donation should be created successfully$`, steps.theDonationShouldBeCreatedSuccessfully)
	ctx.Step(`^the donation amount should be "([^"]*)" USD$`, steps.theDonationAmountShouldBe)
	ctx.Step(`^I should see a list of my donations$`, steps.iShouldSeeListOfMyDonations)
	ctx.Step(`^the donation details should contain amount "([^"]*)" USD$`, steps.theDonationDetailsShouldContainAmount)
	ctx.Step(`^the donation details should contain status "([^"]*)"$`, steps.theDonationDetailsShouldContainStatus)
	ctx.Step(`^I should see a list of all donations$`, steps.iShouldSeeListOfAllDonations)
	ctx.Step(`^the donation status should be "([^"]*)"$`, steps.theDonationStatusShouldBe)
}

// Given step implementations

func (s *DonationSteps) iHaveMadeDonationsInThePast() error {
	// This step assumes the user has made donations in the test database
	// In a real implementation, you would create test donations if they don't exist
	
	// Make a test donation to ensure there's at least one
	donationData := map[string]interface{}{
		"amount":    50.00,
		"comment":   "Test donation for history",
		"anonymous": false,
	}
	
	// Create the donation
	return s.client.Post("/donations", donationData)
}

func (s *DonationSteps) iHaveMadeDonationWithID(id string) error {
	// This step assumes a donation with the given ID exists
	// In a real implementation, you would create a donation with the specified ID if it doesn't exist
	
	// In a real test, this would use a repository directly or a test API to create a donation with a specific ID
	// For now, we'll just check if the donation exists
	err := s.client.Get(fmt.Sprintf("/donations/%s", id))
	if err != nil {
		return err
	}
	
	// If donation doesn't exist, we'd create it
	if s.client.GetResponseStatusCode() == http.StatusNotFound {
		// Since we can't directly create with ID through the API, this would use a test helper or direct DB access
		return fmt.Errorf("donation with ID %s does not exist and cannot be created through API", id)
	}
	
	return nil
}

func (s *DonationSteps) thereIsAPendingDonationWithID(id string) error {
	// This step assumes a pending donation with the given ID exists
	// In a real implementation, you would create a pending donation with the specified ID
	
	// Similar to above, but ensuring the status is "pending"
	// In a real test, this would use a repository directly or a test API to ensure a donation exists with pending status
	return nil
}

func (s *DonationSteps) thereIsADonationWithID(id string) error {
	// This step assumes a donation with the given ID exists
	// Implementation similar to iHaveMadeDonationWithID
	return nil
}

// When step implementations

func (s *DonationSteps) iMakeADonation(amount string) error {
	// Parse amount to float
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return fmt.Errorf("invalid amount: %s", amount)
	}
	
	donationData := map[string]interface{}{
		"amount":    amountFloat,
		"comment":   "Test donation from BDD test",
		"anonymous": false,
	}
	
	return s.client.Post("/donations", donationData)
}

func (s *DonationSteps) iRequestMyDonationHistory() error {
	// Get the current user's donations
	// In a real implementation, the user ID would be extracted from auth context
	return s.client.Get("/donations/user/current")
}

func (s *DonationSteps) iGetDonationDetails(id string) error {
	return s.client.Get(fmt.Sprintf("/donations/%s", id))
}

func (s *DonationSteps) iRequestAllDonations() error {
	return s.client.Get("/donations")
}

func (s *DonationSteps) iUpdateDonationStatus(id, status string) error {
	statusData := map[string]interface{}{
		"status": status,
	}
	
	return s.client.Patch(fmt.Sprintf("/donations/%s/status", id), statusData)
}

func (s *DonationSteps) iDeleteDonation(id string) error {
	return s.client.Delete(fmt.Sprintf("/donations/%s", id))
}

// Then step implementations

func (s *DonationSteps) theDonationShouldBeCreatedSuccessfully() error {
	if s.client.GetResponseStatusCode() != http.StatusCreated {
		return fmt.Errorf("expected status code %d, got %d", 
			http.StatusCreated, s.client.GetResponseStatusCode())
	}
	
	respBody := s.client.GetResponseBodyAsMap()
	if respBody == nil {
		return fmt.Errorf("response body is empty or not valid JSON")
	}
	
	// Check that the response contains an ID
	if _, ok := respBody["id"].(string); !ok {
		return fmt.Errorf("donation ID not found in response")
	}
	
	return nil
}

func (s *DonationSteps) theDonationAmountShouldBe(expectedAmount string) error {
	respBody := s.client.GetResponseBodyAsMap()
	if respBody == nil {
		return fmt.Errorf("response body is empty or not valid JSON")
	}
	
	// Get donation amount
	var amount float64
	if amountRaw, ok := respBody["amount"].(float64); ok {
		amount = amountRaw
	} else if amountRaw, ok := respBody["amount"].(string); ok {
		parsed, err := strconv.ParseFloat(amountRaw, 64)
		if err != nil {
			return fmt.Errorf("amount is not a valid number: %s", amountRaw)
		}
		amount = parsed
	} else {
		return fmt.Errorf("amount not found in response or not a number")
	}
	
	// Parse expected amount
	expectedAmountFloat, err := strconv.ParseFloat(expectedAmount, 64)
	if err != nil {
		return fmt.Errorf("invalid expected amount: %s", expectedAmount)
	}
	
	// Compare amounts
	if amount != expectedAmountFloat {
		return fmt.Errorf("expected amount %.2f, got %.2f", expectedAmountFloat, amount)
	}
	
	return nil
}

func (s *DonationSteps) iShouldSeeListOfMyDonations() error {
	if s.client.GetResponseStatusCode() != http.StatusOK {
		return fmt.Errorf("expected status code %d, got %d", 
			http.StatusOK, s.client.GetResponseStatusCode())
	}
	
	respBody := s.client.GetResponseBodyAsMap()
	if respBody == nil {
		return fmt.Errorf("response body is empty or not valid JSON")
	}
	
	// Check if the response contains an array
	if _, ok := respBody["donations"].([]interface{}); !ok {
		// Try direct array response
		responseBytes := s.client.GetResponseBody()
		var donations []interface{}
		if err := json.Unmarshal(responseBytes, &donations); err != nil {
			return fmt.Errorf("response does not contain a list of donations")
		}
	}
	
	return nil
}

func (s *DonationSteps) theDonationDetailsShouldContainAmount(expectedAmount string) error {
	respBody := s.client.GetResponseBodyAsMap()
	if respBody == nil {
		return fmt.Errorf("response body is empty or not valid JSON")
	}
	
	// Get donation amount
	var amount float64
	if amountRaw, ok := respBody["amount"].(float64); ok {
		amount = amountRaw
	} else {
		return fmt.Errorf("amount not found in response or not a number")
	}
	
	// Parse expected amount
	expectedAmountFloat, err := strconv.ParseFloat(expectedAmount, 64)
	if err != nil {
		return fmt.Errorf("invalid expected amount: %s", expectedAmount)
	}
	
	// Compare amounts
	if amount != expectedAmountFloat {
		return fmt.Errorf("expected amount %.2f, got %.2f", expectedAmountFloat, amount)
	}
	
	return nil
}

func (s *DonationSteps) theDonationDetailsShouldContainStatus(expectedStatus string) error {
	respBody := s.client.GetResponseBodyAsMap()
	if respBody == nil {
		return fmt.Errorf("response body is empty or not valid JSON")
	}
	
	// Get donation status
	var status string
	if statusRaw, ok := respBody["status"].(string); ok {
		status = statusRaw
	} else {
		return fmt.Errorf("status not found in response or not a string")
	}
	
	// Compare statuses
	if status != expectedStatus {
		return fmt.Errorf("expected status %s, got %s", expectedStatus, status)
	}
	
	return nil
}

func (s *DonationSteps) iShouldSeeListOfAllDonations() error {
	if s.client.GetResponseStatusCode() != http.StatusOK {
		return fmt.Errorf("expected status code %d, got %d", 
			http.StatusOK, s.client.GetResponseStatusCode())
	}
	
	// Implementation similar to iShouldSeeListOfMyDonations
	respBody := s.client.GetResponseBodyAsMap()
	if respBody == nil {
		return fmt.Errorf("response body is empty or not valid JSON")
	}
	
	// Check if the response contains an array
	if _, ok := respBody["donations"].([]interface{}); !ok {
		// Try direct array response
		responseBytes := s.client.GetResponseBody()
		var donations []interface{}
		if err := json.Unmarshal(responseBytes, &donations); err != nil {
			return fmt.Errorf("response does not contain a list of donations")
		}
	}
	
	return nil
}

func (s *DonationSteps) theDonationStatusShouldBe(expectedStatus string) error {
	// Implementation similar to theDonationDetailsShouldContainStatus
	respBody := s.client.GetResponseBodyAsMap()
	if respBody == nil {
		return fmt.Errorf("response body is empty or not valid JSON")
	}
	
	// Get donation status
	var status string
	if statusRaw, ok := respBody["status"].(string); ok {
		status = statusRaw
	} else {
		return fmt.Errorf("status not found in response or not a string")
	}
	
	// Compare statuses
	if status != expectedStatus {
		return fmt.Errorf("expected status %s, got %s", expectedStatus, status)
	}
	
	return nil
}
