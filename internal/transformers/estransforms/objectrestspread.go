package estransforms

import (
	"context"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/transformers"
)

type objectRestSpreadTransformer struct {
	transformers.Transformer
	compilerOptions *core.CompilerOptions
}

func (ch *objectRestSpreadTransformer) visit(node *ast.Node) *ast.Node {
	if node.SubtreeFacts()&ast.SubtreeContainsESObjectRestOrSpread == 0 {
		return node
	}
	switch node.Kind {
	case ast.KindObjectLiteralExpression:
		return ch.visitObjectLiteralExpression(node.AsObjectLiteralExpression())
	case ast.KindBinaryExpression:
		return ch.visitBinaryExpression(node.AsBinaryExpression())
	// case ast.KindForOfStatement:
	// 	return ch.visitForOftatement(node.AsForInOrOfStatement())
	// case ast.KindVariableStatement:
	// 	return ch.visitVariableStatement(node.AsVariableStatement())
	default:
		return ch.Visitor().VisitEachChild(node)
	}
}

func (ch *objectRestSpreadTransformer) visitBinaryExpression(node *ast.BinaryExpression) *ast.Node {
	if !ast.IsDestructuringAssignment(node.AsNode()) || (node.Left.SubtreeFacts()&ast.SubtreeContainsObjectRestOrSpread) == 0 {
		return ch.Visitor().VisitEachChild(node.AsNode())
	}
	// !!!
	return node.AsNode()
}

func (ch *objectRestSpreadTransformer) visitObjectLiteralExpression(node *ast.ObjectLiteralExpression) *ast.Node {
	if (node.SubtreeFacts() & ast.SubtreeContainsObjectRestOrSpread) == 0 {
		return ch.Visitor().VisitEachChild(node.AsNode())
	}
	// spread elements emit like so:
	// non-spread elements are chunked together into object literals, and then all are passed to __assign:
	//     { a, ...o, b } => __assign(__assign({a}, o), {b});
	// If the first element is a spread element, then the first argument to __assign is {}:
	//     { ...o, a, b, ...o2 } => __assign(__assign(__assign({}, o), {a, b}), o2)
	//
	// We cannot call __assign with more than two elements, since any element could cause side effects. For
	// example:
	//      var k = { a: 1, b: 2 };
	//      var o = { a: 3, ...k, b: k.a++ };
	//      // expected: { a: 1, b: 1 }
	// If we translate the above to `__assign({ a: 3 }, k, { b: k.a++ })`, the `k.a++` will evaluate before
	// `k` is spread and we end up with `{ a: 2, b: 1 }`.
	//
	// This also occurs for spread elements, not just property assignments:
	//      var k = { a: 1, get b() { l = { z: 9 }; return 2; } };
	//      var l = { c: 3 };
	//      var o = { ...k, ...l };
	//      // expected: { a: 1, b: 2, z: 9 }
	// If we translate the above to `__assign({}, k, l)`, the `l` will evaluate before `k` is spread and we
	// end up with `{ a: 1, b: 2, c: 3 }`

	objects := ch.chunkObjectLiteralElements(node.Properties)
	if len(objects) > 0 && objects[0].Kind != ast.KindObjectLiteralExpression {
		objects = append([]*ast.Node{ch.Factory().NewObjectLiteralExpression(ch.Factory().NewNodeList(nil), false)}, objects...)
	}
	expression := objects[0]
	if len(objects) > 0 {
		for i, obj := range objects {
			if i == 0 {
				continue
			}
			expression = ch.Factory().NewAssignHelper([]*ast.Node{expression, obj}, ch.compilerOptions.GetEmitScriptTarget())
		}
		return expression
	}
	return ch.Factory().NewAssignHelper(objects, ch.compilerOptions.GetEmitScriptTarget())
}

func (ch *objectRestSpreadTransformer) chunkObjectLiteralElements(list *ast.NodeList) []*ast.Node {
	if list == nil || len(list.Nodes) == 0 {
		return nil
	}
	elements := list.Nodes
	var chunkObject []*ast.Node
	objects := make([]*ast.Node, 0, 1)
	for _, e := range elements {
		if e.Kind == ast.KindSpreadAssignment {
			if len(chunkObject) > 0 {
				objects = append(objects, ch.Factory().NewObjectLiteralExpression(ch.Factory().NewNodeList(chunkObject), false))
				chunkObject = nil
			}
			target := e.Expression()
			objects = append(objects, ch.Visitor().Visit(target))
		} else {
			var elem *ast.Node
			if e.Kind == ast.KindPropertyAssignment {
				elem = ch.Factory().NewPropertyAssignment(nil, e.Name(), nil, nil, ch.Visitor().Visit(e.Initializer()))
			} else {
				elem = ch.Visitor().Visit(e)
			}
			chunkObject = append(chunkObject, elem)
		}
	}
	if len(chunkObject) > 0 {
		objects = append(objects, ch.Factory().NewObjectLiteralExpression(ch.Factory().NewNodeList(chunkObject), false))
	}
	return objects
}

func newObjectRestSpreadTransformer(ctx context.Context) *transformers.Transformer {
	tx := &objectRestSpreadTransformer{compilerOptions: transformers.GetCompilerOptionsFromContext(ctx)}
	return tx.NewTransformer(tx.visit, transformers.GetEmitContextFromContext(ctx))
}
