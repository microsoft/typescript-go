package api_test

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/microsoft/typescript-go/internal/api"
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
	sourceFile := parser.ParseSourceFile("/test.ts", "/test.ts", "export const x = 1; foo();", core.ScriptTargetESNext, scanner.JSDocParsingModeParseAll)
	ast.SetParentInChildren(sourceFile.AsNode())

	t.Run("baseline", func(t *testing.T) {
		t.Parallel()
		buf, err := api.EncodeSourceFile(sourceFile)
		assert.NilError(t, err)

		decoded, err := decodeInt32s(buf)
		assert.NilError(t, err)

		str := api.FormatEncodedSourceFile(decoded)
		baseline.Run(t, "encodeSourceFile.txt", str, baseline.Options{
			Subfolder: "api",
		})
	})

	t.Run("verify next", func(t *testing.T) {
		t.Parallel()
		buf, err := api.EncodeSourceFile(sourceFile)
		assert.NilError(t, err)

		decoded, err := decodeInt32s(buf)
		assert.NilError(t, err)

		for i := api.EncodedNodeLength; i < len(decoded); i += api.EncodedNodeLength {
			next := decoded[i+api.EncodedNext]
			parent := decoded[i+api.EncodedParent]
			if next != 0 {
				for j := i + api.EncodedNodeLength; j < int(next)*api.EncodedNodeLength; j += api.EncodedNodeLength {
					if j == int(next) {
						// Ensure 'next' has the same parent
						assert.Assert(t, decoded[j+api.EncodedParent] == parent)
					} else {
						// ...and no others on the way also do
						assert.Assert(t, decoded[j+api.EncodedParent] != parent)
					}
				}
			} else {
				// No subsequent node should have the same parent
				for j := i + api.EncodedNodeLength; j < len(decoded); j += api.EncodedNodeLength {
					assert.Assert(t, decoded[j+api.EncodedParent] != parent)
				}
			}
		}
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

func decodeInt32s(buf []byte) ([]int32, error) {
	count := len(buf) / 4
	result := make([]int32, count)
	if err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, result); err != nil {
		return nil, err
	}
	return result, nil
}
