package sql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
)

var (
	ErrConnectionNil         = errors.New("connection could not be nil")
	ErrTransactionNotStarted = errors.New("transaction not started")
)

type (
	Connection interface {
		Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error)
		NamedQuery(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
		Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		Begin(ctx context.Context) (context.Context, error)
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
		Close() error
	}
	connection struct {
		db *sqlx.DB
	}
	DatabaseConfig struct {
		DriverName         string
		User               string
		Password           string
		Host               string
		Port               string
		DBName             string
		SSLMode            string
		MaxOpenConnections int
		MaxIdleConnections int
		MaxConnLifetime    time.Duration
	}
)

func newConnection(db *sqlx.DB) (Connection, error) {
	if db == nil {
		return nil, ErrConnectionNil
	}
	return &connection{db: db}, nil
}

func (c *connection) Select(_ context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.db.Select(dest, query, args...)
}

func (c *connection) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	exec := c.db.Exec
	tx := extractTx(ctx)
	if tx != nil {
		exec = tx.Exec
	}
	return exec(query, args...)
}

func (c *connection) NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	namedExec := c.db.NamedExec
	tx := extractTx(ctx)
	if tx != nil {
		namedExec = tx.NamedExec
	}
	return namedExec(query, arg)
}

func (c *connection) NamedQuery(_ context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	return c.db.NamedQuery(query, arg)
}

func (c *connection) Get(_ context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.db.Get(dest, query, args...)
}

func (c *connection) Begin(ctx context.Context) (context.Context, error) {
	tx, err := c.db.Beginx()
	if err != nil {
		return ctx, err
	}
	return injectTx(ctx, tx), nil
}

func (c *connection) Commit(ctx context.Context) error {
	tx := extractTx(ctx)
	if tx == nil {
		return ErrTransactionNotStarted
	}
	return tx.Commit()
}

func (c *connection) Rollback(ctx context.Context) error {
	tx := extractTx(ctx)
	if tx == nil {
		return ErrTransactionNotStarted
	}
	return tx.Rollback()
}

func (c *connection) Close() error {
	return c.db.Close()
}
