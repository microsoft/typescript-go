package lsp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"runtime/debug"
	"slices"
	"strings"
	"sync"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type ServerOptions struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer

	Cwd                string
	NewLine            core.NewLineKind
	FS                 vfs.FS
	DefaultLibraryPath string
}

func NewServer(opts *ServerOptions) *Server {
	if opts.Cwd == "" {
		panic("Cwd is required")
	}
	return &Server{
		r:                     lsproto.NewBaseReader(opts.In),
		w:                     lsproto.NewBaseWriter(opts.Out),
		stderr:                opts.Err,
		fatalErrChan:          make(chan error, 1),
		requestQueue:          make(chan *lsproto.RequestMessage, 100),
		outgoingQueue:         make(chan *lsproto.Message, 100),
		pendingClientRequests: make(map[lsproto.ID]pendingClientRequest),
		pendingServerRequests: make(map[lsproto.ID]chan *lsproto.ResponseMessage),
		cwd:                   opts.Cwd,
		newLine:               opts.NewLine,
		fs:                    opts.FS,
		defaultLibraryPath:    opts.DefaultLibraryPath,
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

type Server struct {
	r *lsproto.BaseReader
	w *lsproto.BaseWriter

	stderr io.Writer

	clientSeq               int32
	fatalErrChan            chan error
	requestQueue            chan *lsproto.RequestMessage
	outgoingQueue           chan *lsproto.Message
	pendingClientRequests   map[lsproto.ID]pendingClientRequest
	pendingClientRequestsMu sync.Mutex
	pendingServerRequests   map[lsproto.ID]chan *lsproto.ResponseMessage
	pendingServerRequestsMu sync.Mutex

	cwd                string
	newLine            core.NewLineKind
	fs                 vfs.FS
	defaultLibraryPath string

	initializeParams *lsproto.InitializeParams
	positionEncoding lsproto.PositionEncodingKind

	watchEnabled   bool
	watcherID      int
	watchers       core.Set[project.WatcherHandle]
	logger         *project.Logger
	projectService *project.Service
	converters     *ls.Converters
}

// FS implements project.ServiceHost.
func (s *Server) FS() vfs.FS {
	return s.fs
}

// DefaultLibraryPath implements project.ServiceHost.
func (s *Server) DefaultLibraryPath() string {
	return s.defaultLibraryPath
}

// GetCurrentDirectory implements project.ServiceHost.
func (s *Server) GetCurrentDirectory() string {
	return s.cwd
}

// NewLine implements project.ServiceHost.
func (s *Server) NewLine() string {
	return s.newLine.GetNewLineCharacter()
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
func (s *Server) WatchFiles(watchers []*lsproto.FileSystemWatcher) (project.WatcherHandle, error) {
	watcherId := fmt.Sprintf("watcher-%d", s.watcherID)
	respChan, err := s.sendRequest(lsproto.MethodClientRegisterCapability, &lsproto.RegistrationParams{
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

	// TODO: timeout?
	resp := <-respChan
	if resp.Error != nil {
		return "", fmt.Errorf("failed to register file watcher: %s", resp.Error.String())
	}

	handle := project.WatcherHandle(watcherId)
	s.watchers.Add(handle)
	s.watcherID++
	return handle, nil
}

// UnwatchFiles implements project.Client.
func (s *Server) UnwatchFiles(handle project.WatcherHandle) error {
	if s.watchers.Has(handle) {
		respChan, err := s.sendRequest(lsproto.MethodClientUnregisterCapability, &lsproto.UnregistrationParams{
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

		resp := <-respChan
		if resp.Error != nil {
			return fmt.Errorf("failed to unregister file watcher: %s", resp.Error.String())
		}

		s.watchers.Delete(handle)
		return nil
	}

	return fmt.Errorf("no file watcher exists with ID %s", handle)
}

// RefreshDiagnostics implements project.Client.
func (s *Server) RefreshDiagnostics() error {
	if ptrIsTrue(s.initializeParams.Capabilities.Workspace.Diagnostics.RefreshSupport) {
		if err := s.sendRequest(lsproto.MethodWorkspaceDiagnosticRefresh, nil); err != nil {
			return fmt.Errorf("failed to refresh diagnostics: %w", err)
		}
	}
	return nil
}

func (s *Server) Run() error {
	go s.dispatchLoop()
	go s.writeLoop()
	return s.readLoop()
}

func (s *Server) readLoop() error {
	for {
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
				go s.cancelRequest(req.Params.(*lsproto.CancelParams).Id)
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
	data, err := s.r.Read()
	if err != nil {
		return nil, err
	}

	req := &lsproto.Message{}
	if err := json.Unmarshal(data, req); err != nil {
		return nil, fmt.Errorf("%w: %w", lsproto.ErrInvalidRequest, err)
	}

	return req, nil
}

func (s *Server) dispatchLoop() {
	for req := range s.requestQueue {
		ctx := context.Background()

		if req.ID != nil {
			var cancel context.CancelFunc
			ctx, cancel = context.WithCancel(ctx)
			s.pendingClientRequestsMu.Lock()
			s.pendingClientRequests[*req.ID] = pendingClientRequest{
				req:    req,
				cancel: cancel,
			}
			s.pendingClientRequestsMu.Unlock()
		}

		if err := s.handleRequestOrNotification(ctx, req); err != nil {
			s.fatalErrChan <- err
			return
		}
	}
}

func (s *Server) writeLoop() {
	for msg := range s.outgoingQueue {
		data, err := json.Marshal(msg)
		if err != nil {
			s.fatalErrChan <- fmt.Errorf("failed to marshal message: %w", err)
			continue
		}
		if err := s.w.Write(data); err != nil {
			s.fatalErrChan <- fmt.Errorf("failed to write message: %w", err)
			continue
		}
	}
}

func (s *Server) sendRequest(method lsproto.Method, params any) (<-chan *lsproto.ResponseMessage, error) {
	s.clientSeq++
	id := lsproto.NewIDString(fmt.Sprintf("ts%d", s.clientSeq))
	req := lsproto.NewRequestMessage(method, id, params)

	responseChan := make(chan *lsproto.ResponseMessage, 1)
	s.pendingServerRequestsMu.Lock()
	s.pendingServerRequests[*id] = responseChan
	s.pendingServerRequestsMu.Unlock()

	s.outgoingQueue <- req.Message()
	return responseChan, nil
}

func (s *Server) sendNotification(method lsproto.Method, params any) error {
	req := lsproto.NewRequestMessage(method, nil /*id*/, params)
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return s.w.Write(data)
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
		s.handleInitialized(ctx, req)
	case *lsproto.DidOpenTextDocumentParams:
		s.handleDidOpen(ctx, req)
	case *lsproto.DidChangeTextDocumentParams:
		s.handleDidChange(ctx, req)
	case *lsproto.DidSaveTextDocumentParams:
		s.handleDidSave(ctx, req)
	case *lsproto.DidCloseTextDocumentParams:
		s.handleDidClose(ctx, req)
	case *lsproto.DidChangeWatchedFilesParams:
		s.handleDidChangeWatchedFiles(ctx, req)
	case *lsproto.DocumentDiagnosticParams:
		s.handleDocumentDiagnostic(ctx, req)
	case *lsproto.HoverParams:
		s.handleHover(ctx, req)
	case *lsproto.DefinitionParams:
		s.handleDefinition(ctx, req)
	case *lsproto.CompletionParams:
		s.handleCompletion(ctx, req)
	default:
		switch req.Method {
		case lsproto.MethodShutdown:
			s.projectService.Close()
			s.sendResult(req.ID, nil)
			return nil
		case lsproto.MethodExit:
			return nil
		default:
			s.Log("unknown method", req.Method)
			if req.ID != nil {
				s.sendError(req.ID, lsproto.ErrInvalidRequest)
			}
			return nil
		}
	}
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
			Version: ptrTo(core.Version),
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
			DiagnosticProvider: &lsproto.DiagnosticOptionsOrDiagnosticRegistrationOptions{
				DiagnosticOptions: &lsproto.DiagnosticOptions{
					InterFileDependencies: true,
				},
			},
			CompletionProvider: &lsproto.CompletionOptions{
				TriggerCharacters: &ls.TriggerCharacters,
				// !!! other options
			},
		},
	})
}

func (s *Server) handleInitialized(ctx context.Context, req *lsproto.RequestMessage) error {
	if s.initializeParams.Capabilities.Workspace.DidChangeWatchedFiles != nil && *s.initializeParams.Capabilities.Workspace.DidChangeWatchedFiles.DynamicRegistration {
		s.watchEnabled = true
	}

	s.logger = project.NewLogger([]io.Writer{s.stderr}, "" /*file*/, project.LogLevelVerbose)
	s.projectService = project.NewService(s, project.ServiceOptions{
		Logger:           s.logger,
		WatchEnabled:     s.watchEnabled,
		PositionEncoding: s.positionEncoding,
	})

	s.converters = ls.NewConverters(s.positionEncoding, func(fileName string) ls.ScriptInfo {
		return s.projectService.GetScriptInfo(fileName)
	})

	return nil
}

func (s *Server) handleDidOpen(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DidOpenTextDocumentParams)
	s.projectService.OpenFile(ls.DocumentURIToFileName(params.TextDocument.Uri), params.TextDocument.Text, ls.LanguageKindToScriptKind(params.TextDocument.LanguageId), "")
	return nil
}

func (s *Server) handleDidChange(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DidChangeTextDocumentParams)
	scriptInfo := s.projectService.GetScriptInfo(ls.DocumentURIToFileName(params.TextDocument.Uri))
	if scriptInfo == nil {
		s.sendError(req.ID, lsproto.ErrRequestFailed)
		return nil
	}

	changes := make([]ls.TextChange, len(params.ContentChanges))
	for i, change := range params.ContentChanges {
		if partialChange := change.TextDocumentContentChangePartial; partialChange != nil {
			if textChange, err := s.converters.FromLSPTextChange(partialChange, scriptInfo.FileName()); err != nil {
				s.sendError(req.ID, err)
				return nil
			} else {
				changes[i] = textChange
			}
		} else if wholeChange := change.TextDocumentContentChangeWholeDocument; wholeChange != nil {
			changes[i] = ls.TextChange{
				TextRange: core.NewTextRange(0, len(scriptInfo.Text())),
				NewText:   wholeChange.Text,
			}
		} else {
			s.sendError(req.ID, lsproto.ErrInvalidRequest)
			return nil
		}
	}

	s.projectService.ChangeFile(ls.DocumentURIToFileName(params.TextDocument.Uri), changes)
	return nil
}

func (s *Server) handleDidSave(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DidSaveTextDocumentParams)
	s.projectService.MarkFileSaved(ls.DocumentURIToFileName(params.TextDocument.Uri), *params.Text)
	return nil
}

func (s *Server) handleDidClose(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DidCloseTextDocumentParams)
	s.projectService.CloseFile(ls.DocumentURIToFileName(params.TextDocument.Uri))
	return nil
}

func (s *Server) handleDidChangeWatchedFiles(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DidChangeWatchedFilesParams)
	return s.projectService.OnWatchedFilesChanged(params.Changes)
}

func (s *Server) handleDocumentDiagnostic(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DocumentDiagnosticParams)
	file, project := s.getFileAndProject(params.TextDocument.Uri)
	diagnostics := project.LanguageService().GetDocumentDiagnostics(file.FileName())
	lspDiagnostics := make([]*lsproto.Diagnostic, len(diagnostics))
	for i, diag := range diagnostics {
		if lspDiagnostic, err := s.converters.ToLSPDiagnostic(diag); err != nil {
			s.sendError(req.ID, err)
			return nil
		} else {
			lspDiagnostics[i] = lspDiagnostic
		}
	}
	s.sendResult(req.ID, &lsproto.DocumentDiagnosticReport{
		RelatedFullDocumentDiagnosticReport: &lsproto.RelatedFullDocumentDiagnosticReport{
			FullDocumentDiagnosticReport: lsproto.FullDocumentDiagnosticReport{
				Kind:  lsproto.StringLiteralFull{},
				Items: lspDiagnostics,
			},
		},
	})
	return nil
}

func (s *Server) handleHover(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.HoverParams)
	file, project := s.getFileAndProject(params.TextDocument.Uri)
	pos, err := s.converters.LineAndCharacterToPositionForFile(params.Position, file.FileName())
	if err != nil {
		s.sendError(req.ID, err)
		return nil
	}

	hoverText := project.LanguageService().ProvideHover(file.FileName(), pos)
	s.sendResult(req.ID, &lsproto.Hover{
		Contents: lsproto.MarkupContentOrMarkedStringOrMarkedStrings{
			MarkupContent: &lsproto.MarkupContent{
				Kind:  lsproto.MarkupKindMarkdown,
				Value: codeFence("ts", hoverText),
			},
		},
	})

	return nil
}

func (s *Server) handleDefinition(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DefinitionParams)
	file, project := s.getFileAndProject(params.TextDocument.Uri)
	pos, err := s.converters.LineAndCharacterToPositionForFile(params.Position, file.FileName())
	if err != nil {
		s.sendError(req.ID, err)
		return nil
	}

	locations := project.LanguageService().ProvideDefinitions(file.FileName(), pos)
	lspLocations := make([]lsproto.Location, len(locations))
	for i, loc := range locations {
		if lspLocation, err := s.converters.ToLSPLocation(loc); err != nil {
			s.sendError(req.ID, err)
			return nil
		} else {
			lspLocations[i] = lspLocation
		}
	}

	s.sendResult(req.ID, &lsproto.Definition{Locations: &lspLocations})
	return nil
}

func (s *Server) handleCompletion(ctx context.Context, req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.CompletionParams)
	file, project := s.getFileAndProject(params.TextDocument.Uri)
	pos, err := s.converters.LineAndCharacterToPositionForFile(params.Position, file.FileName())
	if err != nil {
		s.sendError(req.ID, err)
		return nil
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
	list := project.LanguageService().ProvideCompletion(
		file.FileName(),
		pos,
		params.Context,
		s.initializeParams.Capabilities.TextDocument.Completion,
		&ls.UserPreferences{})
	s.sendResult(req.ID, list)
	return nil
}

func (s *Server) getFileAndProject(uri lsproto.DocumentUri) (*project.ScriptInfo, *project.Project) {
	fileName := ls.DocumentURIToFileName(uri)
	return s.projectService.EnsureDefaultProjectForFile(fileName)
}

func (s *Server) Log(msg ...any) {
	fmt.Fprintln(s.stderr, msg...)
}

func codeFence(lang string, code string) string {
	if code == "" {
		return ""
	}
	ticks := 3
	for strings.Contains(code, strings.Repeat("`", ticks)) {
		ticks++
	}
	var result strings.Builder
	result.Grow(len(code) + len(lang) + 2*ticks + 2)
	for range ticks {
		result.WriteByte('`')
	}
	result.WriteString(lang)
	result.WriteByte('\n')
	result.WriteString(code)
	result.WriteByte('\n')
	for range ticks {
		result.WriteByte('`')
	}
	return result.String()
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
