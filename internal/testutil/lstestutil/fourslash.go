package lstestutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	err      *bytes.Buffer
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
	var err bytes.Buffer
	fs := vfstest.FromMap(testfs, true /*useCaseSensitiveFileNames*/)
	server := lsp.NewServer(&lsp.ServerOptions{
		In:  inputReader,
		Out: outputWriter,
		Err: &err,

		Cwd:                "/",
		NewLine:            core.NewLineKindLF, // TODO: verify
		FS:                 bundled.WrapFS(fs),
		DefaultLibraryPath: bundled.LibPath(),
	})

	// !!! panic recovery
	go func() {
		if err := server.Run(); err != nil && !errors.Is(err, io.EOF) {
			panic(fmt.Sprintf("server.Run() failed: %v", err))
		}
	}()

	f := &FourslashTest{
		server:   server,
		in:       lsproto.NewBaseWriter(inputWriter),
		out:      lsproto.NewBaseReader(outputReader),
		err:      &err,
		testData: &testData,
	}

	f.initialize(t, capabilities)
	// !!! global compiler options default extracted from tests

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
	// !!! set capabilities inline once that's allowed by the lsp types
	params := &lsproto.InitializeParams{}
	params.Capabilities = capabilities
	f.sendRequest(t, lsproto.MethodInitialize, params)
}

func (f *FourslashTest) sendRequest(t *testing.T, method lsproto.Method, params any) *lsproto.Message {
	id := f.nextID()
	req := lsproto.NewRequestMessage(
		method,
		lsproto.NewID(lsproto.IntegerOrString{Integer: &id}),
		params,
	)
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}
	err = f.in.Write(data)
	if err != nil {
		t.Fatalf("failed to write request: %v", err)
	}

	// !!! read error
	// !!! filter out response
	// !!! handle out of order responses etc
	data, err = f.out.Read()
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

func (f *FourslashTest) VerifyCompletions(t *testing.T, marker string, expected any) {
	// !!! completion arguments
	params := &lsproto.CompletionParams{}
	res := f.sendRequest(t, lsproto.MethodTextDocumentCompletion, params)
	if res == nil {
		// !!! handle response etc
	}
	// !!! verify result
	// !!! test failure should indicate which marker failed via some sort of prefix msg
}
