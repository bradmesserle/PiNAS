package structs

import (
	"encoding/json"
	"log/slog"
	"os"
)

type WizardInfo struct {
	Step       int
	NasOptions NasOptions
}

type NasOptions struct {
	InstallDns    bool
	InstallCa     bool
	InstallNvmeFa bool
}

type SetupInfo struct {
	IsRunning bool
}

// SaveSetupOptions saves the setup options and returns an error if it fails.
func SaveSetupOptions(options NasOptions) error {

	jsonData, err := json.MarshalIndent(options, "", "    ")
	if err != nil {
		slog.Error("Error marshaling to JSON", "error", err)
		return err
	}

	err = os.WriteFile("/opt/pinas/ui-services/options.json", jsonData, 0644)
	if err != nil {
		slog.Error("Error writing file: %v\n", "error", err)
		return err
	}

	return nil
}

// GetSetupOptions configures the NAS setup options based on the provided NasOptions and returns an error if the setup fails.
func GetSetupOptions() (NasOptions, error) {
	var options NasOptions

	fileBytes, err := os.ReadFile("/opt/pinas/ui-services/options.json")
	if err != nil {
		slog.Error("Failed to read file: %v", err)
	}

	err = json.Unmarshal(fileBytes, &options)
	if err != nil {
		slog.Error("Failed to unmarshal JSON: %v", err)
		return options, err
	}

	return options, nil

}
