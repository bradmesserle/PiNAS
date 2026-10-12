package install_features

import (
	"log/slog"
	"net/http"
	"os/exec"

	"github.com/pinas/rest-services/internal/utilities"
)

func InstallBind9DNS(w http.ResponseWriter) error {

	// Install Bind DNS
	installDnsErr := installDns(w)
	if installDnsErr != nil {
		return installDnsErr
	}

	//Install DNS Utils
	installDnsUtilErr := installDnsUtils(w)
	if installDnsUtilErr != nil {
		return installDnsUtilErr
	}

	//Install Systemd Resolve
	installSystemdResolveErr := installSystemdResolve(w)
	if installSystemdResolveErr != nil {
		return installSystemdResolveErr
	}

	// Updating /etc/bind permissions
	updateEtcBindErr := updateEtcBindPermissions(w)
	if updateEtcBindErr != nil {
		return updateEtcBindErr
	}

	return nil
}

func installDns(w http.ResponseWriter) error {

	if err := utilities.SendEventData("Installing DNS", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	cmd := exec.Command("apt", "install", "bind9", "bind9utils", "bind9-doc", "-y")
	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while installing DNS", "Value", err)
		return err
	}

	return nil

}

// installDnsUtils installs DNS Utils.
func installDnsUtils(w http.ResponseWriter) error {

	if err := utilities.SendEventData("Installing DNS Utils", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	cmd := exec.Command("apt", "install", "dnsutils", "-y")
	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while installing DNS Utils", "Value", err)
		return err
	}

	return nil

}

// installSystemdResolve installs the systemd-resolved package and streams output to the HTTP client.
func installSystemdResolve(w http.ResponseWriter) error {

	if err := utilities.SendEventData("Installing SsystemD Resolve", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	cmd := exec.Command("apt", "install", "systemd-resolved", "-y")
	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while installing DNS Utils", "Value", err)
		return err
	}

	return nil

}

// updateEtcBindPermissions updates the ownership of the /etc/bind directory to the user and group "bind".
func updateEtcBindPermissions(w http.ResponseWriter) error {

	if err := utilities.SendEventData("Updating permissions for /etc/bind", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	cmd := exec.Command("chown", "-R", "bind:bind", "/etc/bind")
	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while updating permissions for /etc/bind", "Value", err)
		return err
	}

	return nil

}

// updateVarCacheBindPermissions updates the ownership of the /var/cache/bind directory to the user and group "bind".
func updateVarCacheBindPermissions(w http.ResponseWriter) error {

	if err := utilities.SendEventData("Updating permissions for /var/cache/bind", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	cmd := exec.Command("chown", "-R", "bind:bind", "/var/cache/bind")
	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while updating permissions for /var/cache/bind", "Value", err)
		return err
	}

	return nil

}

// createLogDir ensures the creation of a log directory, returning an error if the operation fails.
func createLogDir(w http.ResponseWriter) error {

	if err := utilities.SendEventData("Creating log directory", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	cmd := exec.Command("mkdir", "-p", "/var/log/bind")
	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while creating log directory - /var/log/bind", "Value", err)
		return err
	}

	return nil

}
