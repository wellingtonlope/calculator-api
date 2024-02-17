package sql

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/wellingtonlope/calculator-api/internal/app/repository"
	"github.com/wellingtonlope/calculator-api/internal/domain"
	sqlpkg "github.com/wellingtonlope/calculator-api/pkg/sql"
)

func TestVariable_Create(t *testing.T) {
	exampleContext := context.TODO()
	type args struct {
		ctx   context.Context
		input domain.Variable
	}
	testCases := []struct {
		name string
		db   sqlpkg.Connection
		args args
		err  error
	}{
		{
			name: "should fail when db fails",
			db: func() sqlpkg.Connection {
				sqlMock, db := sqlpkg.NewConnectionMock()
				sqlMock.ExpectExec(`INSERT INTO variables \(name\, value\)`).
					WithArgs("PI", float64(3.14)).
					WillReturnError(assert.AnError)
				return db
			}(),
			args: args{
				ctx:   exampleContext,
				input: domain.Variable{Name: "PI", Value: 3.14},
			},
			err: assert.AnError,
		},
		{
			name: "should create variable",
			db: func() sqlpkg.Connection {
				sqlMock, db := sqlpkg.NewConnectionMock()
				sqlMock.ExpectExec(`INSERT INTO variables \(name\, value\)`).
					WithArgs("PI", float64(3.14)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				return db
			}(),
			args: args{
				ctx:   exampleContext,
				input: domain.Variable{Name: "PI", Value: 3.14},
			},
			err: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewVariable(tc.db)
			err := repo.Create(exampleContext, tc.args.input)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestVariable_GetByName(t *testing.T) {
	exampleContext := context.TODO()
	type args struct {
		ctx  context.Context
		name string
	}
	testCases := []struct {
		name   string
		db     sqlpkg.Connection
		args   args
		result domain.Variable
		err    error
	}{
		{
			name: "should fail when db fails",
			db: func() sqlpkg.Connection {
				sqlMock, db := sqlpkg.NewConnectionMock()
				sqlMock.ExpectQuery(`SELECT name, value FROM variables WHERE name = ?`).
					WithArgs("PI").
					WillReturnError(assert.AnError)
				return db
			}(),
			args: args{
				ctx:  exampleContext,
				name: "PI",
			},
			result: domain.Variable{},
			err:    assert.AnError,
		},
		{
			name: "should fail when variable not found",
			db: func() sqlpkg.Connection {
				sqlMock, db := sqlpkg.NewConnectionMock()
				sqlMock.ExpectQuery(`SELECT name, value FROM variables WHERE name = ?`).
					WithArgs("PI").
					WillReturnError(sql.ErrNoRows)
				return db
			}(),
			args: args{
				ctx:  exampleContext,
				name: "PI",
			},
			result: domain.Variable{},
			err:    repository.ErrVariableNotFound,
		},
		{
			name: "should get variable by name",
			db: func() sqlpkg.Connection {
				sqlMock, db := sqlpkg.NewConnectionMock()
				sqlMock.ExpectQuery(`SELECT name, value FROM variables WHERE name = ?`).
					WithArgs("PI").
					WillReturnRows(
						sqlmock.NewRows([]string{"name", "value"}).
							AddRow("PI", 3.14),
					)
				return db
			}(),
			args: args{
				ctx:  exampleContext,
				name: "PI",
			},
			result: domain.Variable{Name: "PI", Value: 3.14},
			err:    nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewVariable(tc.db)
			result, err := repo.GetByName(exampleContext, tc.args.name)
			assert.Equal(t, tc.result, result)
			assert.Equal(t, tc.err, err)
		})
	}
}
