package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarItemsInsideMethodsAndConstructors(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class Class {
    constructor() {
        function LocalFunctionInConstructor() {}
        interface LocalInterfaceInConstrcutor {}
        enum LocalEnumInConstructor { LocalEnumMemberInConstructor }
    }

    method() {
        function LocalFunctionInMethod() {
            function LocalFunctionInLocalFunctionInMethod() {}
        }
        interface LocalInterfaceInMethod {}
        enum LocalEnumInMethod { LocalEnumMemberInMethod }
    }

    emptyMethod() { } // Non child functions method should not be duplicated
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name: "Class",
			Kind: lsproto.SymbolKindClass,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name: "constructor",
					Kind: lsproto.SymbolKindConstructor,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name: "LocalEnumInConstructor",
							Kind: lsproto.SymbolKindEnum,
							Children: PtrTo([]*lsproto.DocumentSymbol{
								{
									Name:     "LocalEnumMemberInConstructor",
									Kind:     lsproto.SymbolKindEnumMember,
									Children: nil,
								},
							}),
						},
						{
							Name:     "LocalFunctionInConstructor",
							Kind:     lsproto.SymbolKindFunction,
							Children: nil,
						},
						{
							Name:     "LocalInterfaceInConstrcutor",
							Kind:     lsproto.SymbolKindInterface,
							Children: nil,
						},
					}),
				},
				{
					Name:     "emptyMethod",
					Kind:     lsproto.SymbolKindMethod,
					Children: nil,
				},
				{
					Name: "method",
					Kind: lsproto.SymbolKindMethod,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name: "LocalEnumInMethod",
							Kind: lsproto.SymbolKindEnum,
							Children: PtrTo([]*lsproto.DocumentSymbol{
								{
									Name:     "LocalEnumMemberInMethod",
									Kind:     lsproto.SymbolKindEnumMember,
									Children: nil,
								},
							}),
						},
						{
							Name: "LocalFunctionInMethod",
							Kind: lsproto.SymbolKindFunction,
							Children: PtrTo([]*lsproto.DocumentSymbol{
								{
									Name:     "LocalFunctionInLocalFunctionInMethod",
									Kind:     lsproto.SymbolKindFunction,
									Children: nil,
								},
							}),
						},
						{
							Name:     "LocalInterfaceInMethod",
							Kind:     lsproto.SymbolKindInterface,
							Children: nil,
						},
					}),
				},
			}),
		},
	})
}
