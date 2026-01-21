package project

import "github.com/microsoft/typescript-go/internal/tspath"

func isDynamicFileName(fileName string) bool {
	return tspath.IsDynamicFileName(fileName)
}
