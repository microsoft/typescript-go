package lstestutil

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

// TODO: move this to a fourslash package

type FourslashTest struct {
	server   *lsp.Server
	in       *lsproto.BaseWriter
	out      *lsproto.BaseReader
	id       int32
	testData *TestData
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
	}
	return f, done
}

func (f *FourslashTest) nextID() int32 {
	id := f.id
	f.id++
	return id
}

func (f *FourslashTest) initialize(t *testing.T, capabilities *lsproto.ClientCapabilities) {
	capabilities.General = &lsproto.GeneralClientCapabilities{
		PositionEncodings: &[]lsproto.PositionEncodingKind{lsproto.PositionEncodingKindUTF8},
	}
	// capabilities.Workspace = &lsproto.WorkspaceClientCapabilities{}
	// !!! set capabilities inline once that's allowed by the lsp types
	params := &lsproto.InitializeParams{}
	params.Capabilities = capabilities
	// !!! check for errors?
	f.sendRequest(t, lsproto.MethodInitialize, params)
	f.sendNotification(t, lsproto.MethodInitialized, &lsproto.InitializedParams{})
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

func (f *FourslashTest) VerifyCompletions(t *testing.T, markerName string, expected any) {
	marker, ok := f.testData.MarkerPositions[markerName]
	if !ok {
		t.Fatalf("Marker %s not found", markerName)
	}
	params := &lsproto.CompletionParams{
		TextDocumentPositionParams: lsproto.TextDocumentPositionParams{
			TextDocument: lsproto.TextDocumentIdentifier{
				Uri: lsproto.DocumentUri(marker.Filename),
			},
			Position: marker.LSPosition,
		},
	}
	resMsg := f.sendRequest(t, lsproto.MethodTextDocumentCompletion, params)
	if resMsg == nil {
		t.Fatalf("Nil response received for completion request at marker %s", markerName)
	}
	response := resMsg.AsResponse()
	switch response.Result.(type) {
	case *lsproto.CompletionList:
		// !!! verify completion list
		// !!! test failure should indicate which marker failed via some sort of prefix msg
		return
	default:
		t.Fatalf("Unexpected response type for completion request at marker %s: %T", markerName, response.Result)
	}
}
