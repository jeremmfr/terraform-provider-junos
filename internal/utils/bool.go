package utils

import (
	"strconv"
	"strings"
)

func ParseTrue(str string) bool {
	switch v, err := strconv.ParseBool(strings.ToLower(str)); {
	case err != nil:
		return false
	default:
		return v
	}
}
