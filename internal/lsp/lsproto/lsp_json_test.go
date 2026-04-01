package lsproto

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/json"
	"gotest.tools/v3/assert"
)

func TestUnmarshalRejectsNullForOptionalNonNullableFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		target  any
		errText string
	}{
		{
			name:    "InlayHint kind null",
			input:   `{"position": {"line": 0, "character": 0}, "label": "foo", "kind": null}`,
			target:  new(InlayHint),
			errText: `null value is not allowed for field "kind"`,
		},
		{
			name:    "InlayHint textEdits null",
			input:   `{"position": {"line": 0, "character": 0}, "label": "foo", "textEdits": null}`,
			target:  new(InlayHint),
			errText: `null value is not allowed for field "textEdits"`,
		},
		{
			name:    "InlayHint paddingLeft null",
			input:   `{"position": {"line": 0, "character": 0}, "label": "foo", "paddingLeft": null}`,
			target:  new(InlayHint),
			errText: `null value is not allowed for field "paddingLeft"`,
		},
		{
			name:    "FoldingRange kind null",
			input:   `{"startLine": 0, "endLine": 10, "kind": null}`,
			target:  new(FoldingRange),
			errText: `null value is not allowed for field "kind"`,
		},
		{
			name:    "FoldingRange startCharacter null",
			input:   `{"startLine": 0, "endLine": 10, "startCharacter": null}`,
			target:  new(FoldingRange),
			errText: `null value is not allowed for field "startCharacter"`,
		},
		{
			name:    "CompletionItem insertTextFormat null",
			input:   `{"label": "test", "insertTextFormat": null}`,
			target:  new(CompletionItem),
			errText: `null value is not allowed for field "insertTextFormat"`,
		},
		{
			name:    "Hover range null",
			input:   `{"contents": {"kind": "plaintext", "value": "hi"}, "range": null}`,
			target:  new(Hover),
			errText: `null value is not allowed for field "range"`,
		},
		{
			name:    "WorkDoneProgressOptions workDoneProgress null",
			input:   `{"workDoneProgress": null}`,
			target:  new(WorkDoneProgressOptions),
			errText: `null value is not allowed for field "workDoneProgress"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := json.Unmarshal([]byte(tt.input), tt.target)
			assert.ErrorContains(t, err, tt.errText)
		})
	}
}

func TestUnmarshalAcceptsNullForNullableFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		target any
	}{
		{
			name:   "InitializeParams rootUri null",
			input:  `{"processId": null, "rootUri": null, "capabilities": {}}`,
			target: new(InitializeParams),
		},
		{
			name:   "InitializeParams workspaceFolders null",
			input:  `{"processId": null, "rootUri": null, "capabilities": {}, "workspaceFolders": null}`,
			target: new(InitializeParams),
		},
		{
			name:   "InitializeParams processId null",
			input:  `{"processId": null, "rootUri": null, "capabilities": {}}`,
			target: new(InitializeParams),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := json.Unmarshal([]byte(tt.input), tt.target)
			assert.NilError(t, err)
		})
	}
}

func TestUnmarshalAcceptsOmittedOptionalFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		target any
		check  func(t *testing.T, target any)
	}{
		{
			name:   "InlayHint with only required fields",
			input:  `{"position": {"line": 1, "character": 5}, "label": "test"}`,
			target: new(InlayHint),
			check: func(t *testing.T, target any) {
				t.Helper()
				hint := target.(*InlayHint)
				assert.Assert(t, hint.Kind == nil)
				assert.Assert(t, hint.TextEdits == nil)
				assert.Assert(t, hint.Tooltip == nil)
				assert.Assert(t, hint.PaddingLeft == nil)
				assert.Assert(t, hint.PaddingRight == nil)
				assert.Assert(t, hint.Data == nil)
				assert.Equal(t, hint.Position.Line, uint32(1))
				assert.Equal(t, hint.Position.Character, uint32(5))
			},
		},
		{
			name:   "FoldingRange with only required fields",
			input:  `{"startLine": 5, "endLine": 10}`,
			target: new(FoldingRange),
			check: func(t *testing.T, target any) {
				t.Helper()
				fr := target.(*FoldingRange)
				assert.Assert(t, fr.Kind == nil)
				assert.Assert(t, fr.StartCharacter == nil)
				assert.Assert(t, fr.EndCharacter == nil)
				assert.Assert(t, fr.CollapsedText == nil)
				assert.Equal(t, fr.StartLine, uint32(5))
				assert.Equal(t, fr.EndLine, uint32(10))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := json.Unmarshal([]byte(tt.input), tt.target)
			assert.NilError(t, err)
			tt.check(t, tt.target)
		})
	}
}

func TestUnmarshalRejectsIncompleteObjects(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		target  any
		errText string
	}{
		{
			name:    "InlayHint missing position",
			input:   `{"label": "test"}`,
			target:  new(InlayHint),
			errText: "missing required properties: position",
		},
		{
			name:    "InlayHint missing label",
			input:   `{"position": {"line": 0, "character": 0}}`,
			target:  new(InlayHint),
			errText: "missing required properties: label",
		},
		{
			name:    "Location missing uri",
			input:   `{"range": {"start": {"line": 0, "character": 0}, "end": {"line": 0, "character": 0}}}`,
			target:  new(Location),
			errText: "missing required properties: uri",
		},
		{
			name:    "Location empty object",
			input:   `{}`,
			target:  new(Location),
			errText: "missing required properties: uri, range",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := json.Unmarshal([]byte(tt.input), tt.target)
			assert.ErrorContains(t, err, tt.errText)
		})
	}
}

func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
	}{
		{
			name: "InlayHint with kind",
			value: &InlayHint{
				Position: Position{Line: 1, Character: 5},
				Label:    StringOrInlayHintLabelParts{String: new("param")},
				Kind:     new(InlayHintKindParameter),
			},
		},
		{
			name: "InlayHint minimal",
			value: &InlayHint{
				Position: Position{Line: 0, Character: 0},
				Label:    StringOrInlayHintLabelParts{String: new("x")},
			},
		},
		{
			name: "FoldingRange with all fields",
			value: &FoldingRange{
				StartLine:      1,
				StartCharacter: new(uint32(0)),
				EndLine:        10,
				EndCharacter:   new(uint32(5)),
				Kind:           new(FoldingRangeKindRegion),
				CollapsedText:  new("..."),
			},
		},
		{
			name: "Location",
			value: &Location{
				Uri: "file:///test.ts",
				Range: Range{
					Start: Position{Line: 1, Character: 2},
					End:   Position{Line: 3, Character: 4},
				},
			},
		},
		{
			name: "InitializeParams with null processId",
			value: &InitializeParams{
				ProcessId:    IntegerOrNull{},
				RootUri:      DocumentUriOrNull{DocumentUri: new(DocumentUri("file:///workspace"))},
				Capabilities: &ClientCapabilities{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data, err := json.Marshal(tt.value)
			assert.NilError(t, err)

			// Unmarshal into a new value of the same type
			switch v := tt.value.(type) {
			case *InlayHint:
				var result InlayHint
				err = json.Unmarshal(data, &result)
				assert.NilError(t, err)
				assert.DeepEqual(t, *v, result)
			case *FoldingRange:
				var result FoldingRange
				err = json.Unmarshal(data, &result)
				assert.NilError(t, err)
				assert.DeepEqual(t, *v, result)
			case *Location:
				var result Location
				err = json.Unmarshal(data, &result)
				assert.NilError(t, err)
				assert.DeepEqual(t, *v, result)
			case *InitializeParams:
				var result InitializeParams
				err = json.Unmarshal(data, &result)
				assert.NilError(t, err)
				assert.DeepEqual(t, *v, result)
			default:
				t.Fatalf("unhandled type %T", tt.value)
			}
		})
	}
}

func TestUnmarshalUnionTypes(t *testing.T) {
	t.Parallel()

	t.Run("IntegerOrString with integer", func(t *testing.T) {
		t.Parallel()
		var v IntegerOrString
		err := json.Unmarshal([]byte(`42`), &v)
		assert.NilError(t, err)
		assert.Assert(t, v.Integer != nil)
		assert.Equal(t, *v.Integer, int32(42))
		assert.Assert(t, v.String == nil)
	})

	t.Run("IntegerOrString with string", func(t *testing.T) {
		t.Parallel()
		var v IntegerOrString
		err := json.Unmarshal([]byte(`"hello"`), &v)
		assert.NilError(t, err)
		assert.Assert(t, v.String != nil)
		assert.Equal(t, *v.String, "hello")
		assert.Assert(t, v.Integer == nil)
	})

	t.Run("IntegerOrNull with integer", func(t *testing.T) {
		t.Parallel()
		var v IntegerOrNull
		err := json.Unmarshal([]byte(`42`), &v)
		assert.NilError(t, err)
		assert.Assert(t, v.Integer != nil)
		assert.Equal(t, *v.Integer, int32(42))
	})

	t.Run("IntegerOrNull with null", func(t *testing.T) {
		t.Parallel()
		var v IntegerOrNull
		err := json.Unmarshal([]byte(`null`), &v)
		assert.NilError(t, err)
		assert.Assert(t, v.Integer == nil)
	})

	t.Run("DocumentUriOrNull with string", func(t *testing.T) {
		t.Parallel()
		var v DocumentUriOrNull
		err := json.Unmarshal([]byte(`"file:///test.ts"`), &v)
		assert.NilError(t, err)
		assert.Assert(t, v.DocumentUri != nil)
		assert.Equal(t, *v.DocumentUri, DocumentUri("file:///test.ts"))
	})

	t.Run("DocumentUriOrNull with null", func(t *testing.T) {
		t.Parallel()
		var v DocumentUriOrNull
		err := json.Unmarshal([]byte(`null`), &v)
		assert.NilError(t, err)
		assert.Assert(t, v.DocumentUri == nil)
	})
}

func TestMarshalUnionTypes(t *testing.T) {
	t.Parallel()

	t.Run("IntegerOrNull with value", func(t *testing.T) {
		t.Parallel()
		v := IntegerOrNull{Integer: new(int32(42))}
		data, err := json.Marshal(&v)
		assert.NilError(t, err)
		assert.Equal(t, string(data), "42")
	})

	t.Run("IntegerOrNull with null", func(t *testing.T) {
		t.Parallel()
		v := IntegerOrNull{}
		data, err := json.Marshal(&v)
		assert.NilError(t, err)
		assert.Equal(t, string(data), "null")
	})

	t.Run("IntegerOrString with integer", func(t *testing.T) {
		t.Parallel()
		v := IntegerOrString{Integer: new(int32(7))}
		data, err := json.Marshal(&v)
		assert.NilError(t, err)
		assert.Equal(t, string(data), "7")
	})

	t.Run("IntegerOrString with string", func(t *testing.T) {
		t.Parallel()
		v := IntegerOrString{String: new("tok")}
		data, err := json.Marshal(&v)
		assert.NilError(t, err)
		assert.Equal(t, string(data), `"tok"`)
	})
}
