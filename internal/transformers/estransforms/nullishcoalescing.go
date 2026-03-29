package estransforms

import (
	"sync"
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/transformers"
)

type nullishCoalescingTransformer struct {
	transformers.Transformer
}

func (ch *nullishCoalescingTransformer) visit(node *ast.Node) *ast.Node {
	if node.SubtreeFacts()&ast.SubtreeContainsNullishCoalescing == 0 {
		return node
	}
	switch node.Kind {
	case ast.KindBinaryExpression:
		return ch.visitBinaryExpression(node.AsBinaryExpression())
	default:
		return ch.Visitor().VisitEachChild(node)
	}
}

func (ch *nullishCoalescingTransformer) visitBinaryExpression(node *ast.BinaryExpression) *ast.Node {
	switch node.OperatorToken.Kind {
	case ast.KindQuestionQuestionToken:
		left := ch.Visitor().VisitNode(node.Left)
		right := left
		if !transformers.IsSimpleCopiableExpression(left) {
			right = ch.Factory().NewTempVariable()
			ch.EmitContext().AddVariableDeclaration(right)
			left = ch.Factory().NewAssignmentExpression(right, left)
		}
		return ch.Factory().NewConditionalExpression(
			createNotNullCondition(ch.EmitContext(), left, right, false),
			ch.Factory().NewToken(ast.KindQuestionToken),
			right,
			ch.Factory().NewToken(ast.KindColonToken),
			ch.Visitor().VisitNode(node.Right),
		)
	default:
		return ch.Visitor().VisitEachChild(node.AsNode())
	}
}

var nullishCoalescingTransformerPool = sync.Pool{New: func() any { return &nullishCoalescingTransformer{} }}

func getNullishCoalescingTransformer() *nullishCoalescingTransformer {
return nullishCoalescingTransformerPool.Get().(*nullishCoalescingTransformer)
}

func putNullishCoalescingTransformer(tx *nullishCoalescingTransformer) {
dispose, visit := tx.SaveState()
*tx = nullishCoalescingTransformer{}
tx.RestoreState(dispose, visit)
nullishCoalescingTransformerPool.Put(tx)
}

func newNullishCoalescingTransformer(opts *transformers.TransformOptions) *transformers.Transformer {
	tx := getNullishCoalescingTransformer()
	if tx.GetDispose() == nil {
		tx.SetDispose(func() { putNullishCoalescingTransformer(tx) })
	}
	return tx.NewTransformer(tx.visit, opts.Context)
}
