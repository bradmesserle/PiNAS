package utilities

import (
	"log"
	"os"
	"strings"
)

// Hardcoded file path for now.
var Path = "/boot/firmware/config.txt"

func IsPCIeEx1SetToGen3() (bool, error) {

	content, err := os.ReadFile(Path)
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	//See if we have property pciex1
	searchString := "dtparam=pciex1_gen=3"
	if strings.Contains(string(content), searchString) {
		return true, err
	}

	return false, err

}
