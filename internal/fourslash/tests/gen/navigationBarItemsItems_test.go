package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationBarItemsItems(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// Interface
interface IPoint {
    getDist(): number;
    new(): IPoint;
    (): any;
    [x:string]: number;
    prop: string;
}

/// Module
module Shapes {

    // Class
    export class Point implements IPoint {
        constructor (public x: number, public y: number) { }

        // Instance member
        getDist() { return Math.sqrt(this.x * this.x + this.y * this.y); }

        // Getter
        get value(): number { return 0; }

        // Setter
        set value(newValue: number) { return; }

        // Static member
        static origin = new Point(0, 0);

        // Static method
        private static getOrigin() { return Point.origin; }
    }

    enum Values { value1, value2, value3 }
}

// Local variables
var p: IPoint = new Shapes.Point(3, 4);
var dist = p.getDist();`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyStradaDocumentSymbol(t, []*lsproto.DocumentSymbol{
		{
			Name:     "dist",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name: "IPoint",
			Kind: lsproto.SymbolKindInterface,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name:     "()",
					Kind:     lsproto.SymbolKindFunction,
					Children: nil,
				},
				{
					Name:     "new()",
					Kind:     lsproto.SymbolKindConstructor,
					Children: nil,
				},
				{
					Name:     "[]",
					Kind:     lsproto.SymbolKindFunction,
					Children: nil,
				},
				{
					Name:     "getDist",
					Kind:     lsproto.SymbolKindMethod,
					Children: nil,
				},
				{
					Name:     "prop",
					Kind:     lsproto.SymbolKindProperty,
					Children: nil,
				},
			}),
		},
		{
			Name:     "p",
			Kind:     lsproto.SymbolKindVariable,
			Children: nil,
		},
		{
			Name: "Shapes",
			Kind: lsproto.SymbolKindModule,
			Children: PtrTo([]*lsproto.DocumentSymbol{
				{
					Name: "Point",
					Kind: lsproto.SymbolKindClass,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "constructor",
							Kind:     lsproto.SymbolKindConstructor,
							Children: nil,
						},
						{
							Name:     "getDist",
							Kind:     lsproto.SymbolKindMethod,
							Children: nil,
						},
						{
							Name:     "getOrigin",
							Kind:     lsproto.SymbolKindMethod,
							Children: nil,
						},
						{
							Name:     "origin",
							Kind:     lsproto.SymbolKindProperty,
							Children: nil,
						},
						{
							Name:     "value",
							Kind:     lsproto.SymbolKindProperty,
							Children: nil,
						},
						{
							Name:     "value",
							Kind:     lsproto.SymbolKindProperty,
							Children: nil,
						},
						{
							Name:     "x",
							Kind:     lsproto.SymbolKindProperty,
							Children: nil,
						},
						{
							Name:     "y",
							Kind:     lsproto.SymbolKindProperty,
							Children: nil,
						},
					}),
				},
				{
					Name: "Values",
					Kind: lsproto.SymbolKindEnum,
					Children: PtrTo([]*lsproto.DocumentSymbol{
						{
							Name:     "value1",
							Kind:     lsproto.SymbolKindEnumMember,
							Children: nil,
						},
						{
							Name:     "value2",
							Kind:     lsproto.SymbolKindEnumMember,
							Children: nil,
						},
						{
							Name:     "value3",
							Kind:     lsproto.SymbolKindEnumMember,
							Children: nil,
						},
					}),
				},
			}),
		},
	})
}
