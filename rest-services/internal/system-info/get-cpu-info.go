package system_info

import (
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pinas/common-structs"
)

// Cpu Info Field Names
const processor = "processor"
const bogoMips = "BogoMIPS"
const architecture = "CPU architecture"
const features = "Features"
const revision = "CPU revision"

// GetCpuInfo  Get the CPU Info
func GetCpuInfo(c *echo.Context) error {

	out, err := exec.Command("cat", "/proc/cpuinfo").Output()
	if err != nil {
		log.Println(err)
		return err
	}

	cpuInfo := common_structs.CpuInfo{}

	//Get Number of CPUs
	cpuInfo.NumberOfCpus = strings.Count(string(out), processor)

	//Get MIPS
	cpuInfo.BogoMIPS = GetLastFieldValue(string(out), bogoMips)

	//Get Architecture
	cpuInfo.Architecture = GetLastFieldValue(string(out), architecture)

	//Get Features
	cpuInfo.Features = GetLastFieldValue(string(out), features)

	//Get Revision
	cpuInfo.Revision = GetLastFieldValue(string(out), revision)

	return c.JSON(http.StatusOK, cpuInfo)

}
