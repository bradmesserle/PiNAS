package endpoints

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/pinas/ui/internal/components/setup"
)

var WizardData = setup.WizardData{Step: 0}

func Setup(c echo.Context) error {

	var cmp templ.Component = setup.Setup(WizardData)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
