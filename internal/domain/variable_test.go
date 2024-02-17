package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVariable(t *testing.T) {
	type args struct {
		name  string
		value float64
	}
	testCases := []struct {
		name   string
		args   args
		result Variable
		err    error
	}{
		{
			name: "should fail when name is empty",
			args: args{
				name:  "",
				value: 3.14,
			},
			result: Variable{},
			err:    fmt.Errorf("%w: name", ErrVariableInvalidInput),
		},
		{
			name: "should create variable",
			args: args{
				name:  "pi",
				value: 3.14,
			},
			result: Variable{Name: "PI", Value: 3.14},
			err:    nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := NewVariable(tc.args.name, tc.args.value)
			assert.Equal(t, tc.result, result)
			assert.Equal(t, tc.err, err)
		})
	}
}
