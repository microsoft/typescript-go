package stringutil

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func ToLowerJS(str string) string {
	if ascii, ok := toLowerASCII(str); ok {
		return ascii
	}

	runes := []rune(str)
	var builder strings.Builder
	builder.Grow(len(str))
	for i, r := range runes {
		if mapping, ok := specialCasingMappings[r]; ok {
			if mapping.condition == specialCasingConditionFinalSigma && !isFinalSigmaContext(runes, i) {
				builder.WriteRune(unicode.ToLower(r))
				continue
			}
			builder.WriteString(mapping.lower)
			continue
		}
		builder.WriteRune(unicode.ToLower(r))
	}
	return builder.String()
}

func ToUpperJS(str string) string {
	if ascii, ok := toUpperASCII(str); ok {
		return ascii
	}

	var builder strings.Builder
	builder.Grow(len(str))
	for _, r := range str {
		if mapping, ok := specialCasingMappings[r]; ok {
			builder.WriteString(mapping.upper)
			continue
		}
		builder.WriteRune(unicode.ToUpper(r))
	}
	return builder.String()
}

func toLowerASCII(str string) (string, bool) {
	needsMapping := false
	for i := range len(str) {
		ch := str[i]
		if ch >= utf8.RuneSelf {
			return "", false
		}
		needsMapping = needsMapping || ('A' <= ch && ch <= 'Z')
	}
	if !needsMapping {
		return str, true
	}

	buf := []byte(str)
	for i, ch := range buf {
		if 'A' <= ch && ch <= 'Z' {
			buf[i] = ch + ('a' - 'A')
		}
	}
	return string(buf), true
}

func toUpperASCII(str string) (string, bool) {
	needsMapping := false
	for i := range len(str) {
		ch := str[i]
		if ch >= utf8.RuneSelf {
			return "", false
		}
		needsMapping = needsMapping || ('a' <= ch && ch <= 'z')
	}
	if !needsMapping {
		return str, true
	}

	buf := []byte(str)
	for i, ch := range buf {
		if 'a' <= ch && ch <= 'z' {
			buf[i] = ch - ('a' - 'A')
		}
	}
	return string(buf), true
}

func isFinalSigmaContext(runes []rune, index int) bool {
	return hasCasedLetterBefore(runes, index) && !hasCasedLetterAfter(runes, index)
}

func hasCasedLetterBefore(runes []rune, index int) bool {
	for i := index - 1; i >= 0; i-- {
		if isUnicodeCaseIgnorable(runes[i]) {
			continue
		}
		return isUnicodeCased(runes[i])
	}
	return false
}

func hasCasedLetterAfter(runes []rune, index int) bool {
	for i := index + 1; i < len(runes); i++ {
		if isUnicodeCaseIgnorable(runes[i]) {
			continue
		}
		return isUnicodeCased(runes[i])
	}
	return false
}

func isUnicodeCased(r rune) bool {
	return unicode.IsLower(r) || unicode.IsUpper(r) || unicode.IsTitle(r)
}

func isUnicodeCaseIgnorable(r rune) bool {
	return IsInRuneRanges(r, unicodeCaseIgnorableRanges)
}
