package os_functions

import (
	"log"
	"log/slog"
	"net/http"
	"os/exec"

	"github.com/pinas/rest-services/internal/utilities"
)

// UpdateCaCertificates performs updates to CA certificates and communicates progress using Server-Sent Events.
func UpdateCaCertificates(w http.ResponseWriter) error {

	_, err := exec.Command("update-ca-certificates", "--fresh").Output()
	if err != nil {
		log.Println(err)
		return err
	}

	if err := utilities.SendEventData("Updating CA certificates", "consoleOutput", w); err != nil {
		slog.Error("Error while updating CA certificates", "Value", err)
	}

	return nil

}
