package repository

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/wellingtonlope/calculator-api/internal/domain"
)

type VariableMock struct {
	mock.Mock
}

func NewVariableMock() *VariableMock {
	return new(VariableMock)
}

func (m *VariableMock) Create(ctx context.Context, variable domain.Variable) error {
	args := m.Called(ctx, variable)
	return args.Error(0)
}

func (m *VariableMock) GetByName(ctx context.Context, name string) (domain.Variable, error) {
	args := m.Called(ctx, name)
	var result domain.Variable
	if _, ok := args.Get(0).(domain.Variable); ok {
		result = args.Get(0).(domain.Variable)
	}
	return result, args.Error(1)
}
