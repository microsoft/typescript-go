package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFoldingRangeLineFoldingOnly(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `if (EMPTY_TAGs.has(tag)) {
  output += "/>";
} else {
  output += ">";

  if (!html && kidcount > 0) {
    //
  }
}

export function use<T>(ctx: any): T | undefined {
  //
}`
	ptrTrue := true
	capabilities := &lsproto.ClientCapabilities{
		TextDocument: &lsproto.TextDocumentClientCapabilities{
			FoldingRange: &lsproto.FoldingRangeClientCapabilities{
				LineFoldingOnly: &ptrTrue,
				FoldingRange: &lsproto.ClientFoldingRangeOptions{
					CollapsedText: &ptrTrue,
				},
			},
		},
	}
	f, done := fourslash.NewFourslash(t, capabilities, content)
	defer done()

	// With lineFoldingOnly, end lines should be adjusted so closing brackets stay visible.
	// Line 0: if (EMPTY_TAGs.has(tag)) {
	// Line 1:   output += "/>";
	// Line 2: } else {
	// Line 3:   output += ">";
	// Line 4:
	// Line 5:   if (!html && kidcount > 0) {
	// Line 6:     //
	// Line 7:   }
	// Line 8: }
	// Line 9:
	// Line 10: export function use<T>(ctx: any): T | undefined {
	// Line 11:   //
	// Line 12: }
	f.VerifyFoldingRangeLines(t, []fourslash.FoldingRangeLineExpected{
		{StartLine: 0, EndLine: 1},   // if block: end adjusted from line 2 to 1
		{StartLine: 2, EndLine: 7},   // else block: end adjusted from line 8 to 7
		{StartLine: 5, EndLine: 6},   // inner if block: end adjusted from line 7 to 6
		{StartLine: 10, EndLine: 11}, // function: end adjusted from line 12 to 11
	})
}
