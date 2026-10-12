package rest_client

import (
	"encoding/json"
	"time"

	"github.com/pinas/common-structs"
)

// CreateZfsPool Create ZFS Pool
func CreateZfsPool(poolInfo common_structs.ZfsPool) (err error) {
	// Marshal struct to JSON
	jsonData, _ := json.Marshal(poolInfo)
	return ExecutePostSseCall("http://localhost:9090/createZfsPool", jsonData, 30*time.Minute)
}
