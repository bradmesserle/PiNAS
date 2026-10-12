package rest_client

import (
	"encoding/json"
	"time"

	"github.com/pinas/common-structs"
)

// CreateZfsDataset  Create ZFS Dataset
func CreateZfsDataset(zfsDataset common_structs.ZfsDataset) (err error) {
	// Marshal struct to JSON
	jsonData, _ := json.Marshal(zfsDataset)
	return ExecutePostSseCall("http://localhost:9090/createZfsDataset", jsonData, 15*time.Minute)
}
