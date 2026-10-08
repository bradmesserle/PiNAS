package structs

type Checkpoint string

const (
	FreshInstall Checkpoint = ""
	ZfsInstalled Checkpoint = "zfs"
	ZfsCreated   Checkpoint = "zfs-created"
	NvmeFaPart1  Checkpoint = "nvme-fa-part1"
	PartDone     Checkpoint = "part-done"
	Completed    Checkpoint = "completed"
)
