package stringutil

import (
	"strings"
	"unicode"
)

func ToLowerJS(str string) string {
	runes := []rune(str)
	var builder strings.Builder
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
	var builder strings.Builder
	for _, r := range str {
		if mapping, ok := specialCasingMappings[r]; ok {
			builder.WriteString(mapping.upper)
			continue
		}
		builder.WriteRune(unicode.ToUpper(r))
	}
	return builder.String()
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
	return IsInRuneRanges(r, unicodeCasedRanges)
}

func isUnicodeCaseIgnorable(r rune) bool {
	return IsInRuneRanges(r, unicodeCaseIgnorableRanges)
}
