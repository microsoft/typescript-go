package ls

import (
	"context"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/scanner"
)

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
			lenses = append(lenses, l.newCodeLensForNode(documentURI, file, node, lsproto.CodeLensKindReferences))
		}

		if userPrefs.ImplementationsCodeLensEnabled && isValidImplementationsCodeLensNode(node, userPrefs) {
			lenses = append(lenses, l.newCodeLensForNode(documentURI, file, node, lsproto.CodeLensKindImplementations))
		}

		node.ForEachChild(visit)
		return false
	}

	visit(file.AsNode())

	return lsproto.CodeLensResponse{
		CodeLenses: &lenses,
	}, nil
}

func (l *LanguageService) ResolveCodeLens(ctx context.Context, codeLens *lsproto.CodeLens, showLocationsCommandName *string) (*lsproto.CodeLens, error) {
	uri := codeLens.Data.Uri
	textDoc := lsproto.TextDocumentIdentifier{
		Uri: uri,
	}

	var locs []lsproto.Location
	var lensTitle string
	switch codeLens.Data.Kind {
	case lsproto.CodeLensKindReferences:
		origNode, symbolsAndEntries, ok := l.ProvideSymbolsAndEntries(ctx, uri, codeLens.Range.Start, false /*isRename*/)
		if ok {
			references, err := l.ProvideReferencesFromSymbolAndEntries(
				ctx,
				&lsproto.ReferenceParams{
					TextDocument: textDoc,
					Position:     codeLens.Range.Start,
					Context: &lsproto.ReferenceContext{
						// Don't include the declaration in the references count.
						IncludeDeclaration: false,
					},
				},
				origNode,
				symbolsAndEntries,
			)
			if err != nil {
				return nil, err
			}

			if references.Locations != nil {
				locs = *references.Locations
			}
		}

		if len(locs) == 1 {
			lensTitle = diagnostics.X_1_reference.Message()
		} else {
			lensTitle = diagnostics.X_0_references.Format(len(locs))
		}
	case lsproto.CodeLensKindImplementations:
		// "Force" link support to be false so that we only get `Locations` back,
		// and don't include the "current" node in the results.
		findImplsOptions := provideImplementationsOpts{
			requireLocationsResult: true,
			dropOriginNodes:        true,
		}
		implementations, err := l.provideImplementationsEx(
			ctx,
			&lsproto.ImplementationParams{
				TextDocument: textDoc,
				Position:     codeLens.Range.Start,
			},
			findImplsOptions,
		)
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
	if len(locs) > 0 && showLocationsCommandName != nil {
		cmd.Command = *showLocationsCommandName
		cmd.Arguments = &[]any{
			uri,
			codeLens.Range.Start,
			locs,
		}
	}

	codeLens.Command = cmd
	return codeLens, nil
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
