# file: sum.feature
Feature: sum feature
  API endpoint to sum numbers

  Scenario: should sum 2 numbers
    When I send a "GET" request for "/sum" with the query params "numbers" filled with the values "1,2"
    Then the response should have status code 200
    And the response should match json:
      """
      {
        "result": 3
      }
      """
  Scenario: should sum 10 numbers
    When I send a "GET" request for "/sum" with the query params "numbers" filled with the values "1,2,3,4,5,6,7,8,9,10"
    Then the response should have status code 200
    And the response should match json:
      """
      {
        "result": 55
      }
      """
  Scenario: should return 0
    When I send a "GET" request for "/sum" with the query params "numbers" filled with the values ""
    Then the response should have status code 200
    And the response should match json:
      """
      {
        "result": 0
      }
      """
  Scenario: should return error when params is not a number
    When I send a "GET" request for "/sum" with the query params "numbers" filled with the values "a,b"
    Then the response should have status code 400
    And the response should match json:
      """
      {
        "message": "numbers values must be numbers"
      }
      """