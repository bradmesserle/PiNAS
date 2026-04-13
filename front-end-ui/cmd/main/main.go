package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pinas/ui/internal"
	"github.com/pinas/ui/internal/components"
	"github.com/pinas/ui/internal/endpoints"
)

func main() {
	app := echo.New()

	defer func(e *echo.Echo) {
		err := e.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(app)

	//Static Files
	app.StaticFS("/", echo.MustSubFS(internal.StaticFiles, ""))

	//logger
	app.Use(middleware.RequestLogger())

	//app.Use(middleware.Recover())
	//app.Use(middleware.CORS())

	app.GET("/", func(c echo.Context) error {
		return endpoints.Home(c, components.Home())
	})

	// Start the server
	app.Logger.Fatal(app.Start(":8080"))

}
