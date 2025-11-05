package regexpchecker

import "unicode/utf8"

// utf16.go contains utilities for handling UTF-16 surrogate pairs and encoding.
// JavaScript regular expressions use UTF-16 internally, so we need to mimic this
// behavior when validating regex patterns. This includes handling surrogate pairs
// and preserving surrogate values that would otherwise be invalid in Go strings.

// UTF-16 surrogate pair constants
const (
	surrogateMin     = 0xD800  // Start of surrogate range
	surrogateMax     = 0xDFFF  // End of surrogate range
	highSurrogateMin = 0xD800  // Start of high surrogate range
	highSurrogateMax = 0xDBFF  // End of high surrogate range
	lowSurrogateMin  = 0xDC00  // Start of low surrogate range
	lowSurrogateMax  = 0xDFFF  // End of low surrogate range
	supplementaryMin = 0x10000 // First code point requiring surrogate pair
)

// isSurrogate returns true if the code point is in the surrogate range
func isSurrogate(r rune) bool {
	return r >= surrogateMin && r <= surrogateMax
}

// isHighSurrogate returns true if the code point is a high surrogate
func isHighSurrogate(r rune) bool {
	return r >= highSurrogateMin && r <= highSurrogateMax
}

// isLowSurrogate returns true if the code point is a low surrogate
func isLowSurrogate(r rune) bool {
	return r >= lowSurrogateMin && r <= lowSurrogateMax
}

// combineSurrogatePair combines a high and low surrogate into a code point
func combineSurrogatePair(high, low rune) rune {
	return supplementaryMin + ((high - highSurrogateMin) << 10) + (low - lowSurrogateMin)
}

// splitToSurrogatePair splits a supplementary code point into high and low surrogates
func splitToSurrogatePair(r rune) (high, low rune) {
	r -= supplementaryMin
	high = highSurrogateMin + ((r >> 10) & 0x3FF)
	low = lowSurrogateMin + (r & 0x3FF)
	return high, low
}

// encodeSurrogate encodes a UTF-16 surrogate value as a 2-byte UTF-16BE sequence.
// This preserves the surrogate value (which would otherwise be invalid in UTF-8/UTF-32).
// Go's string(rune) converts invalid surrogates to U+FFFD, so we use this encoding
// to preserve the exact surrogate value for JavaScript regex semantics.
func encodeSurrogate(surrogate rune) string {
	return string([]byte{byte(surrogate >> 8), byte(surrogate & 0xFF)})
}

// decodeSurrogate decodes a UTF-16BE encoded surrogate from a 2-byte string.
// Returns the surrogate value and true if successful, or 0 and false otherwise.
func decodeSurrogate(s string) (rune, bool) {
	if len(s) == 2 {
		code := rune((uint16(s[0]) << 8) | uint16(s[1]))
		if isSurrogate(code) {
			return code, true
		}
	}
	return 0, false
}

// decodeCodePoint returns the code point value from a character string.
// The string can be either a UTF-8 encoded character or a UTF-16BE encoded surrogate.
// Surrogates from escape sequences are encoded as 2-byte UTF-16BE sequences.
func decodeCodePoint(s string) rune {
	if len(s) == 0 {
		return 0
	}
	// Check if this is a UTF-16BE encoded surrogate
	if code, ok := decodeSurrogate(s); ok {
		return code
	}
	first, _ := utf8.DecodeRuneInString(s)
	return first
}

// charSize returns the number of UTF-16 code units needed to represent a code point.
// This matches JavaScript's internal string representation.
func charSize(ch rune) int {
	if ch >= supplementaryMin {
		// Code points >= 0x10000 require surrogate pairs in UTF-16 (2 code units)
		return 2
	}
	if ch == 0 {
		return 0
	}
	return 1
}

// utf16Length returns the UTF-16 length of a string, matching JavaScript's string.length.
// This counts UTF-16 code units, where surrogate pairs count as 2 units.
// Handles both UTF-8 encoded strings and special 2-byte UTF-16BE surrogate encodings.
func utf16Length(s string) int {
	// Check if this is a UTF-16BE surrogate encoding
	// These are used to preserve surrogate values in patterns like \uD835
	if _, ok := decodeSurrogate(s); ok {
		return 1
	}

	// Otherwise, count UTF-16 code units from UTF-8 runes
	length := 0
	for _, r := range s {
		length += charSize(r)
	}
	return length
}

// regExpChar represents a single "character" in a regex pattern.
// In Unicode mode, this is a single code point.
// In non-Unicode mode, this matches JavaScript's UTF-16 representation,
// where supplementary characters are represented as surrogate pairs.
type regExpChar struct {
	// The code point value. For surrogates, this is the surrogate value itself (0xD800-0xDFFF).
	codePoint rune
	// The UTF-16 length (1 for most characters, 2 for supplementary characters in Unicode mode)
	utf16Length int
}

// makeRegExpChar creates a regExpChar from a code point
func makeRegExpChar(codePoint rune) regExpChar {
	return regExpChar{
		codePoint:   codePoint,
		utf16Length: charSize(codePoint),
	}
}

// String returns a string representation suitable for comparison.
// For surrogates, uses UTF-16BE encoding to preserve the surrogate value.
// For normal code points, uses standard Go string encoding.
func (c regExpChar) String() string {
	if isSurrogate(c.codePoint) {
		return encodeSurrogate(c.codePoint)
	}
	return string(c.codePoint)
}
