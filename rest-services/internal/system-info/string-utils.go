package system_info

import (
	"strings"
)

// GetLastFieldValue Get the value of a field from the data
func GetLastFieldValue(data string, field string) string {
	startIndex := strings.LastIndex(data, field)
	endIndex := strings.Index(data[startIndex:], "\n")

	if len(strings.Split(data[startIndex:endIndex+startIndex], ":")) > 0 {
		return strings.Split(data[startIndex:endIndex+startIndex], ":")[1]
	}

	return ""
}

// GetFieldValue Get the value of a field from the data
func GetFieldValue(data string, field string) string {
	startIndex := strings.Index(data, field)
	endIndex := strings.Index(data[startIndex:], "\n")

	if len(strings.Split(data[startIndex:endIndex+startIndex], ":")) > 0 {
		return strings.Split(data[startIndex:endIndex+startIndex], ":")[1]
	}

	return ""
}
