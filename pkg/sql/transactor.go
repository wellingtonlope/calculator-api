package sql

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// Transactor runs logic inside a single database transaction
type Transactor interface {
	// WithinTransaction runs function within transaction
	//
	// The transaction commits when function were finished without error
	WithinTransaction(context.Context, func(ctx context.Context) error) error
}

type transactor struct {
	db *sqlx.DB
}

func newTransactor(db *sqlx.DB) (Transactor, error) {
	if db == nil {
		return nil, ErrConnectionNil
	}
	return &transactor{
		db: db,
	}, nil
}

func (t *transactor) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	defer purgeTx(ctx)
	tx, err := t.db.Beginx()
	if err != nil {
		return fmt.Errorf("sql.transactor begin fails: %w", err)
	}

	errApi := tFunc(injectTx(ctx, tx))
	if errApi != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return fmt.Errorf("sql.transactor rollback fails: %w", errRollback)
		}
		return errApi
	}

	if errCommit := tx.Commit(); errCommit != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return fmt.Errorf("sql.transactor commit and rollback fails: %w", errRollback)
		}
		return fmt.Errorf("sql.transactor rollback done after commit fails: %w", errCommit)
	}
	return nil
}

type keyTx struct{}

// injectTx injects transaction to context
func injectTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, keyTx{}, tx)
}

// extractTx extracts transaction from context
func extractTx(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(keyTx{}).(*sqlx.Tx); ok {
		return tx
	}
	return nil
}

// purgeTx purge transaction from context
func purgeTx(ctx context.Context) context.Context {
	return context.WithValue(ctx, keyTx{}, nil)
}
