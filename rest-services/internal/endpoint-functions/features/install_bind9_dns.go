package features

import (
	"log"
	"log/slog"

	"github.com/labstack/echo/v5"
	"github.com/pinas/rest-services/internal/install_features"
)

func InstallBind9Dns(c *echo.Context) error {

	log.Printf("SSE client connected, ip: %v", c.RealIP())

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Install DNS CLI
	err := install_features.InstallBind9DNS(w)
	if err != nil {
		slog.Error("Error while installing bind9 DNS", "Value", err)
		return err
	}

	return nil
}
