package utils

import (
	"strconv"
)

func ParseQueryInt(value string, defaultValue int) int {
	if parsedValue, err := strconv.Atoi(value); err == nil {
		return parsedValue
	}
	return defaultValue
}
