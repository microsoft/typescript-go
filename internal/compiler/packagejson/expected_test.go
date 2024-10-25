package packagejson_test

import (
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/microsoft/typescript-go/internal/compiler/packagejson"
	"gotest.tools/v3/assert"
)

func TestExpected(t *testing.T) {
	type packageJson struct {
		Name    packagejson.Expected[string] `json:"name"`
		Version packagejson.Expected[string] `json:"version"`
	}

	var p packageJson

	jsonString := `{
		"name": "test",
		"version": 2,
		"exports": null
	}`

	err := json.Unmarshal([]byte(jsonString), &p)
	assert.NilError(t, err)

	assert.Equal(t, p.Name.Valid, true)
	assert.Equal(t, p.Name.Value, "test")

	assert.Equal(t, p.Version.Valid, false)
	assert.Equal(t, p.Version.Value, "")
}
