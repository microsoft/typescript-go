package compiler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/testutil/baseline"
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
	program := NewProgram(ProgramOptions{
		RootPath: repo.TypeScriptSubmodulePath,
		Options:  &CompilerOptions{Target: ScriptTargetESNext, ModuleKind: ModuleKindNodeNext}})
	for _, sourceFile := range program.SourceFiles() {
		path := filepath.ToSlash(sourceFile.fileName)
		path = path[len(repo.TypeScriptSubmodulePath)+1:]
		path = strings.ReplaceAll(path, "/", "_")
		outputFileName := path + ".ast"
		if outputFileName == "tests_cases_compiler_binderBinaryExpressionStressJs.ts.ast" ||
			outputFileName == "tests_cases_compiler_binderBinaryExpressionStress.ts.ast" {
			continue
		}
		baseline.Run(outputFileName, printAST(sourceFile), baseline.Options{})
	}
}

// prefix specifies the directory to write the baseline
func printAST(sourceFile *SourceFile) string {
	sb := strings.Builder{}
	var visit func(node *Node, indentation int) bool
	visit = func(node *Node, indentation int) bool {
		offset := 1
		skind := node.kind.String()[len("SyntaxKind"):]
		switch node.kind {
		case SyntaxKindModifierList, SyntaxKindTypeParameterList, SyntaxKindTypeArgumentList, SyntaxKindSyntaxList:
			offset = 0
		case SyntaxKindIdentifier:
			indent := strings.Repeat("  ", indentation)
			sb.WriteString(fmt.Sprintf("%s%s: '%s'\n", indent, skind, sourceFile.text[node.loc.pos:node.loc.end]))
		default:
			if isOmittedExpression(node) {
				skind = "OmittedExpression"
			}
			indent := strings.Repeat("  ", indentation)
			sb.WriteString(fmt.Sprintf("%s%s\n", indent, skind))
		}
		// TODO: Include trivia in a more structured way than GetFullText// Visit child nodes with increased indentation
		return node.ForEachChild(func(child *Node) bool {
			if node.kind == SyntaxKindShorthandPropertyAssignment && node.AsShorthandPropertyAssignment().objectAssignmentInitializer == child {
				indent := strings.Repeat("  ", indentation+offset)
				sb.WriteString(fmt.Sprintf("%sEqualsToken\n", indent)) // print an extra line for the EqualsToken
			}
			visit(child, indentation+offset)
			return false
		})
	}
	visit(sourceFile.AsNode(), 0)
	return sb.String()
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
