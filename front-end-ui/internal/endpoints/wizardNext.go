package endpoints

import (
	"github.com/labstack/echo/v4"
	"github.com/pinas/ui/internal/structs"
)

func WizardNext(c echo.Context, wizardInfo *structs.WizardInfo) error {

	wizardInfo.Step = wizardInfo.Step + 1
	return WizardNavigation(c, wizardInfo)

}
