package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/testutil"
	"github.com/microsoft/typescript-go/internal/testutil/contentmappertest"
)

func TestContentMapperRejectsAutoImportInSynthesizedSupplementalPrefix(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// The mapper produces these virtual files for app.astro:
	//
	//   canonical:    export {};
	//   supplemental: /* generated */
	//                 const value = help;
	//
	// Only the original `const value = help;` text is mapped. An auto-import would be inserted at
	// supplemental offset 0, inside the synthesized prefix, so it cannot be represented in app.astro.
	// Requesting completions exercises the full completion -> import adder -> change tracker path. The
	// completion must be omitted rather than returning an edit under the supplemental's virtual filename,
	// which would evade rejection under the physical filename and modify the original file incorrectly.
	f, done := newContentMapperFourslash(t, `// @Filename: /dep.ts
export const helper = 1;

// @Filename: /app.astro
const value = help/**/;
`, contentmappertest.PrefixedSupplementalMapper, ".astro")
	defer done()

	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		UserPreferences: &lsutil.UserPreferences{
			IncludeCompletionsForModuleExports:    core.TSTrue,
			IncludeCompletionsForImportStatements: core.TSTrue,
		},
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{Excludes: []string{"helper"}},
	})
}

func TestContentMapperDropsUnmappedFoldingRanges(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// The mapper replaces app.fold with this fully synthesized virtual TypeScript and returns an empty
	// span map:
	//
	//   import "a";
	//   import "b";
	//   /*
	//    * generated
	//    */
	//   export {};
	//
	// The import group and multiline comment both produce virtual folding candidates, but neither can be
	// mapped into the original `host markup` document. The LSP folding request must return no ranges. In
	// particular, the rejected candidates must not be appended as nil entries and later dereferenced while
	// sorting the response.
	f, done := newContentMapperFourslash(t, `// @Filename: /app.fold
host markup
`, contentmappertest.UnmappedFoldingMapper, ".fold")
	defer done()

	f.GoToFile(t, "/app.fold")
	f.VerifyFoldingRangeLines(t, nil)
}
