package controller

import (
	"net/http"
	"os"
	"time"

	jwt2 "github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
	jwt.RegisteredClaims
}

func Authenticate(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	if username != os.Getenv("USER_NAME") || password != os.Getenv("BLOGGER_PW") {
		return echo.ErrUnauthorized
	}

	claims := &JwtCustomClaims{
		Name: "Blogger Echo",
		UUID: os.Getenv("BLOGGER_UUID"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour * time.Duration(1))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func Welcome(c echo.Context) error {
	user := c.Get("user").(*jwt2.Token)
	claims := user.Claims.(jwt2.MapClaims)
	name, _ := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}
