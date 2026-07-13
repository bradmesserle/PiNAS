package structs

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
