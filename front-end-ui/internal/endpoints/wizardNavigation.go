package endpoints

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	common_structs "github.com/pinas/common-structs"
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
		systemInfoData.Model = GetModelInfo()
		// Get CPU Info
		systemInfoData.CpuInfo = GetCpuInfo()

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

// GetModelInfo Get Model Info
func GetModelInfo() common_structs.PIModel {
	//Get System Info
	var modelInfo, err = rest_client.GetModelInfo()
	if err != nil {
		log.Error(err.Error())
	}

	return modelInfo
}

// GetCpuInfo Get CPU Info
func GetCpuInfo() common_structs.CpuInfo {
	//Get CPU Info
	var cpuInfo, err = rest_client.GetCpuInfo()
	if err != nil {
		log.Error(err.Error())
	}

	return cpuInfo
}
