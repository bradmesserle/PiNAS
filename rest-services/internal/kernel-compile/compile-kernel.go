package kernel_compile

import (
	"bytes"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/exec"

	"github.com/labstack/echo/v5"
	"github.com/pinas/rest-services/internal/utilities"
)

var buildDir = "/zfspool/work-area/kernel-build"
var linuxDir = buildDir + "/linux"

// CompileLinuxKernel Compile the linux kernel we need to enable the nvme-fa options
// Once compiled and installed, we will need to remove any header packages that might be left behind and
// reboot the system
func CompileLinuxKernel(c *echo.Context) error {

	log.Printf("SSE client connected, ip: %v", c.RealIP())

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	//Create a work area
	if err := os.Mkdir(buildDir, os.ModePerm); err != nil {
		// os.ModePerm is equivalent to 0777 on most systems, then modified by umask
		log.Println(err)
		//return c.JSON(http.StatusInternalServerError, "")
	}

	//Install the dev tools needed to compile the kernel
	aptErr := installDevTools(w)
	if aptErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error installing dev tools")
	}

	//Clone the kernel repo
	cloneErr := cloneLinuxKernel(w)
	if cloneErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error cloning linux kernel repo")
	}

	//Prep the kernel build
	prepareErr := prepareKernelBuild(w)
	if prepareErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error preparing kernel build")
	}

	//Configure the kernel build config.
	configErr := updateBuildConfig(w)
	if configErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error configuring kernel build")
	}

	//Build the kernel
	buildErr := buildKernel(w)
	if buildErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error building kernel")
	}

	//Install the kernel
	installErr := installKernel(w)
	if installErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error installing kernel")
	}

	//Copy the files to the boot directory
	copyErr := copyFiles(w)
	if copyErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error copying files")
	}

	//Build deb packages
	createDebPackageErr := createDebPackage(w)
	if createDebPackageErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error creating deb package")
	}

	//Copy the header file to the home directory
	copyHeaderErr := copyHeaderFile(w)
	if copyHeaderErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error copying header file to home directory")
	}

	return c.JSON(http.StatusOK, "Built and installed the linux kernel successfully")
}

// InstallDevTools Install the dev tools needed to compile the kernel
// sudo apt install bc bison flex libssl-dev make git libncurses-dev
func installDevTools(w http.ResponseWriter) error {

	log.Printf("Installing DEV Tools")
	if err := utilities.SendEventData("Installing DEV Tools", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	cmd := exec.Command("apt", "install", "bc", "bison", "flex", "libssl-dev", "make", "git", "libncurses-dev", "debhelper-compat", "libdw-dev", "libelf-dev", "-y")

	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while running apt update", "Value", err)
		return err
	}

	return nil

}

// CloneLinuxKernel Clone the linux kernel repo
func cloneLinuxKernel(w http.ResponseWriter) error {

	log.Printf("Cloning Linux Kernel")
	if err := utilities.SendEventData("Cloning Linux Kernel", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	//Need to check if the directory exists, if so, delete it
	if _, err := os.Stat(linuxDir); err == nil {
		err := os.RemoveAll(linuxDir)
		if err != nil {
			slog.Error("Error while removing linux directory", "Value", err)
		}
	}

	cmd := exec.Command("git", "clone", "--depth=1", "--branch", "rpi-6.18.y", "https://github.com/raspberrypi/linux")
	cmd.Dir = buildDir
	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while running apt update", "Value", err)
		return err
	}

	return nil
}

// PrepareKernelBuild Configure the kernel build config.
func prepareKernelBuild(w http.ResponseWriter) error {

	log.Printf("Preparing the Linux Kernel")
	if err := utilities.SendEventData("Preparing the Linux Kernel", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	cmd := exec.Command("make", "bcm2712_defconfig")
	cmd.Dir = linuxDir
	cmd.Env = append(os.Environ(), "KERNEL=kernel_2712")

	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while running preparing the linux kernel", "Value", err)
		return err
	}

	return nil

}

// UpdateBuildConfig Enable nvme-fa options
func updateBuildConfig(w http.ResponseWriter) error {

	log.Printf("Update Build Config")
	if err := utilities.SendEventData("Update Build Config", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	var buf bytes.Buffer
	buf.WriteString("CONFIG_TLS=m\n")
	buf.WriteString("CONFIG_STREAM_PARSER=y\n")
	buf.WriteString("CONFIG_NVME_KEYRING=m\n")
	buf.WriteString("CONFIG_NVME_CORE=y\n")
	buf.WriteString("CONFIG_BLK_DEV_NVME=y\n")
	buf.WriteString("CONFIG_NVME_HWMON=y\n")
	buf.WriteString("CONFIG_NVME_TARGET=y\n")
	buf.WriteString("CONFIG_NVME_TARGET_TCP=y\n")
	buf.WriteString("CONFIG_NVME_TARGET_TCP_TLS=y\n")
	buf.WriteString("CONFIG_NVME_FABRICS=m\n")
	buf.WriteString("CONFIG_NVME_TCP=m\n")
	buf.WriteString("CONFIG_NVME_TCP_TLS=y\n")

	//Create a config fragment file
	fileErr := os.WriteFile(linuxDir+"/kernel/configs/nvme_fa.config", buf.Bytes(), 0644)
	if fileErr != nil {
		return fileErr
	}

	//Update the kernel build config with the nvme-fa options
	cmd := exec.Command("make", "nvme_fa.config")
	cmd.Dir = linuxDir
	cmd.Env = append(os.Environ(), "KERNEL=kernel_2712")

	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while updating the build config", "Value", err)
		return err
	}

	return nil
}

// buildKernel Build the kernel
func buildKernel(w http.ResponseWriter) error {

	log.Printf("Build the Kernel")
	if err := utilities.SendEventData("Build the Kernel", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	cmd := exec.Command("make", "-j4", "Image.gz", "modules", "dtbs")
	cmd.Dir = linuxDir
	cmd.Env = append(os.Environ(), "KERNEL=kernel_2712")
	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while building the kernel", err)
		return err
	}

	return nil

}

// installKernel Install the kernel
func installKernel(w http.ResponseWriter) error {

	// Install Kernel modules
	log.Printf("Installing the Kernel")
	if err := utilities.SendEventData("Installing the Kernel", "consoleOutput", w); err != nil {
		slog.Error("Error while sending event data", "Value", err)
	}

	cmd := exec.Command("make", "-j4", "modules_install")
	cmd.Dir = linuxDir
	cmd.Env = append(os.Environ(), "KERNEL=kernel_2712")

	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while installing the kernel", "Value", err)
		return err
	}

	return nil
}

func copyFiles(w http.ResponseWriter) error {

	//Copy files
	backupImageErr := copyFile("/boot/firmware/kernel_2712.img", "/boot/firmware/kernel_2712-backup.img", w)
	if backupImageErr != nil {
		return backupImageErr
	}

	copyImageErr := copyFile(linuxDir+"/arch/arm64/boot/Image.gz", "/boot/firmware/kernel_2712.img", w)
	if copyImageErr != nil {
		return copyImageErr
	}

	copyDTBErr := copyFile(linuxDir+"/arch/arm64/boot/dts/broadcom/*.dtb", "/boot/firmware/.", w)
	if copyDTBErr != nil {
		return copyDTBErr
	}

	copyOverlaysErr := copyFile(linuxDir+"/arch/arm64/boot/dts/overlays/*.dtb*", "/boot/firmware/overlays/.", w)
	if copyOverlaysErr != nil {
		return copyOverlaysErr
	}

	copyOverlaysReadMeErr := copyFile(linuxDir+"/arch/arm64/boot/dts/overlays/README", "/boot/firmware/overlays/.", w)
	if copyOverlaysReadMeErr != nil {
		return copyOverlaysReadMeErr
	}

	return nil
}

// createDebPackage Create deb packages
func createDebPackage(w http.ResponseWriter) error {
	// Create deb packages
	log.Printf("Creating deb packages")
	if err := utilities.SendEventData("Creating deb packages", "consoleOutput", w); err != nil {
		slog.Error("Error while creating deb packages", "Value", err)
	}

	//make deb-pkg
	cmd := exec.Command("make", "deb-pkg")
	cmd.Dir = linuxDir
	cmd.Env = append(os.Environ(), "KERNEL=kernel_2712")

	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while creating deb packages", "Value", err)
		return err
	}

	log.Printf("Created deb packages successfully")
	if err := utilities.SendEventData("Created deb packages successfully", "consoleOutput", w); err != nil {
		slog.Error("Error while creating deb packages", "Value", err)
	}

	return nil
}

// copyHeaderFile Copy the header package file to the home directory
func copyHeaderFile(w http.ResponseWriter) error {

	// Create deb packages
	log.Printf("Coping header deb package to home directory")
	if err := utilities.SendEventData("Coping header deb package to home directory", "consoleOutput", w); err != nil {
		slog.Error("Error while coping header deb package to home directory", "Value", err)
	}

	copyOverlaysReadMeErr := copyFile(buildDir+"/linux-headers*.deb", "~/.", w)
	if copyOverlaysReadMeErr != nil {
		return copyOverlaysReadMeErr
	}

	return nil
}

// CopyFile Command line copy function
func copyFile(src string, dst string, w http.ResponseWriter) error {

	log.Println("Copying file ", src, " to ", dst)
	sendEventErr := utilities.SendEventData("Copying file "+src+" to "+dst, "consoleOutput", w)
	if sendEventErr != nil {
		return sendEventErr
	}

	cmd := exec.Command("/bin/sh", "-c", "cp "+src+" "+dst)
	cmd.Dir = "/"

	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while running apt update", "Value", err)
		return err
	}

	return nil

}
