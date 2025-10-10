package lsptestutil

import (
	"io"
	"strings"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/vfs"
	"gotest.tools/v3/assert"
)

type lspReader struct {
	c <-chan *lsproto.Message
}

func (r *lspReader) Read() (*lsproto.Message, error) {
	msg, ok := <-r.c
	if !ok {
		return nil, io.EOF
	}
	return msg, nil
}

type lspWriter struct {
	c chan<- *lsproto.Message
}

func (w *lspWriter) Write(msg *lsproto.Message) error {
	w.c <- msg
	return nil
}

func (r *lspWriter) Close() {
	close(r.c)
}

var (
	_ lsp.Reader = (*lspReader)(nil)
	_ lsp.Writer = (*lspWriter)(nil)
)

func newLSPPipe() (*lspReader, *lspWriter) {
	c := make(chan *lsproto.Message, 100)
	return &lspReader{c: c}, &lspWriter{c: c}
}

var (
	ptrTrue                       = ptrTo(true)
	defaultCompletionCapabilities = &lsproto.CompletionClientCapabilities{
		CompletionItem: &lsproto.ClientCompletionItemOptions{
			SnippetSupport:          ptrTrue,
			CommitCharactersSupport: ptrTrue,
			PreselectSupport:        ptrTrue,
			LabelDetailsSupport:     ptrTrue,
			InsertReplaceSupport:    ptrTrue,
		},
		CompletionList: &lsproto.CompletionListCapabilities{
			ItemDefaults: &[]string{"commitCharacters", "editRange"},
		},
	}
)

type TestLspServer struct {
	Server *lsp.Server
	in     *lspWriter
	out    *lspReader
	id     int32
	Vfs    vfs.FS
}

func NewTestLspServer(
	t *testing.T,
	fs vfs.FS,
	parseCache *project.ParseCache,
	optionsForInferredProject *core.CompilerOptions,
	capabilities *lsproto.ClientCapabilities,
) *TestLspServer {
	t.Helper()

	inputReader, inputWriter := newLSPPipe()
	outputReader, outputWriter := newLSPPipe()

	var err strings.Builder
	server := lsp.NewServer(&lsp.ServerOptions{
		In:  inputReader,
		Out: outputWriter,
		Err: &err,

		Cwd:                "/",
		FS:                 fs,
		DefaultLibraryPath: bundled.LibPath(),

		ParseCache: parseCache,
	})

	go func() {
		defer func() {
			outputWriter.Close()
		}()
		err := server.Run()
		if err != nil {
			t.Error("server error:", err)
		}
	}()

	s := &TestLspServer{
		Server: server,
		in:     inputWriter,
		out:    outputReader,
		Vfs:    fs,
	}

	// !!! temporary; remove when we have `handleDidChangeConfiguration`/implicit project config support
	// !!! replace with a proper request *after initialize*
	s.Server.SetCompilerOptionsForInferredProjects(t.Context(), optionsForInferredProject)
	s.initialize(t, capabilities)

	t.Cleanup(func() {
		inputWriter.Close()
	})
	return s
}

func (s *TestLspServer) nextID() int32 {
	id := s.id
	s.id++
	return id
}

func (s *TestLspServer) initialize(t *testing.T, capabilities *lsproto.ClientCapabilities) {
	params := &lsproto.InitializeParams{
		Locale: ptrTo("en-US"),
	}
	params.Capabilities = getCapabilitiesWithDefaults(capabilities)
	// !!! check for errors?
	SendRequest(t, s, lsproto.InitializeInfo, params)
	SendNotification(t, s, lsproto.InitializedInfo, &lsproto.InitializedParams{})
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

func SendRequest[Params, Resp any](t *testing.T, s *TestLspServer, info lsproto.RequestInfo[Params, Resp], params Params) (*lsproto.Message, Resp, bool) {
	id := s.nextID()
	req := lsproto.NewRequestMessage(
		info.Method,
		lsproto.NewID(lsproto.IntegerOrString{Integer: &id}),
		params,
	)
	s.writeMsg(t, req.Message())
	resp := s.readMsg(t)
	if resp == nil {
		return nil, *new(Resp), false
	}
	result, ok := resp.AsResponse().Result.(Resp)
	return resp, result, ok
}

func SendNotification[Params any](t *testing.T, s *TestLspServer, info lsproto.NotificationInfo[Params], params Params) {
	notification := lsproto.NewNotificationMessage(
		info.Method,
		params,
	)
	s.writeMsg(t, notification.Message())
}

func (s *TestLspServer) writeMsg(t *testing.T, msg *lsproto.Message) {
	assert.NilError(t, json.MarshalWrite(io.Discard, msg), "failed to encode message as JSON")
	if err := s.in.Write(msg); err != nil {
		t.Fatalf("failed to write message: %v", err)
	}
}

func (s *TestLspServer) readMsg(t *testing.T) *lsproto.Message {
	// !!! filter out response by id etc
	msg, err := s.out.Read()
	if err != nil {
		t.Fatalf("failed to read message: %v", err)
	}
	assert.NilError(t, json.MarshalWrite(io.Discard, msg), "failed to encode message as JSON")
	return msg
}

func ptrTo[T any](v T) *T {
	return &v
}
