package ls

import (
	"context"
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/scanner"
)

type CodeLensKind string

const (
	codeLensReferencesKind      CodeLensKind = "references"
	codeLensImplementationsKind CodeLensKind = "implementations"
)

type CodeLensData struct {
	Kind CodeLensKind        `json:"kind"`
	Uri  lsproto.DocumentUri `json:"uri"`
}

func GetCodeLensData(item *lsproto.CodeLens) (*CodeLensData, error) {
	bytes, err := json.Marshal(item.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal completion item data: %w", err)
	}
	var itemData CodeLensData
	if err := json.Unmarshal(bytes, &itemData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal completion item data: %w", err)
	}
	return &itemData, nil
}

func (l *LanguageService) ProvideCodeLenses(ctx context.Context, documentURI lsproto.DocumentUri) (lsproto.CodeLensResponse, error) {
	_, file := l.getProgramAndFile(documentURI)

	userPrefs := l.UserPreferences()
	if !userPrefs.ReferencesCodeLensEnabled && !userPrefs.ImplementationsCodeLensEnabled {
		return lsproto.CodeLensResponse{}, nil
	}

	var lenses []*lsproto.CodeLens
	var visit func(node *ast.Node) bool
	visit = func(node *ast.Node) bool {
		if ctx.Err() != nil {
			return true
		}

		if userPrefs.ReferencesCodeLensEnabled && isValidReferenceLensNode(node, userPrefs) {
			lenses = append(lenses, l.newCodeLensForNode(documentURI, file, node, codeLensReferencesKind))
		}

		if userPrefs.ImplementationsCodeLensEnabled && isValidImplementationsCodeLensNode(node, userPrefs) {
			lenses = append(lenses, l.newCodeLensForNode(documentURI, file, node, codeLensImplementationsKind))
		}

		node.ForEachChild(visit)
		return false
	}

	visit(file.AsNode())

	return lsproto.CodeLensResponse{
		CodeLenss: &lenses,
	}, nil
}

func (l *LanguageService) ResolveCodeLens(ctx context.Context, codeLens *lsproto.CodeLens, codeLensData *CodeLensData) (*lsproto.CodeLens, error) {
	uri := codeLensData.Uri
	textDoc := lsproto.TextDocumentIdentifier{
		Uri: uri,
	}

	var locs []lsproto.Location
	var lensTitle string
	switch codeLensData.Kind {
	case codeLensReferencesKind:
		references, err := l.ProvideReferences(ctx, &lsproto.ReferenceParams{
			TextDocument: textDoc,
			Position:     codeLens.Range.Start,
			Context: &lsproto.ReferenceContext{
				// Don't include the declaration in the references count.
				IncludeDeclaration: false,
			},
		})
		if err != nil {
			return nil, err
		}

		if references.Locations != nil {
			locs = *references.Locations
		}

		if len(locs) == 1 {
			lensTitle = diagnostics.X_1_reference.Message()
		} else {
			lensTitle = diagnostics.X_0_references.Format(len(locs))
		}
	case codeLensImplementationsKind:
		// "Force" link support to be false so that we only get `Locations` back,
		// and don't include the "current" node in the results.
		findImplsOptions := provideImplementationsOpts{
			requireLocationsResult: true,
			dropOriginNodes:        true,
		}
		implementations, err := l.provideImplementationsEx(ctx, &lsproto.ImplementationParams{
			TextDocument: textDoc,
			Position:     codeLens.Range.Start,
		}, findImplsOptions)
		if err != nil {
			return nil, err
		}

		if implementations.Locations != nil {
			locs = *implementations.Locations
		}

		if len(locs) == 1 {
			lensTitle = diagnostics.X_1_implementation.Message()
		} else {
			lensTitle = diagnostics.X_0_implementations.Format(len(locs))
		}
	}

	cmd := &lsproto.Command{
		Title: lensTitle,
	}
	if len(locs) > 0 {
		cmd.Command = "typescript.codeLens.showLocations"
		cmd.Arguments = &[]any{
			uri,
			codeLens.Range.Start,
			locs,
		}
	}

	codeLens.Command = cmd
	return codeLens, nil
}

func (l *LanguageService) newCodeLensForNode(fileUri lsproto.DocumentUri, file *ast.SourceFile, node *ast.Node, kind CodeLensKind) *lsproto.CodeLens {
	nodeForRange := node
	nodeName := node.Name()
	if nodeName != nil {
		nodeForRange = nodeName
	}
	pos := scanner.SkipTrivia(file.Text(), nodeForRange.Pos())

	var data any = &CodeLensData{
		Kind: kind,
		Uri:  fileUri,
	}

	return &lsproto.CodeLens{
		Range: lsproto.Range{
			Start: l.converters.PositionToLineAndCharacter(file, core.TextPos(pos)),
			End:   l.converters.PositionToLineAndCharacter(file, core.TextPos(node.End())),
		},
		Data: &data,
	}
}

func isValidImplementationsCodeLensNode(node *ast.Node, userPrefs *lsutil.UserPreferences) bool {
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

func isValidReferenceLensNode(node *ast.Node, userPrefs *lsutil.UserPreferences) bool {
	switch node.Kind {
	case ast.KindFunctionDeclaration, ast.KindFunctionExpression:
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
