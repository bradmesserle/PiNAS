package endpoints

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v5"
	"github.com/pinas/ui/internal"
	"github.com/starfederation/datastar-go/datastar"
)

type Signals struct {
	Content   string `json:"content"`
	Streaming bool   `json:"streaming"`
}

// ConsoleLogStreamHandler streams console log messages to the client using Server-Sent Events (SSE).
// It initializes SSE, enables the streaming signal, subscribes to the "consoleLog" event, and streams updates.
func ConsoleLogStreamHandler(c *echo.Context) error {

	var in Signals
	_ = datastar.ReadSignals(c.Request(), &in)

	// NewSSE sets the SSE headers and returns a generator bound to this request.
	sse := datastar.NewSSE(c.Response(), c.Request())

	// Flip the `streaming` signal on so front end can process the streaming data.
	// This is a datastar-patch-signals SSE event.
	_ = sse.MarshalAndPatchSignals(map[string]any{"streaming": true})

	//Subscribe to the topic and post on the stream
	_ = internal.EventBus.Subscribe("consoleLog", func(msg string) {
		fmt.Printf("Receiving Data --->: %s\n", msg)
		err := sse.ExecuteScript(fmt.Sprintf(`updateText("%s")`, msg))
		if err != nil {
			log.Println(err)
		}
	})

	for {
		select {
		case <-c.Request().Context().Done():
			return sse.MarshalAndPatchSignals(map[string]any{"streaming": false})
		}
	}

}
