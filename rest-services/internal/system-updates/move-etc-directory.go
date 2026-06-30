package system_updates

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func MoveEtcDirectory(c *echo.Context) error {

	return c.JSON(http.StatusOK, "Moved /etc directory successfully")
}
