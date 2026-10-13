package rest_client

import "time"

// InstallSmallStep installs SmallStep CA
func InstallSmallStep() error {
	return ExecuteGetSseCall("http://localhost:9090/installStepCa", 1*time.Hour)
}
