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
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	stdOutScanner := bufio.NewScanner(stdout)
	for stdOutScanner.Scan() {
		if err := SendEventData(stdOutScanner.Text(), "consoleOutput", w); err != nil {
			return err
		}
	}

	stdErrScanner := bufio.NewScanner(stderr)
	for stdErrScanner.Scan() {
		if err := SendEventData(stdErrScanner.Text(), "consoleErrOutput", w); err != nil {
			return err
		}
	}

	if err := stdOutScanner.Err(); err != nil {
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
