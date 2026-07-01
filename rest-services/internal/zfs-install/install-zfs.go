package zfs_install

import (
	"bufio"
	"bytes"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v5"
)

func InstallZfs(c *echo.Context) error {

	// Check to see if kernel header packages are installed and remove them
	headersRemovedErr := CheckForInstalledHeadersAndRemoveThem()

	if headersRemovedErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error removing kernel headers")
	}

	//Install ZFS

	return c.JSON(http.StatusOK, "Installed ZFS successfully")
}

// CheckForInstalledHeadersAndRemoveThem checks if the kernel headers packages are installed
func CheckForInstalledHeadersAndRemoveThem() error {

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
