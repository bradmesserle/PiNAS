package structs

type Checkpoint string

const (
	FreshInstall Checkpoint = ""
	ZfsInstalled Checkpoint = "zfs"
	NvmeFaPart1  Checkpoint = "nvme-fa-part1"
	NvmeFaPart2  Checkpoint = "nvme-fa-part2"
	Completed    Checkpoint = "completed"
)
