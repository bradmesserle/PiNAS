package endpoints

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/pinas/ui/internal/components"
	"github.com/pinas/ui/internal/components/setup"
	"github.com/pinas/ui/internal/structs"
)

func WizardNavigation(c echo.Context, wizardInfo *structs.WizardInfo) error {

	if wizardInfo.Step == 0 {
		var cmp templ.Component = components.Home(*wizardInfo)
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
		return cmp.Render(c.Request().Context(), c.Response().Writer)
	}

	if wizardInfo.Step == 1 {

		//Get System Info

		//Render System Info Page
		return SystemInfo(c, wizardInfo)
	}

	if wizardInfo.Step == 2 {

		//Render NAS Options Page
		return NasOptions(c, wizardInfo)
	}

	var cmp templ.Component = setup.SetupWizard(*wizardInfo)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return cmp.Render(c.Request().Context(), c.Response().Writer)

}
