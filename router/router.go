package router

import (
	"net/http"
	"sample_blog/controller"

	"github.com/labstack/echo"
)

func Init() *echo.Echo {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Entries of sample blog")
	})
	e.GET("/entries_list", controller.GetEntryList())
	e.GET("/entries/:uuid", controller.GetEntryByUuId())

	return e
}
