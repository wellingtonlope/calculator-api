package domain

import (
	"errors"
	"fmt"
	"strings"
)

var ErrVariableInvalidInput = errors.New("invalid input")

type (
	Variable struct {
		Name  string
		Value float64
	}
)

func NewVariable(name string, value float64) (Variable, error) {
	name = strings.ToUpper(strings.Trim(name, " "))
	if name == "" {
		return Variable{}, fmt.Errorf("%w: name", ErrVariableInvalidInput)
	}
	return Variable{
		Name:  name,
		Value: value,
	}, nil
}
