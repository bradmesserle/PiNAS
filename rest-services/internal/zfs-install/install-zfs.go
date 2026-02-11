package zfs_install

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func InstallZfs(c echo.Context) error {

	return c.JSON(http.StatusOK, "Installed ZFS successfully")
}
