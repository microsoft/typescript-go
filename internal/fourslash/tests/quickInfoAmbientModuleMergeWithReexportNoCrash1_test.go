package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestQuickInfoAmbientModuleMergeWithReexportNoCrash1(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /node_modules/foo/index.d.ts
declare function foo(): void;
declare namespace foo { export const items: string[]; }
export = foo;
// @Filename: /a.d.ts
declare module 'mymod' { import * as foo from 'foo'; export { foo }; }
// @Filename: /b.d.ts
declare module 'mymod' { export const foo: number; }
// @Filename: /index.ts
const x/*m1*/ = 1;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyQuickInfoAt(t, "m1", "const x: 1", "")
}
