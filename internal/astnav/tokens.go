package astnav

import (
	"fmt"
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/scanner"
)

var factoryPool = sync.Pool{
	New: func() any {
		return &ast.NodeFactory{}
	},
}

func getNodeFactory() *ast.NodeFactory {
	return factoryPool.Get().(*ast.NodeFactory)
}

func putNodeFactory(factory *ast.NodeFactory) {
	factoryPool.Put(factory)
}

func GetTouchingPropertyName(sourceFile *ast.SourceFile, position int) *ast.Node {
	return getTokenAtPosition(sourceFile, position, false /*allowPositionInLeadingTrivia*/, func(node *ast.Node) bool {
		return ast.IsPropertyNameLiteral(node) || ast.IsKeywordKind(node.Kind) || ast.IsPrivateIdentifier(node)
	})
}

func GetTokenAtPosition(sourceFile *ast.SourceFile, position int) *ast.Node {
	return getTokenAtPosition(sourceFile, position, true /*allowPositionInLeadingTrivia*/, nil)
}

func getTokenAtPosition(
	sourceFile *ast.SourceFile,
	position int,
	allowPositionInLeadingTrivia bool,
	includePrecedingTokenAtEndPosition func(node *ast.Node) bool,
) *ast.Node {
	var next, prevSubtree *ast.Node
	factory := getNodeFactory()
	current := sourceFile.AsNode()
	left := 0
	right := -1

	testNode := func(node *ast.Node) int {
		if node.End() < position {
			return -1
		}

		start := getPosition(node, sourceFile, allowPositionInLeadingTrivia)
		if start > position {
			return 1
		}

		match, endMatch := nodeContainsPosition(node, start, position, sourceFile)
		if endMatch != nil && includePrecedingTokenAtEndPosition != nil {
			prevSubtree = endMatch
		}
		if match {
			return 0
		}

		return -1
	}

	visitNodes := func(nodes []*ast.Node) {
		index, match := core.BinarySearchUniqueFunc(nodes, position, func(middle int, node *ast.Node) int {
			cmp := testNode(node)
			if cmp < 0 {
				left = node.End()
			}
			return cmp
		})

		if match {
			next = nodes[index]
		} else if index < len(nodes) {
			right = getPosition(nodes[index], sourceFile, allowPositionInLeadingTrivia)
		}
	}

	nodeVisitor := &ast.NodeVisitor{
		Visit: func(node *ast.Node) *ast.Node {
			return node
		},
		Hooks: ast.NodeVisitorHooks{
			VisitNode: func(node *ast.Node, visitor *ast.NodeVisitor) *ast.Node {
				if node != nil && next == nil && right < 0 {
					switch testNode(node) {
					case -1:
						if !ast.IsJSDocKind(node.Kind) {
							// We can't move the left boundary into or beyond JSDoc,
							// because we may end up returning the token after this JSDoc,
							// constructing it with the scanner, and we need to include
							// all its leading trivia in its position.
							left = node.End()
						}
					case 0:
						if node.Flags&ast.NodeFlagsHasJSDoc != 0 {
							visitNodes(node.JSDoc(sourceFile))
						}
						if next == nil {
							next = node
						}
					case 1:
						right = getPosition(node, sourceFile, allowPositionInLeadingTrivia)
					}
				}
				return node
			},
			VisitNodes: func(nodeList *ast.NodeList, visitor *ast.NodeVisitor) *ast.NodeList {
				if nodeList != nil && len(nodeList.Nodes) > 0 && next == nil && right < 0 {
					start := nodeList.Pos()
					if !ast.IsJSDocKind(nodeList.Nodes[0].Kind) {
						start = getPosition(nodeList.Nodes[0], sourceFile, allowPositionInLeadingTrivia)
					}
					if start > position {
						right = start
					} else if nodeList.End() == position && includePrecedingTokenAtEndPosition != nil && nodeList.Nodes[len(nodeList.Nodes)-1].End() == position {
						left = nodeList.End()
						prevSubtree = nodeList.Nodes[len(nodeList.Nodes)-1]
					} else if nodeList.End() <= position {
						left = nodeList.End()
					} else {
						visitNodes(nodeList.Nodes)
					}
				}
				return nodeList
			},
			VisitModifiers: func(modifiers *ast.ModifierList, visitor *ast.NodeVisitor) *ast.ModifierList {
				if modifiers != nil && next == nil && right < 0 {
					start := getPosition(modifiers.Nodes[0], sourceFile, allowPositionInLeadingTrivia)
					if start > position {
						right = start
					} else if modifiers.End() == position && includePrecedingTokenAtEndPosition != nil {
						left = modifiers.End()
						prevSubtree = modifiers.Nodes[len(modifiers.Nodes)-1]
					} else if modifiers.End() < position {
						left = modifiers.End()
					} else {
						visitNodes(modifiers.Nodes)
					}
				}
				return modifiers
			},
			VisitToken: func(token *ast.TokenNode, visitor *ast.NodeVisitor) *ast.Node {
				if token != nil && next == nil && right < 0 {
					switch testNode(token) {
					case -1:
						left = token.End()
					case 0:
						if next == nil {
							next = token
						}
					case 1:
						right = getPosition(token, sourceFile, allowPositionInLeadingTrivia)
					}
				}
				return token
			},
		},
	}

	for {
		visitEachChildAndJSDoc(current, sourceFile, nodeVisitor)
		if prevSubtree != nil {
			child := findRightmostNode(prevSubtree, sourceFile)
			if child.End() == position && includePrecedingTokenAtEndPosition(child) {
				// Optimization: includePrecedingTokenAtEndPosition only ever returns true
				// for real AST nodes, so we don't run the scanner here.
				return child
			}
			prevSubtree = nil
		}
		if next == nil {
			if ast.IsTokenKind(current.Kind) || ast.IsJSDocCommentContainingNode(current) {
				return current
			}
			if right < 0 {
				right = current.End()
			}
			scanner := scanner.GetScannerForSourceFile(sourceFile, left)
			for left < right {
				token := scanner.Token()
				tokenFullStart := scanner.TokenFullStart()
				tokenStart := core.IfElse(allowPositionInLeadingTrivia, tokenFullStart, scanner.TokenStart())
				tokenEnd := scanner.TokenEnd()
				if tokenStart <= position && (position < tokenEnd) {
					if token == ast.KindIdentifier || !ast.IsTokenKind(token) {
						if ast.IsJSDocKind(current.Kind) {
							return current
						}
						panic(fmt.Sprintf("did not expect %s to have %s in its trivia", current.Kind.String(), token.String()))
					}
					tokenNode := factory.NewToken(token)
					tokenNode.Loc = core.NewTextRange(tokenFullStart, tokenEnd)
					tokenNode.Parent = current
					return tokenNode
				}
				if includePrecedingTokenAtEndPosition != nil && tokenEnd == position {
					prevToken := factory.NewToken(token)
					prevToken.Loc = core.NewTextRange(tokenFullStart, tokenEnd)
					prevToken.Parent = current
					if includePrecedingTokenAtEndPosition(prevToken) {
						return prevToken
					}
				}
				left = tokenEnd
				scanner.Scan()
			}
			return current
		}
		current = next
		left = current.Pos()
		right = -1
		next = nil
	}
}

func getPosition(node *ast.Node, sourceFile *ast.SourceFile, allowPositionInLeadingTrivia bool) int {
	if allowPositionInLeadingTrivia {
		return node.Pos()
	}
	return scanner.GetTokenPosOfNode(node, sourceFile, true /*includeJsDoc*/)
}

func nodeContainsPosition(node *ast.Node, nodeStart int, position int, sourceFile *ast.SourceFile) (result bool, prevSubtree *ast.Node) {
	if nodeStart > position {
		return false, nil
	}
	if position < node.End() {
		return true, nil
	}
	return false, core.IfElse(position == node.End(), node, nil)
}

func findRightmostNode(node *ast.Node, sourceFile *ast.SourceFile) *ast.Node {
	var next *ast.Node
	current := node
	visitNode := func(node *ast.Node, _ *ast.NodeVisitor) *ast.Node {
		if node != nil {
			next = node
		}
		return node
	}
	visitor := &ast.NodeVisitor{
		Visit: func(node *ast.Node) *ast.Node {
			return node
		},
		Hooks: ast.NodeVisitorHooks{
			VisitNode:  visitNode,
			VisitToken: visitNode,
			VisitNodes: func(nodeList *ast.NodeList, visitor *ast.NodeVisitor) *ast.NodeList {
				if nodeList != nil && len(nodeList.Nodes) > 0 {
					next = nodeList.Nodes[len(nodeList.Nodes)-1]
				}
				return nodeList
			},
			VisitModifiers: func(modifiers *ast.ModifierList, visitor *ast.NodeVisitor) *ast.ModifierList {
				if modifiers != nil && len(modifiers.Nodes) > 0 {
					next = modifiers.Nodes[len(modifiers.Nodes)-1]
				}
				return modifiers
			},
		},
	}

	for {
		current.VisitEachChild(visitor)
		if next == nil {
			return current
		}
		current = next
		next = nil
	}
}

func visitEachChildAndJSDoc(node *ast.Node, sourceFile *ast.SourceFile, visitor *ast.NodeVisitor) {
	if node.Flags&ast.NodeFlagsHasJSDoc != 0 {
		for _, jsDoc := range node.JSDoc(sourceFile) {
			if visitor.Hooks.VisitNode != nil {
				visitor.Hooks.VisitNode(jsDoc, visitor)
			} else {
				visitor.VisitNode(jsDoc)
			}
		}
	}
	node.VisitEachChild(visitor)
}
