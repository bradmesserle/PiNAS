package rest_client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/starfederation/datastar-go/datastar"
)

type Signals struct {
	Content   string `json:"content"`
	Streaming bool   `json:"streaming"`
}

func AptUpdate(wg *sync.WaitGroup) error {

	defer wg.Done()

	// 1. Setup a cancellable context to close the stream when needed
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// 2. Create the long-lived HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", "http://192.168.0.22:9090/aptUpdate", nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
	}

	// 3. Set the required SSE headers
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	// 4. Execute request using default client (ensure no response timeouts are set)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	//http.HandleFunc("/aptUpdate", streamHandler)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Failed to close body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	// 5. Scan the body line by line
	scanner := bufio.NewScanner(resp.Body)
	var currentData strings.Builder

	fmt.Println("Connected to SSE stream...")
	for scanner.Scan() {
		line := scanner.Text()

		// Empty line denotes the end of an event block
		if line == "" {
			if currentData.Len() > 0 {
				fmt.Printf("%s\n", strings.TrimSpace(currentData.String()))
				currentData.Reset()
			}
			continue
		}

		// Parse the SSE fields (data:, event:, id:, retry:)
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}

		key := parts[0]
		value := strings.TrimSpace(parts[1])

		switch key {
		case "data":
			currentData.WriteString(value + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Stream error: %v", err)
	}

	return nil
}

func StreamHandler(c *echo.Context) error {

	var in Signals
	_ = datastar.ReadSignals(c.Request(), &in)

	// NewSSE sets the SSE headers and returns a generator bound to this request.
	sse := datastar.NewSSE(c.Response(), c.Request())

	// Flip the `streaming` signal on so the buttons disable and the status
	// indicator lights up. This is a datastar-patch-signals SSE event.
	_ = sse.MarshalAndPatchSignals(map[string]any{"streaming": true})

	message := `Datastar streams this text one token at a time, straight from ` +
		`the Go server into the textarea using Server-Sent Events. Each chunk ` +
		`patches the "output" signal, and because the textarea is bound with ` +
		`data-bind, its value updates live — no custom JavaScript required.`

	//var b strings.Builder
	for _, tok := range strings.Fields(message) {
		// Stop early if the browser closes the connection (tab closed, etc.).
		select {
		case <-c.Request().Context().Done():
			return nil
		default:
		}

		//if b.Len() > 0 {
		//	b.WriteByte(' ')
		//}
		//b.WriteString(tok)

		//scriptText = fmt.Printf("%s",tok)
		sse.ExecuteScript(fmt.Sprintf(`updateText("%s")`, tok))
		//sse.ExecuteScript(`console.log(bodyText)`)
		//if err != nil {
		//	//return err
		//}

		// Patch just the `output` signal with the accumulated text so far.
		// A signals patch is a merge, so `streaming` is left untouched.
		//if err := sse.MarshalAndPatchSignals(map[string]any{"output": b.String()}); err != nil {
		//	return nil // client went away
		//}

		time.Sleep(150 * time.Millisecond)
	}

	// Done: turn the `streaming` signal back off.
	return sse.MarshalAndPatchSignals(map[string]any{"streaming": false})

}
