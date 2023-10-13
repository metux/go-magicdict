package utils

import (
	"strconv"
	"strings"
)

func StrToBool(s string, dflt bool) bool {
	switch strings.ToLower(s) {
	case "1", "y", "yes", "true", "on":
		return true
	case "0", "n", "no", "false", "off":
		return false
	default:
		return dflt
	}
}

func StrToInt(s string, dflt int) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	} else {
		return dflt
	}
}
