package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pinas/ui/internal"
	"github.com/pinas/ui/internal/endpoints"
	setup_navigation "github.com/pinas/ui/internal/setup-navigation"
	setup_process "github.com/pinas/ui/internal/setup-process"
	"github.com/pinas/ui/internal/structs"
)

func main() {

	//Setup web server
	setupWebServer()

}

// setupWebServer Setup web server
func setupWebServer() {

	//Echo web server
	app := echo.New()

	//Step Status
	status := new(structs.SetupInfo)

	//Get current checkpoint
	checkpoint, err := setup_process.GetStatus()
	if err != nil {
		slog.Info("Error while getting checkpoint setting to fresh install", "err", err)
		checkpoint = structs.FreshInstall
	}

	if checkpoint == structs.FreshInstall {
		status.IsRunning = false
	} else {
		status.IsRunning = true
	}

	//Static Files
	app.StaticFS("/", echo.MustSubFS(internal.StaticFiles, ""))

	//logger
	app.Use(middleware.RequestLogger())

	//app.Use(middleware.Recover())
	//app.Use(middleware.CORS())

	wizardInfo := new(structs.WizardInfo)

	app.GET("/", func(c *echo.Context) error { return endpoints.Home(c, wizardInfo, status) })

	app.POST("/next", func(c *echo.Context) error { return setup_navigation.WizardNext(c, wizardInfo, status) })

	app.POST("/back", func(c *echo.Context) error { return setup_navigation.WizardBack(c, wizardInfo, status) })

	app.GET("/setup", func(c *echo.Context) error { return endpoints.Setup(c, wizardInfo, status) })

	//Console output SSE
	app.GET("/consoleStream", func(c *echo.Context) error { return endpoints.ConsoleLogStreamHandler(c) })

	// Start the server
	sc := echo.StartConfig{
		Address: ":8080",
		BeforeServeFunc: func(s *http.Server) error {
			s.WriteTimeout = 0 // IMPORTANT: disable for SSE
			return nil
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM) // start shutdown process on ctrl+c
	defer cancel()

	// Start the server
	if err := sc.Start(ctx, app); err != nil {
		app.Logger.Error("failed to start server", "error", err)
	}

}
