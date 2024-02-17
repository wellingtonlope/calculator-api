package repository

import (
	"context"
	"errors"

	"github.com/wellingtonlope/calculator-api/internal/domain"
)

var ErrVariableNotFound = errors.New("variable not found")

type Variable interface {
	Create(context.Context, domain.Variable) error
	GetByName(context.Context, string) (domain.Variable, error)
}
