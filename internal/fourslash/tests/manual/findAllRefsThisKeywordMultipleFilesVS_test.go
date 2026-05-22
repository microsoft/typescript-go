package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFindAllRefsThisKeywordMultipleFilesVS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: file1.ts
/*1*/this; /*2*/this;
// @Filename: file2.ts
/*3*/this;
/*4*/this;
// @Filename: file3.ts
 ((x = /*5*/this, y) => /*6*/this)(/*7*/this, /*8*/this);
 // different 'this'
 function f(this) { return this; }`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "1", "2", "3", "4", "5", "6", "7", "8")
}
