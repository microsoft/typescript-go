package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
	"github.com/microsoft/typescript-go/internal/testutil"
)

const content = `// @Filename: /project/helper/index.ts
export const helperFunc = () => {};
// @Filename: /project/index.ts
helper/**/`

func TestImportModuleSpecifierEndingAuto(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.Configure(t, &lsutil.UserPreferences{
		IncludeCompletionsForModuleExports:    core.TSTrue,
		IncludeCompletionsForImportStatements: core.TSTrue,
		ImportModuleSpecifierEnding:           modulespecifiers.ImportModuleSpecifierEndingPreferenceAuto,
	})
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{"helperFunc"},
		},
	})
	f.BaselineAutoImportsCompletions(t, []string{""})
}

func TestImportModuleSpecifierEndingMinimal(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.Configure(t, &lsutil.UserPreferences{
		IncludeCompletionsForModuleExports:    core.TSTrue,
		IncludeCompletionsForImportStatements: core.TSTrue,
		ImportModuleSpecifierEnding:           modulespecifiers.ImportModuleSpecifierEndingPreferenceMinimal,
	})
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{"helperFunc"},
		},
	})
	f.BaselineAutoImportsCompletions(t, []string{""})
}

func TestImportModuleSpecifierEndingIndex(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.Configure(t, &lsutil.UserPreferences{
		IncludeCompletionsForModuleExports:    core.TSTrue,
		IncludeCompletionsForImportStatements: core.TSTrue,
		ImportModuleSpecifierEnding:           modulespecifiers.ImportModuleSpecifierEndingPreferenceIndex,
	})
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{"helperFunc"},
		},
	})
	f.BaselineAutoImportsCompletions(t, []string{""})
}

func TestImportModuleSpecifierEndingJs(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.Configure(t, &lsutil.UserPreferences{
		IncludeCompletionsForModuleExports:    core.TSTrue,
		IncludeCompletionsForImportStatements: core.TSTrue,
		ImportModuleSpecifierEnding:           modulespecifiers.ImportModuleSpecifierEndingPreferenceJs,
	})
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{"helperFunc"},
		},
	})
	f.BaselineAutoImportsCompletions(t, []string{""})
}
