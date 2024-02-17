package server

import (
	"github.com/labstack/echo/v4"
	"github.com/wellingtonlope/calculator-api/internal/app/usecase"
	http_infra "github.com/wellingtonlope/calculator-api/internal/infra/http"
)

func New() *echo.Echo {
	server := echo.New()
	sumNumbers := usecase.NewSumNumbers()
	numbers := http_infra.NewNumbers(sumNumbers)
	server.GET("/sum", numbers.Sum)
	return server
}
