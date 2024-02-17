package sql

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"github.com/stretchr/testify/assert"
)

func TestNewConnection(t *testing.T) {
	t.Run("should return error when nil is passed in parameter", func(t *testing.T) {
		result, err := newConnection(nil)

		assert.Nil(t, result)
		assert.Equal(t, ErrConnectionNil, err)
	})

	t.Run("should create connection with db connection", func(t *testing.T) {
		db, _, _ := sqlmock.New()
		expectedConnDb := sqlx.NewDb(db, "sqlmock")

		result, err := newConnection(expectedConnDb)

		assert.Nil(t, err)
		resultImpl := result.(*connection)
		assert.Equal(t, expectedConnDb, resultImpl.db)
	})
}

func TestConnection_Select(t *testing.T) {
	t.Run("should execute select", func(t *testing.T) {
		ctx := context.TODO()
		sql, conn, _ := newConnectionMockWithDB()
		expectedRow := []driver.Value{1}
		rows := sqlmock.NewRows([]string{"id"}).AddRow(expectedRow...)
		sql.ExpectQuery("select id from users").WillReturnRows(rows)

		var result []string
		err := conn.Select(ctx, &result, `select id from users`)

		assert.Nil(t, err)
		assert.Equal(t, "1", result[0])
	})
}

func TestConnection_Exec(t *testing.T) {
	t.Run("should execute exec", func(t *testing.T) {
		ctx := context.TODO()
		sql, conn, _ := newConnectionMockWithDB()
		expectedResult := sqlmock.NewResult(2, 1)
		sql.ExpectExec("insert into users").WithArgs(1).WillReturnResult(expectedResult)

		result, err := conn.Exec(ctx, `insert into users(id) values (?)`, 1)

		assert.Nil(t, err)
		lastInserted, _ := result.LastInsertId()
		assert.Equal(t, int64(2), lastInserted)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)
	})

	t.Run("should execute transaction exec", func(t *testing.T) {
		ctx := context.TODO()
		sql, conn, db := newConnectionMockWithDB()
		sql.ExpectBegin()
		tx, err := db.Beginx()
		injectTx(ctx, tx)
		expectedResult := sqlmock.NewResult(2, 1)
		sql.ExpectExec("insert into users").WithArgs(1).WillReturnResult(expectedResult)

		result, err := conn.Exec(ctx, `insert into users(id) values (?)`, 1)

		assert.Nil(t, err)
		lastInserted, _ := result.LastInsertId()
		assert.Equal(t, int64(2), lastInserted)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)
	})
}

func TestConnection_NamedExec(t *testing.T) {
	t.Run("should execute named exec", func(t *testing.T) {
		ctx := context.TODO()
		sql, conn, _ := newConnectionMockWithDB()
		expectedResult := sqlmock.NewResult(2, 1)
		expectedPassed := struct {
			ID string `db:"id"`
		}{ID: "1"}
		sql.ExpectExec("insert into users").WithArgs("1").WillReturnResult(expectedResult)

		result, err := conn.NamedExec(ctx, `insert into users(id) values (:id)`, expectedPassed)

		assert.Nil(t, err)
		lastInserted, _ := result.LastInsertId()
		assert.Equal(t, int64(2), lastInserted)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)
	})

	t.Run("should execute named transaction exec", func(t *testing.T) {
		ctx := context.TODO()
		sql, conn, db := newConnectionMockWithDB()
		sql.ExpectBegin()
		tx, err := db.Beginx()
		injectTx(ctx, tx)
		expectedResult := sqlmock.NewResult(2, 1)
		expectedPassed := struct {
			ID string `db:"id"`
		}{ID: "1"}
		sql.ExpectExec("insert into users").WithArgs("1").WillReturnResult(expectedResult)

		result, err := conn.NamedExec(ctx, `insert into users(id) values (:id)`, expectedPassed)

		assert.Nil(t, err)
		lastInserted, _ := result.LastInsertId()
		assert.Equal(t, int64(2), lastInserted)
		rowsAffected, _ := result.RowsAffected()
		assert.Equal(t, int64(1), rowsAffected)
	})
}

func TestConnection_NamedQuery(t *testing.T) {
	t.Run("should execute named query", func(t *testing.T) {
		ctx := context.TODO()
		sql, conn, _ := newConnectionMockWithDB()
		expectedPassed := struct {
			ID string `db:"id"`
		}{ID: "1"}
		expectedRow := []driver.Value{1}
		rows := sqlmock.NewRows([]string{"id"}).AddRow(expectedRow...)
		sql.ExpectQuery("insert into users").WithArgs("1").WillReturnRows(rows)

		result, err := conn.NamedQuery(ctx, `insert into users(id) values (:id)`, expectedPassed)

		assert.Nil(t, err)
		values := map[string]interface{}{}
		result.Next()
		_ = result.MapScan(values)
		assert.Equal(t, map[string]interface{}{"id": int64(1)}, values)
	})
}

func TestConnection_Get(t *testing.T) {
	t.Run("should execute get", func(t *testing.T) {
		ctx := context.TODO()
		sql, conn, _ := newConnectionMockWithDB()
		expectedRow := []driver.Value{1}
		rows := sqlmock.NewRows([]string{"id"}).AddRow(expectedRow...)
		sql.ExpectQuery("select id from users").WillReturnRows(rows)

		var result int64
		err := conn.Get(ctx, &result, `select id from users`)

		assert.Nil(t, err)
		assert.Equal(t, int64(1), result)
	})
}

func TestConnection_Close(t *testing.T) {
	t.Run("should execute close", func(t *testing.T) {
		sql, conn, _ := newConnectionMockWithDB()
		sql.ExpectClose()

		err := conn.Close()
		assert.Nil(t, err)
	})
}

func TestConnection_Begin(t *testing.T) {
	t.Run("should execute begin", func(t *testing.T) {
		sql, conn, _ := newConnectionMockWithDB()
		sql.ExpectBegin()

		ctx, err := conn.Begin(context.TODO())

		assert.NotNil(t, ctx.Value(keyTx{}))
		assert.Nil(t, err)
	})

	t.Run("should fail to execute begin", func(t *testing.T) {
		sql, conn, _ := newConnectionMockWithDB()
		sql.ExpectBegin().WillReturnError(errors.New("i'm an error"))

		ctx, err := conn.Begin(context.TODO())

		assert.Nil(t, ctx.Value(keyTx{}))
		assert.Equal(t, errors.New("i'm an error"), err)
	})
}

func TestConnection_Commit(t *testing.T) {
	t.Run("should execute commit", func(t *testing.T) {
		sql, conn, _ := newConnectionMockWithDB()
		sql.ExpectBegin()
		sql.ExpectCommit()
		ctx, err := conn.Begin(context.TODO())

		err = conn.Commit(ctx)

		assert.Nil(t, err)
	})

	t.Run("should fail to execute commit", func(t *testing.T) {
		sql, conn, _ := newConnectionMockWithDB()
		sql.ExpectBegin()
		sql.ExpectCommit().WillReturnError(errors.New("i'm an error"))
		ctx, err := conn.Begin(context.TODO())

		err = conn.Commit(ctx)

		assert.Equal(t, errors.New("i'm an error"), err)
	})

	t.Run("should fail when transaction not started", func(t *testing.T) {
		_, conn, _ := newConnectionMockWithDB()

		err := conn.Commit(context.TODO())

		assert.Equal(t, ErrTransactionNotStarted, err)
	})
}

func TestConnection_Rollback(t *testing.T) {
	t.Run("should execute rollback", func(t *testing.T) {
		sql, conn, _ := newConnectionMockWithDB()
		sql.ExpectBegin()
		sql.ExpectRollback()
		ctx, err := conn.Begin(context.TODO())

		err = conn.Rollback(ctx)

		assert.Nil(t, err)
	})

	t.Run("should fail to execute rollback", func(t *testing.T) {
		sql, conn, _ := newConnectionMockWithDB()
		sql.ExpectBegin()
		sql.ExpectRollback().WillReturnError(errors.New("i'm an error"))
		ctx, err := conn.Begin(context.TODO())

		err = conn.Rollback(ctx)

		assert.Equal(t, errors.New("i'm an error"), err)
	})

	t.Run("should fail when transaction not started", func(t *testing.T) {
		_, conn, _ := newConnectionMockWithDB()

		err := conn.Rollback(context.TODO())

		assert.Equal(t, ErrTransactionNotStarted, err)
	})
}
