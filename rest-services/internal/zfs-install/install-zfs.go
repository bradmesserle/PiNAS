package zfs_install

import (
	"bufio"
	"log"
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

	cmd := exec.Command("dpkg", "--get-selections", "|", "grep", "linux-headers")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		line := scanner.Text()
		slog.Info("-->", line)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	return nil
}
