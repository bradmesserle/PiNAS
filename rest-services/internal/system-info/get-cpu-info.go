package system_info

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetCpu(c echo.Context) error {

	fmt.Println("CPU Info")

	return c.JSON(http.StatusOK, "CPU Info")
}
