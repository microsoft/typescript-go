package compiler

import (
	"fmt"
	"os"
	"path/filepath"
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
	root := "../../_submodules/TypeScript"
	local := "../../tests/baselines/local"
	program := NewProgram(ProgramOptions{
		RootPath: root,
		Options:  &CompilerOptions{Target: ScriptTargetESNext, ModuleKind: ModuleKindNodeNext}})
	// Ensure the target directory exists
	if err := os.MkdirAll(local, os.ModePerm); err != nil {
		t.Errorf("Error creating directory %s: %v\n", local, err)
		return
	}
	for _, sourceFile := range program.SourceFiles() {
		printAST(t, root+"/", local, sourceFile)
	}
}

// prefix specifies the directory to write the baseline
func printAST(t *testing.T, sourceRoot string, targetRoot string, sourceFile *SourceFile) {
	path := filepath.ToSlash(sourceFile.FileName())
	if !strings.HasPrefix(path, sourceRoot) {
		t.Errorf("Error: sourcePrefix %s not found in fileName %s\n", sourceRoot, path)
		return
	}
	path = strings.TrimPrefix(path, sourceRoot)
	path = strings.ReplaceAll(path, "/", "_")
	outputFileName := targetRoot + "/" + path + ".ast.txt"
	file, err := os.Create(outputFileName)
	if err != nil {
		t.Errorf("Error creating file %s: %v\n", outputFileName, err)
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
	sourceFile.ForEachChild(func(n *Node) bool { return v(n, 0) })
}
