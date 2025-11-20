package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	. "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestNavigationItemsSpecialPropertyAssignment(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @noLib: true
// @allowJs: true
// @Filename: /a.js
[|exports.x = 0|];
[|exports.y = function() {}|];
function Cls() {
    [|this.instanceProp = 0|];
}
[|Cls.staticMethod = function() {}|];
[|Cls.staticProperty = 0|];
[|Cls.prototype.instanceMethod = function() {}|];`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyWorkspaceSymbol(t, []*fourslash.VerifyWorkspaceSymbolCase{
		{
			Preferences: nil,
			Includes: []*lsproto.SymbolInformation{
				{
					Name:     "x",
					Kind:     lsproto.SymbolKindVariable,
					Location: f.Ranges()[0].LSLocation(),
				},
			},
		}, {
			Preferences: nil,
			Includes: []*lsproto.SymbolInformation{
				{
					Name:     "y",
					Kind:     lsproto.SymbolKindFunction,
					Location: f.Ranges()[1].LSLocation(),
				},
			},
		}, {
			Preferences: nil,
			Includes: []*lsproto.SymbolInformation{
				{
					Name:          "instanceProp",
					Kind:          lsproto.SymbolKindProperty,
					Location:      f.Ranges()[2].LSLocation(),
					ContainerName: PtrTo("Cls"),
				},
			},
		}, {
			Preferences: nil,
			Includes: []*lsproto.SymbolInformation{
				{
					Name:     "staticMethod",
					Kind:     lsproto.SymbolKindMethod,
					Location: f.Ranges()[3].LSLocation(),
				},
			},
		}, {
			Preferences: nil,
			Includes: []*lsproto.SymbolInformation{
				{
					Name:     "staticProperty",
					Kind:     lsproto.SymbolKindProperty,
					Location: f.Ranges()[4].LSLocation(),
				},
			},
		}, {
			Preferences: nil,
			Includes: []*lsproto.SymbolInformation{
				{
					Name:     "instanceMethod",
					Kind:     lsproto.SymbolKindMethod,
					Location: f.Ranges()[5].LSLocation(),
				},
			},
		},
	})
}
