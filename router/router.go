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
	e.POST("/authenticate", controller.Authenticate)
	r := e.Group("/restricted")
	config := middleware.JWTConfig{
		Claims:     &controller.JwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	r.Use(middleware.JWTWithConfig(config))
	r.GET("/entries_list", controller.GetEntryList)
	r.GET("/entries/:id", controller.GetEntryById)

	return e
}
