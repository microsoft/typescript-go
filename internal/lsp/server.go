package lsp

import (
	"io"

	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

type Server struct {
	r *lsproto.BaseReader
	w *lsproto.BaseWriter
}

func NewServer(r io.Reader, w io.Writer) *Server {
	return &Server{
		r: lsproto.NewBaseReader(r),
		w: lsproto.NewBaseWriter(w),
	}
}
