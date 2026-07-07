package os_functions

import (
	"log"
	"log/slog"
	"net/http"
	"os/exec"

	"github.com/pinas/common-structs"
	"github.com/pinas/rest-services/internal/utilities"
)

func CreateZfsDataset(zfsDataset common_structs.ZfsDataset, w http.ResponseWriter) error {

	//Dataset name and Pool name cannot be empty
	if zfsDataset.PoolName != "" || zfsDataset.DatasetName != "" {

		_, err := exec.Command("zfs", "create", zfsDataset.PoolName+"/"+zfsDataset.DatasetName).Output()
		if err != nil {
			log.Println(err)
			return err
		}

		if err := utilities.SendEventData("Created ZFS dataset: "+zfsDataset.PoolName+"/"+zfsDataset.DatasetName, "consoleOutput", w); err != nil {
			slog.Error("Error while sendingCreating ZFS Destroying Pool String", "Value", err)
		}
	}

	return nil

}
