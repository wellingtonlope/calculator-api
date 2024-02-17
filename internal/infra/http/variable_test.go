package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wellingtonlope/calculator-api/internal/app/usecase"
)

func TestVariable_Create(t *testing.T) {
	testCases := []struct {
		name           string
		createVariable *usecase.CreateVariableMock
		requestBody    string
		responseStatus int
		err            error
	}{
		{
			name:           "should fail when input is invalid",
			createVariable: usecase.NewCreateVariableMock(),
			requestBody:    `{`,
			err: usecase.NewError("invalid input JSON", func() error {
				var input struct {
					Name  string  `json:"name"`
					Value float64 `json:"value"`
				}
				e := echo.New()
				req := httptest.NewRequest("POST", "/variable", strings.NewReader(`{`))
				req.Header.Add("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				return e.Binder.Bind(&input, c)
			}(), usecase.ErrorTypeInvalid),
		},
		{
			name: "should fail when usecase fails",
			createVariable: func() *usecase.CreateVariableMock {
				m := usecase.NewCreateVariableMock()
				m.On("Handle", mock.Anything, usecase.CreateVariableInput{Name: "PI", Value: 3.14}).
					Return(assert.AnError).Once()
				return m
			}(),
			requestBody: `{"name":"PI","value":3.14}`,
			err:         assert.AnError,
		},
		{
			name: "should create variable",
			createVariable: func() *usecase.CreateVariableMock {
				m := usecase.NewCreateVariableMock()
				m.On("Handle", mock.Anything, usecase.CreateVariableInput{Name: "PI", Value: 3.14}).
					Return(nil).Once()
				return m
			}(),
			requestBody:    `{"name":"PI","value":3.14}`,
			responseStatus: http.StatusCreated,
			err:            nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest("POST", "/variable", strings.NewReader(tc.requestBody))
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			handle := NewVariable(tc.createVariable)

			err := handle.Create(c)

			if tc.err != nil {
				assert.Equal(t, tc.err, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tc.responseStatus, rec.Code)
		})
	}
}
