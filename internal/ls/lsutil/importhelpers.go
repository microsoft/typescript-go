package lsutil

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/debug"
)

// GetTypeKeywordOfTypeOnlyImport returns the 'type' keyword token from a type-only import clause.
// This is equivalent to TypeScript's getTypeKeywordOfTypeOnlyImport helper.
func GetTypeKeywordOfTypeOnlyImport(importClause *ast.ImportClause, sourceFile *ast.SourceFile) *ast.Node {
	debug.Assert(importClause.IsTypeOnly(), "import clause must be type-only")
	// The first child of a type-only import clause is the 'type' keyword
	// import type { foo } from './bar'
	//        ^^^^
	typeKeyword := astnav.FindChildOfKind(importClause.AsNode(), ast.KindTypeKeyword, sourceFile)
	debug.Assert(typeKeyword != nil, "type-only import clause should have a type keyword")
	return typeKeyword
}
