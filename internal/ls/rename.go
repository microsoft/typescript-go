package ls

import (
	"context"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// RenameInfo represents the result of a rename validation check.
// It is used by the `textDocument/prepareRename` LSP handler.
type RenameInfo struct {
	CanRename bool
	// !!! LocalizedErrorMessage is not currently surfaced via the LSP prepareRename response,
	// !!! which only supports returning null to indicate failure. If the LSP spec adds error
	// !!! message support to prepareRename, this field and the diagnostic messages in
	// !!! renameBlockedReason/wouldRenameInOtherNodeModules should be restored.
	DisplayName string
	TriggerSpan lsproto.Range
}

func (l *LanguageService) ProvideRename(ctx context.Context, params *lsproto.RenameParams, orchestrator CrossProjectOrchestrator) (lsproto.WorkspaceEditOrNull, error) {
	return handleCrossProject(
		l,
		ctx,
		params,
		orchestrator,
		(*LanguageService).symbolAndEntriesToRename,
		combineRenameResponse,
		true,  /*isRename*/
		false, /*implementations*/
		symbolEntryTransformOptions{},
	)
}

func (l *LanguageService) GetRenameInfo(ctx context.Context, documentURI lsproto.DocumentUri, position lsproto.Position) RenameInfo {
	program, sourceFile := l.getProgramAndFile(documentURI)
	pos := int(l.converters.LineAndCharacterToPosition(sourceFile, position))

	node := astnav.GetTouchingPropertyName(sourceFile, pos)
	node = getAdjustedLocation(node, true /*forRename*/, sourceFile)

	if nodeIsEligibleForRename(node) {
		if renameInfo, ok := l.getRenameInfoForNode(ctx, node, sourceFile, program); ok {
			return renameInfo
		}
	}
	// !!! diagnostics.You_cannot_rename_this_element
	return RenameInfo{}
}

func (l *LanguageService) symbolAndEntriesToRename(ctx context.Context, params *lsproto.RenameParams, data SymbolAndEntriesData, options symbolEntryTransformOptions) (lsproto.WorkspaceEditOrNull, error) {
	if !nodeIsEligibleForRename(data.OriginalNode) {
		return lsproto.WorkspaceEditOrNull{}, nil
	}

	program := l.GetProgram()

	// Validate rename eligibility even if the client skipped textDocument/prepareRename.
	if !l.canRenameNode(ctx, data.OriginalNode, program) {
		return lsproto.WorkspaceEditOrNull{}, nil
	}

	entries := core.FlatMap(data.SymbolsAndEntries, func(s *SymbolAndEntries) []*ReferenceEntry { return s.references })
	changes := make(map[lsproto.DocumentUri][]*lsproto.TextEdit)
	ch, done := program.GetTypeChecker(ctx)
	defer done()

	for _, entry := range entries {
		uri := l.getFileNameOfEntry(entry)
		if l.UserPreferences().AllowRenameOfImportPath != core.TSTrue && entry.node != nil && ast.IsStringLiteralLike(entry.node) && ast.TryGetImportFromModuleSpecifier(entry.node) != nil {
			continue
		}
		textEdit := &lsproto.TextEdit{
			Range:   *l.getRangeOfEntry(entry),
			NewText: l.getTextForRename(data.OriginalNode, entry, params.NewName, ch),
		}
		changes[uri] = append(changes[uri], textEdit)
	}
	return lsproto.WorkspaceEditOrNull{
		WorkspaceEdit: &lsproto.WorkspaceEdit{
			Changes: &changes,
		},
	}, nil
}

// getRenameInfoForNode performs detailed validation for a rename operation on a specific node.
func (l *LanguageService) getRenameInfoForNode(ctx context.Context, node *ast.Node, sourceFile *ast.SourceFile, program *compiler.Program) (RenameInfo, bool) {
	ch, done := program.GetTypeChecker(ctx)
	defer done()

	symbol := ch.GetSymbolAtLocation(node)
	if symbol == nil {
		if ast.IsStringLiteralLike(node) {
			// Allow renaming of string literal types with contextual string literal types
			typ := getContextualTypeFromParent(node, ch, checker.ContextFlagsNone)
			if typ != nil && (typ.Flags()&checker.TypeFlagsStringLiteral != 0 ||
				(typ.Flags()&checker.TypeFlagsUnion != 0 && core.Every(typ.Types(), func(t *checker.Type) bool {
					return t.Flags()&checker.TypeFlagsStringLiteral != 0
				}))) {
				return getRenameInfoSuccess(node, sourceFile, node.Text(), l.converters), true
			}
		} else if ast.IsLabelName(node) {
			name := node.Text()
			return getRenameInfoSuccess(node, sourceFile, name, l.converters), true
		}
		return RenameInfo{}, false
	}

	// Only allow a symbol to be renamed if it actually has at least one declaration.
	if len(symbol.Declarations) == 0 {
		return RenameInfo{}, false
	}

	if l.isRenameBlocked(sourceFile, node, symbol, ch, program) {
		return RenameInfo{}, true
	}

	if ast.IsStringLiteralLike(node) && ast.TryGetImportFromModuleSpecifier(node) != nil {
		if l.UserPreferences().AllowRenameOfImportPath.IsTrue() {
			return l.getRenameInfoForModule(node, sourceFile, symbol)
		}
		return RenameInfo{}, false
	}

	return getRenameInfoSuccess(node, sourceFile, ch.SymbolToString(symbol), l.converters), true
}

func nodeIsEligibleForRename(node *ast.Node) bool {
	switch node.Kind {
	case ast.KindIdentifier,
		ast.KindPrivateIdentifier,
		ast.KindStringLiteral,
		ast.KindNoSubstitutionTemplateLiteral:
		return true
	// !!! The reference code also allows ThisKeyword and (conditionally) NumericLiteral.
	default:
		return false
	}
}

func (l *LanguageService) canRenameNode(ctx context.Context, node *ast.Node, program *compiler.Program) bool {
	ch, done := program.GetTypeChecker(ctx)
	defer done()

	symbol := ch.GetSymbolAtLocation(node)
	if symbol == nil {
		return true // Let the rename flow handle no-symbol cases (labels, string literals)
	}

	if l.isRenameBlocked(ast.GetSourceFileOfNode(node), node, symbol, ch, program) {
		return false
	}

	return true
}

// isRenameBlocked returns true if the rename should be blocked
// because the symbol is a library definition, a default keyword, or would cross node_modules boundaries.
func (l *LanguageService) isRenameBlocked(sourceFile *ast.SourceFile, node *ast.Node, symbol *ast.Symbol, ch *checker.Checker, program *compiler.Program) bool {
	for _, declaration := range symbol.Declarations {
		if isDefinedInLibraryFile(program, declaration) {
			// !!! diagnostics.You_cannot_rename_elements_that_are_defined_in_the_standard_TypeScript_library
			return true
		}
	}

	// Cannot rename `default` as in `import { default as foo } from "./someModule"`
	if ast.IsIdentifier(node) && node.Text() == "default" && symbol.Parent != nil && symbol.Parent.Flags&ast.SymbolFlagsModule != 0 {
		// !!! diagnostics.You_cannot_rename_this_element
		return true
	}

	if wouldRenameInOtherNodeModules(sourceFile, symbol, ch, l.UserPreferences()) {
		return true
	}

	return false
}

// isDefinedInLibraryFile checks if a declaration is from a default library file (e.g., lib.d.ts).
func isDefinedInLibraryFile(program *compiler.Program, declaration *ast.Node) bool {
	declSourceFile := ast.GetSourceFileOfNode(declaration)
	return program.IsSourceFileDefaultLibrary(declSourceFile.Path()) && tspath.IsDeclarationFileName(declSourceFile.FileName())
}

// wouldRenameInOtherNodeModules checks if renaming the symbol would affect node_modules.
func wouldRenameInOtherNodeModules(originalFile *ast.SourceFile, symbol *ast.Symbol, ch *checker.Checker, preferences *lsutil.UserPreferences) bool {
	sym := symbol
	if !preferences.UseAliasesForRename.IsTrue() && sym.Flags&ast.SymbolFlagsAlias != 0 {
		importSpecifier := core.Find(sym.Declarations, ast.IsImportSpecifier)
		if importSpecifier != nil && importSpecifier.AsImportSpecifier().PropertyName == nil {
			sym = ch.GetAliasedSymbol(sym)
		}
	}

	declarations := sym.Declarations
	if len(declarations) == 0 {
		return false
	}

	originalPackage := getPackagePathComponents(originalFile.FileName())
	if originalPackage == nil {
		// Original source file is not in node_modules.
		for _, declaration := range declarations {
			if isInsideNodeModules(ast.GetSourceFileOfNode(declaration).FileName()) {
				// !!! diagnostics.You_cannot_rename_elements_that_are_defined_in_a_node_modules_folder
				return true
			}
		}
		return false
	}

	// Original source file is in node_modules.
	for _, declaration := range declarations {
		declPackage := getPackagePathComponents(ast.GetSourceFileOfNode(declaration).FileName())
		if declPackage != nil {
			length := min(len(originalPackage), len(declPackage))
			for i := 0; i <= length; i++ {
				var origComp, declComp string
				if i < len(originalPackage) {
					origComp = originalPackage[i]
				}
				if i < len(declPackage) {
					declComp = declPackage[i]
				}
				if origComp != declComp {
					// !!! diagnostics.You_cannot_rename_elements_that_are_defined_in_another_node_modules_folder
					return true
				}
			}
		}
	}
	return false
}

// getPackagePathComponents returns the path components up to and including the package name
// within node_modules, or nil if the path is not inside node_modules.
func getPackagePathComponents(filePath string) []string {
	components := tspath.GetPathComponents(filePath, "")
	nodeModulesIdx := -1
	for i := len(components) - 1; i >= 0; i-- {
		if components[i] == "node_modules" {
			nodeModulesIdx = i
			break
		}
	}
	if nodeModulesIdx == -1 {
		return nil
	}
	end := min(nodeModulesIdx+2, len(components))
	return components[:end]
}

// getRenameInfoForModule handles rename validation for module specifiers.
func (l *LanguageService) getRenameInfoForModule(node *ast.Node, sourceFile *ast.SourceFile, moduleSymbol *ast.Symbol) (RenameInfo, bool) {
	if !tspath.IsExternalModuleNameRelative(node.Text()) {
		// !!! diagnostics.You_cannot_rename_a_module_via_a_global_import
		return RenameInfo{}, true
	}

	moduleSourceFile := core.Find(moduleSymbol.Declarations, ast.IsSourceFile)
	if moduleSourceFile == nil {
		return RenameInfo{}, false
	}

	fileName := moduleSourceFile.AsSourceFile().FileName()
	withoutIndex := ""
	if !strings.HasSuffix(node.Text(), "/index") && !strings.HasSuffix(node.Text(), "/index.js") {
		candidate := tspath.RemoveFileExtension(fileName)
		if trimmed, ok := strings.CutSuffix(candidate, "/index"); ok {
			withoutIndex = trimmed
		}
	}

	displayName := fileName
	if withoutIndex != "" {
		displayName = withoutIndex
	}

	// Span should only be the last component of the path. + 1 to account for the quote character.
	indexAfterLastSlash := strings.LastIndex(node.Text(), "/") + 1
	start := node.Pos() + 1 + indexAfterLastSlash
	length := len(node.Text()) - indexAfterLastSlash

	return RenameInfo{
		CanRename:   true,
		DisplayName: displayName,
		TriggerSpan: l.converters.ToLSPRange(sourceFile, core.NewTextRange(start, start+length)),
	}, true
}

func (l *LanguageService) getTextForRename(originalNode *ast.Node, entry *ReferenceEntry, newText string, ch *checker.Checker) string {
	if entry.kind != entryKindRange && (ast.IsIdentifier(originalNode) || ast.IsStringLiteralLike(originalNode)) {
		node := ast.GetReparsedNodeForNode(entry.node)
		kind := entry.kind
		parent := node.Parent
		name := originalNode.Text()
		isShorthandAssignment := ast.IsShorthandPropertyAssignment(parent)
		switch {
		case isShorthandAssignment || (isObjectBindingElementWithoutPropertyName(parent) && parent.Name() == node && parent.AsBindingElement().DotDotDotToken == nil):
			if kind == entryKindSearchedLocalFoundProperty {
				return name + ": " + newText
			}
			if kind == entryKindSearchedPropertyFoundLocal {
				return newText + ": " + name
			}
			// In `const o = { x }; o.x`, symbolAtLocation at `x` in `{ x }` is the property symbol.
			// For a binding element `const { x } = o;`, symbolAtLocation at `x` is the property symbol.
			if isShorthandAssignment {
				grandParent := parent.Parent
				if ast.IsObjectLiteralExpression(grandParent) && ast.IsBinaryExpression(grandParent.Parent) && ast.IsModuleExportsAccessExpression(grandParent.Parent.AsBinaryExpression().Left) {
					return name + ": " + newText
				}
				return newText + ": " + name
			}
			return name + ": " + newText
		case ast.IsImportSpecifier(parent) && parent.PropertyName() == nil:
			// If the original symbol was using this alias, just rename the alias.
			var originalSymbol *ast.Symbol
			if ast.IsExportSpecifier(originalNode.Parent) {
				originalSymbol = ch.GetExportSpecifierLocalTargetSymbol(originalNode.Parent)
			} else {
				originalSymbol = ch.GetSymbolAtLocation(originalNode)
			}
			if slices.Contains(originalSymbol.Declarations, parent) {
				return name + " as " + newText
			}
			return newText
		case ast.IsExportSpecifier(parent) && parent.PropertyName() == nil:
			// If the symbol for the node is same as declared node symbol use prefix text
			if originalNode == entry.node || ch.GetSymbolAtLocation(originalNode) == ch.GetSymbolAtLocation(entry.node) {
				return name + " as " + newText
			}
			return newText + " as " + name
		}
	}
	return newText
}

func getRenameInfoSuccess(node *ast.Node, sourceFile *ast.SourceFile, displayName string, converters *lsconv.Converters) RenameInfo {
	start := node.Pos()
	end := node.End()
	if ast.IsStringLiteralLike(node) {
		// Exclude the quotes
		start++
		end--
	}
	return RenameInfo{
		CanRename:   true,
		DisplayName: displayName,
		TriggerSpan: converters.ToLSPRange(sourceFile, core.NewTextRange(start, end)),
	}
}
