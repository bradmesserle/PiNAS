package endpoints

import (
	"github.com/labstack/echo/v5"
	"github.com/pinas/ui/internal/components"
	"github.com/pinas/ui/internal/structs"
)

func Home(c *echo.Context, wizardInfo *structs.WizardInfo, status *structs.SetupInfo) error {

	if status.IsRunning {

	}

	var cmp = components.Home(*wizardInfo, *status)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return cmp.Render(c.Request().Context(), c.Response())
}
