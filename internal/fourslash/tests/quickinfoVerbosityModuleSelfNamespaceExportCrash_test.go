package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestQuickinfoVerbosityModuleSelfNamespaceExportCrash(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @module: esnext
// @filename: /runner.ts
export interface Runner<A> { readonly state: A }
export declare const make: <A>() => Runner<A>;
export * as Runner from "./runner";

// @filename: /test.ts
import { Runner } from "./runner";
const runner = Runner.make<string>();
runner.state/*state*/;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineHoverWithVerbosity(t, map[string][]int{"state": {5}})
}
