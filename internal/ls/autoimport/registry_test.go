package autoimport_test

import (
	"context"
	"testing"

	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
	"gotest.tools/v3/assert"
)

func BenchmarkAutoImportRegistry(b *testing.B) {
	checkerURI := lsconv.FileNameToDocumentURI(tspath.CombinePaths(repo.TypeScriptSubmodulePath, "src/compiler/checker.ts"))
	checkerContent, ok := osvfs.FS().ReadFile(checkerURI.FileName())
	assert.Assert(b, ok, "failed to read checker.ts")

	for b.Loop() {
		b.StopTimer()
		session, _ := projecttestutil.SetupWithRealFS()
		session.DidOpenFile(context.Background(), checkerURI, 1, checkerContent, lsproto.LanguageKindTypeScript)
		b.StartTimer()

		_, err := session.GetLanguageServiceWithAutoImports(context.Background(), checkerURI)
		assert.NilError(b, err)
	}
}
