Feature: Donations
  As a supporter of the Pet Paradise system
  I want to make donations
  So that I can help animals in need

  Background:
    Given I am authenticated as a "user"

  Scenario: Make a donation
    When I make a donation of "100.00" USD
    Then I should receive a 201 status code
    And the donation should be created successfully
    And the donation amount should be "100.00" USD

  Scenario: View my donation history
    Given I have made donations in the past
    When I request my donation history
    Then I should receive a 200 status code
    And I should see a list of my donations

  Scenario: Get donation details
    Given I have made a donation with ID "donation-123"
    When I get donation details for ID "donation-123"
    Then I should receive a 200 status code
    And the donation details should contain amount "100.00" USD
    And the donation details should contain status "completed"

  Scenario: View all donations (admin only)
    Given I am authenticated as an "admin"
    When I request all donations
    Then I should receive a 200 status code
    And I should see a list of all donations

  Scenario: Update donation status (admin only)
    Given I am authenticated as an "admin"
    And there is a pending donation with ID "donation-123"
    When I update donation "donation-123" status to "completed"
    Then I should receive a 200 status code
    And the donation status should be "completed"

  Scenario: Delete a donation (admin only)
    Given I am authenticated as an "admin"
    And there is a donation with ID "donation-123"
    When I delete donation with ID "donation-123"
    Then I should receive a 204 status code

  Scenario: Regular user cannot view all donations
    Given I am authenticated as a "user"
    When I request all donations
    Then I should receive a 403 status code
