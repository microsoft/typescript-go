package checker

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// toUpperCase converts a string to uppercase using the full Unicode case mapping,
// matching JavaScript's String.prototype.toUpperCase() behavior. Unlike Go's
// strings.ToUpper which uses simple case mapping (1:1), this function handles
// special case mappings where a single character maps to multiple characters
// (e.g., 'ß' → "SS").
func toUpperCase(s string) string {
	// Fast path: check if any special casing characters exist
	hasSpecial := false
	for _, r := range s {
		if _, ok := upperSpecialCasings[r]; ok {
			hasSpecial = true
			break
		}
	}
	if !hasSpecial {
		return strings.ToUpper(s)
	}

	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if mapped, ok := upperSpecialCasings[r]; ok {
			b.WriteString(mapped)
		} else {
			b.WriteRune(unicode.ToUpper(r))
		}
	}
	return b.String()
}

// toLowerCase converts a string to lowercase using the full Unicode case mapping,
// matching JavaScript's String.prototype.toLowerCase() behavior.
func toLowerCase(s string) string {
	// Fast path: check if any special casing characters exist
	hasSpecial := false
	for _, r := range s {
		if _, ok := lowerSpecialCasings[r]; ok {
			hasSpecial = true
			break
		}
	}
	if !hasSpecial {
		return strings.ToLower(s)
	}

	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if mapped, ok := lowerSpecialCasings[r]; ok {
			b.WriteString(mapped)
		} else {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}

// toUpperCaseFirstRune converts the first rune to uppercase using full Unicode
// case mapping, returning the result and the byte size of the original first rune.
func toUpperCaseFirstRune(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError {
		return s
	}
	if mapped, ok := upperSpecialCasings[r]; ok {
		return mapped + s[size:]
	}
	return strings.ToUpper(s[:size]) + s[size:]
}

// toLowerCaseFirstRune converts the first rune to lowercase using full Unicode
// case mapping.
func toLowerCaseFirstRune(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError {
		return s
	}
	if mapped, ok := lowerSpecialCasings[r]; ok {
		return mapped + s[size:]
	}
	return strings.ToLower(s[:size]) + s[size:]
}

// upperSpecialCasings contains unconditional special case mappings for toUpperCase
// from Unicode SpecialCasing.txt. These are cases where a single code point maps
// to multiple code points when uppercased, matching JavaScript's behavior.
var upperSpecialCasings = map[rune]string{
	// Latin
	0x00DF: "SS",      // ß → SS
	0x0149: "\u02BCN", // ŉ → ʼN
	0x01F0: "J\u030C", // ǰ → J̌

	// Greek
	0x0390: "\u0399\u0308\u0301", // ΐ → Ϊ́
	0x03B0: "\u03A5\u0308\u0301", // ΰ → Ϋ́

	// Armenian
	0x0587: "\u0535\u0552", // և → ԵՒ

	// Latin extended
	0x1E96: "H\u0331", // ẖ → H̱
	0x1E97: "T\u0308", // ẗ → T̈
	0x1E98: "W\u030A", // ẘ → W̊
	0x1E99: "Y\u030A", // ẙ → Y̊
	0x1E9A: "A\u02BE", // ẚ → Aʾ

	// Greek extended
	0x1F50: "\u03A5\u0313",       // ὐ → Υ̓
	0x1F52: "\u03A5\u0313\u0300", // ὒ → Υ̓̀
	0x1F54: "\u03A5\u0313\u0301", // ὔ → Υ̓́
	0x1F56: "\u03A5\u0313\u0342", // ὖ → Υ̓͂

	// Greek extended - with iota subscript (prosgegrammeni)
	0x1F80: "\u1F08\u0399",       // ᾀ
	0x1F81: "\u1F09\u0399",       // ᾁ
	0x1F82: "\u1F0A\u0399",       // ᾂ
	0x1F83: "\u1F0B\u0399",       // ᾃ
	0x1F84: "\u1F0C\u0399",       // ᾄ
	0x1F85: "\u1F0D\u0399",       // ᾅ
	0x1F86: "\u1F0E\u0399",       // ᾆ
	0x1F87: "\u1F0F\u0399",       // ᾇ
	0x1F88: "\u1F08\u0399",       // ᾈ
	0x1F89: "\u1F09\u0399",       // ᾉ
	0x1F8A: "\u1F0A\u0399",       // ᾊ
	0x1F8B: "\u1F0B\u0399",       // ᾋ
	0x1F8C: "\u1F0C\u0399",       // ᾌ
	0x1F8D: "\u1F0D\u0399",       // ᾍ
	0x1F8E: "\u1F0E\u0399",       // ᾎ
	0x1F8F: "\u1F0F\u0399",       // ᾏ
	0x1F90: "\u1F28\u0399",       // ᾐ
	0x1F91: "\u1F29\u0399",       // ᾑ
	0x1F92: "\u1F2A\u0399",       // ᾒ
	0x1F93: "\u1F2B\u0399",       // ᾓ
	0x1F94: "\u1F2C\u0399",       // ᾔ
	0x1F95: "\u1F2D\u0399",       // ᾕ
	0x1F96: "\u1F2E\u0399",       // ᾖ
	0x1F97: "\u1F2F\u0399",       // ᾗ
	0x1F98: "\u1F28\u0399",       // ᾘ
	0x1F99: "\u1F29\u0399",       // ᾙ
	0x1F9A: "\u1F2A\u0399",       // ᾚ
	0x1F9B: "\u1F2B\u0399",       // ᾛ
	0x1F9C: "\u1F2C\u0399",       // ᾜ
	0x1F9D: "\u1F2D\u0399",       // ᾝ
	0x1F9E: "\u1F2E\u0399",       // ᾞ
	0x1F9F: "\u1F2F\u0399",       // ᾟ
	0x1FA0: "\u1F68\u0399",       // ᾠ
	0x1FA1: "\u1F69\u0399",       // ᾡ
	0x1FA2: "\u1F6A\u0399",       // ᾢ
	0x1FA3: "\u1F6B\u0399",       // ᾣ
	0x1FA4: "\u1F6C\u0399",       // ᾤ
	0x1FA5: "\u1F6D\u0399",       // ᾥ
	0x1FA6: "\u1F6E\u0399",       // ᾦ
	0x1FA7: "\u1F6F\u0399",       // ᾧ
	0x1FA8: "\u1F68\u0399",       // ᾨ
	0x1FA9: "\u1F69\u0399",       // ᾩ
	0x1FAA: "\u1F6A\u0399",       // ᾪ
	0x1FAB: "\u1F6B\u0399",       // ᾫ
	0x1FAC: "\u1F6C\u0399",       // ᾬ
	0x1FAD: "\u1F6D\u0399",       // ᾭ
	0x1FAE: "\u1F6E\u0399",       // ᾮ
	0x1FAF: "\u1F6F\u0399",       // ᾯ
	0x1FB2: "\u1FBA\u0399",       // ᾲ
	0x1FB3: "\u0391\u0399",       // ᾳ
	0x1FB4: "\u0386\u0399",       // ᾴ
	0x1FB6: "\u0391\u0342",       // ᾶ
	0x1FB7: "\u0391\u0342\u0399", // ᾷ
	0x1FBC: "\u0391\u0399",       // ᾼ
	0x1FC2: "\u1FCA\u0399",       // ῂ
	0x1FC3: "\u0397\u0399",       // ῃ
	0x1FC4: "\u0389\u0399",       // ῄ
	0x1FC6: "\u0397\u0342",       // ῆ
	0x1FC7: "\u0397\u0342\u0399", // ῇ
	0x1FCC: "\u0397\u0399",       // ῌ
	0x1FD2: "\u0399\u0308\u0300", // ῒ
	0x1FD3: "\u0399\u0308\u0301", // ΐ
	0x1FD6: "\u0399\u0342",       // ῖ
	0x1FD7: "\u0399\u0308\u0342", // ῗ
	0x1FE2: "\u03A5\u0308\u0300", // ῢ
	0x1FE3: "\u03A5\u0308\u0301", // ΰ
	0x1FE4: "\u03A1\u0313",       // ῤ
	0x1FE6: "\u03A5\u0342",       // ῦ
	0x1FE7: "\u03A5\u0308\u0342", // ῧ
	0x1FF2: "\u1FFA\u0399",       // ῲ
	0x1FF3: "\u03A9\u0399",       // ῳ
	0x1FF4: "\u038F\u0399",       // ῴ
	0x1FF6: "\u03A9\u0342",       // ῶ
	0x1FF7: "\u03A9\u0342\u0399", // ῷ
	0x1FFC: "\u03A9\u0399",       // ῼ

	// Latin ligatures
	0xFB00: "FF",  // ﬀ → FF
	0xFB01: "FI",  // ﬁ → FI
	0xFB02: "FL",  // ﬂ → FL
	0xFB03: "FFI", // ﬃ → FFI
	0xFB04: "FFL", // ﬄ → FFL
	0xFB05: "ST",  // ﬅ → ST
	0xFB06: "ST",  // ﬆ → ST

	// Armenian ligatures
	0xFB13: "\u0544\u0546", // ﬓ → ՄՆ
	0xFB14: "\u0544\u0535", // ﬔ → ՄԵ
	0xFB15: "\u0544\u053B", // ﬕ → ՄԻ
	0xFB16: "\u054E\u0546", // ﬖ → ՎՆ
	0xFB17: "\u0544\u053D", // ﬗ → ՄԽ
}

// lowerSpecialCasings contains unconditional special case mappings for toLowerCase
// from Unicode SpecialCasing.txt. There is only one unconditional special lower-case
// mapping in Unicode.
var lowerSpecialCasings = map[rune]string{
	0x0130: "i\u0307", // İ → i + combining dot above
}
