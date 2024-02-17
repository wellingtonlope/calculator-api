# file: create variable.feature
Feature: create variable feature
  API endpoint to create variable

  Scenario: should create variable
    When I send a "POST" request for "/variable" with json:
      """
      {
        "name": "PI",
        "value": 3.14
      }
      """
    Then the response should have status code 201
  Scenario: should fail when name is invalid
    When I send a "POST" request for "/variable" with json:
      """
      {
        "name": "",
        "value": 3.14
      }
      """
    Then the response should have status code 400
    And the response should match json:
      """
      {
        "message": "invalid input: name"
      }
      """

