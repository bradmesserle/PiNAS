package zfs

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pinas/common-structs"
	os_functions "github.com/pinas/rest-services/internal/os-functions"
)

// CreateZfsDataset Create ZFS Dataset
func CreateZfsDataset(c *echo.Context) error {

	log.Printf("SSE client connected, ip: %v", c.RealIP())

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	zfsDataset := new(common_structs.ZfsDataset)

	if err := c.Bind(zfsDataset); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	//Destroy any existing pools
	err := os_functions.CreateZfsDataset(*zfsDataset, w)
	if err != nil {
		slog.Error("Error while creating zfs dataset", "Value", err)
		return err
	}

	return nil
}
