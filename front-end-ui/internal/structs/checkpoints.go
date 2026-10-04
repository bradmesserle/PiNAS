package structs

type Checkpoint string

const (
	FreshInstall Checkpoint = ""
	ZfsInstalled Checkpoint = "zfs"
)
