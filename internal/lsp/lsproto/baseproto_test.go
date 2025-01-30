package lsproto_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"gotest.tools/v3/assert"
)

func TestBaseReader(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		data  []byte
		value any
		err   string
	}{
		{
			name: "empty",
			data: []byte("Content-Length: 0\r\n\r\n"),
			err:  "lsp: no content length",
		},
		{
			name: "early end",
			data: []byte("oops"),
			err:  "EOF",
		},
		{
			name: "negative length",
			data: []byte("Content-Length: -1\r\n\r\n"),
			err:  "lsp: invalid content length: negative value -1",
		},
		{
			name: "invalid content",
			data: []byte("Content-Length: 1\r\n\r\n{"),
			err:  "lsp: unmarshal content: unexpected end of JSON input",
		},
		{
			name:  "valid content",
			data:  []byte("Content-Length: 2\r\n\r\n{}"),
			value: map[string]any{},
		},
		{
			name:  "extra header values",
			data:  []byte("Content-Length: 2\r\nExtra: 1\r\n\r\n{}"),
			value: map[string]any{},
		},
		{
			name: "too long content length",
			data: []byte("Content-Length: 100\r\n\r\n{}"),
			err:  "lsp: read content: unexpected EOF",
		},
		{
			name: "missing content length",
			data: []byte("Content-Length: \r\n\r\n{}"),
			err:  "lsp: invalid content length: parse error: strconv.ParseInt: parsing \"\": invalid syntax",
		},
		{
			name: "invalid header",
			data: []byte("Nope\r\n\r\n{}"),
			err:  "lsp: invalid header: \"Nope\\r\\n\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := lsproto.NewBaseReader(bytes.NewReader(tt.data))

			var v any
			err := r.Read(&v)
			if tt.err != "" {
				assert.Error(t, err, tt.err)
			}
			assert.DeepEqual(t, v, tt.value)
		})
	}
}

func TestBaseReaderMultipleReads(t *testing.T) {
	t.Parallel()

	data := []byte(
		"Content-Length: 4\r\n\r\n1234" +
			"Content-Length: 2\r\n\r\n{}",
	)
	r := lsproto.NewBaseReader(bytes.NewReader(data))

	var v1 any
	err := r.Read(&v1)
	assert.NilError(t, err)
	assert.DeepEqual(t, v1, 1234.0)

	var v2 any
	err = r.Read(&v2)
	assert.NilError(t, err)
	assert.DeepEqual(t, v2, map[string]any{})

	var v3 any
	err = r.Read(&v3)
	assert.Error(t, err, "EOF")
}

func TestBaseReaderUnmarshalError(t *testing.T) {
	t.Parallel()

	data := []byte("Content-Length: 2\r\n\r\n{}")
	r := lsproto.NewBaseReader(bytes.NewReader(data))
	var v typeWithUnmarshalError
	err := r.Read(&v)
	assert.Error(t, err, "EOF")
}

type typeWithUnmarshalError struct{}

func (*typeWithUnmarshalError) UnmarshalJSON([]byte) error {
	return errors.New("test error")
}

func TestBaseReaderReadError(t *testing.T) {
	t.Parallel()

	r := lsproto.NewBaseReader(&errorReader{})
	var v any
	err := r.Read(&v)
	assert.Error(t, err, "lsp: read header: test error")
}

type errorReader struct{}

func (*errorReader) Read([]byte) (int, error) {
	return 0, errors.New("test error")
}

func TestBaseWriter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		value any
		data  []byte
	}{
		{
			name:  "empty",
			value: map[string]any{},
			data:  []byte("Content-Length: 2\r\n\r\n{}"),
		},
		{
			name: "bigger object",
			value: map[string]any{
				"key": "value",
			},
			data: []byte("Content-Length: 15\r\n\r\n{\"key\":\"value\"}"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var b bytes.Buffer
			w := lsproto.NewBaseWriter(&b)
			err := w.Write(tt.value)
			assert.NilError(t, err)
			assert.DeepEqual(t, b.Bytes(), tt.data)
		})
	}
}

func TestBaseWriterMarshalError(t *testing.T) {
	t.Parallel()

	var b bytes.Buffer
	w := lsproto.NewBaseWriter(&b)
	err := w.Write(&typeWithMarshalError{})
	assert.Error(t, err, "lsp: marshal: json: error calling MarshalJSON for type *lsp_test.typeWithMarshalError: test error")
}

type typeWithMarshalError struct{}

func (*typeWithMarshalError) MarshalJSON() ([]byte, error) {
	return nil, errors.New("test error")
}

func TestBaseWriterWriteError(t *testing.T) {
	t.Parallel()

	w := lsproto.NewBaseWriter(&errorWriter{})
	err := w.Write(map[string]any{})
	assert.Error(t, err, "lsp: write: test error")
}

type errorWriter struct{}

func (*errorWriter) Write([]byte) (int, error) {
	return 0, errors.New("test error")
}
