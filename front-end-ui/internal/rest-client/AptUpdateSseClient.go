package rest_client

import "time"

func AptUpdate() error {
	return ExecuteGetSseCall("http://localhost:9090/aptUpdate", 5*time.Minute)
}
