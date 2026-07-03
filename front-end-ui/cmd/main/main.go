package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pinas/ui/internal"
	"github.com/pinas/ui/internal/components"
	"github.com/pinas/ui/internal/endpoints"
	"github.com/pinas/ui/internal/rest-client"
	"github.com/pinas/ui/internal/structs"
)

func main() {
	app := echo.New()

	//Static Files
	app.StaticFS("/", echo.MustSubFS(internal.StaticFiles, ""))

	//logger
	app.Use(middleware.RequestLogger())

	//app.Use(middleware.Recover())
	//app.Use(middleware.CORS())

	wizardInfo := new(structs.WizardInfo)

	app.GET("/", func(c *echo.Context) error {
		return endpoints.Home(c, components.Home(*wizardInfo))
	})

	app.POST("/next", func(c *echo.Context) error { return endpoints.WizardNext(c, wizardInfo) })

	app.POST("/back", func(c *echo.Context) error { return endpoints.WizardBack(c, wizardInfo) })

	app.GET("/install", func(c *echo.Context) error { return endpoints.Install(c, wizardInfo) })

	//Console output SSE
	app.GET("/consoleStream", func(c *echo.Context) error { return rest_client.StreamHandler(c) })

	// Start the server
	sc := echo.StartConfig{
		Address: ":8080",
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM) // start shutdown process on ctrl+c
	defer cancel()

	// Start the server
	if err := sc.Start(ctx, app); err != nil {
		app.Logger.Error("failed to start server", "error", err)
	}
}
