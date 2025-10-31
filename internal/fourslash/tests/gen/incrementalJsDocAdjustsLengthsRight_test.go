package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestIncrementalJsDocAdjustsLengthsRight(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @noLib: true

/**
 * Pad ` + "`" + `str` + "`" + ` to ` + "`" + `width` + "`" + `.
 *
 * @param {String} str
 * @param {Number} wid/*1*/`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.GoToMarker(t, "1")
	f.Insert(t, "th\n@")
}
