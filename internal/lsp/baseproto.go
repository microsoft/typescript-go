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
			return fmt.Errorf("%w: %q", ErrInvalidHeader, line)
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

	// TODO: pool these buffers
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

type ErrorCode struct { //nolint:errname
	Name string
	Code Integer
}

func (e *ErrorCode) Error() string {
	return e.Name
}

var (
	// Defined by JSON-RPC
	ErrParseError     = &ErrorCode{"ParseError", -32700}
	ErrInvalidRequest = &ErrorCode{"InvalidRequest", -32600}
	ErrMethodNotFound = &ErrorCode{"MethodNotFound", -32601}
	ErrInvalidParams  = &ErrorCode{"InvalidParams", -32602}
	ErrInternalError  = &ErrorCode{"InternalError", -32603}

	// Error code indicating that a server received a notification or
	// request before the server has received the `initialize` request.
	ErrServerNotInitialized = &ErrorCode{"ServerNotInitialized", -32002}
	ErrUnknownErrorCode     = &ErrorCode{"UnknownErrorCode", -32001}

	// A request failed but it was syntactically correct, e.g the
	// method name was known and the parameters were valid. The error
	// message should contain human readable information about why
	// the request failed.
	ErrRequestFailed = &ErrorCode{"RequestFailed", -32803}

	// The server cancelled the request. This error code should
	// only be used for requests that explicitly support being
	// server cancellable.
	ErrServerCancelled = &ErrorCode{"ServerCancelled", -32802}

	// The server detected that the content of a document got
	// modified outside normal conditions. A server should
	// NOT send this error code if it detects a content change
	// in it unprocessed messages. The result even computed
	// on an older state might still be useful for the client.
	//
	// If a client decides that a result is not of any use anymore
	// the client should cancel the request.
	ErrContentModified = &ErrorCode{"ContentModified", -32801}

	// The client has canceled a request and a server has detected
	// the cancel.
	ErrRequestCancelled = &ErrorCode{"RequestCancelled", -32800}
)
