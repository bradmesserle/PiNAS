package rest_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/labstack/gommon/log"
	"github.com/pinas/common-structs"
)

var enableGen3PcieUrl = "http://192.168.0.22:9090/verifyUpdateConfig"
var getPCIeInfoUrl = "http://192.168.0.22:9090/pcieInfo"
var getModelInfoUrl = "http://192.168.0.22:9090/modelInfo"
var getCpuInfoUrl = "http://192.168.0.22:9090/cpuInfo"
var getMemoryInfoUrl = "http://192.168.0.22:9090/memoryInfo"
var verifyUpdateConfigUrl = "http://192.168.0.22:9090/verifyUpdateConfig"
var rebootUrl = "http://192.168.0.22:9090/reboot"
var healthCheckUrl = "http://192.168.0.22:9090/healthCheck"

// EnableGen3Pcie sends a GET request to enable Gen3 PCIe on the specified URL
func EnableGen3Pcie() (resp *http.Response, err error) {
	return http.Get(enableGen3PcieUrl)
}

// GetModelInfo Get Model Info
func GetModelInfo() (model common_structs.PIModel, err error) {
	var modelInfo common_structs.PIModel
	err = execRestCall(&modelInfo, getModelInfoUrl)
	return modelInfo, err
}

// GetCpuInfo Get Cpu Info
func GetCpuInfo() (model common_structs.CpuInfo, err error) {
	var cpuInfo common_structs.CpuInfo
	err = execRestCall(&cpuInfo, getCpuInfoUrl)
	return cpuInfo, err
}

// GetPCIeInfo Get PCIe Info
func GetPCIeInfo() (model common_structs.PCIeInfo, err error) {
	var pcieInfo common_structs.PCIeInfo
	err = execRestCall(&pcieInfo, getPCIeInfoUrl)
	return pcieInfo, err
}

// GetMemoryInfo Get Memory Info
func GetMemoryInfo() (model common_structs.MemoryInfo, err error) {
	var memoryInfo common_structs.MemoryInfo
	err = execRestCall(&memoryInfo, getMemoryInfoUrl)
	return memoryInfo, err
}

// Reboot system
func Reboot(wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	err = execRestCall(nil, rebootUrl)
	return err
}

// HealthCheck Validate the rest-api is running
func HealthCheck() (err error) {

	client := &http.Client{}
	resp, err := client.Get(healthCheckUrl)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	return err
}

// Execute Rest Call and handle errors
func execRestCall(object interface{}, url string) error {
	client := &http.Client{}
	resp, err := client.Get(url)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Info(err.Error())
		}
	}(resp.Body)

	err = json.NewDecoder(resp.Body).Decode(&object)

	if err != nil {

		// Safely extract the type error details
		if unmarshalTypeError, ok := errors.AsType[*json.UnmarshalTypeError](err); ok {
			return fmt.Errorf(
				"JSON type mismatch: cannot parse raw value %q into Go struct field %s.%s (expected Go type: %s) at byte offset %d",
				unmarshalTypeError.Value,
				unmarshalTypeError.Struct,
				unmarshalTypeError.Field,
				unmarshalTypeError.Type,
				unmarshalTypeError.Offset,
			)
		}

		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	return err

}
