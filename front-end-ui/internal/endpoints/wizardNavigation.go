package endpoints

import (
	"github.com/labstack/echo/v4"
	"github.com/pinas/ui/internal/components"
	"github.com/pinas/ui/internal/components/setup"
	"github.com/pinas/ui/internal/rest-client"
	"github.com/pinas/ui/internal/structs"
)

func WizardNavigation(c echo.Context, wizardInfo *structs.WizardInfo) error {

	if wizardInfo.Step == 0 {
		cmp := components.Home(*wizardInfo)
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
		return cmp.Render(c.Request().Context(), c.Response().Writer)
	}

	if wizardInfo.Step == 1 {

		systemInfoData := new(structs.SystemInfo)

		//Get System Info
		var modelInfo, err = rest_client.GetModelInfo()
		systemInfoData.Model = modelInfo

		if err != nil {
			return err
		}

		//Render System Info Page
		return SystemInfo(c, wizardInfo, systemInfoData)
	}

	if wizardInfo.Step == 2 {

		//Render NAS Options Page
		return NasOptions(c, wizardInfo)
	}

	cmp := setup.SetupWizard(*wizardInfo)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return cmp.Render(c.Request().Context(), c.Response().Writer)

}
