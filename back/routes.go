package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func Routes(e *echo.Echo) {
	e.GET("/projects", func(c echo.Context) error {
		// TODO list projects
		return c.String(http.StatusOK, "")
	})
	e.GET("/projects/:id", func(c echo.Context) error {
		// id := c.Param("id")
		// TODO return project
		return c.String(http.StatusOK, "")
	})
	e.POST("/projects", func(c echo.Context) error {
		// TODO process the project request, extract the repository url, add a hook using addHook
		// addHook(user, repo)
		return c.String(http.StatusOK, "")
	})
	e.PUT("/projects", func(c echo.Context) error {
		// TODO update project in db (no need to add hook again)
		return c.String(http.StatusOK, "")
	})
	e.GET("/healthCheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
}
