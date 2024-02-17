package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wellingtonlope/calculator-api/internal/app/usecase"
)

type (
	Variable interface {
		Create(c echo.Context) error
	}
	variable struct {
		createVariable usecase.CreateVariable
	}
)

func NewVariable(createVariable usecase.CreateVariable) Variable {
	return &variable{createVariable: createVariable}
}

func (h variable) Create(c echo.Context) error {
	var input struct {
		Name  string  `json:"name"`
		Value float64 `json:"value"`
	}
	err := c.Bind(&input)
	if err != nil {
		return usecase.NewError("invalid input JSON", err, usecase.ErrorTypeInvalid)
	}
	err = h.createVariable.Handle(c.Request().Context(), usecase.CreateVariableInput(input))
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}
