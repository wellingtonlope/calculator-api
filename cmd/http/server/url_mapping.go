package server

import (
	"github.com/wellingtonlope/calculator-api/cmd/http/server/middleware"
)

func urlMapping(c components) {
	c.server.Use(middleware.Error)
	c.server.GET("/sum", c.controllers.numbers.Sum)
	c.server.POST("/variable", c.controllers.variable.Create)
}
