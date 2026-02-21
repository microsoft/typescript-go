package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// Regression test: tsgo formatter should move the closing brace of a multiline
// import to its own line, matching TypeScript's formatter behavior.
func TestFormatMultilineImportBrace(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `import {
	basename,
	extname, joinPath } from '../base/resources.js';/*0*/
import { URI } from '../base/uri.js';/*1*/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.FormatDocument(t, "")
	f.GoToMarker(t, "0")
	f.VerifyCurrentLineContent(t, `} from '../base/resources.js';`)
	f.GoToMarker(t, "1")
	f.VerifyCurrentLineContent(t, `import { URI } from '../base/uri.js';`)
}
