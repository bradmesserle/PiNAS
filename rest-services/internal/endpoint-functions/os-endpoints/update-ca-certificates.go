package os_endpoints

import (
	"log"
	"log/slog"

	"github.com/labstack/echo/v5"
	"github.com/pinas/rest-services/internal/os-functions"
)

// UpdateCaCertificates updates the CA certificates and sends status updates as Server-Sent Events to the client.
func UpdateCaCertificates(c *echo.Context) error {

	log.Printf("SSE client connected, ip: %v", c.RealIP())

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Install NVMe CLI
	err := os_functions.UpdateCaCertificates(w)
	if err != nil {
		slog.Error("Error while updating CA certificates", "Value", err)
		return err
	}

	return nil
}
