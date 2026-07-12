package os_functions

import (
	"log/slog"
	"net/http"
	"os/exec"

	"github.com/pinas/rest-services/internal/utilities"
)

// UpdateCaCertificates performs updates to CA certificates and communicates progress using Server-Sent Events.
func UpdateCaCertificates(w http.ResponseWriter) error {

	if err := utilities.SendEventData("Updating CA certificates", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	cmd := exec.Command("update-ca-certificates", "--fresh")
	err := utilities.ExecCmdSseStdoutText(cmd, w)

	if err != nil {
		slog.Error("Error while updating CA certificates", "Value", err)
		return err
	}

	return nil

}
