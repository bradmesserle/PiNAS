// Package system_updates Update Raspberry PI system configuration files
package system_updates

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// VerifyUpdateConfig Verify and Update config.txt enable pcie gen3
func VerifyUpdateConfig(c echo.Context) error {

	fmt.Println("CPU Info")

	return c.JSON(http.StatusOK, "CPU Info")
}
