package server

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/wellingtonlope/calculator-api/internal/app/repository"
	"github.com/wellingtonlope/calculator-api/internal/app/usecase"
	"github.com/wellingtonlope/calculator-api/internal/infra/http"
	sql_impl "github.com/wellingtonlope/calculator-api/internal/infra/sql"
	"github.com/wellingtonlope/calculator-api/pkg/sql"
)

type (
	components struct {
		server     *echo.Echo
		connection sql.Connection
		transactor sql.Transactor

		repositories
		usecases
		controllers
	}

	repositories struct {
		variable repository.Variable
	}
	usecases struct {
		sumNumbers     usecase.SumNumbers
		createVariable usecase.CreateVariable
	}
	controllers struct {
		numbers  http.Numbers
		variable http.Variable
	}
)

func bootstrapComponents() components {
	connection, transactor, err := sql.NewConnectionAndTransactorMemory()
	fatalError(err)
	c := components{
		server:     echo.New(),
		connection: connection,
		transactor: transactor,
	}

	c.repositories = bootstrapRepositories(c)
	c.usecases = bootstrapUsecases(c)
	c.controllers = bootstrapControllers(c)

	return c
}

func execMigrations(c components) {
	_, file, _, _ := runtime.Caller(1)
	migrationsDir := strings.Replace(file, "components.go", "../../../migrations", 1)
	fs, err := os.ReadDir(migrationsDir)
	fatalError(err)
	var sqls []byte
	for _, f := range fs {
		sql, err := os.ReadFile(fmt.Sprintf("%s/%s", migrationsDir, f.Name()))
		fatalError(err)
		sqls = append(sqls, sql...)
	}
	_, err = c.connection.Exec(context.TODO(), string(sqls))
	fatalError(err)
}

func bootstrapRepositories(c components) repositories {
	execMigrations(c)
	return repositories{
		variable: sql_impl.NewVariable(c.connection),
	}
}

func bootstrapUsecases(c components) usecases {
	return usecases{
		sumNumbers:     usecase.NewSumNumbers(),
		createVariable: usecase.NewCreateVariable(c.repositories.variable),
	}
}

func bootstrapControllers(c components) controllers {
	return controllers{
		numbers:  http.NewNumbers(c.usecases.sumNumbers),
		variable: http.NewVariable(c.usecases.createVariable),
	}
}
