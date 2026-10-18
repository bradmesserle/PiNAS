package install_features

import (
	"log/slog"
	"net/http"
	"os/exec"

	common_structs "github.com/pinas/common-structs"
	"github.com/pinas/rest-services/internal/utilities"
)

func InstallBind9DNS(dnsSettings common_structs.DnsSettings, w http.ResponseWriter) error {

	// Install Bind DNS
	installDnsErr := installDns(w)
	if installDnsErr != nil {
		return nil
	}

	//Install DNS Utils
	installDnsUtilErr := installDnsUtils(w)
	if installDnsUtilErr != nil {
		return nil
	}

	//Install Systemd Resolve
	installSystemdResolveErr := installSystemdResolve(w)
	if installSystemdResolveErr != nil {
		return nil
	}

	// Updating /etc/bind permissions
	updateEtcBindErr := updateEtcBindPermissions(w)
	if updateEtcBindErr != nil {
		return nil
	}

	// Updating /var/cache/bind permissions
	updateVarCacheErr := updateVarCacheBindPermissions(w)
	if updateVarCacheErr != nil {
		return nil
	}

	// Generate DNS Key
	dnsKeyErr := generateDnsKey(w)
	if dnsKeyErr != nil {
		return nil
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

// generateDnsKey generates a DNS key, returning an error if the operation fails.
func generateDnsKey(w http.ResponseWriter) error {

	if err := utilities.SendEventData("Generating DNS Key", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	// tsig-keygen -a hmac-sha256 externaldns
	cmd := exec.Command("tsig-keygen", "-a", "hmac-sha256", "externaldns")
	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while creating DNS key", "Value", err)
		return err
	}

	return nil
}
