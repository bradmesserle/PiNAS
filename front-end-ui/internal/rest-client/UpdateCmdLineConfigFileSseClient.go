package rest_client

import (
	"time"
)

// UpdateCmdLineFile requests and processes updates to the command line configuration file from a remote server.
func UpdateCmdLineFile() error {
	return ExecuteGetSseCall("http://localhost:9090/updateCmdlineConfigFile", 30*time.Minute)
}
