package compiler

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

	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/testutil/baseline"
	"gotest.tools/v3/assert"
)

func BenchmarkParse(b *testing.B) {
	for _, f := range benchFixtures {
		b.Run(f.Name(), func(b *testing.B) {
			f.SkipIfNotExist(b)

			fileName := f.Path()
			sourceText := f.ReadFile(b)

			for i := 0; i < b.N; i++ {
				ParseSourceFile(fileName, sourceText, ScriptTargetESNext)
			}
		})
	}
}

// TODO: Need to be able to compare against local/reference *locally*
// and against local/gold as part of a test.
// this is two different test cases really, but running them both as part of a full test run is a bad idea

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
		if isIgnoredTestFile(fileName) {
			return nil
		}
		testName, _ := filepath.Rel(repo.TypeScriptSubmodulePath, fileName)
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			// if isIgnoredTestFile(fileName) {
			// 	t.Skip()
			// }
			sourceText, err := os.ReadFile(fileName)
			assert.NilError(t, err)
			sourceFile := ParseSourceFile(fileName, string(sourceText), ScriptTargetESNext)
			err = baseline.Run(generateOutputFileName(t, fileName), printAST(sourceFile), baseline.Options{Gold: true})
			assert.NilError(t, err)
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
			strings.Contains(name, "reference/tsc") ||
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
func printAST(sourceFile *SourceFile) string {
	var sb strings.Builder
	var visit func(node *Node, indentation int) bool
	visit = func(node *Node, indentation int) bool {
		offset := 1
		skind, _ := strings.CutPrefix(node.kind.String(), "SyntaxKind")
		switch node.kind {
		case SyntaxKindModifierList, SyntaxKindTypeParameterList, SyntaxKindTypeArgumentList, SyntaxKindSyntaxList:
			offset = 0
		case SyntaxKindIdentifier:
			indent := getIndentation(indentation)
			sb.WriteString(fmt.Sprintf("%s%s: '%s'\n", indent, skind, sourceFile.text[node.loc.pos:node.loc.end]))
		default:
			if isOmittedExpression(node) {
				skind = "OmittedExpression"
			}
			indent := strings.Repeat("  ", indentation)
			sb.WriteString(fmt.Sprintf("%s%s\n", indent, skind))
		}
		// TODO: Include trivia in a more structured way than GetFullText
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
