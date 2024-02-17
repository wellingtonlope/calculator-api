package sql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/wellingtonlope/calculator-api/internal/app/repository"
	"github.com/wellingtonlope/calculator-api/internal/domain"
	sqlpkg "github.com/wellingtonlope/calculator-api/pkg/sql"
)

type variable struct {
	db sqlpkg.Connection
}

func NewVariable(db sqlpkg.Connection) repository.Variable {
	return variable{db: db}
}

func (r variable) Create(ctx context.Context, input domain.Variable) error {
	_, err := r.db.NamedExec(ctx, `
    INSERT INTO variables (name, value) VALUES (:name, :value);
  `, struct {
		Name  string  `db:"name"`
		Value float64 `db:"value"`
	}(input))
	return err
}

func (r variable) GetByName(ctx context.Context, name string) (domain.Variable, error) {
	var variableDTO struct {
		Name  string  `db:"name"`
		Value float64 `db:"value"`
	}
	err := r.db.Get(ctx, &variableDTO, `
    SELECT name, value FROM variables WHERE name = ?
  `, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Variable{}, repository.ErrVariableNotFound
		}
		return domain.Variable{}, err
	}
	return domain.Variable(variableDTO), nil
}
