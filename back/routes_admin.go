package main

import (
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Defines somes admin routes for the echo server
func RoutesAdmin(e *echo.Echo, dao TyphoonDAO) {

	// Configure middleware with the custom claims type for JWT
	jwtConfig := middleware.JWTConfig{
		Claims:     &JwtCustomClaims{},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}

	a := e.Group("/admin")
	a.Use(middleware.JWTWithConfig(jwtConfig))

	// List users or admins
	a.GET("/list", func(c echo.Context) error {
		// Parse the JWT and check if admin
		claims := c.Get("user").(*jwt.Token).Claims.(*JwtCustomClaims)
		if "admin" != claims.Scope {
			return c.String(http.StatusUnauthorized, "You are not admin")
		}
		qpa, ok := c.QueryParams()["admin"]
		var users []ProjectUser
		var err error
		if ok && (qpa[0] == "" || qpa[0] == "true") {
			users, err = dao.FindAllAdminUsers()
		} else {
			users, err = dao.FindAllUsers()
		}
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, users)
	})

	// Change scope of user
	a.PUT("/scope/:id", func(c echo.Context) error {
		// Parse id and scope, JWT and check if admin
		id := c.Param("id")
		scope := c.QueryParam("scope")
		claims := c.Get("user").(*jwt.Token).Claims.(*JwtCustomClaims)
		if "admin" != claims.Scope {
			return c.String(http.StatusUnauthorized, "You are not admin")
		}
		user, err := dao.FindUserById(id)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Could not find user in the database: "+err.Error())
		}
		user.Scope = scope
		err = dao.UpdateUser(user)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Could not update user: "+err.Error())
		}
		return c.JSON(http.StatusOK, user)
	})

	// Update user
	a.PUT("/user/:id", func(c echo.Context) error {
		// Parse id, JWT and check if admin
		id := c.Param("id")
		claims := c.Get("user").(*jwt.Token).Claims.(*JwtCustomClaims)
		if "admin" != claims.Scope {
			return c.String(http.StatusUnauthorized, "You are not admin")
		}
		// Parse the body to find the new user info
		user := new(ProjectUser)
		if err := c.Bind(user); err != nil {
			return c.String(http.StatusBadRequest, "Invalid user info: "+err.Error())
		}
		// Check if the id given in url is the same as id in the body
		if id != user.Id.Hex() {
			return c.String(http.StatusBadRequest, "Users id mismatch")
		}
		// Update user in database
		if err := dao.UpdateUser(*user); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, user)
	})

	// Delete user
	a.DELETE("/user/:id", func(c echo.Context) error {
		// Parse id, JWT and check if admin
		id := c.Param("id")
		claims := c.Get("user").(*jwt.Token).Claims.(*JwtCustomClaims)
		if "admin" != claims.Scope {
			return c.String(http.StatusUnauthorized, "You are not admin")
		}
		if err := dao.DeleteUser(id); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, id)
	})
}
