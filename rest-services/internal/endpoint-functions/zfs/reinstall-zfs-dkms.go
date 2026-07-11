package zfs

import (
	"bufio"
	"log"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pinas/rest-services/internal/utilities"
)

// ReinstallZfsDkms reinstalls the zfs-dkms package and sends server-sent events to the client regarding the operation status.
func ReinstallZfsDkms(c *echo.Context) error {

	log.Printf("SSE client connected, ip: %v", c.RealIP())

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Install new kernel headers
	headerErr := installKernelHeaders(w)
	if headerErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error installing new header package")
	}
	//Reinstall zfs-dkms
	aptErr := aptReinstallZfsDkms(w)
	if aptErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error re-installing zfs-dkms")
	}

	// Send Success Message
	log.Printf("Reinstalled zfs-dkms Successfully")
	if err := utilities.SendEventData("Reinstalled zfs-dkms Successfully", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	return c.JSON(http.StatusOK, "Re-Installed zfs-dkms successfully")
}

// installKernelHeaders installs the required kernel headers on the system and returns an error if the installation fails.
func installKernelHeaders(w http.ResponseWriter) error {

	log.Printf("Installing New Kernel Headers")
	if err := utilities.SendEventData("Installing New Kernel Headers", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	out, errListCmd := exec.Command("sudo", "ls", "-tr", "/root").Output()
	if errListCmd != nil {
		log.Println(errListCmd)
		return errListCmd
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	var headerFile string

	for scanner.Scan() {
		headerFile = scanner.Text()
	}

	cmd := exec.Command("sudo", "apt", "install", "/root/"+headerFile, "-y")
	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while installing new header package", "Value", err)
		return err
	}

	return nil
}

// aptReinstallZfsDkms Install ZFS using apt installer
func aptReinstallZfsDkms(w http.ResponseWriter) error {

	log.Printf("Reinstalling ZFS-DKMS")
	if err := utilities.SendEventData("Reinstalling zfs-dkms", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	cmd := exec.Command("sudo", "DEBIAN_FRONTEND=noninteractive", "apt", "reinstall", "zfs-dkms", "-y")
	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while re-installing zfs-dkms", "Value", err)
		return err
	}

	return nil
}
