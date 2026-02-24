package estransforms

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/transformers"
)

func convertClassDeclarationToClassExpression(emitContext *printer.EmitContext, node *ast.ClassDeclaration) *ast.Expression {
	updated := emitContext.Factory.NewClassExpression(
		transformers.ExtractModifiers(emitContext, node.Modifiers(), ^ast.ModifierFlagsExportDefault),
		node.Name(),
		node.TypeParameters,
		node.HeritageClauses,
		node.Members,
	)
	emitContext.SetOriginal(updated, node.AsNode())
	updated.Loc = node.Loc
	return updated
}

func createNotNullCondition(emitContext *printer.EmitContext, left *ast.Node, right *ast.Node, invert bool) *ast.Node {
	token := ast.KindExclamationEqualsEqualsToken
	op := ast.KindAmpersandAmpersandToken
	if invert {
		token = ast.KindEqualsEqualsEqualsToken
		op = ast.KindBarBarToken
	}

	return emitContext.Factory.NewBinaryExpression(
		nil,
		emitContext.Factory.NewBinaryExpression(
			nil,
			left,
			nil,
			emitContext.Factory.NewToken(token),
			emitContext.Factory.NewKeywordExpression(ast.KindNullKeyword),
		),
		nil,
		emitContext.Factory.NewToken(op),
		emitContext.Factory.NewBinaryExpression(
			nil,
			right,
			nil,
			emitContext.Factory.NewToken(token),
			emitContext.Factory.NewVoidZeroExpression(),
		),
	)
}

// superAccessState tracks super property/element accesses and super property assignments
// within async function or async generator bodies. It is embedded by both asyncTransformer
// and forawaitTransformer to share the tracking logic.
type superAccessState struct {
	// Keeps track of property names accessed on super (`super.x`) within async functions.
	capturedSuperProperties *collections.OrderedSet[string]
	// Whether the async function contains an element access on super (`super[x]`).
	hasSuperElementAccess      bool
	hasSuperPropertyAssignment bool
}

// trackSuperAccess records super property/element accesses and super property assignments
// for the enclosing async method body. Called from both the main visitor and auxiliary
// visitors to ensure super accesses are tracked regardless of whether the node has
// transform flags.
func (s *superAccessState) trackSuperAccess(node *ast.Node) {
	if s.capturedSuperProperties == nil {
		return
	}
	switch node.Kind {
	case ast.KindPropertyAccessExpression:
		if node.Expression().Kind == ast.KindSuperKeyword {
			s.capturedSuperProperties.Add(node.Name().Text())
		}
	case ast.KindElementAccessExpression:
		if node.Expression().Kind == ast.KindSuperKeyword {
			s.hasSuperElementAccess = true
		}
	case ast.KindBinaryExpression:
		if ast.IsAssignmentOperator(node.AsBinaryExpression().OperatorToken.Kind) && assignmentTargetContainsSuperProperty(node.AsBinaryExpression().Left) {
			s.hasSuperPropertyAssignment = true
		}
	case ast.KindPrefixUnaryExpression:
		if isUpdateExpression(node) && assignmentTargetContainsSuperProperty(node.AsPrefixUnaryExpression().Operand) {
			s.hasSuperPropertyAssignment = true
		}
	case ast.KindPostfixUnaryExpression:
		if isUpdateExpression(node) && assignmentTargetContainsSuperProperty(node.AsPostfixUnaryExpression().Operand) {
			s.hasSuperPropertyAssignment = true
		}
	}
}
