package packagejson_test

import (
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/microsoft/typescript-go/internal/compiler/packagejson"
	"gotest.tools/v3/assert"
)

func TestNullable(t *testing.T) {
	type packageJson struct {
		Main    packagejson.Nullable[string]         `json:"main"`
		Types   packagejson.Nullable[string]         `json:"types"`
		Typings packagejson.Nullable[string]         `json:"typings"`
		Exports packagejson.Nullable[any]            `json:"exports"`
		Imports packagejson.Nullable[map[string]any] `json:"imports"`
	}

	var p packageJson

	jsonString := `{
		"main": null,
		"types": "test",
		"exports": null
	}`

	err := json.Unmarshal([]byte(jsonString), &p)
	assert.NilError(t, err)

	assert.Equal(t, p.Main.IsNull(), true)
	assert.Equal(t, p.Main.IsMissing(), false)
	assert.Equal(t, p.Main.IsPresent(), false)
	assert.Equal(t, p.Main.Value, "")

	assert.Equal(t, p.Types.IsNull(), false)
	assert.Equal(t, p.Types.IsMissing(), false)
	assert.Equal(t, p.Types.IsPresent(), true)
	assert.Equal(t, p.Types.Value, "test")

	assert.Equal(t, p.Typings.IsNull(), false)
	assert.Equal(t, p.Typings.IsMissing(), true)
	assert.Equal(t, p.Typings.IsPresent(), false)
	assert.Equal(t, p.Typings.Value, "")

	assert.Equal(t, p.Exports.IsNull(), true)
	assert.Equal(t, p.Exports.IsMissing(), false)
	assert.Equal(t, p.Exports.IsPresent(), false)
	assert.Equal(t, p.Exports.Value, nil)

	assert.ErrorContains(t, json.Unmarshal([]byte(`{"imports": 1}`), &p), "json: cannot unmarshal")
}
