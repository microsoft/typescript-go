package logging

import (
	"testing"
)

// Verify LogTree implements the expected interface
type testLogger interface {
	Log(msg ...any)
	Write(msg string)
}

func TestLogTreeImplementsLogger(t *testing.T) {
	var _ testLogger = &LogTree{}
}

func TestLogTree(t *testing.T) {
}
