package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/testutil"
	"github.com/microsoft/typescript-go/internal/testutil/contentmappertest"
)

func TestContentMapperDuplicateMappings(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	f, done := newContentMapperFourslash(t, `// @Filename: /value.dup
[|val/*query*/ue|]
`, contentmappertest.DuplicateMapper, ".dup")
	defer done()

	f.VerifyQuickInfoAt(t, "query", "const value: 1", "")
	f.VerifyBaselineGoToDefinition(t, false, "query")
	f.VerifyBaselineFindAllReferences(t, "query")
}

func TestContentMapperDisabledPurposes(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	f, done := newContentMapperFourslash(t, `// @Filename: /disabled.dup
val/*query*/ue
`, contentmappertest.DuplicateMapper, ".dup")
	defer done()

	f.GoToMarker(t, "query")
	f.VerifyNotQuickInfoExists(t)
	f.VerifyBaselineGoToDefinition(t, false, "query")
	f.VerifyBaselineFindAllReferences(t, "query")
}
