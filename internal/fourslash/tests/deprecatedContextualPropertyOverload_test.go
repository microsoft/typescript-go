package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestDeprecatedContextualPropertyOverload(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `interface DeprecatedOptions {
    kind: "deprecated";
    /** @deprecated */
    value: number;
}
interface CurrentOptions {
    kind: "current";
    value: number;
}
declare function select(options: DeprecatedOptions): void;
declare function select(options: CurrentOptions): void;

select({ kind: "current", value: 1 });

interface FirstDeprecatedOptions {
    kind: "first";
    /** @deprecated */
    value: number;
}
interface SecondDeprecatedOptions {
    kind: "second";
    /** @deprecated */
    value: number;
}
declare function selectDeprecated(options: FirstDeprecatedOptions): void;
declare function selectDeprecated(options: SecondDeprecatedOptions): void;

selectDeprecated({ kind: "second", [|value|]: 1 });`

	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.VerifySuggestionDiagnostics(t, []*lsproto.Diagnostic{
		{
			Code:    &lsproto.IntegerOrString{Integer: new(int32(6385))},
			Message: lsproto.StringOrMarkupContent{String: new("'value' is deprecated.")},
			Tags:    &[]lsproto.DiagnosticTag{lsproto.DiagnosticTagDeprecated},
			Range:   f.Ranges()[0].LSRange,
		},
	})
}
