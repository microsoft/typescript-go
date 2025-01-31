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
	s := &server{
		r:                  lsproto.NewBaseReader(opts.Stdin),
		w:                  lsproto.NewBaseWriter(opts.Stdout),
		stderr:             opts.Stderr,
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

	stderr io.Writer

	fs                 vfs.FS
	currentDirectory   string
	defaultLibraryPath string

	initialize *lsproto.InitializeParams
}

func (s *server) run() error {
	enc := json.NewEncoder(s.stderr)

	for {
		var req lsproto.RequestMessage
		err := s.r.Read(&req)
		if err != nil {
			fmt.Fprintln(s.stderr, err)
			continue
		}

		// TODO(jakebailey): temporary debug logging
		enc.SetIndent("", "    ")
		enc.SetEscapeHTML(false)
		enc.Encode(req)

		if s.initialize == nil {
			if req.Method != lsproto.MethodInitialize {
				if err := s.sendResponseError(req.ID, lsproto.ErrServerNotInitialized); err != nil {
					// TODO(jakebailey): need to continue on error?
					return err
				}
				continue
			}
			s.initialize = req.Params.(*lsproto.InitializeParams)

			// TODO(jakebailey): handle initialize
		}

		// TODO(jakebailey): respond with cancellations for now
	}
}

func (s *server) sendResponse(id *lsproto.ID, result any) error {
	var resultPtr *any
	if result != nil {
		result = &result
	}
	m := &lsproto.ResponseMessage{
		ID:     id,
		Result: resultPtr,
	}
	return s.w.Write(m)
}

func (s *server) sendResponseError(id *lsproto.ID, err error) error {
	code := lsproto.ErrInternalError.Code
	if errCode := (*lsproto.ErrorCode)(nil); errors.As(err, &errCode) {
		code = errCode.Code
	}

	m := &lsproto.ResponseMessage{
		ID: id,
		Error: &lsproto.ResponseError{
			Code:    code,
			Message: err.Error(),
		},
	}
	return s.w.Write(m)
}
