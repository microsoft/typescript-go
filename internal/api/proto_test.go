package api_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/api"
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/json"
	"github.com/microsoft/typescript-go/internal/parser"
	"gotest.tools/v3/assert"
)

func TestOrderedJSONValueUnmarshalJSON(t *testing.T) {
	t.Parallel()

	var value api.OrderedJSONValue
	err := json.Unmarshal([]byte(`{"z":1,"a":{"y":2,"x":3},"m":[{"b":4,"a":5}],"e":[]}`), &value)
	assert.NilError(t, err)

	root := value.Value.(*collections.OrderedMap[string, any])
	assert.DeepEqual(t, slices.Collect(root.Keys()), []string{"z", "a", "m", "e"})
	nested := root.GetOrZero("a").(*collections.OrderedMap[string, any])
	assert.DeepEqual(t, slices.Collect(nested.Keys()), []string{"y", "x"})
	array := root.GetOrZero("m").([]any)
	arrayObject := array[0].(*collections.OrderedMap[string, any])
	assert.DeepEqual(t, slices.Collect(arrayObject.Keys()), []string{"b", "a"})
	empty := root.GetOrZero("e").([]any)
	assert.Assert(t, empty != nil)
	assert.Equal(t, len(empty), 0)
}

func TestNewWatchOptionsResponsePreservesEmptyArrays(t *testing.T) {
	t.Parallel()

	response := api.NewWatchOptionsResponse(&core.WatchOptions{ExcludeFiles: []string{}})
	assert.Assert(t, response != nil)
	data, err := json.Marshal(response)
	assert.NilError(t, err)
	assert.Equal(t, string(data), `{"excludeFiles":[]}`)
}

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
			err:   "expected string or object, got number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var d api.DocumentIdentifier
			err := json.Unmarshal([]byte(tt.input), &d)
			if tt.err != "" {
				assert.ErrorContains(t, err, tt.err)
				return
			}
			assert.NilError(t, err)
			assert.Equal(t, d.FileName, tt.fileName)
			assert.Equal(t, string(d.URI), tt.uri)
		})
	}
}

func TestNewDiagnosticResponseUsesUTF16Offsets(t *testing.T) {
	t.Parallel()

	text := "const 💩 = 1;"
	file := parser.ParseSourceFile(ast.SourceFileParseOptions{FileName: "/unicode.ts"}, text, core.ScriptKindTS)
	pos := strings.Index(text, "=")
	assert.Assert(t, pos > 0)
	end := pos + len("=")

	diag := ast.NewDiagnostic(file, core.NewTextRange(pos, end), diagnostics.Expression_expected)
	resp := api.NewDiagnosticResponse(diag)

	assert.Equal(t, resp.Pos, 9)
	assert.Equal(t, resp.End, 10)
	assert.Equal(t, resp.Pos, file.GetPositionMap().UTF8ToUTF16(pos))
	assert.Equal(t, resp.End, file.GetPositionMap().UTF8ToUTF16(end))
}
