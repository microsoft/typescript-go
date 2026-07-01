package project_test

import (
	"context"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"gotest.tools/v3/assert"
)

func TestOpenUnknownFileTypeDoesNotCrash(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	session, _ := projecttestutil.Setup(map[string]any{})

	ctx := projecttestutil.WithRequestID(context.Background())
	uri := lsproto.DocumentUri("file:///component.vue")
	session.DidOpenFile(ctx, uri, 1, "let x = 1;", "vue")

	languageService, err := session.GetLanguageService(ctx, uri)
	assert.NilError(t, err)
	assert.Assert(t, languageService.GetProgram().GetSourceFile("/component.vue") != nil)
}
