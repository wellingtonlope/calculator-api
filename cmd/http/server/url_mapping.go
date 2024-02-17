package server

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wellingtonlope/calculator-api/internal/app/usecase"
)

func urlMapping(c components) {
	// TODO improving this middleware
	c.server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}
			c.Logger().Errorf("Error Middleware: %v", err)
			var errorUC usecase.Error
			code := http.StatusInternalServerError
			message := "internal error"
			if errors.As(err, &errorUC) {
				message = errorUC.Message
				value := map[usecase.ErrorType]int{
					usecase.ErrorTypeInvalid: http.StatusBadRequest,
					usecase.ErrorTypeUnknown: http.StatusInternalServerError,
				}
				newCode, ok := value[errorUC.Type]
				if ok {
					code = newCode
				}
			}
			return c.JSON(code, map[string]string{
				"message": message,
			})
		}
	})
	c.server.GET("/sum", c.controllers.numbers.Sum)
	c.server.POST("/variable", c.controllers.variable.Create)
}
