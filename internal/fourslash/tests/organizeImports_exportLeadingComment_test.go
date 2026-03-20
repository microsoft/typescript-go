package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestOrganizeImports_exportLeadingComment_notDuplicated(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// a
export { a } from "a";
console.log(a);`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyOrganizeImports(t,
		`// a
export { a } from "a";
console.log(a);`,
		lsproto.CodeActionKindSourceSortImports,
		nil,
	)
}

func TestOrganizeImports_exportLeadingComment_multipleComments_notDuplicated(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// a
// a
export { a } from "a";
console.log(a);`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyOrganizeImports(t,
		`// a
// a
export { a } from "a";
console.log(a);`,
		lsproto.CodeActionKindSourceSortImports,
		nil,
	)
}

func TestOrganizeImports_exportLeadingComment_secondExport_notDuplicated(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `export { a } from "a";
// b
export { b } from "b";
console.log(a, b);`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyOrganizeImports(t,
		`export { a } from "a";
// b
export { b } from "b";
console.log(a, b);`,
		lsproto.CodeActionKindSourceSortImports,
		nil,
	)
}

func TestOrganizeImports_exportLeadingComment_secondExport_withBlankLine_notDuplicated(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `export { a } from "a";

// b
export { b } from "b";
console.log(a, b);`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyOrganizeImports(t,
		`export { a } from "a";

// b
export { b } from "b";
console.log(a, b);`,
		lsproto.CodeActionKindSourceSortImports,
		nil,
	)
}

func TestOrganizeImports_exportLeadingComment_withBlankLine_notDuplicated(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// a

export { a } from "a";
console.log(a);`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyOrganizeImports(t,
		`// a

export { a } from "a";
console.log(a);`,
		lsproto.CodeActionKindSourceSortImports,
		nil,
	)
}
