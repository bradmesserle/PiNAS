package rest_client

import (
	"time"

	"github.com/pinas/common-structs"
)

// CreateZfsDataset  Create ZFS Dataset
func CreateZfsDataset(zfsDataset common_structs.ZfsDataset) (err error) {
	return ExecuteGetSseCall("http://localhost:9090/createZfsDataset", 15*time.Minute)
}
