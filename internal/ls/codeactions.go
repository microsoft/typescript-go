package ls

import (
	"context"
	"fmt"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

// Import-related diagnostic codes that can be fixed with auto-import
var importFixDiagnosticCodes = map[int32]struct{}{
	diagnostics.Cannot_find_name_0.Code():                                          {},
	diagnostics.Cannot_find_name_0_Did_you_mean_1.Code():                           {},
	diagnostics.Cannot_find_namespace_0.Code():                                     {},
	diagnostics.Module_0_has_no_exported_member_1.Code():                           {},
	diagnostics.Cannot_find_module_0_or_its_corresponding_type_declarations.Code(): {},
}

// ProvideCodeActions returns code actions for the given range and context
func (l *LanguageService) ProvideCodeActions(ctx context.Context, params *lsproto.CodeActionParams) (lsproto.CodeActionResponse, error) {
	program, file := l.getProgramAndFile(params.TextDocument.Uri)

	var actions []lsproto.CommandOrCodeAction

	// Process diagnostics in the context to generate quick fixes
	if params.Context != nil && params.Context.Diagnostics != nil {
		for _, diag := range params.Context.Diagnostics {
			// Check if this is an import-related diagnostic
			if diag.Code.Integer != nil {
				if _, isImportDiag := importFixDiagnosticCodes[*diag.Code.Integer]; isImportDiag {
					importActions := l.getImportCodeActionsForDiagnostic(ctx, program, file, diag, params)
					actions = append(actions, importActions...)
				}
			}
		}
	}

	if len(actions) == 0 {
		return lsproto.CommandOrCodeActionArrayOrNull{CommandOrCodeActionArray: &[]lsproto.CommandOrCodeAction{}}, nil
	}

	return lsproto.CommandOrCodeActionArrayOrNull{CommandOrCodeActionArray: &actions}, nil
}

// getImportCodeActionsForDiagnostic generates auto-import code actions for a diagnostic
func (l *LanguageService) getImportCodeActionsForDiagnostic(
	ctx context.Context,
	program *compiler.Program,
	file *ast.SourceFile,
	diag *lsproto.Diagnostic,
	params *lsproto.CodeActionParams,
) []lsproto.CommandOrCodeAction {
	// Get the type checker
	ch, done := program.GetTypeCheckerForFile(ctx, file)
	if done != nil {
		defer done()
	}

	// Get the position from the diagnostic
	position := l.converters.LineAndCharacterToPosition(file, diag.Range.Start)

	// Get the token at the diagnostic position
	token := astnav.GetTokenAtPosition(file, int(position))
	if token == nil {
		return nil
	}

	// Token should be an identifier
	if !ast.IsIdentifier(token) {
		return nil
	}

	// Extract the symbol name from the identifier
	symbolName := token.Text()
	if symbolName == "" {
		return nil
	}

	// Get the meaning from the token location to filter exports properly
	meaning := getMeaningFromLocation(token)

	// Search for all exports matching this symbol name
	exportInfos := l.searchExportInfosForCodeAction(ctx, ch, file, symbolName, meaning)
	if len(exportInfos) == 0 {
		return nil
	}

	// Get all import fixes for the symbol, which will be sorted by preference
	isValidTypeOnlyUseSite := ast.IsValidTypeOnlyAliasUseSite(token)
	useRequire := getShouldUseRequire(file, program)
	_, fixes := l.getImportFixes(ch, exportInfos, &diag.Range.Start, &isValidTypeOnlyUseSite, &useRequire, file, false)

	var actions []lsproto.CommandOrCodeAction

	// Limit to top 3 import suggestions to avoid overwhelming the user
	maxActions := 3
	for i, fix := range fixes {
		if i >= maxActions {
			break
		}

		// Create the code action using the fix
		internalAction := l.codeActionForFix(ctx, file, symbolName, fix, false)

		// Convert internal code action to LSP CodeAction
		if len(internalAction.changes) > 0 {
			kind := lsproto.CodeActionKindQuickFix
			title := fmt.Sprintf("Import %s from \"%s\"", symbolName, fix.moduleSpecifier)

			changes := map[lsproto.DocumentUri][]*lsproto.TextEdit{
				params.TextDocument.Uri: internalAction.changes,
			}
			diagnostics := []*lsproto.Diagnostic{diag}
			codeAction := &lsproto.CodeAction{
				Title: title,
				Kind:  &kind,
				Edit: &lsproto.WorkspaceEdit{
					Changes: &changes,
				},
				Diagnostics: &diagnostics,
			}
			actions = append(actions, lsproto.CommandOrCodeAction{CodeAction: codeAction})
		}
	}

	return actions
}

// searchExportInfosForCodeAction searches for exports that match the given symbol name
// This follows the same pattern as TypeScript's getExportInfos in importFixes.ts
func (l *LanguageService) searchExportInfosForCodeAction(
	ctx context.Context,
	ch *checker.Checker,
	importingFile *ast.SourceFile,
	symbolName string,
	meaning ast.SemanticMeaning,
) []*SymbolExportInfo {
	var results []*SymbolExportInfo
	moduleCount := 0

	// Iterate through all available modules to find exports matching the symbol name
	forEachExternalModuleToImportFrom(
		ch,
		l.GetProgram(),
		func(moduleSymbol *ast.Symbol, moduleFile *ast.SourceFile, ch *checker.Checker, isFromPackageJson bool) {
			if moduleCount = moduleCount + 1; moduleCount%100 == 0 && ctx.Err() != nil {
				return
			}

			moduleFileName := ""
			if moduleFile != nil {
				moduleFileName = moduleFile.FileName()
			}

			// Check if the module has an export with the exact symbol name
			exportedSymbol := ch.TryGetMemberInModuleExportsAndProperties(symbolName, moduleSymbol)
			if exportedSymbol != nil && isImportableSymbol(exportedSymbol, ch) {
				// Filter by semantic meaning to match what the token location expects
				if symbolFlagsHaveMeaning(exportedSymbol.Flags, meaning) {
					results = append(results, &SymbolExportInfo{
						symbol:            exportedSymbol,
						moduleSymbol:      moduleSymbol,
						moduleFileName:    moduleFileName,
						exportKind:        ExportKindNamed,
						targetFlags:       ch.SkipAlias(exportedSymbol).Flags,
						isFromPackageJson: isFromPackageJson,
					})
				}
			}

			// Also check default exports
			defaultInfo := getDefaultLikeExportInfo(moduleSymbol, ch)
			if defaultInfo != nil && isImportableSymbol(defaultInfo.exportingModuleSymbol, ch) {
				// Check if the default export's name matches and has the right meaning
				if defaultInfo.exportingModuleSymbol.Name == symbolName &&
					symbolFlagsHaveMeaning(defaultInfo.exportingModuleSymbol.Flags, meaning) {
					results = append(results, &SymbolExportInfo{
						symbol:            defaultInfo.exportingModuleSymbol,
						moduleSymbol:      moduleSymbol,
						moduleFileName:    moduleFileName,
						exportKind:        defaultInfo.exportKind,
						targetFlags:       ch.SkipAlias(defaultInfo.exportingModuleSymbol).Flags,
						isFromPackageJson: isFromPackageJson,
					})
				}
			}
		},
	)

	return results
}
