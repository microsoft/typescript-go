package lsp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

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

	if err := s.run(); err != nil {
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
		req := &lsproto.RequestMessage{}
		err := s.r.Read(req)

		// TODO(jakebailey): temporary debug logging
		if _, err := s.stderr.Write([]byte("REQUEST\n")); err != nil {
			return err
		}
		if err := s.stderrJSON.Encode(req); err != nil {
			return err
		}

		if err != nil {
			fmt.Fprintln(s.stderr, err)
			continue
		}

		if s.initializeParams == nil {
			if req.Method != lsproto.MethodInitialize {
				if err := s.sendError(req.ID, lsproto.ErrServerNotInitialized); err != nil {
					// TODO(jakebailey): need to continue on error?
					return err
				}
				continue
			}
			if err := s.handleInitialize(req); err != nil {
				return err
			}
		}

		// TODO(jakebailey): respond with cancellations for now
	}
}

func (s *server) sendResult(id *lsproto.ID, result any) error {
	var resultPtr *any
	if result != nil {
		resultPtr = &result
	}
	return s.sendResponse(&lsproto.ResponseMessage{
		ID:     id,
		Result: resultPtr,
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
	return s.w.Write(resp)
}

func (s *server) handleInitialize(req *lsproto.RequestMessage) error {
	s.initializeParams = req.Params.(*lsproto.InitializeParams)
	return s.sendResult(req.ID, &lsproto.InitializeResult{
		ServerInfo: &lsproto.ServerInfo{
			Name: "typescript-go",
			// Version: core.Version, // TODO(jakebailey): put version in package other than core
		},
		Capabilities: map[string]any{ // TODO(jakebailey): do something here
			"textDocumentSync": 1, // TextDocumentSyncKind.Full
		},
	})
}
