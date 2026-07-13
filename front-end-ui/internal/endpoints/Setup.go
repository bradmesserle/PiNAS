package endpoints

import (
	"sync"

	"github.com/labstack/echo/v5"
	"github.com/pinas/ui/internal/components/setup"
	setup_process "github.com/pinas/ui/internal/setup-process"
	"github.com/pinas/ui/internal/structs"
)

func Setup(c *echo.Context, wizardInfo *structs.WizardInfo, status *structs.SetupInfo) error {

	//Check to see if we are already running
	if !status.IsRunning {
		//Kick off the installation process
		var wg sync.WaitGroup

		//Kick off a background process to start the setup process.
		wg.Go(func() {
			status.IsRunning = true
			setup_process.PerformSetup(wizardInfo)
			status.IsRunning = false
		})
	}

	var cmp = setup.InstallProgressPage()
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return cmp.Render(c.Request().Context(), c.Response())

}
