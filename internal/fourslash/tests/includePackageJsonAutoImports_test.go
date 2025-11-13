package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestIncludePackageJsonImportsAutoImports(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /project/package.json
{
  "imports": {
    "#foo": "./foo.ts"
  }
}
// @Filename: /project/tsconfig.json
{
  "compilerOptions": {
    "baseUrl": "."
  }
}
// @Filename: /project/foo.ts
export const foo = 123;
// @Filename: /project/index.ts
foo/**/`

	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.Configure(t, &lsutil.UserPreferences{
		IncludeCompletionsForModuleExports:    core.TSTrue,
		IncludeCompletionsForImportStatements: core.TSTrue,
		IncludePackageJsonAutoImports:         lsutil.IncludePackageJsonAutoImportsOn,
	})
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{"foo"},
		},
	})
	f.BaselineAutoImportsCompletions(t, []string{""})
}
