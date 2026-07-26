package ls

import (
	"context"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/locale"
	"github.com/microsoft/typescript-go/internal/ls/change"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/nodebuilder"
)

const inferReturnTypeRefactorKind = lsproto.CodeActionKind("refactor.rewrite.function.returnType")

// InferReturnTypeProvider is a RefactorProvider that adds return type annotations
// by inferring the type from the function body.
var InferReturnTypeProvider = &RefactorProvider{
	RefactorActions: []RefactorAction{
		{
			Title:      "Infer function return type",
			ID:         "inferReturnType",
			Kinds:      []lsproto.CodeActionKind{inferReturnTypeRefactorKind},
			GetActions: getInferReturnTypeCodeActions,
		},
	},
}

// RefactorContext contains the context needed to generate refactoring actions.
type RefactorContext struct {
	SourceFile *ast.SourceFile
	Range      core.TextRange
	Program    *compiler.Program
	LS         *LanguageService
	Params     *lsproto.CodeActionParams
}

// convertibleDeclaration returns true if the node is a function-like declaration
// that can have its return type inferred (matches TypeScript's ConvertibleDeclaration).
func convertibleDeclaration(node *ast.Node) bool {
	return ast.IsFunctionLikeDeclaration(node) &&
		!ast.IsConstructorDeclaration(node) &&
		!ast.IsGetAccessorDeclaration(node) &&
		!ast.IsSetAccessorDeclaration(node)
}

// getInferReturnTypeCodeActions provides the "Infer Return Type" refactoring action.
func getInferReturnTypeCodeActions(ctx context.Context, refactorContext *RefactorContext, refactorID string) ([]*CodeAction, error) {
	if ast.IsInJSFile(refactorContext.SourceFile.AsNode()) {
		return nil, nil
	}

	token := astnav.GetTouchingPropertyName(refactorContext.SourceFile, int(refactorContext.Range.Pos()))

	declaration := findConvertibleAncestor(token)
	if declaration == nil || !hasBody(declaration) || declaration.Type() != nil {
		return nil, nil
	}

	ch, done := refactorContext.Program.GetTypeCheckerForFile(ctx, refactorContext.SourceFile)
	defer done()

	typeNode := getInferredReturnTypeNode(ch, declaration, refactorContext.SourceFile)
	if typeNode == nil {
		return nil, nil
	}

	formatOptions := refactorContext.LS.FormatOptions()
	changeTracker := change.NewTracker(ctx, refactorContext.Program.Options(), formatOptions, refactorContext.LS.converters)

	if ast.IsArrowFunction(declaration) {
		changeTracker.ParenthesizeArrowParameters(refactorContext.SourceFile, declaration)
	}
	changeTracker.TryInsertTypeAnnotation(refactorContext.SourceFile, declaration, typeNode)

	changes := changeTracker.GetChanges()
	if len(changes) == 0 {
		return nil, nil
	}

	title := diagnostics.Infer_function_return_type.Localize(locale.FromContext(ctx))

	actions := []*CodeAction{{
		Description: title,
		Changes:     changes[refactorContext.SourceFile.FileName()],
		FixID:       refactorID,
		Kind:        inferReturnTypeRefactorKind,
	}}
	return actions, nil
}

// getInferredReturnTypeNode analyzes a function-like declaration and returns its
// inferred return type node, handling overloads and type predicates.
func getInferredReturnTypeNode(ch *checker.Checker, declaration *ast.Node, sourceFile *ast.SourceFile) *ast.Node {
	idToSymbol := make(map[*ast.IdentifierNode]*ast.Symbol)

	// Handle overloaded function implementations: union return types of all signatures.
	if ch.GetEmitResolver().IsImplementationOfOverload(declaration) {
		fnType := ch.GetTypeAtLocation(declaration)
		if fnType != nil {
			signatures := ch.GetCallSignatures(fnType)
			if len(signatures) > 1 {
				returnTypes := make([]*checker.Type, 0, len(signatures))
				for _, sig := range signatures {
					rt := ch.GetReturnTypeOfSignature(sig)
					if rt != nil {
						returnTypes = append(returnTypes, rt)
					}
				}

				if len(returnTypes) > 0 {
					return ch.TypeToTypeNodeEx(ch.GetUnionType(returnTypes), declaration, nodebuilder.FlagsNoTruncation, nodebuilder.InternalFlagsAllowUnresolvedNames, idToSymbol)
				}
			}
		}
	}

	// Check for type predicate (e.g., `x is T`): build the predicate node directly.
	signature := ch.GetSignatureFromDeclaration(declaration)
	if signature != nil {
		typePredicate := ch.GetTypePredicateOfSignature(signature)
		if typePredicate != nil && typePredicate.Type() != nil {
			enclosingDecl := ast.FindAncestor(declaration, ast.IsDeclaration)
			if enclosingDecl == nil {
				enclosingDecl = sourceFile.AsNode()
			}
			return ch.TypePredicateToTypePredicateNodeEx(typePredicate, enclosingDecl, nodebuilder.FlagsNoTruncation, nodebuilder.InternalFlagsAllowUnresolvedNames, idToSymbol)
		}
	}

	// Normal case: get the return type of the signature.
	if signature != nil {
		return ch.TypeToTypeNodeEx(ch.GetReturnTypeOfSignature(signature), declaration, nodebuilder.FlagsNoTruncation, nodebuilder.InternalFlagsAllowUnresolvedNames, idToSymbol)
	}

	return nil
}

// findConvertibleAncestor walks up from node to find the nearest convertible declaration.
func findConvertibleAncestor(node *ast.Node) *ast.Node {
	for node != nil {
		if ast.IsBlock(node) {
			return nil
		}

		if node.Parent != nil && ast.IsArrowFunction(node.Parent) {
			if node.Kind == ast.KindEqualsGreaterThanToken || node.Parent.AsArrowFunction().Body == node {
				return nil
			}
		}

		if convertibleDeclaration(node) {
			return node
		}

		node = node.Parent
	}

	return nil
}

// hasBody returns true if the node has a function body.
func hasBody(node *ast.Node) bool {
	switch node.Kind {
	case ast.KindFunctionDeclaration:
		return node.AsFunctionDeclaration().Body != nil
	case ast.KindFunctionExpression:
		return node.AsFunctionExpression().Body != nil
	case ast.KindArrowFunction:
		return node.AsArrowFunction().Body != nil
	case ast.KindMethodDeclaration:
		return node.AsMethodDeclaration().Body != nil
	}
	return false
}
