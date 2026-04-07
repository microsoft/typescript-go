package ls

import (
	"context"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/change"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
	"github.com/microsoft/typescript-go/internal/scanner"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type pathUpdater func(path string) (string, bool)

type toImport struct {
	newFileName string
	updated     bool
}

func (l *LanguageService) GetEditsForFileRename(ctx context.Context, oldURI lsproto.DocumentUri, newURI lsproto.DocumentUri) map[lsproto.DocumentUri][]*lsproto.TextEdit {
	program := l.GetProgram()
	oldPath := oldURI.FileName()
	newPath := newURI.FileName()

	oldToNew := l.createPathUpdater(oldPath, newPath)
	newToOld := l.createPathUpdater(newPath, oldPath)

	changeTracker := change.NewTracker(ctx, program.Options(), l.FormatOptions(), l.converters)
	l.updateTsconfigFiles(program, changeTracker, oldToNew, oldPath, newPath)
	l.updateImportsForFileRename(program, changeTracker, oldToNew, newToOld)

	result := map[lsproto.DocumentUri][]*lsproto.TextEdit{}
	for fileName, edits := range changeTracker.GetChanges() {
		result[lsconv.FileNameToDocumentURI(fileName)] = edits
	}
	return result
}

func (l *LanguageService) createPathUpdater(oldPath string, newPath string) pathUpdater {
	compareOptions := tspath.ComparePathsOptions{UseCaseSensitiveFileNames: l.UseCaseSensitiveFileNames()}
	getUpdatedPath := func(path string) (string, bool) {
		if tspath.ComparePaths(path, oldPath, compareOptions) == 0 {
			return newPath, true
		}
		if tspath.StartsWithDirectory(path, oldPath, l.UseCaseSensitiveFileNames()) {
			return newPath + path[len(oldPath):], true
		}
		return "", false
	}

	return func(path string) (string, bool) {
		if original := l.tryGetSourcePosition(path, 0); original != nil {
			if updated, ok := getUpdatedPath(original.FileName); ok {
				return makeCorrespondingRelativeChange(original.FileName, updated, path, compareOptions), true
			}
		}
		return getUpdatedPath(path)
	}
}

func makeCorrespondingRelativeChange(a0 string, b0 string, a1 string, compareOptions tspath.ComparePathsOptions) string {
	rel := tspath.GetRelativePathFromFile(a0, b0, compareOptions)
	return tspath.CombinePaths(tspath.GetDirectoryPath(a1), rel)
}

func (l *LanguageService) updateTsconfigFiles(program *compiler.Program, changeTracker *change.Tracker, oldToNew pathUpdater, oldPath string, newPath string) {
	commandLine := program.CommandLine()
	if commandLine == nil || commandLine.ConfigFile == nil {
		return
	}

	configFile := commandLine.ConfigFile.SourceFile
	if configFile == nil {
		return
	}
	configDir := tspath.GetDirectoryPath(configFile.FileName())
	jsonObjectLiteral := getTsConfigObjectLiteralExpression(configFile)
	if jsonObjectLiteral == nil {
		return
	}

	forEachObjectProperty(jsonObjectLiteral, func(property *ast.PropertyAssignment, propertyName string) {
		switch propertyName {
		case "files", "include", "exclude":
			foundExactMatch := updatePathsProperty(configFile, configDir, property, changeTracker, oldToNew, l.converters, l.UseCaseSensitiveFileNames())
			if foundExactMatch || propertyName != "include" || !ast.IsArrayLiteralExpression(property.Initializer) {
				return
			}
			if oldSpec, isDefault := commandLine.GetMatchedIncludeSpec(oldPath); oldSpec != "" && !isDefault {
				if newSpec, _ := commandLine.GetMatchedIncludeSpec(newPath); newSpec == "" {
					elements := property.Initializer.Elements()
					if len(elements) > 0 {
						changeTracker.InsertNodeAfter(
							configFile,
							elements[len(elements)-1],
							changeTracker.NodeFactory.NewStringLiteral(relativePathFromDirectory(configDir, newPath, l.UseCaseSensitiveFileNames()), ast.TokenFlagsNone),
						)
					}
				}
			}
		case "compilerOptions":
			if !ast.IsObjectLiteralExpression(property.Initializer) {
				return
			}
			forEachObjectProperty(property.Initializer.AsObjectLiteralExpression(), func(property *ast.PropertyAssignment, propertyName string) {
				option := tsoptions.CommandLineCompilerOptionsMap.Get(propertyName)
				if option != nil {
					elementOption := option.Elements()
					if option.IsFilePath || (option.Kind == tsoptions.CommandLineOptionTypeList && elementOption != nil && elementOption.IsFilePath) {
						updatePathsProperty(configFile, configDir, property, changeTracker, oldToNew, l.converters, l.UseCaseSensitiveFileNames())
						return
					}
				}

				if propertyName != "paths" || !ast.IsObjectLiteralExpression(property.Initializer) {
					return
				}
				forEachObjectProperty(property.Initializer.AsObjectLiteralExpression(), func(pathsProperty *ast.PropertyAssignment, _ string) {
					if !ast.IsArrayLiteralExpression(pathsProperty.Initializer) {
						return
					}
					for _, element := range pathsProperty.Initializer.Elements() {
						tryUpdateConfigString(configFile, configDir, element, changeTracker, oldToNew, l.converters, l.UseCaseSensitiveFileNames())
					}
				})
			})
		}
	})
}

func updatePathsProperty(configFile *ast.SourceFile, configDir string, property *ast.PropertyAssignment, changeTracker *change.Tracker, oldToNew pathUpdater, converters *lsconv.Converters, useCaseSensitiveFileNames bool) bool {
	elements := []*ast.Node{property.Initializer}
	if ast.IsArrayLiteralExpression(property.Initializer) {
		elements = property.Initializer.Elements()
	}

	foundExactMatch := false
	for _, element := range elements {
		foundExactMatch = tryUpdateConfigString(configFile, configDir, element, changeTracker, oldToNew, converters, useCaseSensitiveFileNames) || foundExactMatch
	}
	return foundExactMatch
}

func tryUpdateConfigString(configFile *ast.SourceFile, configDir string, element *ast.Node, changeTracker *change.Tracker, oldToNew pathUpdater, converters *lsconv.Converters, useCaseSensitiveFileNames bool) bool {
	if !ast.IsStringLiteral(element) {
		return false
	}

	elementFileName := tspath.NormalizePath(tspath.CombinePaths(configDir, element.Text()))
	updated, ok := oldToNew(elementFileName)
	if !ok {
		return false
	}

	changeTracker.ReplaceRangeWithText(configFile, lsproto.Range{
		Start: converters.PositionToLineAndCharacter(configFile, core.TextPos(scanner.GetTokenPosOfNode(element, configFile, false)+1)),
		End:   converters.PositionToLineAndCharacter(configFile, core.TextPos(element.End()-1)),
	}, relativePathFromDirectory(configDir, updated, useCaseSensitiveFileNames))
	return true
}

func (l *LanguageService) updateImportsForFileRename(program *compiler.Program, changeTracker *change.Tracker, oldToNew pathUpdater, newToOld pathUpdater) {
	allFiles := program.GetSourceFiles()
	checker, done := program.GetTypeChecker(context.Background())
	defer done()
	moduleSpecifierPreferences := l.UserPreferences().ModuleSpecifierPreferences()

	for _, sourceFile := range allFiles {
		newFromOld, hasNewFromOld := oldToNew(sourceFile.FileName())
		oldFromNew, hasOldFromNew := newToOld(sourceFile.FileName())
		newImportFromPath := sourceFile.FileName()
		if hasNewFromOld {
			newImportFromPath = newFromOld
		}
		oldImportFromPath := sourceFile.FileName()
		if hasOldFromNew {
			oldImportFromPath = oldFromNew
		}
		importingSourceFileMoved := hasNewFromOld || hasOldFromNew

		for _, ref := range sourceFile.ReferencedFiles {
			if !tspath.IsExternalModuleNameRelative(ref.FileName) {
				continue
			}
			oldAbsolute := tspath.NormalizePath(tspath.CombinePaths(tspath.GetDirectoryPath(oldImportFromPath), ref.FileName))
			newAbsolute, ok := oldToNew(oldAbsolute)
			if !ok {
				continue
			}
			updated := relativeImportPathFromDirectory(tspath.GetDirectoryPath(newImportFromPath), newAbsolute, l.UseCaseSensitiveFileNames())
			if updated != ref.FileName {
				changeTracker.ReplaceRangeWithText(sourceFile, l.converters.ToLSPRange(sourceFile, ref.TextRange), updated)
			}
		}

		for _, importStringLiteral := range sourceFile.Imports() {
			updated := l.getUpdatedImportSpecifier(program, checker, sourceFile, importStringLiteral, oldToNew, newToOld, newImportFromPath, oldImportFromPath, importingSourceFileMoved, moduleSpecifierPreferences)
			if updated != "" && updated != importStringLiteral.Text() {
				changeTracker.ReplaceRangeWithText(sourceFile, l.converters.ToLSPRange(sourceFile, createStringTextRange(sourceFile, importStringLiteral)), updated)
			}
		}
	}
}

func (l *LanguageService) getUpdatedImportSpecifier(program *compiler.Program, checker interface {
	GetSymbolAtLocation(node *ast.Node) *ast.Symbol
}, sourceFile *ast.SourceFile, importLiteral *ast.StringLiteralLike, oldToNew pathUpdater, newToOld pathUpdater, newImportFromPath string, oldImportFromPath string, importingSourceFileMoved bool, userPreferences modulespecifiers.UserPreferences,
) string {
	importedModuleSymbol := checker.GetSymbolAtLocation(importLiteral)
	if isAmbientModuleSymbol(importedModuleSymbol) {
		return ""
	}

	if updated := getUpdatedImportSpecifierFromMovedSourceFiles(program, sourceFile, importLiteral, oldToNew, newImportFromPath, userPreferences); updated != "" && updated != importLiteral.Text() {
		return updated
	}

	var target *toImport
	if _, hasOldFromNew := newToOld(sourceFile.FileName()); hasOldFromNew {
		resolutionMode := program.GetModeForUsageLocation(sourceFile, importLiteral)
		target = getSourceFileToImportFromResolved(importLiteral, program.ResolveModuleName(importLiteral.Text(), oldImportFromPath, resolutionMode), oldToNew, program.GetSourceFiles())
	} else {
		target = getSourceFileToImport(program, importedModuleSymbol, sourceFile, importLiteral, oldToNew, userPreferences)
	}

	if target == nil {
		if importingSourceFileMoved && tspath.IsExternalModuleNameRelative(importLiteral.Text()) {
			absoluteTarget := tspath.NormalizePath(tspath.CombinePaths(tspath.GetDirectoryPath(sourceFile.FileName()), importLiteral.Text()))
			return relativeImportPathFromDirectory(tspath.GetDirectoryPath(newImportFromPath), absoluteTarget, l.UseCaseSensitiveFileNames())
		}
		return ""
	}

	if !target.updated && !(importingSourceFileMoved && tspath.IsExternalModuleNameRelative(importLiteral.Text())) {
		return ""
	}

	updated := modulespecifiers.UpdateModuleSpecifier(
		program.Options(),
		program,
		sourceFile,
		newImportFromPath,
		importLiteral.Text(),
		target.newFileName,
		userPreferences,
		modulespecifiers.ModuleSpecifierOptions{
			OverrideImportMode: program.GetModeForUsageLocation(sourceFile, importLiteral),
		},
	)
	return updated
}

func getSourceFileToImport(program *compiler.Program, importedModuleSymbol *ast.Symbol, sourceFile *ast.SourceFile, importLiteral *ast.StringLiteralLike, oldToNew pathUpdater, userPreferences modulespecifiers.UserPreferences) *toImport {
	if importedModuleSymbol != nil {
		if moduleSourceFile := core.Find(importedModuleSymbol.Declarations, ast.IsSourceFile); moduleSourceFile != nil {
			oldFileName := moduleSourceFile.AsSourceFile().FileName()
			if newFileName, ok := oldToNew(oldFileName); ok {
				return &toImport{newFileName: newFileName, updated: true}
			}
			return &toImport{newFileName: oldFileName, updated: false}
		}
	}

	if resolved := program.GetResolvedModuleFromModuleSpecifier(sourceFile, importLiteral); resolved != nil {
		return getSourceFileToImportFromResolved(importLiteral, resolved, oldToNew, program.GetSourceFiles())
	}

	resolutionMode := program.GetModeForUsageLocation(sourceFile, importLiteral)
	if resolved := program.ResolveModuleName(importLiteral.Text(), sourceFile.FileName(), resolutionMode); resolved != nil {
		return getSourceFileToImportFromResolved(importLiteral, resolved, oldToNew, program.GetSourceFiles())
	}

	return getSourceFileToImportFromMovedSourceFiles(program, sourceFile, importLiteral, oldToNew, resolutionMode, userPreferences)
}

func getSourceFileToImportFromResolved(importLiteral *ast.StringLiteralLike, resolved *module.ResolvedModule, oldToNew pathUpdater, sourceFiles []*ast.SourceFile) *toImport {
	if resolved == nil {
		return nil
	}

	if resolved.IsResolved() {
		if result := tryChange(resolved.ResolvedFileName, oldToNew); result != nil {
			return result
		}
	}

	for _, oldFileName := range resolved.FailedLookupLocations {
		if result := tryChangeWithIgnoringPackageJSONExisting(oldFileName, oldToNew, sourceFiles); result != nil {
			return result
		}
	}

	if tspath.IsExternalModuleNameRelative(importLiteral.Text()) {
		for _, oldFileName := range resolved.FailedLookupLocations {
			if result := tryChangeWithIgnoringPackageJSON(oldFileName, oldToNew); result != nil {
				return result
			}
		}
	}

	if resolved.IsResolved() {
		return &toImport{newFileName: resolved.ResolvedFileName, updated: false}
	}
	return nil
}

func tryChangeWithIgnoringPackageJSONExisting(oldFileName string, oldToNew pathUpdater, sourceFiles []*ast.SourceFile) *toImport {
	newFileName, ok := oldToNew(oldFileName)
	if !ok || !sourceFileExists(sourceFiles, newFileName) {
		return nil
	}
	return tryChangeWithIgnoringPackageJSON(oldFileName, oldToNew)
}

func tryChangeWithIgnoringPackageJSON(oldFileName string, oldToNew pathUpdater) *toImport {
	if strings.HasSuffix(oldFileName, "/package.json") {
		return nil
	}
	return tryChange(oldFileName, oldToNew)
}

func tryChange(oldFileName string, oldToNew pathUpdater) *toImport {
	if newFileName, ok := oldToNew(oldFileName); ok {
		return &toImport{newFileName: newFileName, updated: true}
	}
	return nil
}

func sourceFileExists(sourceFiles []*ast.SourceFile, fileName string) bool {
	for _, sourceFile := range sourceFiles {
		if sourceFile.FileName() == fileName {
			return true
		}
	}
	return false
}

func getSourceFileToImportFromMovedSourceFiles(program *compiler.Program, sourceFile *ast.SourceFile, importLiteral *ast.StringLiteralLike, oldToNew pathUpdater, resolutionMode core.ResolutionMode, userPreferences modulespecifiers.UserPreferences) *toImport {
	for _, candidate := range program.GetSourceFiles() {
		newFileName, ok := oldToNew(candidate.FileName())
		if !ok {
			continue
		}

		moduleSpecifier := modulespecifiers.UpdateModuleSpecifier(
			program.Options(),
			program,
			sourceFile,
			sourceFile.FileName(),
			importLiteral.Text(),
			candidate.FileName(),
			userPreferences,
			modulespecifiers.ModuleSpecifierOptions{
				OverrideImportMode: resolutionMode,
			},
		)
		if moduleSpecifier == importLiteral.Text() {
			return &toImport{newFileName: newFileName, updated: true}
		}
	}
	return nil
}

func getUpdatedImportSpecifierFromMovedSourceFiles(program *compiler.Program, sourceFile *ast.SourceFile, importLiteral *ast.StringLiteralLike, oldToNew pathUpdater, importingSourceFileName string, userPreferences modulespecifiers.UserPreferences) string {
	resolutionMode := program.GetModeForUsageLocation(sourceFile, importLiteral)
	for _, candidate := range program.GetSourceFiles() {
		newFileName, ok := oldToNew(candidate.FileName())
		if !ok {
			continue
		}

		oldSpecifier := modulespecifiers.UpdateModuleSpecifier(
			program.Options(),
			program,
			sourceFile,
			importingSourceFileName,
			importLiteral.Text(),
			candidate.FileName(),
			userPreferences,
			modulespecifiers.ModuleSpecifierOptions{
				OverrideImportMode: resolutionMode,
			},
		)
		if oldSpecifier != importLiteral.Text() {
			continue
		}

		return modulespecifiers.UpdateModuleSpecifier(
			program.Options(),
			program,
			sourceFile,
			importingSourceFileName,
			importLiteral.Text(),
			newFileName,
			userPreferences,
			modulespecifiers.ModuleSpecifierOptions{
				OverrideImportMode: resolutionMode,
			},
		)
	}
	return ""
}

func createStringTextRange(sourceFile *ast.SourceFile, node *ast.LiteralLikeNode) core.TextRange {
	return core.NewTextRange(scanner.GetTokenPosOfNode(node, sourceFile, false)+1, node.End()-1)
}

func getTsConfigObjectLiteralExpression(tsConfigSourceFile *ast.SourceFile) *ast.ObjectLiteralExpression {
	if tsConfigSourceFile != nil && tsConfigSourceFile.Statements != nil && len(tsConfigSourceFile.Statements.Nodes) > 0 {
		expression := tsConfigSourceFile.Statements.Nodes[0].Expression()
		if ast.IsObjectLiteralExpression(expression) {
			return expression.AsObjectLiteralExpression()
		}
	}
	return nil
}

func forEachObjectProperty(objectLiteral *ast.ObjectLiteralExpression, cb func(property *ast.PropertyAssignment, propertyName string)) {
	if objectLiteral == nil {
		return
	}
	for _, property := range objectLiteral.Properties.Nodes {
		if !ast.IsPropertyAssignment(property) {
			continue
		}
		if name, ok := ast.TryGetTextOfPropertyName(property.Name()); ok {
			cb(property.AsPropertyAssignment(), name)
		}
	}
}

func relativePathFromDirectory(fromDirectory string, to string, useCaseSensitiveFileNames bool) string {
	return tspath.GetRelativePathFromDirectory(fromDirectory, to, tspath.ComparePathsOptions{UseCaseSensitiveFileNames: useCaseSensitiveFileNames})
}

func relativeImportPathFromDirectory(fromDirectory string, to string, useCaseSensitiveFileNames bool) string {
	return tspath.EnsurePathIsNonModuleName(relativePathFromDirectory(fromDirectory, to, useCaseSensitiveFileNames))
}

func isAmbientModuleSymbol(symbol *ast.Symbol) bool {
	if symbol == nil {
		return false
	}
	return slices.ContainsFunc(symbol.Declarations, ast.IsModuleWithStringLiteralName)
}
