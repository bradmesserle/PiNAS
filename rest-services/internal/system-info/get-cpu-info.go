package system_info

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v4"
)

// Cpu Info Field Names
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

// GetCpu Get the CPU Info
func GetCpu(c echo.Context) error {

	fmt.Println("CPU Info")
	out, err := exec.Command("cat", "/proc/cpuinfo").Output()
	if err != nil {
		log.Println(err)
		return err
	}

	data := string(out)

	cpuInfo := CpuInfo{}

	//Get Number of CPUs
	cpuInfo.NumberOfCpus = strings.Count(data, "processor")

	//Get MIPS
	cpuInfo.BogoMIPS = GetFieldValue(data, bogoMips)

	//Get Architecture
	cpuInfo.Architecture = GetFieldValue(data, architecture)

	//Get Features
	cpuInfo.Features = GetFieldValue(data, features)

	//Get Revision
	cpuInfo.Revision = GetFieldValue(data, revision)

	fmt.Println(cpuInfo)

	jsonData, marshalErr := json.Marshal(cpuInfo)
	if marshalErr != nil {
		log.Println(marshalErr)
		return marshalErr
	}

	return c.JSON(http.StatusOK, jsonData)

}

// GetFieldValue Get the value of a field from the data
func GetFieldValue(data string, field string) string {
	startIndex := strings.LastIndex(data, field)
	endIndex := strings.Index(data[startIndex:], "\n")

	if len(strings.Split(data[startIndex:endIndex+startIndex], ":")) > 0 {
		return strings.Split(data[startIndex:endIndex+startIndex], ":")[1]
	}

	return ""
}

//func GetMockData() string {
//	var data string = "processor       : 0\nBogoMIPS        : 108.00\nFeatures        : fp asimd evtstrm aes pmull sha1 sha2 crc32 atomics fphp asimdhp cpuid asimdrdm lrcpc dcpop asimddp\nCPU implementer : 0x41\nCPU architecture: 8\nCPU variant     : 0x4\nCPU part        : 0xd0b\nCPU revision    : 1\n\nprocessor       : 1\nBogoMIPS        : 108.00\nFeatures        : fp asimd evtstrm aes pmull sha1 sha2 crc32 atomics fphp asimdhp cpuid asimdrdm lrcpc dcpop asimddp\nCPU implementer : 0x41\nCPU architecture: 8\nCPU variant     : 0x4\nCPU part        : 0xd0b\nCPU revision    : 1\n\nprocessor       : 2\nBogoMIPS        : 108.00\nFeatures        : fp asimd evtstrm aes pmull sha1 sha2 crc32 atomics fphp asimdhp cpuid asimdrdm lrcpc dcpop asimddp\nCPU implementer : 0x41\nCPU architecture: 8\nCPU variant     : 0x4\nCPU part        : 0xd0b\nCPU revision    : 1\n\nprocessor       : 3\nBogoMIPS        : 108.00\nFeatures        : fp asimd evtstrm aes pmull sha1 sha2 crc32 atomics fphp asimdhp cpuid asimdrdm lrcpc dcpop asimddp\nCPU implementer : 0x41\nCPU architecture: 8\nCPU variant     : 0x4\nCPU part        : 0xd0b\nCPU revision    : 1\n\nRevision        : e04171\nSerial          : 4e34a94cdf59ebcf\nModel           : Raspberry Pi 5 Model B Rev 1.1"
//	return data
//}
