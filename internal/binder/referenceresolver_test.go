package binder

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler/diagnostics"
	"gotest.tools/v3/assert"
	"testing"
)

const Identifier = "some-identifier"

func TestNewReferenceResolver(t *testing.T) {
	hooks := ReferenceResolverHooks{}
	resolver := NewReferenceResolver(hooks)
	assert.Assert(t, resolver != nil)
}

func TestGetResolvedSymbol(t *testing.T) {
	hooks := ReferenceResolverHooks{
		GetResolvedSymbol: func(node *ast.Node) *ast.Symbol {
			return &ast.Symbol{}
		},
	}
	resolver := NewReferenceResolver(hooks)
	node := &ast.Node{}
	symbol := resolver.(*referenceResolver).getResolvedSymbol(node)
	assert.Assert(t, symbol != nil)
}

func TestGetMergedSymbol(t *testing.T) {
	hooks := ReferenceResolverHooks{
		GetMergedSymbol: func(symbol *ast.Symbol) *ast.Symbol {
			return &ast.Symbol{}
		},
	}
	resolver := NewReferenceResolver(hooks)
	symbol := &ast.Symbol{}
	mergedSymbol := resolver.(*referenceResolver).getMergedSymbol(symbol)
	assert.Assert(t, mergedSymbol != nil)
}

func TestGetParentOfSymbol(t *testing.T) {
	hooks := ReferenceResolverHooks{
		GetParentOfSymbol: func(symbol *ast.Symbol) *ast.Symbol {
			return &ast.Symbol{}
		},
	}
	resolver := NewReferenceResolver(hooks)
	symbol := &ast.Symbol{}
	parentSymbol := resolver.(*referenceResolver).getParentOfSymbol(symbol)
	assert.Assert(t, parentSymbol != nil)
}

func TestGetSymbolOfDeclaration(t *testing.T) {
	hooks := ReferenceResolverHooks{
		GetSymbolOfDeclaration: func(declaration *ast.Declaration) *ast.Symbol {
			return &ast.Symbol{}
		},
	}
	resolver := NewReferenceResolver(hooks)
	declaration := &ast.Declaration{}
	symbol := resolver.(*referenceResolver).getSymbolOfDeclaration(declaration)
	assert.Assert(t, symbol != nil)
}

func TestGetReferencedValueSymbol(t *testing.T) {
	hooks := ReferenceResolverHooks{
		ResolveName: func(location *ast.Node, name string, meaning ast.SymbolFlags, nameNotFoundMessage *diagnostics.Message, isUse bool, excludeGlobals bool) *ast.Symbol {
			return &ast.Symbol{}
		},
	}
	resolver := NewReferenceResolver(hooks)

	nf := ast.NodeFactory{}
	reference := nf.NewIdentifier("sometext")
	symbol := resolver.(*referenceResolver).getReferencedValueSymbol(reference, false)
	assert.Assert(t, symbol != nil)
}

func TestIsTypeOnlyAliasDeclaration(t *testing.T) {
	hooks := ReferenceResolverHooks{
		GetTypeOnlyAliasDeclaration: func(symbol *ast.Symbol, include ast.SymbolFlags) *ast.Declaration {
			return &ast.Declaration{}
		},
	}
	resolver := NewReferenceResolver(hooks)
	symbol := &ast.Symbol{}
	isTypeOnly := resolver.(*referenceResolver).isTypeOnlyAliasDeclaration(symbol)
	assert.Assert(t, isTypeOnly)
}

func TestGetDeclarationOfAliasSymbol(t *testing.T) {
	resolver := NewReferenceResolver(ReferenceResolverHooks{})
	symbol := &ast.Symbol{
		Declarations: []*ast.Declaration{
			{Kind: ast.KindImportEqualsDeclaration},
		},
	}
	declaration := resolver.(*referenceResolver).getDeclarationOfAliasSymbol(symbol)
	assert.Assert(t, declaration != nil)
}

func TestGetExportSymbolOfValueSymbolIfExported(t *testing.T) {
	hooks := ReferenceResolverHooks{
		GetExportSymbolOfValueSymbolIfExported: func(symbol *ast.Symbol) *ast.Symbol {
			return &ast.Symbol{}
		},
	}
	resolver := NewReferenceResolver(hooks)
	symbol := &ast.Symbol{}
	exportSymbol := resolver.(*referenceResolver).getExportSymbolOfValueSymbolIfExported(symbol)
	assert.Assert(t, exportSymbol != nil)
}

func TestGetReferencedExportContainer(t *testing.T) {
	hooks := ReferenceResolverHooks{
		ResolveName: func(location *ast.Node, name string, meaning ast.SymbolFlags, nameNotFoundMessage *diagnostics.Message, isUse bool, excludeGlobals bool) *ast.Symbol {
			return &ast.Symbol{}
		},
	}
	resolver := NewReferenceResolver(hooks)
	NF := ast.NodeFactory{}
	node := NF.NewIdentifier("sometext")
	container := resolver.GetReferencedExportContainer(node, false)
	assert.Assert(t, container == nil)
}

func TestGetReferencedImportDeclaration(t *testing.T) {
	hooks := ReferenceResolverHooks{
		ResolveName: func(location *ast.Node, name string, meaning ast.SymbolFlags, nameNotFoundMessage *diagnostics.Message, isUse bool, excludeGlobals bool) *ast.Symbol {
			return &ast.Symbol{}
		},
	}
	resolver := NewReferenceResolver(hooks)
	NF := ast.NodeFactory{}
	node := NF.NewIdentifier("sometext")
	declaration := resolver.GetReferencedImportDeclaration(node)
	assert.Assert(t, declaration == nil)
}

func TestGetReferencedValueDeclaration(t *testing.T) {
	hooks := ReferenceResolverHooks{
		ResolveName: func(location *ast.Node, name string, meaning ast.SymbolFlags, nameNotFoundMessage *diagnostics.Message, isUse bool, excludeGlobals bool) *ast.Symbol {
			return &ast.Symbol{}
		},
	}
	resolver := NewReferenceResolver(hooks)
	NF := ast.NodeFactory{}
	node := NF.NewIdentifier("sometext")
	declaration := resolver.GetReferencedValueDeclaration(node)
	assert.Assert(t, declaration == nil)
}

func TestGetReferencedValueDeclarations(t *testing.T) {
	hooks := ReferenceResolverHooks{
		ResolveName: func(location *ast.Node, name string, meaning ast.SymbolFlags, nameNotFoundMessage *diagnostics.Message, isUse bool, excludeGlobals bool) *ast.Symbol {
			return &ast.Symbol{}
		},
	}
	resolver := NewReferenceResolver(hooks)
	NF := ast.NodeFactory{}
	node := NF.NewIdentifier(Identifier)
	declarations := resolver.GetReferencedValueDeclarations(node)
	assert.Assert(t, len(declarations) == 0)
}
