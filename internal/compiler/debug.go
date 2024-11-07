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

// TODO(rbuckton): Ideally this would be in a subpackage except that it would introduce a circular dependency with
// ast.go. If necessary, these could be made to be generic, e.g.:
//
//	type NodeLike interface {
//		Pos() int
//		End() int
//		Kind() SyntaxKind
//	}
//	type NodeLikeTest[T NodeLike] func (node T) bool
//
//	func assertNode[T NodeLike](node T, fns ...NodeLikeTest[T]) {
//		assert(!isNil(node), "Expected node to be present")
//		if len(fns) > 0 {
//			for _, fn := range fns {
//				if fn(node) {
//					return
//				}
//			}
//			assertFail("Node of %v did not pass the expected node tests: %v", node.Kind(), strings.Join(mapf(fns, getFunctionName), ", "))
//		}
//	}

func isNil[T any](v T) bool {
	switch r := reflect.ValueOf(v); r.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return r.IsNil()
	default:
		return false
	}
}

// TODO(rbuckton): Rename or make _assert public as it now conflicts with "gotest.tools/v3/assert"

// Asserts 'cond' is true or panics with a formatted message
func _assert(cond bool, a ...any) {
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
	_assert(false, a...)
}

// Asserts 'x' is nil
func assertMissing(x any) {
	_assert(isNil(x), "Expected node to be missing")
}

func assertDefined(x any) {
	_assert(!isNil(x), "Expected node to be present")
}

func checkDefined[T any](x T) T {
	assertDefined(x)
	return x
}

// Asserts 'kind' is one of the expected values
//
//	assertKind(kind, SyntaxKindQuestionToken, SyntaxKindExclamationToken)
//	assertKind(kind, SyntaxKindCaseClause, SyntaxKindDefaultClause)
func assertKind(kind SyntaxKind, expected ...SyntaxKind) {
	_assert(slices.Contains(expected, kind), "kind '%v' was not one of the expected values: %v", kind, expected)
}

// Asserts 'node' is not nil and is a Token with one of the expected 'kinds'
//
//	assertToken(asteriskToken)
//	assertToken(postfixToken, SyntaxKindQuestionToken, SyntaxKindExclamationToken)
func assertToken(node *Node, kinds ...SyntaxKind) {
	assertDefined(node)
	_assert(node.kind >= SyntaxKindFirstToken && node.kind <= SyntaxKindLastToken, "Node of %v was not a token", node.kind)
	_assert(len(kinds) == 0 || nodeKindIs(node, kinds...), "Token '%v' did not have one of the expected kinds: %v", kinds)
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
	assertDefined(node)
	if len(fns) > 0 {
		for _, fn := range fns {
			if fn(node) {
				return
			}
		}
		assertFail("Node of %v did not pass the expected node tests: %v", node.kind, strings.Join(mapf(fns, getFunctionName), ", "))
	}
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
	_assert(nodes != nil, "Expected nodes to be present")
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
	_assert(nodes != nil, "Expected nodes to be present")
	_assert(len(nodes) > 0, "Expected nodes to be non-empty")
	assertNodes(nodes, fns...)
}

// Asserts 'node' is not nil and matches one of the provided node tests
//
//	assertNode(node, isExpression)
//	assertNode(node, isIdentifierName, isStringLiteral)
func assertNotNode(node *Node, fns ...NodeTest) {
	assertDefined(node)
	_assert(len(fns) > 0, "Expected one or more test functions")
	for _, fn := range fns {
		if fn(node) {
			assertFail("Node of %v did not pass the expected node tests: %v", node.kind, strings.Join(mapf(fns, getFunctionName), ", "))
		}
	}
}

func getFunctionName[T any](fn T) string {
	v := reflect.ValueOf(fn)
	_assert(v.Kind() == reflect.Func, "Expected fn to be a Func")
	name := runtime.FuncForPC(v.Pointer()).Name()
	return name[strings.LastIndex(name, ".")+1:]
}
