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

// EnableGen3Pcie sends a GET request to enable Gen3 PCIe on the specified URL
func EnableGen3Pcie() (resp *http.Response, err error) {
	return http.Get(enableGen3PcieUrl)
}

func GetModelInfo() (model common_structs.PIModel, err error) {

	var modelInfo common_structs.PIModel

	client := &http.Client{}
	resp, err := client.Get(getModelInfoUrl)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Info(err.Error())
		}
	}(resp.Body)

	err = json.NewDecoder(resp.Body).Decode(&modelInfo)

	if err != nil {

		// Safely extract the type error details
		if unmarshalTypeError, ok := errors.AsType[*json.UnmarshalTypeError](err); ok {
			return modelInfo, fmt.Errorf(
				"JSON type mismatch: cannot parse raw value %q into Go struct field %s.%s (expected Go type: %s) at byte offset %d",
				unmarshalTypeError.Value,
				unmarshalTypeError.Struct,
				unmarshalTypeError.Field,
				unmarshalTypeError.Type,
				unmarshalTypeError.Offset,
			)
		}

		return modelInfo, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return modelInfo, err
}

func GetPCIeInfo() (resp *http.Response, err error) {
	return http.Get(getPCIeInfoUrl)
}
