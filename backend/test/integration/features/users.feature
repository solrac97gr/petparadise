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
    Then I should receive a 409 status code
    And the response should contain "email already in use"

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

  Scenario: Delete user as admin
    Given I am authenticated as an "admin"
    And a user exists in the system
    When I delete the user
    Then I should receive a 204 status code
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
