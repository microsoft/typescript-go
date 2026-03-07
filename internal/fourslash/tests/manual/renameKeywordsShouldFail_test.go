package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestRenameKeywordsShouldFail(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @noLib: true
/*classKw*/class /*className*/Foo {
    /*publicKw*/public bar: string = "";
}
/*funcKw*/function /*funcName*/baz() {
    /*returnKw*/return 1;
}
/*constKw*/const /*constName*/x = 1;
/*ifKw*/if (x) /*openBrace*/{ }
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	// Keywords that are adjustable (shift to declaration name) should succeed
	adjustableMarkers := []string{"classKw", "funcKw", "constKw", "publicKw"}
	for _, marker := range adjustableMarkers {
		f.GoToMarker(t, marker)
		f.VerifyRenameSucceeded(t, nil /*preferences*/)
	}

	// Keywords and tokens that cannot be renamed should fail
	failMarkers := []string{"returnKw", "ifKw", "openBrace"}
	for _, marker := range failMarkers {
		f.GoToMarker(t, marker)
		f.VerifyRenameFailed(t, nil /*preferences*/)
	}

	// Identifiers should be renameable
	identifierMarkers := []string{"className", "funcName", "constName"}
	for _, marker := range identifierMarkers {
		f.GoToMarker(t, marker)
		f.VerifyRenameSucceeded(t, nil /*preferences*/)
	}
}
