package system_updates

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"

	"github.com/labstack/echo/v5"
)

// AptUpdate Update the apt cache
func AptUpdate(c *echo.Context) error {

	log.Printf("SSE client connected, ip: %v", c.RealIP())

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	cmd := exec.Command("sudo", "apt-get", "update")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		// Prints each line immediately as it is generated
		fmt.Println("Streamed line:", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	//ticker := time.NewTicker(1 * time.Second)
	//defer ticker.Stop()
	//count := uint64(0)
	//for {
	//	select {
	//	case <-c.Request().Context().Done():
	//		log.Printf("SSE client disconnected, ip: %v", c.RealIP())
	//		return nil
	//	case <-ticker.C:
	//		count++
	//		event := utilities.Event{
	//			Data: []byte(fmt.Sprintf("count: %d, time: %s\n\n", count, time.Now().Format(time.RFC3339Nano))),
	//		}
	//		if err := event.MarshalTo(w); err != nil {
	//			return err
	//		}
	//		if err := http.NewResponseController(w).Flush(); err != nil {
	//			return err
	//		}
	//	}
	//}

	return nil

}
