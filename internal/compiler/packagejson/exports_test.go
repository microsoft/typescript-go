package packagejson_test

import (
	"encoding/json"
	"testing"

	json2 "github.com/go-json-experiment/json"
	"github.com/microsoft/typescript-go/internal/compiler/packagejson"
	"gotest.tools/v3/assert"
)

func TestExports(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		unmarshal func([]byte, any) error
	}{
		{"UnmarshalJSON", json.Unmarshal},
		{"UnmarshalJSONV2", func(in []byte, out any) error { return json2.Unmarshal(in, out) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			type Exports struct {
				Exports packagejson.Exports `json:"exports"`
			}

			var e Exports

			jsonString := `{
				"exports": {
					".": {
						"import": "./test.ts",
						"default": "./test.ts"
					},
					"./test": [
						"./test1.ts",
						"./test2.ts"
					]
				}
			}`

			err := tt.unmarshal([]byte(jsonString), &e)
			assert.NilError(t, err)

			assert.Assert(t, e.Exports.IsSubpaths())
			assert.Equal(t, e.Exports.AsObject().Size(), 2)
			assert.Assert(t, e.Exports.AsObject().GetOrZero(".").IsConditions())
		})
	}
}
