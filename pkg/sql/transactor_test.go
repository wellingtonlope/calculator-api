package sql

import (
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestNewTransactor(t *testing.T) {
	t.Run("should fail when sqlx.DB is nil", func(t *testing.T) {
		service, err := newTransactor(nil)

		assert.Nil(t, service)
		assert.Equal(t, ErrConnectionNil, err)
	})

	t.Run("should create transactor", func(t *testing.T) {
		_, _, sqlDB := newConnectionMockWithDB()

		service, err := newTransactor(sqlDB)

		assert.NotEmpty(t, service)
		assert.Nil(t, err)
	})
}

func TestInjectTx(t *testing.T) {
	sqlMock, _, sqlDB := newConnectionMockWithDB()
	sqlMock.ExpectBegin()
	expectedTX, _ := sqlDB.Beginx()
	expectedContext := context.WithValue(context.TODO(), keyTx{}, expectedTX)
	inputContext := context.TODO()

	result := injectTx(inputContext, expectedTX)

	assert.Equal(t, expectedContext, result)
}

func TestExtractTx(t *testing.T) {
	t.Run("should return transaction", func(t *testing.T) {
		sqlMock, _, sqlDB := newConnectionMockWithDB()
		sqlMock.ExpectBegin()
		expectedTX, _ := sqlDB.Beginx()
		expectedContext := context.WithValue(context.TODO(), keyTx{}, expectedTX)

		result := extractTx(expectedContext)

		assert.Equal(t, expectedTX, result)
	})

	t.Run("should return nil when transaction not exists", func(t *testing.T) {
		expectedContext := context.TODO()

		result := extractTx(expectedContext)

		assert.Nil(t, result)
	})
}

func TestTransactor_WithinTransaction(t *testing.T) {
	testCases := []struct {
		name      string
		write     func() *sqlx.DB
		input     context.Context
		funcInput func(ctx context.Context) error
		err       error
	}{
		{
			name: "should return error when it fails to create transaction",
			write: func() *sqlx.DB {
				sqlMock, _, sqlDB := newConnectionMockWithDB()
				sqlMock.ExpectBegin().WillReturnError(errors.New("i'm an error"))
				return sqlDB
			},
			input: context.TODO(),
			funcInput: func(ctx context.Context) error {
				return nil
			},
			err: fmt.Errorf("sql.transactor begin fails: %w", errors.New("i'm an error")),
		},
		{
			name: "should return error when it fails to rollback transaction",
			write: func() *sqlx.DB {
				sqlMock, _, sqlDB := newConnectionMockWithDB()
				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback().WillReturnError(errors.New("i'm an error"))
				return sqlDB
			},
			input: context.TODO(),
			funcInput: func(ctx context.Context) error {
				return errors.New("i'm an error")
			},
			err: fmt.Errorf("sql.transactor rollback fails: %w", errors.New("i'm an error")),
		},
		{
			name: "should execute rollback without error",
			write: func() *sqlx.DB {
				sqlMock, _, sqlDB := newConnectionMockWithDB()
				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()
				return sqlDB
			},
			input: context.TODO(),
			funcInput: func(ctx context.Context) error {
				return errors.New("i'm an error")
			},
			err: errors.New("i'm an error"),
		},
		{
			name: "should execute commit without error",
			write: func() *sqlx.DB {
				sqlMock, _, sqlDB := newConnectionMockWithDB()
				sqlMock.ExpectBegin()
				sqlMock.ExpectCommit()
				return sqlDB
			},
			input: context.TODO(),
			funcInput: func(ctx context.Context) error {
				return nil
			},
			err: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service, _ := newTransactor(tc.write())

			err := service.WithinTransaction(tc.input, tc.funcInput)

			assert.Equal(t, tc.err, err)
		})
	}
}

func TestTransactor_WithinTransaction_After(t *testing.T) {
	t.Run("should be possible execute query after Transaction", func(t *testing.T) {
		expectedContext := context.TODO()
		sqlMock, conn, sqlDB := newConnectionMockWithDB()
		sqlMock.ExpectBegin()
		sqlMock.ExpectCommit()
		sqlMock.ExpectExec("SELECT").WillReturnResult(sqlmock.NewResult(1, 1))
		service, _ := newTransactor(sqlDB)

		_ = service.WithinTransaction(expectedContext, func(ctx context.Context) error {
			return nil
		})

		_, err := conn.Exec(expectedContext, "SELECT 1;")
		assert.Nil(t, err)
	})
}
