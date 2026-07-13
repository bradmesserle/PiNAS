package zfs

import (
	"log"

	"github.com/labstack/echo/v5"
	"github.com/pinas/rest-services/internal/system-info"
)

func GetZfsPoolStatus(c *echo.Context) error {

	log.Printf("SSE client connected, ip: %v", c.RealIP())

	err := system_info.GetZfsPoolStatus(c)
	if err != nil {
		log.Printf("Error getting ZFS pool status: %v", err)
	}

	return nil
}
