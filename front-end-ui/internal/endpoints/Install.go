package endpoints

import (
	"sync"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
	"github.com/pinas/ui/internal/components/setup"
	"github.com/pinas/ui/internal/setup-process"
	"github.com/pinas/ui/internal/structs"
)

func Install(c *echo.Context, wizardInfo *structs.WizardInfo) error {

	//Kick off the installation process
	var wg sync.WaitGroup

	var cmp templ.Component = setup.InstallProgressPage(*wizardInfo)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	renderPage := cmp.Render(c.Request().Context(), c.Response())

	//Kick off a background process to start the setup process.
	wg.Go(func() {
		setup_process.PerformSetup(wizardInfo)
	})

	return renderPage
}
