package os_functions

import (
	"bufio"
	"log"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"

	"github.com/pinas/rest-services/internal/utilities"
)

// DestroyPools destroys all ZFS pools on the system.
func DestroyPools(w http.ResponseWriter) error {

	out, err := exec.Command("zpool", "list", "-H").Output()
	if err != nil {
		log.Println(err)
		return err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	for scanner.Scan() {
		line := scanner.Text()

		if line != "no pools available" {

			poolName := strings.Split(line, "\t")[0]

			if err := utilities.SendEventData("Destroying ZFS Pool: "+poolName, "consoleOutput", w); err != nil {
				slog.Error("Error while sendingCreating ZFS Destroying Pool String", "Value", err)
			}

			err := exec.Command("zpool", "destroy", poolName, "-f").Run()
			if err != nil {
				log.Println(err)
				return err
			}

			if err := utilities.SendEventData("Successfully Destroyed ZFS Pool: "+poolName, "consoleOutput", w); err != nil {
				slog.Error("Error while sendingCreating ZFS Destroying Pool String", "Value", err)
			}

		}

	}

	return nil

}
