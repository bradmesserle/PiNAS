// Package system_updates Update Raspberry PI system configuration files
package system_updates

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

// VerifyUpdateConfig Verify and Update config.txt enable pcie gen3
func VerifyUpdateConfig(c echo.Context) error {

	// Hardcoded file path for now.
	filePath := "/boot/firmware/config.txt"

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
		return err
	}

	//See if we have property pciex1
	searchString := "dtparam=pciex1"
	if !strings.Contains(string(content), searchString) {
		log.Println("String not found in file")
		err := updateFile(filePath, searchString)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Error while updating config.txt file"+err.Error())
		}
	}

	//See if we have property pciex1_gen=3
	searchString = "dtparam=pciex1_gen=3"
	if !strings.Contains(string(content), searchString) {
		log.Println("String not found in file")
		err := updateFile(filePath, searchString)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Error while updating config.txt file"+err.Error())
		}
	}

	return c.JSON(http.StatusOK, "Updated config.txt file")
}

func updateFile(filePath string, stringData string) error {

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	if _, err := file.WriteString("\n" + stringData); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
