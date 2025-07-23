package lsp

import (
	"context"
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
	"github.com/microsoft/typescript-go/internal/projectv2"
	"github.com/microsoft/typescript-go/internal/vfs"
	"golang.org/x/sync/errgroup"
)

func NewProjectV2Server(opts ServerOptions) *ProjectV2Server {
	if opts.Cwd == "" {
		panic("Cwd is required")
	}
	return &ProjectV2Server{
		r:                     opts.In,
		w:                     opts.Out,
		stderr:                opts.Err,
		requestQueue:          make(chan *lsproto.RequestMessage, 100),
		outgoingQueue:         make(chan *lsproto.Message, 100),
		logQueue:              make(chan string, 100),
		pendingClientRequests: make(map[lsproto.ID]pendingClientRequest),
		pendingServerRequests: make(map[lsproto.ID]chan *lsproto.ResponseMessage),
		cwd:                   opts.Cwd,
		fs:                    opts.FS,
		defaultLibraryPath:    opts.DefaultLibraryPath,
		typingsLocation:       opts.TypingsLocation,
		parsedFileCache:       opts.ParsedFileCache,
	}
}

type ProjectV2Server struct {
	r Reader
	w Writer

	stderr io.Writer

	clientSeq               atomic.Int32
	requestQueue            chan *lsproto.RequestMessage
	outgoingQueue           chan *lsproto.Message
	logQueue                chan string
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
	watchers     collections.SyncSet[projectv2.WatcherID]

	session *projectv2.Session

	// enables tests to share a cache of parsed source files
	parsedFileCache project.ParsedFileCache

	// !!! temporary; remove when we have `handleDidChangeConfiguration`/implicit project config support
	compilerOptionsForInferredProjects *core.CompilerOptions
}

// FS implements project.ServiceHost.
func (s *ProjectV2Server) FS() vfs.FS {
	return s.fs
}

// DefaultLibraryPath implements project.ServiceHost.
func (s *ProjectV2Server) DefaultLibraryPath() string {
	return s.defaultLibraryPath
}

// TypingsLocation implements project.ServiceHost.
func (s *ProjectV2Server) TypingsLocation() string {
	return s.typingsLocation
}

// GetCurrentDirectory implements project.ServiceHost.
func (s *ProjectV2Server) GetCurrentDirectory() string {
	return s.cwd
}

// Trace implements project.ServiceHost.
func (s *ProjectV2Server) Trace(msg string) {
	s.Log(msg)
}

// Client implements project.ServiceHost.
func (s *ProjectV2Server) Client() projectv2.Client {
	if !s.watchEnabled {
		return nil
	}
	return s
}

// WatchFiles implements project.Client.
func (s *ProjectV2Server) WatchFiles(ctx context.Context, id projectv2.WatcherID, watchers []*lsproto.FileSystemWatcher) error {
	_, err := s.sendRequest(ctx, lsproto.MethodClientRegisterCapability, &lsproto.RegistrationParams{
		Registrations: []*lsproto.Registration{
			{
				Id:     string(id),
				Method: string(lsproto.MethodWorkspaceDidChangeWatchedFiles),
				RegisterOptions: ptrTo(any(lsproto.DidChangeWatchedFilesRegistrationOptions{
					Watchers: watchers,
				})),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to register file watcher: %w", err)
	}

	s.watchers.Add(id)
	return nil
}

// UnwatchFiles implements project.Client.
func (s *ProjectV2Server) UnwatchFiles(ctx context.Context, id projectv2.WatcherID) error {
	if s.watchers.Has(id) {
		_, err := s.sendRequest(ctx, lsproto.MethodClientUnregisterCapability, &lsproto.UnregistrationParams{
			Unregisterations: []*lsproto.Unregistration{
				{
					Id:     string(id),
					Method: string(lsproto.MethodWorkspaceDidChangeWatchedFiles),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to unregister file watcher: %w", err)
		}

		s.watchers.Delete(id)
		return nil
	}

	return fmt.Errorf("no file watcher exists with ID %s", id)
}

// RefreshDiagnostics implements project.Client.
func (s *ProjectV2Server) RefreshDiagnostics(ctx context.Context) error {
	if ptrIsTrue(s.initializeParams.Capabilities.Workspace.Diagnostics.RefreshSupport) {
		if _, err := s.sendRequest(ctx, lsproto.MethodWorkspaceDiagnosticRefresh, nil); err != nil {
			return fmt.Errorf("failed to refresh diagnostics: %w", err)
		}
	}
	return nil
}

func (s *ProjectV2Server) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return s.dispatchLoop(ctx) })
	g.Go(func() error { return s.writeLoop(ctx) })
	g.Go(func() error { return s.logLoop(ctx) })

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

func (s *ProjectV2Server) readLoop(ctx context.Context) error {
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

func (s *ProjectV2Server) cancelRequest(rawID lsproto.IntegerOrString) {
	id := lsproto.NewID(rawID)
	s.pendingClientRequestsMu.Lock()
	defer s.pendingClientRequestsMu.Unlock()
	if pendingReq, ok := s.pendingClientRequests[*id]; ok {
		pendingReq.cancel()
		delete(s.pendingClientRequests, *id)
	}
}

func (s *ProjectV2Server) read() (*lsproto.Message, error) {
	return s.r.Read()
}

func (s *ProjectV2Server) dispatchLoop(ctx context.Context) error {
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
						// !!! send something back to client
						lspExit()
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

func (s *ProjectV2Server) writeLoop(ctx context.Context) error {
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

func (s *ProjectV2Server) logLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case logMessage := <-s.logQueue:
			if _, err := fmt.Fprintln(s.stderr, logMessage); err != nil {
				return fmt.Errorf("failed to write log message: %w", err)
			}
		}
	}
}

func (s *ProjectV2Server) sendRequest(ctx context.Context, method lsproto.Method, params any) (any, error) {
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

func (s *ProjectV2Server) sendResult(id *lsproto.ID, result any) {
	s.sendResponse(&lsproto.ResponseMessage{
		ID:     id,
		Result: result,
	})
}

func (s *ProjectV2Server) sendError(id *lsproto.ID, err error) {
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

func (s *ProjectV2Server) sendResponse(resp *lsproto.ResponseMessage) {
	s.outgoingQueue <- resp.Message()
}

func (s *ProjectV2Server) handleRequestOrNotification(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params
	switch params.(type) {
	case *lsproto.InitializeParams:
		s.sendError(req.ID, lsproto.ErrInvalidRequest)
		return nil
	case *lsproto.InitializedParams:
		return s.handleInitialized(ctx, req)
	case *lsproto.DidOpenTextDocumentParams:
		return s.handleDidOpen(ctx, req)
	case *lsproto.DidChangeTextDocumentParams:
		return s.handleDidChange(ctx, req)
	case *lsproto.DidSaveTextDocumentParams:
		return s.handleDidSave(ctx, req)
	case *lsproto.DidCloseTextDocumentParams:
		return s.handleDidClose(ctx, req)
	case *lsproto.DidChangeWatchedFilesParams:
		return s.handleDidChangeWatchedFiles(ctx, req)
	case *lsproto.DocumentDiagnosticParams:
		return s.handleDocumentDiagnostic(ctx, req)
	case *lsproto.HoverParams:
		return s.handleHover(ctx, req)
	case *lsproto.DefinitionParams:
		return s.handleDefinition(ctx, req)
	case *lsproto.CompletionParams:
		return s.handleCompletion(ctx, req)
	case *lsproto.ReferenceParams:
		return s.handleReferences(ctx, req)
	case *lsproto.SignatureHelpParams:
		return s.handleSignatureHelp(ctx, req)
	case *lsproto.DocumentFormattingParams:
		return s.handleDocumentFormat(ctx, req)
	case *lsproto.DocumentRangeFormattingParams:
		return s.handleDocumentRangeFormat(ctx, req)
	case *lsproto.DocumentOnTypeFormattingParams:
		return s.handleDocumentOnTypeFormat(ctx, req)
	default:
		switch req.Method {
		case lsproto.MethodShutdown:
			s.session.Close()
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

func (s *ProjectV2Server) handleInitialize(req *lsproto.RequestMessage) {
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
			ReferencesProvider: &lsproto.BooleanOrReferenceOptions{
				Boolean: ptrTo(true),
			},
			DiagnosticProvider: &lsproto.DiagnosticOptionsOrDiagnosticRegistrationOptions{
				DiagnosticOptions: &lsproto.DiagnosticOptions{
					InterFileDependencies: true,
				},
			},
			CompletionProvider: &lsproto.CompletionOptions{
				TriggerCharacters: &ls.TriggerCharacters,
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
		},
	})
}

func (s *ProjectV2Server) handleInitialized(ctx context.Context, req *lsproto.RequestMessage) error {
	if shouldEnableWatch(s.initializeParams) {
		s.watchEnabled = true
	}

	s.session = projectv2.NewSession(projectv2.SessionOptions{
		CurrentDirectory:   s.cwd,
		DefaultLibraryPath: s.defaultLibraryPath,
		TypingsLocation:    s.typingsLocation,
		PositionEncoding:   s.positionEncoding,
		WatchEnabled:       s.watchEnabled,
		LoggingEnabled:     true,
	}, s.fs, s.Client(), s)

	return nil
}

func (s *ProjectV2Server) handleDidOpen(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DidOpenTextDocumentParams)
	s.session.DidOpenFile(ctx, params.TextDocument.Uri, params.TextDocument.Version, params.TextDocument.Text, params.TextDocument.LanguageId)
	return nil
}

func (s *ProjectV2Server) handleDidChange(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DidChangeTextDocumentParams)
	s.session.DidChangeFile(ctx, params.TextDocument.Uri, params.TextDocument.Version, params.ContentChanges)
	return nil
}

func (s *ProjectV2Server) handleDidSave(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DidSaveTextDocumentParams)
	s.session.DidSaveFile(ctx, params.TextDocument.Uri)
	return nil
}

func (s *ProjectV2Server) handleDidClose(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DidCloseTextDocumentParams)
	s.session.DidCloseFile(ctx, params.TextDocument.Uri)
	return nil
}

func (s *ProjectV2Server) handleDidChangeWatchedFiles(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DidChangeWatchedFilesParams)
	s.session.DidChangeWatchedFiles(ctx, params.Changes)
	return nil
}

func (s *ProjectV2Server) handleDocumentDiagnostic(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DocumentDiagnosticParams)
	languageService, err := s.session.GetLanguageService(ctx, params.TextDocument.Uri)
	if err != nil {
		return err
	}
	diagnostics, err := languageService.GetDocumentDiagnostics(ctx, params.TextDocument.Uri)
	if err != nil {
		return err
	}
	s.sendResult(req.ID, diagnostics)
	return nil
}

func (s *ProjectV2Server) handleHover(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.HoverParams)
	languageService, err := s.session.GetLanguageService(ctx, params.TextDocument.Uri)
	if err != nil {
		return err
	}
	hover, err := languageService.ProvideHover(ctx, params.TextDocument.Uri, params.Position)
	if err != nil {
		return err
	}
	s.sendResult(req.ID, hover)
	return nil
}

func (s *ProjectV2Server) handleSignatureHelp(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.SignatureHelpParams)
	languageService, err := s.session.GetLanguageService(ctx, params.TextDocument.Uri)
	if err != nil {
		return err
	}
	signatureHelp := languageService.ProvideSignatureHelp(
		ctx,
		params.TextDocument.Uri,
		params.Position,
		params.Context,
		s.initializeParams.Capabilities.TextDocument.SignatureHelp,
		&ls.UserPreferences{},
	)
	s.sendResult(req.ID, signatureHelp)
	return nil
}

func (s *ProjectV2Server) handleDefinition(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DefinitionParams)
	languageService, err := s.session.GetLanguageService(ctx, params.TextDocument.Uri)
	if err != nil {
		return err
	}
	definition, err := languageService.ProvideDefinition(ctx, params.TextDocument.Uri, params.Position)
	if err != nil {
		return err
	}
	s.sendResult(req.ID, definition)
	return nil
}

func (s *ProjectV2Server) handleReferences(ctx context.Context, req *lsproto.RequestMessage) error {
	// findAllReferences
	params := req.Params.(*lsproto.ReferenceParams)
	languageService, err := s.session.GetLanguageService(ctx, params.TextDocument.Uri)
	if err != nil {
		return err
	}
	// !!! remove this after find all references is fully ported/tested
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			s.Log("panic obtaining references:", r, string(stack))
			s.sendResult(req.ID, []*lsproto.Location{})
		}
	}()

	locations := languageService.ProvideReferences(params)
	s.sendResult(req.ID, locations)
	return nil
}

func (s *ProjectV2Server) handleCompletion(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.CompletionParams)
	languageService, err := s.session.GetLanguageService(ctx, params.TextDocument.Uri)
	if err != nil {
		return err
	}
	// !!! remove this after completions is fully ported/tested
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			s.Log("panic obtaining completions:", r, string(stack))
			s.sendResult(req.ID, &lsproto.CompletionList{})
		}
	}()
	// !!! get user preferences
	list, err := languageService.ProvideCompletion(
		ctx,
		params.TextDocument.Uri,
		params.Position,
		params.Context,
		getCompletionClientCapabilities(s.initializeParams),
		&ls.UserPreferences{})
	if err != nil {
		return err
	}
	s.sendResult(req.ID, list)
	return nil
}

func (s *ProjectV2Server) handleDocumentFormat(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DocumentFormattingParams)
	languageService, err := s.session.GetLanguageService(ctx, params.TextDocument.Uri)
	if err != nil {
		return err
	}
	// !!! remove this after formatting is fully ported/tested
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			s.Log("panic on document format:", r, string(stack))
			s.sendResult(req.ID, []*lsproto.TextEdit{})
		}
	}()

	res, err := languageService.ProvideFormatDocument(
		ctx,
		params.TextDocument.Uri,
		params.Options,
	)
	if err != nil {
		return err
	}
	s.sendResult(req.ID, res)
	return nil
}

func (s *ProjectV2Server) handleDocumentRangeFormat(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DocumentRangeFormattingParams)
	languageService, err := s.session.GetLanguageService(ctx, params.TextDocument.Uri)
	if err != nil {
		return err
	}
	// !!! remove this after formatting is fully ported/tested
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			s.Log("panic on document range format:", r, string(stack))
			s.sendResult(req.ID, []*lsproto.TextEdit{})
		}
	}()

	res, err := languageService.ProvideFormatDocumentRange(
		ctx,
		params.TextDocument.Uri,
		params.Options,
		params.Range,
	)
	if err != nil {
		return err
	}
	s.sendResult(req.ID, res)
	return nil
}

func (s *ProjectV2Server) handleDocumentOnTypeFormat(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DocumentOnTypeFormattingParams)
	languageService, err := s.session.GetLanguageService(ctx, params.TextDocument.Uri)
	if err != nil {
		return err
	}
	// !!! remove this after formatting is fully ported/tested
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			s.Log("panic on type format:", r, string(stack))
			s.sendResult(req.ID, []*lsproto.TextEdit{})
		}
	}()

	res, err := languageService.ProvideFormatDocumentOnType(
		ctx,
		params.TextDocument.Uri,
		params.Options,
		params.Position,
		params.Ch,
	)
	if err != nil {
		return err
	}
	s.sendResult(req.ID, res)
	return nil
}

// Log implements projectv2.Logger interface
func (s *ProjectV2Server) Log(msg ...any) {
	s.logQueue <- fmt.Sprint(msg...)
}
