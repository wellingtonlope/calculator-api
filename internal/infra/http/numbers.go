package http

import (
	"github.com/labstack/echo/v4"
	"github.com/wellingtonlope/calculator-api/internal/app/usecase"
	"net/http"
	"strconv"
	"strings"
)

type (
	Numbers interface {
		Sum(c echo.Context) error
	}
	numbers struct {
		sum usecase.SumNumbers
	}
)

func NewNumbers(sum usecase.SumNumbers) Numbers {
	return &numbers{sum: sum}
}

func (h numbers) Sum(c echo.Context) error {
	numberParam := c.QueryParam("numbers")
	if len(numberParam) == 0 {
		return c.JSON(http.StatusOK, map[string]float64{"result": 0})
	}
	numbersStringSlice := strings.Split(numberParam, ",")
	numbersSlice := make([]float64, 0, len(numbersStringSlice))
	for _, numberS := range numbersStringSlice {
		number, err := strconv.ParseFloat(numberS, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "numbers values must be numbers"})
		}
		numbersSlice = append(numbersSlice, number)
	}
	result := h.sum.Handle(c.Request().Context(), numbersSlice)
	return c.JSON(http.StatusOK, map[string]float64{"result": result})
}
