package stringutil

// IsUnicodeIdentifierStart reports whether ch may begin an ECMAScript
// identifier, i.e. whether it has the Unicode ID_Start (or Other_ID_Start)
// property. The range table is generated; see generate-unicode-data.mts.
func IsUnicodeIdentifierStart(ch rune) bool {
	return isInRuneRanges(ch, unicodeESNextIdentifierStart)
}

// IsUnicodeIdentifierPart reports whether ch may appear after the first
// character of an ECMAScript identifier, i.e. whether it has the Unicode
// ID_Continue (or Other_ID_Continue) property, which also includes ID_Start.
func IsUnicodeIdentifierPart(ch rune) bool {
	return isInRuneRanges(ch, unicodeESNextIdentifierPart)
}
