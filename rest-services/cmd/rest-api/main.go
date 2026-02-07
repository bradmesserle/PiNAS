package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pinas/rest-services/internal/system-info"
)

func main() {

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.GET("/cpuInfo", system_info.GetCpu)

	// Start the server
	e.Logger.Fatal(e.Start("localhost:9090"))
}
