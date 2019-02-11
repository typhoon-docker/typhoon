package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Get project with given id,
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
func Routes(e *echo.Echo, dao TyphoonDAO) {

	// Configure middleware with the custom claims type for JWT
	jwtConfig := middleware.JWTConfig{
		Claims:     &JwtCustomClaims{},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}
	log.Println("JWT_SECRET is: '" + os.Getenv("JWT_SECRET") + "'")

	p := e.Group("/projects")
	p.Use(middleware.JWTWithConfig(jwtConfig))

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
		log.Println(project)
		if err := c.Bind(project); err != nil {
			return c.String(http.StatusBadRequest, "Invalid Project info: "+err.Error())
		}
		log.Println(project)
		project.Sanitize()
		log.Println(project)

		// Check if the requested name is available
		if _, err := dao.FindProjectByName(project.Name); err == nil {
			return c.String(http.StatusConflict, "This project name seems to already exist")
		}

		// Get user info (with its id found in JWT)
		user, err := dao.FindUserById(claims.TyphoonId)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Could not find you in the user database: "+err.Error())
		}

		// The project is attributed to the user that requested it
		project.Id = bson.NewObjectId()
		project.BelongsToId = claims.TyphoonId
		project.BelongsTo = &user

		// Insert the project into the database
		if err := dao.InsertProject(*project); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// TODO process the project request, extract the repository url, add a hook using addHook
		// addHook(project)
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
			return c.String(http.StatusBadRequest, "Invalid Project info: "+err.Error())
		}
		project.Sanitize()

		// Check if the id given in url is the same as id in the body
		if id != project.Id.Hex() {
			return c.String(http.StatusBadRequest, "Projects id mismatch")
		}
		// Only continue if project belongs to the user that requested the info (JWT)
		if "admin" != claims.Scope && project.BelongsToId != claims.TyphoonId {
			return c.String(http.StatusUnauthorized, "The project does not belong to you")
		}
		// Check if the requested name is available
		if curProject, err := dao.FindProjectByName(project.Name); curProject.Id.Hex() != id && err == nil {
			return c.String(http.StatusConflict, "This project name seems to already exist")
		}
		// TODO: Not sure of about the belongs_to behaviour
		// Update project in database
		if err := dao.UpdateProject(*project); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
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

	/////////////////////////////
	/////////// TEMP? ///////////
	d := e.Group("/docker")
	d.Use(middleware.JWTWithConfig(jwtConfig))

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

	// Activate the other routes
	RoutesAdmin(e, dao)
	RoutesMisc(e, dao)

	////////////////////////////  // TEMP Make JWTs for tests, and allow routes to list and delete users
	/////////// TEMP ///////////  // Those routes are open for everyone, should not be accessible as is in prod
	e.GET("/token/:login", func(c echo.Context) error {
		userLoginToTest := c.Param("login")
		scope := c.QueryParam("scope")
		if scope == "" {
			scope = "user"
		}
		// Get user from mongoDB, create the entry in db if not found. Get its Id and Scope.
		pUser, err := dao.FindUserByLogin(userLoginToTest)
		if err == mgo.ErrNotFound {
			tUser := ProjectUser{Login: userLoginToTest, FirstName: "foo", LastName: "bar", Email: "nope@nope.fr", Scope: scope}
			nUser, nErr := dao.InsertUser(tUser)
			if nErr != nil {
				log.Println("InsertUser error: " + nErr.Error())
			}
			pUser = nUser
		} else if err != nil {
			log.Println("FindUserByLogin error for " + userLoginToTest + ": " + err.Error())
		}
		// Now user Id and Scope should have the right value

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtCustomClaims{
			pUser.OauthId, pUser.Login, pUser.FirstName, pUser.LastName,
			pUser.Email, pUser.Id.Hex(), pUser.Scope, jwt.StandardClaims{},
		})
		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			log.Println("Could not make JWT: " + err.Error())
		}
		log.Println("Made JWT for " + pUser.Login + ": " + tokenString)
		return c.String(http.StatusOK, tokenString)
	})
	/////////// /TEMP ///////////
	/////////////////////////////
}
