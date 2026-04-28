package core

import (
	"os"
	"runtime/debug"
	"strconv"
)

// ApplyDebugStackLimit reads TSC_DEBUG_STACK_LIMIT and, if it parses as a
// positive integer, applies it via runtime/debug.SetMaxStack. Useful for
// catching runaway recursion. Should be called from program entry points.
func ApplyDebugStackLimit() {
	v := os.Getenv("TSC_DEBUG_STACK_LIMIT") //nolint:forbidigo
	if v == "" {
		return
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return
	}
	debug.SetMaxStack(n)
}
