package ls

import (
	"fmt"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/scanner"
)

func getTouchingPropertyName(sourceFile *ast.SourceFile, position int) *ast.Node {
	return getTokenAtPosition(sourceFile, position, false, false, func(node *ast.Node) bool {
		return ast.IsPropertyNameLiteral(node) || ast.IsKeywordKind(node.Kind) || ast.IsPrivateIdentifier(node)
	})
}

func getTouchingPropertyName_fast(sourceFile *ast.SourceFile, position int) *ast.Node {
	return getTokenAtPosition_fast(sourceFile, position, false, false, func(node *ast.Node) bool {
		return ast.IsPropertyNameLiteral(node) || ast.IsKeywordKind(node.Kind) || ast.IsPrivateIdentifier(node)
	})
}

func getTokenAtPosition(
	sourceFile *ast.SourceFile,
	position int,
	allowPositionInLeadingTrivia bool,
	includeEndPosition bool,
	includePrecedingTokenAtEndPosition func(node *ast.Node) bool,
) *ast.Node {
	var foundToken *ast.Node
	current := sourceFile.AsNode()
	factory := getNodeFactory()
	for {
		children := getNodeChildren(current, sourceFile, factory)
		index, match := core.BinarySearchUniqueFunc(children, position, func(middle int, node *ast.Node) int {
			// This last callback is more of a selector than a comparator -
			// `0` causes the `node` result to be returned
			// `1` causes recursion on the left of the middle
			// `-1` causes recursion on the right of the middle

			// Let's say you have 3 nodes, spanning positons
			// pos: 1, end: 3
			// pos: 3, end: 3
			// pos: 3, end: 5
			// and you're looking for the token at positon 3 - all 3 of these nodes are overlapping with position 3.
			// In fact, there's a _good argument_ that node 2 shouldn't even be allowed to exist - depending on if
			// the start or end of the ranges are considered inclusive, it's either wholly subsumed by the first or the last node.
			// Unfortunately, such nodes do exist. :( - See fourslash/completionsImport_tsx.tsx - empty jsx attributes create
			// a zero-length node.
			// What also you may not expect is that which node we return depends on the includePrecedingTokenAtEndPosition flag.
			// Specifically, if includePrecedingTokenAtEndPosition is set, we return the 1-3 node, while if it's unset, we
			// return the 3-5 node. (The zero length node is never correct.) This is because the includePrecedingTokenAtEndPosition
			// flag causes us to return the first node whose end position matches the position and which produces and acceptable token
			// kind. Meanwhile, if includePrecedingTokenAtEndPosition is unset, we look for the first node whose start is <= the
			// position and whose end is greater than the position.

			// There are more sophisticated end tests later, but this one is very fast
			// and allows us to skip a bunch of work
			if node.End() < position {
				return -1
			}

			start := getPosition(node, sourceFile, allowPositionInLeadingTrivia)
			if start > position {
				return 1
			}

			// first element whose start position is before the input and whose end position is after or equal to the input
			if result, found := nodeContainsPosition(node, start, position, sourceFile, includeEndPosition, factory, includePrecedingTokenAtEndPosition); result {
				if found != nil {
					foundToken = found
				}
				if middle > 0 {
					// we want the _first_ element that contains the position, so left-recur if the prior node also contains the position
					prevNode := children[middle-1]
					prevNodeStart := getPosition(prevNode, sourceFile, allowPositionInLeadingTrivia)
					if result, found := nodeContainsPosition(prevNode, prevNodeStart, position, sourceFile, includeEndPosition, factory, includePrecedingTokenAtEndPosition); result {
						if found != nil {
							foundToken = found
						}
						return 1
					}
				}
				return 0
			}

			// this complex condition makes us left-recur around a zero-length node when includePrecedingTokenAtEndPosition is set, rather than right-recur on it
			if includePrecedingTokenAtEndPosition != nil && start == position && middle > 0 && children[middle-1].End() == position {
				prevNode := children[middle-1]
				prevNodeStart := getPosition(prevNode, sourceFile, allowPositionInLeadingTrivia)
				if result, found := nodeContainsPosition(prevNode, prevNodeStart, position, sourceFile, includeEndPosition, factory, includePrecedingTokenAtEndPosition); result {
					if found != nil {
						foundToken = found
					}
					return 1
				}
			}
			return -1
		})

		if foundToken != nil {
			return foundToken
		}
		if match {
			current = children[index]
			continue
		}
		return current
	}
}

func getTokenAtPosition_fast(
	sourceFile *ast.SourceFile,
	position int,
	allowPositionInLeadingTrivia bool,
	includeEndPosition bool,
	includePrecedingTokenAtEndPosition func(node *ast.Node) bool,
) *ast.Node {
	var next, foundToken, prevSubtree *ast.Node
	factory := getNodeFactory()
	current := sourceFile.AsNode()
	left := 0
	right := -1

	testNode := func(node *ast.Node) int {
		if node.End() < position {
			return -1
		}

		start := getPosition(node, sourceFile, allowPositionInLeadingTrivia)
		if node.Kind == ast.KindJSDoc && allowPositionInLeadingTrivia {
			// There is a bug where JSDoc nodes don't include their leading trivia in their start position,
			// which breaks binary searching, since the token *after* them has a position *before* them.
			start = node.Parent.Pos()
		}
		if start > position {
			return 1
		}

		if result, found := nodeContainsPosition(node, start, position, sourceFile, includeEndPosition, factory, includePrecedingTokenAtEndPosition); result {
			if found != nil {
				foundToken = found
			}
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
					} else if nodeList.End() == position && (includeEndPosition || includePrecedingTokenAtEndPosition != nil) && nodeList.Nodes[len(nodeList.Nodes)-1].End() == position {
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
					} else if modifiers.End() == position && (includeEndPosition || includePrecedingTokenAtEndPosition != nil) {
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
		if foundToken != nil {
			return foundToken
		}
		if next == nil {
			if prevSubtree != nil {
				child := findRightmostNode(prevSubtree, sourceFile)
				if child.End() == position && includePrecedingTokenAtEndPosition(child) {
					// Optimization: includePrecedingTokenAtEndPosition only ever returns true
					// for real AST nodes, so we don't run the scanner here.
					return child
				}
			}
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
				if tokenStart <= position && (position < tokenEnd || position == tokenEnd && includeEndPosition) {
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
		if current.Kind == ast.KindJSDoc {
			left = current.Parent.Pos()
		}
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

func maybeSkipTrivia(pos int, sourceFile *ast.SourceFile, allowPositionInLeadingTrivia bool) int {
	if allowPositionInLeadingTrivia {
		return pos
	}
	return scanner.SkipTrivia(sourceFile.Text, pos)
}

func nodeContainsPosition(node *ast.Node, nodeStart int, position int, sourceFile *ast.SourceFile, includeEndPosition bool, factory *ast.NodeFactory, includePrecedingTokenAtEndPosition func(node *ast.Node) bool) (result bool, foundToken *ast.Node) {
	if nodeStart > position {
		// If this child begins after position, then all subsequent children will as well.
		return false, nil
	}
	if position < node.End() || position == node.End() && includeEndPosition {
		return true, nil
	}
	if includePrecedingTokenAtEndPosition != nil && position == node.End() {
		previousToken := findPrecedingToken(position, sourceFile, node, false /*excludeJsDoc*/, factory)
		if previousToken != nil && includePrecedingTokenAtEndPosition(previousToken) {
			return true, previousToken
		}
	}
	return false, nil
}

func findPrecedingToken(position int, sourceFile *ast.SourceFile, startNode *ast.Node, excludeJsDoc bool, factory *ast.NodeFactory) *ast.Node {
	node := startNode
	if node == nil {
		node = sourceFile.AsNode()
	}

	if ast.IsNonWhitespaceToken(node) {
		return node
	}

	children := getNodeChildren(node, sourceFile, factory)
	index, match := core.BinarySearchUniqueFunc(children, position, func(index int, middle *ast.Node) int {
		// This last callback is more of a selector than a comparator -
		// `0` causes the `middle` result to be returned
		// `1` causes recursion on the left of the middle
		// `-1` causes recursion on the right of the middle
		if position < middle.End() {
			// first element whose end position is greater than the input position
			if index == 0 || position >= children[index-1].End() {
				return 0
			}
			return 1
		}
		return -1
	})

	if match {
		child := children[index]
		// Note that the span of a node's tokens is [scanner.GetTokenPosOfNode(node), node.End()).
		// Given that `position < child.End()` and child has constituent tokens, we distinguish these cases:
		// 1) `position` precedes `child`'s tokens or `child` has no tokens (i.e., in a comment or whitespace preceding `child`):
		// we need to find the last token in a previous child.
		// 2) `position` is within the same span: we recurse on `child`.
		if position < child.End() {
			start := scanner.GetTokenPosOfNode(child, sourceFile, !excludeJsDoc)
			lookInPreviousChild := start >= position || // cursor in the leading trivia
				!nodeHasTokens(child, sourceFile) ||
				ast.IsWhitespaceOnlyJsxText(child)
			if lookInPreviousChild {
				// actual start of the node is past the position - previous token should be at the end of previous child
				candidate := findRightmostChildNodeWithTokens(children, index, sourceFile, node.Kind)
				if candidate != nil {
					// Ensure we recurse into JSDoc nodes with children.
					if !excludeJsDoc && ast.IsJSDocCommentContainingNode(candidate) && len(getNodeChildren(candidate, sourceFile, factory)) > 0 {
						return findPrecedingToken(position, sourceFile, candidate, excludeJsDoc, factory)
					}
					return findRightmostToken(candidate, sourceFile, factory)
				}
				return nil
			} else {
				// candidate should be in this node
				return findPrecedingToken(position, sourceFile, child, excludeJsDoc, factory)
			}
		}
	}

	// Here we know that none of child token nodes embrace the position,
	// the only known case is when position is at the end of the file.
	// Try to find the rightmost token in the file without filtering.
	// Namely we are skipping the check: 'position < node.End()'
	if candidate := findRightmostChildNodeWithTokens(children, len(children), sourceFile, node.Kind); candidate != nil {
		return findRightmostToken(candidate, sourceFile, factory)
	}
	return nil
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

func nodeHasTokens(node *ast.Node, sourceFile *ast.SourceFile) bool {
	return scanner.GetTokenPosOfNode(node, sourceFile, false /*includeJsDoc*/) < node.End()
}

func findRightmostToken(n *ast.Node, sourceFile *ast.SourceFile, factory *ast.NodeFactory) *ast.Node {
	if ast.IsNonWhitespaceToken(n) {
		return n
	}

	children := getNodeChildren(n, sourceFile, factory)
	if len(children) == 0 {
		return n
	}

	candidate := findRightmostChildNodeWithTokens(children, len(children), sourceFile, n.Kind)
	if candidate != nil {
		return findRightmostToken(candidate, sourceFile, factory)
	}
	return nil
}

// findRightmostChildNodeWithTokens finds the rightmost child to the left of `children[exclusiveStartPosition]` which is a non-all-whitespace token or has constituent tokens.
func findRightmostChildNodeWithTokens(children []*ast.Node, exclusiveStartPosition int, sourceFile *ast.SourceFile, parentKind ast.Kind) *ast.Node {
	for i := exclusiveStartPosition - 1; i >= 0; i-- {
		child := children[i]

		if ast.IsWhitespaceOnlyJsxText(child) {
			if i == 0 && (parentKind == ast.KindJsxText || parentKind == ast.KindJsxSelfClosingElement) {
				panic("`JsxText` tokens should not be the first child of `JsxElement | JsxSelfClosingElement`")
			}
		} else if nodeHasTokens(children[i], sourceFile) {
			return children[i]
		}
	}
	return nil
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
