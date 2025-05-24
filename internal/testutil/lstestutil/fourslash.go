package lstestutil

import (
	"bytes"
	"encoding/json"
	"errors"
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
func NewFourslash(t *testing.T, content string, fileName string) *FourslashTest {
	rootDir := "/"
	testfs := make(map[string]string)
	testData := ParseTestData(t, content, fileName)
	for _, file := range testData.Files {
		filePath := tspath.GetNormalizedAbsolutePath(file.Filename, rootDir)
		testfs[filePath] = file.Content
	}
	var in, out, err bytes.Buffer
	fs := vfstest.FromMap(testfs, true /*useCaseSensitiveFileNames*/)
	server := lsp.NewServer(&lsp.ServerOptions{
		In:  &in,
		Out: &out,
		Err: &err,

		Cwd:                "/",
		NewLine:            core.NewLineKindLF, // TODO: verify
		FS:                 bundled.WrapFS(fs),
		DefaultLibraryPath: bundled.LibPath(),
	})

	// !!! panic recovery
	go func() {
		if err := server.Run(); err != nil && !errors.Is(err, io.EOF) {
			// !!! do something with the error
			// t.Fatalf("server.Run() failed: %v", err)
		}
	}()
	// !!! send initialize request to server
	// !!! receive initialize response
	// !!! receive file watching stuff?
	// !!! global compiler options default extracted from tests

	// !!! return cleanup function that closes the server
	return &FourslashTest{
		server:   server,
		in:       lsproto.NewBaseWriter(&in),
		out:      lsproto.NewBaseReader(&out),
		err:      &err,
		testData: &testData,
	}
}

func (f *FourslashTest) nextID() int32 {
	id := f.id
	f.id++
	return id
}

func (f *FourslashTest) VerifyCompletions(t *testing.T, marker string, expected any) {
	// !!! completion arguments
	params := &lsproto.CompletionParams{}
	id := f.nextID()
	req := lsproto.NewRequestMessage(
		lsproto.MethodTextDocumentCompletion,
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
	data, err = f.out.Read()
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}
	res := &lsproto.Message{}
	err = json.Unmarshal(data, res)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// !!! verify result
	// !!! test failure should indicate which marker failed via some sort of prefix msg
}
