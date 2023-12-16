package usecase

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type SumNumbersMock struct {
	mock.Mock
}

func NewSumNumbersMock() *SumNumbersMock {
	return new(SumNumbersMock)
}

func (m *SumNumbersMock) Handle(ctx context.Context, numbers []float64) float64 {
	args := m.Called(ctx, numbers)
	var result float64
	if _, ok := args.Get(0).(float64); ok {
		result = args.Get(0).(float64)
	}
	return result
}
