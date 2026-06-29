package endpoints

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pinas/common-structs"
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

		// Get System Info
		systemInfoData.Model = getModelInfo()
		// Get CPU Info
		systemInfoData.CpuInfo = getCpuInfo()
		// Get Pcie Info
		systemInfoData.PcieInfo = getPcieInfo()
		// Get Memory Info
		systemInfoData.MemoryInfo = getMemoryInfo()

		//Render System Info Page
		return SystemInfo(c, wizardInfo, systemInfoData)
	}

	if wizardInfo.Step == 2 {
		//Render Drive Setup Page
		return DriveSetup(c, wizardInfo)
	}

	if wizardInfo.Step == 3 {
		//Render NAS Options Page
		return NasOptions(c, wizardInfo)
	}

	if wizardInfo.Step == 4 {
		//Render Install Summary Page
		return InstallSummary(c, wizardInfo)
	}

	cmp := setup.SetupWizard(*wizardInfo)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return cmp.Render(c.Request().Context(), c.Response().Writer)

}

// getModelInfo Get Model Info
func getModelInfo() common_structs.PIModel {
	//Get System Info
	var modelInfo, err = rest_client.GetModelInfo()
	if err != nil {
		log.Error(err.Error())
	}

	return modelInfo
}

// getCpuInfo Get CPU Info
func getCpuInfo() common_structs.CpuInfo {
	//Get CPU Info
	var cpuInfo, err = rest_client.GetCpuInfo()
	if err != nil {
		log.Error(err.Error())
	}

	return cpuInfo
}

// getPcieInfo Get PCIe Info
func getPcieInfo() common_structs.PCIeInfo {
	//Get Pcie Info
	var pcieInfo, err = rest_client.GetPCIeInfo()
	if err != nil {
		log.Error(err.Error())
	}

	return pcieInfo
}

// Get Memory Info
func getMemoryInfo() common_structs.MemoryInfo {
	//Get Memory Info
	var memoryInfo, err = rest_client.GetMemoryInfo()
	if err != nil {
		log.Error(err.Error())
	}

	return memoryInfo
}
