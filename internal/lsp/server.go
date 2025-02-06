package lsp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

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

	FS                 vfs.FS
	DefaultLibraryPath string
}

func NewServer(opts *ServerOptions) *Server {
	return &Server{
		r:                  lsproto.NewBaseReader(opts.In),
		w:                  lsproto.NewBaseWriter(opts.Out),
		stderr:             opts.Err,
		fs:                 opts.FS,
		defaultLibraryPath: opts.DefaultLibraryPath,
	}
}

type Server struct {
	r *lsproto.BaseReader
	w *lsproto.BaseWriter

	stderr io.Writer

	fs                 vfs.FS
	defaultLibraryPath string

	initializeParams *lsproto.InitializeParams

	projectService *project.ProjectService
}

// FS implements project.ProjectServiceHost.
func (s *Server) FS() vfs.FS {
	return s.fs
}

// GetCurrentDirectory implements project.ProjectServiceHost.
func (s *Server) GetCurrentDirectory() string {
	return "/"
}

// NewLine implements project.ProjectServiceHost.
func (s *Server) NewLine() string {
	return "\n"
}

// Trace implements project.ProjectServiceHost.
func (s *Server) Trace(msg string) {
	panic("unimplemented")
}

var _ project.ProjectServiceHost = (*Server)(nil)

func (s *Server) Run() error {
	for {
		req, err := s.read()
		if err != nil {
			if errors.Is(err, lsproto.ErrInvalidRequest) {
				if err := s.sendError(nil, err); err != nil {
					return err
				}
				continue
			}
			return err
		}

		if s.initializeParams == nil {
			if req.Method == lsproto.MethodInitialize {
				if err := s.handleInitialize(req); err != nil {
					return err
				}
			} else {
				if err := s.sendError(req.ID, lsproto.ErrServerNotInitialized); err != nil {
					return err
				}
			}
			continue
		}

		if err := s.handleMessage(req); err != nil {
			return err
		}
	}
}

func (s *Server) read() (*lsproto.RequestMessage, error) {
	data, err := s.r.Read()
	if err != nil {
		return nil, err
	}

	req := &lsproto.RequestMessage{}
	if err := json.Unmarshal(data, req); err != nil {
		return nil, fmt.Errorf("%w: %w", lsproto.ErrInvalidRequest, err)
	}

	return req, err
}

func (s *Server) sendResult(id *lsproto.ID, result any) error {
	return s.sendResponse(&lsproto.ResponseMessage{
		ID:     id,
		Result: result,
	})
}

func (s *Server) sendError(id *lsproto.ID, err error) error {
	code := lsproto.ErrInternalError.Code
	if errCode := (*lsproto.ErrorCode)(nil); errors.As(err, &errCode) {
		code = errCode.Code
	}
	// TODO(jakebailey): error data
	return s.sendResponse(&lsproto.ResponseMessage{
		ID: id,
		Error: &lsproto.ResponseError{
			Code:    code,
			Message: err.Error(),
		},
	})
}

func (s *Server) sendResponse(resp *lsproto.ResponseMessage) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	return s.w.Write(data)
}

func ptrTo[T any](v T) *T {
	return &v
}

func (s *Server) handleInitialize(req *lsproto.RequestMessage) error {
	s.initializeParams = req.Params.(*lsproto.InitializeParams)
	s.projectService = project.NewProjectService(s, project.ProjectServiceOptions{
		DefaultLibraryPath: s.defaultLibraryPath,
		Logger:             project.NewLogger([]io.Writer{s.stderr}, project.LogLevelVerbose),
	})
	return s.sendResult(req.ID, &lsproto.InitializeResult{
		ServerInfo: &lsproto.ServerInfo{
			Name:    "typescript-go",
			Version: ptrTo(core.Version),
		},
		Capabilities: lsproto.ServerCapabilities{
			TextDocumentSync: &lsproto.TextDocumentSyncOptionsOrTextDocumentSyncKind{
				TextDocumentSyncKind: ptrTo(lsproto.TextDocumentSyncKindIncremental),
			},
			HoverProvider: &lsproto.BooleanOrHoverOptions{
				Boolean: ptrTo(true),
			},
		},
	})
}

func (s *Server) handleDidOpen(req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DidOpenTextDocumentParams)
	s.projectService.OpenClientFile(strings.Replace(string(params.TextDocument.Uri), "file://", "", 1), params.TextDocument.Text, LanguageIDToScriptKind(params.TextDocument.LanguageId), "")
	return s.sendResult(req.ID, nil)
}

func (s *Server) handleDidChange(req *lsproto.RequestMessage) error {
	params := req.Params.(*lsproto.DidChangeTextDocumentParams)
	scriptInfo := s.projectService.GetScriptInfo(strings.Replace(string(params.TextDocument.Uri), "file://", "", 1))
	if scriptInfo == nil {
		return s.sendError(req.ID, lsproto.ErrRequestFailed)
	}

	changes := make([]ls.TextChange, len(params.ContentChanges))
	for i, change := range params.ContentChanges {
		if partialChange := change.TextDocumentContentChangePartial; partialChange != nil {
			changes[i] = ls.TextChange{
				TextRange: core.NewTextRange(
					LineAndCharacterToPosition(partialChange.Range.Start, scriptInfo.LineMap()),
					LineAndCharacterToPosition(partialChange.Range.End, scriptInfo.LineMap()),
				),
				NewText: partialChange.Text,
			}
		} else if wholeChange := change.TextDocumentContentChangeWholeDocument; wholeChange != nil {
			changes[i] = ls.TextChange{
				TextRange: core.NewTextRange(0, len(scriptInfo.Text())),
				NewText:   wholeChange.Text,
			}
		} else {
			return s.sendError(req.ID, lsproto.ErrInvalidRequest)
		}
	}

	s.projectService.ApplyChangesInOpenFiles(
		nil, /*openFiles*/
		[]project.ChangeFileArguments{{
			FileName: strings.Replace(string(params.TextDocument.Uri), "file://", "", 1),
			Changes:  changes,
		}},
		nil, /*closedFiles*/
	)

	return s.sendResult(req.ID, nil)
}

func (s *Server) handleMessage(req *lsproto.RequestMessage) error {
	params := req.Params
	switch params := params.(type) {
	case *lsproto.InitializeParams:
		return s.sendError(req.ID, lsproto.ErrInvalidRequest)
	case *lsproto.DidOpenTextDocumentParams:
		return s.handleDidOpen(req)
	case *lsproto.DidChangeTextDocumentParams:
		return s.handleDidChange(req)
	case *lsproto.HoverParams:
		file, project := s.GetFileAndProject(params.TextDocument.Uri)
		hoverText := project.LanguageService().ProvideHover(
			file.FileName(),
			LineAndCharacterToPosition(params.Position, file.LineMap()),
		)
		return s.sendResult(req.ID, &lsproto.Hover{
			Contents: lsproto.MarkupContentOrMarkedStringOrMarkedStrings{
				MarkupContent: &lsproto.MarkupContent{
					Kind:  lsproto.MarkupKindPlainText,
					Value: hoverText,
				},
			},
		})
	default:
		fmt.Fprintln(s.stderr, "unknown method", req.Method)
		if req.ID != nil {
			return s.sendError(req.ID, lsproto.ErrInvalidRequest)
		}
		return nil
	}
}

func (s *Server) GetFileAndProject(uri lsproto.DocumentUri) (*project.ScriptInfo, *project.Project) {
	fileName := strings.Replace(string(uri), "file://", "", 1)
	return s.projectService.EnsureDefaultProjectForFile(fileName)
}

func (s *Server) Log(msg ...any) {
	fmt.Fprintln(s.stderr, msg...)
}
