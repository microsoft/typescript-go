package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestEnumComputedPropertyNameDeprecated(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = "// @strict: true\n" +
		"// @Filename: a.ts\n" +
		"enum CHAR {\n" +
		"    [|['\\t']|] = 0x09,\n" +
		"    [|['\\n']|] = 0x0A,\n" +
		"    [|[`\\r`]|] = 0x0D,\n" +
		"    'space' = 0x20,\n" +
		"}\n\n" +
		"enum NoWarning {\n" +
		"    A = 1,\n" +
		"    B = 2,\n" +
		"    \"quoted\" = 3,\n" +
		"}"
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	ranges := f.Ranges()
	if len(ranges) != 3 {
		t.Fatalf("Expected 3 ranges, got %d", len(ranges))
	}

	message := "Using a string literal as an enum member name via a computed property is deprecated. Use a simple string literal instead."
	f.VerifySuggestionDiagnostics(t, []*lsproto.Diagnostic{
		{
			Code:    &lsproto.IntegerOrString{Integer: PtrTo[int32](1550)},
			Message: message,
			Range:   ranges[0].LSRange,
			Tags:    &[]lsproto.DiagnosticTag{lsproto.DiagnosticTagDeprecated},
		},
		{
			Code:    &lsproto.IntegerOrString{Integer: PtrTo[int32](1550)},
			Message: message,
			Range:   ranges[1].LSRange,
			Tags:    &[]lsproto.DiagnosticTag{lsproto.DiagnosticTagDeprecated},
		},
		{
			Code:    &lsproto.IntegerOrString{Integer: PtrTo[int32](1550)},
			Message: message,
			Range:   ranges[2].LSRange,
			Tags:    &[]lsproto.DiagnosticTag{lsproto.DiagnosticTagDeprecated},
		},
	})
}
