package checker

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/printer"
)

func (c *Checker) runWithoutResolvedSignatureCaching(node *ast.Node, fn func() *Signature) *Signature {
	ancestorNode := ast.FindAncestor(node, func(n *ast.Node) bool {
		return ast.IsCallLikeOrFunctionLikeExpression(n)
	})
	if ancestorNode != nil {
		cachedResolvedSignatures := make(map[*SignatureLinks]*Signature)
		cachedTypes := make(map[*ValueSymbolLinks]*Type)
		for ancestorNode != nil {
			signatureLinks := c.signatureLinks.Get(ancestorNode)
			cachedResolvedSignatures[signatureLinks] = signatureLinks.resolvedSignature
			signatureLinks.resolvedSignature = nil
			if ast.IsFunctionExpressionOrArrowFunction(ancestorNode) {
				symbolLinks := c.valueSymbolLinks.Get(c.getSymbolOfDeclaration(ancestorNode))
				resolvedType := symbolLinks.resolvedType
				cachedTypes[symbolLinks] = resolvedType
				symbolLinks.resolvedType = nil
			}
			ancestorNode = ast.FindAncestor(ancestorNode.Parent, ast.IsCallLikeOrFunctionLikeExpression)
		}
		result := fn()
		for signatureLinks, resolvedSignature := range cachedResolvedSignatures {
			signatureLinks.resolvedSignature = resolvedSignature
		}
		for symbolLinks, resolvedType := range cachedTypes {
			symbolLinks.resolvedType = resolvedType
		}
		return result
	}
	return fn()
}

func (c *Checker) getResolvedSignatureWorker(node *ast.Node, candidatesOutArray *[]*Signature, checkMode CheckMode) *Signature {
	parsedNode := printer.NewEmitContext().ParseNode(node)
	var res *Signature = nil
	if parsedNode != nil {
		res = c.getResolvedSignature(parsedNode, candidatesOutArray, checkMode)
	}
	return res
}

func (c *Checker) GetResolvedSignatureForSignatureHelp(node *ast.Node, candidatesOutArray *[]*Signature) *Signature {
	return c.runWithoutResolvedSignatureCaching(node, func() *Signature {
		return c.getResolvedSignatureWorker(node, candidatesOutArray, CheckModeIsForSignatureHelp)
	})
}
