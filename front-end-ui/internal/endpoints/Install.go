package endpoints

import (
	"log"
	"sync"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
	"github.com/pinas/common-structs"
	"github.com/pinas/ui/internal/components/setup"
	"github.com/pinas/ui/internal/rest-client"
	"github.com/pinas/ui/internal/structs"
)

func Install(c *echo.Context, wizardInfo *structs.WizardInfo) error {

	//Kick off the installation process
	var wg sync.WaitGroup
	wg.Add(15)

	var cmp templ.Component = setup.InstallProgressPage(*wizardInfo)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	renderPage := cmp.Render(c.Request().Context(), c.Response())

	go func() {

		//Install Part 1
		//Update system and install ZFS

		//Update repos
		errUpdate := rest_client.AptUpdate(&wg)
		if errUpdate != nil {
			log.Println(errUpdate.Error())
		}

		//Upgrade System
		errUpgrade := rest_client.AptUpgrade(&wg)
		if errUpgrade != nil {
			log.Println(errUpgrade.Error())
		}

		//Install ZFS
		errInstallZfs := rest_client.InstallZfs(&wg)
		if errInstallZfs != nil {
			log.Println(errInstallZfs.Error())
		}

		//Reboot
		errReboot := rest_client.Reboot(&wg)
		if errReboot != nil {
			log.Println(errReboot.Error())
		}

		//Wait for the reboot to complete
		waitForReboot(&wg)

		//Install Part 2
		// Create the ZFS Pool, standard work-area dataset.

		//Create ZFS Pool
		createZfsPool(&wg)

		//Create Work Area Dataset
		createWorkAreaDataset(&wg)

		//Check if nvme-fa is enabled if so compile the kernel and enable it
		if wizardInfo.NasOptions.InstallNvmeFa {
			err := installNvmeFa(&wg)
			if err != nil {
				log.Println(errInstallZfs.Error())
			}
		}

	}()

	return renderPage
}

func waitForReboot(wg *sync.WaitGroup) {

	defer wg.Done()
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

func createZfsPool(wg *sync.WaitGroup) {
	defer wg.Done()

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

	errCreatePool := rest_client.CreateZfsPool(wg, *poolInfo)
	if errCreatePool != nil {
		log.Println(errCreatePool.Error())
	}

}

// createWorkAreaDataset sets up a working area ZFS dataset with specified options and synchronizes with a wait group.
func createWorkAreaDataset(wg *sync.WaitGroup) {
	defer wg.Done()

	dataset := new(common_structs.ZfsDataset)
	dataset.DatasetName = "work-area"
	dataset.PoolName = "zfspool"

	err := rest_client.CreateZfsDataset(wg, *dataset)
	if err != nil {
		log.Println(err.Error())
	}

}

// installNvmeFa checks if nvme-fa is enabled and compiles the kernel and enables it if so.
func installNvmeFa(wg *sync.WaitGroup) error {
	defer wg.Done()
	log.Println("Enable nvme-fa")

	// Compile the kernel
	err := rest_client.CompileKernel(wg)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	//Reinstall ZFS so it can compile the headers
	errReinstallZfsDkms := rest_client.ReinstallZfsDkms(wg)
	if errReinstallZfsDkms != nil {
		log.Println(errReinstallZfsDkms.Error())
		return errReinstallZfsDkms
	}

	//Reboot
	errReboot := rest_client.Reboot(wg)
	if errReboot != nil {
		log.Println(errReboot.Error())
	}

	//Wait for the reboot to complete
	waitForReboot(wg)

	return nil

}
