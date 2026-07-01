package zfs_install

import (
	"bytes"
	"fmt"
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
	grepCmd := exec.Command("grep", "headers")

	pipe, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("Error while getting the kernel headers", err)
	}

	grepCmd.Stdin = pipe

	var out bytes.Buffer
	grepCmd.Stdout = &out

	//stdout, _ := grepCmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := grepCmd.Run(); err != nil {
		slog.Error("Error while getting the kernel headers", err)
	}

	if err := cmd.Wait(); err != nil {
		slog.Error("Error while getting the kernel headers", err)
	}

	fmt.Println(out.String())

	return nil
}
