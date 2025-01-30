package lsp

import "io"

type Server struct {
	r *BaseReader
	w *BaseWriter
}

func NewServer(r io.Reader, w io.Writer) *Server {
	return &Server{
		r: NewBaseReader(r),
		w: NewBaseWriter(w),
	}
}
