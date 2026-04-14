package rest_client

import (
	"net/http"
)

var enableGen3PcieUrl = "http://192.168.0.13:9090/api/v1/verifyUpdateConfig"
var getPCIeInfoUrl = "http://192.168.0.13:9090/pcieInfo"
var getModelInfoUrl = "http://192.168.0.13:9090/modelInfo"

// EnableGen3Pcie sends a GET request to enable Gen3 PCIe on the specified URL
func EnableGen3Pcie(url string) (resp *http.Response, err error) {
	return http.Get(url)
}
