package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavbar_exportDefault(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: a.ts
export default class { }
// @Filename: b.ts
export default class C { }
// @Filename: c.ts
export default function { }
// @Filename: d.ts
export default function Func { }`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.GoToFile(t, "a.ts")
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "default",
			Kind:     lsproto.SymbolKindClass,
			Children: nil,
		},
	})
	f.GoToFile(t, "b.ts")
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "C",
			Kind:     lsproto.SymbolKindClass,
			Children: nil,
		},
	})
	f.GoToFile(t, "c.ts")
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "default",
			Kind:     lsproto.SymbolKindFunction,
			Children: nil,
		},
	})
	f.GoToFile(t, "d.ts")
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "Func",
			Kind:     lsproto.SymbolKindFunction,
			Children: nil,
		},
	})
}
