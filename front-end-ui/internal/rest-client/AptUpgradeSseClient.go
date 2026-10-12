package rest_client

import (
	"time"
)

func AptUpgrade() error {
	return ExecuteGetSseCall("http://localhost:9090/aptUpgrade", 30*time.Minute)
}
