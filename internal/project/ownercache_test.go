package project

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestOwnerCachePeek(t *testing.T) {
	t.Parallel()

	cache := NewOwnerCache(
		func(key string, value int) int { return value },
		nil,
	)

	cache.Acquire("key", 1, 42)
	value, ok := cache.Peek("key")
	assert.Assert(t, ok)
	assert.Equal(t, value, 42)

	cache.Release("key", 1)
	_, ok = cache.Peek("key")
	assert.Assert(t, !ok)
}
