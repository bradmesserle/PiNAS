package system_info

import (
	"bufio"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pinas/common-structs"
)

func GetDriveTelemetry(c *echo.Context) error {

	out, err := exec.Command("lsblk", "-d", "-o", "NAME,MODEL,SIZE").Output()
	if err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(http.StatusOK, getDrives(string(out)))
}

// Parse string to extract the drive info
func getDrives(data string) []common_structs.NvmeDrive {

	var drives []common_structs.NvmeDrive
	scanner := bufio.NewScanner(strings.NewReader(data))

	for scanner.Scan() {
		line := scanner.Text()

		drive := createDrive(line)
		if drive != nil {
			drives = append(drives, drive.(common_structs.NvmeDrive))
		}

	}

	return drives
}

// Parse device string block and create drive object
func createDrive(driveString string) interface{} {

	var fields = strings.Split(driveString, " ")

	if strings.Contains(fields[0], "nvme") {
		//create pcie device
		drive := common_structs.NvmeDrive{}

		drive.DeviceIdentifier = strings.TrimSpace(fields[0])
		drive.Name = strings.TrimSpace(fields[1])
		drive.Model = strings.TrimSpace(fields[2])
		drive.Size = strings.TrimSpace(fields[3])

		return drive
	}

	return nil
}
