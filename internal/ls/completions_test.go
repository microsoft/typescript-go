package ls_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"gotest.tools/v3/assert"
)

var defaultCommitCharacters = []string{".", ",", ";"}

type testCase struct {
	name     string
	content  string
	expected map[string]*lsproto.CompletionList
}

func TestCompletions(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		// Without embedding, we'd need to read all of the lib files out from disk into the MapFS.
		// Just skip this for now.
		t.Skip("bundled files are not embedded")
	}

	itemDefaults := &lsproto.CompletionItemDefaults{
		CommitCharacters: &defaultCommitCharacters,
	}
	insertTextFormatPlainText := ptrTo(lsproto.InsertTextFormatPlainText)
	sortTextLocationPriority := ptrTo(string(ls.SortTextLocationPriority))
	sortTextLocalDeclarationPriority := ptrTo(string(ls.SortTextLocalDeclarationPriority))
	sortTextDeprecatedLocationPriority := ptrTo(string(ls.DeprecateSortText(ls.SortTextLocationPriority)))
	fieldKind := ptrTo(lsproto.CompletionItemKindField)
	methodKind := ptrTo(lsproto.CompletionItemKindMethod)
	functionKind := ptrTo(lsproto.CompletionItemKindFunction)
	variableKind := ptrTo(lsproto.CompletionItemKindVariable)

	stringMembers := []*lsproto.CompletionItem{
		{Label: "charAt", Kind: methodKind, SortText: sortTextLocationPriority, InsertTextFormat: insertTextFormatPlainText},
		{Label: "charCodeAt", Kind: methodKind, SortText: sortTextLocationPriority, InsertTextFormat: insertTextFormatPlainText},
		{Label: "concat", Kind: methodKind, SortText: sortTextLocationPriority, InsertTextFormat: insertTextFormatPlainText},
		{Label: "indexOf", Kind: methodKind, SortText: sortTextLocationPriority, InsertTextFormat: insertTextFormatPlainText},
		{Label: "lastIndexOf", Kind: methodKind, SortText: sortTextLocationPriority, InsertTextFormat: insertTextFormatPlainText},
		{Label: "length", Kind: fieldKind, SortText: sortTextLocationPriority, InsertTextFormat: insertTextFormatPlainText},
		{Label: "localeCompare", Kind: methodKind, SortText: sortTextLocationPriority, InsertTextFormat: insertTextFormatPlainText},
		{Label: "match", Kind: methodKind, SortText: sortTextLocationPriority, InsertTextFormat: insertTextFormatPlainText},
		{Label: "replace", Kind: methodKind, SortText: sortTextLocationPriority, InsertTextFormat: insertTextFormatPlainText},
		{Label: "search", Kind: methodKind, SortText: sortTextLocationPriority, InsertTextFormat: insertTextFormatPlainText},
		{Label: "slice", Kind: methodKind, SortText: sortTextLocationPriority, InsertTextFormat: insertTextFormatPlainText},
		{Label: "split", Kind: methodKind, SortText: sortTextLocationPriority, InsertTextFormat: insertTextFormatPlainText},
		{Label: "substring", Kind: methodKind, SortText: sortTextLocationPriority, InsertTextFormat: insertTextFormatPlainText},
		{Label: "toLocaleLowerCase", Kind: methodKind, SortText: sortTextLocationPriority, InsertTextFormat: insertTextFormatPlainText},
		{Label: "toLocaleUpperCase", Kind: methodKind, SortText: sortTextLocationPriority, InsertTextFormat: insertTextFormatPlainText},
		{Label: "toLowerCase", Kind: methodKind, SortText: sortTextLocationPriority, InsertTextFormat: insertTextFormatPlainText},
		{Label: "toString", Kind: methodKind, SortText: sortTextLocationPriority, InsertTextFormat: insertTextFormatPlainText},
		{Label: "toUpperCase", Kind: methodKind, SortText: sortTextLocationPriority, InsertTextFormat: insertTextFormatPlainText},
		{Label: "trim", Kind: methodKind, SortText: sortTextLocationPriority, InsertTextFormat: insertTextFormatPlainText},
		{Label: "valueOf", Kind: methodKind, SortText: sortTextLocationPriority, InsertTextFormat: insertTextFormatPlainText},
		{Label: "substr", Kind: methodKind, SortText: sortTextDeprecatedLocationPriority, InsertTextFormat: insertTextFormatPlainText},
	}

	testCases := []testCase{
		{
			name: "basicInterfaceMembers",
			content: `export {};
interface Point {
    x: number;
    y: number;
}
declare const p: Point;
p./*a*/`,
			expected: map[string]*lsproto.CompletionList{
				"a": {
					IsIncomplete: false,
					ItemDefaults: itemDefaults,
					Items: []*lsproto.CompletionItem{
						{
							Label:            "x",
							Kind:             fieldKind,
							SortText:         sortTextLocationPriority,
							InsertTextFormat: insertTextFormatPlainText,
						},
						{
							Label:            "y",
							Kind:             fieldKind,
							SortText:         sortTextLocationPriority,
							InsertTextFormat: insertTextFormatPlainText,
						},
					},
				},
			},
		},
		{
			name: "objectLiteralType",
			content: `export {};
let x = { foo: 123 };
x./*a*/`,
			expected: map[string]*lsproto.CompletionList{
				"a": {
					IsIncomplete: false,
					ItemDefaults: itemDefaults,
					Items: []*lsproto.CompletionItem{
						{
							Label:            "foo",
							Kind:             fieldKind,
							SortText:         sortTextLocationPriority,
							InsertTextFormat: ptrTo(lsproto.InsertTextFormatPlainText),
						},
					},
				},
			},
		},
		{
			name: "basicClassMembers",
			content: `
class n {
    constructor (public x: number, public y: number, private z: string) { }
}
var t = new n(0, 1, '');t./*a*/`,
			expected: map[string]*lsproto.CompletionList{
				"a": {
					IsIncomplete: false,
					ItemDefaults: itemDefaults,
					Items: []*lsproto.CompletionItem{
						{
							Label:            "x",
							Kind:             fieldKind,
							SortText:         sortTextLocationPriority,
							InsertTextFormat: insertTextFormatPlainText,
						},
						{
							Label:            "y",
							Kind:             fieldKind,
							SortText:         sortTextLocationPriority,
							InsertTextFormat: insertTextFormatPlainText,
						},
					},
				},
			},
		},
		{
			name: "cloduleAsBaseClass",
			content: `
class A {
    constructor(x: number) { }
    foo() { }
    static bar() { }
}

module A {
    export var x = 1;
    export function baz() { }
}

class D extends A {
    constructor() {
        super(1);
    }
    foo2() { }
    static bar2() { }
}

D./*a*/`,
			expected: map[string]*lsproto.CompletionList{
				"a": {
					IsIncomplete: false,
					ItemDefaults: itemDefaults,
					Items: []*lsproto.CompletionItem{ // !!! `funcionMembersPlus`
						{
							Label:            "bar",
							Kind:             methodKind,
							SortText:         sortTextLocalDeclarationPriority,
							InsertTextFormat: insertTextFormatPlainText,
						},
						{
							Label:            "bar2",
							Kind:             methodKind,
							SortText:         sortTextLocalDeclarationPriority,
							InsertTextFormat: insertTextFormatPlainText,
						},
						{
							Label:            "apply",
							Kind:             methodKind,
							SortText:         sortTextLocationPriority,
							InsertTextFormat: insertTextFormatPlainText,
						},
						{
							Label:            "arguments",
							Kind:             fieldKind,
							SortText:         sortTextLocationPriority,
							InsertTextFormat: insertTextFormatPlainText,
						},
						{
							Label:            "baz",
							Kind:             functionKind,
							SortText:         sortTextLocationPriority,
							InsertTextFormat: insertTextFormatPlainText,
						},
						{
							Label:            "bind",
							Kind:             methodKind,
							SortText:         sortTextLocationPriority,
							InsertTextFormat: insertTextFormatPlainText,
						},
						{
							Label:            "call",
							Kind:             methodKind,
							SortText:         sortTextLocationPriority,
							InsertTextFormat: insertTextFormatPlainText,
						},
						{
							Label:            "caller",
							Kind:             fieldKind,
							SortText:         sortTextLocationPriority,
							InsertTextFormat: insertTextFormatPlainText,
						},
						{
							Label:            "length",
							Kind:             fieldKind,
							SortText:         sortTextLocationPriority,
							InsertTextFormat: insertTextFormatPlainText,
						},
						{
							Label:            "prototype",
							Kind:             fieldKind,
							SortText:         sortTextLocationPriority,
							InsertTextFormat: insertTextFormatPlainText,
						},
						{
							Label:            "toString",
							Kind:             methodKind,
							SortText:         sortTextLocationPriority,
							InsertTextFormat: insertTextFormatPlainText,
						},
						{
							Label:            "x",
							Kind:             variableKind,
							SortText:         sortTextLocationPriority,
							InsertTextFormat: insertTextFormatPlainText,
						},
					},
				},
			},
		},
		{
			name: "forwardReference",
			content: `function f() {
    var x = new t();
    x./*a*/
}
class t {
    public n: number;
}`,
			expected: map[string]*lsproto.CompletionList{
				"a": {
					IsIncomplete: false,
					ItemDefaults: itemDefaults,
					Items: []*lsproto.CompletionItem{
						{
							Label:            "n",
							Kind:             fieldKind,
							SortText:         sortTextLocationPriority,
							InsertTextFormat: insertTextFormatPlainText,
						},
					},
				},
			},
		},
		{
			name: "lambdaThisMembers",
			content: `class Foo {
    a: number;
    b() {
        var x = () => {
            this./**/;
        }
    }
}`,
			expected: map[string]*lsproto.CompletionList{
				"": {
					IsIncomplete: false,
					ItemDefaults: itemDefaults,
					Items: []*lsproto.CompletionItem{
						{
							Label:            "a",
							Kind:             fieldKind,
							SortText:         sortTextLocationPriority,
							InsertTextFormat: insertTextFormatPlainText,
						},
						{
							Label:            "b",
							Kind:             methodKind,
							SortText:         sortTextLocationPriority,
							InsertTextFormat: insertTextFormatPlainText,
						},
					},
				},
			},
		},
		{
			name: "memberCompletionInForEach1",
			content: `var x: string[] = [];
x.forEach(function (y) { y./*1*/`,
			expected: map[string]*lsproto.CompletionList{
				"1": {
					IsIncomplete: false,
					ItemDefaults: itemDefaults,
					Items:        stringMembers,
				},
			},
		},
		{
			name: "completionsTuple",
			content: `declare const x: [number, number];
x./**/;`,
			expected: map[string]*lsproto.CompletionList{
				"": {
					IsIncomplete: false,
					ItemDefaults: itemDefaults,
					Items: []*lsproto.CompletionItem{
						{
							Label:            "0",
							Kind:             fieldKind,
							SortText:         sortTextLocationPriority,
							InsertText:       ptrTo("[0]"),
							InsertTextFormat: insertTextFormatPlainText,
							TextEdit: &lsproto.TextEditOrInsertReplaceEdit{
								TextEdit: &lsproto.TextEdit{
									NewText: "[0]",
									Range:   lsproto.Range{},
								},
							},
						},
						{
							Label:            "1",
							Kind:             fieldKind,
							SortText:         sortTextLocationPriority,
							InsertText:       ptrTo("[1]"),
							InsertTextFormat: insertTextFormatPlainText,
							TextEdit: &lsproto.TextEditOrInsertReplaceEdit{
								TextEdit: &lsproto.TextEdit{
									NewText: "[1]",
									Range:   lsproto.Range{},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			runTest(t, testCase.content, testCase.expected)
		})
	}
}

func runTest(t *testing.T, content string, expected map[string]*lsproto.CompletionList) {
	testData := ls.ParseTestdata("/index.ts", content, "/index.ts")
	files := map[string]string{
		"/index.ts": testData.Files[0].Content,
	}
	languageService := createLanguageService("/index.ts", files)
	context := &lsproto.CompletionContext{
		TriggerKind: lsproto.CompletionTriggerKindInvoked,
	}
	capabilities := &lsproto.CompletionClientCapabilities{
		CompletionItem: &lsproto.ClientCompletionItemOptions{
			SnippetSupport:          ptrTo(true),
			CommitCharactersSupport: ptrTo(true),
			PreselectSupport:        ptrTo(true),
			LabelDetailsSupport:     ptrTo(true),
		},
		CompletionList: &lsproto.CompletionListCapabilities{
			ItemDefaults: &[]string{"commitCharacters"},
		},
	}
	preferences := &ls.UserPreferences{}

	for markerName, expectedResult := range expected {
		marker, ok := testData.MarkerPositions[markerName]
		if !ok {
			t.Fatalf("No marker found for '%s'", markerName)
		}
		completionList := languageService.ProvideCompletion(
			"/index.ts",
			marker.Position,
			context,
			capabilities,
			preferences)
		assert.DeepEqual(t, completionList, expectedResult)
	}
}

func createLanguageService(fileName string, files map[string]string) *ls.LanguageService {
	projectService, _ := projecttestutil.Setup(files)
	projectService.OpenFile(fileName, files[fileName], core.ScriptKindTS, "")
	project := projectService.Projects()[0]
	return project.LanguageService()
}

func ptrTo[T any](v T) *T {
	return &v
}
