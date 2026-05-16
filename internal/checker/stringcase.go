package checker

import (
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// upperCaser and lowerCaser use golang.org/x/text/cases with the "und" (undetermined)
// locale to perform full Unicode case mapping, matching JavaScript's
// String.prototype.toUpperCase() / toLowerCase() behavior.
// Unlike Go's strings.ToUpper/ToLower which use simple case mapping (1:1),
// these handle special case mappings where a single character maps to multiple
// characters (e.g., 'ß' → "SS", 'İ' → "i̇").
// The mapping tables come from Unicode's SpecialCasing.txt and are kept up to
// date via the golang.org/x/text module.
var (
	upperCaser = cases.Upper(language.Und)
	lowerCaser = cases.Lower(language.Und)
)

// toUpperCase converts a string to uppercase using the full Unicode case mapping,
// matching JavaScript's String.prototype.toUpperCase() behavior.
func toUpperCase(s string) string {
	return upperCaser.String(s)
}

// toLowerCase converts a string to lowercase using the full Unicode case mapping,
// matching JavaScript's String.prototype.toLowerCase() behavior.
func toLowerCase(s string) string {
	return lowerCaser.String(s)
}

// toUpperCaseFirstRune converts the first rune to uppercase using full Unicode
// case mapping, leaving the rest of the string unchanged.
func toUpperCaseFirstRune(s string) string {
	if s == "" {
		return s
	}
	_, size := utf8.DecodeRuneInString(s)
	return upperCaser.String(s[:size]) + s[size:]
}

// toLowerCaseFirstRune converts the first rune to lowercase using full Unicode
// case mapping, leaving the rest of the string unchanged.
func toLowerCaseFirstRune(s string) string {
	if s == "" {
		return s
	}
	_, size := utf8.DecodeRuneInString(s)
	return lowerCaser.String(s[:size]) + s[size:]
}
