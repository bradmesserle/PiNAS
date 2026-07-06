package os_functions

import (
	"bufio"
	"bytes"
	"log"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"

	"github.com/pinas/common-structs"
	"github.com/pinas/rest-services/internal/utilities"
)

// CreateZfsPool Create ZFS Pool Zfspool object should contain the pool name and the nvme drives you want to create the pool with
func CreateZfsPool(zfsPool common_structs.ZfsPool, w http.ResponseWriter) error {

	var drivePaths []string

	// Need to get drive paths. by serial number
	for _, drive := range zfsPool.NvmeDrives {
		drivePath, err := getDrivePath(drive.Serial)
		if err != nil {
			slog.Error("Error while getting drive path", "Error", err)
		}

		//If we have a valid path, add it to the array
		if len(drivePath) > 0 {
			drivePaths = append(drivePaths, drivePath)
		}
	}

	//Create a mirror pool if we have two drives
	if len(drivePaths) == 2 {
		//Generate the create pool command
		//zpool create {pool name} mirror /dev/disk/by-id/nvme-KINGSTON_SNV3S1000G_50026B768714C0F8 /dev/disk/by-id/nvme-KINGSTON_SNV3S1000G_50026B768714C15B
		sendInfoMessage(zfsPool, w)

		cmd := exec.Command("zpool", "create", zfsPool.PoolName, "mirror", drivePaths[0], drivePaths[1])
		err := utilities.ExecCmdSseStdoutText(cmd, w)
		if err != nil {
			slog.Error("Error while creating the zfs pool", "Value", err)
			return err
		}

	}

	//Create a single pool if we have one drive
	if len(drivePaths) == 1 {
		//Generate the create pool command
		//zpool create {pool name} /dev/disk/by-id/nvme-KINGSTON_SNV3S1000G_50026B768714C0F8
		sendInfoMessage(zfsPool, w)

		cmd := exec.Command("zpool", "create", zfsPool.PoolName, drivePaths[0])
		err := utilities.ExecCmdSseStdoutText(cmd, w)
		if err != nil {
			slog.Error("Error while creating the zfs pool", "Value", err)
			return err
		}

	}

	return nil
}

// getDrivePath Get drive path by serial number
func getDrivePath(serialNumber string) (string, error) {

	var id string

	//ls -al /dev/disk/by-id | grep serialNumber
	cmd := exec.Command("ls", "-al", "/dev/disk/by-id/")
	grepCmd := exec.Command("grep", serialNumber)

	pipe, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("Error while getting running ls -al /dev/disk/by-id command", "Error", err)
	}

	grepCmd.Stdin = pipe

	var out bytes.Buffer
	grepCmd.Stdout = &out

	if err := cmd.Start(); err != nil {
		return "", err
	}
	if err := grepCmd.Run(); err != nil {
		slog.Error("Error while getting the kernel headers", "Error", err)
	}
	if err := cmd.Wait(); err != nil {
		slog.Error("Error while getting the kernel headers", "Error", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(out.String()))

	// Loop through each line
	for scanner.Scan() {
		line := scanner.Text()
		slog.Info("Stdout: ", "Value", line)

		if strings.HasSuffix(line, serialNumber) {
			id = line
			break
		}
	}

	return "/dev/disk/by-id/" + id, nil
}

// sendInfoMessage Send an info message to the client
func sendInfoMessage(zfsPool common_structs.ZfsPool, w http.ResponseWriter) {
	log.Printf("Creating ZFS Pool: " + zfsPool.PoolName)
	if err := utilities.SendEventData("Creating ZFS Pool: "+zfsPool.PoolName, "consoleOutput", w); err != nil {
		slog.Error("Error while sending creating ZFS Pool String", "Value", err)
	}
}
