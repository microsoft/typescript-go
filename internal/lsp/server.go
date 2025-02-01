package lsp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type MainOptions struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	FS                 vfs.FS
	CurrentDirectory   string
	DefaultLibraryPath string
}

func Main(opts *MainOptions) int {
	stderrJSON := json.NewEncoder(opts.Stderr)
	stderrJSON.SetIndent("", "    ")
	stderrJSON.SetEscapeHTML(false)

	s := &server{
		r:                  lsproto.NewBaseReader(opts.Stdin),
		w:                  lsproto.NewBaseWriter(opts.Stdout),
		stderr:             opts.Stderr,
		stderrJSON:         stderrJSON,
		fs:                 opts.FS,
		currentDirectory:   opts.CurrentDirectory,
		defaultLibraryPath: opts.DefaultLibraryPath,
	}

	if err := s.run(); err != nil && !errors.Is(err, io.EOF) {
		return 1
	}
	return 0
}

type server struct {
	r *lsproto.BaseReader
	w *lsproto.BaseWriter

	stderr     io.Writer
	stderrJSON *json.Encoder

	fs                 vfs.FS
	currentDirectory   string
	defaultLibraryPath string

	initializeParams *lsproto.InitializeParams
}

func (s *server) run() error {
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

func (s *server) read() (*lsproto.RequestMessage, error) {
	data, err := s.r.Read()
	if err != nil {
		return nil, err
	}

	req := &lsproto.RequestMessage{}
	if err := json.Unmarshal(data, req); err != nil {
		return nil, fmt.Errorf("%w: %w", lsproto.ErrInvalidRequest, err)
	}

	// TODO(jakebailey): temporary debug logging
	if _, err := s.stderr.Write([]byte("REQUEST\n")); err != nil {
		return nil, err
	}
	if err := s.stderrJSON.Encode(req); err != nil {
		return nil, err
	}

	return req, err
}

func (s *server) sendResult(id *lsproto.ID, result any) error {
	return s.sendResponse(&lsproto.ResponseMessage{
		ID:     id,
		Result: result,
	})
}

func (s *server) sendError(id *lsproto.ID, err error) error {
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

func (s *server) sendResponse(resp *lsproto.ResponseMessage) error {
	// TODO(jakebailey): temporary debug logging
	if _, err := s.stderr.Write([]byte("RESPONSE\n")); err != nil {
		return err
	}
	if err := s.stderrJSON.Encode(resp); err != nil {
		return err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	return s.w.Write(data)
}

func (s *server) handleInitialize(req *lsproto.RequestMessage) error {
	s.initializeParams = req.Params.(*lsproto.InitializeParams)
	return s.sendResult(req.ID, &lsproto.InitializeResult{
		ServerInfo: &lsproto.ServerInfo{
			Name:    "typescript-go",
			Version: core.Version,
		},
		Capabilities: map[string]any{
			"textDocumentSync": lsproto.TextDocumentSyncKindIncremental,
			"hoverProvider":    true,
		},
	})
}

func (s *server) handleMessage(req *lsproto.RequestMessage) error {
	params := req.Params
	switch params.(type) {
	case *lsproto.InitializeParams:
		return s.sendError(req.ID, lsproto.ErrInvalidRequest)
	case *lsproto.HoverParams:
		return s.sendResult(req.ID, &lsproto.Hover{
			Contents: lsproto.MarkupContent{
				Kind:  lsproto.MarkupKindPlaintext,
				Value: "It works!",
			},
		})
	default:
		if req.ID != nil {
			return s.sendError(req.ID, lsproto.ErrInvalidRequest)
		}
		return nil
	}
}
