package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var githubOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/callback",
	ClientID:     "",
	ClientSecret: "",
	Scopes:       []string{"user"},
	Endpoint:     github.Endpoint,
}

func main() {
	e := echo.New()
	e.GET("/", handleMain)
	e.GET("/login", handleLogin)
	e.GET("/callback", handleCallback)
	e.Logger.Fatal(e.Start(":8080"))
}

func handleMain(c echo.Context) error {
	return c.String(http.StatusOK, "visit /login to use Oauth")
}

func handleLogin(c echo.Context) error {
	url := githubOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func handleCallback(c echo.Context) error {
	code := c.QueryParam("code")
	token, err := githubOauthConfig.Exchange(c.Request().Context(), code)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to exchange token: "+err.Error())
	}

	return c.String(http.StatusOK, "Authenticated with GitHub! Token: "+token.AccessToken)
}
