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
	fieldKind := ptrTo(lsproto.CompletionItemKindField)
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
