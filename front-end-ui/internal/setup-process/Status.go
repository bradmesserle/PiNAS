package setup_process

import (
	"log/slog"
	"os"

	"github.com/pinas/ui/internal/structs"
)

func GetStatus() structs.Checkpoint {

	return structs.FreshInstall
}

// SaveStatus saves the current status to a file.
func SaveStatus(checkPoint structs.Checkpoint) error {
	data := []byte(checkPoint)
	err := os.WriteFile("/tmp/installProcess.txt", data, 0644)
	if err != nil {
		slog.Error("Failed to write file", "err", err)
		return err
	}

	return nil

}
