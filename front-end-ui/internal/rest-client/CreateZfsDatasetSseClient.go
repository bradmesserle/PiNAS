package rest_client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/pinas/common-structs"
	"github.com/pinas/ui/internal"
)

// CreateZfsDataset  Create ZFS Dataset
func CreateZfsDataset(wg *sync.WaitGroup, zfsDataset common_structs.ZfsDataset) (err error) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Marshal struct to JSON
	jsonData, _ := json.Marshal(zfsDataset)

	req, err := http.NewRequestWithContext(ctx, "POST", "http://192.168.0.22:9090/createZfsDataset", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to connect: %v", err)
	}

	//Check to see if we have a response.
	if resp == nil {
		return nil
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close body: %v", err)
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
