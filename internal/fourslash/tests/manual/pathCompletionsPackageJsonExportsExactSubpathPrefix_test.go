package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// Regression test for https://github.com/microsoft/typescript-go/issues/3981
//
// When a package has only exact (non-wildcard) exports under a shared prefix
// (e.g. "./unstable/sync", "./unstable/ast") and the user types the prefix
// followed by "/", completions must name only the suffix component ("sync",
// "ast"), not repeat the whole subpath ("unstable/sync", "unstable/ast").
// Accepting a completion should yield e.g. "pkg/unstable/sync", not
// "pkg/unstable/unstable/sync".
func TestPathCompletionsPackageJsonExportsExactSubpathPrefix(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @module: node18
// @Filename: /node_modules/pkg/package.json
{
  "name": "pkg",
  "exports": {
    "./unstable/sync":  { "default": "./dist/sync.js" },
    "./unstable/async": { "default": "./dist/async.js" },
    "./unstable/ast":   { "default": "./dist/ast.js" }
  }
}
// @Filename: /node_modules/pkg/dist/sync.d.ts
export const sync = 0;
// @Filename: /node_modules/pkg/dist/async.d.ts
export const async = 0;
// @Filename: /node_modules/pkg/dist/ast.d.ts
export const ast = 0;
// @Filename: /index.mts
import { } from "pkg/unstable//**/";`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &[]string{},
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			// Completions must name only the suffix after "unstable/" — not repeat the full
			// subpath — so that accepting "sync" yields "pkg/unstable/sync", not
			// "pkg/unstable/unstable/sync".
			Includes: []fourslash.CompletionsExpectedItem{
				"sync",
				"async",
				"ast",
			},
			// These duplicated-prefix names must not appear.
			Excludes: []string{"unstable/sync", "unstable/async", "unstable/ast"},
		},
	})
}
