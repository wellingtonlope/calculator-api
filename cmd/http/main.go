package main

import (
	"github.com/labstack/echo/v4"
	"github.com/wellingtonlope/calculator-api/internal/app/usecase"
	http_infra "github.com/wellingtonlope/calculator-api/internal/infra/http"
)

func main() {
	e := echo.New()
	sumNumbers := usecase.NewSumNumbers()
	numbers := http_infra.NewNumbers(sumNumbers)
	e.GET("/sum", numbers.Sum)
	e.Logger.Fatal(e.Start(":8080"))
}
