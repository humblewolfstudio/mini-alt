package utils

import "strings"

func ClearInput(input string) string {
	return strings.TrimPrefix(input, "/")
}
