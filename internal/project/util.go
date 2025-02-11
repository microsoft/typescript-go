package project

import "strings"

func isDynamicFileName(fileName string) bool {
	return strings.HasPrefix(fileName, "^")
}

func isBundledUri(fileName string) bool {
	return strings.HasPrefix(fileName, "bundled://")
}
