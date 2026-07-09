package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pinas/rest-services/internal/endpoint-functions/zfs"
	"github.com/pinas/rest-services/internal/kernel-compile"
	"github.com/pinas/rest-services/internal/os-functions"
	"github.com/pinas/rest-services/internal/system-info"
	"github.com/pinas/rest-services/internal/system-updates"
	"github.com/pinas/rest-services/internal/utilities"
)

func main() {

	var e = echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	//TODO: Endpoints
	// Verify and Update config.txt enable pcie gen3 - Done
	// Compile the linux kernel -
	//     a. download linux - Done
	//     b. update the kernel config file --> enable nvme-fa options -Done
	//     c. stretch goal - create a slim version of the kernel config txt file.. get rid of stuff we dont need.
	// Check see if any header packages are installed if so removed them - this messes up zfs install -Done
	// Install ZFS -Done
	// Get Drive Telemetry - get the installed nvme drive data with serial number so we can build the pool from drive SN - Done
	// Create PiNAS working partition(1-5G?).. We will need a working partition to store data, install dns, step-ca and move /etc off the micro-sd
	// Create zpool api - Done
	// Create datasets - Need to look at the options and support what we need. nvme-fa block storage we will need
	// Restart the system - Done
	// Wipe Drives Clean - for before creating pools - Done
	// Save the installation progress to a file so we can restart from where we left off
	// Create endpoint to install DNS Server
	// Create endpoint to install Step CA

	//Get CPU Info
	e.GET("/cpuInfo", system_info.GetCpuInfo)

	//Get Model Info
	e.GET("/modelInfo", system_info.GetPiModelInfo)

	//Get Memory Info
	e.GET("/memoryInfo", system_info.GetMemoryInfo)

	//Get PCIe Info
	e.GET("/pcieInfo", system_info.GetPcieInfo)

	//Verify the config.txt has pcie gen3 enabled. If not, update the config.txt file
	e.GET("/verifyUpdateConfig", system_updates.VerifyUpdateConfig)

	//Compile linux kernel we need to enable the nvme-fa options
	e.GET("/compileKernel", kernel_compile.CompileLinuxKernel)

	//Install ZFS
	e.GET("/installZFS", zfs.InstallZfs)

	//Re-Install zfs-dkms
	e.GET("/reinstallZFSDkms", zfs.ReinstallZfsDkms)

	//Get Drive Telemetry
	e.GET("/getDriveTelemetry", system_info.GetDriveTelemetry)

	//Create ZFS Pool
	e.POST("/createZfsPool", zfs.CreateZfsPool)

	//Create ZFS Dataset
	e.POST("/createZfsDataset", zfs.CreateZfsDataset)

	//Move the /etc directory off the micro-sd card
	e.GET("/moveEtcDirectory", system_updates.MoveEtcDirectory)

	//Apt Update
	e.GET("/aptUpdate", system_updates.AptUpdate)

	//Apt Upgrade
	e.GET("/aptUpgrade", system_updates.AptUpgrade)

	//Health Check
	e.GET("/healthCheck", utilities.GetHealthCheck)

	//Reboot the system
	e.GET("/reboot", os_functions.Reboot)

	// Start the server
	sc := echo.StartConfig{
		Address: ":9090",
		BeforeServeFunc: func(s *http.Server) error {
			s.WriteTimeout = 0 // IMPORTANT: disable for SSE
			return nil
		},
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM) // start shutdown process on ctrl+c
	defer cancel()

	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
