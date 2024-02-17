package server

import (
	"log"

	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	c := bootstrapComponents()
	urlMapping(c)
	return c.server
}

func fatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
