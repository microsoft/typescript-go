package astnav_test

import (
"fmt"
"os"
"path/filepath"
"testing"

"github.com/microsoft/typescript-go/internal/ast"
"github.com/microsoft/typescript-go/internal/astnav"
"github.com/microsoft/typescript-go/internal/core"
"github.com/microsoft/typescript-go/internal/parser"
"github.com/microsoft/typescript-go/internal/repo"
)

func TestDebugPrecedingToken(t *testing.T) {
repo.SkipIfNoTypeScriptSubmodule(t)
fileName := filepath.Join(repo.TypeScriptSubmodulePath(), "src/services/mapCode.ts")
fileText, err := os.ReadFile(fileName)
if err != nil { t.Fatal(err) }
file := parser.ParseSourceFile(ast.SourceFileParseOptions{
FileName: "/file.ts", Path: "/file.ts",
}, string(fileText), core.ScriptKindTS)

// Find FunctionDeclaration[1607,3727)
var stmt *ast.Node
file.AsNode().ForEachChild(func(node *ast.Node) bool {
if node.Pos() <= 1733 && node.End() > 1733 {
stmt = node
return true
}
return false
})
if stmt == nil { t.Fatal("no stmt"); return }
fmt.Printf("stmt: %s [%d,%d)\n", stmt.Kind, stmt.Pos(), stmt.End())

// Check ALL children (not just first)
fmt.Println("Children of stmt:")
stmt.ForEachChild(func(c *ast.Node) bool {
fmt.Printf("  child: %s [%d,%d)\n", c.Kind, c.Pos(), c.End())
return false
})

// FindPrecedingToken with ExEx to trace
for _, pos := range []int{1732, 1733} {
tok := astnav.FindPrecedingToken(file, pos)
if tok == nil {
fmt.Printf("FindPrecedingToken(%d) = nil\n", pos)
} else {
fmt.Printf("FindPrecedingToken(%d) = %s [%d,%d)\n", pos, tok.Kind, tok.Pos(), tok.End())
}
}
}
