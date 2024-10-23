package repo

import (
	"fmt"
	"os"
	"path/filepath"
)

var RepositoryRootPath string
var TypeScriptSubmodulePath string
var TestDataPath string

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("failed to get current working directory: %v", err))
	}
	RepositoryRootPath = findGoMod(cwd)
	TypeScriptSubmodulePath = filepath.Join(RepositoryRootPath, "_submodules", "TypeScript")
	TestDataPath = filepath.Join(RepositoryRootPath, "testdata")
}

func findGoMod(dir string) string {
	root := filepath.VolumeName(dir)
	for dir != root {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	panic("could not find go.mod")
}
