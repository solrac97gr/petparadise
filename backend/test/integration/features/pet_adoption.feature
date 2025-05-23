Feature: Pet Adoption
  As a user of the Pet Paradise system
  I want to adopt pets
  So that I can give them a loving home

  Background:
    Given I am authenticated as a "user"
    And there is an available pet with ID "pet-123"

  Scenario: View available pets
    When I request the list of available pets
    Then I should receive a 200 status code
    And I should see a list of pets
    And all pets should have status "available"

  Scenario: View pet details
    When I get pet details for pet with ID "pet-123"
    Then I should receive a 200 status code
    And the pet details should contain name "Fluffy"
    And the pet details should contain species "Dog"

  Scenario: Create an adoption request
    When I create an adoption request for pet "pet-123"
    Then I should receive a 201 status code
    And the adoption request should be created successfully
    And the pet status should be "pending"

  Scenario: View my adoption requests
    Given I have submitted an adoption request
    When I request my adoption requests
    Then I should receive a 200 status code
    And I should see a list of my adoption requests

  Scenario: Cancel an adoption request
    Given I have submitted an adoption request with ID "adoption-123"
    When I cancel the adoption request with ID "adoption-123"
    Then I should receive a 200 status code
    And the adoption request status should be "cancelled"
    And the pet status should be "available"

  Scenario: Approve an adoption request (staff only)
    Given I am authenticated as a "volunteer"
    And there is a pending adoption request with ID "adoption-123"
    When I approve the adoption request with ID "adoption-123"
    Then I should receive a 200 status code
    And the adoption request status should be "approved"
    And the pet status should be "adopted"

  Scenario: Reject an adoption request (staff only)
    Given I am authenticated as a "volunteer"
    And there is a pending adoption request with ID "adoption-123"
    When I reject the adoption request with ID "adoption-123"
    Then I should receive a 200 status code
    And the adoption request status should be "rejected"
    And the pet status should be "available"
