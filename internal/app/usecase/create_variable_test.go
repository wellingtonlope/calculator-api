package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wellingtonlope/calculator-api/internal/app/repository"
	"github.com/wellingtonlope/calculator-api/internal/domain"
)

func TestCreateVariable_Handle(t *testing.T) {
	exampleContext := context.TODO()
	type args struct {
		ctx   context.Context
		input CreateVariableInput
	}
	testCases := []struct {
		name               string
		variableRepository *repository.VariableMock
		args               args
		err                error
	}{
		{
			name:               "should fail when input is invalid",
			variableRepository: repository.NewVariableMock(),
			args: args{
				ctx:   exampleContext,
				input: CreateVariableInput{},
			},
			err: NewError(
				fmt.Errorf("%w: name", domain.ErrVariableInvalidInput).Error(),
				fmt.Errorf("%w: name", domain.ErrVariableInvalidInput),
				ErrorTypeInvalid,
			),
		},
		{
			name: "should fail when repository fails",
			variableRepository: func() *repository.VariableMock {
				m := repository.NewVariableMock()
				m.On("Create", exampleContext, domain.Variable{Name: "PI", Value: 3.14}).
					Return(assert.AnError).Once()
				return m
			}(),
			args: args{
				ctx: exampleContext,
				input: CreateVariableInput{
					Name:  "PI",
					Value: 3.14,
				},
			},
			err: NewError(
				"error creating variable",
				assert.AnError,
				ErrorTypeUnknown,
			),
		},
		{
			name: "should create variable",
			variableRepository: func() *repository.VariableMock {
				m := repository.NewVariableMock()
				m.On("Create", exampleContext, domain.Variable{Name: "PI", Value: 3.14}).
					Return(nil).Once()
				return m
			}(),
			args: args{
				ctx: exampleContext,
				input: CreateVariableInput{
					Name:  "PI",
					Value: 3.14,
				},
			},
			err: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc := NewCreateVariable(tc.variableRepository)
			err := uc.Handle(tc.args.ctx, tc.args.input)
			assert.Equal(t, tc.err, err)
			tc.variableRepository.AssertExpectations(t)
		})
	}
}
