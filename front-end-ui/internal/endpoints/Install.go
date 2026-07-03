package endpoints

import (
	"log"
	"sync"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
	"github.com/pinas/ui/internal/components/setup"
	"github.com/pinas/ui/internal/rest-client"
	"github.com/pinas/ui/internal/structs"
)

func Install(c *echo.Context, wizardInfo *structs.WizardInfo) error {

	//Kick off the installation process
	var wg sync.WaitGroup
	wg.Add(3)

	var cmp templ.Component = setup.InstallProgressPage(*wizardInfo)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	renderPage := cmp.Render(c.Request().Context(), c.Response())

	go func() {
		errUpdate := rest_client.AptUpdate(&wg)
		if errUpdate != nil {
			log.Println(errUpdate.Error())
		}

		errUpgrade := rest_client.AptUpgrade(&wg)
		if errUpgrade != nil {
			log.Println(errUpgrade.Error())
		}

		errInstallZfs := rest_client.InstallZfs(&wg)
		if errInstallZfs != nil {
			log.Println(errInstallZfs.Error())
		}

	}()

	return renderPage
}
