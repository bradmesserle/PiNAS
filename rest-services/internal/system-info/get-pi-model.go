package system_info

import (
	"encoding/json"
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

	jsonData, marshalErr := json.Marshal(modelInfo)

	if marshalErr != nil {
		log.Println(marshalErr)
		return marshalErr
	}

	return c.JSON(http.StatusOK, string(jsonData))
}
