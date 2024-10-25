package packagejson_test

import (
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/microsoft/typescript-go/internal/compiler/packagejson"
	"gotest.tools/v3/assert"
)

func TestJSONValue(t *testing.T) {
	type packageJson struct {
		Name    packagejson.JSONValue `json:"name"`
		Version packagejson.JSONValue `json:"version"`
	}

	var p packageJson

	jsonString := `{
		"name": "test",
		"version": 2
	}`

	err := json.Unmarshal([]byte(jsonString), &p)
	assert.NilError(t, err)

	assert.Equal(t, p.Name.Type, packagejson.JSONValueTypeString)
	assert.Equal(t, p.Name.Value, "test")

	assert.Equal(t, p.Version.Type, packagejson.JSONValueTypeNumber)
	assert.Equal(t, p.Version.Value, float64(2))
}
