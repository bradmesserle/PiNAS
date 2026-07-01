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

var buildDir = "/kernel-build"
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
	aptErr := InstallDevTools(w)
	if aptErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error installing dev tools")
	}

	//Clone the kernel repo
	cloneErr := CloneLinuxKernel(w)
	if cloneErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error cloning linux kernel repo")
	}

	//Prep the kernel build
	prepareErr := PrepareKernelBuild(w)
	if prepareErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error preparing kernel build")
	}

	//Configure the kernel build config.
	configErr := UpdateBuildConfig(w)
	if configErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error configuring kernel build")
	}

	//Build the kernel
	buildErr := BuildKernel(w)
	if buildErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error building kernel")
	}

	//Install the kernel
	installErr := InstallKernel(w)
	if installErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error installing kernel")
	}

	return c.JSON(http.StatusOK, "Built and installed the linux kernel successfully")
}

// InstallDevTools Install the dev tools needed to compile the kernel
// sudo apt install bc bison flex libssl-dev make git libncurses-dev
func InstallDevTools(w http.ResponseWriter) error {
	cmd := exec.Command("apt", "install", "bc", "bison", "flex", "libssl-dev", "make", "git", "libncurses-dev", "-y")

	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while running apt update", err)
		return err
	}

	return nil

}

// CloneLinuxKernel Clone the linux kernel repo
func CloneLinuxKernel(w http.ResponseWriter) error {

	//Need to check if the directory exists, if so, delete it
	if _, err := os.Stat(linuxDir); err == nil {
		err := os.RemoveAll(linuxDir)
		if err != nil {
			log.Println(err)
		}
	}

	cmd := exec.Command("git", "clone", "--depth=1", "https://github.com/raspberrypi/linux")
	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while running apt update", err)
		return err
	}

	return nil
}

// PrepareKernelBuild Configure the kernel build config.
func PrepareKernelBuild(w http.ResponseWriter) error {
	cmd := exec.Command("make", "bcm2712_defconfig")
	cmd.Dir = linuxDir
	cmd.Env = append(os.Environ(), "KERNEL=kernel_2712")

	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while running apt update", err)
		return err
	}

	return nil

}

// UpdateBuildConfig Enable nvme-fa options
func UpdateBuildConfig(w http.ResponseWriter) error {

	var buf bytes.Buffer
	buf.WriteString("CONFIG_TLS=m\n")
	buf.WriteString("CONFIG_STREAM_PARSER=y\n")
	buf.WriteString("CONFIG_NVME_KEYRING=m\n")
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
		slog.Error("Error while running apt update", err)
		return err
	}

	return nil
}

// BuildKernel Build the kernel
func BuildKernel(w http.ResponseWriter) error {
	cmd := exec.Command("make", "-j6", "Image.gz", "modules", "dtbs")
	cmd.Dir = linuxDir
	cmd.Env = append(os.Environ(), "KERNEL=kernel_2712")
	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while running apt update", err)
		return err
	}

	return nil

}

// InstallKernel Install the kernel
func InstallKernel(w http.ResponseWriter) error {

	// Install Kernel modules
	cmd := exec.Command("make", "-j6", "modules_install")
	cmd.Dir = linuxDir
	cmd.Env = append(os.Environ(), "KERNEL=kernel_2712")

	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while running apt update", err)
		return err
	}

	//Copy files
	backupImageErr := CopyFile("/boot/firmware/kernel_2712.img", "/boot/firmware/kernel_2712-backup.img", w)
	if backupImageErr != nil {
		return err
	}

	copyImageErr := CopyFile(linuxDir+"/arch/arm64/boot/Image.gz", "/boot/firmware/kernel_2712.img", w)
	if copyImageErr != nil {
		return err
	}

	copyDTBErr := CopyFile(linuxDir+"/arch/arm64/boot/dts/broadcom/*.dtb", "/boot/firmware/", w)
	if copyDTBErr != nil {
		return err
	}

	copyOverlaysErr := CopyFile(linuxDir+"/arch/arm64/boot/dts/overlays/*.dtb*", "/boot/firmware/overlays/", w)
	if copyOverlaysErr != nil {
		return err
	}

	copyOverlaysReadMeErr := CopyFile(linuxDir+"/arch/arm64/boot/dts/overlays/README", "/boot/firmware/overlays/", w)
	if copyOverlaysReadMeErr != nil {
		return err
	}

	return nil
}

// CopyFile Command line copy function
func CopyFile(src string, dst string, w http.ResponseWriter) error {

	log.Println("Copying file ", src, " to ", dst)
	cmd := exec.Command("cp", src, dst)

	err := utilities.ExecCmdSseStdoutText(cmd, w)
	if err != nil {
		slog.Error("Error while running apt update", err)
		return err
	}

	return nil

}
