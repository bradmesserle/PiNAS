package endpoints

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func HealthCheck(c *echo.Context) error {
	return c.String(http.StatusOK, "")
}
