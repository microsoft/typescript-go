package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSemanticModernClassificationCallableVariables2(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `import "node";
var fs = require("fs")
require.resolve('react');
require.resolve.paths;
interface LanguageMode { getFoldingRanges?: (d: string) => number[]; };
function (mode: LanguageMode | undefined) { if (mode && mode.getFoldingRanges) { return mode.getFoldingRanges('a'); }};
function b(a: () => void) { a(); };`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{{Type: "variable.declaration", Text: "fs"}, {Type: "interface.declaration", Text: "LanguageMode"}, {Type: "member.declaration", Text: "getFoldingRanges"}, {Type: "parameter.declaration", Text: "d"}, {Type: "parameter.declaration", Text: "mode"}, {Type: "interface", Text: "LanguageMode"}, {Type: "parameter", Text: "mode"}, {Type: "parameter", Text: "mode"}, {Type: "member", Text: "getFoldingRanges"}, {Type: "parameter", Text: "mode"}, {Type: "member", Text: "getFoldingRanges"}, {Type: "function.declaration", Text: "b"}, {Type: "function.declaration", Text: "a"}, {Type: "function", Text: "a"}})
}
