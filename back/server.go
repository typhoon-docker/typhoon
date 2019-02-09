package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo"
	"github.com/imroc/req"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type githubHookCreate struct {
	Name   string                 `json:"name"`
	Config githubHookCreateConfig `json:"config"`
}

type githubHookCreateConfig struct {
	Url         string `json:"url"`
	ContentType string `json:"content_type"`
}

type githubTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
}

type githubUserResponse struct {
	Login string `json:"login"`
}

type viarezoTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresAt    int    `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type viarezoUserResponse struct {
	Id        int    `json:"id"`
	Login     string `json:"login"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type oauthService struct {
	Authorize  string
	Token      string
	Parameters map[string]string
}

var (
	oauthServices = map[string]oauthService{
		"VIAREZO": {
			Authorize: "https://auth.viarezo.fr/oauth/authorize",
			Token:     "https://auth.viarezo.fr/oauth/token",
			Parameters: map[string]string{
				"scope":         "default",
				"response_type": "code",
			},
		},
		"GITHUB": {
			Authorize: "https://github.com/login/oauth/authorize",
			Token:     "https://github.com/login/oauth/access_token",
			Parameters: map[string]string{
				"scope": "repo",
			},
		},
	}
)

type hook struct {
	ref      string
	cloneUrl string
	user     string
}

type githubHook struct {
	GitRef     string `json:"ref"`
	Repository struct {
		CloneUrl string `json:"clone_url"`
		Owner    struct {
			Login string `json:"login"`
		} `json:"owner"`
	} `json:"repository"`
}

func authorizeURL(rawOauth string) (string, error) {
	upperOauth := strings.ToUpper(rawOauth)
	lowerOauth := strings.ToLower(rawOauth)
	service, ok := oauthServices[upperOauth]
	if !ok {
		return "", errors.New("oauth: unknown oauth service: " + lowerOauth)
	}

	req, err := http.NewRequest("GET", service.Authorize, nil)
	if err != nil {
		return "", err
	}

	query := req.URL.Query()
	for param, value := range service.Parameters {
		query.Add(param, value)
	}
	query.Add("redirect_uri", os.Getenv("BACKEND_URL")+"/callback/"+lowerOauth)
	query.Add("client_id", os.Getenv(upperOauth+"_CLIENT_ID"))
	query.Add("state", "connection-to-oauth")
	req.URL.RawQuery = query.Encode()

	return req.URL.String(), nil
}

func getToken(user string) string {
	// TODO load from database
	return "###tokentodo###"
}

func setToken(user string, token string) {
	// TODO store into database
}

func addHook(user string, repo string) error {
	buf, err := json.Marshal(githubHookCreate{
		Name: "url",
		Config: githubHookCreateConfig{
			Url:         os.Getenv("BACKEND_URL") + "/hook",
			ContentType: "json",
		},
	})
	if err != nil {
		return err
	}
	r, err := http.NewRequest(http.MethodPost, "https://api.github.com/repos/"+user+"/"+repo+"/hooks", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	r.Header.Add("Authorization", "token "+getToken(user))
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return errors.New("add hook: github: http: unexpected status code " + strconv.Itoa(res.StatusCode) + "!")
	}
	return nil
}

// DAO to access data from the database
var dao = TyphoonDAO{}

func main() {
	loadEnv()

	// Create the DAO object and connect it to the mongo server
	dao.Server = "mongodb://root:example@mongo:27017/"
	dao.Database = "typhoon"
	dao.Connect()

	// echo web server
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS restricted
	// Allows requests from those origins with GET, PUT, POST or DELETE method
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://typhoon.viarezo.fr", "http://typhoon-dev.viarezo.fr:1234/"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// echo routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	})

	e.GET("/callback/viarezo", func(c echo.Context) error {
		body := req.Param{
			"grant_type":    "authorization_code",
			"code":          c.QueryParam("code"),
			"redirect_uri":  os.Getenv("BACKEND_URL") + "/callback/viarezo",
			"client_id":     os.Getenv("VIAREZO_CLIENT_ID"),
			"client_secret": os.Getenv("VIAREZO_CLIENT_SECRET"),
		}
		res, err := req.Post(
			oauthServices["VIAREZO"].Token,
			body,
		)
		if err != nil {
			log.Println("Failed to get token")
			log.Println(body)

			log.Println(err)
			return c.String(http.StatusInternalServerError, "server error")
		}
		var viarezoToken viarezoTokenResponse
		err = res.ToJSON(&viarezoToken)
		if err != nil {
			log.Println("Failed to get token")
			log.Println(res)

			log.Println(err)
			return c.String(http.StatusInternalServerError, "server error")
		}

		res, err = req.Get("https://auth.viarezo.fr/api/user/show/me", req.Header{"Authorization": "Bearer " + viarezoToken.AccessToken})
		if err != nil {
			log.Println("Failed to get my infos")
			log.Println(err)
			return c.String(http.StatusInternalServerError, "server error")
		}
		var user viarezoUserResponse
		err = res.ToJSON(&user)
		if err != nil {
			log.Println("Failed to parse my infos")
			log.Println(res)
			log.Println(err)
			return c.String(http.StatusInternalServerError, "server error")
		}

		// Get user from mongoDB, create the entry in db if not found. Get its Id and Scope.
		pUser, err := dao.FindUserByLogin(user.Login)
		if err == mgo.ErrNotFound {
			log.Println("New user will be made with login: " + user.Login)
			tUser := ProjectUser{OauthId: user.Id, Login: user.Login, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Scope: "user"}
			nUser, nErr := dao.InsertUser(tUser)
			if nErr != nil {
				log.Println("InsertUser error: " + nErr.Error())
				return c.String(http.StatusInternalServerError, "server error")
			}
			pUser = nUser
		} else if err != nil {
			log.Println("FindUserByLogin error for " + user.Login + ": " + err.Error())
			return c.String(http.StatusInternalServerError, "server error")
		}
		// Now user Id and Scope should have the right value

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtCustomClaims{
			user.Id,
			user.Login,
			user.FirstName,
			user.LastName,
			user.Email,
			pUser.Id.Hex(),
			pUser.Scope,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			},
		})

		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			log.Println("Error while using SignedString(): " + err.Error())
			return c.String(http.StatusInternalServerError, "server error")
		}
		values := url.Values{}
		values.Add("token", tokenString)
		return c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_URL")+"/callback/viarezo?"+values.Encode())
	})

	e.GET("/callback/github", func(c echo.Context) error {
		body := req.Param{
			"code":          c.QueryParam("code"),
			"redirect_uri":  os.Getenv("BACKEND_URL") + "/callback/github",
			"client_id":     os.Getenv("GITHUB_CLIENT_ID"),
			"client_secret": os.Getenv("GITHUB_CLIENT_SECRET"),
		}
		res, err := req.Post(
			oauthServices["GITHUB"].Token,
			req.BodyJSON(&body),
			req.Header{"Accept": "application/json"},
		)
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "server error")
		}
		var tokenResponse githubTokenResponse
		err = res.ToJSON(&tokenResponse)
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "server error")
		}
		if tokenResponse.Scope != oauthServices["GITHUB"].Parameters["scope"] {
			// user didn't authorize the repo scope
			// TODO prompt the user to authorize again instead of throwing 500
			return c.String(http.StatusInternalServerError, "server error")
		}
		res, err = req.Get(
			"https://api.github.com/user",
			req.Header{"Authorization": tokenResponse.AccessToken},
		)
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "server error")
		}
		var userReponse githubUserResponse
		err = res.ToJSON(&userReponse)
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "server error")
		}

		values := url.Values{}
		values.Add("token", tokenResponse.AccessToken)
		return c.Redirect(
			http.StatusTemporaryRedirect,
			os.Getenv("FRONTEND_URL")+"/callback/github?"+values.Encode(),
		)
	})

	e.GET("/login/viarezo", func(c echo.Context) error {
		url, err := authorizeURL("VIAREZO")
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "server error")
		}
		return c.Redirect(http.StatusTemporaryRedirect, url)
	})

	e.GET("/login/github", func(c echo.Context) error {
		url, err := authorizeURL("GITHUB")
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "server error")
		}
		return c.Redirect(http.StatusTemporaryRedirect, url)
	})

	e.POST("/hook", func(c echo.Context) error {
		func() {
			var h hook
			if c.Request().Header.Get("X-GitHub-Event") == "push" {
				var gh githubHook
				if err := c.Bind(&h); err != nil {
					e.Logger.Warn(err)
					return
				}
				h = hook{
					ref:      gh.GitRef,
					cloneUrl: gh.Repository.CloneUrl,
					user:     gh.Repository.Owner.Login,
				}
			} else {
				return
			}
			if h.ref != "refs/heads/master" {
				return
			}
			if i := strings.Index(h.cloneUrl, "//"); i != -1 {
				h.cloneUrl = h.cloneUrl[i+len("//"):]
			}
			dir, err := ioutil.TempDir("", "typhoon-clone")
			if err != nil {
				log.Println(err)
				return
			}
			defer os.RemoveAll(dir)
			path, err := filepath.Abs(dir)
			if err != nil {
				log.Println(err)
				return
			}
			cmdGit := exec.Command("git", "clone", "-q", "--depth", "1", "--", h.cloneUrl, path)
			cmdGit.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
			if err := cmdGit.Run(); err != nil {
				log.Fatal(err)
			}
			// TODO run build and install commands
			// TODO push the image to a docker image repository
			// TODO notify the docker slave to restart the container (and use the latest image)
		}()
		return c.String(http.StatusOK, "")
	})

	Routes(e, dao)
	// test()
	// TemplateTest()

	e.Logger.Fatal(e.Start(":80"))
}
