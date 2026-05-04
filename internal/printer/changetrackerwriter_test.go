package printer

import (
	"testing"
)

func TestSetLastNonTriviaPositionEmptyString(t *testing.T) {
	t.Parallel()
	ctw := NewChangeTrackerWriter("\n", 4)
	// WriteLiteral calls setLastNonTriviaPosition(s, true) which would panic
	// on an empty string due to out-of-bounds access in the trailing whitespace loop.
	ctw.WriteLiteral("")
}

func TestSetLastNonTriviaPositionAllWhitespace(t *testing.T) {
	t.Parallel()
	ctw := NewChangeTrackerWriter("\n", 4)
	// Verify no panic when string is all whitespace and force=true.
	ctw.WriteLiteral("   ")
}
