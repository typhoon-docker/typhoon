package main

import (
	"log"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Defines somes routes for the echo server
func Routes(e *echo.Echo, dao TyphoonDAO) {

	// Configure middleware with the custom claims type for JWT
	jwtConfig := middleware.JWTConfig{
		Claims:     &JwtCustomClaims{},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}
	m := middleware.JWTWithConfig(jwtConfig)
	log.Println("JWT_SECRET is: '" + os.Getenv("JWT_SECRET") + "'")

	// Activate the other routes
	RoutesProjects(e, m, dao)
	RoutesDocker(e, m, dao)
	RoutesAdmin(e, m, dao)
	RoutesMisc(e, m, dao)

	// JWTs for tests, and allow routes to list and delete users
	e.GET("/token/:login", func(c echo.Context) error {
		// Close this route in production mode
		if os.Getenv("GO_ENV") == "production" {
			return c.String(http.StatusUnauthorized, "You cannot do that. At least not in production.")
		}

		// Gather requested info
		userLoginToTest := c.Param("login")
		scope := c.QueryParam("scope")
		if scope == "" {
			scope = "user"
		}

		// Get user from mongoDB, create the entry in db if not found.
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

		// Now we have all the info we need about the user
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
}
