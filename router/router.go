package router

import (
	"net/http"

	"github.com/kohitsujijess/sample_blog/controller"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Entries of sample blog")
	})
	e.GET("/entries_list", controller.GetEntryList)
	e.GET("/entries/:id", controller.GetEntryById)
	e.POST("/authenticate", controller.Authenticate)
	r := e.Group("/restricted")
	config := middleware.JWTConfig{
		Claims:     &controller.JwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	r.Use(middleware.JWTWithConfig(config))
	r.GET("/welcome", controller.Restricted)

	return e
}
