package core_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"gotest.tools/v3/assert"
)

func TestFindBestPatternMatch(t *testing.T) {
	t.Parallel()

	t.Run("exact match beats wildcard", func(t *testing.T) {
		t.Parallel()
		// When both an exact pattern and a wildcard pattern match the
		// candidate, the exact match must win regardless of iteration order.
		type entry struct {
			name    string
			pattern core.Pattern
		}
		values := []entry{
			{"wildcard", core.Pattern{Text: "src/foo/*", StarIndex: 8}},
			{"exact", core.Pattern{Text: "src/foo/bar", StarIndex: -1}},
		}
		result := core.FindBestPatternMatch(values, func(v entry) core.Pattern { return v.pattern }, "src/foo/bar")
		assert.Equal(t, result.name, "exact", "exact match must beat wildcard")
	})

	t.Run("exact match beats wildcard even when wildcard comes first", func(t *testing.T) {
		t.Parallel()
		type entry struct {
			name    string
			pattern core.Pattern
		}
		values := []entry{
			{"exact", core.Pattern{Text: "src/foo/bar", StarIndex: -1}},
			{"wildcard", core.Pattern{Text: "src/foo/*", StarIndex: 8}},
		}
		result := core.FindBestPatternMatch(values, func(v entry) core.Pattern { return v.pattern }, "src/foo/bar")
		assert.Equal(t, result.name, "exact", "exact match must beat wildcard regardless of order")
	})

	t.Run("longer prefix wildcard beats shorter prefix wildcard", func(t *testing.T) {
		t.Parallel()
		type entry struct {
			name    string
			pattern core.Pattern
		}
		values := []entry{
			{"short-prefix", core.Pattern{Text: "src/*", StarIndex: 4}},
			{"long-prefix", core.Pattern{Text: "src/foo/*", StarIndex: 8}},
		}
		result := core.FindBestPatternMatch(values, func(v entry) core.Pattern { return v.pattern }, "src/foo/bar")
		assert.Equal(t, result.name, "long-prefix", "longer prefix wildcard must win")
	})

	t.Run("no match returns zero value", func(t *testing.T) {
		t.Parallel()
		type entry struct {
			name    string
			pattern core.Pattern
		}
		values := []entry{
			{"a", core.Pattern{Text: "src/*", StarIndex: 4}},
			{"b", core.Pattern{Text: "lib/*", StarIndex: 4}},
		}
		result := core.FindBestPatternMatch(values, func(v entry) core.Pattern { return v.pattern }, "other/baz")
		assert.Equal(t, result.name, "", "no match should return zero value")
	})

	t.Run("single exact match", func(t *testing.T) {
		t.Parallel()
		type entry struct {
			name    string
			pattern core.Pattern
		}
		values := []entry{
			{"exact", core.Pattern{Text: "src/foo/bar", StarIndex: -1}},
		}
		result := core.FindBestPatternMatch(values, func(v entry) core.Pattern { return v.pattern }, "src/foo/bar")
		assert.Equal(t, result.name, "exact")
	})

	t.Run("single wildcard match", func(t *testing.T) {
		t.Parallel()
		type entry struct {
			name    string
			pattern core.Pattern
		}
		values := []entry{
			{"wildcard", core.Pattern{Text: "src/foo/*", StarIndex: 8}},
		}
		result := core.FindBestPatternMatch(values, func(v entry) core.Pattern { return v.pattern }, "src/foo/baz")
		assert.Equal(t, result.name, "wildcard")
	})

	t.Run("exact match with multiple wildcards", func(t *testing.T) {
		t.Parallel()
		type entry struct {
			name    string
			pattern core.Pattern
		}
		values := []entry{
			{"wc1", core.Pattern{Text: "src/*", StarIndex: 4}},
			{"wc2", core.Pattern{Text: "src/foo/*", StarIndex: 8}},
			{"exact", core.Pattern{Text: "src/foo/bar", StarIndex: -1}},
			{"wc3", core.Pattern{Text: "src/foo/bar/*", StarIndex: 12}},
		}
		result := core.FindBestPatternMatch(values, func(v entry) core.Pattern { return v.pattern }, "src/foo/bar")
		assert.Equal(t, result.name, "exact", "exact match must beat all wildcards")
	})
}
