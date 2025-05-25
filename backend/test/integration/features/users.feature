Feature: User Management
  As a user of the Pet Paradise system
  I want to manage user accounts
  So that I can perform CRUD operations on users

  Background:
    Given the system is initialized

  Scenario: Register a new user
    Given I have valid registration data
    When I register a new user
    Then I should receive a 201 status code
    And the response should contain user details
    And the user should be created in the database

  Scenario: Register a user with invalid data
    Given I have invalid registration data
    When I register a new user
    Then I should receive a 400 status code
    And the response should contain an error message

  Scenario: Register a user with existing email
    Given I have registration data with an existing email
    When I register a new user
    Then I should receive a 400 status code
    And the response should contain "Email already exists"

  Scenario: Get all users as admin
    Given I am authenticated as an "admin"
    When I request all users
    Then I should receive a 200 status code
    And the response should contain a list of users

  Scenario: Get all users as regular user
    Given I am authenticated as a "user"
    When I request all users
    Then I should receive a 200 status code
    And the response should contain a list of users

  Scenario: Get user by ID as admin
    Given I am authenticated as an "admin"
    And a user exists in the system
    When I request the user by ID
    Then I should receive a 200 status code
    And the response should contain the user details

  Scenario: Get user by ID as same user
    Given I am authenticated as a "user"
    When I request my own user details by ID
    Then I should receive a 200 status code
    And the response should contain my user details

  Scenario: Get user by ID as different user
    Given I am authenticated as a "user"
    And another user exists in the system
    When I request the other user by ID
    Then I should receive a 403 status code

  Scenario: Get user by email as admin
    Given I am authenticated as an "admin"
    And a user exists in the system
    When I request the user by email
    Then I should receive a 200 status code
    And the response should contain the user details

  Scenario: Get users by status as admin
    Given I am authenticated as an "admin"
    And users with different statuses exist
    When I request users by status "active"
    Then I should receive a 200 status code
    And the response should contain only active users

  Scenario: Update user information as admin
    Given I am authenticated as an "admin"
    And a user exists in the system
    When I update the user information
    Then I should receive a 200 status code
    And the user information should be updated

  Scenario: Update own user information
    Given I am authenticated as a "user"
    When I update my own user information
    Then I should receive a 200 status code
    And my user information should be updated

  Scenario: Update other user information as regular user
    Given I am authenticated as a "user"
    And another user exists in the system
    When I try to update the other user information
    Then I should receive a 403 status code

  Scenario: Update user role as admin
    Given I am authenticated as an "admin"
    And a user exists in the system
    When I update the user role to "admin"
    Then I should receive a 200 status code
    And the user role should be updated

  Scenario: Update user role as regular user
    Given I am authenticated as a "user"
    And another user exists in the system
    When I try to update the user role to "admin"
    Then I should receive a 403 status code

  Scenario: Update user status as admin
    Given I am authenticated as an "admin"
    And a user exists in the system
    When I update the user status to "inactive"
    Then I should receive a 200 status code
    And the user status should be updated

  Scenario: Update user status as regular user
    Given I am authenticated as a "user"
    And another user exists in the system
    When I try to update the user status to "inactive"
    Then I should receive a 403 status code

  Scenario: Change user password
    Given I am authenticated as a "user"
    When I change my password
    Then I should receive a 200 status code
    And I should be able to login with the new password

  Scenario: Change user password with wrong current password
    Given I am authenticated as a "user"
    When I try to change my password with wrong current password
    Then I should receive a 400 status code
    And the response should contain "Invalid current password"

  Scenario: Delete user as admin
    Given I am authenticated as an "admin"
    And a user exists in the system
    When I delete the user
    Then I should receive a 200 status code
    And the user should be removed from the system

  Scenario: Delete user as regular user
    Given I am authenticated as a "user"
    And another user exists in the system
    When I try to delete the other user
    Then I should receive a 403 status code

  Scenario: Access user endpoints without authentication
    Given I am not authenticated
    When I try to access user endpoints without authentication
    Then I should receive a 401 status code

  Scenario: Revoke all user tokens
    Given I am authenticated as a "user"
    When I revoke all my tokens
    Then I should receive a 200 status code
    And all my tokens should be invalidated
