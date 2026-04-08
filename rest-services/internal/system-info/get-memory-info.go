package system_info

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"

	"github.com/labstack/echo/v4"
)

const memTotal = "MemTotal"
const memFree = "MemFree"
const available = "MemAvailable"
const buffers = "Buffers"
const cached = "Cached"
const swapTotal = "SwapTotal"
const swapFree = "SwapFree"

type MemoryInfo struct {
	Total     string `json:"total"`
	Free      string `json:"free"`
	Available string `json:"available"`
	Buffers   string `json:"buffers"`
	Cached    string `json:"cached"`
	SwapTotal string `json:"swapTotal"`
	SwapFree  string `json:"swapFree"`
}

// GetMemoryInfo Get System Memory Info
func GetMemoryInfo(c echo.Context) error {

	out, err := exec.Command("cat", "/proc/meminfo").Output()
	if err != nil {
		log.Println(err)
		return err
	}

	memoryInfo := MemoryInfo{}
	memoryInfo.Total = GetFieldValue(string(out), memTotal)
	memoryInfo.Free = GetFieldValue(string(out), memFree)
	memoryInfo.Available = GetFieldValue(string(out), available)
	memoryInfo.Buffers = GetFieldValue(string(out), buffers)
	memoryInfo.Cached = GetFieldValue(string(out), cached)
	memoryInfo.SwapTotal = GetFieldValue(string(out), swapTotal)
	memoryInfo.SwapFree = GetFieldValue(string(out), swapFree)

	jsonData, marshalErr := json.Marshal(memoryInfo)
	if marshalErr != nil {
		log.Println(marshalErr)
		return marshalErr
	}

	return c.JSON(http.StatusOK, string(jsonData))
}
