package ls

import (
	"context"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
)

func (l *LanguageService) getExportInfoMap(
	ctx context.Context,
	ch *checker.Checker,
	importingFile *ast.SourceFile,
	preferences *UserPreferences,
) *exportInfoMap {
	// Pulling the AutoImportProvider project will trigger its updateGraph if pending,
	// which will invalidate the export map cache if things change, so pull it before
	// checking the cache.
	// l.GetPackageJsonAutoImportProvider?.();
	// cache := host.getCachedExportInfoMap()
	// || createCacheableExportInfoMap({
	//     getCurrentProgram: () => program,
	//     getPackageJsonAutoImportProvider: () => host.getPackageJsonAutoImportProvider?.(),
	//     getGlobalTypingsCacheLocation: () => host.getGlobalTypingsCacheLocation?.(),
	// });

	// if (cache.isUsableByFile(importingFile.path)) {
	//     host.log?.("getExportInfoMap: cache hit");
	//     return cache;
	// }

	// host.log?.("getExportInfoMap: expInfoMap miss or empty; calculating new results");
	expInfoMap := NewExportInfoMap(l.GetProgram().GetGlobalTypingsCacheLocation())
	moduleCount := 0
	forEachExternalModuleToImportFrom(
		ch,
		l.GetProgram(),
		preferences,
		// /*useAutoImportProvider*/ true,
		func(moduleSymbol *ast.Symbol, moduleFile *ast.SourceFile, ch *checker.Checker, isFromPackageJson bool) {
			if moduleCount = moduleCount + 1; moduleCount%100 == 0 && ctx.Err() != nil {
				return
			}
			seenExports := collections.Set[string]{}
			defaultInfo := getDefaultLikeExportInfo(moduleSymbol, ch)
			// Note: I think we shouldn't actually see resolved module symbols here, but weird merges
			// can cause it to happen: see 'completionsImport_mergedReExport.ts'
			if defaultInfo != nil && isImportableSymbol(defaultInfo.exportingModuleSymbol, ch) {
				expInfoMap.add(
					importingFile.Path(),
					defaultInfo.exportingModuleSymbol,
					core.IfElse(defaultInfo.exportKind == ExportKindDefault, ast.InternalSymbolNameDefault, ast.InternalSymbolNameExportEquals),
					moduleSymbol,
					moduleFile,
					defaultInfo.exportKind,
					isFromPackageJson,
					ch,
					nil,
					nil,
				)
			}
			var exportingModuleSymbol *ast.Symbol
			if defaultInfo != nil {
				exportingModuleSymbol = defaultInfo.exportingModuleSymbol
			}
			ch.ForEachExportAndPropertyOfModule(moduleSymbol, func(exported *ast.Symbol, key string) {
				if exported != exportingModuleSymbol && isImportableSymbol(exported, ch) && seenExports.AddIfAbsent(key) {
					expInfoMap.add(
						importingFile.Path(),
						exported,
						key,
						moduleSymbol,
						moduleFile,
						ExportKindNamed,
						isFromPackageJson,
						ch,
						nil,
						nil,
					)
				}
			})
		})

	// catch (err) {
	//     // Ensure cache is reset if operation is cancelled
	//     cache.clear();
	//     throw err;
	// }

	// host.log?.(`getExportInfoMap: done in ${timestamp() - start} ms`);
	return expInfoMap
}

// func (l *LanguageService) getAllExportInfoForSymbol(ctx context.Context, ch *checker.Checker, importingFile *ast.SourceFile, symbol *ast.Symbol, symbolName string, moduleSymbol *ast.Symbol, preferCapitalized bool, preferences *UserPreferences) []*SymbolExportInfo {
// 	// !!! isFileExcluded := len(preferences.AutoImportFileExcludePatterns) != 0 && getIsFileExcluded(host, preferences);
// 	// mergedModuleSymbol := ch.GetMergedSymbol(moduleSymbol)
// 	// moduleSourceFile := isFileExcluded && len(mergedModuleSymbol.Declarations) > 0 && ast.GetDeclarationOfKind(mergedModuleSymbol, SyntaxKind.SourceFile)
// 	// moduleSymbolExcluded := moduleSourceFile && isFileExcluded(moduleSourceFile.AsSourceFile());
// 	moduleSymbolExcluded := false
// 	return l.getExportInfoMap(ctx, ch, importingFile, preferences).search(
// 		ch,
// 		importingFile.Path(),
// 		preferCapitalized,
// 		func(name string, _ ast.SymbolFlags) bool { return name == symbolName },
// 		func(info []*SymbolExportInfo, symbolName string, isFromAmbientModule bool, key ExportInfoMapKey) []*SymbolExportInfo {
// 			if ch.GetMergedSymbol(ch.SkipAlias(info[0].symbol)) == symbol && (moduleSymbolExcluded || core.Some(info, func(i *SymbolExportInfo) bool {
// 				return ch.GetMergedSymbol(i.moduleSymbol) == moduleSymbol || i.symbol.Parent == moduleSymbol
// 			})) {
// 				return info
// 			}
// 			return nil
// 		},
// 	)
// }

func (l *LanguageService) getSingleExportInfoForSymbol(ch *checker.Checker, symbol *ast.Symbol, symbolName string, moduleSymbol *ast.Symbol) *SymbolExportInfo {
	getInfoWithChecker := func(program *compiler.Program, isFromPackageJson bool) *SymbolExportInfo {
		defaultInfo := getDefaultLikeExportInfo(moduleSymbol, ch)
		if defaultInfo != nil && ch.SkipAlias(defaultInfo.exportingModuleSymbol) == symbol {
			return &SymbolExportInfo{
				symbol:            defaultInfo.exportingModuleSymbol,
				moduleSymbol:      moduleSymbol,
				moduleFileName:    "",
				exportKind:        defaultInfo.exportKind,
				targetFlags:       ch.SkipAlias(symbol).Flags,
				isFromPackageJson: isFromPackageJson,
			}
		}
		if named := ch.TryGetMemberInModuleExportsAndProperties(symbolName, moduleSymbol); named != nil && ch.SkipAlias(named) == symbol {
			return &SymbolExportInfo{
				symbol:            named,
				moduleSymbol:      moduleSymbol,
				moduleFileName:    "",
				exportKind:        ExportKindNamed,
				targetFlags:       ch.SkipAlias(symbol).Flags,
				isFromPackageJson: isFromPackageJson,
			}
		}
		return nil
	}

	if mainProgramInfo := getInfoWithChecker(l.GetProgram() /*isFromPackageJson*/, false); mainProgramInfo != nil {
		return mainProgramInfo
	}
	// !!! autoImportProvider := host.getPackageJsonAutoImportProvider?.()?
	// return debug.CheckDefined(autoImportProvider && getInfoWithChecker(autoImportProvider, /*isFromPackageJson*/ true), `Could not find symbol in specified module for code actions`);
	return nil
}
