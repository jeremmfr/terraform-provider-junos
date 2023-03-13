package utils

import "strconv"

func ConvI64toa(i int64) string {
	return strconv.FormatInt(i, 10)
}

//nolint:wrapcheck
func ConvAtoi64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
