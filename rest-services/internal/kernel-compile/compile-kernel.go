package kernel_compile

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/labstack/echo/v4"
)

var buildDir = "/kernel-build"
var linuxDir = buildDir + "/linux"

// CompileLinuxKernel Compile the linux kernel we need to enable the nvme-fa options
// Once compiled and installed, we will need to remove any header packages that might be left behind and
// reboot the system
func CompileLinuxKernel(c echo.Context) error {

	//Create a work area
	if err := os.Mkdir(buildDir, os.ModePerm); err != nil {
		// os.ModePerm is equivalent to 0777 on most systems, then modified by umask
		log.Println(err)
		//return c.JSON(http.StatusInternalServerError, "")
	}
	//Install the dev tools needed to compile the kernel
	aptErr := InstallDevTools()
	if aptErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error installing dev tools")
	}

	//Clone the kernel repo
	cloneErr := CloneLinuxKernel()
	if cloneErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error cloning linux kernel repo")
	}

	//Prep the kernel build
	prepareErr := PrepareKernelBuild()
	if prepareErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error preparing kernel build")
	}

	//Configure the kernel build config.
	configErr := UpdateBuildConfig()
	if configErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error configuring kernel build")
	}

	//Build the kernel
	buildErr := BuildKernel()
	if buildErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error building kernel")
	}

	//Install the kernel
	installErr := InstallKernel()
	if installErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error installing kernel")
	}

	return c.JSON(http.StatusOK, "Built and installed the linux kernel successfully")
}

// InstallDevTools Install the dev tools needed to compile the kernel
// sudo apt install bc bison flex libssl-dev make git libncurses-dev
func InstallDevTools() error {
	aptCmd := exec.Command("apt", "install", "bc", "bison", "flex", "libssl-dev", "make", "git", "libncurses-dev", "-y")
	aptCmd.Stdout = os.Stdout
	aptCmd.Stderr = os.Stderr
	aptErr := aptCmd.Run()
	if aptErr != nil {
		log.Println(aptErr)
		return aptErr
	}

	return nil
}

// CloneLinuxKernel Clone the linux kernel repo
func CloneLinuxKernel() error {

	//Need to check if the directory exists, if so, delete it
	if _, err := os.Stat(linuxDir); err == nil {
		err := os.RemoveAll(linuxDir)
		if err != nil {
			log.Println(err)
		}
	}

	cloneCmd := exec.Command("git", "clone", "--depth=1", "https://github.com/raspberrypi/linux")
	cloneCmd.Dir = buildDir
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	cloneErr := cloneCmd.Run()
	if cloneErr != nil {
		log.Println(cloneErr)
		return cloneErr
	}
	return nil
}

// PrepareKernelBuild Configure the kernel build config.
func PrepareKernelBuild() error {
	cmd := exec.Command("make", "bcm2712_defconfig")
	cmd.Dir = linuxDir
	cmd.Env = append(os.Environ(), "KERNEL=kernel_2712")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// UpdateBuildConfig Enable nvme-fa options
func UpdateBuildConfig() error {

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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	runErr := cmd.Run()
	if runErr != nil {
		return runErr
	}

	return nil
}

// BuildKernel Build the kernel
func BuildKernel() error {
	cmd := exec.Command("make", "-j6", "Image.gz", "modules", "dtbs")
	cmd.Dir = linuxDir
	cmd.Env = append(os.Environ(), "KERNEL=kernel_2712")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// InstallKernel Install the kernel
func InstallKernel() error {

	return nil
}
