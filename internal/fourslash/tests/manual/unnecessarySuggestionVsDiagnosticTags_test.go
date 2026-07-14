package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

const diagnosticTagVsHiddenInEditor = lsproto.DiagnosticTag(2147483641)

func TestUnnecessarySuggestionDiagnosticTags(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `function f([|p|]: any) {
    const [|x|] = 0;
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifySuggestionDiagnostics(t, []*lsproto.Diagnostic{
		{
			Message: "'p' is declared but its value is never read.",
			Range:   f.Ranges()[0].LSRange,
			Code:    &lsproto.IntegerOrString{Integer: new(int32(6133))},
			Tags:    &[]lsproto.DiagnosticTag{lsproto.DiagnosticTagUnnecessary},
		},
		{
			Message: "'x' is declared but its value is never read.",
			Range:   f.Ranges()[1].LSRange,
			Code:    &lsproto.IntegerOrString{Integer: new(int32(6133))},
			Tags:    &[]lsproto.DiagnosticTag{lsproto.DiagnosticTagUnnecessary},
		},
	})
}

func TestUnnecessarySuggestionDiagnosticTagsVS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `function f([|p|]: any) {
    const [|x|] = 0;
}`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifySuggestionDiagnostics(t, []*lsproto.Diagnostic{
		{
			Message: "'p' is declared but its value is never read.",
			Range:   f.Ranges()[0].LSRange,
			Code:    &lsproto.IntegerOrString{String: new("TS6133")},
			Tags:    &[]lsproto.DiagnosticTag{lsproto.DiagnosticTagUnnecessary, diagnosticTagVsHiddenInEditor},
		},
		{
			Message: "'x' is declared but its value is never read.",
			Range:   f.Ranges()[1].LSRange,
			Code:    &lsproto.IntegerOrString{String: new("TS6133")},
			Tags:    &[]lsproto.DiagnosticTag{lsproto.DiagnosticTagUnnecessary, diagnosticTagVsHiddenInEditor},
		},
	})
}
