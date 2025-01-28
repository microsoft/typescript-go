package lsp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
)

// https://microsoft.github.io/language-server-protocol/specifications/base/0.9/specification/

var (
	ErrInvalidHeader        = errors.New("lsp: invalid header")
	ErrInvalidContentLength = errors.New("lsp: invalid content length")
	ErrNoContentLength      = errors.New("lsp: no content length")
)

type BaseReader struct {
	r *bufio.Reader
}

func NewBaseReader(r io.Reader) *BaseReader {
	return &BaseReader{
		r: bufio.NewReader(r),
	}
}

func (r *BaseReader) Read(v any) error {
	var contentLength int64

	for {
		line, err := r.r.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return io.EOF
			}
			return fmt.Errorf("lsp: read header: %w", err)
		}

		if bytes.Equal(line, []byte("\r\n")) {
			break
		}

		key, value, ok := bytes.Cut(line, []byte(":"))
		if !ok {
			return fmt.Errorf("%w: %s", ErrInvalidHeader, line)
		}

		if bytes.Equal(key, []byte("Content-Length")) {
			contentLength, err = strconv.ParseInt(string(bytes.TrimSpace(value)), 10, 64)
			if err != nil {
				return fmt.Errorf("%w: parse error: %w", ErrInvalidContentLength, err)
			}
			if contentLength < 0 {
				return fmt.Errorf("%w: negative value %d", ErrInvalidContentLength, contentLength)
			}
		}
	}

	if contentLength <= 0 {
		return ErrNoContentLength
	}

	buf := make([]byte, contentLength)
	if _, err := io.ReadFull(r.r, buf); err != nil {
		return fmt.Errorf("lsp: read content: %w", err)
	}

	if err := json.Unmarshal(buf, v); err != nil {
		return fmt.Errorf("lsp: unmarshal content: %w", err)
	}

	return nil
}

type BaseWriter struct {
	w io.Writer
}

func NewBaseWriter(w io.Writer) *BaseWriter {
	return &BaseWriter{
		w: w,
	}
}

func (w *BaseWriter) Write(v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("lsp: marshal: %w", err)
	}

	if _, err := fmt.Fprintf(w.w, "Content-Length: %d\r\n\r\n%s", len(b), b); err != nil {
		return fmt.Errorf("lsp: write: %w", err)
	}

	return nil
}
