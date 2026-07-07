package zfs

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pinas/common-structs"
	"github.com/pinas/rest-services/internal/os-functions"
	"github.com/pinas/rest-services/internal/utilities"
)

// CreateZfsPool Create ZFS Pool
func CreateZfsPool(c *echo.Context) error {

	log.Printf("SSE client connected, ip: %v", c.RealIP())

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	poolInfo := new(common_structs.ZfsPool)

	if err := c.Bind(poolInfo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	log.Printf("Creating ZFS Pool: " + poolInfo.PoolName)
	if err := utilities.SendEventData("Creating ZFS Pool: "+poolInfo.PoolName, "consoleOutput", w); err != nil {
		slog.Error("Error while sendingCreating ZFS Pool String", "Value", err)
	}

	//Wipe Filesystem data
	if poolInfo.WifeFilesystem {

		//Destroy any existing pools
		err := os_functions.DestroyPools(w)
		if err != nil {
			slog.Error("Error while destroying zfs pools", "Value", err)
			return err
		}

		for _, drive := range poolInfo.NvmeDrives {
			err := os_functions.WipeFileSystem("/dev/"+drive.DeviceIdentifier, w)
			if err != nil {
				slog.Error("Error while wiping filesystem", "Value", err)
			}
		}

	}

	//Create ZFS Pool
	err := os_functions.CreateZfsPool(*poolInfo, w)
	if err != nil {
		slog.Error("Error while creating ZFS Pool", "Value", err)
	}

	return nil
}
