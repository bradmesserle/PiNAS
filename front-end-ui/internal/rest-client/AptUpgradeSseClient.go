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

	"github.com/pinas/ui/internal"
)

func AptUpgrade(wg *sync.WaitGroup) error {

	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://192.168.0.22:9090/aptUpgrade", nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//log.Fatalf("Failed to connect: %v", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//log.Fatalf("Failed to close body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		//log.("Unexpected status code: %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	var currentData strings.Builder

	fmt.Println("Connected to SSE stream...")
	for scanner.Scan() {
		line := scanner.Text()

		// Empty line denotes the end of an event block
		if line == "" {
			if currentData.Len() > 0 {
				internal.EventBus.Publish("consoleLog", strings.TrimSpace(currentData.String()))
				//fmt.Printf("%s\n", strings.TrimSpace(currentData.String()))
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
