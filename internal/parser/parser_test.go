package parser_test

import (
	"io/fs"
	"iter"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/testrunner"
	"github.com/microsoft/typescript-go/internal/testutil/fixtures"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
	"gotest.tools/v3/assert"
)

func BenchmarkParse(b *testing.B) {
	for _, f := range fixtures.BenchFixtures {
		b.Run(f.Name(), func(b *testing.B) {
			f.SkipIfNotExist(b)

			fileName := tspath.GetNormalizedAbsolutePath(f.Path(), "/")
			path := tspath.ToPath(fileName, "/", osvfs.FS().UseCaseSensitiveFileNames())
			sourceText := f.ReadFile(b)
			scriptKind := core.GetScriptKindFromFileName(fileName)

			opts := ast.SourceFileParseOptions{
				FileName: fileName,
				Path:     path,
			}

			for b.Loop() {
				parser.ParseSourceFile(opts, sourceText, scriptKind)
			}
		})
	}
}

func TestDeeplyNestedParserConstructs(t *testing.T) {
	t.Parallel()

	const depth = 16_000
	tests := []struct {
		name               string
		sourceText         string
		expectedStatements int
	}{
		{
			name:               "expression",
			sourceText:         strings.Repeat("(", depth) + "x" + strings.Repeat(")", depth) + ";",
			expectedStatements: 1,
		},
		{
			name:               "unterminated expression",
			sourceText:         strings.Repeat("(", depth) + "x;",
			expectedStatements: 1,
		},
		{
			name:               "unterminated expression with trailing statement",
			sourceText:         strings.Repeat("(", depth) + "x; y;",
			expectedStatements: 2,
		},
		{
			name:               "unterminated expression with close paren in comment",
			sourceText:         strings.Repeat("(", depth) + "x; // )",
			expectedStatements: 1,
		},
		{
			name:               "unterminated expression before later parens",
			sourceText:         strings.Repeat("(", depth) + "x; f();",
			expectedStatements: 2,
		},
		{
			name:               "unterminated expression before close paren string",
			sourceText:         strings.Repeat("(", depth) + `x; ")";`,
			expectedStatements: 2,
		},
		{
			name:               "unterminated array expression",
			sourceText:         strings.Repeat("(", depth) + "[x",
			expectedStatements: 1,
		},
		{
			name:               "unterminated object expression",
			sourceText:         strings.Repeat("(", depth) + "{a: x",
			expectedStatements: 1,
		},
		{
			name:               "unterminated division expression",
			sourceText:         strings.Repeat("(", depth) + "x / y;",
			expectedStatements: 1,
		},
		{
			name:               "unterminated regular expression",
			sourceText:         strings.Repeat("(", depth) + `x === /[)]/;`,
			expectedStatements: 1,
		},
		{
			name:               "type",
			sourceText:         "type T = " + strings.Repeat("(", depth) + "string" + strings.Repeat(")", depth) + ";",
			expectedStatements: 1,
		},
		{
			name:               "unterminated type",
			sourceText:         "type T = " + strings.Repeat("(", depth) + "string;",
			expectedStatements: 1,
		},
		{
			name:               "prefix unary expression",
			sourceText:         strings.Repeat("!", depth) + "x;",
			expectedStatements: 1,
		},
		{
			name:               "type assertion expression",
			sourceText:         strings.Repeat("<T>", depth) + "x;",
			expectedStatements: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			opts := ast.SourceFileParseOptions{
				FileName: "/index.ts",
				Path:     "/index.ts",
			}
			file := parser.ParseSourceFile(opts, test.sourceText, core.ScriptKindTS)
			assert.Equal(t, len(file.Statements.Nodes), test.expectedStatements)
		})
	}
}

func TestDeeplyNestedParserStackBound(t *testing.T) {
	t.Parallel()

	const depth = 16_000
	tests := []string{
		"expression",
		"comment",
		"laterParens",
		"string",
		"array",
		"object",
		"division",
		"regularExpression",
		"type",
		"prefixUnary",
		"delete",
		"typeof",
		"void",
		"await",
		"typeAssertion",
	}

	if testName := os.Getenv("TSGO_DEEP_PARSER_STACK_CHILD"); testName != "" {
		debug.SetMaxStack(512 << 10)
		opts := ast.SourceFileParseOptions{
			FileName: "/index.ts",
			Path:     "/index.ts",
		}
		parser.ParseSourceFile(opts, deeplyNestedParserSource(testName, depth), core.ScriptKindTS)
		return
	}

	for _, testName := range tests {
		cmd := exec.Command(os.Args[0], "-test.run=^TestDeeplyNestedParserStackBound$")
		cmd.Env = append(os.Environ(), "TSGO_DEEP_PARSER_STACK_CHILD="+testName)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("%s exceeded the parser stack bound: %v\n%s", testName, err, output)
		}
	}
}

func deeplyNestedParserSource(testName string, depth int) string {
	switch testName {
	case "expression":
		return strings.Repeat("(", depth) + "x;"
	case "comment":
		return strings.Repeat("(", depth) + "x; // )"
	case "laterParens":
		return strings.Repeat("(", depth) + "x; f();"
	case "string":
		return strings.Repeat("(", depth) + `x; ")";`
	case "array":
		return strings.Repeat("(", depth) + "[x"
	case "object":
		return strings.Repeat("(", depth) + "{a: x"
	case "division":
		return strings.Repeat("(", depth) + "x / y;"
	case "regularExpression":
		return strings.Repeat("(", depth) + `x === /[)]/;`
	case "type":
		return "type T = " + strings.Repeat("(", depth) + "string;"
	case "prefixUnary":
		return strings.Repeat("!~+-", depth/4) + "x;"
	case "delete":
		return strings.Repeat("delete ", depth) + "x;"
	case "typeof":
		return strings.Repeat("typeof ", depth) + "x;"
	case "void":
		return strings.Repeat("void ", depth) + "x;"
	case "await":
		return "async function f() { return " + strings.Repeat("await ", depth) + "x; }"
	case "typeAssertion":
		return strings.Repeat("<T>", depth) + "x;"
	default:
		panic("unknown deeply nested parser test")
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

func FuzzParser(f *testing.F) {
	repo.SkipIfNoTypeScriptSubmodule(f)

	tests := []string{
		"src",
		"scripts",
		"Herebyfile.mjs",
	}

	var extensions collections.Set[string]
	for _, es := range tspath.AllSupportedExtensionsWithJson {
		for _, e := range es {
			extensions.Add(e)
		}
	}

	for _, test := range tests {
		root := filepath.Join(repo.TypeScriptSubmodulePath(), test)

		for file := range allParsableFiles(f, root) {
			sourceText, err := os.ReadFile(file.path)
			assert.NilError(f, err)
			extension := tspath.TryGetExtensionFromPath(file.path)
			f.Add(extension, string(sourceText), false, false)
		}
	}

	testDirs := []string{
		filepath.Join(repo.TypeScriptSubmodulePath(), "tests/cases/compiler"),
		filepath.Join(repo.TypeScriptSubmodulePath(), "tests/cases/conformance"),
		filepath.Join(repo.TestDataPath(), "tests/cases/compiler"),
	}

	for _, testDir := range testDirs {
		if _, err := os.Stat(testDir); os.IsNotExist(err) {
			continue
		}

		for file := range allParsableFiles(f, testDir) {
			sourceText, err := os.ReadFile(file.path)
			assert.NilError(f, err)

			type testFile struct {
				content string
				name    string
			}

			testUnits, _, _, _, err := testrunner.ParseTestFilesAndSymlinks(
				string(sourceText),
				file.path,
				func(filename string, content string, fileOptions map[string]string) (testFile, error) {
					return testFile{content: content, name: filename}, nil
				},
			)
			assert.NilError(f, err)

			for _, unit := range testUnits {
				extension := tspath.TryGetExtensionFromPath(unit.name)
				if extension == "" {
					continue
				}
				f.Add(extension, unit.content, false, false)
			}
		}
	}

	f.Fuzz(func(t *testing.T, extension string, sourceText string, externalModuleIndicatorOptionsJSX bool, externalModuleIndicatorOptionsForce bool) {
		if !extensions.Has(extension) {
			t.Skip()
		}

		fileName := "/index" + extension
		path := tspath.Path(fileName)

		opts := ast.SourceFileParseOptions{
			FileName: fileName,
			Path:     path,
			ExternalModuleIndicatorOptions: ast.ExternalModuleIndicatorOptions{
				JSX:   externalModuleIndicatorOptionsJSX,
				Force: externalModuleIndicatorOptionsForce,
			},
		}

		parser.ParseSourceFile(opts, sourceText, core.GetScriptKindFromFileName(fileName))
	})
}

func TestJSDocImportTypeParentChain(t *testing.T) {
	t.Parallel()
	sourceText := `test("", async function () {
  ;(/** @type {typeof import("a")} */ ({}))
})

test("", async function () {
  ;(/** @type {typeof import("a")} */ a)
})

test("", async function () {
  (/** @type {typeof import("a")} */ ({}))
  ;(/** @type {typeof import("a")} */ ({}))
})

test("", async function () {
  (/** @type {typeof import("a")} */ a)
  ;(/** @type {typeof import("a")} */ a)
})

test("", async function () {
  (/** @type {typeof import("a")} */ ({}))
  ;(/** @type {typeof import("a")} */ ({}))
})
`
	opts := ast.SourceFileParseOptions{
		FileName: "/index.js",
		Path:     "/index.js",
	}

	file := parser.ParseSourceFile(opts, sourceText, core.ScriptKindJS)

	for i := 1; i < len(file.ReparsedClones); i++ {
		a, b := file.ReparsedClones[i-1], file.ReparsedClones[i]
		if a.Pos() == b.Pos() && a.End() == b.End() && a.Kind == b.Kind {
			t.Errorf("duplicate ReparsedClones at [%d] and [%d]: %s pos=%d end=%d", i-1, i, a.Kind.String(), a.Pos(), a.End())
		}
	}

	for _, imp := range file.Imports() {
		reparsed := ast.GetReparsedNodeForNode(imp)
		if ast.GetSourceFileOfNode(reparsed) == nil {
			t.Errorf("reparsed import at pos=%d has broken parent chain", imp.Pos())
		}
	}
}

func TestSourceFileContainsNonASCIIInStringLiteralFastPath(t *testing.T) {
	t.Parallel()
	sourceText := `const x = "─";

namespace N {
  export const y = x;
}
`
	opts := ast.SourceFileParseOptions{
		FileName: "/index.ts",
		Path:     "/index.ts",
	}

	file := parser.ParseSourceFile(opts, sourceText, core.ScriptKindTS)

	assert.Assert(t, file.ContainsNonASCII)
	positionMap := file.GetPositionMap()
	assert.Assert(t, !positionMap.IsAsciiOnly())
	afterBoxDrawingCharacter := strings.Index(sourceText, "─") + len("─")
	assert.Equal(t, positionMap.UTF8ToUTF16(afterBoxDrawingCharacter), afterBoxDrawingCharacter-2)
	assert.Equal(t, positionMap.UTF8ToUTF16(len(sourceText)), len(sourceText)-2)
}
