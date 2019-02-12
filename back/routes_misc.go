package main

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// Defines somes routes for the echo server
func RoutesMisc(e *echo.Echo, m echo.MiddlewareFunc, dao TyphoonDAO) {

	// Just return "OK", showing that the server is up
	e.GET("/healthCheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Check if a project with this name exists
	e.GET("/checkProject", func(c echo.Context) error {
		name := c.QueryParam("name")
		_, err := dao.FindProjectByName(name)
		if err != nil {
			return c.String(http.StatusOK, "false")
		}
		return c.String(http.StatusOK, "true")
	})

	// Use the given middleware (JWT)
	sm := e.Group("/showme")
	sm.Use(m)

	// Get my token info
	sm.GET("", func(c echo.Context) error {
		claims := c.Get("user").(*jwt.Token).Claims.(*JwtCustomClaims)
		return c.JSON(http.StatusOK, claims)
	})
}
