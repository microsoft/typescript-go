package estransforms

import (
	"github.com/microsoft/typescript-go/internal/ast"
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

// For example, += -> +
func getNonAssignmentOperatorForCompoundAssignment(emitContext *printer.EmitContext, tokenNode *ast.TokenNode) *ast.TokenNode {
	switch tokenNode.Kind {
	case ast.KindPlusEqualsToken:
		return emitContext.Factory.NewToken(ast.KindPlusToken)
	case ast.KindMinusEqualsToken:
		return emitContext.Factory.NewToken(ast.KindMinusToken)
	case ast.KindAsteriskEqualsToken:
		return emitContext.Factory.NewToken(ast.KindAsteriskToken)
	case ast.KindAsteriskAsteriskEqualsToken:
		return emitContext.Factory.NewToken(ast.KindAsteriskAsteriskToken)
	case ast.KindSlashEqualsToken:
		return emitContext.Factory.NewToken(ast.KindSlashToken)
	case ast.KindPercentEqualsToken:
		return emitContext.Factory.NewToken(ast.KindPercentToken)
	case ast.KindLessThanLessThanEqualsToken:
		return emitContext.Factory.NewToken(ast.KindLessThanLessThanToken)
	case ast.KindGreaterThanGreaterThanEqualsToken:
		return emitContext.Factory.NewToken(ast.KindGreaterThanGreaterThanToken)
	case ast.KindGreaterThanGreaterThanGreaterThanEqualsToken:
		return emitContext.Factory.NewToken(ast.KindGreaterThanGreaterThanGreaterThanToken)
	case ast.KindAmpersandEqualsToken:
		return emitContext.Factory.NewToken(ast.KindAmpersandToken)
	case ast.KindBarEqualsToken:
		return emitContext.Factory.NewToken(ast.KindBarToken)
	case ast.KindCaretEqualsToken:
		return emitContext.Factory.NewToken(ast.KindCaretToken)
	case ast.KindBarBarEqualsToken:
		return emitContext.Factory.NewToken(ast.KindBarBarToken)
	case ast.KindAmpersandAmpersandEqualsToken:
		return emitContext.Factory.NewToken(ast.KindAmpersandAmpersandToken)
	case ast.KindQuestionQuestionEqualsToken:
		return emitContext.Factory.NewToken(ast.KindQuestionQuestionToken)
	}
	return nil
}
