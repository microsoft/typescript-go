package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestAutoImportCompletion1(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: a.ts
export const someVar = 10;

// @Filename: b.ts
export const anotherVar = 10;

// @Filename: c.ts
import {someVar} from "./a.ts";
someVar;
a/**/
`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		UserPreferences: &ls.UserPreferences{
			IncludeCompletionsForModuleExports:    PtrTo(true),
			IncludeCompletionsForImportStatements: PtrTo(true),
			AllowIncompleteCompletions:            PtrTo(true),
		},
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{"someVar", "anotherVar"},
		},
	})
	f.BaselineAutoImportsCompletions(t, "")
}

func TestAutoImportCompletion2(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: a.ts
export const someVar = 10;
export const anotherVar = 10;

// @Filename: c.ts
import {someVar} from "./a.ts";
someVar;
a/**/
`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		UserPreferences: &ls.UserPreferences{
			IncludeCompletionsForModuleExports:    PtrTo(true),
			IncludeCompletionsForImportStatements: PtrTo(true),
			AllowIncompleteCompletions:            PtrTo(true),
		},
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{"someVar", "anotherVar"},
		},
	})
	f.BaselineAutoImportsCompletions(t, "")
}

func TestAutoImportCompletion3(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: a.ts
export const aa = "asdf";
export const someVar = 10;
export const bb = 10;

// @Filename: c.ts
import { aa, someVar } from "./a.ts";
someVar;
b/**/
`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		UserPreferences: &ls.UserPreferences{
			IncludeCompletionsForModuleExports:    PtrTo(true),
			IncludeCompletionsForImportStatements: PtrTo(true),
			AllowIncompleteCompletions:            PtrTo(true),
		},
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{"bb"},
		},
	})
	f.BaselineAutoImportsCompletions(t, "")
}
