package module

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
)

const InferredTypesContainingFile = "__inferred type names__.ts"

func ParseNodeModuleFromPath(resolved string, isFolder bool) string {
	path := tspath.NormalizePath(resolved)
	idx := strings.LastIndex(path, "/node_modules/")
	if idx == -1 {
		return ""
	}

	indexAfterNodeModules := idx + len("/node_modules/")
	indexAfterPackageName := moveToNextDirectorySeparatorIfAvailable(path, indexAfterNodeModules, isFolder)
	if path[indexAfterNodeModules] == '@' {
		indexAfterPackageName = moveToNextDirectorySeparatorIfAvailable(path, indexAfterPackageName, isFolder)
	}
	return path[:indexAfterPackageName]
}

func ParsePackageName(moduleName string) (packageName, rest string) {
	idx := strings.Index(moduleName, "/")
	if len(moduleName) > 0 && moduleName[0] == '@' {
		idx = strings.Index(moduleName[idx+1:], "/") + idx + 1
	}
	if idx == -1 {
		return moduleName, ""
	}
	return moduleName[:idx], moduleName[idx+1:]
}

func GetEffectiveTypeRoots(options *core.CompilerOptions, currentDirectory string) (result []string, fromConfig bool) {
	if options.TypeRoots != nil {
		return options.TypeRoots, true
	}
	var baseDir string
	if options.ConfigFilePath != "" {
		baseDir = tspath.GetDirectoryPath(options.ConfigFilePath)
	} else {
		baseDir = currentDirectory
		if baseDir == "" {
			// This was accounted for in the TS codebase, but only for third-party API usage
			// where the module resolution host does not provide a getCurrentDirectory().
			panic("cannot get effective type roots without a config file path or current directory")
		}
	}

	typeRoots := make([]string, 0, strings.Count(baseDir, "/"))
	tspath.ForEachAncestorDirectory(baseDir, func(dir string) (any, bool) {
		typeRoots = append(typeRoots, tspath.CombinePaths(dir, "node_modules", "@types"))
		return nil, false
	})
	return typeRoots, false
}

func MangleScopedPackageName(packageName string) string {
	if packageName[0] == '@' {
		idx := strings.Index(packageName, "/")
		if idx == -1 {
			return packageName
		}
		return packageName[1:idx] + "__" + packageName[idx+1:]
	}
	return packageName
}

func UnmangleScopedPackageName(packageName string) string {
	idx := strings.Index(packageName, "__")
	if idx != -1 {
		return "@" + packageName[:idx] + "/" + packageName[idx+2:]
	}
	return packageName
}
