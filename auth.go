package main

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

const (
	ClientID     = "yourClientID"
	ClientSecret = "yourClientSecret"
)

func main() {
	e := echo.New()

	e.GET("/auth", func(c echo.Context) error {
		clientID := c.QueryParam("client_id")
		redirectURI := c.QueryParam("redirect_uri")
		if clientID != ClientID {
			return c.String(http.StatusBadRequest, "無効なクライアントID")
		}
		return c.Redirect(http.StatusFound, redirectURI+"?code=authorizationCode")
	})

	e.POST("/token", func(c echo.Context) error {
		clientID := c.FormValue("client_id")
		clientSecret := c.FormValue("client_secret")
		code := c.FormValue("code")

		if clientID != ClientID || clientSecret != ClientSecret || code != "authorizationCode" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid_request",
			})
		}

		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["aud"] = clientID
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		t, err := token.SignedString([]byte("yourSecretKey"))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{
			"access_token": t,
			"token_type":   "Bearer",
			"expires_in":   "3600",
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
