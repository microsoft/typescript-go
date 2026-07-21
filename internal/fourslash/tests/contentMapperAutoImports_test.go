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

func TestContentMapperAutoImports(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	f, done := newContentMapperFourslash(t, `// @Filename: /ProfileCard.vue
<component name="ProfileCard">
<script lang="ts">
export const profileTitle = "Profile";
</script>

// @Filename: /main.ts
profileTi/**/
`, contentmappertest.ComponentMapper, ".vue")
	defer done()

	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		UserPreferences: &lsutil.UserPreferences{
			IncludeCompletionsForModuleExports:    core.TSTrue,
			IncludeCompletionsForImportStatements: core.TSTrue,
		},
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{"profileTitle"},
		},
	})
	f.BaselineAutoImportsCompletions(t, []string{""})
}

func TestContentMapperAutoImportsIntoMappedFile(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	f, done := newContentMapperFourslash(t, `// @Filename: /dep.ts
export const existing = 1;
export const helper = 2;

// @Filename: /ProfileCard.vue
<component name="ProfileCard">
<script lang="ts">
import { existing } from "./dep";
export const profileTitle = help/**/;
</script>
`, contentmappertest.ComponentMapper, ".vue")
	defer done()

	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		UserPreferences: &lsutil.UserPreferences{
			IncludeCompletionsForModuleExports:    core.TSTrue,
			IncludeCompletionsForImportStatements: core.TSTrue,
		},
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{"helper"},
		},
	})
	f.BaselineAutoImportsCompletions(t, []string{""})
}

func TestContentMapperNodeModulesAutoImports(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	f, done := newContentMapperFourslash(t, `// @Filename: /package.json
{ "dependencies": { "profile-package": "1.0.0" } }

// @Filename: /node_modules/profile-package/package.json
{ "name": "profile-package", "version": "1.0.0" }

// @Filename: /node_modules/profile-package/ProfileCard.vue
<component name="ProfileCard">
<script lang="ts">
export const profileTitle = "Profile";
</script>

// @Filename: /node_modules/profile-package/HiddenCard.vue
<component name="HiddenCard">
<script lang="ts">
export const hiddenTitle = "Hidden";
</script>

// @Filename: /load.ts
import "profile-package/ProfileCard.vue";

// @Filename: /main.ts
profileTi/**/
`, contentmappertest.ComponentMapper, ".vue")
	defer done()

	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		UserPreferences: &lsutil.UserPreferences{
			IncludeCompletionsForModuleExports:    core.TSTrue,
			IncludeCompletionsForImportStatements: core.TSTrue,
		},
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{"profileTitle"},
			Excludes: []string{"hiddenTitle"},
		},
	})
	f.BaselineAutoImportsCompletions(t, []string{""})
}
