package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/testutil"
	"github.com/microsoft/typescript-go/internal/testutil/contentmappertest"
)

func TestContentMapperSynthesizedDocumentSymbols(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	f, done := newContentMapperFourslash(t, `// @Filename: /app.vue
component source with no direct TypeScript span/**/
`, contentmappertest.SynthesizingMapper, ".vue")
	defer done()

	f.GoToMarker(t, "")
	f.VerifyBaselineDocumentSymbol(t)
}
