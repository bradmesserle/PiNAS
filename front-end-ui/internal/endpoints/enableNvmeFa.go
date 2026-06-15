package endpoints

import (
	"github.com/labstack/echo/v4"
)

func EnableNvmeFa(c echo.Context) error {
	return c.Redirect(302, "/")
}
