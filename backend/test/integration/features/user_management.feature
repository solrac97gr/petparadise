Feature: User Management
  As an administrator of the Pet Paradise system
  I want to manage user accounts
  So that I can control who has access to the system

  Scenario: Create a new user
    Given I am not authenticated
    When I create a new user with name "John Doe", email "john@example.com", password "password123"
    Then the user should be created successfully
    And I should see user details with email "john@example.com"
    And I should see user role is "user"
    And I should see user status is "active"

  Scenario: Get user details by ID
    Given I am authenticated as an "admin"
    When I get user details for user with ID "user-123"
    Then I should receive a 200 status code
    And I should see user details with email "test@example.com"

  Scenario: Get user details by email
    Given I am authenticated as an "admin"
    When I get user details for email "test@example.com"
    Then I should receive a 200 status code
    And I should see user details with email "test@example.com"

  Scenario: Update user information
    Given I am authenticated as an "admin"
    When I update user "user-123" with name "Updated Name"
    Then the user should be updated successfully
    And I should see user details with email "test@example.com"

  Scenario: Update user role (admin only)
    Given I am authenticated as an "admin"
    When I update user "user-123" role to "volunteer"
    Then the user should be updated successfully
    And I should see user role is "volunteer"

  Scenario: Update user status (admin only)
    Given I am authenticated as an "admin"
    When I update user "user-123" status to "inactive"
    Then the user should be updated successfully
    And I should see user status is "inactive"

  Scenario: Change user password
    Given I am authenticated as a "user"
    When I change password for user "user-123" from "oldpassword" to "newpassword"
    Then I should receive a 200 status code
    And the response should contain "Password changed successfully"

  Scenario: Delete user (admin only)
    Given I am authenticated as an "admin"
    When I delete user "user-123"
    Then the user should be deleted successfully

  Scenario: List all users (admin only)
    Given I am authenticated as an "admin"
    When I get user details for user with ID ""
    Then I should receive a 200 status code
    And I should see a list of users

  Scenario: Regular user cannot access admin endpoints
    Given I am authenticated as a "user"
    When I get user details for user with ID "another-user-id"
    Then I should receive a 403 status code
