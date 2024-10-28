package packagejson_test

import (
	"encoding/json"
	"testing"

	jsonExp "github.com/go-json-experiment/json"
	"github.com/microsoft/typescript-go/internal/compiler/packagejson"
	"gotest.tools/v3/assert"
)

func TestExports(t *testing.T) {
	testExportsWithUnmarshal(t, json.Unmarshal)
	testExportsWithUnmarshal(t, func(in []byte, out any) error { return jsonExp.Unmarshal(in, out) })
}

func testExportsWithUnmarshal(t *testing.T, unmarshal func([]byte, any) error) {
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

	err := unmarshal([]byte(jsonString), &e)
	assert.NilError(t, err)

	assert.Assert(t, e.Exports.IsSubpaths())
	assert.Equal(t, e.Exports.GetSubpaths().Size(), 2)
	assert.Assert(t, e.Exports.GetSubpaths().GetOrZero(".").IsConditions())
}
