package format_test

import (
"testing"
"github.com/microsoft/typescript-go/internal/ast"
"github.com/microsoft/typescript-go/internal/core"
"github.com/microsoft/typescript-go/internal/format"
"github.com/microsoft/typescript-go/internal/parser"
"fmt"
)

func TestReproduceIssue(t *testing.T) {
ctx := format.WithFormatCodeSettings(t.Context(), &format.FormatCodeSettings{
EditorSettings: format.EditorSettings{
TabSize:                4,
IndentSize:             4,
BaseIndentSize:         0,
NewLineCharacter:       "\n",
ConvertTabsToSpaces:    false,
IndentStyle:            format.IndentStyleSmart,
TrimTrailingWhitespace: true,
},
InsertSpaceBeforeTypeAnnotation: core.TSTrue,
}, "\n")

// Test case 1: console.log with comment
originalText1 := `console.log(
"a",
// the second arg
"b"
);`

sourceFile1 := parser.ParseSourceFile(ast.SourceFileParseOptions{
FileName: "/test.ts",
Path:     "/test.ts",
}, originalText1, core.ScriptKindTS)

edits1 := format.FormatDocument(ctx, sourceFile1)
formatted1 := applyBulkEdits(originalText1, edits1)

fmt.Println("Test case 1: console.log with comment")
fmt.Println("Original:")
fmt.Println(originalText1)
fmt.Println("\nFormatted:")
fmt.Println(formatted1)
fmt.Println()

// Test case 2: chained method calls with comment
originalText2 := `foo
.bar()
// A second call
.baz();`

sourceFile2 := parser.ParseSourceFile(ast.SourceFileParseOptions{
FileName: "/test.ts",
Path:     "/test.ts",
}, originalText2, core.ScriptKindTS)

edits2 := format.FormatDocument(ctx, sourceFile2)
formatted2 := applyBulkEdits(originalText2, edits2)

fmt.Println("Test case 2: chained method calls with comment")
fmt.Println("Original:")
fmt.Println(originalText2)
fmt.Println("\nFormatted:")
fmt.Println(formatted2)
}
