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
	// ECMAScript points at Unicode Default Case Conversion for toLowerCase, and
	// modern V8 reaches that behavior through Intl::ConvertToLower, which uses
	// ICU root-locale lowercasing for non-Latin1 strings like Greek sigma.
	// We intentionally do not delegate this to golang.org/x/text/cases: x/text
	// is a general Unicode casing library, but its root-locale behavior is not
	// an exact match for the JS semantics exercised by String.prototype
	// .toLowerCase(), especially around Final_Sigma context. TypeScript needs the
	// JS behavior itself here, so we keep the context-sensitive part explicit.
	// SpiderMonkey models Final_Sigma with a more explicit context walk, while
	// Unicode Table 3-17 describes it in terms of Cased and Case_Ignorable.
	// We model the exposed V8/ICU behavior directly here: skip Case_Ignorable
	// code points and then look for lowercase/uppercase/titlecase code points,
	// including the DerivedCoreProperties Lowercase/Uppercase extras such as ª,
	// º, and Roman numerals.
	return hasSigmaCasedBefore(runes, index) && !hasSigmaCasedAfter(runes, index)
}

func hasSigmaCasedBefore(runes []rune, index int) bool {
	for i := index - 1; i >= 0; i-- {
		if isUnicodeCaseIgnorable(runes[i]) {
			continue
		}
		return isSigmaCased(runes[i])
	}
	return false
}

func hasSigmaCasedAfter(runes []rune, index int) bool {
	for i := index + 1; i < len(runes); i++ {
		if isUnicodeCaseIgnorable(runes[i]) {
			continue
		}
		return isSigmaCased(runes[i])
	}
	return false
}

func isSigmaCased(r rune) bool {
	return unicode.IsLower(r) || unicode.IsUpper(r) || unicode.IsTitle(r) ||
		IsInRuneRanges(r, unicodeLowercaseRanges) || IsInRuneRanges(r, unicodeUppercaseRanges)
}

func isUnicodeCaseIgnorable(r rune) bool {
	return IsInRuneRanges(r, unicodeCaseIgnorableRanges)
}
