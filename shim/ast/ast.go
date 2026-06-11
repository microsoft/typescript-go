// Package ast re-exports the internal AST API for external consumers of this
// fork. See package github.com/microsoft/typescript-go/shim for the governing
// policy.
package ast

import (
	"github.com/microsoft/typescript-go/internal/ast"
)

type (
	Node       = ast.Node
	NodeList   = ast.NodeList
	SourceFile = ast.SourceFile
	Symbol     = ast.Symbol
	Diagnostic = ast.Diagnostic
	Kind       = ast.Kind
)

// Node kinds the schema walker matches on. Expand as walker code demands.
const (
	KindNumericLiteral           = ast.KindNumericLiteral
	KindStringLiteral            = ast.KindStringLiteral
	KindIdentifier               = ast.KindIdentifier
	KindDecorator                = ast.KindDecorator
	KindPropertyDeclaration      = ast.KindPropertyDeclaration
	KindMethodDeclaration        = ast.KindMethodDeclaration
	KindTypeReference            = ast.KindTypeReference
	KindArrayType                = ast.KindArrayType
	KindUnionType                = ast.KindUnionType
	KindObjectLiteralExpression  = ast.KindObjectLiteralExpression
	KindPropertyAccessExpression = ast.KindPropertyAccessExpression
	KindCallExpression           = ast.KindCallExpression
	KindClassDeclaration         = ast.KindClassDeclaration
	KindInterfaceDeclaration     = ast.KindInterfaceDeclaration
	KindTypeAliasDeclaration     = ast.KindTypeAliasDeclaration
	KindEnumDeclaration          = ast.KindEnumDeclaration
	KindHeritageClause           = ast.KindHeritageClause
	KindSourceFile               = ast.KindSourceFile
)
