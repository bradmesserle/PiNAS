package system_info

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v4"
)

// Cpu Info Field Names
const processor = "processor"
const bogoMips = "BogoMIPS"
const architecture = "CPU architecture"
const features = "Features"
const revision = "CPU revision"

type CpuInfo struct {
	NumberOfCpus int    `json:"numberOfCpus"`
	BogoMIPS     string `json:"bogoMIPS"`
	Architecture string `json:"architecture"`
	Revision     string `json:"revision"`
	Features     string `json:"features"`
}

// GetCpuInfo  Get the CPU Info
func GetCpuInfo(c echo.Context) error {

	out, err := exec.Command("cat", "/proc/cpuinfo").Output()
	if err != nil {
		log.Println(err)
		return err
	}

	cpuInfo := CpuInfo{}

	//Get Number of CPUs
	cpuInfo.NumberOfCpus = strings.Count(string(out), processor)

	//Get MIPS
	cpuInfo.BogoMIPS = GetFieldValue(string(out), bogoMips)

	//Get Architecture
	cpuInfo.Architecture = GetFieldValue(string(out), architecture)

	//Get Features
	cpuInfo.Features = GetFieldValue(string(out), features)

	//Get Revision
	cpuInfo.Revision = GetFieldValue(string(out), revision)

	jsonData, marshalErr := json.Marshal(cpuInfo)
	if marshalErr != nil {
		log.Println(marshalErr)
		return marshalErr
	}

	return c.JSON(http.StatusOK, string(jsonData))

}
