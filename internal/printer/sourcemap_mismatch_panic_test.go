package printer

import (
    "testing"

    "github.com/microsoft/typescript-go/internal/ast"
    "github.com/microsoft/typescript-go/internal/core"
    "github.com/microsoft/typescript-go/internal/parser"
    "github.com/microsoft/typescript-go/internal/sourcemap"
    "github.com/microsoft/typescript-go/internal/tspath"
)

// This test intentionally reproduces the panic seen when the source map range
// refers to a different file (or position space) than the file used for text lookups.
// It constructs a small source file and forces a SourceMapRange with a huge Pos,
// then triggers printing with source maps enabled. The buggy implementation slices
// text[start:pos] with pos > len(text), panicking.
func TestSourceMapMismatchPanics(t *testing.T) {
    t.Parallel()

    // Small source text (to be the current file)
    sourceText := "export const x = 1;\n"
    sf := parser.ParseSourceFile(ast.SourceFileParseOptions{
        FileName:         "/home/src/workspaces/project/index.ts",
        Path:             "/home/src/workspaces/project/index.ts",
        JSDocParsingMode: ast.JSDocParsingModeParseAll,
    }, sourceText, core.ScriptKindTS)

    // Choose a node to print (the first statement)
    if len(sf.Statements.Nodes) == 0 {
        t.Fatalf("expected at least one statement")
    }
    stmt := sf.Statements.Nodes[0]

    // Create a second, much longer file to act as the "original" mapping source
    longPrefix := make([]byte, 2500)
    for i := range longPrefix {
        longPrefix[i] = ' '
    }
    longText := string(longPrefix) + "export const y = 2;\n"
    sf2 := parser.ParseSourceFile(ast.SourceFileParseOptions{
        FileName:         "/home/src/workspaces/project/long.ts",
        Path:             "/home/src/workspaces/project/long.ts",
        JSDocParsingMode: ast.JSDocParsingModeParseAll,
    }, longText, core.ScriptKindTS)
    if len(sf2.Statements.Nodes) == 0 {
        t.Fatalf("expected at least one statement in long file")
    }
    stmt2 := sf2.Statements.Nodes[0]

    // Build a printer with source maps enabled
    emitCtx := NewEmitContext()
    p := NewPrinter(PrinterOptions{SourceMap: true}, PrintHandlers{}, emitCtx)
    writer := NewTextWriter("\n")
    gen := sourcemap.NewGenerator("index.js", "", "/home/src/workspaces/project", tspath.ComparePathsOptions{})

    // Map the current node to the original node from the longer file.
    // This simulates transformation attaching original positions from a different file tree.
    emitCtx.SetOriginal(stmt, stmt2)
    emitCtx.SetSourceMapRange(stmt, stmt2.Loc)

    // This should not panic if source selection is correct (mapping against original file).
    p.Write(stmt, sf, writer, gen)
}


