package rest_client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"time"

	"github.com/pinas/ui/internal"
)

func AptUpdate() error {

	//defer wg.Done()

	// 1. Setup a cancellable context to close the stream when needed
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
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
		//log.Fatalf("Failed to connect: %v", err)
	}

	if resp == nil {
		log.Println("Response body is nil")
		return nil
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//log.Fatalf("Failed to close body: %v", err)
			log.Printf("Failed to close body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		//log.Fatalf("Unexpected status code: %d", resp.StatusCode)
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
				internal.EventBus.Publish("consoleLog", strings.TrimSpace(currentData.String()))
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
