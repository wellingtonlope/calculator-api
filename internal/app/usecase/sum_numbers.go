package usecase

import "context"

type (
	SumNumbers interface {
		Handle(ctx context.Context, numbers []float64) float64
	}
	sumNumbers struct{}
)

func NewSumNumbers() SumNumbers {
	return &sumNumbers{}
}

func (s sumNumbers) Handle(_ context.Context, numbers []float64) (sum float64) {
	for _, number := range numbers {
		sum += number
	}
	return sum
}
