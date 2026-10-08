package endpoints

import (
	"log/slog"
	"sync"

	"github.com/labstack/echo/v5"
	"github.com/pinas/ui/internal/components/setup"
	"github.com/pinas/ui/internal/setup-process"
	"github.com/pinas/ui/internal/structs"
)

func Setup(c *echo.Context, wizardInfo *structs.WizardInfo, status *structs.SetupInfo) error {

	slog.Info("Running Status : %s", status.IsRunning)

	//Check to see if we are already running
	if !status.IsRunning {
		//Kick off the installation process
		var wg sync.WaitGroup

		//Kick off a background process to start the setup process.
		wg.Go(func() {

			//Save the setup options
			err := structs.SaveSetupOptions(wizardInfo.NasOptions)
			if err != nil {
				slog.Error("Failed to save setup options: %v", err)
				return
			}

			status.IsRunning = true
			setup_process.PerformSetup(c, wizardInfo)
			status.IsRunning = false
		})
	}

	var cmp = setup.InstallProgressPage()
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return cmp.Render(c.Request().Context(), c.Response())

}
