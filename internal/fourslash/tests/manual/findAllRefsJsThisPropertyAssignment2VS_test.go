package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFindAllRefsJsThisPropertyAssignment2VS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowJs: true
// @noImplicitThis: true
// @Filename: infer.d.ts
export declare function infer(o: { m: Record<string, Function> } & ThisType<{ x: number }>): void;
// @Filename: a.js
import { infer } from "./infer";
infer({
    m: {
        initData() {
            this.x = 1;
            this./*1*/x;
        },
    }
});
// @Filename: b.ts
import { infer } from "./infer";
infer({
    m: {
        initData() {
            this.x = 1;
            this./*2*/x;
        },
    }
});`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "1", "2")
}
