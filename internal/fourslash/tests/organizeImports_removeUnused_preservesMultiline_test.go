package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestOrganizeImports_removeUnused_preservesMultiline(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `import {
    a,
    b,
    c,
} from "module";

export { a, b, c };`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyOrganizeImports(
		t,
		`import {
    a,
    b,
    c,
} from "module";

export { a, b, c };`,
		lsproto.CodeActionKindSourceRemoveUnusedImports,
		nil,
	)
}

func TestOrganizeImports_removeUnused_preservesMultilineWithRemoval(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `import {
    a,
    b,
    c,
} from "module";

export { a, c };`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyOrganizeImports(
		t,
		`import {
    a,
    c
} from "module";

export { a, c };`,
		lsproto.CodeActionKindSourceRemoveUnusedImports,
		nil,
	)
}

func TestOrganizeImports_removeUnusedUsesFormatOptions(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `import {
    a,
    b,
    c,
} from "module";

export { a, c };`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyOrganizeImportsWithFormattingOptions(
		t,
		"import {\n\ta,\n\tc\n} from \"module\";\n\nexport { a, c };",
		lsproto.CodeActionKindSourceRemoveUnusedImports,
		&lsproto.FormattingOptions{TabSize: 2, InsertSpaces: false},
	)
}
