package structs

import "github.com/pinas/common-structs"

type SystemInfo struct {
	Model       common_structs.PIModel
	pcieEnabled bool
	cpuInfo     common_structs.CpuInfo
}
