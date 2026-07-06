package os_functions

import (
	"log"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"

	"github.com/pinas/rest-services/internal/utilities"
)

// WipeFileSystem Wipe File System for a drive path
func WipeFileSystem(drivePath string, w http.ResponseWriter) error {

	if len(strings.TrimSpace(drivePath)) > 0 {

		//TODO: May want to do more validation..
		log.Printf("Wiping file system: " + drivePath)
		if err := utilities.SendEventData("Wiping filesystem on drive "+drivePath, "consoleOutput", w); err != nil {
			slog.Error("Error while sending wiping file system log info string", "Value", err)
		}

		cmd := exec.Command("wipefs", "-all", "--force", drivePath)
		err := utilities.ExecCmdSseStdoutText(cmd, w)
		if err != nil {
			slog.Error("Error while wiping filesystem: "+drivePath, "Value", err)
			return err
		}
	}

	return nil
}
