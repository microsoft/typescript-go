package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormatInterfaceWithMissingBraceAndLaterTemplateString2(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = "// @Filename: /FormCheck.tsx\ninterface FormCheckProps {\n\nconst FormCheck: DynamicRefForwardingComponent<'input', FormCheckProps> =\n  React.forwardRef(\n    () => {\n      return <div className={`${bsPrefix}-reverse`} />;\n    },\n  );\n\nFormCheck.displayName = 'FormCheck';\n"
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.FormatDocument(t, "")
	f.VerifyCurrentFileContent(t, "interface FormCheckProps {\n\nconst FormCheck: DynamicRefForwardingComponent<'input', FormCheckProps> =\n    React.forwardRef(\n        () => {\n            return <div className={`${bsPrefix}-reverse`} />;\n        },\n    );\n\nFormCheck.displayName = 'FormCheck';\n")
}
