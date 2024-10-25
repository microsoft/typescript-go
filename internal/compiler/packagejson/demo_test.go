package packagejson_test

import (
	"encoding/json"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestBuiltInStuff(t *testing.T) {
	someJson := `{ "name": "test", "version": "1.0.0" }`

	// We can deserialize into a struct...
	type MyJSON struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	var myJSON MyJSON
	err := json.Unmarshal([]byte(someJson), &myJSON)
	assert.NilError(t, err)

	// But if the types don't match, the whole thing fails
	someJson = `{ "name": "test", "version": 1.0 }`
	err = json.Unmarshal([]byte(someJson), &myJSON)
	assert.ErrorContains(t, err, "json: cannot unmarshal number into Go struct field MyJSON.version of type string")

	// Also, we can't tell the difference between a missing, null, and zero-value field this way
	someJson = `{ "name": "test", "exports": null }`
	type MyJSONWithExports struct {
		MyJSON
		Exports string `json:"exports"`
	}
	var myJSONWithExports MyJSONWithExports
	json.Unmarshal([]byte(someJson), &myJSONWithExports)
	assert.Equal(t, myJSONWithExports.Version, "")
	assert.Equal(t, myJSONWithExports.Exports, "")

	// We can deserialize into a map (or even an interface{})
	var myMap map[string]any
	json.Unmarshal([]byte(someJson), &myMap)
	exports, hasExports := myMap["exports"]
	assert.Assert(t, hasExports)
	assert.Equal(t, exports, nil)

	// But then we have to do type assertions or switches to get the values we want
	assert.Equal(t, strings.ToUpper(myMap["name"].(string)), "TEST")

	// Also, maps don't preserve order, which is important for package.json exports!
}

func TestLetsCustomizeStuff(t *testing.T) {
	someJson := `{ "private": "yes", "name": "test", "version": "1.0.0" }`

	type MyJSON struct {
		Private bool   `json:"private"`
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	var myJSON MyJSON
	err := json.Unmarshal([]byte(someJson), &myJSON)
	assert.NilError(t, err)
}
