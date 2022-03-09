package router

import (
	"fmt"
	"net/http"
	"sample_blog/controller"

	"github.com/labstack/echo"
)

func Init() {
	fmt.Println("Init router: Hello from docker")

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Entries of sample blog")
	})

	e.GET("/entries_list", controller.GetEntryList())

	e.Start(":1323")
}
