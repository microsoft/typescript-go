package lsp_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/microsoft/typescript-go/internal/lsp"
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
			err:  "lsp: decode content: unexpected EOF",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := lsp.NewBaseReader(bytes.NewReader(tt.data))

			var v any
			err := r.Read(&v)
			if tt.err != "" {
				assert.Error(t, err, tt.err)
			}
			assert.DeepEqual(t, v, tt.value)
		})
	}
}

func TestBaseReaderUnmarshalError(t *testing.T) {
	t.Parallel()

	data := []byte("Content-Length: 2\r\n\r\n{}")
	r := lsp.NewBaseReader(bytes.NewReader(data))
	var v typeWithUnmarshalError
	err := r.Read(&v)
	assert.Error(t, err, "lsp: decode content: test error")
}

type typeWithUnmarshalError struct{}

func (*typeWithUnmarshalError) UnmarshalJSON([]byte) error {
	return errors.New("test error")
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
			w := lsp.NewBaseWriter(&b)
			err := w.Write(tt.value)
			assert.NilError(t, err)
			assert.DeepEqual(t, b.Bytes(), tt.data)
		})
	}
}

func TestBaseWriterMarshalError(t *testing.T) {
	t.Parallel()

	var b bytes.Buffer
	w := lsp.NewBaseWriter(&b)
	err := w.Write(&typeWithMarshalError{})
	assert.Error(t, err, "lsp: marshal: json: error calling MarshalJSON for type *lsp_test.typeWithMarshalError: test error")
}

type typeWithMarshalError struct{}

func (*typeWithMarshalError) MarshalJSON() ([]byte, error) {
	return nil, errors.New("test error")
}
