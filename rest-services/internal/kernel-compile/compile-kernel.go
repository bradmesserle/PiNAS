package kernel_compile

import (
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/labstack/echo/v4"
)

// CompileLinuxKernel Compile the linux kernel we need to enable the nvme-fa options
// Once compiled and installed, we will need to remove any header packages that might be left behind and
// reboot the system
func CompileLinuxKernel(c echo.Context) error {

	//Create a work area
	if err := os.Mkdir("/kernel-build", os.ModePerm); err != nil {
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

	//Configure the kernel build config.

	//Prep the kernel build
	prepareErr := PrepareKernelBuild()
	if prepareErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error preparing kernel build")
	}

	//Build the kernel
	buildErr := BuildKernel()
	if buildErr != nil {
		return c.JSON(http.StatusInternalServerError, "Error building kernel")
	}

	//Install the kernel

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
	if _, err := os.Stat("/kernel-build/linux"); err == nil {
		err := os.RemoveAll("/kernel-build/linux")
		if err != nil {
			log.Println(err)
		}
	}

	cloneCmd := exec.Command("git", "clone", "--depth=1", "https://github.com/raspberrypi/linux")
	cloneCmd.Dir = "/kernel-build"
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
	cmd.Dir = "/kernel-build/linux"
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

// BuildKernel Build the kernel
func BuildKernel() error {
	cmd := exec.Command("make", "-j6", "Image.gz", "modules", "dtbs")
	cmd.Dir = "/kernel-build/linux"
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
