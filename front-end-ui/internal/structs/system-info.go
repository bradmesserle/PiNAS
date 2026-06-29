package structs

import "github.com/pinas/common-structs"

type SystemInfo struct {
	Model       common_structs.PIModel
	PcieEnabled bool
	CpuInfo     common_structs.CpuInfo
	MemoryInfo  common_structs.MemoryInfo
	PcieInfo    common_structs.PCIeInfo
}
