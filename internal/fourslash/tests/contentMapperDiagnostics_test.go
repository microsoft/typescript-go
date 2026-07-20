package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/testutil"
	"github.com/microsoft/typescript-go/internal/testutil/contentmappertest"
)

func TestContentMapperDiagnostics(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	f, done := newContentMapperFourslash(t, `// @Filename: /ProfileCard.vue
<component name="ProfileCard">
<template>
  <h1>{{ [|missingTitle|] }}</h1>
  <p>{{ takesNumber([|title + suffix|]) }}</p>
</template>
<script lang="ts">
export const [|bad|]: number = "wrong";
const title = "Profile";
const suffix = "!";
function takesNumber(value: number) { return value; }
</script>
`, contentmappertest.ComponentMapper, ".vue")
	defer done()

	f.VerifyBaselineNonSuggestionDiagnostics(t)
}

func TestContentMapperSynthesizedDiagnostics(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	f, done := newContentMapperFourslash(t, `// @Filename: /ProfileCard.vue
<template><h1>Profile</h1></template>
`, contentmappertest.SynthesizingMapper, ".vue")
	defer done()

	f.VerifyBaselineNonSuggestionDiagnostics(t)
}
