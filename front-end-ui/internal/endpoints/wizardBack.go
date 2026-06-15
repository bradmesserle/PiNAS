package endpoints

import (
	"github.com/labstack/echo/v4"
	"github.com/pinas/ui/internal/structs"
)

func WizardBack(c echo.Context, wizardInfo *structs.WizardInfo) error {

	wizardInfo.Step = wizardInfo.Step - 1
	return WizardNavigation(c, wizardInfo)

	//if wizardInfo.Step == 0 {
	//	var cmp templ.Component = components.Home(*wizardInfo)
	//	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	//	return cmp.Render(c.Request().Context(), c.Response().Writer)
	//}
	//
	//var cmp templ.Component = setup.SetupWizard(*wizardInfo)
	//c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	//return cmp.Render(c.Request().Context(), c.Response().Writer)

}
