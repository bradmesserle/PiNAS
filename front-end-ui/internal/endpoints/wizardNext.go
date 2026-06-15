package endpoints

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/pinas/ui/internal/components/setup"
	"github.com/pinas/ui/internal/structs"
)

func WizardNext(c echo.Context, wizardInfo *structs.WizardInfo) error {

	wizardInfo.Step = wizardInfo.Step + 1

	var cmp templ.Component = setup.SetupWizard(*wizardInfo)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
