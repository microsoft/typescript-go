package core

import "github.com/microsoft/typescript-go/internal/tspath"

func GetOutputDeclarationFileName(
	inputFileName string,
	configFile *ParsedOptions,
	getCommonSourceDirectory func() string,
	options tspath.ComparePathsOptions,
) string {
	return tspath.ChangeExtension(
		getOutputPathWithoutChangingExt(
			inputFileName,
			IfElse(configFile.CompilerOptions.DeclarationDir != "", configFile.CompilerOptions.DeclarationDir, configFile.CompilerOptions.OutDir),
			getCommonSourceDirectory,
			options,
		),
		tspath.GetDeclarationEmitExtensionForPath(inputFileName),
	)
}

func getOutputPathWithoutChangingExt(
	inputFileName string,
	outputDir string,
	getCommonSourceDirectory func() string,
	options tspath.ComparePathsOptions,
) string {
	if outputDir != "" {
		return tspath.CombinePaths(
			outputDir,
			tspath.GetRelativePathFromDirectory(getCommonSourceDirectory(), inputFileName, options),
		)
	}
	return inputFileName
}
