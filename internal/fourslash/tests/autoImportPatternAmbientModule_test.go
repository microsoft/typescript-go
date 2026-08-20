package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// Auto-imports ignores
func TestAutoImportMergedPatternAmbientModule(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /tsconfig.json
{ "compilerOptions": { "module": "preserve", "moduleResolution": "bundler" } }

// @Filename: /first.d.ts
declare module "*.asset" with { type: "css" } {
    export const styles: string;
}

// @Filename: /second.d.ts
declare module "*.asset" with { type: "css" } {
    export const styleTokens: string;
}

// @Filename: /index.ts
sty/**/`

	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Excludes: []string{"styles", "styleTokens"},
		},
	})
}
