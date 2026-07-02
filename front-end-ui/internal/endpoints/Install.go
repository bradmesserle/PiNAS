package endpoints

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
	"github.com/pinas/ui/internal/components/setup"
	"github.com/pinas/ui/internal/structs"
)

func Install(c *echo.Context, wizardInfo *structs.WizardInfo) error {

	//Kick off the install process

	var cmp templ.Component = setup.InstallProgressPage(*wizardInfo)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return cmp.Render(c.Request().Context(), c.Response())
}
