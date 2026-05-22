package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestReferencesInEmptyFileWithMultipleProjectsVS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /home/src/workspaces/project/a/tsconfig.json
{ "files": ["a.ts"], "compilerOptions": { "lib": ["es5"] } }
// @Filename: /home/src/workspaces/project/a/a.ts
/// <reference path="../b/b.ts" />
/*1*/;
// @Filename: /home/src/workspaces/project/b/tsconfig.json
{ "files": ["b.ts"], "compilerOptions": { "lib": ["es5"] } }
// @Filename: /home/src/workspaces/project/b/b.ts
/*2*/;`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "1", "2")
}
