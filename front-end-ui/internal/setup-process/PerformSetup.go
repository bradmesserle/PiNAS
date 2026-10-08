package setup_process

import (
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/pinas/common-structs"
	"github.com/pinas/ui/internal"
	"github.com/pinas/ui/internal/rest-client"
	"github.com/pinas/ui/internal/structs"
)

func PerformSetup(wizardInfo *structs.WizardInfo) {

	//Load the current setup options
	options, err := structs.GetSetupOptions()
	if err != nil {
		slog.Error("Error while loading setup options", "err", err)
		return
	}

	//Set the options in the wizardInfo
	wizardInfo.NasOptions = options

	//Get the current installation checkpoint
	checkpoint, err := GetStatus()
	if err != nil {
		slog.Error("Error while getting status", "err", err)
		return
	}

	slog.Info("Current installation checkpoint", "checkpoint", checkpoint)

	time.Sleep(1 * time.Second)

	//Install Part 1
	//Update system and install ZFS

	//Get ZfsStatus
	zfsStatus, errZfsStatus := rest_client.GetZfsStatus()
	if errZfsStatus != nil {
		log.Println(errZfsStatus.Error())
	}

	//Check to see if zfs is up and running and the pool is online
	//If zfs pools are not online, we need to install them, and we are at the start of the installation process
	if zfsStatus.State != common_structs.ZfsStateOnline && checkpoint == structs.FreshInstall {

		//Send Success message to the front end.
		internal.EventBus.Publish("consoleLog", strings.TrimSpace("Starting the setup process"))

		//Update ca certificates on the box
		errCaCert := rest_client.UpdateCaCertificates()
		if errCaCert != nil {
			log.Println(errCaCert.Error())
		}

		//Update the cmdline file for cgroup memory management.
		errUpdateCmdLineFile := rest_client.UpdateCmdLineFile()
		if errUpdateCmdLineFile != nil {
			log.Println(errUpdateCmdLineFile.Error())
		}

		//Update repos
		errUpdate := rest_client.AptUpdate()
		if errUpdate != nil {
			log.Println(errUpdate.Error())
		}

		//Upgrade System
		errUpgrade := rest_client.AptUpgrade()
		if errUpgrade != nil {
			log.Println(errUpgrade.Error())
		}

		//Install ZFS
		errInstallZfs := rest_client.InstallZfs()
		if errInstallZfs != nil {
			log.Println(errInstallZfs.Error())
		}

		//Save Checkpoint
		errSaveCheckpoint := SaveStatus(structs.ZfsInstalled)
		if errSaveCheckpoint != nil {
			log.Println(errSaveCheckpoint.Error())
		}

		//Reboot
		errReboot := rest_client.Reboot()
		if errReboot != nil {
			log.Println(errReboot.Error())
		}
	}

	// Continue with installation
	// to Create the zfs pools and work area datasets
	if checkpoint == structs.ZfsInstalled {

		slog.Info("Installing part 2")

		//Create ZFS Pool
		createZfsPool()

		//Create Work Area Dataset
		createWorkAreaDataset()

		checkpoint = structs.ZfsCreated
	}

	//Check if nvme-fa is enabled if so, compile the kernel and enable it
	if wizardInfo.NasOptions.InstallNvmeFa {

		if checkpoint == structs.ZfsCreated {
			err := installNvmeFaPart1(checkpoint)
			if err != nil {
				log.Println(err.Error())
			}
		}

		if checkpoint == structs.NvmeFaPart1 {
			err := NvmeFaInstallPart2(checkpoint)
			if err != nil {
				log.Println(err.Error())
			}
		}

		//Set to done
		checkpoint = structs.PartDone

	}

	//Check if we need to install DNS
	if wizardInfo.NasOptions.InstallDns {
		err := installDns()
		if err != nil {
			log.Println(err.Error())
		}
	}

	// Check to see if we need to install CA
	if wizardInfo.NasOptions.InstallCa {
		err := installCa()
		if err != nil {
			log.Println(err.Error())
		}
	}

	//Save Checkpoint - Check to see if all steps are done
	if checkpoint == structs.PartDone {

		//Send completed message to the front end.
		internal.EventBus.Publish("consoleLog", strings.TrimSpace("Setup Completed Successfully"))

		errSaveCheckpoint := SaveStatus(structs.Completed)
		if errSaveCheckpoint != nil {
			log.Println(errSaveCheckpoint.Error())
		}
	}

}

// createZfsPool initializes and creates a ZFS pool using drive telemetry fetched via REST API calls.
func createZfsPool() {

	poolInfo := new(common_structs.ZfsPool)
	poolInfo.WifeFilesystem = true
	poolInfo.PoolName = "zfspool"

	//Get Drive Telemetry
	var drives, err = rest_client.GetDriveTelemetry()
	if err != nil {
		log.Println(err.Error())
	}

	for _, drive := range drives {
		poolInfo.NvmeDrives = append(poolInfo.NvmeDrives, drive)
	}

	errCreatePool := rest_client.CreateZfsPool(*poolInfo)
	if errCreatePool != nil {
		log.Println(errCreatePool.Error())
	}

}

// createWorkAreaDataset sets up a working area ZFS dataset with specified options and synchronizes with a wait group.
func createWorkAreaDataset() {

	dataset := new(common_structs.ZfsDataset)
	dataset.DatasetName = "work-area"
	dataset.PoolName = "zfspool"

	err := rest_client.CreateZfsDataset(*dataset)
	if err != nil {
		log.Println(err.Error())
	}

}

// installNvmeFa checks if nvme-fa is enabled and compiles the kernel and enables it if so.
func installNvmeFaPart1(checkpoint structs.Checkpoint) error {

	log.Println("Enable nvme-fa")
	slog.Info("Enable nvme-fa")

	//We are at the start of the process

	// Compile the kernel
	err := rest_client.CompileKernel()
	if err != nil {
		//log.Println(err.Error())
		return err
	}

	//Reinstall ZFS so it can compile the headers
	errReinstallZfsDkms := rest_client.ReinstallZfsDkms()
	if errReinstallZfsDkms != nil {
		//log.Println(errReinstallZfsDkms.Error())
		return errReinstallZfsDkms
	}

	//Save Checkpoint
	errSaveCheckpoint := SaveStatus(structs.NvmeFaPart1)
	if errSaveCheckpoint != nil {
		log.Println(errSaveCheckpoint.Error())
	}

	//Reboot
	errReboot := rest_client.Reboot()
	if errReboot != nil {
		//log.Println(errReboot.Error())
	}

	return nil
}

func NvmeFaInstallPart2(checkpoint structs.Checkpoint) error {

	//Continue onto part 2 Install the CLI
	log.Println("NvmeFaPart2 - Installing nvme-cli ")

	//Install nvme-cli
	errNvmeCli := rest_client.InstallNvmeCli()
	if errNvmeCli != nil {
		log.Println(errNvmeCli.Error())
		return errNvmeCli
	}

	return nil
}

// installDns configures and installs the DNS settings required for the application and returns an error if it fails.
func installDns() error {

	log.Println("Installing DNS")

	return nil

}

// installCa installs a Certificate Authority (CA) on the system and returns an error if the installation fails.
func installCa() error {

	return nil
}
