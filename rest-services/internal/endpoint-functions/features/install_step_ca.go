package features

import (
	"log"
	"log/slog"

	"github.com/labstack/echo/v5"
	"github.com/pinas/rest-services/internal/install_features"
)

func InstallStepCa(c *echo.Context) error {

	log.Printf("SSE client connected, ip: %v", c.RealIP())

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Install Step CA
	err := install_features.InstallStepCA(w)
	if err != nil {
		slog.Error("Error while installing Step CA", "Value", err)
	}

	return nil

}
