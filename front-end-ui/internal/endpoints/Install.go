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
	wg.Add(7)

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

		//Install Part 2
		//Wait for the reboot to complete
		waitForReboot(&wg)

		//Create ZFS Pool
		createZfsPool(&wg)

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
	poolInfo.PoolName = "ZfsPool"

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
