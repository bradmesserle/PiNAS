package system_info

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"

	"github.com/labstack/echo/v4"
)

type PIModel struct {
	PiModel string `json:"piModel"`
}

// GetPiModelInfo Retrieves the Raspberry Pi model information
func GetPiModelInfo(c echo.Context) error {

	out, err := exec.Command("cat", "/proc/cpuinfo").Output()

	if err != nil {
		log.Println(err)
		return err
	}

	modelInfo := PIModel{}
	modelInfo.PiModel = GetFieldValue(string(out), "Model")

	jsonData, marshalErr := json.Marshal(modelInfo)

	if marshalErr != nil {
		log.Println(marshalErr)
		return marshalErr
	}

	return c.JSON(http.StatusOK, string(jsonData))
}
