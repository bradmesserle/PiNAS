package zfs_install

import (
	"bufio"
	"log"
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

	out, err := exec.Command("sudo", "apt", "list", "--installed", "|", "grep", "linux-headers").Output()
	if err != nil {
		log.Println(err)
		return err
	}
	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	for scanner.Scan() {
		line := scanner.Text()
		slog.Info("-->", line)
	}

	return nil
}
