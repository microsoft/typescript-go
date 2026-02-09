package regexpchecker

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/scanner"
	"github.com/microsoft/typescript-go/internal/stringutil"
)

// regExpFlags represents regexp flags (e.g., 'g', 'i', 'm', etc.)
type regExpFlags uint32

const (
	regExpFlagsNone           regExpFlags = 0
	regExpFlagsGlobal         regExpFlags = 1 << 0 // g
	regExpFlagsIgnoreCase     regExpFlags = 1 << 1 // i
	regExpFlagsMultiline      regExpFlags = 1 << 2 // m
	regExpFlagsDotAll         regExpFlags = 1 << 3 // s
	regExpFlagsUnicode        regExpFlags = 1 << 4 // u
	regExpFlagsSticky         regExpFlags = 1 << 5 // y
	regExpFlagsHasIndices     regExpFlags = 1 << 6 // d
	regExpFlagsUnicodeSets    regExpFlags = 1 << 7 // v
	regExpFlagsModifiers      regExpFlags = regExpFlagsIgnoreCase | regExpFlagsMultiline | regExpFlagsDotAll
	regExpFlagsAnyUnicodeMode regExpFlags = regExpFlagsUnicode | regExpFlagsUnicodeSets
)

var charCodeToRegExpFlag = map[rune]regExpFlags{
	'd': regExpFlagsHasIndices,
	'g': regExpFlagsGlobal,
	'i': regExpFlagsIgnoreCase,
	'm': regExpFlagsMultiline,
	's': regExpFlagsDotAll,
	'u': regExpFlagsUnicode,
	'v': regExpFlagsUnicodeSets,
	'y': regExpFlagsSticky,
}

// regExpValidator is used to validate regular expressions
type regExpValidator struct {
	text                           string
	pos                            int
	end                            int
	languageVersion                core.ScriptTarget
	languageVariant                core.LanguageVariant
	onError                        scanner.ErrorCallback
	regExpFlags                    regExpFlags
	annexB                         bool
	unicodeSetsMode                bool
	anyUnicodeMode                 bool
	anyUnicodeModeOrNonAnnexB      bool
	namedCaptureGroups             bool
	numberOfCapturingGroups        int
	groupSpecifiers                map[string]bool
	groupNameReferences            []namedReference
	decimalEscapes                 []decimalEscape
	namedCapturingGroupsScopeStack []map[string]bool
	topNamedCapturingGroupsScope   map[string]bool
	mayContainStrings              bool
	isCharacterComplement          bool
	tokenValue                     string
	surrogateState                 *surrogatePairState // For non-Unicode mode: tracks partial surrogate pair
}

// surrogatePairState tracks when we're in the middle of emitting a surrogate pair
// in non-Unicode mode (where literal characters >= U+10000 must be split into two UTF-16 code units)
type surrogatePairState struct {
	lowSurrogate rune // The low surrogate value to return next
	utf8Size     int  // Size of the UTF-8 character to advance past
}

type namedReference struct {
	pos  int
	end  int
	name string
}

type decimalEscape struct {
	pos   int
	end   int
	value int
}

func Check(
	node *ast.RegularExpressionLiteral,
	sourceFile *ast.SourceFile,
	languageVersion core.ScriptTarget,
	onError scanner.ErrorCallback,
) {
	text := node.Text
	v := &regExpValidator{
		text:            text,
		languageVersion: languageVersion,
		languageVariant: sourceFile.LanguageVariant,
		onError:         onError,
	}

	// Similar to the original scanRegularExpressionWorker, but since we are outside the scanner,
	// we need to rescan for some information that the scanner previously calculated.

	bodyEnd := strings.LastIndexByte(text, '/')
	if bodyEnd <= 0 {
		panic("regexpchecker: regex must have closing '/' (scanner should have validated)")
	}

	v.pos = bodyEnd + 1
	v.end = len(text)
	v.regExpFlags = v.scanFlags(regExpFlagsNone, false)
	v.pos = 1
	v.end = bodyEnd

	v.unicodeSetsMode = v.regExpFlags&regExpFlagsUnicodeSets != 0
	v.anyUnicodeMode = v.regExpFlags&regExpFlagsAnyUnicodeMode != 0
	v.annexB = true
	v.anyUnicodeModeOrNonAnnexB = v.anyUnicodeMode || !v.annexB
	v.namedCaptureGroups = v.detectNamedCaptureGroups()

	v.scanDisjunction(false)
	v.validateGroupReferences()
	v.validateDecimalEscapes()
}

// detectNamedCaptureGroups performs a quick scan of the pattern to detect
// if it contains any named capture groups (?<name>...). This is needed because
// the presence of named groups changes the interpretation of \k escapes:
// - Without named groups: \k is an identity escape (matches literal 'k')
// - With named groups: \k must be followed by <name> or it's a syntax error
// This matches the behavior in scanner.ts's reScanSlashToken.
func (v *regExpValidator) detectNamedCaptureGroups() bool {
	inEscape := false
	inCharacterClass := false
	text := v.text[v.pos:v.end]

	for i, ch := range text {
		// Only check ASCII characters for the pattern (?<
		if ch >= 0x80 {
			continue
		}

		if inEscape {
			inEscape = false
		} else if ch == '\\' {
			inEscape = true
		} else if ch == '[' {
			inCharacterClass = true
		} else if ch == ']' {
			inCharacterClass = false
		} else if !inCharacterClass &&
			ch == '(' &&
			i+3 < len(text) &&
			text[i+1] == '?' &&
			text[i+2] == '<' &&
			text[i+3] != '=' &&
			text[i+3] != '!' {
			// Found (?< that's not (?<= or (?<! - this is a named capture group
			return true
		}
	}
	return false
}

func (v *regExpValidator) charAndSize() (rune, int) {
	if v.pos >= v.end {
		return 0, 0
	}
	// Simple ASCII fast path
	if ch := v.text[v.pos]; ch < 0x80 {
		return rune(ch), 1
	}
	// Decode multi-byte UTF-8 character
	r, size := utf8.DecodeRuneInString(v.text[v.pos:])
	return r, size
}

func (v *regExpValidator) charAtOffset(offset int) rune {
	if v.pos+offset >= v.end {
		return 0
	}
	// Simple ASCII fast path
	if ch := v.text[v.pos+offset]; ch < 0x80 {
		return rune(ch)
	}
	// Decode multi-byte UTF-8 character
	r, _ := utf8.DecodeRuneInString(v.text[v.pos+offset:])
	return r
}

func (v *regExpValidator) error(message *diagnostics.Message, start, length int, args ...any) {
	v.onError(message, start, length, args...)
}

func (v *regExpValidator) checkRegularExpressionFlagAvailability(flag regExpFlags, size int) {
	var availableFrom core.ScriptTarget
	switch flag {
	case regExpFlagsHasIndices:
		availableFrom = core.ScriptTargetES2022
	case regExpFlagsDotAll:
		availableFrom = core.ScriptTargetES2018
	case regExpFlagsUnicodeSets:
		availableFrom = core.ScriptTargetES2024
	default:
		return
	}

	if v.languageVersion < availableFrom {
		// !!! Old compiler lowercases these names.
		v.error(diagnostics.This_regular_expression_flag_is_only_available_when_targeting_0_or_later, v.pos, size, strings.ToLower(availableFrom.String()))
	}
}

func (v *regExpValidator) scanDisjunction(isInGroup bool) {
	for {
		v.namedCapturingGroupsScopeStack = append(v.namedCapturingGroupsScopeStack, v.topNamedCapturingGroupsScope)
		v.topNamedCapturingGroupsScope = nil
		v.scanAlternative(isInGroup)
		v.topNamedCapturingGroupsScope = v.namedCapturingGroupsScopeStack[len(v.namedCapturingGroupsScopeStack)-1]
		v.namedCapturingGroupsScopeStack = v.namedCapturingGroupsScopeStack[:len(v.namedCapturingGroupsScopeStack)-1]

		if v.charAtOffset(0) != '|' {
			return
		}
		v.pos++
	}
}

func (v *regExpValidator) scanAlternative(isInGroup bool) {
	isPreviousTermQuantifiable := false
	for {
		start := v.pos
		ch := v.charAtOffset(0)
		switch ch {
		case 0:
			return
		case '^', '$':
			v.pos++
			isPreviousTermQuantifiable = false
		case '\\':
			v.pos++
			switch v.charAtOffset(0) {
			case 'b', 'B':
				v.pos++
				isPreviousTermQuantifiable = false
			default:
				v.scanAtomEscape()
				isPreviousTermQuantifiable = true
			}
		case '(':
			v.pos++
			if v.charAtOffset(0) == '?' {
				v.pos++
				switch v.charAtOffset(0) {
				case '=', '!':
					v.pos++
					isPreviousTermQuantifiable = !v.anyUnicodeModeOrNonAnnexB
				case '<':
					groupNameStart := v.pos
					v.pos++
					switch v.charAtOffset(0) {
					case '=', '!':
						v.pos++
						isPreviousTermQuantifiable = false
					default:
						v.scanGroupName(false)
						v.scanExpectedChar('>')
						if v.languageVersion < core.ScriptTargetES2018 {
							v.error(diagnostics.Named_capturing_groups_are_only_available_when_targeting_ES2018_or_later, groupNameStart, v.pos-groupNameStart)
						}
						v.numberOfCapturingGroups++
						isPreviousTermQuantifiable = true
					}
				default:
					start := v.pos
					setFlags := v.scanPatternModifiers(regExpFlagsNone)
					if v.charAtOffset(0) == '-' {
						v.pos++
						v.scanPatternModifiers(setFlags)
						if v.pos == start+1 {
							v.error(diagnostics.Subpattern_flags_must_be_present_when_there_is_a_minus_sign, start, v.pos-start)
						}
					}
					v.scanExpectedChar(':')
					isPreviousTermQuantifiable = true
				}
			} else {
				v.numberOfCapturingGroups++
				isPreviousTermQuantifiable = true
			}
			v.scanDisjunction(true)
			v.scanExpectedChar(')')
		case '{':
			v.pos++
			digitsStart := v.pos
			v.scanDigits()
			minVal := v.tokenValue
			if !v.anyUnicodeModeOrNonAnnexB && minVal == "" {
				isPreviousTermQuantifiable = true
				break
			}
			if v.charAtOffset(0) == ',' {
				v.pos++
				v.scanDigits()
				maxVal := v.tokenValue
				if minVal == "" {
					if maxVal != "" || v.charAtOffset(0) == '}' {
						v.error(diagnostics.Incomplete_quantifier_Digit_expected, digitsStart, 0)
					} else {
						v.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, start, 1, string(ch))
						isPreviousTermQuantifiable = true
						break
					}
				} else if maxVal != "" {
					minInt := 0
					maxInt := 0
					for _, c := range minVal {
						minInt = minInt*10 + int(c-'0')
					}
					for _, c := range maxVal {
						maxInt = maxInt*10 + int(c-'0')
					}
					if minInt > maxInt && (v.anyUnicodeModeOrNonAnnexB || v.charAtOffset(0) == '}') {
						v.error(diagnostics.Numbers_out_of_order_in_quantifier, digitsStart, v.pos-digitsStart)
					}
				}
			} else if minVal == "" {
				if v.anyUnicodeModeOrNonAnnexB {
					v.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, start, 1, string(ch))
				}
				isPreviousTermQuantifiable = true
				break
			}
			if v.charAtOffset(0) != '}' {
				if v.anyUnicodeModeOrNonAnnexB {
					v.error(diagnostics.X_0_expected, v.pos, 0, "}")
					v.pos--
				} else {
					isPreviousTermQuantifiable = true
					break
				}
			}
			fallthrough
		case '*', '+', '?':
			v.pos++
			if v.charAtOffset(0) == '?' {
				v.pos++
			}
			if !isPreviousTermQuantifiable {
				v.error(diagnostics.There_is_nothing_available_for_repetition, start, v.pos-start)
			}
			isPreviousTermQuantifiable = false
		case '.':
			v.pos++
			isPreviousTermQuantifiable = true
		case '[':
			v.pos++
			if v.unicodeSetsMode {
				v.scanClassSetExpression()
			} else {
				v.scanClassRanges()
			}
			v.scanExpectedChar(']')
			isPreviousTermQuantifiable = true
		case ')':
			if isInGroup {
				return
			}
			fallthrough
		case ']', '}':
			if v.anyUnicodeModeOrNonAnnexB || ch == ')' {
				v.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, v.pos, 1, string(ch))
			}
			v.pos++
			isPreviousTermQuantifiable = true
		case '/', '|':
			return
		default:
			v.scanSourceCharacter()
			isPreviousTermQuantifiable = true
		}
	}
}

func (v *regExpValidator) validateGroupReferences() {
	for _, ref := range v.groupNameReferences {
		if !v.groupSpecifiers[ref.name] {
			v.error(diagnostics.There_is_no_capturing_group_named_0_in_this_regular_expression, ref.pos, ref.end-ref.pos, ref.name)
			// Provide spelling suggestions
			if len(v.groupSpecifiers) > 0 {
				// Convert map keys to slice
				candidates := make([]string, 0, len(v.groupSpecifiers))
				for name := range v.groupSpecifiers {
					candidates = append(candidates, name)
				}
				suggestion := core.GetSpellingSuggestion(ref.name, candidates, core.Identity[string])
				if suggestion != "" {
					v.error(diagnostics.Did_you_mean_0, ref.pos, ref.end-ref.pos, suggestion)
				}
			}
		}
	}
}

func (v *regExpValidator) validateDecimalEscapes() {
	for _, escape := range v.decimalEscapes {
		if escape.value > v.numberOfCapturingGroups {
			if v.numberOfCapturingGroups > 0 {
				v.error(diagnostics.This_backreference_refers_to_a_group_that_does_not_exist_There_are_only_0_capturing_groups_in_this_regular_expression, escape.pos, escape.end-escape.pos, v.numberOfCapturingGroups)
			} else {
				v.error(diagnostics.This_backreference_refers_to_a_group_that_does_not_exist_There_are_no_capturing_groups_in_this_regular_expression, escape.pos, escape.end-escape.pos)
			}
		}
	}
}

func (v *regExpValidator) scanDigits() {
	start := v.pos
	for v.pos < v.end && stringutil.IsDigit(v.charAtOffset(0)) {
		v.pos++
	}
	v.tokenValue = v.text[start:v.pos]
}

func (v *regExpValidator) scanExpectedChar(expected rune) {
	if v.charAtOffset(0) == expected {
		v.pos++
	} else {
		v.error(diagnostics.X_0_expected, v.pos, 0, string(expected))
	}
}

// scanFlags scans regexp flags and validates them.
// If checkModifiers is true, only allows modifier flags (i, m, s).
func (v *regExpValidator) scanFlags(currFlags regExpFlags, checkModifiers bool) regExpFlags {
	for {
		ch, size := v.charAndSize()
		if ch == 0 || !scanner.IsIdentifierPart(ch) {
			break
		}
		flag, ok := charCodeToRegExpFlag[ch]
		if !ok {
			v.error(diagnostics.Unknown_regular_expression_flag, v.pos, size)
		} else if currFlags&flag != 0 {
			v.error(diagnostics.Duplicate_regular_expression_flag, v.pos, size)
		} else if (currFlags|flag)&regExpFlagsAnyUnicodeMode == regExpFlagsAnyUnicodeMode {
			v.error(diagnostics.The_Unicode_u_flag_and_the_Unicode_Sets_v_flag_cannot_be_set_simultaneously, v.pos, size)
		} else if checkModifiers && flag&regExpFlagsModifiers == 0 {
			v.error(diagnostics.This_regular_expression_flag_cannot_be_toggled_within_a_subpattern, v.pos, size)
		} else {
			currFlags |= flag
			v.checkRegularExpressionFlagAvailability(flag, size)
		}
		v.pos += size
	}
	return currFlags
}

func (v *regExpValidator) scanPatternModifiers(currFlags regExpFlags) regExpFlags {
	return v.scanFlags(currFlags, true)
}

func (v *regExpValidator) scanAtomEscape() {
	switch v.charAtOffset(0) {
	case 'k':
		v.pos++
		if v.charAtOffset(0) == '<' {
			v.pos++
			v.scanGroupName(true)
			v.scanExpectedChar('>')
		} else if v.anyUnicodeModeOrNonAnnexB || v.namedCaptureGroups {
			v.error(diagnostics.X_k_must_be_followed_by_a_capturing_group_name_enclosed_in_angle_brackets, v.pos-2, 2)
		}
	case 'q':
		if v.unicodeSetsMode {
			v.pos++
			v.error(diagnostics.X_q_is_only_available_inside_character_class, v.pos-2, 2)
			break
		}
		fallthrough
	default:
		if !v.scanCharacterClassEscape() && !v.scanDecimalEscape() {
			v.scanCharacterEscape(true)
		}
	}
}

func (v *regExpValidator) scanDecimalEscape() bool {
	ch := v.charAtOffset(0)
	if ch >= '1' && ch <= '9' {
		start := v.pos
		v.scanDigits()
		value := 0
		for _, c := range v.tokenValue {
			value = value*10 + int(c-'0')
		}
		v.decimalEscapes = append(v.decimalEscapes, decimalEscape{pos: start, end: v.pos, value: value})
		return true
	}
	return false
}

func (v *regExpValidator) scanCharacterClassEscape() bool {
	ch := v.charAtOffset(0)
	isCharacterComplement := false
	switch ch {
	case 'd', 'D', 's', 'S', 'w', 'W':
		v.pos++
		return true
	case 'P':
		isCharacterComplement = true
		fallthrough
	case 'p':
		v.pos++
		if v.charAtOffset(0) == '{' {
			v.pos++
			v.scanUnicodePropertyValueExpression(isCharacterComplement)
		} else {
			if v.anyUnicodeModeOrNonAnnexB {
				v.error(diagnostics.X_0_must_be_followed_by_a_Unicode_property_value_expression_enclosed_in_braces, v.pos-2, 2, string(ch))
			} else {
				v.pos--
			}
		}
		return true
	}
	return false
}

func (v *regExpValidator) scanUnicodePropertyValueExpression(isCharacterComplement bool) {
	// start is at the first character after '{', so start-3 points to '\' before 'p' or 'P'
	start := v.pos - 3

	propertyNameOrValueStart := v.pos
	v.scanIdentifier(v.charAtOffset(0))
	propertyNameOrValue := v.tokenValue

	if v.charAtOffset(0) == '=' {
		// property=value syntax
		propertyNameValid := true
		if v.pos == propertyNameOrValueStart {
			v.error(diagnostics.Expected_a_Unicode_property_name, propertyNameOrValueStart, 0)
			propertyNameValid = false
		} else if !isValidNonBinaryUnicodePropertyName(propertyNameOrValue) {
			v.error(diagnostics.Unknown_Unicode_property_name, propertyNameOrValueStart, v.pos-propertyNameOrValueStart)
			// Provide spelling suggestion
			candidates := make([]string, 0, len(nonBinaryUnicodePropertyNames))
			for key := range nonBinaryUnicodePropertyNames {
				candidates = append(candidates, key)
			}
			suggestion := core.GetSpellingSuggestion(propertyNameOrValue, candidates, core.Identity[string])
			if suggestion != "" {
				v.error(diagnostics.Did_you_mean_0, propertyNameOrValueStart, v.pos-propertyNameOrValueStart, suggestion)
			}
			propertyNameValid = false
		}
		v.pos++
		propertyValueStart := v.pos
		v.scanIdentifier(v.charAtOffset(0))
		propertyValue := v.tokenValue
		if v.pos == propertyValueStart {
			v.error(diagnostics.Expected_a_Unicode_property_value, propertyValueStart, 0)
		} else if propertyNameValid && !isValidUnicodeProperty(propertyNameOrValue, propertyValue) {
			v.error(diagnostics.Unknown_Unicode_property_value, propertyValueStart, v.pos-propertyValueStart)
			// Provide spelling suggestion based on the property name
			canonicalName := nonBinaryUnicodePropertyNames[propertyNameOrValue]
			var candidates []string
			if canonicalName == "General_Category" {
				candidates = generalCategoryValues.KeysSlice()
			} else if canonicalName == "Script" || canonicalName == "Script_Extensions" {
				candidates = scriptValues.KeysSlice()
			}
			if len(candidates) > 0 {
				suggestion := core.GetSpellingSuggestion(propertyValue, candidates, core.Identity[string])
				if suggestion != "" {
					v.error(diagnostics.Did_you_mean_0, propertyValueStart, v.pos-propertyValueStart, suggestion)
				}
			}
		}
	} else {
		// property name alone
		if v.pos == propertyNameOrValueStart {
			v.error(diagnostics.Expected_a_Unicode_property_name_or_value, propertyNameOrValueStart, 0)
		} else if binaryUnicodePropertiesOfStrings.Has(propertyNameOrValue) {
			// Properties that match more than one character (strings)
			if !v.unicodeSetsMode {
				v.error(diagnostics.Any_Unicode_property_that_would_possibly_match_more_than_a_single_character_is_only_available_when_the_Unicode_Sets_v_flag_is_set, propertyNameOrValueStart, v.pos-propertyNameOrValueStart)
			} else if isCharacterComplement {
				v.error(diagnostics.Anything_that_would_possibly_match_more_than_a_single_character_is_invalid_inside_a_negated_character_class, propertyNameOrValueStart, v.pos-propertyNameOrValueStart)
			} else {
				v.mayContainStrings = true
			}
		} else if !isValidUnicodePropertyName(propertyNameOrValue) {
			v.error(diagnostics.Unknown_Unicode_property_name_or_value, propertyNameOrValueStart, v.pos-propertyNameOrValueStart)
			// Provide spelling suggestion from general category values, binary properties, and binary properties of strings
			candidates := make([]string, 0, generalCategoryValues.Len()+binaryUnicodeProperties.Len()+binaryUnicodePropertiesOfStrings.Len())
			candidates = slices.AppendSeq(candidates, maps.Keys(generalCategoryValues.M))
			candidates = slices.AppendSeq(candidates, maps.Keys(binaryUnicodeProperties.M))
			candidates = slices.AppendSeq(candidates, maps.Keys(binaryUnicodePropertiesOfStrings.M))
			suggestion := core.GetSpellingSuggestion(propertyNameOrValue, candidates, core.Identity)
			if suggestion != "" {
				v.error(diagnostics.Did_you_mean_0, propertyNameOrValueStart, v.pos-propertyNameOrValueStart, suggestion)
			}
		}
	}

	// Scan the expected closing brace
	v.scanExpectedChar('}')

	// Report the "only available when unicode mode" error AFTER validation
	if !v.anyUnicodeMode {
		v.error(diagnostics.Unicode_property_value_expressions_are_only_available_when_the_Unicode_u_flag_or_the_Unicode_Sets_v_flag_is_set, start, v.pos-start)
	}
}

func (v *regExpValidator) scanIdentifier(ch rune) {
	start := v.pos
	if ch != 0 && (scanner.IsIdentifierStart(ch) || ch == '_' || ch == '$') {
		v.pos++
		for v.pos < v.end {
			ch = v.charAtOffset(0)
			if scanner.IsIdentifierPart(ch) || ch == '_' || ch == '$' {
				v.pos++
			} else {
				break
			}
		}
	}
	v.tokenValue = v.text[start:v.pos]
}

func (v *regExpValidator) scanCharacterEscape(atomEscape bool) string {
	ch := v.charAtOffset(0)
	switch ch {
	case 0:
		v.error(diagnostics.Undetermined_character_escape, v.pos-1, 1)
		return "\\"
	case 'c':
		v.pos++
		ch = v.charAtOffset(0)
		if stringutil.IsASCIILetter(ch) {
			v.pos++
			return string(ch & 0x1f)
		}
		if v.anyUnicodeModeOrNonAnnexB {
			v.error(diagnostics.X_c_must_be_followed_by_an_ASCII_letter, v.pos-2, 2)
		} else if atomEscape {
			v.pos--
			return "\\"
		}
		return string(ch)
	case '^', '$', '/', '\\', '.', '*', '+', '?', '(', ')', '[', ']', '{', '}', '|':
		v.pos++
		return string(ch)
	default:
		return v.scanEscapeSequence(atomEscape)
	}
}

func (v *regExpValidator) scanEscapeSequence(atomEscape bool) string {
	// start points to the backslash (before the escape character)
	start := v.pos - 1
	ch := v.charAtOffset(0)
	if ch == 0 {
		v.error(diagnostics.Unexpected_end_of_text, start, 1)
		return ""
	}
	v.pos++

	switch ch {
	case '0':
		// '\0' - null character, but check if followed by digit
		if v.pos >= v.end || !stringutil.IsDigit(v.charAtOffset(0)) {
			return "\x00"
		}
		// This is an octal escape starting with \0
		// falls through to handle as octal
		if stringutil.IsOctalDigit(v.charAtOffset(0)) {
			v.pos++
		}
		fallthrough

	case '1', '2', '3':
		// Can be up to 3 octal digits
		if v.pos < v.end && stringutil.IsOctalDigit(v.charAtOffset(0)) {
			v.pos++
		}
		fallthrough

	case '4', '5', '6', '7':
		// Can be 1 or 2 octal digits (already consumed one above for 1-3)
		if v.pos < v.end && stringutil.IsOctalDigit(v.charAtOffset(0)) {
			v.pos++
		}
		// Always report errors for octal escapes in regexp mode
		code := 0
		for i := start + 1; i < v.pos; i++ {
			code = code*8 + int(v.text[i]-'0')
		}
		hexCode := fmt.Sprintf("\\x%02x", code)
		if !atomEscape && ch != '0' {
			v.error(diagnostics.Octal_escape_sequences_and_backreferences_are_not_allowed_in_a_character_class_If_this_was_intended_as_an_escape_sequence_use_the_syntax_0_instead, start, v.pos-start, hexCode)
		} else {
			v.error(diagnostics.Octal_escape_sequences_are_not_allowed_Use_the_syntax_0, start, v.pos-start, hexCode)
		}
		return string(ch)

	case '8', '9':
		// Invalid decimal escapes - always report in regexp mode
		if !atomEscape {
			v.error(diagnostics.Decimal_escape_sequences_and_backreferences_are_not_allowed_in_a_character_class, start, v.pos-start)
		} else {
			v.error(diagnostics.Escape_sequence_0_is_not_allowed, start, v.pos-start, v.text[start:v.pos])
		}
		return string(ch)

	case 'b':
		return "\b"
	case 't':
		return "\t"
	case 'n':
		return "\n"
	case 'v':
		return "\v"
	case 'f':
		return "\f"
	case 'r':
		return "\r"

	case 'x':
		// Hex escape '\xDD'
		hexStart := v.pos
		validHex := true
		for range 2 {
			if v.pos >= v.end || !stringutil.IsHexDigit(v.charAtOffset(0)) {
				validHex = false
				break
			}
			v.pos++
		}
		if !validHex {
			v.error(diagnostics.Hexadecimal_digit_expected, hexStart, v.pos-hexStart)
			return v.text[start:v.pos]
		}
		code := parseHexValue(v.text, hexStart, v.pos)
		return string(rune(code))

	case 'u':
		// Unicode escape '\uDDDD' or '\u{DDDDDD}'
		if v.charAtOffset(0) == '{' {
			// Extended unicode escape \u{DDDDDD}
			v.pos++
			hexStart := v.pos
			hasDigits := false
			for v.pos < v.end && stringutil.IsHexDigit(v.charAtOffset(0)) {
				hasDigits = true
				v.pos++
			}
			if !hasDigits {
				v.error(diagnostics.Hexadecimal_digit_expected, hexStart, 0)
				return v.text[start:v.pos]
			}
			if v.charAtOffset(0) == '}' {
				v.pos++
			} else if hasDigits {
				v.error(diagnostics.Unterminated_Unicode_escape_sequence, start, v.pos-start)
				return v.text[start:v.pos]
			}
			// Parse hex value (-1 to skip closing brace)
			code := parseHexValue(v.text, hexStart, v.pos-1)
			// Validate the code point is within valid Unicode range
			if code > 0x10FFFF {
				v.error(diagnostics.An_extended_Unicode_escape_value_must_be_between_0x0_and_0x10FFFF_inclusive, hexStart, v.pos-1-hexStart)
			}
			if !v.anyUnicodeMode {
				v.error(diagnostics.Unicode_escape_sequences_are_only_available_when_the_Unicode_u_flag_or_the_Unicode_Sets_v_flag_is_set, start, v.pos-start)
			}
			return string(rune(code))
		} else {
			// Standard unicode escape '\uDDDD'
			hexStart := v.pos
			validHex := true
			for range 4 {
				if v.pos >= v.end || !stringutil.IsHexDigit(v.charAtOffset(0)) {
					validHex = false
					break
				}
				v.pos++
			}
			if !validHex {
				v.error(diagnostics.Hexadecimal_digit_expected, hexStart, v.pos-hexStart)
				return v.text[start:v.pos]
			}
			code := parseHexValue(v.text, hexStart, v.pos)
			// For surrogates, we need to preserve the actual value since string(rune(surrogate))
			// converts to 0xFFFD. We encode the surrogate as UTF-16BE bytes.
			var escapedValueString string
			if isSurrogate(rune(code)) {
				// Surrogate - encode as 2-byte sequence (UTF-16BE)
				escapedValueString = encodeSurrogate(rune(code))
			} else {
				escapedValueString = string(rune(code))
			}
			// In Unicode mode, check for surrogate pairs
			if v.anyUnicodeMode && isHighSurrogate(rune(code)) &&
				v.pos+6 <= v.end && v.text[v.pos:v.pos+2] == "\\u" {
				// High surrogate followed by potential low surrogate
				nextStart := v.pos
				nextPos := v.pos + 2
				validNext := true
				for range 4 {
					if nextPos >= v.end || !stringutil.IsHexDigit(rune(v.text[nextPos])) {
						validNext = false
						break
					}
					nextPos++
				}
				if validNext {
					// Parse the next escape
					nextCode := parseHexValue(v.text, nextStart+2, nextPos)
					// Check if it's a low surrogate
					if isLowSurrogate(rune(nextCode)) {
						// Combine surrogates into a single code point
						combinedCodePoint := combineSurrogatePair(rune(code), rune(nextCode))
						v.pos = nextPos
						return string(combinedCodePoint)
					}
				}
			}
			return escapedValueString
		}

	default:
		// Identity escape or invalid escape
		// Report error if:
		// - In any Unicode mode, OR
		// - In regexp mode, not Annex B, and ch is an identifier part
		if v.anyUnicodeMode || (v.anyUnicodeModeOrNonAnnexB && scanner.IsIdentifierPart(ch)) {
			v.error(diagnostics.This_character_cannot_be_escaped_in_a_regular_expression, start, v.pos-start)
		}
		return string(ch)
	}
}

// parseHexValue parses hexadecimal digits from text and returns the integer value
func parseHexValue(text string, start, end int) int {
	code := 0
	for i := start; i < end; i++ {
		digit := text[i]
		if digit >= '0' && digit <= '9' {
			code = code*16 + int(digit-'0')
		} else if digit >= 'a' && digit <= 'f' {
			code = code*16 + int(digit-'a'+10)
		} else if digit >= 'A' && digit <= 'F' {
			code = code*16 + int(digit-'A'+10)
		}
	}
	return code
}

func (v *regExpValidator) scanGroupName(isReference bool) {
	tokenStart := v.pos
	v.scanIdentifier(v.charAtOffset(0))
	if v.pos == tokenStart {
		v.error(diagnostics.Expected_a_capturing_group_name, v.pos, 0)
	} else if isReference {
		v.groupNameReferences = append(v.groupNameReferences, namedReference{pos: tokenStart, end: v.pos, name: v.tokenValue})
	} else {
		// Check for duplicate names in scope
		if v.topNamedCapturingGroupsScope != nil && v.topNamedCapturingGroupsScope[v.tokenValue] {
			v.error(diagnostics.Named_capturing_groups_with_the_same_name_must_be_mutually_exclusive_to_each_other, tokenStart, v.pos-tokenStart)
		} else {
			for _, scope := range v.namedCapturingGroupsScopeStack {
				if scope != nil && scope[v.tokenValue] {
					v.error(diagnostics.Named_capturing_groups_with_the_same_name_must_be_mutually_exclusive_to_each_other, tokenStart, v.pos-tokenStart)
					break
				}
			}
		}
		if v.topNamedCapturingGroupsScope == nil {
			v.topNamedCapturingGroupsScope = make(map[string]bool)
		}
		v.topNamedCapturingGroupsScope[v.tokenValue] = true
		if v.groupSpecifiers == nil {
			v.groupSpecifiers = make(map[string]bool)
		}
		v.groupSpecifiers[v.tokenValue] = true
	}
}

// scanSourceCharacter scans and returns a single "character" from the source.
// In Unicode mode (u or v flags), returns complete Unicode code points.
// In non-Unicode mode, mimics JavaScript's UTF-16 behavior where literal characters
// >= U+10000 are treated as surrogate pairs and consumed across two sequential calls.
func (v *regExpValidator) scanSourceCharacter() string {
	// Check if we have a pending low surrogate from the previous call
	if v.surrogateState != nil {
		low := v.surrogateState.lowSurrogate
		size := v.surrogateState.utf8Size
		v.surrogateState = nil
		v.pos += size
		// Return the low surrogate encoded as UTF-16BE
		return encodeSurrogate(low)
	}

	// Decode the next UTF-8 character from the source
	r, s := v.charAndSize()

	if v.anyUnicodeMode || r < supplementaryMin {
		// In Unicode mode, or for BMP characters, consume and return normally
		v.pos += s
		return v.text[v.pos-s : v.pos]
	}

	// In non-Unicode mode with a supplementary character (>= U+10000):
	// JavaScript represents these as surrogate pairs (two UTF-16 code units).
	// Return the high surrogate now and save the low surrogate for the next call.
	high, low := splitToSurrogatePair(r)

	v.surrogateState = &surrogatePairState{
		lowSurrogate: low,
		utf8Size:     s,
	}

	return encodeSurrogate(high)
}

// ClassRanges ::= ClassAtom ('-' ClassAtom)?
// Scans character class content like [a-z] or [^0-9].
// Follows ECMAScript regexp grammar
func (v *regExpValidator) scanClassRanges() {
	isNegated := v.charAtOffset(0) == '^'
	if isNegated {
		v.pos++
	}
	oldIsCharacterComplement := v.isCharacterComplement
	v.isCharacterComplement = isNegated
	defer func() {
		v.isCharacterComplement = oldIsCharacterComplement
	}()
	for {
		ch := v.charAtOffset(0)
		if v.isClassContentExit(ch) {
			return
		}
		atomStart := v.pos
		atom := v.scanClassAtom()
		if v.charAtOffset(0) == '-' && v.charAtOffset(1) != ']' {
			v.pos++
			if v.isClassContentExit(v.charAtOffset(0)) {
				return
			}
			// Check if min side of range is a character class escape
			if atom == "" && v.anyUnicodeModeOrNonAnnexB {
				v.error(diagnostics.A_character_class_range_must_not_be_bounded_by_another_character_class, atomStart, v.pos-1-atomStart)
			}
			rangeEndStart := v.pos
			rangeEnd := v.scanClassAtom()
			// Check if max side of range is a character class escape
			if rangeEnd == "" && v.anyUnicodeModeOrNonAnnexB {
				v.error(diagnostics.A_character_class_range_must_not_be_bounded_by_another_character_class, rangeEndStart, v.pos-rangeEndStart)
			}
			// Check range order
			if atom != "" && rangeEnd != "" {
				minCodePoint := decodeCodePoint(atom)
				maxCodePoint := decodeCodePoint(rangeEnd)

				// Get the expected sizes (in UTF-16 code units)
				minExpectedSize := charSize(minCodePoint)
				maxExpectedSize := charSize(maxCodePoint)

				// Check if both are "complete" characters in JavaScript's UTF-16 representation.
				// A character is complete if its UTF-16 length matches the expected size.
				// In JavaScript, string.length returns the UTF-16 code unit count.
				minUTF16Length := utf16Length(atom)
				maxUTF16Length := utf16Length(rangeEnd)

				minIsComplete := minUTF16Length == minExpectedSize
				maxIsComplete := maxUTF16Length == maxExpectedSize

				if minIsComplete && maxIsComplete && minCodePoint > maxCodePoint {
					// In non-Unicode mode, literal characters >= 0x10000 are scanned
					// as individual surrogates by scanSourceCharacter(), so the code
					// points being compared are already the surrogate values (0xD800-0xDFFF).
					// Escape sequences like \u{1D608} return the full character, so the
					// code points are the actual values (>= 0x10000).

					// If there's a pending low surrogate from scanning the second atom,
					// we need to account for its UTF-8 size in the error range.
					errorEnd := v.pos
					if v.surrogateState != nil {
						errorEnd += v.surrogateState.utf8Size
					}
					length := errorEnd - atomStart
					v.error(diagnostics.Range_out_of_order_in_character_class, atomStart, length)
				}
			}
		}
	}
}

func (v *regExpValidator) isClassContentExit(ch rune) bool {
	return ch == ']' || ch == 0 || v.pos >= v.end
}

// ClassAtom ::=
//
//	| SourceCharacter but not one of '\' or ']'
//	| '\' ClassEscape
//
// ClassEscape ::=
//
//	| 'b'
//	| '-'
//	| CharacterClassEscape
//	| CharacterEscape
func (v *regExpValidator) scanClassAtom() string {
	if v.charAtOffset(0) == '\\' {
		v.pos++
		ch := v.charAtOffset(0)
		switch ch {
		case 'b':
			v.pos++
			return "\b" // backspace character
		case '-':
			v.pos++
			return string(ch) // hyphen character
		default:
			if v.scanCharacterClassEscape() {
				return ""
			}
			return v.scanCharacterEscape(false)
		}
	}
	return v.scanSourceCharacter()
}

type classSetExpressionType int

const (
	classSetExpressionNone classSetExpressionType = iota
	classSetExpressionSubtraction
	classSetExpressionIntersection
)

func (v *regExpValidator) scanClassSetExpression() {
	isCharacterComplement := false
	if v.charAtOffset(0) == '^' {
		v.pos++
		isCharacterComplement = true
	}

	oldIsCharacterComplement := v.isCharacterComplement
	v.isCharacterComplement = isCharacterComplement
	defer func() {
		v.isCharacterComplement = oldIsCharacterComplement
	}()

	expressionMayContainStrings := false
	ch := v.charAtOffset(0)
	if v.isClassContentExit(ch) {
		return
	}

	start := v.pos
	var operand string

	// Check for operators at the start
	slice := v.text[v.pos:min(v.pos+2, v.end)]
	if slice == "--" || slice == "&&" {
		v.error(diagnostics.Expected_a_class_set_operand, v.pos, 0)
		v.mayContainStrings = false
	} else {
		operand = v.scanClassSetOperand()
	}

	// Check what follows the first operand
	switch v.charAtOffset(0) {
	case '-':
		if v.charAtOffset(1) == '-' {
			if isCharacterComplement && v.mayContainStrings {
				v.error(diagnostics.Anything_that_would_possibly_match_more_than_a_single_character_is_invalid_inside_a_negated_character_class, start, v.pos-start)
			}
			expressionMayContainStrings = v.mayContainStrings
			v.scanClassSetSubExpression(classSetExpressionSubtraction)
			v.mayContainStrings = !isCharacterComplement && expressionMayContainStrings
			return
		}
	case '&':
		if v.charAtOffset(1) == '&' {
			v.scanClassSetSubExpression(classSetExpressionIntersection)
			if isCharacterComplement && v.mayContainStrings {
				v.error(diagnostics.Anything_that_would_possibly_match_more_than_a_single_character_is_invalid_inside_a_negated_character_class, start, v.pos-start)
			}
			expressionMayContainStrings = v.mayContainStrings
			v.mayContainStrings = !isCharacterComplement && expressionMayContainStrings
			return
		} else {
			v.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, v.pos, 1, string(ch))
		}
	default:
		if isCharacterComplement && v.mayContainStrings {
			v.error(diagnostics.Anything_that_would_possibly_match_more_than_a_single_character_is_invalid_inside_a_negated_character_class, start, v.pos-start)
		}
		expressionMayContainStrings = v.mayContainStrings
	}

	// Continue scanning operands
	for {
		ch = v.charAtOffset(0)
		if ch == 0 {
			break
		}

		switch ch {
		case '-':
			v.pos++
			ch = v.charAtOffset(0)
			if v.isClassContentExit(ch) {
				v.mayContainStrings = !isCharacterComplement && expressionMayContainStrings
				return
			}
			if ch == '-' {
				v.pos++
				v.error(diagnostics.Operators_must_not_be_mixed_within_a_character_class_Wrap_it_in_a_nested_class_instead, v.pos-2, 2)
				start = v.pos - 2
				operand = v.text[start:v.pos]
				continue
			} else {
				if operand == "" {
					v.error(diagnostics.A_character_class_range_must_not_be_bounded_by_another_character_class, start, v.pos-1-start)
				}
				secondStart := v.pos
				secondOperand := v.scanClassSetOperand()
				if isCharacterComplement && v.mayContainStrings {
					v.error(diagnostics.Anything_that_would_possibly_match_more_than_a_single_character_is_invalid_inside_a_negated_character_class, secondStart, v.pos-secondStart)
				}
				expressionMayContainStrings = expressionMayContainStrings || v.mayContainStrings
				if secondOperand == "" {
					v.error(diagnostics.A_character_class_range_must_not_be_bounded_by_another_character_class, secondStart, v.pos-secondStart)
					break
				}
				if operand == "" {
					break
				}
				// Check range order
				minRune, minSize := utf8.DecodeRuneInString(operand)
				maxRune, maxSize := utf8.DecodeRuneInString(secondOperand)
				if len(operand) == minSize && len(secondOperand) == maxSize && minRune > maxRune {
					v.error(diagnostics.Range_out_of_order_in_character_class, start, v.pos-start)
				}
			}

		case '&':
			start = v.pos
			v.pos++
			if v.charAtOffset(0) == '&' {
				v.pos++
				v.error(diagnostics.Operators_must_not_be_mixed_within_a_character_class_Wrap_it_in_a_nested_class_instead, v.pos-2, 2)
				if v.charAtOffset(0) == '&' {
					v.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, v.pos, 1, string(ch))
					v.pos++
				}
			} else {
				v.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, v.pos-1, 1, string(ch))
			}
			operand = v.text[start:v.pos]
			continue
		}

		if v.isClassContentExit(v.charAtOffset(0)) {
			break
		}

		start = v.pos
		slice = v.text[v.pos:min(v.pos+2, v.end)]
		if slice == "--" || slice == "&&" {
			v.error(diagnostics.Operators_must_not_be_mixed_within_a_character_class_Wrap_it_in_a_nested_class_instead, v.pos, 2)
			v.pos += 2
			operand = v.text[start:v.pos]
		} else {
			operand = v.scanClassSetOperand()
		}
	}
	v.mayContainStrings = !isCharacterComplement && expressionMayContainStrings
}

func (v *regExpValidator) scanClassSetSubExpression(expressionType classSetExpressionType) {
	expressionMayContainStrings := v.mayContainStrings
	for {
		ch := v.charAtOffset(0)
		if v.isClassContentExit(ch) {
			break
		}

		// Provide user-friendly diagnostic messages
		switch ch {
		case '-':
			v.pos++
			if v.charAtOffset(0) == '-' {
				v.pos++
				if expressionType != classSetExpressionSubtraction {
					v.error(diagnostics.Operators_must_not_be_mixed_within_a_character_class_Wrap_it_in_a_nested_class_instead, v.pos-2, 2)
				}
			} else {
				v.error(diagnostics.Operators_must_not_be_mixed_within_a_character_class_Wrap_it_in_a_nested_class_instead, v.pos-1, 1)
			}
		case '&':
			v.pos++
			if v.charAtOffset(0) == '&' {
				v.pos++
				if expressionType != classSetExpressionIntersection {
					v.error(diagnostics.Operators_must_not_be_mixed_within_a_character_class_Wrap_it_in_a_nested_class_instead, v.pos-2, 2)
				}
				if v.charAtOffset(0) == '&' {
					v.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, v.pos, 1, string(ch))
					v.pos++
				}
			} else {
				v.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, v.pos-1, 1, string(ch))
			}
		default:
			switch expressionType {
			case classSetExpressionSubtraction:
				v.error(diagnostics.X_0_expected, v.pos, 0, "--")
			case classSetExpressionIntersection:
				v.error(diagnostics.X_0_expected, v.pos, 0, "&&")
			}
		}

		ch = v.charAtOffset(0)
		if v.isClassContentExit(ch) {
			v.error(diagnostics.Expected_a_class_set_operand, v.pos, 0)
			break
		}
		v.scanClassSetOperand()
		// Used only if expressionType is Intersection
		expressionMayContainStrings = expressionMayContainStrings && v.mayContainStrings
	}
	v.mayContainStrings = expressionMayContainStrings
}

// ClassSetOperand ::=
//
//	| '[' ClassSetExpression ']'
//	| '\' CharacterClassEscape
//	| '\q{' ClassStringDisjunctionContents '}'
//	| ClassSetCharacter
func (v *regExpValidator) scanClassSetOperand() string {
	v.mayContainStrings = false
	switch v.charAtOffset(0) {
	case 0:
		return ""
	case '[':
		v.pos++
		v.scanClassSetExpression()
		v.scanExpectedChar(']')
		return ""
	case '\\':
		v.pos++
		if v.scanCharacterClassEscape() {
			return ""
		} else if v.charAtOffset(0) == 'q' {
			v.pos++
			if v.charAtOffset(0) == '{' {
				v.pos++
				v.scanClassStringDisjunctionContents()
				v.scanExpectedChar('}')
				return ""
			} else {
				v.error(diagnostics.X_q_must_be_followed_by_string_alternatives_enclosed_in_braces, v.pos-2, 2)
				return "q"
			}
		}
		v.pos--
		// falls through
	}
	return v.scanClassSetCharacter()
}

// ClassStringDisjunctionContents ::= ClassSetCharacter* ('|' ClassSetCharacter*)*
func (v *regExpValidator) scanClassStringDisjunctionContents() {
	characterCount := 0
	for {
		ch := v.charAtOffset(0)
		switch ch {
		case 0:
			return
		case '}':
			if characterCount != 1 {
				v.mayContainStrings = true
			}
			return
		case '|':
			if characterCount != 1 {
				v.mayContainStrings = true
			}
			v.pos++
			characterCount = 0
		default:
			v.scanClassSetCharacter()
			characterCount++
		}
	}
}

// ClassSetCharacter ::=
//
//	| SourceCharacter -- ClassSetSyntaxCharacter -- ClassSetReservedDoublePunctuator
//	| '\' (CharacterEscape | ClassSetReservedPunctuator | 'b')
func (v *regExpValidator) scanClassSetCharacter() string {
	ch := v.charAtOffset(0)
	if ch == 0 {
		return ""
	}

	if ch == '\\' {
		v.pos++
		ch = v.charAtOffset(0)
		switch ch {
		case 'b':
			v.pos++
			return "\b"
		case '&', '-', '!', '#', '%', ',', ':', ';', '<', '=', '>', '@', '`', '~':
			v.pos++
			return string(ch)
		default:
			return v.scanCharacterEscape(false)
		}
	} else if ch == v.charAtOffset(1) {
		// Check for reserved double punctuators
		switch ch {
		case '&', '!', '#', '%', '*', '+', ',', '.', ':', ';', '<', '=', '>', '?', '@', '`', '~':
			v.error(diagnostics.A_character_class_must_not_contain_a_reserved_double_punctuator_Did_you_mean_to_escape_it_with_backslash, v.pos, 2)
			v.pos += 2
			return v.text[v.pos-2 : v.pos]
		}
	}

	switch ch {
	case '/', '(', ')', '[', ']', '{', '}', '-', '|':
		v.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, v.pos, 1, string(ch))
		v.pos++
		return string(ch)
	}

	return v.scanSourceCharacter()
}
