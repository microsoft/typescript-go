package outputpaths

import (
	"fmt"
	"testing"
)

func TestZZReproPanic(t *testing.T) {
	currentDirectory := "/a"
	commonSourceDirectory := "../../../../../../../../../../a/src"
	fileName := "/a/src/x.ts"
	newDirPath := "/a/out"

	fmt.Println("commonSourceDirectory (raw)+sep len:", len(commonSourceDirectory)+1)
	fmt.Println("fileName len:", len(fileName))

	result := GetSourceFilePathInNewDir(fileName, newDirPath, currentDirectory, commonSourceDirectory, true)
	fmt.Println("result:", result)
}
