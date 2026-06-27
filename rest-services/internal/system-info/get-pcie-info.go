package system_info

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pinas/common-structs"
)

// GetPcieInfo Get PCIe Device Info
func GetPcieInfo(c echo.Context) error {

	//Get list of devices
	out, err := exec.Command("lspci", "-mmv").Output()
	if err != nil {
		log.Println(err)
		return err
	}

	pcieInfo := common_structs.PCIeInfo{}

	//Create devices object
	pcieInfo.PCIeDevices = getDevices(string(out))

	//Get PCIe Bus details
	detailsOut, err := exec.Command("sudo", "lspci", "-vv").Output()
	if err != nil {
		log.Println(err)
	}

	//Get Link Status for each Device
	for i, device := range pcieInfo.PCIeDevices {
		pcieInfo.PCIeDevices[i] = getLinkStatus(string(detailsOut), device)
	}

	jsonData, marshalErr := json.Marshal(pcieInfo)
	if marshalErr != nil {
		log.Println(marshalErr)
		return marshalErr
	}

	return c.JSON(http.StatusOK, string(jsonData))
}

// Parse string to extract the devices
func getDevices(data string) []common_structs.PCIeDevice {

	var devices []common_structs.PCIeDevice

	scanner := bufio.NewScanner(strings.NewReader(data))
	var currentBlock []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			if len(currentBlock) > 0 {
				devices = append(devices, createDevice(currentBlock))
				currentBlock = nil
			}
			continue
		}

		currentBlock = append(currentBlock, line)
	}

	return devices
}

// Parse device string block and create device object
func createDevice(deviceStringBlock []string) common_structs.PCIeDevice {

	//create pcie device
	device := common_structs.PCIeDevice{}

	//Set Slot
	if len(deviceStringBlock) >= 1 {
		device.Slot = strings.TrimSpace(strings.Split(deviceStringBlock[0], "Slot:")[1])
	}

	//Set Vendor
	if len(deviceStringBlock) >= 3 {
		device.Vendor = strings.TrimSpace(strings.Split(deviceStringBlock[2], "Vendor:")[1])
	}

	//Set Description
	if len(deviceStringBlock) >= 4 {
		device.Description = strings.TrimSpace(strings.Split(deviceStringBlock[3], "Device:")[1])
	}

	return device
}

func getLinkStatus(busDetails string, device common_structs.PCIeDevice) common_structs.PCIeDevice {

	scanner := bufio.NewScanner(strings.NewReader(busDetails))
	var deviceFound bool = false
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, string(device.Slot)) {
			deviceFound = true
		}

		if deviceFound {
			if strings.Contains(line, "LnkSta") {
				device.Speed = strings.TrimSpace(strings.Split(line, "LnkSta:")[1])
				break
			}
		}

	}
	return device
}
