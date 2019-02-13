package main

import (
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// Defines somes routes for the echo server
func RoutesDocker(e *echo.Echo, m echo.MiddlewareFunc, dao TyphoonDAO) {

	// Use the given middleware (JWT)
	d := e.Group("/docker")
	d.Use(m)

	// Clone and apply templates for given project id
	d.POST("/templates/:id", func(c echo.Context) error {
		// Parse id and JWT
		id := c.Param("id")
		claims := c.Get("user").(*jwt.Token).Claims.(*JwtCustomClaims)

		// Get project if authorized
		project, err := getProjectIfAuthorized(c, id, claims)
		if err != nil {
			return err
		}

		// Clone the source code
		if err := GetSourceCode(&project); err != nil {
			return c.String(http.StatusInternalServerError, "Could not clone: "+err.Error())
		}

		// Write the templates
		res, err := WriteFromTemplates(&project)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Could not write from templates: "+err.Error())
		}

		// Return ok
		return c.JSON(http.StatusOK, res)
	})

	// Clone and apply templates for given project id
	d.POST("/apply/:id", func(c echo.Context) error {
		// Parse id and JWT
		id := c.Param("id")
		claims := c.Get("user").(*jwt.Token).Claims.(*JwtCustomClaims)

		// Get project if authorized
		project, err := getProjectIfAuthorized(c, id, claims)
		if err != nil {
			return err
		}

		// Clone the source code
		if err := GetSourceCode(&project); err != nil {
			return c.String(http.StatusInternalServerError, "Could not clone: "+err.Error())
		}

		// Write the templates
		res, err := WriteFromTemplates(&project)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Could not write from templates: "+err.Error())
		}

		// Build images
		if err := BuildImages(&project); err != nil {
			return c.String(http.StatusInternalServerError, "Could not build: "+err.Error())
		}

		// Docker-compose up
		if err := DockerUp(&project); err != nil {
			return c.String(http.StatusInternalServerError, "Could not up: "+err.Error())
		}

		// Return ok
		return c.JSON(http.StatusOK, res)
	})

	// Up the deployment
	d.POST("/up/:id", func(c echo.Context) error {
		// Parse id and JWT
		id := c.Param("id")
		claims := c.Get("user").(*jwt.Token).Claims.(*JwtCustomClaims)

		// Get project if authorized
		project, err := getProjectIfAuthorized(c, id, claims)
		if err != nil {
			return err
		}

		// Docker-compose up
		if err := DockerUp(&project); err != nil {
			return c.String(http.StatusInternalServerError, "Could not up: "+err.Error())
		}

		// Return ok
		return c.String(http.StatusOK, "OK")
	})

	// Down the deployment
	d.POST("/down/:id", func(c echo.Context) error {
		// Parse id and JWT
		id := c.Param("id")
		claims := c.Get("user").(*jwt.Token).Claims.(*JwtCustomClaims)

		// Get project if authorized
		project, err := getProjectIfAuthorized(c, id, claims)
		if err != nil {
			return err
		}

		// Docker-compose down
		if err := DockerDown(&project); err != nil {
			return c.String(http.StatusInternalServerError, "Could not down: "+err.Error())
		}

		// Return ok
		return c.String(http.StatusOK, "OK")
	})

	// Get deployment status
	d.GET("/status/:id", func(c echo.Context) error {
		// Parse id and JWT
		id := c.Param("id")
		claims := c.Get("user").(*jwt.Token).Claims.(*JwtCustomClaims)

		// Get project if authorized
		project, err := getProjectIfAuthorized(c, id, claims)
		if err != nil {
			return err
		}

		// Docker-compose ps
		out, err := DockerStatus(&project)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Could not check status: "+err.Error())
		}

		// Return output
		return c.String(http.StatusOK, out)
	})

	// Get logs
	d.GET("/logs/:id", func(c echo.Context) error {
		// Parse id, query params, and JWT
		id := c.Param("id")
		lines_ := c.QueryParam("lines")
		lines, err := strconv.Atoi(lines_)
		if err != nil {
			lines = 30
		}
		claims := c.Get("user").(*jwt.Token).Claims.(*JwtCustomClaims)

		// Get project if authorized
		project, err := getProjectIfAuthorized(c, id, claims)
		if err != nil {
			return err
		}

		// Get Logs
		stdout, err := GetLogsByName(&project, lines)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Could not check logs: "+err.Error())
		}

		// Return output
		return c.String(http.StatusOK, stdout)
	})
}
