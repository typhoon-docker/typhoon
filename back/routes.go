package main

import (
	"net/http"
	"strings"

	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
)

var dao = TyphoonDAO{}

func Routes(e *echo.Echo) {

	e.GET("/healthCheck", func(c echo.Context) error {
		// Just return "OK", showing that the server is up
		return c.String(http.StatusOK, "OK")
	})

	// Create the DAO object and connect it to the mongo server
	dao.Server = "mongodb://root:example@mongo:27017/"
	dao.Database = "typhoon"
	dao.Connect()

	// TODO check JWT stuff

	e.GET("/projects", func(c echo.Context) error {
		// List projects

		if strings.ToLower(c.QueryParam("all")) == "true" {
			// User asks for all projects
			// TODO check if admin (jwt)
			projects, err := dao.FindAllProjects()
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
			return c.JSON(http.StatusOK, projects)

		} else {
			// User asks for his projects (jwt)
			// TODO get user info from jwt
			userId := "TODO"
			projects, err := dao.FindProjectsOfUser(userId)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
			return c.JSON(http.StatusOK, projects)
		}
	})

	e.GET("/projects/:id", func(c echo.Context) error {
		// Return the project with the specified id
		id := c.Param("id")
		project, err := dao.FindProjectById(id)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid Project ID:"+err.Error())
		}
		// TODO only give project if it belongs to the user that requested the info (jwt)
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
		project := new(Project)
		if err := c.Bind(project); err != nil {
			return c.String(http.StatusBadRequest, "Invalid Project info: "+err.Error())
		}
		if _, err := dao.FindProjectByName(project.Name); err == nil {
			return c.String(http.StatusConflict, "This project name seems to already exist")
		}
		project.Id = bson.NewObjectId()
		if err := dao.InsertProject(*project); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		// TODO process the project request, extract the repository url, add a hook using addHook
		// addHook(user, repo)
		return c.JSON(http.StatusCreated, project)
	})

	e.PUT("/projects/:id", func(c echo.Context) error {
		// Update project in db (no need to add hook again)
		id := c.Param("id")
		project := new(Project)
		if err := c.Bind(project); err != nil {
			return c.String(http.StatusBadRequest, "Invalid Project info: "+err.Error())
		}
		if id != project.Id.Hex() {
			return c.String(http.StatusBadRequest, "Projects id mismatch")
		}
		// TODO only continue if project belongs to the user (jwt)
		if _, err := dao.FindProjectByName(project.Name); err == nil {
			return c.String(http.StatusConflict, "This project name seems to already exist")
		}
		if err := dao.UpdateProject(*project); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, project)
	})

	e.DELETE("/projects/:id", func(c echo.Context) error {
		// Delete the project in db
		id := c.Param("id")
		project := new(Project)
		if err := c.Bind(project); err != nil {
			return c.String(http.StatusBadRequest, "Invalid Project info: "+err.Error())
		}
		if id != project.Id.Hex() {
			return c.String(http.StatusBadRequest, "Projects id mismatch")
		}
		// TODO only continue if project belongs to the user (jwt)
		if err := dao.DeleteProject(*project); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, project)
	})
}
