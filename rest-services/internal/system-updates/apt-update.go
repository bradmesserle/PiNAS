package system_updates

import (
	"log"
	"log/slog"
	"os/exec"

	"github.com/labstack/echo/v5"
	"github.com/pinas/rest-services/internal/utilities"
)

// AptUpdate Update the apt repos
func AptUpdate(c *echo.Context) error {

	log.Printf("SSE client connected, ip: %v", c.RealIP())

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	log.Printf("Updating APT Repositories")
	if err := utilities.SendEventData("Updating APT Repositories", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	cmd := exec.Command("sudo", "apt-get", "update", "-y")

	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while running apt update", "Value", err)
		return err
	}

	return nil

}
