package encoder_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/microsoft/typescript-go/internal/api/encoder"
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/scanner"
	"github.com/microsoft/typescript-go/internal/testutil/baseline"
	"gotest.tools/v3/assert"
)

func TestEncodeSourceFile(t *testing.T) {
	t.Parallel()
	sourceFile := parser.ParseSourceFile("/test.ts", "/test.ts", "export function foo<T, U>(a: string, b: string): any {}\nfoo();", core.ScriptTargetESNext, scanner.JSDocParsingModeParseAll)
	ast.SetParentInChildren(sourceFile.AsNode())

	t.Run("baseline", func(t *testing.T) {
		t.Parallel()
		buf, err := encoder.EncodeSourceFile(sourceFile)
		assert.NilError(t, err)

		str := encoder.FormatEncodedSourceFile(buf)
		baseline.Run(t, "encodeSourceFile.txt", str, baseline.Options{
			Subfolder: "api",
		})
	})
}

func BenchmarkEncodeSourceFile(b *testing.B) {
	repo.SkipIfNoTypeScriptSubmodule(b)
	filePath := filepath.Join(repo.TypeScriptSubmodulePath, "src/compiler/checker.ts")
	fileContent, err := os.ReadFile(filePath)
	assert.NilError(b, err)
	sourceFile := parser.ParseSourceFile(
		"/checker.ts",
		"/checker.ts",
		string(fileContent),
		core.ScriptTargetESNext,
		scanner.JSDocParsingModeParseAll,
	)

	for b.Loop() {
		_, err := api.EncodeSourceFile(sourceFile)
		assert.NilError(b, err)
	}
}
