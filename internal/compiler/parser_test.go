package compiler_test

import (
	"bytes"
	"fmt"
	"io/fs"
	"iter"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/testutil/baseline"
	"github.com/microsoft/typescript-go/internal/tspath"
	"gotest.tools/v3/assert"
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

// compare current code's tsgo AST with tsc's AST, but only write local baselines for tsgo's AST.
// How to use:
// 1. In _submodules/TypeScript, run `npm install` and `npx hereby services --no-typecheck`
// 2. Run this test manually (you might not need 50 minutes, or you might need more on Windows)
//    TEST_ALL=ALL go test ./... -run TestParseAgainstTSC -timeout 50m
// 3. If all tests pass, you're done! If not, you can look at the local output to see if it looks wrong.
// 4. If there are lots of failures, or the failure isn't obvious, run
//    node internal/compiler/testdata/baselineAST.js -r _submodules/TypeScript testdata/baselines/gold
//    This writes the tsc output to disk.
// 5. Now diff gold/ and local/
// 6. To run a single file, 
//    TEST_ALL=ALL go test ./... -run TestParseSingleAgainstTSC -args -filename=tests/baselines/reference/parserVariableDeclaration1.js
func TestParseAgainstTSC(t *testing.T) {
	if os.Getenv("TEST_ALL") == "" {
		t.Skip()
	}
	t.Parallel()
	// TODO: Either build tsc first or document that you have to build it yourself
	err := filepath.WalkDir(repo.TypeScriptSubmodulePath, parseTestComparisonWorker(t))
	if err != nil {
		t.Fatalf("Error walking the path %q: %v", repo.TypeScriptSubmodulePath, err)
	}
}

func TestParseSingleAgainstTSC(t *testing.T) {
	if os.Getenv("TEST_ALL") == "" {
		t.Skip()
	}
	t.Parallel()
	parseTestComparisonWorker(t)(filepath.Join(repo.TypeScriptSubmodulePath, "tests/baselines/reference/dynamicImportsDeclaration.js"), nil, nil)
}

func parseTestComparisonWorker(t *testing.T) func(fileName string, d fs.DirEntry, err error) error {
	return func(fileName string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d != nil && d.IsDir() {
			return nil
		}
		testName, _ := filepath.Rel(repo.TypeScriptSubmodulePath, fileName)
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			if isIgnoredTestFile(fileName) {
				t.Skip()
			}
			sourceText, err := os.ReadFile(fileName)
			outputFilename := generateOutputFileName(t, fileName)
			assert.NilError(t, err)
			var expected string
			cachetext, err := os.ReadFile(filepath.Join(repo.TestDataPath, "baselines", "gold", outputFilename))
			if err == nil {
				expected = string(cachetext)
			} else {
				cmd := exec.Command("node", "testdata/baselineAST.js", fileName)
				var stderr bytes.Buffer
				var stdout bytes.Buffer
				cmd.Stderr = &stderr
				cmd.Stdout = &stdout
				err = cmd.Run()
				if err != nil {
					t.Fatalf("Error running the command %q: %v\nStderr: %s", cmd.String(), err, stderr.String())
				}
				expected = stdout.String()
			}
			actual := printAST(compiler.ParseSourceFile(fileName, string(sourceText), core.ScriptTargetESNext))
			baseline.RunFromText(t, outputFilename, expected, actual, baseline.Options{})
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
		skind, _ := strings.CutPrefix(node.Kind.String(), "Kind")
		if node.Kind == ast.KindImportSpecifier {
			parent = node
		}
		switch node.Kind {
		case ast.KindSyntaxList:
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

func TestParseTypeScriptRepo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		ignoreErrors bool
	}{
		{"src", false},
		{"scripts", false},
		{"Herebyfile.mjs", false},
		{"tests/cases", true},
	}

	for _, test := range tests {
		root := filepath.Join(repo.TypeScriptSubmodulePath, test.name)
		if _, err := os.Stat(root); os.IsNotExist(err) {
			t.Skipf("%q does not exist", root)
		}

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			for f := range allParsableFiles(t, root) {
				t.Run(f.name, func(t *testing.T) {
					t.Parallel()

					// !!! TODO: Fix this bug
					if f.name == "compiler/unicodeEscapesInNames01.ts" {
						t.Skip("times out")
					}

					sourceText, err := os.ReadFile(f.path)
					assert.NilError(t, err)

					var sourceFile *ast.SourceFile

					if strings.HasSuffix(f.name, ".json") {
						sourceFile = compiler.ParseJSONText(f.path, string(sourceText))
					} else {
						sourceFile = compiler.ParseSourceFile(f.path, string(sourceText), core.ScriptTargetESNext)
					}

					if !test.ignoreErrors {
						assert.Equal(t, len(sourceFile.Diagnostics()), 0)
					}
				})
			}
		})
	}
}

type parsableFile struct {
	path string
	name string
}

func allParsableFiles(tb testing.TB, root string) iter.Seq[parsableFile] {
	tb.Helper()
	return func(yield func(parsableFile) bool) {
		tb.Helper()
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() || tspath.TryGetExtensionFromPath(path) == "" {
				return nil
			}

			testName, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			testName = filepath.ToSlash(testName)

			if !yield(parsableFile{path, testName}) {
				return filepath.SkipAll
			}
			return nil
		})
		assert.NilError(tb, err)
	}
}
