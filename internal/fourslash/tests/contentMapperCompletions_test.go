package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/testutil"
	"github.com/microsoft/typescript-go/internal/testutil/contentmappertest"
)

func TestContentMapperCompletions(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	f, done := newContentMapperFourslash(t, `// @Filename: /settings.ts
export const settings = { color: "blue", size: 2 };

// @Filename: /ProfileCard.vue
<component name="ProfileCard">
<template>
  <h1>{{ ti/*atom*/tle }}</h1>
  <div class="card/*markup*/">Profile</div>
</template>
<script lang="ts">
import { settings } from "./settings";
export const title = "Profile";
export const card = { title, settings };
settings.co/*outgoing*/lor;
card.ti/*exact*/tle;
</script>

// @Filename: /main.ts
import { card } from "./ProfileCard.vue";
card.se/*incoming*/ttings;
`, contentmappertest.ComponentMapper, ".vue")
	defer done()

	propertyCompletions := func(items ...fourslash.CompletionsExpectedItem) *fourslash.CompletionsExpectedList {
		return &fourslash.CompletionsExpectedList{
			ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{CommitCharacters: &DefaultCommitCharacters, EditRange: Ignored},
			Items:        &fourslash.CompletionsExpectedItems{Includes: items},
		}
	}
	f.VerifyCompletions(t, "exact", propertyCompletions("title"))
	f.VerifyCompletions(t, "outgoing", propertyCompletions("color"))
	f.VerifyCompletions(t, "incoming", propertyCompletions("settings"))
	f.VerifyCompletions(t, "atom", nil)
	f.VerifyCompletions(t, "markup", nil)
}
