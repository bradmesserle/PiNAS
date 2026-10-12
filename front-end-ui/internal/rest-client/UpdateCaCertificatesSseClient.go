package rest_client

import (
	"time"
)

// UpdateCaCertificates retrieves and processes CA certificate updates from an SSE stream endpoint.
// It handles the HTTP request, monitors the stream, and publishes parsed events via an internal event bus.
// The function uses a wait group to signal completion and enforces a timeout for the operation.
func UpdateCaCertificates() error {
	return ExecuteGetSseCall("http://localhost:9090/updateCaCertificates", 30*time.Minute)
}
