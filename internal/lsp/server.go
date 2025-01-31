package lsp

import (
	"encoding/json"
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
}

func (s *server) run() error {
	for {
		var req lsproto.RequestMessage
		err := s.r.Read(&req)
		if err != nil {
			fmt.Fprintln(s.stderr, err)
			continue
		}

		// temporary debug logging
		enc := json.NewEncoder(s.stderr)
		enc.SetIndent("", "    ")
		enc.SetEscapeHTML(false)
		enc.Encode(req)
	}
}
