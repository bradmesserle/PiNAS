package common_structs

type CpuInfo struct {
	NumberOfCpus int    `json:"numberOfCpus"`
	BogoMIPS     string `json:"bogoMIPS"`
	Architecture string `json:"architecture"`
	Revision     string `json:"revision"`
	Features     string `json:"features"`
}

type MemoryInfo struct {
	Total     string `json:"total"`
	Free      string `json:"free"`
	Available string `json:"available"`
	Buffers   string `json:"buffers"`
	Cached    string `json:"cached"`
	SwapTotal string `json:"swapTotal"`
	SwapFree  string `json:"swapFree"`
}

type PCIeInfo struct {
	PCIeDevices         []PCIeDevice `json:"nvmeDrives"`
	GenVersion          string       `json:"genVersion"`
	Gen3EnabledInConfig bool         `json:"gen3EnabledInConfig"`
}

type PCIeDevice struct {
	Slot        string `json:"slot"`
	Vendor      string `json:"vendor"`
	Description string `json:"description"`
	Speed       string `json:"speed"`
}

type PIModel struct {
	PiModel string `json:"piModel"`
}

type NvmeDrive struct {
	DeviceIdentifier string `json:"deviceIdentifier"`
	Name             string `json:"name"`
	Model            string `json:"model"`
	Size             string `json:"size"`
	Serial           string `json:"serial"`
}

type InstallOptions struct {
	InstallDns    bool `json:"installDns"`
	InstallCa     bool `json:"installCa"`
	InstallNvmeFa bool `json:"installNvmeFa"`
}

type InstallStatus struct {
	DnsInstalled    bool   `json:"dnsInstalled"`
	CaInstalled     bool   `json:"caInstalled"`
	NvmeFaInstalled bool   `json:"nvmeFaInstalled"`
	SystemUpdated   bool   `json:"systemUpdated"`
	KernelCompiled  bool   `json:"kernelCompiled"`
	CurrentProcess  string `json:"currentProcess"`
	InstallComplete bool   `json:"installComplete"`
}

type ZfsPool struct {
	PoolName       string      `json:"poolName"`
	WifeFilesystem bool        `json:"wipeFilesystem"`
	NvmeDrives     []NvmeDrive `json:"nvmeDrives"`
}
