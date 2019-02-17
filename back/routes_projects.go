package main

import (
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
)

// Get project with given id, if authorized (admin or broject belongs to claimed user)
func getProjectIfAuthorized(c echo.Context, project_id string, claims *JwtCustomClaims) (Project, error) {
	// Get the project from database
	project, err := dao.FindProjectById(project_id)
	if err != nil {
		return project, c.String(http.StatusBadRequest, "Invalid Project ID: "+err.Error())
	}
	// Only give project if it belongs to the user that requested the info (JWT)
	if "admin" != claims.Scope && project.BelongsToId != claims.TyphoonId {
		return project, c.String(http.StatusUnauthorized, "The project does not belong to you")
	}
	// Return ok
	return project, nil
}

// Defines somes routes for the echo server
func RoutesProjects(e *echo.Echo, m echo.MiddlewareFunc, dao TyphoonDAO) {

	p := e.Group("/projects")
	p.Use(m)

	// List projects
	p.GET("", func(c echo.Context) error {
		// Parse the JWT
		claims := c.Get("user").(*jwt.Token).Claims.(*JwtCustomClaims)

		// User asks for all projects ?
		qpa, ok := c.QueryParams()["all"]
		if ok && (qpa[0] == "" || qpa[0] == "true") {
			// Check if admin (JWT)
			if "admin" != claims.Scope {
				return c.String(http.StatusUnauthorized, "You are not admin")
			}
			// Find all projects in database
			projects, err := dao.FindAllProjects()
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
			return c.JSON(http.StatusOK, projects)

		} else {
			// User asks for his projects. Login is taken fron the JWT
			projects, err := dao.FindProjectsOfUser(claims.TyphoonId)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
			return c.JSON(http.StatusOK, projects)
		}
	})

	// Return the project with the specified id
	p.GET("/:id", func(c echo.Context) error {
		// Parse id and JWT
		id := c.Param("id")
		claims := c.Get("user").(*jwt.Token).Claims.(*JwtCustomClaims)

		// Get the project from database
		project, err := dao.FindProjectById(id)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid Project ID: "+err.Error())
		}

		// Only give project if it belongs to the user that requested the info (JWT)
		if "admin" != claims.Scope && project.BelongsToId != claims.TyphoonId {
			return c.String(http.StatusUnauthorized, "The project does not belong to you")
		}
		return c.JSON(http.StatusOK, project)
	})

	// Create a new project
	p.POST("", func(c echo.Context) error {
		// Parse the JWT
		claims := c.Get("user").(*jwt.Token).Claims.(*JwtCustomClaims)

		// Parse the body to find the new project info
		project := new(Project)
		if err := c.Bind(project); err != nil {
			log.Println("Invalid Project info: " + err.Error())
			return c.String(http.StatusBadRequest, "Invalid Project info: "+err.Error())
		}
		project.Sanitize()

		// Check if the requested name is available
		if _, err := dao.FindProjectByName(project.Name); err == nil {
			log.Println("This project name seems to already exist")
			return c.String(http.StatusConflict, "This project name seems to already exist")
		}

		// Get user info (with its id found in JWT)
		user, err := dao.FindUserById(claims.TyphoonId)
		if err != nil {
			log.Println("Could not find you in the user database: " + err.Error())
			return c.String(http.StatusInternalServerError, "Could not find you in the user database: "+err.Error())
		}

		// The project is attributed to the user that requested it
		project.Id = bson.NewObjectId()
		project.BelongsToId = claims.TyphoonId
		project.BelongsTo = &user

		// Insert the project into the database
		if err := dao.InsertProject(*project); err != nil {
			log.Println("Could not insert project in database: " + err.Error())
			return c.String(http.StatusInternalServerError, "Could not insert project in database: "+err.Error())
		}

		addHook(project)
		return c.JSON(http.StatusCreated, project)
	})

	// Update project in db (no need to add hook again)
	p.PUT("/:id", func(c echo.Context) error {
		// Parse id and JWT
		id := c.Param("id")
		claims := c.Get("user").(*jwt.Token).Claims.(*JwtCustomClaims)

		// Parse the body to find the new project info
		project := new(Project)
		if err := c.Bind(project); err != nil {
			log.Println("Invalid Project info: " + err.Error())
			return c.String(http.StatusBadRequest, "Invalid Project info: "+err.Error())
		}
		project.Sanitize()
		project.BelongsToId = project.BelongsTo.Id.Hex()

		// Check if the id given in url is the same as id in the body
		if id != project.Id.Hex() {
			log.Println("Projects id mismatch")
			return c.String(http.StatusBadRequest, "Projects id mismatch")
		}
		// Only continue if project belongs to the user that requested the info (JWT)
		if "admin" != claims.Scope && project.BelongsToId != claims.TyphoonId {
			log.Println("The project does not belong to you")
			return c.String(http.StatusUnauthorized, "The project does not belong to you")
		}
		// Check if the requested name is available
		curProject, err := dao.FindProjectByName(project.Name)
		if curProject.Id.Hex() != id {
			if err == nil {
				return c.String(http.StatusConflict, "This project name seems to already exist")
			} else if err != mgo.ErrNotFound {
				return c.String(http.StatusBadRequest, "Error when using FindProjectByName: "+err.Error())
			}
		}

		// Take care of BelongsTo
		if project.BelongsToId == curProject.BelongsToId {
			project.BelongsTo = curProject.BelongsTo
		} else {
			// Get user info (with its id found in the given BelongsTo data) to assign the project to him
			user, err := dao.FindUserById(project.BelongsTo.Id.Hex())
			if err != nil {
				log.Println("Could not find the user in the database: " + err.Error())
				return c.String(http.StatusInternalServerError, "Could not find the user in the database: "+err.Error())
			}
			project.BelongsTo = &user
		}

		// Update project in database
		if err := dao.UpdateProject(*project); err != nil {
			log.Println("Error while updating project: " + err.Error())
			return c.String(http.StatusInternalServerError, "Error while updating project: "+err.Error())
		}
		return c.JSON(http.StatusOK, project)
	})

	// Delete the project in db
	p.DELETE("/:id", func(c echo.Context) error {
		// Parse id and JWT
		id := c.Param("id")
		claims := c.Get("user").(*jwt.Token).Claims.(*JwtCustomClaims)

		// Get project if authorized
		project, err := getProjectIfAuthorized(c, id, claims)
		if err != nil {
			return err
		}

		// Docker-compose down (in case the project is still running)
		DockerDown(&project)

		// Delete project in database
		if err := dao.DeleteProject(id); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, project)
	})
}
