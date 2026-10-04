package setup_process

import (
	"bufio"
	"log/slog"
	"os"

	"github.com/pinas/ui/internal/structs"
)

func GetStatus() (structs.Checkpoint, error) {

	file, err := os.Open("/opt/pinas/ui-services/installProgress.txt")
	if err != nil {
		slog.Error("Failed to read file", "err", err)
		return structs.FreshInstall, nil
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			slog.Error("Failed to close file", "err", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text() // Retrieves the current line without the newline character
		return structs.Checkpoint(line), nil
	}

	if err := scanner.Err(); err != nil {
		slog.Error("Error while reading file", "err", err)
		return structs.FreshInstall, err
	}

	return structs.FreshInstall, nil
}

// SaveStatus saves the current status to a file.
func SaveStatus(checkPoint structs.Checkpoint) error {
	data := []byte(checkPoint)
	err := os.WriteFile("/opt/pinas/ui-services/installProgress.txt", data, 0644)
	if err != nil {
		slog.Error("Failed to write file", "err", err)
		return err
	}

	return nil

}
