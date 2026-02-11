package kernel_compile

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// CompileLinuxKernel Compile the linux kernel we need to enable the nvme-fa options
// Once compiled and installed, we will need to remove any header packages that might be left behind and
// reboot the system
func CompileLinuxKernel(c echo.Context) error {

	//Create a work area

	//Install the dev tools needed to compile the kernel

	//Clone the kernel repo

	//Configure the kernel build config.

	//Prep the kernel build

	//Build the kernel

	//Install the kernel

	return c.JSON(http.StatusOK, "Built and installed the linux kernel successfully")
}
