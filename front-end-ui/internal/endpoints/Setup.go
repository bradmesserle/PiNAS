package endpoints

import (
	"fmt"
	"log"
	"log/slog"
	"sync"

	"github.com/labstack/echo/v5"
	"github.com/pinas/ui/internal/components/setup"
	"github.com/pinas/ui/internal/setup-process"
	"github.com/pinas/ui/internal/structs"
	"github.com/starfederation/datastar-go/datastar"
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
			startSSEMonitor(c)
			status.IsRunning = false
		})
	}

	var cmp = setup.InstallProgressPage()
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return cmp.Render(c.Request().Context(), c.Response())

}

// Start the SSE monitor. This to detect reboots and reconnect when the server comes back up
func startSSEMonitor(c *echo.Context) {

	// NewSSE sets the SSE headers and returns a generator bound to this request.
	sse := datastar.NewSSE(c.Response(), c.Request())
	err := sse.ExecuteScript(fmt.Sprintf(`startServerMonitor()`))
	if err != nil {
		log.Println(err)
	}

}
