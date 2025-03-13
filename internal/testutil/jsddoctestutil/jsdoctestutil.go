// Package jsdoctestutil provides helper functions for testing JSDoc comments.
package jsdoctestutil

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/scanner"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// ParseJSDocFile parses a TypeScript file with JSDoc comments.
func ParseJSDocFile(source string) *ast.SourceFile {
	fileName := "/test.ts"
	return parser.ParseSourceFile(fileName, tspath.Path(fileName), source, core.ScriptTargetESNext, scanner.JSDocParsingModeParseAll)
}

// GetJSDocFromNode returns JSDoc comments from a node.
// Note: This function is for testing purposes only and uses its own logic
// to access private functionality in the parser package.
func GetJSDocFromNode(file *ast.SourceFile, node *ast.Node) []*ast.Node {
	if node.Flags&ast.NodeFlagsHasJSDoc == 0 {
		return nil
	}

	// Let's make this simple check for testing purposes
	// In a real application, we would need access to the JSDoc cache in the parser package
	return nil
}

// HasTag checks if a JSDoc comment has a specific tag.
func HasTag(jsDoc *ast.JSDoc, tagName string) bool {
	if jsDoc.Tags == nil {
		return false
	}

	for _, tag := range jsDoc.Tags.Nodes {
		// Multiple steps are needed to convert to JSDocTag
		if tag.Kind >= ast.KindFirstJSDocTagNode && tag.Kind <= ast.KindLastJSDocTagNode {
			// Define a safe way to get the tag name
			if ident := getTagNameFromJSDocNode(tag); ident != nil {
				if ident.Text() == tagName {
					return true
				}
			}
		}
	}

	return false
}

// FindTag finds and returns a specific tag in a JSDoc comment.
func FindTag(jsDoc *ast.JSDoc, tagName string) *ast.Node {
	if jsDoc.Tags == nil {
		return nil
	}

	for _, tag := range jsDoc.Tags.Nodes {
		// Multiple steps are needed to convert to JSDocTag
		if tag.Kind >= ast.KindFirstJSDocTagNode && tag.Kind <= ast.KindLastJSDocTagNode {
			// Define a safe way to get the tag name
			if ident := getTagNameFromJSDocNode(tag); ident != nil {
				if ident.Text() == tagName {
					return tag
				}
			}
		}
	}

	return nil
}

// getTagNameFromJSDocNode helper function gets the tag name from a JSDoc node.
func getTagNameFromJSDocNode(node *ast.Node) *ast.IdentifierNode {
	// All nodes of type JSDocXXXTag have a TagName field
	// But it can't be accessed directly from the Node type, so we need to check each type individually

	switch node.Kind {
	case ast.KindJSDocParameterTag:
		return node.AsJSDocParameterTag().TagName
	case ast.KindJSDocReturnTag:
		return node.AsJSDocReturnTag().TagName
	case ast.KindJSDocTypeTag:
		return node.AsJSDocTypeTag().TagName
	case ast.KindJSDocDeprecatedTag:
		return node.AsJSDocDeprecatedTag().TagName
	// Similar checks can be added for other JSDoc tag types
	default:
		// Unknown tag type
		return nil
	}
}

// AssertHasJSDoc asserts that a node has a JSDoc comment.
func AssertHasJSDoc(t *testing.T, node *ast.Node) {
	t.Helper()
	if node.Flags&ast.NodeFlagsHasJSDoc == 0 {
		t.Error("Node should have JSDoc comment")
	}
}

// AssertHasTag asserts that a JSDoc comment has a specific tag.
func AssertHasTag(t *testing.T, jsDoc *ast.JSDoc, tagName string) {
	t.Helper()
	if !HasTag(jsDoc, tagName) {
		t.Errorf("JSDoc should have tag '%s'", tagName)
	}
}

// AssertIsDeprecated asserts that a node is marked as deprecated.
func AssertIsDeprecated(t *testing.T, node *ast.Node) {
	t.Helper()
	if node.Flags&ast.NodeFlagsDeprecated == 0 {
		t.Error("Node should be marked as deprecated")
	}
}

// GetJSDocComment gets the text content of a JSDoc comment from a node.
func GetJSDocComment(jsDoc *ast.JSDoc) string {
	if jsDoc.Comment == nil || len(jsDoc.Comment.Nodes) == 0 {
		return ""
	}

	comment := ""
	for _, part := range jsDoc.Comment.Nodes {
		if part.Kind == ast.KindJSDocText {
			comment += part.AsJSDocText().Text
		}
	}

	return comment
}
