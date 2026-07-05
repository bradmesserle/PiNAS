package os_functions

import (
	"log"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pinas/rest-services/internal/utilities"
)

// WipeFileSystem Wipe File System for a drive path
func WipeFileSystem(c *echo.Context) error {

	drivePath := c.QueryParam("drivePath")

	log.Printf("SSE client connected, ip: %v", c.RealIP())

	if len(strings.TrimSpace(drivePath)) > 0 {

		//TODO: May want to do more validation..

		w := c.Response()
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		log.Printf("Wiping file system: " + drivePath)
		if err := utilities.SendEventData("Wiping filesystem on drive "+drivePath, "consoleOutput", w); err != nil {
			slog.Error("Error while sending wiping file system log info", "Value", err)
		}

		cmd := exec.Command("wipefs", "-all", "--force", drivePath)
		err := utilities.ExecCmdSseStdoutText(cmd, w)
		if err != nil {
			slog.Error("Error while wiping filesystem: "+drivePath, "Value", err)
			return err
		}
	}

	return c.JSON(http.StatusOK, nil)
}
