package usecase

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSumNumbers_Handle(t *testing.T) {
	testCases := []struct {
		name    string
		ctx     context.Context
		numbers []float64
		result  float64
	}{
		{
			name:    "should return 0 when slice is empty",
			ctx:     context.TODO(),
			numbers: []float64{},
			result:  0,
		},
		{
			name:    "should return the sum when slice is not empty",
			ctx:     context.TODO(),
			numbers: []float64{1.1, 2},
			result:  3.1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc := NewSumNumbers()
			result := uc.Handle(tc.ctx, tc.numbers)
			assert.Equal(t, tc.result, result)
		})
	}
}
