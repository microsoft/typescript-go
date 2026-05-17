package parser

import (
"fmt"
"testing"
"github.com/microsoft/typescript-go/internal/ast"
"github.com/microsoft/typescript-go/internal/core"
)

func TestDebugAwaitParsing(t *testing.T) {
source := "export {}\nconst foo = await { bar: 42 }\n"
opts := ast.SourceFileParseOptions{
FileName: "/test.ts",
}
file := ParseSourceFile(opts, source, core.ScriptKindTS)

fmt.Println("Statements:", len(file.Statements.Nodes))
for i, stmt := range file.Statements.Nodes {
fmt.Printf("  [%d] Kind=%v Flags=%v Pos=%d End=%d\n", i, stmt.Kind, stmt.Flags, stmt.Pos(), stmt.End())
if stmt.Kind == ast.KindVariableStatement {
vs := stmt.AsVariableStatement()
dl := vs.DeclarationList.AsVariableDeclarationList()
fmt.Printf("       DeclarationList declarations: %d\n", len(dl.Declarations.Nodes))
for j, decl := range dl.Declarations.Nodes {
vd := decl.AsVariableDeclaration()
fmt.Printf("       [%d] Name=%v\n", j, vd.Name().Kind)
if vd.Initializer != nil {
fmt.Printf("            InitKind=%v\n", vd.Initializer.Kind)
}
}
}
}
fmt.Println("\nDiagnostics:", len(file.Diagnostics()))
for _, d := range file.Diagnostics() {
fmt.Printf("  pos=%d\n", d.Pos())
}
fmt.Println("\nExternalModuleIndicator:", file.ExternalModuleIndicator != nil)
}
