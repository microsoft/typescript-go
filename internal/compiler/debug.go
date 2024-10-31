package compiler

import (
	"fmt"
	"reflect"
	"runtime"
	"slices"
	"strings"
)

// TODO(rbuckton): If we want some of these assertions to only be part of a checked/debug build, we could
// use build tags (i.e., `go build --tags DEBUG ...`) and split out debug-only assertions into two files:
//
// debug_checked.go:
//
//		// +BUILD DEBUG
//		func assertFoo() { /* implementation*/ }
//
// debug_release.go:
//
//		// +BUILD !DEBUG
//		func assertFoo() {} // no implementation

// Asserts 'cond' is true or panics with a formatted message
func assert(cond bool, a ...any) {
	if !cond {
		if len(a) > 0 {
			panic(fmt.Sprintf(fmt.Sprint(a[0]), a[1:]...))
		} else {
			panic("assert failed")
		}
	}
}

// Panics with a formatted message
func assertFail(a ...any) {
	assert(false, a...)
}

// Asserts 'x' is nil
func assertMissing(x any) {
	assert(x != nil, "Expected node to be missing")
}

// Asserts 'kind' is one of the expected values
//
//	assertKind(kind, SyntaxKindQuestionToken, SyntaxKindExclamationToken)
//	assertKind(kind, SyntaxKindCaseClause, SyntaxKindDefaultClause)
func assertKind(kind SyntaxKind, expected ...SyntaxKind) {
	assert(slices.Contains(expected, kind), "kind '%v' was not one of the expected values: %v", kind, expected)
}

// Asserts 'node' is not nil and is a Token with one of the expected 'kinds'
//
//	assertToken(asteriskToken)
//	assertToken(postfixToken, SyntaxKindQuestionToken, SyntaxKindExclamationToken)
func assertToken(node *Node, kinds ...SyntaxKind) {
	assert(node != nil, "Expected node to be present")
	assert(node.kind >= SyntaxKindFirstToken && node.kind <= SyntaxKindLastToken, "Node of %v was not a token", node.kind)
	assert(len(kinds) == 0 || nodeKindIs(node, kinds...), "Token '%v' did not have one of the expected kinds: %v", kinds)
}

// Asserts 'node' is either nil or a Token with one of the expected kinds
//
//	assertTokenOpt(asteriskToken)
//	assertTokenOpt(postfixToken, SyntaxKindQuestionToken, SyntaxKindExclamationToken)
func assertTokenOpt(node *Node, kinds ...SyntaxKind) {
	if node != nil {
		assertToken(node, kinds...)
	}
}

// Asserts 'node' is not nil and matches one of the provided node tests
//
//	assertNode(node, isExpression)
//	assertNode(node, isIdentifierName, isStringLiteral)
func assertNode(node *Node, fns ...NodeTest) {
	assert(node != nil, "Expected node to be present")
	if len(fns) > 0 {
		for _, fn := range fns {
			if fn(node) {
				return
			}
		}
		assertFail("Node of %v did not pass the expected node tests: %v", node.kind, strings.Join(mapf(fns, getFuncName), ", "))
	}
}

func getFuncName[T any](fn T) string {
	v := reflect.ValueOf(fn)
	assert(v.Kind() == reflect.Func, "Expected fn to be a Func")
	name := runtime.FuncForPC(v.Pointer()).Name()
	return name[strings.LastIndex(name, ".")+1:]
}

// Asserts 'node' is either nil or matches one of the provided node tests
//
//	assertNodeOpt(node, isTypeNode)
//	assertNodeOpt(node, isFunctionBody, isExpression)
func assertNodeOpt(node *Node, fns ...NodeTest) {
	if node != nil {
		assertNode(node, fns...)
	}
}

// Asserts 'nodes' is not nil and that every element matches one of the provided node tests
//
//	assertNodes(arguments, isExpression)
//	assertNodes(properties, isObjectLiteralElement)
func assertNodes(nodes []*Node, fns ...NodeTest) {
	assert(nodes != nil, "Expected nodes to be present")
	for _, node := range nodes {
		assertNode(node, fns...)
	}
}

// Asserts 'nodes' is either nil or that every element matches one of the provided node tests
//
//	assertNodesOpt(arguments, isExpression) // e.g., for a NewExpression like `new Foo`
func assertNodesOpt(nodes []*Node, fns ...NodeTest) {
	if nodes != nil {
		assertNodes(nodes, fns...)
	}
}

// Asserts 'nodes' is not nil, has at least one element, and that every element matches one of the provided node tests
//
//	assertNonEmptyNodes(types, isTypeNode)
func assertNonEmptyNodes(nodes []*Node, fns ...NodeTest) {
	assert(nodes != nil, "Expected nodes to be present")
	assert(len(nodes) > 0, "Expected nodes to be non-empty")
	assertNodes(nodes, fns...)
}

// Asserts 'nodes' is either nil or has at least one element and that every element matches one of the provided node tests
func assertNonEmptyNodesOpt(nodes []*Node, fns ...NodeTest) {
	if nodes != nil {
		assertNonEmptyNodes(nodes, fns...)
	}
}

// common assertions

func assertModifiersOpt(node *Node) {
	assertNodeOpt(node, isModifierList)
}

func assertBindingName(node *Node) {
	assertNode(node, isBindingName)
}

func assertIdentifierName(node *Node) {
	assertNode(node, isIdentifierName)
}

func assertIdentifierReference(node *Node) {
	assertNode(node, isIdentifierReference)
}

func assertBindingIdentifier(node *Node) {
	assertNode(node, isBindingIdentifier)
}

func assertBindingIdentifierOpt(node *Node) {
	assertNodeOpt(node, isBindingIdentifier)
}

func assertLabelIdentifier(node *Node) {
	assertNode(node, isLabelIdentifier)
}

func assertEntityName(node *Node) {
	assertNode(node, isEntityName)
}

func assertEntityNameOpt(node *Node) {
	assertNodeOpt(node, isEntityName)
}

func assertPropertyName(node *Node) {
	assertNode(node, isPropertyName)
}

func assertMemberName(node *Node) {
	assertNode(node, isMemberName)
}

func assertModuleExportName(node *Node) {
	assertNode(node, isModuleExportName)
}

func assertModuleExportNameOpt(node *Node) {
	assertNodeOpt(node, isModuleExportName)
}

func assertLeftHandSideExpression(node *Node) {
	assertNode(node, isLeftHandSideExpression)
}

func assertExpression(node *Node) {
	assertNode(node, isExpression)
}

func assertExpressionOpt(node *Node) {
	assertNodeOpt(node, isExpression)
}

func assertJsxTagNameExpression(node *Node) {
	assertNode(node, isIdentifierName, isThisExpression, isPropertyAccessExpression, isJsxNamespacedName)
}

func assertStatement(node *Node) {
	assertNode(node, isStatement)
}

func assertStatements(nodes []*Node) {
	assertNodes(nodes, isStatement)
}

func assertStatementOpt(node *Node) {
	assertNodeOpt(node, isStatement)
}

func assertBlock(node *Node) {
	assertNode(node, isBlock)
}

func assertBlockOpt(node *Node) {
	assertNodeOpt(node, isBlock)
}

func assertFunctionBody(node *Node) {
	assertNode(node, isFunctionBody)
}

func assertFunctionBodyOpt(node *Node) {
	assertNodeOpt(node, isFunctionBody)
}

func assertTypeNode(node *Node) {
	assertNode(node, isTypeNode)
}

func assertTypeNodeOpt(node *Node) {
	assertNodeOpt(node, isTypeNode)
}

func assertTypeParametersOpt(node *Node) {
	assertNodeOpt(node, isTypeParameterList)
}

func assertTypeArgumentsOpt(node *Node) {
	assertNodeOpt(node, isTypeArgumentList)
}

func assertParameters(nodes []*Node) {
	assertNodes(nodes, isParameter)
}

func assertAsteriskTokenOpt(node *Node) {
	assertTokenOpt(node, SyntaxKindAsteriskToken)
}

func assertQuestionTokenOpt(node *Node) {
	assertTokenOpt(node, SyntaxKindQuestionToken)
}

func assertQuestionDotTokenOpt(node *Node) {
	assertTokenOpt(node, SyntaxKindQuestionDotToken)
}

func assertDotDotDotTokenOpt(node *Node) {
	assertTokenOpt(node, SyntaxKindDotDotDotToken)
}

func assertArguments(nodes []*Node) {
	assertNodes(nodes, isExpression)
}

func assertArgumentsOpt(nodes []*Node) {
	assertNodesOpt(nodes, isExpression)
}
