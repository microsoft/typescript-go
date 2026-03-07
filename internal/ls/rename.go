package ls

import (
	"context"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func (l *LanguageService) ProvidePrepareRename(ctx context.Context, documentURI lsproto.DocumentUri, position lsproto.Position) (lsproto.PrepareRenameResponse, error) {
	program, sourceFile := l.getProgramAndFile(documentURI)
	pos := int(l.converters.LineAndCharacterToPosition(sourceFile, position))

	node := getAdjustedRenameLocation(astnav.GetTouchingPropertyName(sourceFile, pos))
	if !isNodeEligibleForRename(node) {
		return lsproto.PrepareRenameResponse{}, nil
	}

	renameInfo := l.getRenameInfo(ctx, node, sourceFile, program)
	if renameInfo == nil {
		return lsproto.PrepareRenameResponse{}, nil
	}

	return *renameInfo, nil
}

func getAdjustedRenameLocation(node *ast.Node) *ast.Node {
	return getAdjustedLocation(node, true /*forRename*/, nil)
}

func (l *LanguageService) getRenameInfo(ctx context.Context, node *ast.Node, sourceFile *ast.SourceFile, program *compiler.Program) *lsproto.PrepareRenameResponse {
	checker, done := program.GetTypeChecker(ctx)
	defer done()

	symbol := checker.GetSymbolAtLocation(node)
	if symbol == nil {
		if ast.IsStringLiteralLike(node) {
			return l.prepareRenameResponse(node, sourceFile)
		}
		if isLabelName(node) {
			return l.prepareRenameResponse(node, sourceFile)
		}
		return nil
	}

	// Only allow a symbol to be renamed if it actually has at least one declaration.
	declarations := symbol.Declarations
	if len(declarations) == 0 {
		return nil
	}

	// Disallow rename for elements that are defined in the standard TypeScript library.
	if isDefinedInLibraryFile(program, declarations) {
		return nil
	}

	// Cannot rename `default` as in `import { default as foo } from "./someModule";
	if ast.IsIdentifier(node) && node.Text() == "default" && symbol.Parent != nil && symbol.Parent.Flags&ast.SymbolFlagsModule != 0 {
		return nil
	}

	if ast.IsStringLiteralLike(node) && ast.TryGetImportFromModuleSpecifier(node) != nil {
		if l.UserPreferences().AllowRenameOfImportPath != core.TSTrue {
			return nil
		}
	}

	// Disallow rename for elements that would rename across node_modules packages.
	if wouldRenameInOtherNodeModules(sourceFile, symbol) {
		return nil
	}

	return l.prepareRenameResponse(node, sourceFile)
}

func (l *LanguageService) prepareRenameResponse(node *ast.Node, sourceFile *ast.SourceFile) *lsproto.PrepareRenameResponse {
	start := astnav.GetStartOfNode(node, sourceFile, false /*includeJSDoc*/)
	end := node.End()

	if ast.IsStringLiteralLike(node) {
		// Exclude the quotes
		start += 1
		end -= 1
	}

	startPos := l.converters.PositionToLineAndCharacter(sourceFile, core.TextPos(start))
	endPos := l.converters.PositionToLineAndCharacter(sourceFile, core.TextPos(end))

	text := sourceFile.Text()[start:end]
	return &lsproto.PrepareRenameResponse{
		PrepareRenamePlaceholder: &lsproto.PrepareRenamePlaceholder{
			Range: lsproto.Range{
				Start: startPos,
				End:   endPos,
			},
			Placeholder: text,
		},
	}
}

func isDefinedInLibraryFile(program *compiler.Program, declarations []*ast.Node) bool {
	for _, declaration := range declarations {
		sourceFile := ast.GetSourceFileOfNode(declaration)
		if program.IsSourceFileDefaultLibrary(sourceFile.Path()) && tspath.FileExtensionIs(sourceFile.FileName(), tspath.ExtensionDts) {
			return true
		}
	}
	return false
}

func wouldRenameInOtherNodeModules(originalFile *ast.SourceFile, symbol *ast.Symbol) bool {
	declarations := symbol.Declarations
	if len(declarations) == 0 {
		return false
	}

	originalIsInNodeModules := strings.Contains(string(originalFile.Path()), "/node_modules/")
	if !originalIsInNodeModules {
		// If the original file is not in node_modules, disallow rename if any declaration is in node_modules
		for _, declaration := range declarations {
			if strings.Contains(string(ast.GetSourceFileOfNode(declaration).Path()), "/node_modules/") {
				return true
			}
		}
		return false
	}
	return false
}

func isLabelName(node *ast.Node) bool {
	return isDeclarationName(node) && ast.IsLabeledStatement(node.Parent)
}

func isDeclarationName(node *ast.Node) bool {
	return !ast.IsSourceFile(node) && node.Parent != nil && ast.GetNameOfDeclaration(node.Parent) == node
}
