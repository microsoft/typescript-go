package lsp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime/debug"
	"slices"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/vfs"
	"golang.org/x/sync/errgroup"
)

type ServerOptions struct {
	In  Reader
	Out Writer
	Err io.Writer

	Cwd                string
	FS                 vfs.FS
	DefaultLibraryPath string
	TypingsLocation    string

	ParsedFileCache project.ParsedFileCache
}

func NewServer(opts *ServerOptions) *Server {
	if opts.Cwd == "" {
		panic("Cwd is required")
	}
	return &Server{
		r:                     opts.In,
		w:                     opts.Out,
		stderr:                opts.Err,
		requestQueue:          make(chan *lsproto.RequestMessage, 100),
		outgoingQueue:         make(chan *lsproto.Message, 100),
		pendingClientRequests: make(map[lsproto.ID]pendingClientRequest),
		pendingServerRequests: make(map[lsproto.ID]chan *lsproto.ResponseMessage),
		cwd:                   opts.Cwd,
		fs:                    opts.FS,
		defaultLibraryPath:    opts.DefaultLibraryPath,
		typingsLocation:       opts.TypingsLocation,
		parsedFileCache:       opts.ParsedFileCache,
	}
}

var (
	_ project.ServiceHost = (*Server)(nil)
	_ project.Client      = (*Server)(nil)
)

type pendingClientRequest struct {
	req    *lsproto.RequestMessage
	cancel context.CancelFunc
}

type Reader interface {
	Read() (*lsproto.Message, error)
}

type Writer interface {
	Write(msg *lsproto.Message) error
}

type lspReader struct {
	r *lsproto.BaseReader
}

type lspWriter struct {
	w *lsproto.BaseWriter
}

func (r *lspReader) Read() (*lsproto.Message, error) {
	data, err := r.r.Read()
	if err != nil {
		return nil, err
	}

	req := &lsproto.Message{}
	if err := json.Unmarshal(data, req); err != nil {
		return nil, fmt.Errorf("%w: %w", lsproto.ErrInvalidRequest, err)
	}

	return req, nil
}

func ToReader(r io.Reader) Reader {
	return &lspReader{r: lsproto.NewBaseReader(r)}
}

func (w *lspWriter) Write(msg *lsproto.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	return w.w.Write(data)
}

func ToWriter(w io.Writer) Writer {
	return &lspWriter{w: lsproto.NewBaseWriter(w)}
}

var (
	_ Reader = (*lspReader)(nil)
	_ Writer = (*lspWriter)(nil)
)

type Server struct {
	r Reader
	w Writer

	stderr io.Writer

	clientSeq               atomic.Int32
	requestQueue            chan *lsproto.RequestMessage
	outgoingQueue           chan *lsproto.Message
	pendingClientRequests   map[lsproto.ID]pendingClientRequest
	pendingClientRequestsMu sync.Mutex
	pendingServerRequests   map[lsproto.ID]chan *lsproto.ResponseMessage
	pendingServerRequestsMu sync.Mutex

	cwd                string
	fs                 vfs.FS
	defaultLibraryPath string
	typingsLocation    string

	initializeParams *lsproto.InitializeParams
	positionEncoding lsproto.PositionEncodingKind

	watchEnabled bool
	watcherID    atomic.Uint32
	watchers     collections.SyncSet[project.WatcherHandle]

	logger         *project.Logger
	projectService *project.Service

	// enables tests to share a cache of parsed source files
	parsedFileCache project.ParsedFileCache

	// !!! temporary; remove when we have `handleDidChangeConfiguration`/implicit project config support
	compilerOptionsForInferredProjects *core.CompilerOptions
}

// FS implements project.ServiceHost.
func (s *Server) FS() vfs.FS {
	return s.fs
}

// DefaultLibraryPath implements project.ServiceHost.
func (s *Server) DefaultLibraryPath() string {
	return s.defaultLibraryPath
}

// TypingsLocation implements project.ServiceHost.
func (s *Server) TypingsLocation() string {
	return s.typingsLocation
}

// GetCurrentDirectory implements project.ServiceHost.
func (s *Server) GetCurrentDirectory() string {
	return s.cwd
}

// Trace implements project.ServiceHost.
func (s *Server) Trace(msg string) {
	s.Log(msg)
}

// Client implements project.ServiceHost.
func (s *Server) Client() project.Client {
	if !s.watchEnabled {
		return nil
	}
	return s
}

// WatchFiles implements project.Client.
func (s *Server) WatchFiles(ctx context.Context, watchers []*lsproto.FileSystemWatcher) (project.WatcherHandle, error) {
	watcherId := fmt.Sprintf("watcher-%d", s.watcherID.Add(1))
	_, err := s.sendRequest(ctx, lsproto.MethodClientRegisterCapability, &lsproto.RegistrationParams{
		Registrations: []*lsproto.Registration{
			{
				Id:     watcherId,
				Method: string(lsproto.MethodWorkspaceDidChangeWatchedFiles),
				RegisterOptions: ptrTo(any(lsproto.DidChangeWatchedFilesRegistrationOptions{
					Watchers: watchers,
				})),
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to register file watcher: %w", err)
	}

	handle := project.WatcherHandle(watcherId)
	s.watchers.Add(handle)
	return handle, nil
}

// UnwatchFiles implements project.Client.
func (s *Server) UnwatchFiles(ctx context.Context, handle project.WatcherHandle) error {
	if s.watchers.Has(handle) {
		_, err := s.sendRequest(ctx, lsproto.MethodClientUnregisterCapability, &lsproto.UnregistrationParams{
			Unregisterations: []*lsproto.Unregistration{
				{
					Id:     string(handle),
					Method: string(lsproto.MethodWorkspaceDidChangeWatchedFiles),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to unregister file watcher: %w", err)
		}

		s.watchers.Delete(handle)
		return nil
	}

	return fmt.Errorf("no file watcher exists with ID %s", handle)
}

// RefreshDiagnostics implements project.Client.
func (s *Server) RefreshDiagnostics(ctx context.Context) error {
	if s.initializeParams.Capabilities == nil ||
		s.initializeParams.Capabilities.Workspace == nil ||
		s.initializeParams.Capabilities.Workspace.Diagnostics == nil ||
		!ptrIsTrue(s.initializeParams.Capabilities.Workspace.Diagnostics.RefreshSupport) {
		return nil
	}

	if _, err := s.sendRequest(ctx, lsproto.MethodWorkspaceDiagnosticRefresh, nil); err != nil {
		return fmt.Errorf("failed to refresh diagnostics: %w", err)
	}

	return nil
}

func (s *Server) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return s.dispatchLoop(ctx) })
	g.Go(func() error { return s.writeLoop(ctx) })

	// Don't run readLoop in the group, as it blocks on stdin read and cannot be cancelled.
	readLoopErr := make(chan error, 1)
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-readLoopErr:
			return err
		}
	})
	go func() { readLoopErr <- s.readLoop(ctx) }()

	if err := g.Wait(); err != nil && !errors.Is(err, io.EOF) && ctx.Err() != nil {
		return err
	}
	return nil
}

func (s *Server) readLoop(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		msg, err := s.read()
		if err != nil {
			if errors.Is(err, lsproto.ErrInvalidRequest) {
				s.sendError(nil, err)
				continue
			}
			return err
		}

		if s.initializeParams == nil && msg.Kind == lsproto.MessageKindRequest {
			req := msg.AsRequest()
			if req.Method == lsproto.MethodInitialize {
				s.handleInitialize(req)
			} else {
				s.sendError(req.ID, lsproto.ErrServerNotInitialized)
			}
			continue
		}

		if msg.Kind == lsproto.MessageKindResponse {
			resp := msg.AsResponse()
			s.pendingServerRequestsMu.Lock()
			if respChan, ok := s.pendingServerRequests[*resp.ID]; ok {
				respChan <- resp
				close(respChan)
				delete(s.pendingServerRequests, *resp.ID)
			}
			s.pendingServerRequestsMu.Unlock()
		} else {
			req := msg.AsRequest()
			if req.Method == lsproto.MethodCancelRequest {
				s.cancelRequest(req.Params.(*lsproto.CancelParams).Id)
			} else {
				s.requestQueue <- req
			}
		}
	}
}

func (s *Server) cancelRequest(rawID lsproto.IntegerOrString) {
	id := lsproto.NewID(rawID)
	s.pendingClientRequestsMu.Lock()
	defer s.pendingClientRequestsMu.Unlock()
	if pendingReq, ok := s.pendingClientRequests[*id]; ok {
		pendingReq.cancel()
		delete(s.pendingClientRequests, *id)
	}
}

func (s *Server) read() (*lsproto.Message, error) {
	return s.r.Read()
}

func (s *Server) dispatchLoop(ctx context.Context) error {
	ctx, lspExit := context.WithCancel(ctx)
	defer lspExit()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case req := <-s.requestQueue:
			requestCtx := ctx
			if req.ID != nil {
				var cancel context.CancelFunc
				requestCtx, cancel = context.WithCancel(core.WithRequestID(requestCtx, req.ID.String()))
				s.pendingClientRequestsMu.Lock()
				s.pendingClientRequests[*req.ID] = pendingClientRequest{
					req:    req,
					cancel: cancel,
				}
				s.pendingClientRequestsMu.Unlock()
			}

			handle := func() {
				defer func() {
					if r := recover(); r != nil {
						stack := debug.Stack()
						s.Log("panic handling request", req.Method, r, string(stack))
						if isBlockingMethod(req.Method) {
							lspExit()
						} else {
							if req.ID != nil {
								s.sendError(req.ID, fmt.Errorf("%w: panic handling request %s: %v", lsproto.ErrInternalError, req.Method, r))
							} else {
								s.Log("unhandled panic in notification", req.Method, r)
							}
						}
					}
				}()
				if err := s.handleRequestOrNotification(requestCtx, req); err != nil {
					if errors.Is(err, io.EOF) {
						lspExit()
					} else {
						s.sendError(req.ID, err)
					}
				}

				if req.ID != nil {
					s.pendingClientRequestsMu.Lock()
					delete(s.pendingClientRequests, *req.ID)
					s.pendingClientRequestsMu.Unlock()
				}
			}

			if isBlockingMethod(req.Method) {
				handle()
			} else {
				go handle()
			}
		}
	}
}

func (s *Server) writeLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-s.outgoingQueue:
			if err := s.w.Write(msg); err != nil {
				return fmt.Errorf("failed to write message: %w", err)
			}
		}
	}
}

func (s *Server) sendRequest(ctx context.Context, method lsproto.Method, params any) (any, error) {
	id := lsproto.NewIDString(fmt.Sprintf("ts%d", s.clientSeq.Add(1)))
	req := lsproto.NewRequestMessage(method, id, params)

	responseChan := make(chan *lsproto.ResponseMessage, 1)
	s.pendingServerRequestsMu.Lock()
	s.pendingServerRequests[*id] = responseChan
	s.pendingServerRequestsMu.Unlock()

	s.outgoingQueue <- req.Message()

	select {
	case <-ctx.Done():
		s.pendingServerRequestsMu.Lock()
		defer s.pendingServerRequestsMu.Unlock()
		if respChan, ok := s.pendingServerRequests[*id]; ok {
			close(respChan)
			delete(s.pendingServerRequests, *id)
		}
		return nil, ctx.Err()
	case resp := <-responseChan:
		if resp.Error != nil {
			return nil, fmt.Errorf("request failed: %s", resp.Error.String())
		}
		return resp.Result, nil
	}
}

func (s *Server) sendResult(id *lsproto.ID, result any) {
	s.sendResponse(&lsproto.ResponseMessage{
		ID:     id,
		Result: result,
	})
}

func (s *Server) sendError(id *lsproto.ID, err error) {
	code := lsproto.ErrInternalError.Code
	if errCode := (*lsproto.ErrorCode)(nil); errors.As(err, &errCode) {
		code = errCode.Code
	}
	// TODO(jakebailey): error data
	s.sendResponse(&lsproto.ResponseMessage{
		ID: id,
		Error: &lsproto.ResponseError{
			Code:    code,
			Message: err.Error(),
		},
	})
}

func (s *Server) sendResponse(resp *lsproto.ResponseMessage) {
	s.outgoingQueue <- resp.Message()
}

func (s *Server) handleRequestOrNotification(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params
	switch params.(type) {
	case *lsproto.InitializeParams:
		s.sendError(req.ID, lsproto.ErrInvalidRequest)
		return nil
	case *lsproto.InitializedParams:
		return handleNotification(s, ctx, lsproto.InitializedHandler, req, (*Server).handleInitialized)
	case *lsproto.DidOpenTextDocumentParams:
		return handleNotification(s, ctx, lsproto.TextDocumentDidOpenHandler, req, (*Server).handleDidOpen)
	case *lsproto.DidChangeTextDocumentParams:
		return handleNotification(s, ctx, lsproto.TextDocumentDidChangeHandler, req, (*Server).handleDidChange)
	case *lsproto.DidSaveTextDocumentParams:
		return handleNotification(s, ctx, lsproto.TextDocumentDidSaveHandler, req, (*Server).handleDidSave)
	case *lsproto.DidCloseTextDocumentParams:
		return handleNotification(s, ctx, lsproto.TextDocumentDidCloseHandler, req, (*Server).handleDidClose)
	case *lsproto.DidChangeWatchedFilesParams:
		return handleNotification(s, ctx, lsproto.WorkspaceDidChangeWatchedFilesHandler, req, (*Server).handleDidChangeWatchedFiles)
	case *lsproto.DocumentDiagnosticParams:
		return handleWithSingleResponse(s, ctx, lsproto.TextDocumentDiagnosticHandler, req, (*Server).handleDocumentDiagnostic)
	case *lsproto.HoverParams:
		return handleWithSingleResponse(s, ctx, lsproto.TextDocumentHoverHandler, req, (*Server).handleHover)
	case *lsproto.DefinitionParams:
		return handleWithSingleResponse(s, ctx, lsproto.TextDocumentDefinitionHandler, req, (*Server).handleDefinition)
	case *lsproto.TypeDefinitionParams:
		return handleWithSingleResponse(s, ctx, lsproto.TextDocumentTypeDefinitionHandler, req, (*Server).handleTypeDefinition)
	case *lsproto.CompletionParams:
		return handleWithSingleResponse(s, ctx, lsproto.TextDocumentCompletionHandler, req, (*Server).handleCompletion)
	case *lsproto.ReferenceParams:
		return handleWithSingleResponse(s, ctx, lsproto.TextDocumentReferencesHandler, req, (*Server).handleReferences)
	case *lsproto.ImplementationParams:
		return handleWithSingleResponse(s, ctx, lsproto.TextDocumentImplementationHandler, req, (*Server).handleImplementations)
	case *lsproto.SignatureHelpParams:
		return handleWithSingleResponse(s, ctx, lsproto.TextDocumentSignatureHelpHandler, req, (*Server).handleSignatureHelp)
	case *lsproto.DocumentFormattingParams:
		return handleWithSingleResponse(s, ctx, lsproto.TextDocumentFormattingHandler, req, (*Server).handleDocumentFormat)
	case *lsproto.DocumentRangeFormattingParams:
		return handleWithSingleResponse(s, ctx, lsproto.TextDocumentRangeFormattingHandler, req, (*Server).handleDocumentRangeFormat)
	case *lsproto.DocumentOnTypeFormattingParams:
		return handleWithSingleResponse(s, ctx, lsproto.TextDocumentOnTypeFormattingHandler, req, (*Server).handleDocumentOnTypeFormat)
	case *lsproto.WorkspaceSymbolParams:
		return handleWithSingleResponse(s, ctx, lsproto.WorkspaceSymbolHandler, req, (*Server).handleWorkspaceSymbol)
	case *lsproto.DocumentSymbolParams:
		return handleWithSingleResponse(s, ctx, lsproto.TextDocumentDocumentSymbolHandler, req, (*Server).handleDocumentSymbol)
	case *lsproto.CompletionItem:
		return handleWithSingleResponse(s, ctx, lsproto.CompletionItemResolveHandler, req, (*Server).handleCompletionItemResolve)
	default:
		switch req.Method {
		case lsproto.MethodShutdown:
			s.projectService.Close()
			s.sendResult(req.ID, nil)
			return nil
		case lsproto.MethodExit:
			return io.EOF
		default:
			s.Log("unknown method", req.Method)
			if req.ID != nil {
				s.sendError(req.ID, lsproto.ErrInvalidRequest)
			}
			return nil
		}
	}
}

func handleNotification[Req any](
	s *Server,
	ctx context.Context,
	info lsproto.NotificationMapping[Req],
	req *lsproto.RequestMessage,
	fn func(*Server, context.Context, Req) error,
) error {
	if req.Method != info.Method {
		panic(fmt.Sprintf("expected method %s, got %s", info.Method, req.Method))
	}
	params := req.Params.(Req)
	if err := fn(s, ctx, params); err != nil {
		return err
	}
	return ctx.Err()
}

func handleWithSingleResponse[Req, Resp any](
	s *Server,
	ctx context.Context,
	info lsproto.RequestToResponseMapping[Req, Resp],
	req *lsproto.RequestMessage,
	fn func(*Server, context.Context, Req) (Resp, error),
) error {
	if req.Method != info.Method {
		panic(fmt.Sprintf("expected method %s, got %s", info.Method, req.Method))
	}
	params := req.Params.(Req)
	resp, err := fn(s, ctx, params)
	if err != nil {
		return err
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	s.sendResult(req.ID, resp)
	return nil
}

func (s *Server) handleInitialize(req *lsproto.RequestMessage) {
	s.initializeParams = req.Params.(*lsproto.InitializeParams)

	s.positionEncoding = lsproto.PositionEncodingKindUTF16
	if genCapabilities := s.initializeParams.Capabilities.General; genCapabilities != nil && genCapabilities.PositionEncodings != nil {
		if slices.Contains(*genCapabilities.PositionEncodings, lsproto.PositionEncodingKindUTF8) {
			s.positionEncoding = lsproto.PositionEncodingKindUTF8
		}
	}

	s.sendResult(req.ID, &lsproto.InitializeResult{
		ServerInfo: &lsproto.ServerInfo{
			Name:    "typescript-go",
			Version: ptrTo(core.Version()),
		},
		Capabilities: &lsproto.ServerCapabilities{
			PositionEncoding: ptrTo(s.positionEncoding),
			TextDocumentSync: &lsproto.TextDocumentSyncOptionsOrTextDocumentSyncKind{
				TextDocumentSyncOptions: &lsproto.TextDocumentSyncOptions{
					OpenClose: ptrTo(true),
					Change:    ptrTo(lsproto.TextDocumentSyncKindIncremental),
					Save: &lsproto.BooleanOrSaveOptions{
						SaveOptions: &lsproto.SaveOptions{
							IncludeText: ptrTo(true),
						},
					},
				},
			},
			HoverProvider: &lsproto.BooleanOrHoverOptions{
				Boolean: ptrTo(true),
			},
			DefinitionProvider: &lsproto.BooleanOrDefinitionOptions{
				Boolean: ptrTo(true),
			},
			TypeDefinitionProvider: &lsproto.BooleanOrTypeDefinitionOptionsOrTypeDefinitionRegistrationOptions{
				Boolean: ptrTo(true),
			},
			ReferencesProvider: &lsproto.BooleanOrReferenceOptions{
				Boolean: ptrTo(true),
			},
			ImplementationProvider: &lsproto.BooleanOrImplementationOptionsOrImplementationRegistrationOptions{
				Boolean: ptrTo(true),
			},
			DiagnosticProvider: &lsproto.DiagnosticOptionsOrDiagnosticRegistrationOptions{
				DiagnosticOptions: &lsproto.DiagnosticOptions{
					InterFileDependencies: true,
				},
			},
			CompletionProvider: &lsproto.CompletionOptions{
				TriggerCharacters: &ls.TriggerCharacters,
				ResolveProvider:   ptrTo(true),
				// !!! other options
			},
			SignatureHelpProvider: &lsproto.SignatureHelpOptions{
				TriggerCharacters: &[]string{"(", ","},
			},
			DocumentFormattingProvider: &lsproto.BooleanOrDocumentFormattingOptions{
				Boolean: ptrTo(true),
			},
			DocumentRangeFormattingProvider: &lsproto.BooleanOrDocumentRangeFormattingOptions{
				Boolean: ptrTo(true),
			},
			DocumentOnTypeFormattingProvider: &lsproto.DocumentOnTypeFormattingOptions{
				FirstTriggerCharacter: "{",
				MoreTriggerCharacter:  &[]string{"}", ";", "\n"},
			},
			WorkspaceSymbolProvider: &lsproto.BooleanOrWorkspaceSymbolOptions{
				Boolean: ptrTo(true),
			},
			DocumentSymbolProvider: &lsproto.BooleanOrDocumentSymbolOptions{
				Boolean: ptrTo(true),
			},
		},
	})
}

func (s *Server) handleInitialized(ctx context.Context, params *lsproto.InitializedParams) error {
	if shouldEnableWatch(s.initializeParams) {
		s.watchEnabled = true
	}

	s.logger = project.NewLogger([]io.Writer{s.stderr}, "" /*file*/, project.LogLevelVerbose)
	s.projectService = project.NewService(s, project.ServiceOptions{
		Logger:           s.logger,
		WatchEnabled:     s.watchEnabled,
		PositionEncoding: s.positionEncoding,
		TypingsInstallerOptions: project.TypingsInstallerOptions{
			ThrottleLimit: 5,
			NpmInstall:    project.NpmInstall,
		},
		ParsedFileCache: s.parsedFileCache,
	})
	// !!! temporary; remove when we have `handleDidChangeConfiguration`/implicit project config support
	if s.compilerOptionsForInferredProjects != nil {
		s.projectService.SetCompilerOptionsForInferredProjects(s.compilerOptionsForInferredProjects)
	}

	return nil
}

func (s *Server) handleDidOpen(ctx context.Context, params *lsproto.DidOpenTextDocumentParams) error {
	s.projectService.OpenFile(ls.DocumentURIToFileName(params.TextDocument.Uri), params.TextDocument.Text, ls.LanguageKindToScriptKind(params.TextDocument.LanguageId), "")
	return nil
}

func (s *Server) handleDidChange(ctx context.Context, params *lsproto.DidChangeTextDocumentParams) error {
	return s.projectService.ChangeFile(params.TextDocument, params.ContentChanges)
}

func (s *Server) handleDidSave(ctx context.Context, params *lsproto.DidSaveTextDocumentParams) error {
	s.projectService.MarkFileSaved(ls.DocumentURIToFileName(params.TextDocument.Uri), *params.Text)
	return nil
}

func (s *Server) handleDidClose(ctx context.Context, params *lsproto.DidCloseTextDocumentParams) error {
	s.projectService.CloseFile(ls.DocumentURIToFileName(params.TextDocument.Uri))
	return nil
}

func (s *Server) handleDidChangeWatchedFiles(ctx context.Context, params *lsproto.DidChangeWatchedFilesParams) error {
	return s.projectService.OnWatchedFilesChanged(ctx, params.Changes)
}

func (s *Server) handleDocumentDiagnostic(ctx context.Context, params *lsproto.DocumentDiagnosticParams) (lsproto.DocumentDiagnosticResponse, error) {
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	return languageService.ProvideDiagnostics(ctx, params.TextDocument.Uri)
}

func (s *Server) handleHover(ctx context.Context, params *lsproto.HoverParams) (lsproto.HoverResponse, error) {
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	return languageService.ProvideHover(ctx, params.TextDocument.Uri, params.Position)
}

func (s *Server) handleSignatureHelp(ctx context.Context, params *lsproto.SignatureHelpParams) (lsproto.SignatureHelpResponse, error) {
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	return languageService.ProvideSignatureHelp(
		ctx,
		params.TextDocument.Uri,
		params.Position,
		params.Context,
		s.initializeParams.Capabilities.TextDocument.SignatureHelp,
		&ls.UserPreferences{},
	), nil
}

func (s *Server) handleDefinition(ctx context.Context, params *lsproto.DefinitionParams) (lsproto.DefinitionResponse, error) {
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	return languageService.ProvideDefinition(ctx, params.TextDocument.Uri, params.Position)
}

func (s *Server) handleTypeDefinition(ctx context.Context, params *lsproto.TypeDefinitionParams) (lsproto.TypeDefinitionResponse, error) {
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	return languageService.ProvideTypeDefinition(ctx, params.TextDocument.Uri, params.Position)
}

func (s *Server) handleReferences(ctx context.Context, params *lsproto.ReferenceParams) (lsproto.ReferencesResponse, error) {
	// findAllReferences
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	locations := languageService.ProvideReferences(params)
	return &locations, nil
}

func (s *Server) handleImplementations(ctx context.Context, params *lsproto.ImplementationParams) (lsproto.ImplementationResponse, error) {
	// goToImplementation
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	locations := languageService.ProvideImplementations(params)
	return &lsproto.LocationOrLocationsOrDefinitionLinks{
		Locations: &locations,
	}, nil
}

func (s *Server) handleCompletion(ctx context.Context, params *lsproto.CompletionParams) (lsproto.CompletionResponse, error) {
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	// !!! get user preferences
	list, err := languageService.ProvideCompletion(
		ctx,
		params.TextDocument.Uri,
		params.Position,
		params.Context,
		getCompletionClientCapabilities(s.initializeParams),
		&ls.UserPreferences{})
	if err != nil {
		return nil, err
	}
	return &lsproto.CompletionItemsOrCompletionList{
		CompletionList: list,
	}, nil
}

func (s *Server) handleCompletionItemResolve(ctx context.Context, params *lsproto.CompletionItem) (lsproto.CompletionResolveResponse, error) {
	data, err := ls.GetCompletionItemData(params)
	if err != nil {
		return nil, err
	}
	_, project := s.projectService.EnsureDefaultProjectForFile(data.FileName)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	return languageService.ResolveCompletionItem(
		ctx,
		params,
		data,
		getCompletionClientCapabilities(s.initializeParams),
		&ls.UserPreferences{},
	)
}

func (s *Server) handleDocumentFormat(ctx context.Context, params *lsproto.DocumentFormattingParams) (lsproto.DocumentFormattingResponse, error) {
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	res, err := languageService.ProvideFormatDocument(
		ctx,
		params.TextDocument.Uri,
		params.Options,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *Server) handleDocumentRangeFormat(ctx context.Context, params *lsproto.DocumentRangeFormattingParams) (lsproto.DocumentRangeFormattingResponse, error) {
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	res, err := languageService.ProvideFormatDocumentRange(
		ctx,
		params.TextDocument.Uri,
		params.Options,
		params.Range,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *Server) handleDocumentOnTypeFormat(ctx context.Context, params *lsproto.DocumentOnTypeFormattingParams) (lsproto.DocumentOnTypeFormattingResponse, error) {
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	res, err := languageService.ProvideFormatDocumentOnType(
		ctx,
		params.TextDocument.Uri,
		params.Options,
		params.Position,
		params.Ch,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *Server) handleWorkspaceSymbol(ctx context.Context, params *lsproto.WorkspaceSymbolParams) (lsproto.WorkspaceSymbolResponse, error) {
	programs := core.Map(s.projectService.Projects(), (*project.Project).GetProgram)
	symbols, err := ls.ProvideWorkspaceSymbols(ctx, programs, s.projectService.Converters(), params.Query)
	if err != nil {
		return nil, err
	}
	return &lsproto.SymbolInformationsOrWorkspaceSymbols{
		SymbolInformations: &symbols,
	}, nil
}

func (s *Server) handleDocumentSymbol(ctx context.Context, params *lsproto.DocumentSymbolParams) (lsproto.DocumentSymbolResponse, error) {
	project := s.projectService.EnsureDefaultProjectForURI(params.TextDocument.Uri)
	languageService, done := project.GetLanguageServiceForRequest(ctx)
	defer done()
	symbols, err := languageService.ProvideDocumentSymbols(ctx, params.TextDocument.Uri)
	if err != nil {
		return nil, err
	}
	return &lsproto.SymbolInformationsOrDocumentSymbols{
		DocumentSymbols: &symbols,
	}, nil
}

func (s *Server) Log(msg ...any) {
	fmt.Fprintln(s.stderr, msg...)
}

// !!! temporary; remove when we have `handleDidChangeConfiguration`/implicit project config support
func (s *Server) SetCompilerOptionsForInferredProjects(options *core.CompilerOptions) {
	s.compilerOptionsForInferredProjects = options
	if s.projectService != nil {
		s.projectService.SetCompilerOptionsForInferredProjects(options)
	}
}

func isBlockingMethod(method lsproto.Method) bool {
	switch method {
	case lsproto.MethodInitialize,
		lsproto.MethodInitialized,
		lsproto.MethodTextDocumentDidOpen,
		lsproto.MethodTextDocumentDidChange,
		lsproto.MethodTextDocumentDidSave,
		lsproto.MethodTextDocumentDidClose,
		lsproto.MethodWorkspaceDidChangeWatchedFiles:
		return true
	}
	return false
}

func ptrTo[T any](v T) *T {
	return &v
}

func ptrIsTrue(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

func shouldEnableWatch(params *lsproto.InitializeParams) bool {
	if params == nil || params.Capabilities == nil || params.Capabilities.Workspace == nil {
		return false
	}
	return params.Capabilities.Workspace.DidChangeWatchedFiles != nil &&
		ptrIsTrue(params.Capabilities.Workspace.DidChangeWatchedFiles.DynamicRegistration)
}

func getCompletionClientCapabilities(params *lsproto.InitializeParams) *lsproto.CompletionClientCapabilities {
	if params == nil || params.Capabilities == nil || params.Capabilities.TextDocument == nil {
		return nil
	}
	return params.Capabilities.TextDocument.Completion
}
