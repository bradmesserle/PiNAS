package utilities

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// GetHealthCheck returns a health check response
func GetHealthCheck(c *echo.Context) error {
	return c.String(http.StatusOK, "Alive")
}
