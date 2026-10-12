package rest_client

import (
	"time"
)

func InstallNvmeCli() error {
	return ExecuteGetSseCall("http://localhost:9090/installNvmeCli", 30*time.Minute)
}
