package lstestutil

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

// !!! move this to a fourslash package

type FourslashTest struct {
	server *lsp.Server
	in     *lsproto.BaseWriter
	out    *lsproto.BaseReader
	id     int32

	testData *TestData

	currentCaretPosition lsproto.Position
	currentFilename      string
	lastKnownMarkerName  string
	activeFilename       string
	// !!! markers
	// !!! ranges
	// !!! files
}

// !!! automatically get fileName from test somehow?
func NewFourslash(t *testing.T, capabilities *lsproto.ClientCapabilities, content string, fileName string) (*FourslashTest, func()) {
	rootDir := "/"
	testfs := make(map[string]string)
	testData := ParseTestData(t, content, fileName)
	for _, file := range testData.Files {
		filePath := tspath.GetNormalizedAbsolutePath(file.Filename, rootDir)
		testfs[filePath] = file.Content
	}
	inputReader, inputWriter := io.Pipe()
	outputReader, outputWriter := io.Pipe()
	fs := vfstest.FromMap(testfs, true /*useCaseSensitiveFileNames*/)
	server := lsp.NewServer(&lsp.ServerOptions{
		In:  inputReader,
		Out: outputWriter,
		Err: os.Stderr,

		Cwd:                "/",
		NewLine:            core.NewLineKindLF, // !!! verify
		FS:                 bundled.WrapFS(fs),
		DefaultLibraryPath: bundled.LibPath(),
	})

	go func() {
		var err error
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("server panicked: %v", r)
			}
			inputReader.CloseWithError(err)
			outputWriter.CloseWithError(err)
		}()
		err = server.Run()
	}()

	f := &FourslashTest{
		server:   server,
		in:       lsproto.NewBaseWriter(inputWriter),
		out:      lsproto.NewBaseReader(outputReader),
		testData: &testData,
	}

	// !!! global compiler options default extracted from tests
	f.initialize(t, capabilities)

	done := func() {
		inputWriter.Close()
		outputReader.Close()
	}
	return f, done
}

func (f *FourslashTest) nextID() int32 {
	id := f.id
	f.id++
	return id
}

func (f *FourslashTest) initialize(t *testing.T, capabilities *lsproto.ClientCapabilities) {
	params := &lsproto.InitializeParams{}
	params.Capabilities = getCapabilitiesWithDefaults(capabilities)
	// !!! check for errors?
	f.sendRequest(t, lsproto.MethodInitialize, params)
	f.sendNotification(t, lsproto.MethodInitialized, &lsproto.InitializedParams{})
}

var ptrTrue = PtrTo(true)
var defaultCompletionCapabilities = &lsproto.CompletionClientCapabilities{
	CompletionItem: &lsproto.ClientCompletionItemOptions{
		SnippetSupport:          ptrTrue,
		CommitCharactersSupport: ptrTrue,
		PreselectSupport:        ptrTrue,
		LabelDetailsSupport:     ptrTrue,
		InsertReplaceSupport:    ptrTrue,
	},
	CompletionList: &lsproto.CompletionListCapabilities{
		ItemDefaults: &[]string{"commitCharacters"},
	},
}

func getCapabilitiesWithDefaults(capabilities *lsproto.ClientCapabilities) *lsproto.ClientCapabilities {
	var capabilitiesWithDefaults lsproto.ClientCapabilities
	if capabilities != nil {
		capabilitiesWithDefaults = *capabilities
	}
	capabilitiesWithDefaults.General = &lsproto.GeneralClientCapabilities{
		PositionEncodings: &[]lsproto.PositionEncodingKind{lsproto.PositionEncodingKindUTF8},
	}
	if capabilitiesWithDefaults.TextDocument == nil {
		capabilitiesWithDefaults.TextDocument = &lsproto.TextDocumentClientCapabilities{}
	}
	if capabilitiesWithDefaults.TextDocument.Completion == nil {
		capabilitiesWithDefaults.TextDocument.Completion = defaultCompletionCapabilities
	}
	return &capabilitiesWithDefaults
}

func (f *FourslashTest) sendRequest(t *testing.T, method lsproto.Method, params any) *lsproto.Message {
	id := f.nextID()
	req := lsproto.NewRequestMessage(
		method,
		lsproto.NewID(lsproto.IntegerOrString{Integer: &id}),
		params,
	)
	f.writeMsg(t, req)
	return f.readMsg(t)
}

func (f *FourslashTest) sendNotification(t *testing.T, method lsproto.Method, params any) {
	notification := lsproto.NewNotificationMessage(
		method,
		params,
	)
	f.writeMsg(t, notification)
}

func (f *FourslashTest) writeMsg(t *testing.T, msg any) {
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}
	err = f.in.Write(data)
	if err != nil {
		t.Fatalf("failed to write message: %v", err)
	}
}

func (f *FourslashTest) readMsg(t *testing.T) *lsproto.Message {
	// !!! filter out response by id
	data, err := f.out.Read()
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}
	res := &lsproto.Message{}
	err = json.Unmarshal(data, res)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	return res
}

// !!! unsorted completions? only used in 47 tests
type VerifyCompletionsResult struct {
	Includes []*lsproto.CompletionItem
	Excludes []string
	Exact    *lsproto.CompletionList
}

// !!! user preferences param
// !!! completion context param
// !!! go to marker: use current marker if none specified
func (f *FourslashTest) VerifyCompletions(t *testing.T, markerName string, expected VerifyCompletionsResult) {
	f.GoToMarker(t, markerName)
	f.verifyCompletionsWorker(t, expected)
}

func (f *FourslashTest) GoToMarker(t *testing.T, markerName string) {
	marker, ok := f.testData.MarkerPositions[markerName]
	if !ok {
		t.Fatalf("Marker %s not found", markerName)
	}
	f.ensureActiveFile(t, marker.Filename)
	f.currentCaretPosition = marker.LSPosition
	f.currentFilename = marker.Filename
	f.lastKnownMarkerName = marker.Name
}

func (f *FourslashTest) ensureActiveFile(t *testing.T, filename string) {
	if f.activeFilename != filename {
		file := core.Find(f.testData.Files, func(f *TestFileInfo) bool {
			return f.Filename == filename
		})
		if file == nil {
			t.Fatalf("File %s not found in test data", filename)
		}
		f.openFile(t, file)
	}
}

func (f *FourslashTest) openFile(t *testing.T, file *TestFileInfo) {
	// !!! normalize file path?
	f.activeFilename = file.Filename
	f.sendNotification(t, lsproto.MethodTextDocumentDidOpen, &lsproto.DidOpenTextDocumentParams{
		TextDocument: &lsproto.TextDocumentItem{
			Uri:        ls.FileNameToDocumentURI(file.Filename),
			LanguageId: getLanguageKind(file.Filename),
			Text:       file.Content,
		},
	})
}

func getLanguageKind(filename string) lsproto.LanguageKind {
	if tspath.FileExtensionIsOneOf(
		filename,
		[]string{tspath.ExtensionTs, tspath.ExtensionMts, tspath.ExtensionCts,
			tspath.ExtensionDmts, tspath.ExtensionDcts, tspath.ExtensionDts}) {
		return lsproto.LanguageKindTypeScript
	}
	if tspath.FileExtensionIsOneOf(filename, []string{tspath.ExtensionJs, tspath.ExtensionMjs, tspath.ExtensionCjs}) {
		return lsproto.LanguageKindJavaScript
	}
	if tspath.FileExtensionIs(filename, tspath.ExtensionJsx) {
		return lsproto.LanguageKindJavaScriptReact
	}
	if tspath.FileExtensionIs(filename, tspath.ExtensionTsx) {
		return lsproto.LanguageKindTypeScriptReact
	}
	if tspath.FileExtensionIs(filename, tspath.ExtensionJson) {
		return lsproto.LanguageKindJSON
	}
	return lsproto.LanguageKindTypeScript // !!! should we error in this case?
}

func (f *FourslashTest) verifyCompletionsWorker(t *testing.T, expected VerifyCompletionsResult) {
	params := &lsproto.CompletionParams{
		TextDocumentPositionParams: lsproto.TextDocumentPositionParams{
			TextDocument: lsproto.TextDocumentIdentifier{
				Uri: ls.FileNameToDocumentURI(f.currentFilename),
			},
			Position: f.currentCaretPosition,
		},
		Context: &lsproto.CompletionContext{},
	}
	resMsg := f.sendRequest(t, lsproto.MethodTextDocumentCompletion, params)
	if resMsg == nil {
		t.Fatalf("Nil response received for completion request at marker %s", f.lastKnownMarkerName)
	}
	result := resMsg.AsResponse().Result
	list := &lsproto.CompletionList{}
	err := fromDataToLsp(result, list)
	if err != nil {
		t.Fatalf("Unexpected response for completion request at marker %s: %v", f.lastKnownMarkerName, result)
	}
	verifyCompletionsResult(t, f.lastKnownMarkerName, list, expected)
}

func verifyCompletionsResult(t *testing.T, markerName string, actual *lsproto.CompletionList, expected VerifyCompletionsResult) {
	prefix := fmt.Sprintf("At marker '%s': ", markerName)
	if expected.Exact != nil {
		if expected.Includes != nil {
			t.Fatal(prefix + "Expected exact completion list but also specified 'includes'.")
		}
		if expected.Excludes != nil {
			t.Fatal(prefix + "Expected exact completion list but also specified 'excludes'.")
		}
		assertDeepEqual(t, actual, expected.Exact, prefix+"Exact completion list mismatch")
		return
	}
	nameToActualItem := make(map[string]*lsproto.CompletionItem)
	if actual != nil {
		for _, item := range actual.Items {
			nameToActualItem[item.Label] = item
		}
	}
	if expected.Includes != nil {
		for _, item := range expected.Includes {
			actualItem, ok := nameToActualItem[item.Label]
			if !ok {
				t.Fatalf("%sLabel %s not found in actual items. Actual items: %v", prefix, item.Label, actual.Items)
			}
			assertDeepEqual(t, actualItem, item, prefix+"Includes completion item mismatch for label "+item.Label)
		}
	}
	for _, exclude := range expected.Excludes {
		if _, ok := nameToActualItem[exclude]; ok {
			t.Fatalf("%sLabel %s should not be in actual items but was found. Actual items: %v", prefix, exclude, actual.Items)
		}
	}
}

// Converts from a generic JSON data structure to a specific LSP type.
func fromDataToLsp(jsonData any, result any) error {
	bytes, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, result)
}

// !!! don't compare properties that are not set in the expected value
func assertDeepEqual(t *testing.T, actual any, expected any, prefix string) {
	t.Helper()

	diff := cmp.Diff(actual, expected)
	if diff != "" {
		t.Fatalf("%s:\n%s", prefix, diff)
	}
}
