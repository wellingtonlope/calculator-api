package usecase

import (
	"context"

	"github.com/wellingtonlope/calculator-api/internal/app/repository"
	"github.com/wellingtonlope/calculator-api/internal/domain"
)

type (
	CreateVariableInput struct {
		Name  string
		Value float64
	}
	CreateVariable interface {
		Handle(ctx context.Context, input CreateVariableInput) error
	}
	createVariable struct {
		variableRepository repository.Variable
	}
)

func NewCreateVariable(variableRepository repository.Variable) CreateVariable {
	return createVariable{variableRepository: variableRepository}
}

func (uc createVariable) Handle(ctx context.Context, input CreateVariableInput) error {
	variable, err := domain.NewVariable(input.Name, input.Value)
	if err != nil {
		return NewError(err.Error(), err, ErrorTypeInvalid)
	}
	err = uc.variableRepository.Create(ctx, variable)
	if err != nil {
		return NewError("error creating variable", err, ErrorTypeUnknown)
	}
	return nil
}
