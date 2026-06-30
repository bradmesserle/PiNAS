package system_updates

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pinas/rest-services/internal/utilities"
)

// AptUpdate Update the apt cache
func AptUpdate(c *echo.Context) error {

	log.Printf("SSE client connected, ip: %v", c.RealIP())

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	count := uint64(0)
	for {
		select {
		case <-c.Request().Context().Done():
			log.Printf("SSE client disconnected, ip: %v", c.RealIP())
			return nil
		case <-ticker.C:
			count++
			event := utilities.Event{
				Data: []byte(fmt.Sprintf("count: %d, time: %s\n\n", count, time.Now().Format(time.RFC3339Nano))),
			}
			if err := event.MarshalTo(w); err != nil {
				return err
			}
			if err := http.NewResponseController(w).Flush(); err != nil {
				return err
			}
		}
	}
}
