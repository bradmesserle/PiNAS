package rest_client

import "time"

func InstallDns() error {
	return ExecuteGetSseCall("http://localhost:9090/installDnsFeature", 15*time.Minute)
}
