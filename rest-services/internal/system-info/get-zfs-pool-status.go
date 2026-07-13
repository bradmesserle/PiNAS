package system_info

import (
	"bufio"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pinas/common-structs"
)

// GetZfsPoolStatus returns the status of the ZFS pool
func GetZfsPoolStatus(c *echo.Context) error {

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	poolInfo := new(common_structs.ZfsPool)

	out, zpoolErr := exec.Command("zpool", "status").Output()
	if zpoolErr != nil {
		log.Println(zpoolErr)
		return zpoolErr
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	for scanner.Scan() {
		line := scanner.Text()

		log.Println(line)
		if strings.Contains(line, "pool:") {
			poolArray := strings.Split(line, " ")
			poolInfo.PoolName = strings.TrimSpace(poolArray[1])
		}

		if strings.Contains(line, "state:") {
			stateArray := strings.Split(line, " ")

			switch strings.TrimSpace(stateArray[1]) {
			case "ONLINE":
				poolInfo.State = common_structs.ZfsStateOnline
			case "OFFLINE":
				poolInfo.State = common_structs.ZfsStateOffline
			case "DEGRADED":
				poolInfo.State = common_structs.ZfsStateDegraded
			default:
				poolInfo.State = common_structs.ZfsStateUnknown
			}

		}

	}

	return c.JSON(http.StatusOK, poolInfo)

}
