package utilities

import (
	"bufio"
	"log"
	"net/http"
	"os/exec"
)

// ExecCmdSseStdoutText Start the command and sends the standard output of a command to the client as Server-Sent Events.
func ExecCmdSseStdoutText(cmd *exec.Cmd, w http.ResponseWriter) error {

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		if err := SendEventData(scanner.Text(), "consoleOutput", w); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	return nil

}

// SendEventData Send the event data to the client
func SendEventData(data string, eventId string, w http.ResponseWriter) error {

	event := Event{
		Data:  []byte(data),
		Event: []byte(eventId),
	}
	if err := event.MarshalTo(w); err != nil {
		return err
	}
	if err := http.NewResponseController(w).Flush(); err != nil {
		return err
	}

	return nil
}
