package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/cucumber/godog"
	server "github.com/wellingtonlope/calculator-api/cmd/http/server"
)

func TestMain(m *testing.M) {
	server := server.New()
	go server.Start(":8080")
	code := m.Run()
	server.Close()
	os.Exit(code)
}

type restFeature struct {
	response *http.Response
}

func (f *restFeature) iSendRequestTo(method, endpoint, paramName, paramValue string) error {
	req, err := http.NewRequest(method, fmt.Sprintf("http://localhost:8080%s?%s=%s", endpoint, paramName, paramValue), nil)
	if err != nil {
		return fmt.Errorf("fail to mount request: %w", err)
	}
	f.response, err = http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("fail to make request: %w", err)
	}
	return nil
}

func (f *restFeature) theResponseCodeShouldBe(code int) error {
	if f.response.StatusCode != code {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, f.response.StatusCode)
	}
	return nil
}

func (f *restFeature) theResponseShouldMatchJSON(body *godog.DocString) (err error) {
	var expected, actual interface{}
	if err = json.Unmarshal([]byte(body.Content), &expected); err != nil {
		return
	}
	bodyResponse, err := io.ReadAll(f.response.Body)
	if err != nil {
		return
	}
	if err = json.Unmarshal(bodyResponse, &actual); err != nil {
		return
	}
	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expected, actual)
	}
	return nil
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	rest := &restFeature{}

	ctx.Step(`^I send a "(GET|POST|PUT|DELETE)" request for "([^"]*)" with the query params "([^"]*)" filled with the values "([^"]*)"$`, rest.iSendRequestTo)
	ctx.Step(`^the response should have status code (\d+)$`, rest.theResponseCodeShouldBe)
	ctx.Step(`^the response should match json:$`, rest.theResponseShouldMatchJSON)
}
