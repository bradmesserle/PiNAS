package endpoints

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
)

func Home(c *echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return cmp.Render(c.Request().Context(), c.Response())
}
