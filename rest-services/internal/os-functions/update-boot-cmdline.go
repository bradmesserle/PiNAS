package os_functions

import (
	"bufio"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/pinas/rest-services/internal/utilities"
)

// UpdateBootCmdlineTextFile updates the cmdline.txt file with cgroup memory settings and sends event data to the client.
func UpdateBootCmdlineTextFile(w http.ResponseWriter) error {

	if err := utilities.SendEventData("Updating cmdline.txt with cgroup memory settings", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	var filepath = "/boot/firmware/cmdline.txt"

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			slog.Error("Error while closing file", "Value", err)
		}
	}(file)

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line+" cgroup_memory=1 cgroup_enable=memory")
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	readFileError := file.Close()
	if readFileError != nil {
		return err
	}

	output := strings.Join(lines, "\n") + "\n"
	writeFileErr := os.WriteFile(filepath, []byte(output), 0644)
	if writeFileErr != nil {
		return err
	}

	return nil

}
