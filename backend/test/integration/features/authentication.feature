Feature: Authentication
  As a user of the Pet Paradise system
  I want to be able to authenticate securely
  So that I can access protected resources

  Scenario: User login with valid credentials
    Given I have valid user credentials
    When I login with my credentials
    Then I should receive a valid token pair
    And I should receive a 200 status code

  Scenario: User login with invalid credentials
    Given I have invalid user credentials
    When I login with my credentials
    Then I should receive an authentication error
    And I should receive a 401 status code

  Scenario: Refresh access token with valid refresh token
    Given I am authenticated as a "user"
    And I have a valid refresh token
    When I request to refresh my tokens
    Then I should receive new valid tokens
    And I should receive a 200 status code

  Scenario: Refresh access token with expired refresh token
    Given I have an expired refresh token
    When I request to refresh my tokens
    Then I should receive an authentication error
    And I should receive a 401 status code
    And the response should contain "Invalid refresh token"

  Scenario: Access protected resource with valid token
    Given I am authenticated as a "user"
    When I use my token to access a protected resource
    Then I should receive a 200 status code

  Scenario: Access protected resource without authentication
    Given I am not authenticated
    When I try to access a protected resource without authentication
    Then I should receive an authentication error
    And I should receive a 401 status code


