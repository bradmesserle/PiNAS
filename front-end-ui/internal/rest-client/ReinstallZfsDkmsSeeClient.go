package rest_client

import (
	"time"
)

// ReinstallZfsDkms sends an HTTP request to reinstall the ZFS DKMS and processes the server-sent events (SSE) stream.
func ReinstallZfsDkms() (err error) {
	return ExecuteGetSseCall("http://localhost:9090/reinstallZFSDkms", 30*time.Minute)
}
