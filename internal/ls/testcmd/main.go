//go:build ignore

package main

import (
	"fmt"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func main() {
	src := `export interface Foo {
  /** JSDoc */
  foo(): void;
}`
	
	fooPos := strings.Index(src, "foo")
	fmt.Printf("'foo' starts at position: %d\n", fooPos)
	
	file := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: tspath.GetNormalizedAbsolutePath("test.ts", "/"),
	}, src, core.ScriptKindTS)
	
	safeText := func(pos, end int) string {
		if pos < 0 { pos = 0 }
		if end > len(src) { end = len(src) }
		if pos > end { return "<invalid>" }
		return fmt.Sprintf("%q", src[pos:end])
	}
	
	// Find the method signature
	var findMethodSig func(node *ast.Node) *ast.Node
	findMethodSig = func(node *ast.Node) *ast.Node {
		if node.Kind == ast.KindMethodSignature { return node }
		var result *ast.Node
		node.ForEachChild(func(child *ast.Node) bool {
			r := findMethodSig(child)
			if r != nil { result = r; return true }
			return false
		})
		return result
	}
	
	methodSig := findMethodSig(file.AsNode())
	fmt.Printf("MethodSig: Pos=%d End=%d\n", methodSig.Pos(), methodSig.End())
	
	jsdocs := methodSig.JSDoc(file)
	fmt.Printf("JSDoc count: %d\n", len(jsdocs))
	if len(jsdocs) > 0 {
		jsdoc := jsdocs[0]
		fmt.Printf("JSDoc: Pos=%d End=%d text=%s\n", jsdoc.Pos(), jsdoc.End(), safeText(jsdoc.Pos(), jsdoc.End()))
		fmt.Printf("JSDoc HasJSDoc flag: %v\n", jsdoc.Flags&ast.NodeFlagsHasJSDoc != 0)
		fmt.Printf("JSDoc IsJSDoc: %v\n", jsdoc.IsJSDoc())
		
		jdoc := jsdoc.AsJSDoc()
		fmt.Printf("JSDoc.Comment: %v\n", jdoc.Comment)
		if jdoc.Comment != nil {
			fmt.Printf("JSDoc.Comment length: %d\n", len(jdoc.Comment.Nodes))
			for i, node := range jdoc.Comment.Nodes {
				fmt.Printf("  Comment[%d]: %s Pos=%d End=%d text=%s\n", i, node.Kind.String(), node.Pos(), node.End(), safeText(node.Pos(), node.End()))
			}
		}
		fmt.Printf("JSDoc.Tags: %v\n", jdoc.Tags)
		
		// Try walking JSDoc children
		fmt.Print("JSDoc ForEachChild:\n")
		jsdoc.ForEachChild(func(child *ast.Node) bool {
			fmt.Printf("  child: %s Pos=%d End=%d text=%s\n", child.Kind.String(), child.Pos(), child.End(), safeText(child.Pos(), child.End()))
			return false
		})
	}
}
