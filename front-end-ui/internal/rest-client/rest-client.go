package rest_client

import (
	"net/http"
)

var enableGen3PcieUrl = "http://localhost:9090/api/v1/verifyUpdateConfig"

// EnableGen3Pcie sends a GET request to enable Gen3 PCIe on the specified URL
func EnableGen3Pcie(url string) (resp *http.Response, err error) {
	return http.Get(url)
}
