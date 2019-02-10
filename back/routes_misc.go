package main

import (
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Defines somes routes for the echo server
func RoutesMisc(e *echo.Echo, dao TyphoonDAO) {

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

	// Configure middleware with the custom claims type for JWT
	jwtConfig := middleware.JWTConfig{
		Claims:     &JwtCustomClaims{},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}

	sm := e.Group("/showme")
	sm.Use(middleware.JWTWithConfig(jwtConfig))

	// Get my token info
	sm.GET("", func(c echo.Context) error {
		claims := c.Get("user").(*jwt.Token).Claims.(*JwtCustomClaims)
		return c.JSON(http.StatusOK, claims)
	})
}
