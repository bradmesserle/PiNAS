package system_info

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetDriveTelemetry(c echo.Context) error {

	return c.JSON(http.StatusOK, "Retrieved drive telemetry successfully")
}
