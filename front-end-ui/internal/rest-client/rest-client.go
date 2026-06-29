package rest_client

import (
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
	err = execRestCall(&pcieInfo, getCpuInfoUrl)
	return pcieInfo, err
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
