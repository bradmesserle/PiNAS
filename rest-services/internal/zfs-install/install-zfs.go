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

	cmd := exec.Command("dpkg", "--get-selections")
	grepCmd := exec.Command("grep", "headers")

	pipe, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("Error while getting the kernel headers", err)
	}

	grepCmd.Stdin = pipe

	var out bytes.Buffer
	grepCmd.Stdout = &out

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := grepCmd.Run(); err != nil {
		slog.Error("Error while getting the kernel headers", err)
	}

	if err := cmd.Wait(); err != nil {
		slog.Error("Error while getting the kernel headers", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	// Loop through each line
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "common") && strings.Contains(line, "install") {

			result := strings.Split(line, "\t")
			slog.Info("Removing kernel headers package: ", result[0])

			//remove the package
			cmd := exec.Command("apt-get", "remove", "-y", result[0])
			if err := cmd.Run(); err != nil {
				slog.Error("Error while removing the kernel headers", err)
			}
		}

	}

	return nil
}
