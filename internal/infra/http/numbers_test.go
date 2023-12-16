package http

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wellingtonlope/calculator-api/internal/app/usecase"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNumbers_Sum(t *testing.T) {
	testCases := []struct {
		name           string
		sum            *usecase.SumNumbersMock
		url            string
		responseBody   string
		responseStatus int
	}{
		{
			name:           "should fail when param numbers are not numbers",
			sum:            usecase.NewSumNumbersMock(),
			url:            "/sum?numbers=a,b",
			responseBody:   `{"message":"numbers values must be numbers"}`,
			responseStatus: http.StatusBadRequest,
		},
		{
			name:           "should return 0 when numbers params is empty",
			sum:            usecase.NewSumNumbersMock(),
			url:            "/sum?numbers=",
			responseBody:   `{"result":0}`,
			responseStatus: http.StatusOK,
		},
		{
			name: "should sum number successfully",
			sum: func() *usecase.SumNumbersMock {
				m := usecase.NewSumNumbersMock()
				m.On("Handle", mock.Anything, []float64{1, 2}).Return(3.0).Once()
				return m
			}(),
			url:            "/sum?numbers=1,2",
			responseBody:   `{"result":3}`,
			responseStatus: http.StatusOK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest("GET", tc.url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			handle := NewNumbers(tc.sum)

			err := handle.Sum(c)

			assert.Nil(t, err)
			assert.Equal(t, tc.responseStatus, rec.Code)
			assert.Equal(t, tc.responseBody, strings.ReplaceAll(rec.Body.String(), "\n", ""))
		})
	}
}
