package ls

import (
	"context"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/scanner"
)

func (l *LanguageService) ProvideCodeLenses(ctx context.Context, documentURI lsproto.DocumentUri) (lsproto.CodeLensResponse, error) {
	_, file := l.getProgramAndFile(documentURI)

	userPrefs := &l.UserPreferences().CodeLens
	if !userPrefs.ReferencesCodeLensEnabled && !userPrefs.ImplementationsCodeLensEnabled {
		return lsproto.CodeLensResponse{}, nil
	}

	// Keeps track of the last symbol to avoid duplicating code lenses across overloads.
	var lastSymbol *ast.Symbol
	var result []*lsproto.CodeLens
	var visit func(node *ast.Node) bool
	visit = func(node *ast.Node) bool {
		if ctx.Err() != nil {
			return true
		}

		if currentSymbol := node.Symbol(); lastSymbol != currentSymbol {
			lastSymbol = currentSymbol

			if userPrefs.ReferencesCodeLensEnabled && isValidReferenceLensNode(node, userPrefs) {
				result = append(result, l.newCodeLensForNode(documentURI, file, node, lsproto.CodeLensKindReferences))
			}

			if userPrefs.ImplementationsCodeLensEnabled && isValidImplementationsCodeLensNode(node, userPrefs) {
				result = append(result, l.newCodeLensForNode(documentURI, file, node, lsproto.CodeLensKindImplementations))
			}
		}

		savedLastSymbol := lastSymbol
		node.ForEachChild(visit)
		lastSymbol = savedLastSymbol
		return false
	}

	visit(file.AsNode())

	return lsproto.CodeLensResponse{
		CodeLenses: &result,
	}, nil
}

func (l *LanguageService) newCodeLensForNode(fileUri lsproto.DocumentUri, file *ast.SourceFile, node *ast.Node, kind lsproto.CodeLensKind) *lsproto.CodeLens {
	nodeForRange := node
	nodeName := node.Name()
	if nodeName != nil {
		nodeForRange = nodeName
	}
	pos := scanner.SkipTrivia(file.Text(), nodeForRange.Pos())

	return &lsproto.CodeLens{
		Range: lsproto.Range{
			Start: l.converters.PositionToLineAndCharacter(file, core.TextPos(pos)),
			End:   l.converters.PositionToLineAndCharacter(file, core.TextPos(node.End())),
		},
		Data: &lsproto.CodeLensData{
			Kind: kind,
			Uri:  fileUri,
		},
	}
}

func isValidImplementationsCodeLensNode(node *ast.Node, userPrefs *lsutil.CodeLensUserPreferences) bool {
	switch node.Kind {
	// Always show on interfaces
	case ast.KindInterfaceDeclaration:
		// TODO: ast.KindTypeAliasDeclaration?
		return true

	// If configured, show on interface methods
	case ast.KindMethodSignature:
		return userPrefs.ImplementationsCodeLensShowOnInterfaceMethods && node.Parent.Kind == ast.KindInterfaceDeclaration

	// If configured, show on all class methods - but not private ones.
	case ast.KindMethodDeclaration:
		if userPrefs.ImplementationsCodeLensShowOnAllClassMethods && node.Parent.Kind == ast.KindClassDeclaration {
			return !ast.HasModifier(node, ast.ModifierFlagsPrivate) && node.Name().Kind != ast.KindPrivateIdentifier
		}
		fallthrough

	// Always show on abstract classes/properties/methods
	case ast.KindClassDeclaration, ast.KindConstructor,
		ast.KindGetAccessor, ast.KindSetAccessor, ast.KindPropertyDeclaration:
		return ast.HasModifier(node, ast.ModifierFlagsAbstract)
	}

	return false
}

func isValidReferenceLensNode(node *ast.Node, userPrefs *lsutil.CodeLensUserPreferences) bool {
	switch node.Kind {
	case ast.KindFunctionDeclaration:
		if userPrefs.ReferencesCodeLensShowOnAllFunctions {
			return true
		}
		fallthrough

	case ast.KindVariableDeclaration:
		return ast.GetCombinedModifierFlags(node)&ast.ModifierFlagsExport != 0

	case ast.KindClassDeclaration, ast.KindInterfaceDeclaration, ast.KindTypeAliasDeclaration, ast.KindEnumDeclaration, ast.KindEnumMember:
		return true

	case ast.KindMethodDeclaration, ast.KindMethodSignature, ast.KindConstructor,
		ast.KindGetAccessor, ast.KindSetAccessor,
		ast.KindPropertyDeclaration, ast.KindPropertySignature:
		// Don't show if child and parent have same start
		// For https://github.com/microsoft/vscode/issues/90396
		// !!!

		switch node.Parent.Kind {
		case ast.KindClassDeclaration, ast.KindInterfaceDeclaration, ast.KindTypeLiteral:
			return true
		}
	}

	return false
}
