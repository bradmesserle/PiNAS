package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	kernel_compile "github.com/pinas/rest-services/internal/kernel-compile"
	"github.com/pinas/rest-services/internal/system-info"
	"github.com/pinas/rest-services/internal/system-updates"
	zfs_install "github.com/pinas/rest-services/internal/zfs-install"
)

func main() {

	e := echo.New()
	e.Use(middleware.RequestLogger())

	//TODO: Endpoints
	// Verify and Update config.txt enable pcie gen3
	// Compile the linux kernel -
	//     a. download linux
	//     b. update the kernel config file --> enable nvme-fa options
	//     c. stretch goal - create a slim version of the kernel config txt file.. get rid of stuff we dont need.
	// Check see if any header packages are installed if so removed them - this messes up zfs install
	// Install ZFS
	// Get Drive Telemetry - get the installed nvme drive data with serial number so we can build the pool from drive SN
	// Create PiNAS working partition(1-5G?).. We will need a working partition to store data, install dns, step-ca and move /etc off the micro-sd
	// Create zpool api
	// Create datasets - Need to look at the options and support what we need. nvme-fa block storage we will need
	// Restart the system

	e.GET("/cpuInfo", system_info.GetCpu)

	//Verify the config.txt has pcie gen3 enabled. If not, update the config.txt file
	e.GET("/verifyUpdateConfig", system_updates.VerifyUpdateConfig)

	//Compile linux kernel we need to enable the nvme-fa options
	e.GET("/compileKernel", kernel_compile.CompileLinuxKernel)

	//Install ZFS
	e.GET("/installZFS", zfs_install.InstallZfs)

	//Get Drive Telemetry
	e.GET("/getDriveTelemetry", system_info.GetDriveTelemetry)

	//Move the /etc directory off the micro-sd card
	e.GET("/moveEtcDirectory", system_updates.MoveEtcDirectory)

	// Start the server
	e.Logger.Fatal(e.Start("localhost:9090"))
}
