package core

import (
	"testing"
)

func TestTrimPanicRecoveryFrames(t *testing.T) {
	t.Parallel()

	input := []byte(`goroutine 42 [running]:
runtime/debug.Stack()
	runtime/debug/stack.go:26 +0x5e
github.com/microsoft/typescript-go/internal/ls.handleCrossProject[...].func1.1()
	github.com/microsoft/typescript-go/internal/ls/crossproject.go:88 +0x70
panic({0xc323a0?, 0x1780b90?})
	runtime/panic.go:783 +0x132
github.com/microsoft/typescript-go/internal/checker.(*Checker).checkExpression(0xc0045a8000)
	github.com/microsoft/typescript-go/internal/checker/checker.go:5000 +0x1a0
github.com/microsoft/typescript-go/internal/ls.handleCrossProject[...].func1()
	github.com/microsoft/typescript-go/internal/ls/crossproject.go:105 +0x150`)

	expected := `github.com/microsoft/typescript-go/internal/checker.(*Checker).checkExpression(0xc0045a8000)
	github.com/microsoft/typescript-go/internal/checker/checker.go:5000 +0x1a0
github.com/microsoft/typescript-go/internal/ls.handleCrossProject[...].func1()
	github.com/microsoft/typescript-go/internal/ls/crossproject.go:105 +0x150`

	result := string(trimPanicRecoveryFrames(input))
	if result != expected {
		t.Errorf("trimPanicRecoveryFrames result mismatch.\nGot:\n%s\n\nExpected:\n%s", result, expected)
	}
}

func TestTrimPanicRecoveryFramesNoPanicFrame(t *testing.T) {
	t.Parallel()

	// If no panic() frame exists, the stack should be returned as-is.
	input := []byte(`github.com/microsoft/typescript-go/internal/checker.(*Checker).checkExpression(0xc0045a8000)
	github.com/microsoft/typescript-go/internal/checker/checker.go:5000 +0x1a0`)

	result := string(trimPanicRecoveryFrames(input))
	if result != string(input) {
		t.Errorf("expected unchanged output when no panic frame, got:\n%s", result)
	}
}
