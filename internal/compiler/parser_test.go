package compiler_test

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/testutil/baseline"
	"gotest.tools/v3/assert"

	"github.com/microsoft/typescript-go/internal/core"
)

func BenchmarkParse(b *testing.B) {
	for _, f := range compiler.BenchFixtures {
		b.Run(f.Name(), func(b *testing.B) {
			f.SkipIfNotExist(b)

			fileName := f.Path()
			sourceText := f.ReadFile(b)

			for i := 0; i < b.N; i++ {
				compiler.ParseSourceFile(fileName, sourceText, core.ScriptTargetESNext)
			}
		})
	}
}

func TestParseAndPrintNodes(t *testing.T) {
	t.Parallel()
	err := filepath.WalkDir(repo.TypeScriptSubmodulePath, parseTestWorker(t, &baseline.Options{}))
	if err != nil {
		t.Fatalf("Error walking the path %q: %v", repo.TypeScriptSubmodulePath, err)
	}
}

func TestParseAgainstTSC(t *testing.T) {
	t.Parallel()
	goldDir := "../../testdata/baselines/gold"
	entries, err := os.ReadDir(goldDir)
	if err != nil || len(entries) == 0 {
		cmd := exec.Command("node", "testdata/baselineAST.js", "../../_submodules/TypeScript/", goldDir)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			t.Fatalf("Error running the command %q: %v\nStderr: %s", cmd.String(), err, stderr.String())
		}
	}
	err = filepath.WalkDir(repo.TypeScriptSubmodulePath, parseTestWorker(t, &baseline.Options{Gold: true}))
	if err != nil {
		t.Fatalf("Error walking the path %q: %v", repo.TypeScriptSubmodulePath, err)
	}
}

func parseTestWorker(t *testing.T, options *baseline.Options) func(fileName string, d fs.DirEntry, err error) error {
	return func(fileName string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		testName, _ := filepath.Rel(repo.TypeScriptSubmodulePath, fileName)
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			if isIgnoredTestFile(fileName) {
				t.Skip()
			}
			sourceText, err := os.ReadFile(fileName)
			assert.NilError(t, err)
			sourceFile := compiler.ParseSourceFile(fileName, string(sourceText), core.ScriptTargetESNext)
			baseline.Run(t, generateOutputFileName(t, fileName), printAST(sourceFile), baseline.Options{Gold: true})
		})
		return nil
	}
}

func isIgnoredTestFile(name string) bool {
	ext := filepath.Ext(name)
	return !(ext == ".ts" || ext == ".js" || ext == ".tsx" || ext == ".jsx") ||
		// Too deep for a simmple Javascript tree walker
		(strings.HasSuffix(name, "binderBinaryExpressionStress.ts") ||
			strings.HasSuffix(name, "binderBinaryExpressionStress.js") ||
			strings.HasSuffix(name, "binderBinaryExpressionStressJs.ts") ||
			strings.HasSuffix(name, "binderBinaryExpressionStressJs.js") ||
			// very large minified code
			strings.Contains(name, "codeMirrorModule") ||
			// not actually .js
			strings.Contains(name, "reference/config/") ||
			strings.Contains(name, "reference/tsc") ||
			strings.Contains(name, "reference/tsserver") ||
			strings.Contains(name, "reference/tsbuild"))
}

func generateOutputFileName(t *testing.T, fileName string) string {
	path, err := filepath.Rel(repo.TypeScriptSubmodulePath, fileName)
	if err != nil {
		t.Errorf("%s is outside of the TypeScript submodule", fileName)
	}
	return strings.ReplaceAll(path, string(filepath.Separator), "_") + ".ast"
}

var (
	indentationCache   map[int]string = make(map[int]string)
	indentationCacheMu sync.Mutex
)

func getIndentation(level int) string {
	indentationCacheMu.Lock()
	defer indentationCacheMu.Unlock()
	if indent, ok := indentationCache[level]; ok {
		return indent
	}
	indent := strings.Repeat("  ", level)
	indentationCache[level] = indent
	return indent
}

// prefix specifies the directory to write the baseline
func printAST(sourceFile *ast.SourceFile) string {
	var sb strings.Builder
	var visit func(node *ast.Node, indentation int) bool
	var parent *ast.Node
	visit = func(node *ast.Node, indentation int) bool {
		offset := 1
		skind, _ := strings.CutPrefix(node.Kind.String(), "SyntaxKind")
		if node.Kind == ast.KindImportSpecifier {
			parent = node
		}
		switch node.Kind {
		case ast.KindModifierList, ast.KindTypeParameterList, ast.KindTypeArgumentList, ast.KindSyntaxList:
			offset = 0
		case ast.KindIdentifier:
			indent := getIndentation(indentation)
			if parent != nil && parent.AsImportSpecifier().Name() == node && node.AsIdentifier().Text == "" && sourceFile.Text[node.Loc.Pos():node.Loc.End()] != "" {
				sb.WriteString(fmt.Sprintf("%s%s: '%s'\n", indent, skind, ""))
			} else {
				sb.WriteString(fmt.Sprintf("%s%s: '%s'\n", indent, skind, sourceFile.Text[node.Loc.Pos():node.Loc.End()]))
			}
		default:
			if isOmittedExpression(node) {
				skind = "OmittedExpression"
			}
			indent := strings.Repeat("  ", indentation)
			sb.WriteString(fmt.Sprintf("%s%s\n", indent, skind))
		}
		// TODO: Include trivia in a more structured way than GetFullText
		return node.ForEachChild(func(child *ast.Node) bool {
			if node.Kind == ast.KindShorthandPropertyAssignment && node.AsShorthandPropertyAssignment().ObjectAssignmentInitializer == child {
				indent := strings.Repeat("  ", indentation+offset)
				sb.WriteString(indent + "EqualsToken\n") // print an extra line for the EqualsToken
			}
			visit(child, indentation+offset)
			return false
		})
	}
	visit(sourceFile.AsNode(), 0)
	return sb.String()
}

func isOmittedExpression(node *ast.Node) bool {
	if node.Kind == ast.KindBindingElement {
		b := node.AsBindingElement()
		if b.Initializer == nil && b.Name() == nil && b.DotDotDotToken == nil {
			return true
		}
	} else {
		return false
	}
	return false
}
