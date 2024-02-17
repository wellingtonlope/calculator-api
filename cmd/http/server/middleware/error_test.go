package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/wellingtonlope/calculator-api/internal/app/usecase"
)

func TestError(t *testing.T) {
	testCases := []struct {
		name           string
		handler        echo.HandlerFunc
		err            error
		responseBody   string
		responseStatus int
	}{
		{
			name: "should do nothing when error is nil",
			handler: func(c echo.Context) error {
				return nil
			},
			responseBody:   ``,
			responseStatus: http.StatusOK,
		},
		{
			name: "should return internal server error when isn't usecase.Error",
			handler: func(c echo.Context) error {
				return assert.AnError
			},
			err:            assert.AnError,
			responseBody:   `{"message":"internal error"}`,
			responseStatus: http.StatusInternalServerError,
		},
		{
			name: "should return correct status and message when is usecase.Error",
			handler: func(c echo.Context) error {
				return usecase.NewError("invalid", assert.AnError, usecase.ErrorTypeInvalid)
			},
			err:            usecase.NewError("invalid", assert.AnError, usecase.ErrorTypeInvalid),
			responseBody:   `{"message":"invalid"}`,
			responseStatus: http.StatusBadRequest,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest("GET", "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := Error(tc.handler)(c)

			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.responseStatus, rec.Code)
			assert.Equal(t, tc.responseBody, strings.ReplaceAll(rec.Body.String(), "\n", ""))
		})
	}
}
