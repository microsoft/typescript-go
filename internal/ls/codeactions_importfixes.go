package ls

import (
	"cmp"
	"context"
	"fmt"
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/ls/change"
	"github.com/microsoft/typescript-go/internal/ls/organizeimports"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/outputpaths"
	"github.com/microsoft/typescript-go/internal/scanner"
	"github.com/microsoft/typescript-go/internal/tspath"
)

var importFixErrorCodes = []int32{
	diagnostics.Cannot_find_name_0.Code(),
	diagnostics.Cannot_find_name_0_Did_you_mean_1.Code(),
	diagnostics.Cannot_find_name_0_Did_you_mean_the_instance_member_this_0.Code(),
	diagnostics.Cannot_find_name_0_Did_you_mean_the_static_member_1_0.Code(),
	diagnostics.Cannot_find_namespace_0.Code(),
	diagnostics.X_0_refers_to_a_UMD_global_but_the_current_file_is_a_module_Consider_adding_an_import_instead.Code(),
	diagnostics.X_0_only_refers_to_a_type_but_is_being_used_as_a_value_here.Code(),
	diagnostics.No_value_exists_in_scope_for_the_shorthand_property_0_Either_declare_one_or_provide_an_initializer.Code(),
	diagnostics.X_0_cannot_be_used_as_a_value_because_it_was_imported_using_import_type.Code(),
	diagnostics.Cannot_find_name_0_Do_you_need_to_install_type_definitions_for_jQuery_Try_npm_i_save_dev_types_Slashjquery.Code(),
	diagnostics.Cannot_find_name_0_Do_you_need_to_change_your_target_library_Try_changing_the_lib_compiler_option_to_1_or_later.Code(),
	diagnostics.Cannot_find_name_0_Do_you_need_to_change_your_target_library_Try_changing_the_lib_compiler_option_to_include_dom.Code(),
	diagnostics.Cannot_find_name_0_Do_you_need_to_install_type_definitions_for_a_test_runner_Try_npm_i_save_dev_types_Slashjest_or_npm_i_save_dev_types_Slashmocha_and_then_add_jest_or_mocha_to_the_types_field_in_your_tsconfig.Code(),
	diagnostics.Cannot_find_name_0_Did_you_mean_to_write_this_in_an_async_function.Code(),
	diagnostics.Cannot_find_name_0_Do_you_need_to_install_type_definitions_for_jQuery_Try_npm_i_save_dev_types_Slashjquery_and_then_add_jquery_to_the_types_field_in_your_tsconfig.Code(),
	diagnostics.Cannot_find_name_0_Do_you_need_to_install_type_definitions_for_a_test_runner_Try_npm_i_save_dev_types_Slashjest_or_npm_i_save_dev_types_Slashmocha.Code(),
	diagnostics.Cannot_find_name_0_Do_you_need_to_install_type_definitions_for_node_Try_npm_i_save_dev_types_Slashnode.Code(),
	diagnostics.Cannot_find_name_0_Do_you_need_to_install_type_definitions_for_node_Try_npm_i_save_dev_types_Slashnode_and_then_add_node_to_the_types_field_in_your_tsconfig.Code(),
	diagnostics.Cannot_find_namespace_0_Did_you_mean_1.Code(),
	diagnostics.Cannot_extend_an_interface_0_Did_you_mean_implements.Code(),
	diagnostics.This_JSX_tag_requires_0_to_be_in_scope_but_it_could_not_be_found.Code(),
}

const (
	importFixID = "fixMissingImport"
)

// ImportFixProvider is the CodeFixProvider for import-related fixes
var ImportFixProvider = &CodeFixProvider{
	ErrorCodes:     importFixErrorCodes,
	GetCodeActions: getImportCodeActions,
	FixIds:         []string{importFixID},
}

type fixInfo struct {
	fix                 *ImportFix
	symbolName          string
	errorIdentifierText string
	isJsxNamespaceFix   bool
}

func getImportCodeActions(ctx context.Context, fixContext *CodeFixContext) []CodeAction {
	info := getFixInfos(ctx, fixContext, fixContext.ErrorCode, fixContext.Span.Pos(), true /* useAutoImportProvider */)
	if len(info) == 0 {
		return nil
	}

	var actions []CodeAction
	for _, fixInfo := range info {
		tracker := change.NewTracker(ctx, fixContext.Program.Options(), fixContext.LS.FormatOptions(), fixContext.LS.converters)
		msg := fixContext.LS.codeActionForFixWorker(
			tracker,
			fixContext.SourceFile,
			fixInfo.symbolName,
			fixInfo.fix,
			fixInfo.symbolName != fixInfo.errorIdentifierText,
		)

		if msg != nil {
			// Convert changes to LSP edits
			changes := tracker.GetChanges()
			var edits []*lsproto.TextEdit
			for _, fileChanges := range changes {
				edits = append(edits, fileChanges...)
			}

			actions = append(actions, CodeAction{
				Description: msg.Message(),
				Changes:     edits,
			})
		}
	}
	return actions
}

func getFixInfos(ctx context.Context, fixContext *CodeFixContext, errorCode int32, pos int, useAutoImportProvider bool) []*fixInfo {
	symbolToken := astnav.GetTokenAtPosition(fixContext.SourceFile, pos)

	var info []*fixInfo

	if errorCode == diagnostics.X_0_refers_to_a_UMD_global_but_the_current_file_is_a_module_Consider_adding_an_import_instead.Code() {
		info = getFixesInfoForUMDImport(ctx, fixContext, symbolToken)
	} else if !ast.IsIdentifier(symbolToken) {
		return nil
	} else if errorCode == diagnostics.X_0_cannot_be_used_as_a_value_because_it_was_imported_using_import_type.Code() {
		// Handle type-only import promotion
		ch, done := fixContext.Program.GetTypeChecker(ctx)
		defer done()
		compilerOptions := fixContext.Program.Options()
		symbolNames := getSymbolNamesToImport(fixContext.SourceFile, ch, symbolToken, compilerOptions)
		if len(symbolNames) != 1 {
			panic("Expected exactly one symbol name for type-only import promotion")
		}
		symbolName := symbolNames[0]
		fix := getTypeOnlyPromotionFix(ctx, fixContext.SourceFile, symbolToken, symbolName, fixContext.Program)
		if fix != nil {
			return []*fixInfo{{fix: fix, symbolName: symbolName, errorIdentifierText: symbolToken.Text()}}
		}
		return nil
	} else {
		info = getFixesInfoForNonUMDImport(ctx, fixContext, symbolToken, useAutoImportProvider)
	}

	// Sort fixes by preference
	return sortFixInfo(info, fixContext)
}

func getFixesInfoForUMDImport(ctx context.Context, fixContext *CodeFixContext, token *ast.Node) []*fixInfo {
	ch, done := fixContext.Program.GetTypeChecker(ctx)
	defer done()

	umdSymbol := getUmdSymbol(token, ch)
	if umdSymbol == nil {
		return nil
	}

	symbol := ch.GetAliasedSymbol(umdSymbol)
	symbolName := umdSymbol.Name
	exportInfo := []*SymbolExportInfo{{
		symbol:            umdSymbol,
		moduleSymbol:      symbol,
		moduleFileName:    "",
		exportKind:        ExportKindUMD,
		targetFlags:       symbol.Flags,
		isFromPackageJson: false,
	}}

	useRequire := shouldUseRequire(fixContext.SourceFile, fixContext.Program)
	// `usagePosition` is undefined because `token` may not actually be a usage of the symbol we're importing.
	// For example, we might need to import `React` in order to use an arbitrary JSX tag. We could send a position
	// for other UMD imports, but `usagePosition` is currently only used to insert a namespace qualification
	// before a named import, like converting `writeFile` to `fs.writeFile` (whether `fs` is already imported or
	// not), and this function will only be called for UMD symbols, which are necessarily an `export =`, not a
	// named export.
	_, fixes := fixContext.LS.getImportFixes(
		ch,
		exportInfo,
		nil, // usagePosition undefined for UMD
		ptrTo(false),
		&useRequire,
		fixContext.SourceFile,
		false, // fromCacheOnly
	)

	var result []*fixInfo
	for _, fix := range fixes {
		errorIdentifierText := ""
		if ast.IsIdentifier(token) {
			errorIdentifierText = token.Text()
		}
		result = append(result, &fixInfo{
			fix:                 fix,
			symbolName:          symbolName,
			errorIdentifierText: errorIdentifierText,
		})
	}
	return result
}

func getUmdSymbol(token *ast.Node, ch *checker.Checker) *ast.Symbol {
	// try the identifier to see if it is the umd symbol
	var umdSymbol *ast.Symbol
	if ast.IsIdentifier(token) {
		umdSymbol = ch.GetResolvedSymbol(token)
	}
	if isUMDExportSymbol(umdSymbol) {
		return umdSymbol
	}

	// The error wasn't for the symbolAtLocation, it was for the JSX tag itself, which needs access to e.g. `React`.
	parent := token.Parent
	if (ast.IsJsxOpeningLikeElement(parent) && parent.TagName() == token) ||
		ast.IsJsxOpeningFragment(parent) {
		var location *ast.Node
		if ast.IsJsxOpeningLikeElement(parent) {
			location = token
		} else {
			location = parent
		}
		jsxNamespace := ch.GetJsxNamespace(parent)
		parentSymbol := ch.ResolveName(jsxNamespace, location, ast.SymbolFlagsValue, false /* excludeGlobals */)
		if isUMDExportSymbol(parentSymbol) {
			return parentSymbol
		}
	}
	return nil
}

func isUMDExportSymbol(symbol *ast.Symbol) bool {
	return symbol != nil && len(symbol.Declarations) > 0 &&
		symbol.Declarations[0] != nil &&
		ast.IsNamespaceExportDeclaration(symbol.Declarations[0])
}

func getFixesInfoForNonUMDImport(ctx context.Context, fixContext *CodeFixContext, symbolToken *ast.Node, useAutoImportProvider bool) []*fixInfo {
	ch, done := fixContext.Program.GetTypeChecker(ctx)
	defer done()
	compilerOptions := fixContext.Program.Options()

	symbolNames := getSymbolNamesToImport(fixContext.SourceFile, ch, symbolToken, compilerOptions)
	var allInfo []*fixInfo

	for _, symbolName := range symbolNames {
		// "default" is a keyword and not a legal identifier for the import
		if symbolName == "default" {
			continue
		}

		isValidTypeOnlyUseSite := ast.IsValidTypeOnlyAliasUseSite(symbolToken)
		useRequire := shouldUseRequire(fixContext.SourceFile, fixContext.Program)
		exportInfos := getExportInfos(
			ctx,
			symbolName,
			isJSXTagName(symbolToken),
			getMeaningFromLocation(symbolToken),
			fixContext.SourceFile,
			fixContext.Program,
			fixContext.LS,
		)
		for exportInfoList := range exportInfos.Values() {
			for _, exportInfo := range exportInfoList {
				usagePos := scanner.GetTokenPosOfNode(symbolToken, fixContext.SourceFile, false)
				lspPos := fixContext.LS.converters.PositionToLineAndCharacter(fixContext.SourceFile, core.TextPos(usagePos))
				_, fixes := fixContext.LS.getImportFixes(
					ch,
					[]*SymbolExportInfo{exportInfo},
					&lspPos,
					&isValidTypeOnlyUseSite,
					&useRequire,
					fixContext.SourceFile,
					false, // fromCacheOnly
				)

				for _, fix := range fixes {
					allInfo = append(allInfo, &fixInfo{
						fix:                 fix,
						symbolName:          symbolName,
						errorIdentifierText: symbolToken.Text(),
						isJsxNamespaceFix:   symbolName != symbolToken.Text(),
					})
				}
			}
		}
	}

	return allInfo
}

func getTypeOnlyPromotionFix(ctx context.Context, sourceFile *ast.SourceFile, symbolToken *ast.Node, symbolName string, program *compiler.Program) *ImportFix {
	ch, done := program.GetTypeChecker(ctx)
	defer done()

	// Get the symbol at the token location
	symbol := ch.ResolveName(symbolName, symbolToken, ast.SymbolFlagsValue, true /* excludeGlobals */)
	if symbol == nil {
		return nil
	}

	// Get the type-only alias declaration
	typeOnlyAliasDeclaration := ch.GetTypeOnlyAliasDeclaration(symbol)
	if typeOnlyAliasDeclaration == nil || ast.GetSourceFileOfNode(typeOnlyAliasDeclaration) != sourceFile {
		return nil
	}

	return &ImportFix{
		kind:                     ImportFixKindPromoteTypeOnly,
		typeOnlyAliasDeclaration: typeOnlyAliasDeclaration,
	}
}

func getSymbolNamesToImport(sourceFile *ast.SourceFile, ch *checker.Checker, symbolToken *ast.Node, compilerOptions *core.CompilerOptions) []string {
	parent := symbolToken.Parent
	if (ast.IsJsxOpeningLikeElement(parent) || ast.IsJsxClosingElement(parent)) &&
		parent.TagName() == symbolToken &&
		jsxModeNeedsExplicitImport(compilerOptions.Jsx) {
		jsxNamespace := ch.GetJsxNamespace(sourceFile.AsNode())
		if needsJsxNamespaceFix(jsxNamespace, symbolToken, ch) {
			needsComponentNameFix := !scanner.IsIntrinsicJsxName(symbolToken.Text()) &&
				ch.ResolveName(symbolToken.Text(), symbolToken, ast.SymbolFlagsValue, false /* excludeGlobals */) == nil
			if needsComponentNameFix {
				return []string{symbolToken.Text(), jsxNamespace}
			}
			return []string{jsxNamespace}
		}
	}
	return []string{symbolToken.Text()}
}

func needsJsxNamespaceFix(jsxNamespace string, symbolToken *ast.Node, ch *checker.Checker) bool {
	if scanner.IsIntrinsicJsxName(symbolToken.Text()) {
		return true
	}
	namespaceSymbol := ch.ResolveName(jsxNamespace, symbolToken, ast.SymbolFlagsValue, true /* excludeGlobals */)
	if namespaceSymbol == nil {
		return true
	}
	// Check if all declarations are type-only
	if slices.ContainsFunc(namespaceSymbol.Declarations, ast.IsTypeOnlyImportOrExportDeclaration) {
		return (namespaceSymbol.Flags & ast.SymbolFlagsValue) == 0
	}
	return false
}

func jsxModeNeedsExplicitImport(jsx core.JsxEmit) bool {
	return jsx == core.JsxEmitReact || jsx == core.JsxEmitReactNative
}

func getExportInfos(
	ctx context.Context,
	symbolName string,
	isJsxTagName bool,
	currentTokenMeaning ast.SemanticMeaning,
	fromFile *ast.SourceFile,
	program *compiler.Program,
	ls *LanguageService,
) *collections.MultiMap[ast.SymbolId, *SymbolExportInfo] {
	// For each original symbol, keep all re-exports of that symbol together
	// Maps symbol id to info for modules providing that symbol (original export + re-exports)
	originalSymbolToExportInfos := &collections.MultiMap[ast.SymbolId, *SymbolExportInfo]{}

	ch, done := program.GetTypeChecker(ctx)
	defer done()

	packageJsonFilter := ls.createPackageJsonImportFilter(fromFile)

	// Helper to add a symbol to the results map
	addSymbol := func(moduleSymbol *ast.Symbol, toFile *ast.SourceFile, exportedSymbol *ast.Symbol, exportKind ExportKind, isFromPackageJson bool) {
		if !ls.isImportable(fromFile, toFile, moduleSymbol, packageJsonFilter) {
			return
		}

		// Get unique ID for the exported symbol
		symbolID := ast.GetSymbolId(exportedSymbol)

		moduleFileName := ""
		if toFile != nil {
			moduleFileName = toFile.FileName()
		}

		originalSymbolToExportInfos.Add(symbolID, &SymbolExportInfo{
			symbol:            exportedSymbol,
			moduleSymbol:      moduleSymbol,
			moduleFileName:    moduleFileName,
			exportKind:        exportKind,
			targetFlags:       ch.SkipAlias(exportedSymbol).Flags,
			isFromPackageJson: isFromPackageJson,
		})
	}

	// Iterate through all external modules
	forEachExternalModuleToImportFrom(
		ch,
		program,
		func(moduleSymbol *ast.Symbol, sourceFile *ast.SourceFile, checker *checker.Checker, isFromPackageJson bool) {
			// Check for cancellation
			if ctx.Err() != nil {
				return
			}

			compilerOptions := program.Options()

			// Check default export
			defaultInfo := getDefaultLikeExportInfo(moduleSymbol, checker)
			if defaultInfo != nil &&
				symbolFlagsHaveMeaning(checker.GetSymbolFlags(defaultInfo.exportingModuleSymbol), currentTokenMeaning) &&
				forEachNameOfDefaultExport(defaultInfo.exportingModuleSymbol, checker, compilerOptions.GetEmitScriptTarget(), func(name, capitalizedName string) string {
					actualName := name
					if isJsxTagName && capitalizedName != "" {
						actualName = capitalizedName
					}
					if actualName == symbolName {
						return actualName
					}
					return ""
				}) != "" {
				addSymbol(moduleSymbol, sourceFile, defaultInfo.exportingModuleSymbol, defaultInfo.exportKind, isFromPackageJson)
			}
			// Check for named export with identical name
			exportSymbol := checker.TryGetMemberInModuleExportsAndProperties(symbolName, moduleSymbol)
			if exportSymbol != nil && symbolFlagsHaveMeaning(checker.GetSymbolFlags(exportSymbol), currentTokenMeaning) {
				addSymbol(moduleSymbol, sourceFile, exportSymbol, ExportKindNamed, isFromPackageJson)
			}
		},
	)

	return originalSymbolToExportInfos
}

func shouldUseRequire(sourceFile *ast.SourceFile, program *compiler.Program) bool {
	// Delegate to the existing implementation in autoimports.go
	return getShouldUseRequire(sourceFile, program)
}

func isJSXTagName(node *ast.Node) bool {
	parent := node.Parent
	if parent == nil {
		return false
	}
	if ast.IsJsxOpeningLikeElement(parent) || ast.IsJsxClosingElement(parent) {
		return parent.TagName() == node
	}
	return false
}

// getModuleSpecifierFromDeclaration gets the module specifier string from a declaration
func getModuleSpecifierFromDeclaration(decl *ast.Declaration) string {
	var moduleSpec *ast.Node

	switch decl.Kind {
	case ast.KindImportDeclaration:
		moduleSpec = decl.AsImportDeclaration().ModuleSpecifier
	case ast.KindImportEqualsDeclaration:
		importEq := decl.AsImportEqualsDeclaration()
		if importEq.ModuleReference != nil && importEq.ModuleReference.Kind == ast.KindExternalModuleReference {
			moduleSpec = importEq.ModuleReference.AsExternalModuleReference().Expression
		}
	case ast.KindImportSpecifier:
		// Walk up to find the import declaration
		if clause := getImportClauseOfSpecifier(decl.AsImportSpecifier()); clause != nil {
			if clause.Parent != nil && clause.Parent.Kind == ast.KindImportDeclaration {
				moduleSpec = clause.Parent.AsImportDeclaration().ModuleSpecifier
			}
		}
	case ast.KindImportClause:
		if decl.Parent != nil && decl.Parent.Kind == ast.KindImportDeclaration {
			moduleSpec = decl.Parent.AsImportDeclaration().ModuleSpecifier
		}
	case ast.KindNamespaceImport:
		if decl.Parent != nil && decl.Parent.Kind == ast.KindImportClause {
			if decl.Parent.Parent != nil && decl.Parent.Parent.Kind == ast.KindImportDeclaration {
				moduleSpec = decl.Parent.Parent.AsImportDeclaration().ModuleSpecifier
			}
		}
	}

	if moduleSpec != nil && ast.IsStringLiteral(moduleSpec) {
		return moduleSpec.AsStringLiteral().Text
	}
	return ""
}

// getImportClauseOfSpecifier gets the import clause containing a specifier
func getImportClauseOfSpecifier(spec *ast.ImportSpecifier) *ast.ImportClause {
	if spec.Parent != nil && spec.Parent.Kind == ast.KindNamedImports {
		if spec.Parent.Parent != nil && spec.Parent.Parent.Kind == ast.KindImportClause {
			return spec.Parent.Parent.AsImportClause()
		}
	}
	return nil
}

func sortFixInfo(fixes []*fixInfo, fixContext *CodeFixContext) []*fixInfo {
	if len(fixes) == 0 {
		return fixes
	}

	// Create a copy to avoid modifying the original
	sorted := make([]*fixInfo, len(fixes))
	copy(sorted, fixes)

	// Create package.json filter for import filtering
	packageJsonFilter := fixContext.LS.createPackageJsonImportFilter(fixContext.SourceFile)

	// Sort by:
	// 1. JSX namespace fixes last
	// 2. Fix kind (UseNamespace and AddToExisting preferred)
	// 3. Module specifier comparison
	slices.SortFunc(sorted, func(a, b *fixInfo) int {
		// JSX namespace fixes should come last
		if cmp := core.CompareBooleans(a.isJsxNamespaceFix, b.isJsxNamespaceFix); cmp != 0 {
			return cmp
		}

		// Compare fix kinds (lower is better)
		if cmp := cmp.Compare(int(a.fix.kind), int(b.fix.kind)); cmp != 0 {
			return cmp
		}

		// Compare module specifiers
		return fixContext.LS.compareModuleSpecifiers(
			a.fix,
			b.fix,
			fixContext.SourceFile,
			packageJsonFilter.allowsImportingSpecifier,
			func(fileName string) tspath.Path { return tspath.Path(fileName) },
		)
	})

	return sorted
}

func promoteFromTypeOnly(
	changes *change.Tracker,
	aliasDeclaration *ast.Declaration,
	program *compiler.Program,
	sourceFile *ast.SourceFile,
	ls *LanguageService,
) *ast.Declaration {
	compilerOptions := program.Options()
	// See comment in `doAddExistingFix` on constant with the same name.
	convertExistingToTypeOnly := compilerOptions.VerbatimModuleSyntax

	switch aliasDeclaration.Kind {
	case ast.KindImportSpecifier:
		spec := aliasDeclaration.AsImportSpecifier()
		if spec.IsTypeOnly {
			// TODO: TypeScript creates a new specifier with isTypeOnly=false, computes insertion index,
			// and if different from current position, deletes and re-inserts at new position.
			// We just delete the type keyword which may result in incorrectly ordered imports.
			// Also, TypeScript uses changes.deleteRange with getTokenPosOfNode for more precise range,
			// we use a simpler approach.

			// If there are multiple specifiers, we might need to move this one
			if spec.Parent != nil && spec.Parent.Kind == ast.KindNamedImports {
				namedImports := spec.Parent.AsNamedImports()
				if len(namedImports.Elements.Nodes) > 1 {
					// For now, just remove the 'type' keyword from this specifier
					// Full implementation would handle reordering
					deleteTypeKeywordFromSpecifier(changes, sourceFile, spec)
				} else {
					// Single specifier - just remove the 'type' keyword
					deleteTypeKeywordFromSpecifier(changes, sourceFile, spec)
				}
			}
			return aliasDeclaration
		} else {
			// The parent import clause is type-only
			if spec.Parent == nil || spec.Parent.Kind != ast.KindNamedImports {
				panic("ImportSpecifier parent must be NamedImports")
			}
			if spec.Parent.Parent == nil || spec.Parent.Parent.Kind != ast.KindImportClause {
				panic("NamedImports parent must be ImportClause")
			}
			promoteImportClause(changes, spec.Parent.Parent.AsImportClause(), program, sourceFile, ls, convertExistingToTypeOnly, aliasDeclaration)
			return spec.Parent.Parent
		}

	case ast.KindImportClause:
		promoteImportClause(changes, aliasDeclaration.AsImportClause(), program, sourceFile, ls, convertExistingToTypeOnly, aliasDeclaration)
		return aliasDeclaration

	case ast.KindNamespaceImport:
		// Promote the parent import clause
		if aliasDeclaration.Parent == nil || aliasDeclaration.Parent.Kind != ast.KindImportClause {
			panic("NamespaceImport parent must be ImportClause")
		}
		promoteImportClause(changes, aliasDeclaration.Parent.AsImportClause(), program, sourceFile, ls, convertExistingToTypeOnly, aliasDeclaration)
		return aliasDeclaration.Parent

	case ast.KindImportEqualsDeclaration:
		// Remove the 'type' keyword (which is the second child: 'import' 'type' name '=' ...)
		deleteTypeKeywordFromImportEquals(changes, sourceFile, aliasDeclaration.AsImportEqualsDeclaration())
		return aliasDeclaration
	default:
		panic(fmt.Sprintf("Unexpected alias declaration kind: %v", aliasDeclaration.Kind))
	}
}

// promoteImportClause removes the type keyword from an import clause
func promoteImportClause(
	changes *change.Tracker,
	importClause *ast.ImportClause,
	program *compiler.Program,
	sourceFile *ast.SourceFile,
	ls *LanguageService,
	convertExistingToTypeOnly core.Tristate,
	aliasDeclaration *ast.Declaration,
) {
	// Delete the 'type' keyword
	if importClause.PhaseModifier == ast.KindTypeKeyword {
		deleteTypeKeywordFromImportClause(changes, sourceFile, importClause)
	}

	// Handle .ts extension conversion to .js if necessary
	compilerOptions := program.Options()
	if compilerOptions.AllowImportingTsExtensions.IsFalse() {
		moduleSpecifier := checker.TryGetModuleSpecifierFromDeclaration(importClause.Parent)
		if moduleSpecifier != nil {
			resolvedModule := program.GetResolvedModuleFromModuleSpecifier(sourceFile, moduleSpecifier)
			if resolvedModule != nil && resolvedModule.ResolvedUsingTsExtension {
				moduleText := moduleSpecifier.AsStringLiteral().Text
				changedExtension := tspath.ChangeExtension(
					moduleText,
					outputpaths.GetOutputExtension(moduleText, compilerOptions.Jsx),
				)
				// Replace the module specifier with the new extension
				// We need to update the string literal, keeping the quotes
				replaceStringLiteral(changes, sourceFile, moduleSpecifier, changedExtension)
			}
		}
	}

	// Handle verbatimModuleSyntax conversion
	// If convertExistingToTypeOnly is true, we need to add 'type' to other specifiers
	// in the same import declaration
	if convertExistingToTypeOnly.IsTrue() {
		namedImports := importClause.NamedBindings
		if namedImports != nil && namedImports.Kind == ast.KindNamedImports {
			namedImportsData := namedImports.AsNamedImports()
			if len(namedImportsData.Elements.Nodes) > 1 {
				// Check if the list is sorted and if we need to reorder
				_, isSorted := organizeimports.GetNamedImportSpecifierComparerWithDetection(
					importClause.Parent,
					sourceFile,
					ls.UserPreferences(),
				)

				// If the alias declaration is an ImportSpecifier and the list is sorted,
				// move it to index 0 (since it will be the only non-type-only import)
				if isSorted.IsFalse() == false && // isSorted !== false
					aliasDeclaration != nil &&
					aliasDeclaration.Kind == ast.KindImportSpecifier {
					// Find the index of the alias declaration
					aliasIndex := -1
					for i, element := range namedImportsData.Elements.Nodes {
						if element == aliasDeclaration {
							aliasIndex = i
							break
						}
					}
					// If not already at index 0, move it there
					if aliasIndex > 0 {
						// Delete the specifier from its current position
						deleteNode(changes, sourceFile, aliasDeclaration, namedImportsData.Elements.Nodes)
						// Insert it at index 0
						changes.InsertImportSpecifierAtIndex(sourceFile, aliasDeclaration, namedImports, 0)
					}
				}

				// Add 'type' keyword to all other import specifiers that aren't already type-only
				for _, element := range namedImportsData.Elements.Nodes {
					spec := element.AsImportSpecifier()
					// Skip the specifier being promoted (if aliasDeclaration is an ImportSpecifier)
					if aliasDeclaration != nil && aliasDeclaration.Kind == ast.KindImportSpecifier {
						if element == aliasDeclaration {
							continue
						}
					}
					// Skip if already type-only
					if !spec.IsTypeOnly {
						insertTypeModifierBefore(changes, sourceFile, element)
					}
				}
			}
		}
	}
}

// deleteTypeKeywordFromImportClause deletes the 'type' keyword from an import clause
func deleteTypeKeywordFromImportClause(changes *change.Tracker, sourceFile *ast.SourceFile, importClause *ast.ImportClause) {
	// The type keyword is the first token in the import clause
	// import type { foo } from "bar"
	//        ^^^^^ - this keyword
	// We need to find and delete "type " (including trailing space)

	// The import clause starts at the position of the "type" keyword
	// Use the scanner to get the token at that position
	scan := scanner.GetScannerForSourceFile(sourceFile, importClause.Pos())
	token := scan.Token()

	if token != ast.KindTypeKeyword {
		panic(fmt.Sprintf("Expected type keyword at import clause start, got %v at pos %d", token, importClause.Pos()))
	}

	// Use TokenStart (not TokenFullStart) to avoid including leading whitespace
	typeStart := scan.TokenStart()
	typeEnd := scan.TokenEnd()

	// Skip whitespace after 'type' to include it in the deletion
	text := sourceFile.Text()
	endPos := typeEnd
	for endPos < len(text) && (text[endPos] == ' ' || text[endPos] == '\t') {
		endPos++
	}

	// Convert text positions to LSP positions
	// Note: changes.Converters() provides access to the converters
	conv := changes.Converters()
	lspRange := conv.ToLSPRange(sourceFile, core.NewTextRange(typeStart, endPos))

	changes.ReplaceRangeWithText(sourceFile, lspRange, "")
}

// deleteTypeKeywordFromSpecifier deletes the 'type' keyword from an import specifier
func deleteTypeKeywordFromSpecifier(changes *change.Tracker, sourceFile *ast.SourceFile, spec *ast.ImportSpecifier) {
	// import { type foo } from "bar"
	//          ^^^^^ - this keyword and space after

	specStart := spec.Pos()

	// The 'type' keyword is at the start of the specifier
	// We want to delete "type " (including the trailing space)
	typeKeywordEnd := specStart + 4 // length of "type"

	// Skip whitespace after 'type'
	text := sourceFile.Text()
	endPos := typeKeywordEnd
	for endPos < len(text) && (text[endPos] == ' ' || text[endPos] == '\t') {
		endPos++
	}

	// Convert text positions to LSP positions
	conv := changes.Converters()
	lspRange := conv.ToLSPRange(sourceFile, core.NewTextRange(specStart, endPos))

	changes.ReplaceRangeWithText(sourceFile, lspRange, "")
}

// deleteTypeKeywordFromImportEquals deletes the 'type' keyword from an import equals declaration
func deleteTypeKeywordFromImportEquals(changes *change.Tracker, sourceFile *ast.SourceFile, decl *ast.ImportEqualsDeclaration) {
	// import type Foo = require("bar")
	//        ^^^^^ - this keyword and space after

	// The 'type' keyword comes after 'import' and before the name
	// We need to find it by looking at the text
	declStart := decl.Pos()
	text := sourceFile.Text()

	// Skip 'import' keyword and whitespace
	pos := declStart + 6 // length of "import"
	for pos < len(text) && (text[pos] == ' ' || text[pos] == '\t') {
		pos++
	}

	// Now we should be at 'type'
	typeStart := pos
	typeEnd := pos + 4 // length of "type"

	// Skip whitespace after 'type'
	for typeEnd < len(text) && (text[typeEnd] == ' ' || text[typeEnd] == '\t') {
		typeEnd++
	}

	// Convert text positions to LSP positions
	conv := changes.Converters()
	lspRange := conv.ToLSPRange(sourceFile, core.NewTextRange(typeStart, typeEnd))

	changes.ReplaceRangeWithText(sourceFile, lspRange, "")
}

func replaceStringLiteral(changes *change.Tracker, sourceFile *ast.SourceFile, stringLiteral *ast.Node, newText string) {
	// Get the position of the string literal content (excluding quotes)
	literalStart := stringLiteral.Pos()
	literalEnd := stringLiteral.End()

	// Determine the quote character used
	text := sourceFile.Text()
	quoteChar := text[literalStart]

	// Create the new string literal with quotes
	newLiteral := string(quoteChar) + newText + string(quoteChar)

	// Convert text positions to LSP positions
	conv := changes.Converters()
	lspRange := conv.ToLSPRange(sourceFile, core.NewTextRange(literalStart, literalEnd))

	changes.ReplaceRangeWithText(sourceFile, lspRange, newLiteral)
}

func insertTypeModifierBefore(changes *change.Tracker, sourceFile *ast.SourceFile, specifier *ast.Node) {
	// Insert "type " before the specifier
	// import { foo } => import { type foo }
	specStart := specifier.Pos()

	// Convert text position to LSP position
	conv := changes.Converters()
	lspPos := conv.PositionToLineAndCharacter(sourceFile, core.TextPos(specStart))

	// Insert "type " at the beginning of the specifier
	changes.InsertText(sourceFile, lspPos, "type ")
}

func deleteNode(changes *change.Tracker, sourceFile *ast.SourceFile, node *ast.Node, containingList []*ast.Node) {
	// Find the node in the list to determine if we need to delete a comma
	nodeIndex := -1
	for i, n := range containingList {
		if n == node {
			nodeIndex = i
			break
		}
	}

	if nodeIndex == -1 {
		return // Node not found in list
	}

	// Determine the range to delete
	start := node.Pos()
	end := node.End()

	// If this is not the last element, we need to include the comma after it
	if nodeIndex < len(containingList)-1 {
		// Find and include the comma after this element
		text := sourceFile.Text()
		pos := end
		// Skip whitespace to find the comma
		for pos < len(text) && (text[pos] == ' ' || text[pos] == '\t' || text[pos] == '\n' || text[pos] == '\r') {
			pos++
		}
		if pos < len(text) && text[pos] == ',' {
			end = pos + 1
			// Also skip trailing whitespace after comma
			for end < len(text) && (text[end] == ' ' || text[end] == '\t') {
				end++
			}
		}
	} else if nodeIndex > 0 {
		// This is the last element - include the comma before it
		text := sourceFile.Text()
		pos := start - 1
		// Skip whitespace backwards to find the comma
		for pos >= 0 && (text[pos] == ' ' || text[pos] == '\t' || text[pos] == '\n' || text[pos] == '\r') {
			pos--
		}
		if pos >= 0 && text[pos] == ',' {
			start = pos
		}
	}

	// Convert text positions to LSP positions
	conv := changes.Converters()
	lspRange := conv.ToLSPRange(sourceFile, core.NewTextRange(start, end))

	changes.ReplaceRangeWithText(sourceFile, lspRange, "")
}
