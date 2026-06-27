package system_info

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"

	"github.com/labstack/echo/v4"
	"github.com/pinas/common-structs"
)

const memTotal = "MemTotal"
const memFree = "MemFree"
const available = "MemAvailable"
const buffers = "Buffers"
const cached = "Cached"
const swapTotal = "SwapTotal"
const swapFree = "SwapFree"

// GetMemoryInfo Get System Memory Info
func GetMemoryInfo(c echo.Context) error {

	out, err := exec.Command("cat", "/proc/meminfo").Output()
	if err != nil {
		log.Println(err)
		return err
	}

	memoryInfo := common_structs.MemoryInfo{}
	memoryInfo.Total = GetLastFieldValue(string(out), memTotal)
	memoryInfo.Free = GetLastFieldValue(string(out), memFree)
	memoryInfo.Available = GetLastFieldValue(string(out), available)
	memoryInfo.Buffers = GetLastFieldValue(string(out), buffers)
	memoryInfo.Cached = GetLastFieldValue(string(out), cached)
	memoryInfo.SwapTotal = GetLastFieldValue(string(out), swapTotal)
	memoryInfo.SwapFree = GetLastFieldValue(string(out), swapFree)

	jsonData, marshalErr := json.Marshal(memoryInfo)
	if marshalErr != nil {
		log.Println(marshalErr)
		return marshalErr
	}

	return c.JSON(http.StatusOK, string(jsonData))
}
