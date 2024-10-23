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
	var visit func(node *Node, indentation int) bool
	visit = func(node *Node, indentation int) bool {
		offset := 1
		skind := node.kind.String()[len("SyntaxKind"):]
		switch node.kind {
		case SyntaxKindModifierList, SyntaxKindTypeParameterList, SyntaxKindTypeArgumentList, SyntaxKindSyntaxList:
			offset = 0
		case SyntaxKindIdentifier:
			indent := strings.Repeat("  ", indentation)
			fmt.Fprintf(file, "%s%s: '%s'\n", indent, skind, sourceFile.text[node.loc.pos:node.loc.end])
		default:
			if isOmittedExpression(node) {
				skind = "OmittedExpression"
			}
			indent := strings.Repeat("  ", indentation)
			fmt.Fprintf(file, "%s%s\n", indent, skind)
		}
		// TODO: Include trivia in a more structured way than GetFullText// Visit child nodes with increased indentation
		return node.ForEachChild(func(child *Node) bool {
			visit(child, indentation+offset)
			return false
		})
	}
	visit(sourceFile.AsNode(), 0)
}

func isOmittedExpression(node *Node) bool {
	if node.kind == SyntaxKindBindingElement {
		b := node.AsBindingElement()
		if b.initializer == nil && b.name == nil && b.dotDotDotToken == nil {
			return true
		}
	} else {
		return false
	}
	return false
}
