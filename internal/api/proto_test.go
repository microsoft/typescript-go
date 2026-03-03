package api_test

import (
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/microsoft/typescript-go/internal/api"
	"gotest.tools/v3/assert"
)

func TestDocumentIdentifierUnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    string
		fileName string
		uri      string
		err      string
	}{
		{
			name:     "plain string",
			input:    `"foo.ts"`,
			fileName: "foo.ts",
		},
		{
			name:  "uri object",
			input: `{"uri":"file:///foo.ts"}`,
			uri:   "file:///foo.ts",
		},
		{
			name:  "uri object with unknown fields",
			input: `{"uri":"file:///foo.ts","extra":true}`,
			uri:   "file:///foo.ts",
		},
		{
			name:  "empty object",
			input: `{}`,
		},
		{
			name:  "invalid type",
			input: `42`,
			err:   "json: cannot unmarshal into Go api.DocumentIdentifier: DocumentIdentifier: expected string or object, got number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var d api.DocumentIdentifier
			err := json.Unmarshal([]byte(tt.input), &d)
			if tt.err != "" {
				assert.Error(t, err, tt.err)
				return
			}
			assert.NilError(t, err)
			assert.Equal(t, d.FileName, tt.fileName)
			assert.Equal(t, string(d.URI), tt.uri)
		})
	}
}
