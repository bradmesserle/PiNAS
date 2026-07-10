package os_functions

import (
	"log"
	"log/slog"
	"net/http"
	"os/exec"

	"github.com/pinas/rest-services/internal/utilities"
)

func InstallNvmeCli(w http.ResponseWriter) error {

	log.Printf("Installing nvme-cli")
	if err := utilities.SendEventData("Installing nvme-cli", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	cmd := exec.Command("apt", "install", "nvme-cli", "-y")
	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while installing nvme-cli", "Value", err)
		return err
	}

	return nil
}
