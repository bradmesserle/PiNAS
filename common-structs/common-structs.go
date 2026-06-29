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
}
