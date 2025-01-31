package lsp

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type server struct {
	r *lsproto.BaseReader
	w *lsproto.BaseWriter
}

func newServer(r io.Reader, w io.Writer) *server {
	return &server{
		r: lsproto.NewBaseReader(r),
		w: lsproto.NewBaseWriter(w),
	}
}

type MainOptions struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	FS                 vfs.FS
	CurrentDirectory   string
	DefaultLibraryPath string
}

func Main(opts *MainOptions) int {
	// TODO: ctx signal cancel exit
	server := newServer(opts.Stdin, opts.Stdout)

	for {
		var req lsproto.RequestMessage
		err := server.r.Read(&req)
		if err != nil {
			fmt.Fprintln(opts.Stderr, err)
			continue
		}

		enc := json.NewEncoder(opts.Stderr)
		enc.SetIndent("", "    ")
		enc.SetEscapeHTML(false)
		enc.Encode(req)
	}

	// return 0
}
