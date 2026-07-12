package rest_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

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
var getDriveTelemetryUrl = "http://192.168.0.22:9090/getDriveTelemetry"
var createZfsPoolUrl = "http://192.168.0.22:9090/createZfsPool"

// EnableGen3Pcie sends a GET request to enable Gen3 PCIe on the specified URL
func EnableGen3Pcie() (resp *http.Response, err error) {
	return http.Get(enableGen3PcieUrl)
}

// GetModelInfo Get Model Info
func GetModelInfo() (model common_structs.PIModel, err error) {
	var modelInfo common_structs.PIModel
	err = execRestGetCall(&modelInfo, getModelInfoUrl)
	return modelInfo, err
}

// GetCpuInfo Get Cpu Info
func GetCpuInfo() (model common_structs.CpuInfo, err error) {
	var cpuInfo common_structs.CpuInfo
	err = execRestGetCall(&cpuInfo, getCpuInfoUrl)
	return cpuInfo, err
}

// GetPCIeInfo Get PCIe Info
func GetPCIeInfo() (model common_structs.PCIeInfo, err error) {
	var pcieInfo common_structs.PCIeInfo
	err = execRestGetCall(&pcieInfo, getPCIeInfoUrl)
	return pcieInfo, err
}

// GetMemoryInfo Get Memory Info
func GetMemoryInfo() (model common_structs.MemoryInfo, err error) {
	var memoryInfo common_structs.MemoryInfo
	err = execRestGetCall(&memoryInfo, getMemoryInfoUrl)
	return memoryInfo, err
}

// Reboot system
func Reboot() (err error) {
	err = execRestGetCall(nil, rebootUrl)
	return err
}

// GetDriveTelemetry Get Drive Telemetry
func GetDriveTelemetry() (drives []common_structs.NvmeDrive, err error) {
	var driveTelemetry []common_structs.NvmeDrive
	err = execRestGetCall(&driveTelemetry, getDriveTelemetryUrl)
	return driveTelemetry, err
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
func execRestGetCall(object interface{}, url string) error {
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

// Execute Rest Call and handle errors
func execRestPostCall(objectBody interface{}, postObject interface{}, url string) error {
	client := &http.Client{}

	//Marshal Object
	jsonStr, err := json.Marshal(postObject)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return fmt.Errorf("failed to execute POST request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Info(err.Error())
		}
	}(resp.Body)

	err = json.NewDecoder(resp.Body).Decode(&objectBody)

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
