package setup_process

import (
	"log"
	"strings"
	"time"

	"github.com/pinas/common-structs"
	"github.com/pinas/ui/internal"
	"github.com/pinas/ui/internal/rest-client"
	"github.com/pinas/ui/internal/structs"
)

func PerformSetup(wizardInfo *structs.WizardInfo) {

	time.Sleep(1 * time.Second)

	//Install Part 1
	//Update system and install ZFS

	//Send Success message to the front end.
	internal.EventBus.Publish("consoleLog", strings.TrimSpace("Starting the setup process"))

	//Get ZfsStatus
	zfsStatus, errZfsStatus := rest_client.GetZfsStatus()
	if errZfsStatus != nil {
		log.Println(errZfsStatus.Error())
	}

	//Check to see if zfs is up and running and the pool is online
	//If zfs pools are not online we need to install them
	if zfsStatus.State != common_structs.ZfsStateOnline {

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

		//Reboot
		errReboot := rest_client.Reboot()
		if errReboot != nil {
			log.Println(errReboot.Error())
		}

		//Wait for the reboot to complete
		waitForReboot()
		internal.EventBus.Publish("consoleLog", strings.TrimSpace("System is back up"))

		//Install Part 2
		// Create the ZFS Pool, standard work-area dataset.
		// Install Setup Options

		//Create ZFS Pool
		createZfsPool()

		//Create Work Area Dataset
		createWorkAreaDataset()
	}

	//Check if nvme-fa is enabled if so compile the kernel and enable it
	if wizardInfo.NasOptions.InstallNvmeFa {
		err := installNvmeFa()
		if err != nil {
			log.Println(err.Error())
		}
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

	//Send Success message to the front end.
	internal.EventBus.Publish("consoleLog", strings.TrimSpace("Setup Completed Successfully"))

}

func waitForReboot() {

	//defer wg.Done()
	rebootStarted := false

	//Check the health check
	for {

		//Check to see if we have an error if show we know the system is rebooting..
		if !rebootStarted {
			errReboot := rest_client.HealthCheck()
			if errReboot != nil {
				rebootStarted = true
			}
		} else {

			//Lets wait until we dont get an error
			errReboot := rest_client.HealthCheck()
			if errReboot == nil {
				log.Println("Reboot is finished")
				break
			}

		}

		time.Sleep(1 * time.Second)
	}

}

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
func installNvmeFa() error {

	log.Println("Enable nvme-fa")

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

	//Reboot
	errReboot := rest_client.Reboot()
	if errReboot != nil {
		//log.Println(errReboot.Error())
	}

	//Wait for the reboot to complete
	waitForReboot()
	internal.EventBus.Publish("consoleLog", strings.TrimSpace("System is back up"))

	//Install nvme-cli
	errNvmeCli := rest_client.InstallNvmeCli()
	if errNvmeCli != nil {
		//log.Println(errNvmeCli.Error())
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
