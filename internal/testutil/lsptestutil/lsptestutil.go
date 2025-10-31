package lsptestutil

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/project/ata"
	"github.com/microsoft/typescript-go/internal/project/logging"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
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
			DocumentationFormat:     &[]lsproto.MarkupKind{lsproto.MarkupKindMarkdown, lsproto.MarkupKindPlainText},
		},
		CompletionList: &lsproto.CompletionListCapabilities{
			ItemDefaults: &[]string{"commitCharacters", "editRange"},
		},
	}
	defaultDefinitionCapabilities = &lsproto.DefinitionClientCapabilities{
		LinkSupport: ptrTrue,
	}
	defaultTypeDefinitionCapabilities = &lsproto.TypeDefinitionClientCapabilities{
		LinkSupport: ptrTrue,
	}
	defaultHoverCapabilities = &lsproto.HoverClientCapabilities{
		ContentFormat: &[]lsproto.MarkupKind{lsproto.MarkupKindMarkdown, lsproto.MarkupKindPlainText},
	}
)

type TestLspServer struct {
	Server          *lsp.Server
	in              *lspWriter
	out             *lspReader
	id              int32
	FS              vfs.FS
	UserPreferences *lsutil.UserPreferences
	TypingsLocation string
}

type TestLspServerOptions struct {
	FS                        vfs.FS
	Client                    project.Client
	Logger                    logging.Logger
	NpmExecutor               ata.NpmExecutor
	ParseCache                *project.ParseCache
	OptionsForInferredProject *core.CompilerOptions
	TypingsLocation           string
	Capabilities              *lsproto.ClientCapabilities
}

func NewTestLspServer(t *testing.T, options *TestLspServerOptions) *TestLspServer {
	t.Helper()

	inputReader, inputWriter := newLSPPipe()
	outputReader, outputWriter := newLSPPipe()

	var npmInstall func(cwd string, args []string) ([]byte, error)
	if options.NpmExecutor != nil {
		npmInstall = func(cwd string, args []string) ([]byte, error) {
			return options.NpmExecutor.NpmInstall(cwd, args)
		}
	}

	var err strings.Builder
	server := lsp.NewServer(&lsp.ServerOptions{
		In:  inputReader,
		Out: outputWriter,
		Err: &err,

		Cwd:                "/",
		FS:                 options.FS,
		DefaultLibraryPath: bundled.LibPath(),
		TypingsLocation:    options.TypingsLocation,

		ParseCache: options.ParseCache,
		Client:     options.Client,
		Logger:     options.Logger,
		NpmInstall: npmInstall,
	})

	go func() {
		defer func() {
			outputWriter.Close()
		}()
		err := server.Run(context.TODO())
		if err != nil {
			t.Error("server error:", err)
		}
	}()

	s := &TestLspServer{
		Server:          server,
		in:              inputWriter,
		out:             outputReader,
		FS:              options.FS,
		UserPreferences: lsutil.NewDefaultUserPreferences(), // !!! parse default preferences for fourslash case?
		TypingsLocation: options.TypingsLocation,
	}

	// !!! temporary; remove when we have `handleDidChangeConfiguration`/implicit project config support
	// !!! replace with a proper request *after initialize*
	if options.OptionsForInferredProject != nil {
		s.Server.SetCompilerOptionsForInferredProjects(t.Context(), options.OptionsForInferredProject)
	}
	s.initialize(t, options.Capabilities)

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
	if capabilitiesWithDefaults.TextDocument.Diagnostic == nil {
		capabilitiesWithDefaults.TextDocument.Diagnostic = &lsproto.DiagnosticClientCapabilities{
			RelatedInformation: ptrTrue,
			TagSupport: &lsproto.ClientDiagnosticsTagOptions{
				ValueSet: []lsproto.DiagnosticTag{
					lsproto.DiagnosticTagUnnecessary,
					lsproto.DiagnosticTagDeprecated,
				},
			},
		}
	}
	if capabilitiesWithDefaults.Workspace == nil {
		capabilitiesWithDefaults.Workspace = &lsproto.WorkspaceClientCapabilities{}
	}
	if capabilitiesWithDefaults.Workspace.Configuration == nil {
		capabilitiesWithDefaults.Workspace.Configuration = ptrTrue
	}
	if capabilitiesWithDefaults.TextDocument.Definition == nil {
		capabilitiesWithDefaults.TextDocument.Definition = defaultDefinitionCapabilities
	}
	if capabilitiesWithDefaults.TextDocument.TypeDefinition == nil {
		capabilitiesWithDefaults.TextDocument.TypeDefinition = defaultTypeDefinitionCapabilities
	}
	if capabilitiesWithDefaults.TextDocument.Hover == nil {
		capabilitiesWithDefaults.TextDocument.Hover = defaultHoverCapabilities
	}
	if capabilitiesWithDefaults.TextDocument.SignatureHelp == nil {
		capabilitiesWithDefaults.TextDocument.SignatureHelp = &lsproto.SignatureHelpClientCapabilities{
			SignatureInformation: &lsproto.ClientSignatureInformationOptions{
				DocumentationFormat: &[]lsproto.MarkupKind{lsproto.MarkupKindMarkdown, lsproto.MarkupKindPlainText},
				ParameterInformation: &lsproto.ClientSignatureParameterInformationOptions{
					LabelOffsetSupport: ptrTrue,
				},
				ActiveParameterSupport: ptrTrue,
			},
			ContextSupport: ptrTrue,
		}
	}
	return &capabilitiesWithDefaults
}

func SendRequest[Params, Resp any](t *testing.T, server *TestLspServer, info lsproto.RequestInfo[Params, Resp], params Params) (*lsproto.Message, Resp, bool) {
	id := server.nextID()
	req := lsproto.NewRequestMessage(
		info.Method,
		lsproto.NewID(lsproto.IntegerOrString{Integer: &id}),
		params,
	)
	server.writeMsg(t, req.Message())
	resp := server.readMsg(t)
	if resp == nil {
		return nil, *new(Resp), false
	}

	// currently, the only request that may be sent by the server during a client request is one `config` request
	// !!! remove if `config` is handled in initialization and there are no other server-initiated requests
	if resp.Kind == lsproto.MessageKindRequest {
		req := resp.AsRequest()
		switch req.Method {
		case lsproto.MethodWorkspaceConfiguration:
			req := lsproto.ResponseMessage{
				ID:      req.ID,
				JSONRPC: req.JSONRPC,
				Result:  []any{server.UserPreferences},
			}
			server.writeMsg(t, req.Message())
			resp = server.readMsg(t)
		default:
			// other types of requests not yet used in fourslash; implement them if needed
			t.Fatalf("Unexpected request received: %s", req.Method)
		}
	}

	if resp == nil {
		return nil, *new(Resp), false
	}
	result, ok := resp.AsResponse().Result.(Resp)
	return resp, result, ok
}

func SendNotification[Params any](t *testing.T, server *TestLspServer, info lsproto.NotificationInfo[Params], params Params) {
	notification := lsproto.NewNotificationMessage(
		info.Method,
		params,
	)
	server.writeMsg(t, notification.Message())
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

func (s *TestLspServer) Session() *project.Session { return s.Server.Session() }

func ptrTo[T any](v T) *T {
	return &v
}

func Setup(t *testing.T, files map[string]any) (*TestLspServer, *projecttestutil.SessionUtils) {
	initOptions, sessionUtils := projecttestutil.GetSessionInitOptions(files, nil, &projecttestutil.TypingsInstallerOptions{})
	initOptions.Options.TypingsLocation = "" // Disable ata
	watchEnabledCapabilities := &lsproto.ClientCapabilities{
		Workspace: &lsproto.WorkspaceClientCapabilities{
			DidChangeWatchedFiles: &lsproto.DidChangeWatchedFilesClientCapabilities{
				DynamicRegistration: ptrTo(true),
			},
		},
	}
	server := NewTestLspServer(t, &TestLspServerOptions{
		FS:              initOptions.FS,
		Client:          initOptions.Client,
		Logger:          initOptions.Logger,
		NpmExecutor:     initOptions.NpmExecutor,
		ParseCache:      initOptions.ParseCache,
		TypingsLocation: initOptions.Options.TypingsLocation,
		Capabilities:    watchEnabledCapabilities,
	})
	return server, sessionUtils
}
