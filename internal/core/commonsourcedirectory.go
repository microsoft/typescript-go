package core

import "github.com/microsoft/typescript-go/internal/tspath"

func computeCommonSourceDirectoryOfFilenames(fileNames []string, currentDirectory string, useCaseSensitiveFileNames bool) string {
	var commonPathComponents []string
	for _, sourceFile := range fileNames {
		// Each file contributes into common source file path
		sourcePathComponents := tspath.GetNormalizedPathComponents(sourceFile, currentDirectory)

		// The base file name is not part of the common directory path
		sourcePathComponents = sourcePathComponents[:len(sourcePathComponents)-1]

		if commonPathComponents == nil {
			// first file
			commonPathComponents = sourcePathComponents
			continue
		}

		n := min(len(commonPathComponents), len(sourcePathComponents))
		for i := range n {
			if tspath.GetCanonicalFileName(commonPathComponents[i], useCaseSensitiveFileNames) != tspath.GetCanonicalFileName(sourcePathComponents[i], useCaseSensitiveFileNames) {
				if i == 0 {
					// Failed to find any common path component
					return ""
				}

				// New common path found that is 0 -> i-1
				commonPathComponents = commonPathComponents[:i]
				break
			}
		}

		// If the sourcePathComponents was shorter than the commonPathComponents, truncate to the sourcePathComponents
		if len(sourcePathComponents) < len(commonPathComponents) {
			commonPathComponents = commonPathComponents[:len(sourcePathComponents)]
		}
	}

	if len(commonPathComponents) == 0 {
		// Can happen when all input files are .d.ts files
		return currentDirectory
	}

	return tspath.GetPathFromPathComponents(commonPathComponents)
}

func GetCommonSourceDirectory(options *CompilerOptions, files func() []string, currentDirectory string, useCaseSensitiveFileNames bool) string {
	var commonSourceDirectory string
	if options.RootDir != "" {
		// If a rootDir is specified use it as the commonSourceDirectory
		commonSourceDirectory = tspath.GetNormalizedAbsolutePath(options.RootDir, currentDirectory)
		// !!! checkSourceFilesBelongToPath?.(options.rootDir);
	} else if options.Composite.IsTrue() && options.ConfigFilePath != "" {
		// Project compilations never infer their root from the input source paths
		commonSourceDirectory = tspath.GetDirectoryPath(options.ConfigFilePath)
		// !!! checkSourceFilesBelongToPath?.(commonSourceDirectory);
	} else {
		commonSourceDirectory = computeCommonSourceDirectoryOfFilenames(files(), currentDirectory, useCaseSensitiveFileNames)
	}

	if len(commonSourceDirectory) > 0 {
		// Make sure directory path ends with directory separator so this string can directly
		// used to replace with "" to get the relative path of the source file and the relative path doesn't
		// start with / making it rooted path
		commonSourceDirectory = tspath.EnsureTrailingDirectorySeparator(commonSourceDirectory)
	}

	return commonSourceDirectory
}

func GetCommonSourceDirectoryOfConfig(options *ParsedOptions, currentDirectory string, useCaseSensitiveFileNames bool) string {
	return GetCommonSourceDirectory(
		options.CompilerOptions,
		func() []string {
			return Filter(
				options.FileNames,
				func(file string) bool {
					return !(options.CompilerOptions.NoEmitForJsFiles.IsTrue() && tspath.HasJSFileExtension(file)) &&
						!tspath.IsDeclarationFileName(file)
				})
		},
		currentDirectory,
		useCaseSensitiveFileNames,
	)
}
