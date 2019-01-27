package main

import (
	"net/http"

	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
)

var dao = TyphoonDAO{}

func Routes(e *echo.Echo) {
	// Create the DAO object and connect it to the mongo server
	dao.Server = "mongodb://root:example@mongo:27017/"
	dao.Database = "typhoon"
	dao.Connect()

	// TODO check JWT stuff

	e.GET("/projects", func(c echo.Context) error {
		// List projects
		projects, err := dao.FindAllProjects()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, projects)
	})
	e.GET("/projects/:id", func(c echo.Context) error {
		// Return the project with the specified id
		id := c.Param("id")
		project, err := dao.FindProjectById(id)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid Project ID:"+err.Error())
		}
		return c.JSON(http.StatusOK, project)
	})
	e.GET("/checkProject", func(c echo.Context) error {
		// Check if a project with this name exists
		name := c.QueryParam("name")
		_, err := dao.FindProjectByName(name)
		if err != nil {
			return c.String(http.StatusOK, "false")
		}
		return c.String(http.StatusOK, "true")
	})
	e.POST("/projects", func(c echo.Context) error {
		// Create a new project
		// TODO process the project request, extract the repository url, add a hook using addHook
		// addHook(user, repo)
		project := new(Project)
		if err := c.Bind(project); err != nil {
			return c.String(http.StatusBadRequest, "Invalid Project info: "+err.Error())
		}
		if _, err := dao.FindProjectByName(project.Name); err == nil {
			return c.String(http.StatusBadRequest, "This project name seems to already exist")
		}
		project.Id = bson.NewObjectId()
		if err := dao.InsertProject(*project); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusCreated, project)
	})
	e.PUT("/projects", func(c echo.Context) error {
		// Update project in db (no need to add hook again)
		project := new(Project)
		if err := c.Bind(project); err != nil {
			return c.String(http.StatusBadRequest, "Invalid Project info: "+err.Error())
		}
		if _, err := dao.FindProjectByName(project.Name); err == nil {
			return c.String(http.StatusBadRequest, "This project name seems to already exist")
		}
		if err := dao.UpdateProject(*project); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, project)
	})
	e.DELETE("/projects", func(c echo.Context) error {
		// Delete the project in db
		project := new(Project)
		if err := c.Bind(project); err != nil {
			return c.String(http.StatusBadRequest, "Invalid Project info: "+err.Error())
		}
		if err := dao.DeleteProject(*project); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, project)
	})

	e.GET("/healthCheck", func(c echo.Context) error {
		// Just return "OK", showing that the server is up
		return c.String(http.StatusOK, "OK")
	})
}
