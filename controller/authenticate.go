package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

const secret = "secret"

type JwtCustomClaims struct {
	Name  string `json:"name"`
	UUID  string `json:"uuid"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func Authenticate(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username != os.Getenv("USER_NAME") || password != os.Getenv("BLOGGER_PW") {
		return echo.ErrUnauthorized
	}

	claims := &JwtCustomClaims{
		Name:  "Blogger Echo",
		UUID:  os.Getenv("BLOGGER_UUID"),
		Admin: true,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
