package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Error(t *testing.T) {
	testCases := []struct {
		name   string
		error  Error
		result string
	}{
		{
			name: "should get cause error",
			error: Error{
				Message: "message",
				Cause:   assert.AnError,
				Type:    ErrorTypeUnknown,
			},
			result: assert.AnError.Error(),
		},
		{
			name: "should get message error",
			error: Error{
				Message: "message",
				Type:    ErrorTypeUnknown,
			},
			result: "message",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.error.Error()
			assert.Equal(t, tc.result, result)
		})
	}
}

func TestNewError(t *testing.T) {
	type args struct {
		message   string
		cause     error
		errorType ErrorType
	}
	testCases := []struct {
		name   string
		args   args
		result Error
	}{
		{
			name: "should create error with message and cause",
			args: args{
				message:   "message",
				cause:     assert.AnError,
				errorType: ErrorTypeUnknown,
			},
			result: Error{
				Message: "message",
				Cause:   assert.AnError,
				Type:    ErrorTypeUnknown,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := NewError(tc.args.message, tc.args.cause, tc.args.errorType)
			assert.Equal(t, tc.result, result)
		})
	}
}
