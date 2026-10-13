package install_features

import (
	"log/slog"
	"net/http"
	"os/exec"

	"github.com/pinas/rest-services/internal/utilities"
)

// InstallStepCA sets up the SmallStep Certificate Authority by installing OS packages, adding the repository, and installing required tools.
func InstallStepCA(w http.ResponseWriter) error {

	// OS Dependencies
	osDepErr := installOsPackages(w)
	if osDepErr != nil {
		return osDepErr
	}

	// Add SmallStep Repo to system
	repoErr := installSmallStepAptRepo(w)
	if repoErr != nil {
		return repoErr
	}

	// Install  installSmallStep
	installErr := installSmallStep(w)
	if installErr != nil {
		return installErr
	}

	return nil
}

// installOsPackages installs required OS packages and streams progress updates via Server-Sent Events to the client.
func installOsPackages(w http.ResponseWriter) error {

	if err := utilities.SendEventData("Installing Small Step CA Dependencies", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	cmd := exec.Command("apt", "install", "-y", "--no-install-recommends", "curl", "vim", "gpg", "ca-certificates")
	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while installing SmallStep CA Dependencies", "Value", err)
		return err
	}

	return nil

}

// curl -fsSL https://packages.smallstep.com/keys/apt/repo-signing-key.gpg -o /etc/apt/trusted.gpg.d/smallstep.asc && echo 'deb [signed-by=/etc/apt/trusted.gpg.d/smallstep.asc] https://packages.smallstep.com/stable/debian debs main'  | tee /etc/apt/sources.list.d/smallstep.list

// installSmallStepAptRepo installs the SmallStep APT repository on a Debian-based system and streams updates via HTTP response.
// It uses utilities to send SSE events and execute commands while providing feedback to the client.
// Returns an error if the repository setup or related operations fail.
func installSmallStepAptRepo(w http.ResponseWriter) error {

	if err := utilities.SendEventData("Installing STEP CA APT Repo", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	cmd := exec.Command("curl", "-fsSL", " https://packages.smallstep.com/keys/apt/repo-signing-key.gpg", "-o", "/etc/apt/trusted.gpg.d/smallstep.asc", "&&", "echo", "'deb [signed-by=/etc/apt/trusted.gpg.d/smallstep.asc] https://packages.smallstep.com/stable/debian debs main'", "|", "tee", "/etc/apt/sources.list.d/smallstep.list")
	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while installing CA APT Repo", "Value", err)
		return err
	}

	return nil
}

// installSmallStep installs the Step CA by updating package lists and executing the installation, sending progress as SSE.
func installSmallStep(w http.ResponseWriter) error {

	if err := utilities.SendEventData("Installing STEP CA", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	cmd := exec.Command("apt-get", "update", "&&", "apt-get", "-y", "install", "step-cli", "step-ca")
	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while Installing Step CA", "Value", err)
		return err
	}

	return nil
}
