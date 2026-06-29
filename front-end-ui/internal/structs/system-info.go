package structs

import "github.com/pinas/common-structs"

type SystemInfo struct {
	Model       common_structs.PIModel
	PcieEnabled bool
	CpuInfo     common_structs.CpuInfo
}
