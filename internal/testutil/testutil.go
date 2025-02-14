package testutil

import (
	"runtime/debug"
	"testing"

	"gotest.tools/v3/assert"
)

func AssertPanics(tb testing.TB, fn func(), expected any, msgAndArgs ...interface{}) {
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
