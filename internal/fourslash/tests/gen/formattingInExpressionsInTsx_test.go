package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormattingInExpressionsInTsx(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()
	t.Skip()
=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: test.tsx
import * as React from "react";
<div
    autoComplete={(function () {
return true/*1*/
    })() }
    >
</div>`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "1")
	f.Insert(t, ";")
<<<<<<< HEAD
	f.VerifyCurrentLineContent(t, `        return true;`)
=======
	f.VerifyCurrentLineContentIs(t, "        return true;")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
