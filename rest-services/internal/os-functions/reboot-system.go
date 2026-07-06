package os_functions

import (
	"log/slog"
	"os/exec"

	"github.com/labstack/echo/v5"
)

// Reboot Reboots the system
func Reboot(c *echo.Context) error {

	_, err := exec.Command("reboot").Output()
	if err != nil {
		slog.Error("Error while rebooting the system", "Value", err)
		return err
	}

	return nil

}
