package endpoints

import (
	"log/slog"
	"sync"

	"github.com/labstack/echo/v5"
	"github.com/pinas/ui/internal/components"
	"github.com/pinas/ui/internal/components/setup"
	"github.com/pinas/ui/internal/setup-process"
	"github.com/pinas/ui/internal/structs"
)

func Home(c *echo.Context, wizardInfo *structs.WizardInfo, status *structs.SetupInfo) error {

	//Get current checkpoint
	checkpoint, err := setup_process.GetStatus()
	if err != nil {
		slog.Info("Error while getting checkpoint setting to fresh install", "err", err)
		checkpoint = structs.FreshInstall
	}

	if checkpoint == structs.FreshInstall {
		status.IsRunning = false
	} else {
		status.IsRunning = true
	}

	// We want to resume where we left off before a reboot
	if status.IsRunning {

		//Kick off the installation process
		var wg sync.WaitGroup

		//Kick off a background process to start the setup process.
		wg.Go(func() {
			setup_process.PerformSetup(wizardInfo)
		})

		var cmp = setup.InstallProgressPage()
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
		return cmp.Render(c.Request().Context(), c.Response())

	}

	var cmp = components.Home(*wizardInfo, *status)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return cmp.Render(c.Request().Context(), c.Response())

}
