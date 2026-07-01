package zfs_install

import (
	"bufio"
	"log/slog"
	"net/http"
	"os/exec"

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
	grepCmd := exec.Command("grep", "linux-headers")

	pipe, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("Error while getting the kernel headers", err)
	}

	stdout, _ := grepCmd.StdoutPipe()

	grepCmd.Stdin = pipe

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := grepCmd.Run(); err != nil {
		slog.Error("Error while getting the kernel headers", err)
	}

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		line := scanner.Text()
		slog.Info("-->", line)
	}

	if err := cmd.Wait(); err != nil {
		slog.Error("Error while getting the kernel headers", err)
	}

	return nil
}
