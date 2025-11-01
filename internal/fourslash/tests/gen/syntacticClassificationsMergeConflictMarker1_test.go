package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSyntacticClassificationsMergeConflictMarker1(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `<<<<<<< HEAD
"AAAA"
=======
"BBBB"
>>>>>>> Feature`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{})
}
