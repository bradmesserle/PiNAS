package rest_client

import (
	"time"

	"github.com/pinas/common-structs"
)

// CreateZfsPool Create ZFS Pool
func CreateZfsPool(poolInfo common_structs.ZfsPool) (err error) {
	return ExecuteGetSseCall("http://localhost:9090/createZfsPool", 15*time.Minute)
}
