package sql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func NewConnectionAndTransactor(dbConfig DatabaseConfig) (Connection, Transactor, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbConfig.Host, dbConfig.Port, dbConfig.User,
		dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)
	conn, err := sqlx.Open(dbConfig.DriverName, connectionString)
	if err != nil {
		return nil, nil, err
	}
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	conn.DB.SetMaxOpenConns(dbConfig.MaxOpenConnections)
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	conn.DB.SetMaxIdleConns(dbConfig.MaxIdleConnections)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	conn.DB.SetConnMaxLifetime(dbConfig.MaxConnLifetime)
	err = conn.Ping()
	if err != nil {
		return nil, nil, err
	}
	return &connection{db: conn}, &transactor{db: conn}, nil
}

func NewConnectionAndTransactorMemory() (Connection, Transactor, error) {
	conn, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, nil, err
	}
	return &connection{db: conn}, &transactor{db: conn}, nil
}
