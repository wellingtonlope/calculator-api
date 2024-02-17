package sql

import (
	"context"
)

type transactorMock struct{}

func NewTransactorMock() Transactor {
	return &transactorMock{}
}

func (t *transactorMock) WithinTransaction(ctx context.Context, tFunc func(context.Context) error) error {
	return tFunc(ctx)
}
