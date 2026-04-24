package core

import (
	"encoding/json"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestCompilerOptionsJSONStrictArrayVariance(t *testing.T) {
	t.Parallel()

	t.Run("TSUnknown omits strictArrayVariance key", func(t *testing.T) {
		t.Parallel()
		o := CompilerOptions{Strict: TSTrue}
		b, err := json.Marshal(&o)
		assert.NilError(t, err)
		assert.Assert(t, !strings.Contains(string(b), "strictArrayVariance"), "json: %s", string(b))
	})

	t.Run("TSTrue encodes and roundtrips", func(t *testing.T) {
		t.Parallel()
		o := CompilerOptions{StrictArrayVariance: TSTrue}
		b, err := json.Marshal(&o)
		assert.NilError(t, err)
		assert.Assert(t, strings.Contains(string(b), `"strictArrayVariance":true`), "json: %s", string(b))
		var out CompilerOptions
		assert.NilError(t, json.Unmarshal(b, &out))
		assert.Equal(t, TSTrue, out.StrictArrayVariance)
	})

	t.Run("TSFalse encodes and roundtrips", func(t *testing.T) {
		t.Parallel()
		o := CompilerOptions{StrictArrayVariance: TSFalse}
		b, err := json.Marshal(&o)
		assert.NilError(t, err)
		assert.Assert(t, strings.Contains(string(b), `"strictArrayVariance":false`), "json: %s", string(b))
		var out CompilerOptions
		assert.NilError(t, json.Unmarshal(b, &out))
		assert.Equal(t, TSFalse, out.StrictArrayVariance)
	})
}
