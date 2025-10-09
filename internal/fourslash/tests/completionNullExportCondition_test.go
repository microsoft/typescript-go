package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// Regression test for panic when completing with null export conditions
// Issue: https://github.com/microsoft/typescript-go/issues/XXXX
func TestCompletionNullExportCondition(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @module: esnext
// @moduleResolution: node16
// @allowJs: true
// @noTypesAndSymbols: true
// @noEmit: true

// @Filename: /node_modules/dep/package.json
{
  "name": "dep",
  "version": "1.0.0",
  "exports": {
    ".": {
        "import": null,
        "require": "./dist/index.js",
        "types": "./dist/index.d.ts"
    }
  }
}

// @Filename: /node_modules/dep/dist/index.d.ts
export const sourceMapsEnabled = true;
export const someOtherExport = "value";

// @Filename: /node_modules/dep/dist/index.mjs
export const sourceMapsEnabled = true;
export const someOtherExport = "value";

// @Filename: /index.mts
// Trigger auto-import completion - should work without panicking
source/**/
`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	// This should not panic when requesting completions
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		UserPreferences: &ls.UserPreferences{
			IncludeCompletionsForModuleExports:    core.TSTrue,
			IncludeCompletionsForImportStatements: core.TSTrue,
		},
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			// The exact items don't matter - we're just testing it doesn't panic
		},
	})
}
