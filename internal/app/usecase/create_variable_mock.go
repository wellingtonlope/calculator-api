package usecase

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type CreateVariableMock struct {
	mock.Mock
}

func NewCreateVariableMock() *CreateVariableMock {
	return new(CreateVariableMock)
}

func (m *CreateVariableMock) Handle(ctx context.Context, input CreateVariableInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}
