package zfs_install

import (
	"bufio"
	"bytes"
	"log"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pinas/rest-services/internal/utilities"
)

func InstallZfs(c *echo.Context) error {

	log.Printf("SSE client connected, ip: %v", c.RealIP())

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Check to see if kernel header packages are installed and remove them
	headersRemovedErr := checkForInstalledHeadersAndRemoveThem()
	if headersRemovedErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error removing kernel headers")
	}

	//Install ZFS
	aptErr := aptInstallZfs(w)
	if aptErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error installing zfs")
	}

	return c.JSON(http.StatusOK, "Installed ZFS successfully")
}

// checkForInstalledHeadersAndRemoveThem checks if the kernel headers packages are installed
func checkForInstalledHeadersAndRemoveThem() error {

	var headerPackageName string
	var commonPackageName string

	cmd := exec.Command("dpkg", "--get-selections")
	grepCmd := exec.Command("grep", "headers")

	pipe, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("Error while getting the kernel headers", "Error", err)
	}

	grepCmd.Stdin = pipe

	var out bytes.Buffer
	grepCmd.Stdout = &out

	if err := cmd.Start(); err != nil {
		return err
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

		//Find the header package name
		if strings.Contains(line, "2712") && strings.Contains(line, "install") {
			headerPackageName = strings.Split(line, "\t")[0]
		}

		//Find the common package name
		if strings.Contains(line, "common") && strings.Contains(line, "install") {
			commonPackageName = strings.Split(line, "\t")[0]
		}

	}

	//Remove the common hearers package
	if commonPackageName != "" {
		slog.Info("Removing kernel headers package: ", "Value", commonPackageName)
		cmd := exec.Command("apt-get", "remove", "-y", commonPackageName)
		if err := cmd.Run(); err != nil {
			slog.Error("Error while removing the kernel headers", "Error", err)
		}
	}

	//Install the correct version of the kernel headers
	if headerPackageName != "" {
		slog.Info("Installing kernel headers package: ", "Value", headerPackageName)
		cmd := exec.Command("apt-get", "install", "-y", headerPackageName)
		if err := cmd.Run(); err != nil {
			slog.Error("Error while installing the kernel headers", "Error", err)
		}
	}

	return nil
}

// aptInstallZfs Install ZFS using apt installer
func aptInstallZfs(w http.ResponseWriter) error {

	log.Printf("Installing ZFS")
	if err := utilities.SendEventData("Installing ZFS", "consoleOutput", w); err != nil {
		slog.Error("Error while installing zfs", "Value", err)
	}

	cmd := exec.Command("DEBIAN_FRONTEND=noninteractive", "apt", "install", "zfs-dkms", "zfsutils-linux", "-y")
	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while installing zfs", "Value", err)
		return err
	}

	return nil
}
