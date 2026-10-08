package endpoints

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
	"github.com/pinas/ui/internal/components/setup"
	"github.com/pinas/ui/internal/structs"
	"github.com/starfederation/datastar-go/datastar"
)

func InstallSummary(c *echo.Context, wizardInfo *structs.WizardInfo) error {
	var cmp templ.Component = setup.InstallSummaryPage(*wizardInfo)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	err := cmp.Render(c.Request().Context(), c.Response())
	//startSSEMonitor(c)
	return err
}

func startSSEMonitor(c *echo.Context) {

	slog.Info("Starting the server monitor")
	// NewSSE sets the SSE headers and returns a generator bound to this request.
	sse := datastar.NewSSE(c.Response(), c.Request())
	err := sse.ExecuteScript(fmt.Sprintf(`startServerMonitor()`))
	if err != nil {
		log.Println(err)
	}

}
