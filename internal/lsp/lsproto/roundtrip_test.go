package lsproto

import (
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/json"
	"gotest.tools/v3/assert"
)

// TestRoundTrip locks the generated codecs: every value must survive
// marshal -> unmarshal unchanged. This guards fidelity so codec changes
// (e.g. pruning or table-driving them) cannot silently corrupt the wire
// format. Cover a representative spread of shapes: required fields,
// nullable/non-nullable optionals, enums, slices, nested objects, and
// unions.
func TestRoundTrip(t *testing.T) {
	t.Parallel()

	roundTrip(t, "Range", &Range{
		Start: Position{Line: 1, Character: 2},
		End:   Position{Line: 3, Character: 4},
	})
	roundTrip(t, "TextEdit", &TextEdit{
		Range:   Range{Start: Position{Line: 1, Character: 2}, End: Position{Line: 3, Character: 4}},
		NewText: "hello",
	})
	roundTrip(t, "MarkupContent", &MarkupContent{Kind: MarkupKindMarkdown, Value: "**x**"})
	roundTrip(t, "DidChangeConfigurationParams object", &DidChangeConfigurationParams{
		Settings: map[string]any{"js/ts": map[string]any{"x": float64(1)}},
	})
	roundTrip(t, "DidChangeConfigurationParams null", &DidChangeConfigurationParams{Settings: nil})
	roundTrip(t, "CompletionItem", &CompletionItem{
		Label:            "pageXOffset",
		Kind:             new(CompletionItemKindField),
		SortText:         new("15"),
		InsertTextFormat: new(InsertTextFormatPlainText),
	})
	// StringOrTuple union (string arm and tuple arm).
	roundTrip(t, "ParameterInformation string label", &ParameterInformation{
		Label: StringOrTuple{String: new("p: number")},
	})
	roundTrip(t, "ParameterInformation tuple label", &ParameterInformation{
		Label: StringOrTuple{Tuple: &[2]uint32{0, 4}},
	})
}

func roundTrip[T any](t *testing.T, name string, v *T) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		t.Parallel()
		data, err := json.Marshal(v)
		assert.NilError(t, err)
		var got T
		assert.NilError(t, json.Unmarshal(data, &got))
		again, err := json.Marshal(&got)
		assert.NilError(t, err)
		assert.Equal(t, string(data), string(again), "re-marshal differs for %s", name)
	})
}

// TestStrictnessMissingRequired confirms required fields are still enforced;
// default reflective decoding would silently accept these.
func TestStrictnessMissingRequired(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		input  string
		target any
	}{
		{"TextEdit missing newText", `{"range":{"start":{"line":0,"character":0},"end":{"line":0,"character":1}}}`, new(TextEdit)},
		{"Range missing end", `{"start":{"line":0,"character":0}}`, new(Range)},
		{"Position missing character", `{"line":0}`, new(Position)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := json.Unmarshal([]byte(tt.input), tt.target)
			assert.ErrorContains(t, err, "missing required properties")
		})
	}
}

// TestStrictnessNotObject confirms a non-object where an object is required
// is rejected rather than coerced.
func TestStrictnessNotObject(t *testing.T) {
	t.Parallel()
	err := json.Unmarshal([]byte(`"oops"`), new(TextEdit))
	assert.Assert(t, err != nil)
	assert.Assert(t, strings.Contains(err.Error(), "object") || strings.Contains(err.Error(), "cannot unmarshal"))
}
