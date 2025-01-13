package utils

import "strings"

func QueryParamExists(rawQuery, param string) bool {
	return strings.Contains(rawQuery, param+"=")
}
