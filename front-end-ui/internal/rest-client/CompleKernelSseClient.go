package rest_client

import (
	"time"
)

func CompileKernel() error {
	return ExecuteGetSseCall("http://localhost:9090/compileKernel", 2*time.Hour)
}
