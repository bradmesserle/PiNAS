package features

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	common_structs "github.com/pinas/common-structs"
	"github.com/pinas/rest-services/internal/install_features"
)

func InstallBind9Dns(c *echo.Context) error {

	log.Printf("SSE client connected, ip: %v", c.RealIP())

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	dnsSettings := new(common_structs.DnsSettings)

	if err := c.Bind(dnsSettings); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload - Unable to unmarshall json data",
		})
	}

	// Install DNS CLI
	err := install_features.InstallBind9DNS(*dnsSettings, w)
	if err != nil {
		slog.Error("Error while installing bind9 DNS", "Value", err)
		return err
	}

	return nil
}
