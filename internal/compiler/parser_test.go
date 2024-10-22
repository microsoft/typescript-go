package compiler

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func BenchmarkParse(b *testing.B) {
	fileName := "../../_submodules/TypeScript/src/compiler/checker.ts"
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
	}
	sourceText := string(bytes)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseSourceFile(fileName, sourceText, ScriptTargetESNext)
	}
}

func TestParseAndPrintNodes(t *testing.T) {
	// TODO: Should be an entire Program's sourcefiles
    fileName := "../../_submodules/TypeScript/src/compiler/checker.ts"
    bytes, err := os.ReadFile(fileName)
    if err != nil {
        t.Fatal(err)
    }
    sourceText := string(bytes)
    sourceFile := ParseSourceFile(fileName, sourceText, ScriptTargetESNext)
	printAST(t, "tests/baselines/local", sourceFile)
}

// prefix specifies the directory to write the baseline
func printAST(t *testing.T, prefix string, sourceFile *SourceFile) {
	flattenedPath := strings.ReplaceAll(sourceFile.FileName(), "/", "_")
	flattenedPath = strings.ReplaceAll(flattenedPath, "\\", "_") // windows
	outputFileName := prefix + "/" + flattenedPath + ".ast.txt"
    file, err := os.Create(outputFileName)
    if err != nil {
		t.Errorf("Error creating file %s: %v\n", outputFileName, err)
        return
    }
    defer file.Close()
    // Recursive function to visit nodes with indentation
	var v func(node *Node, indentation int) bool
    v = func(node *Node, indentation int) bool {
        indent := strings.Repeat("  ", indentation)
        if node.kind == SyntaxKindIdentifier {
            fmt.Fprintf(file, "%s%s: '%s'\n", indent, node.kind, sourceFile.text[node.loc.pos:node.loc.end])
        } else {
            fmt.Fprintf(file, "%s%s\n", indent, node.kind)
        }
        // TODO: Include trivia
        // Visit child nodes with increased indentation
		return node.ForEachChild(func(child *Node) bool {
            v(child, indentation+1)
			return false
        })
    }
	sourceFile.ForEachChild(func (n *Node) bool { return v(n, 0) })
}
