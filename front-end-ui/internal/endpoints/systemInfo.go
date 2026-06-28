package endpoints

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/pinas/ui/internal/components/setup"
	"github.com/pinas/ui/internal/structs"
)

func SystemInfo(c echo.Context, wizardInfo *structs.WizardInfo, systemInfo *structs.SystemInfo) error {

	var cmp templ.Component = setup.SystemInfoPage(*wizardInfo, *systemInfo)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
