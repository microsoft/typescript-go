package core

import (
	"os"
	"runtime/debug"
	"strconv"
)

// DebugStackLimitEnvVar is the name of the environment variable that, when set
// to a positive integer, configures the Go runtime's per-goroutine maximum
// stack size via debug.SetMaxStack. This is intended for debugging runaway
// recursion: with a sufficiently low limit, well-behaved code still completes
// successfully but pathological recursion produces a fatal stack overflow
// instead of consuming gigabytes of memory before being killed.
const DebugStackLimitEnvVar = "TSC_DEBUG_STACK_LIMIT"

// ApplyDebugStackLimit checks the TSC_DEBUG_STACK_LIMIT environment variable
// and, if it is set to an integer greater than zero, applies it as the
// per-goroutine maximum stack size via runtime/debug.SetMaxStack. If the
// variable is unset, empty, not a valid integer, or not positive, no change is
// made and the Go runtime default applies.
//
// This should be called from program entry points (including TestMain).
func ApplyDebugStackLimit() {
	v := os.Getenv(DebugStackLimitEnvVar) //nolint:forbidigo
	if v == "" {
		return
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return
	}
	debug.SetMaxStack(n)
}
