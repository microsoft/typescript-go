package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestIndirectClassInstantiation(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowJs: true
// @Filename: something.js
function TestObj(){
    this.property = "value";
}
var constructor = TestObj;
var instance = new constructor();
instance./*a*/
var class2 = function() { };
class2.prototype.blah = function() { };
var inst2 = new class2();
inst2.blah/*b*/;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.GoToMarker(t, "a")
	f.VerifyCompletions(t, nil, &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &defaultCommitCharacters,
			EditRange:        ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				"property",
				&lsproto.CompletionItem{
					Label:    "blah",
					SortText: ptrTo(string(ls.SortTextJavascriptIdentifiers)),
				},
				&lsproto.CompletionItem{
					Label:    "class2",
					SortText: ptrTo(string(ls.SortTextJavascriptIdentifiers)),
				},
				&lsproto.CompletionItem{
					Label:    "constructor",
					SortText: ptrTo(string(ls.SortTextJavascriptIdentifiers)),
				},
				&lsproto.CompletionItem{
					Label:    "inst2",
					SortText: ptrTo(string(ls.SortTextJavascriptIdentifiers)),
				},
				&lsproto.CompletionItem{
					Label:    "instance",
					SortText: ptrTo(string(ls.SortTextJavascriptIdentifiers)),
				},
				&lsproto.CompletionItem{
					Label:    "prototype",
					SortText: ptrTo(string(ls.SortTextJavascriptIdentifiers)),
				},
				&lsproto.CompletionItem{
					Label:    "TestObj",
					SortText: ptrTo(string(ls.SortTextJavascriptIdentifiers)),
				},
			},
		},
	})
	f.Backspace(t, 1)
	f.GoToMarker(t, "b")
	f.VerifyQuickInfoIs(t, "(method) class2.blah(): void", "")
}
