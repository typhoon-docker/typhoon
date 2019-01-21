package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func Routes(e *echo.Echo) {
	e.GET("/healthCheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
}
