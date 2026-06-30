package system_updates

import (
	"log"
	"log/slog"
	"os/exec"

	"github.com/labstack/echo/v5"
	"github.com/pinas/rest-services/internal/utilities"
)

// AptUpgrade Upgrade the system
func AptUpgrade(c *echo.Context) error {

	log.Printf("SSE client connected, ip: %v", c.RealIP())

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	cmd := exec.Command("sudo", "apt-get", "upgrade", "-y")

	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while running apt upgrade", err)
		return err
	}
	return nil

}
