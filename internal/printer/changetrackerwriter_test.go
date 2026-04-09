package printer

import "testing"

func TestSetLastNonTriviaPositionEmptyString(t *testing.T) {
	t.Parallel()
	ctw := NewChangeTrackerWriter("\n", 4)
	// WriteLiteral calls setLastNonTriviaPosition with force=true.
	// This should not panic when called with an empty string.
	ctw.WriteLiteral("")
}

func TestSetLastNonTriviaPositionAllWhitespace(t *testing.T) {
	t.Parallel()
	ctw := NewChangeTrackerWriter("\n", 4)
	// WriteLiteral calls setLastNonTriviaPosition with force=true.
	// This should not panic when the string is all whitespace.
	ctw.WriteLiteral("   ")
}
