package testutil

import (
	"os"
	"runtime/debug"
	"sync"
	"testing"

	"github.com/microsoft/typescript-go/internal/testutil/race"
	"gotest.tools/v3/assert"
)

func AssertPanics(tb testing.TB, fn func(), expected any, msgAndArgs ...any) {
	tb.Helper()

	var got any

	func() {
		defer func() {
			got = recover()
		}()
		fn()
	}()

	assert.Assert(tb, got != nil, msgAndArgs...)
	assert.Equal(tb, got, expected, msgAndArgs...)
}

func RecoverAndFail(t *testing.T, msg string) {
	if r := recover(); r != nil {
		stack := debug.Stack()
		t.Fatalf("%s:\n%v\n%s", msg, r, string(stack))
	}
}

var testConcurrency = sync.OnceValue(func() string {
	// Leave Program in SingleThreaded mode unless explicitly configured or in race mode.
	if v := os.Getenv("TSGO_TEST_CONCURRENCY"); v != "" {
		return v
	}
	if race.Enabled {
		return "default"
	}
	return "single"
})

func TestConcurrency() string {
	return testConcurrency()
}
