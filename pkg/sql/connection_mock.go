package sql

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func newConnectionMockWithDB() (sqlmock.Sqlmock, Connection, *sqlx.DB) {
	db, mock, _ := sqlmock.New()
	conn := sqlx.NewDb(db, "sqlmock")
	con, _ := newConnection(conn)
	return mock, con, conn
}

func NewConnectionMock() (sqlmock.Sqlmock, Connection) {
	db, mock, _ := sqlmock.New()
	conn := sqlx.NewDb(db, "sqlmock")
	con, _ := newConnection(conn)
	return mock, con
}
