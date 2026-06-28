package system_info

import (
	"log"
	"net/http"
	"os/exec"

	"github.com/labstack/echo/v4"
	"github.com/pinas/common-structs"
)

// GetPiModelInfo Retrieves the Raspberry Pi model information
func GetPiModelInfo(c echo.Context) error {

	out, err := exec.Command("cat", "/proc/cpuinfo").Output()

	if err != nil {
		log.Println(err)
		return err
	}

	modelInfo := common_structs.PIModel{}
	modelInfo.PiModel = GetFieldValue(string(out), "Model")
	return c.JSON(http.StatusOK, modelInfo)
}
