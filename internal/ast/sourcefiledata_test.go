package ast_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
	"gotest.tools/v3/assert"
)

func TestGetOrComputeSourceFileDataComputesOnce(t *testing.T) {
	t.Parallel()

	sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: "/test.ts",
		Path:     "/test.ts",
	}, "const x = 1;", core.ScriptKindTS)

	key := ast.NewSourceFileDataKey[*int]()
	var computes atomic.Int32
	const goroutines = 8

	var wg sync.WaitGroup
	results := make([]*int, goroutines)
	for i := range goroutines {
		wg.Go(func() {
			results[i] = ast.GetOrComputeSourceFileData(sourceFile, key, func(*ast.SourceFile) *int {
				computes.Add(1)
				value := 123
				return &value
			})
		})
	}
	wg.Wait()

	assert.Equal(t, computes.Load(), int32(1))
	for _, result := range results {
		assert.Equal(t, result, results[0])
		assert.Equal(t, *result, 123)
	}
}
