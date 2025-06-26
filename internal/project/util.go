package project

import "strings"

func IsDynamicFileName(fileName string) bool {
	return strings.HasPrefix(fileName, "^")
}
