package checker

import (
	"strings"
	"sync"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// upperCaser and lowerCaser are lazily initialized via sync.OnceValue to avoid
// paying the cost of loading golang.org/x/text/cases tables unless non-ASCII
// input is actually encountered.
var (
	upperCaser = sync.OnceValue(func() *cases.Caser {
		c := cases.Upper(language.Und)
		return &c
	})
	lowerCaser = sync.OnceValue(func() *cases.Caser {
		c := cases.Lower(language.Und)
		return &c
	})
)

// isASCII reports whether s contains only ASCII bytes.
func isASCII(s string) bool {
	for i := range len(s) {
		if s[i] > 0x7F {
			return false
		}
	}
	return true
}

// toUpperCase converts a string to uppercase using full Unicode case mapping,
// matching JavaScript's String.prototype.toUpperCase() behavior.
// ASCII strings use a fast path; non-ASCII strings fall through to
// golang.org/x/text/cases which handles special case mappings where a single
// character maps to multiple characters (e.g., 'ß' → "SS").
func toUpperCase(s string) string {
	if isASCII(s) {
		return strings.ToUpper(s)
	}
	return upperCaser().String(s)
}

// toLowerCase converts a string to lowercase using full Unicode case mapping,
// matching JavaScript's String.prototype.toLowerCase() behavior.
func toLowerCase(s string) string {
	if isASCII(s) {
		return strings.ToLower(s)
	}
	return lowerCaser().String(s)
}

// toUpperCaseFirstRune converts the first rune to uppercase using full Unicode
// case mapping, leaving the rest of the string unchanged.
func toUpperCaseFirstRune(s string) string {
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	if r < 0x80 {
		// ASCII fast path: single byte uppercase
		if r >= 'a' && r <= 'z' {
			b := []byte(s)
			b[0] -= 'a' - 'A'
			return string(b)
		}
		return s
	}
	return upperCaser().String(s[:size]) + s[size:]
}

// toLowerCaseFirstRune converts the first rune to lowercase using full Unicode
// case mapping, leaving the rest of the string unchanged.
func toLowerCaseFirstRune(s string) string {
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	if r < 0x80 {
		// ASCII fast path: single byte lowercase
		if r >= 'A' && r <= 'Z' {
			b := []byte(s)
			b[0] += 'a' - 'A'
			return string(b)
		}
		return s
	}
	return lowerCaser().String(s[:size]) + s[size:]
}
