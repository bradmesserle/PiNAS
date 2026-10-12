package rest_client

import (
	"time"
)

func InstallZfs() error {
	return ExecuteGetSseCall("http://localhost:9090/installZFS", 30*time.Minute)
}
