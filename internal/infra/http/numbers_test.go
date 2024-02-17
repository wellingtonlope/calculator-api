package http

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wellingtonlope/calculator-api/internal/app/usecase"
)

func TestNumbers_Sum(t *testing.T) {
	testCases := []struct {
		name           string
		sum            *usecase.SumNumbersMock
		url            string
		err            error
		responseBody   string
		responseStatus int
	}{
		{
			name: "should fail when param numbers are not numbers",
			sum:  usecase.NewSumNumbersMock(),
			url:  "/sum?numbers=a,b",
			err: usecase.NewError("numbers values must be numbers", func() error {
				_, err := strconv.ParseFloat("a", 64)
				return err
			}(), usecase.ErrorTypeInvalid),
			responseBody:   `{"message":"numbers values must be numbers"}`,
			responseStatus: http.StatusBadRequest,
		},
		{
			name:           "should return 0 when numbers params is empty",
			sum:            usecase.NewSumNumbersMock(),
			url:            "/sum?numbers=",
			err:            nil,
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
			err:            nil,
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

			if tc.err != nil {
				assert.Equal(t, tc.err, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tc.responseStatus, rec.Code)
			assert.Equal(t, tc.responseBody, strings.ReplaceAll(rec.Body.String(), "\n", ""))
		})
	}
}
