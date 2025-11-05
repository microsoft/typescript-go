package scanner

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/stringutil"
)

// RegExpFlags represents regexp flags (e.g., 'g', 'i', 'm', etc.)
type RegExpFlags uint32

const (
	RegExpFlagsNone           RegExpFlags = 0
	RegExpFlagsGlobal         RegExpFlags = 1 << 0 // g
	RegExpFlagsIgnoreCase     RegExpFlags = 1 << 1 // i
	RegExpFlagsMultiline      RegExpFlags = 1 << 2 // m
	RegExpFlagsDotAll         RegExpFlags = 1 << 3 // s
	RegExpFlagsUnicode        RegExpFlags = 1 << 4 // u
	RegExpFlagsSticky         RegExpFlags = 1 << 5 // y
	RegExpFlagsHasIndices     RegExpFlags = 1 << 6 // d
	RegExpFlagsUnicodeSets    RegExpFlags = 1 << 7 // v
	RegExpFlagsModifiers      RegExpFlags = RegExpFlagsIgnoreCase | RegExpFlagsMultiline | RegExpFlagsDotAll
	RegExpFlagsAnyUnicodeMode RegExpFlags = RegExpFlagsUnicode | RegExpFlagsUnicodeSets
)

var charCodeToRegExpFlag = map[rune]RegExpFlags{
	'd': RegExpFlagsHasIndices,
	'g': RegExpFlagsGlobal,
	'i': RegExpFlagsIgnoreCase,
	'm': RegExpFlagsMultiline,
	's': RegExpFlagsDotAll,
	'u': RegExpFlagsUnicode,
	'v': RegExpFlagsUnicodeSets,
	'y': RegExpFlagsSticky,
}

var regExpFlagToCharCode = map[RegExpFlags]rune{
	RegExpFlagsHasIndices:  'd',
	RegExpFlagsGlobal:      'g',
	RegExpFlagsIgnoreCase:  'i',
	RegExpFlagsMultiline:   'm',
	RegExpFlagsDotAll:      's',
	RegExpFlagsUnicode:     'u',
	RegExpFlagsUnicodeSets: 'v',
	RegExpFlagsSticky:      'y',
}

// CharacterCodeToRegularExpressionFlag converts a character code to a regexp flag
func CharacterCodeToRegularExpressionFlag(ch rune) (RegExpFlags, bool) {
	flag, ok := charCodeToRegExpFlag[ch]
	return flag, ok
}

// RegularExpressionFlagToCharacterCode converts a regexp flag to a character code
func RegularExpressionFlagToCharacterCode(f RegExpFlags) (rune, bool) {
	ch, ok := regExpFlagToCharCode[f]
	return ch, ok
}

// RegExpValidator is used to validate regular expressions
type RegExpValidator struct {
	text                           string
	pos                            int
	end                            int
	languageVersion                core.ScriptTarget
	languageVariant                core.LanguageVariant
	onError                        ErrorCallback
	regExpFlags                    RegExpFlags
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
	pendingLowSurrogate            rune // For non-Unicode mode: pending low surrogate to return
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

// ValidateRegularExpressionLiteral validates a regular expression literal node
// This is called from the checker package during semantic analysis
func ValidateRegularExpressionLiteral(
	node *ast.RegularExpressionLiteral,
	sourceFile *ast.SourceFile,
	languageVersion core.ScriptTarget,
	languageVariant core.LanguageVariant,
	onError ErrorCallback,
) {
	// Check if the token has the Unterminated flag
	if node.AsNode().LiteralLikeData().TokenFlags&ast.TokenFlagsUnterminated != 0 {
		return
	}

	text := node.Text
	v := &RegExpValidator{
		text:            text,
		pos:             1, // Skip initial '/'
		end:             len(text),
		languageVersion: languageVersion,
		languageVariant: languageVariant,
		onError:         onError,
	}

	v.validateRegularExpression()
}

func (v *RegExpValidator) validateRegularExpression() {
	// Find the body end (before flags)
	bodyEnd := v.findRegExpBodyEnd()
	if bodyEnd < 0 {
		return // Already unterminated
	}

	// Parse flags
	flagsStart := bodyEnd + 1
	v.pos = flagsStart
	v.regExpFlags = RegExpFlagsNone

	for v.pos < v.end {
		ch, size := v.charAndSize()
		if !IsIdentifierPart(ch) {
			break
		}

		flag, ok := CharacterCodeToRegularExpressionFlag(ch)
		if !ok {
			v.error(diagnostics.Unknown_regular_expression_flag, v.pos, size)
		} else if v.regExpFlags&flag != 0 {
			v.error(diagnostics.Duplicate_regular_expression_flag, v.pos, size)
		} else if (v.regExpFlags|flag)&RegExpFlagsAnyUnicodeMode == RegExpFlagsAnyUnicodeMode {
			v.error(diagnostics.The_Unicode_u_flag_and_the_Unicode_Sets_v_flag_cannot_be_set_simultaneously, v.pos, size)
		} else {
			v.regExpFlags |= flag
			v.checkRegularExpressionFlagAvailability(flag, size)
		}
		v.pos += size
	}

	// Set up validation parameters
	v.unicodeSetsMode = v.regExpFlags&RegExpFlagsUnicodeSets != 0
	v.anyUnicodeMode = v.regExpFlags&RegExpFlagsAnyUnicodeMode != 0
	// TypeScript always uses Annex B mode for regexp validation
	v.annexB = true
	v.anyUnicodeModeOrNonAnnexB = v.anyUnicodeMode || !v.annexB

	// Validate the pattern body
	v.pos = 1 // Reset to start of pattern (after initial '/')
	v.end = bodyEnd
	v.scanDisjunction(false)

	// Post-validation checks
	v.validateGroupReferences()
	v.validateDecimalEscapes()
}

func (v *RegExpValidator) findRegExpBodyEnd() int {
	pos := 1 // Skip initial '/'
	inEscape := false
	inCharacterClass := false

	for pos < len(v.text) {
		ch := v.text[pos]
		if ch == '\\' && !inEscape {
			inEscape = true
			pos++
			continue
		}

		if inEscape {
			inEscape = false
			pos++
			continue
		}

		if ch == '/' && !inCharacterClass {
			return pos
		}

		if ch == '[' {
			inCharacterClass = true
		} else if ch == ']' {
			inCharacterClass = false
		}

		pos++
	}

	return -1 // Unterminated
}

func (v *RegExpValidator) textStart() int {
	return 1 // Offset from node start (which includes the '/')
}

func (v *RegExpValidator) charAndSize() (rune, int) {
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

func (v *RegExpValidator) charAtOffset(offset int) rune {
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

func (v *RegExpValidator) error(message *diagnostics.Message, start, length int, args ...any) {
	if v.onError != nil {
		v.onError(message, start, length, args...)
	}
}

func (v *RegExpValidator) checkRegularExpressionFlagAvailability(flag RegExpFlags, size int) {
	var availableFrom core.ScriptTarget
	switch flag {
	case RegExpFlagsHasIndices:
		availableFrom = core.ScriptTargetES2022
	case RegExpFlagsDotAll:
		availableFrom = core.ScriptTargetES2018
	case RegExpFlagsUnicodeSets:
		availableFrom = core.ScriptTargetES2024
	default:
		return
	}

	if v.languageVersion < availableFrom {
		// Workaround: TypeScript uses lowercase ES version names in error messages (e.g., "es2024" not "ES2024")
		v.error(diagnostics.This_regular_expression_flag_is_only_available_when_targeting_0_or_later, v.pos, size, strings.ToLower(availableFrom.String()))
	}
}

func (v *RegExpValidator) scanDisjunction(isInGroup bool) {
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

func (v *RegExpValidator) scanAlternative(isInGroup bool) {
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
					setFlags := v.scanPatternModifiers(RegExpFlagsNone)
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

func (v *RegExpValidator) validateGroupReferences() {
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

func (v *RegExpValidator) validateDecimalEscapes() {
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

func (v *RegExpValidator) scanDigits() {
	start := v.pos
	for v.pos < v.end && v.charAtOffset(0) >= '0' && v.charAtOffset(0) <= '9' {
		v.pos++
	}
	v.tokenValue = v.text[start:v.pos]
}

func (v *RegExpValidator) scanExpectedChar(expected rune) {
	if v.charAtOffset(0) == expected {
		v.pos++
	} else {
		v.error(diagnostics.X_0_expected, v.pos, 0, string(expected))
	}
}

func (v *RegExpValidator) scanPatternModifiers(currFlags RegExpFlags) RegExpFlags {
	for {
		ch, size := v.charAndSize()
		if ch == 0 || !IsIdentifierPart(ch) {
			break
		}
		flag, ok := CharacterCodeToRegularExpressionFlag(ch)
		if !ok {
			v.error(diagnostics.Unknown_regular_expression_flag, v.pos, size)
		} else if currFlags&flag != 0 {
			v.error(diagnostics.Duplicate_regular_expression_flag, v.pos, size)
		} else if flag&RegExpFlagsModifiers == 0 {
			v.error(diagnostics.This_regular_expression_flag_cannot_be_toggled_within_a_subpattern, v.pos, size)
		} else {
			currFlags |= flag
			v.checkRegularExpressionFlagAvailability(flag, size)
		}
		v.pos += size
	}
	return currFlags
}

func (v *RegExpValidator) scanAtomEscape() {
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

func (v *RegExpValidator) scanDecimalEscape() bool {
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

func (v *RegExpValidator) scanCharacterClassEscape() bool {
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

func (v *RegExpValidator) scanUnicodePropertyValueExpression(isCharacterComplement bool) {
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
		} else if !v.isValidNonBinaryUnicodePropertyName(propertyNameOrValue) {
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
		} else if propertyNameValid && !v.isValidUnicodeProperty(propertyNameOrValue, propertyValue) {
			v.error(diagnostics.Unknown_Unicode_property_value, propertyValueStart, v.pos-propertyValueStart)
			// Provide spelling suggestion based on the property name
			canonicalName := nonBinaryUnicodePropertyNames[propertyNameOrValue]
			var candidates []string
			if canonicalName == "General_Category" {
				candidates = make([]string, 0, len(generalCategoryValues))
				for key := range generalCategoryValues {
					candidates = append(candidates, key)
				}
			} else if canonicalName == "Script" || canonicalName == "Script_Extensions" {
				candidates = make([]string, 0, len(scriptValues))
				for key := range scriptValues {
					candidates = append(candidates, key)
				}
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
		} else if binaryUnicodePropertiesOfStrings[propertyNameOrValue] {
			// Properties that match more than one character (strings)
			if !v.unicodeSetsMode {
				v.error(diagnostics.Any_Unicode_property_that_would_possibly_match_more_than_a_single_character_is_only_available_when_the_Unicode_Sets_v_flag_is_set, propertyNameOrValueStart, v.pos-propertyNameOrValueStart)
			} else if isCharacterComplement {
				v.error(diagnostics.Anything_that_would_possibly_match_more_than_a_single_character_is_invalid_inside_a_negated_character_class, propertyNameOrValueStart, v.pos-propertyNameOrValueStart)
			} else {
				v.mayContainStrings = true
			}
		} else if !v.isValidUnicodePropertyName(propertyNameOrValue) {
			v.error(diagnostics.Unknown_Unicode_property_name_or_value, propertyNameOrValueStart, v.pos-propertyNameOrValueStart)
			// Provide spelling suggestion from general category values, binary properties, and binary properties of strings
			candidates := make([]string, 0, len(generalCategoryValues)+len(binaryUnicodeProperties)+len(binaryUnicodePropertiesOfStrings))
			for key := range generalCategoryValues {
				candidates = append(candidates, key)
			}
			for key := range binaryUnicodeProperties {
				candidates = append(candidates, key)
			}
			for key := range binaryUnicodePropertiesOfStrings {
				candidates = append(candidates, key)
			}
			suggestion := core.GetSpellingSuggestion(propertyNameOrValue, candidates, core.Identity[string])
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

// Table 67: Binary Unicode property aliases and their canonical property names
// https://tc39.es/ecma262/#table-binary-unicode-properties
var binaryUnicodeProperties = map[string]bool{
	"ASCII": true, "ASCII_Hex_Digit": true, "AHex": true, "Alphabetic": true, "Alpha": true, "Any": true, "Assigned": true, "Bidi_Control": true, "Bidi_C": true, "Bidi_Mirrored": true, "Bidi_M": true, "Case_Ignorable": true, "CI": true, "Cased": true, "Changes_When_Casefolded": true, "CWCF": true, "Changes_When_Casemapped": true, "CWCM": true, "Changes_When_Lowercased": true, "CWL": true, "Changes_When_NFKC_Casefolded": true, "CWKCF": true, "Changes_When_Titlecased": true, "CWT": true, "Changes_When_Uppercased": true, "CWU": true, "Dash": true, "Default_Ignorable_Code_Point": true, "DI": true, "Deprecated": true, "Dep": true, "Diacritic": true, "Dia": true, "Emoji": true, "Emoji_Component": true, "EComp": true, "Emoji_Modifier": true, "EMod": true, "Emoji_Modifier_Base": true, "EBase": true, "Emoji_Presentation": true, "EPres": true, "Extended_Pictographic": true, "ExtPict": true, "Extender": true, "Ext": true, "Grapheme_Base": true, "Gr_Base": true, "Grapheme_Extend": true, "Gr_Ext": true, "Hex_Digit": true, "Hex": true, "IDS_Binary_Operator": true, "IDSB": true, "IDS_Trinary_Operator": true, "IDST": true, "ID_Continue": true, "IDC": true, "ID_Start": true, "IDS": true, "Ideographic": true, "Ideo": true, "Join_Control": true, "Join_C": true, "Logical_Order_Exception": true, "LOE": true, "Lowercase": true, "Lower": true, "Math": true, "Noncharacter_Code_Point": true, "NChar": true, "Pattern_Syntax": true, "Pat_Syn": true, "Pattern_White_Space": true, "Pat_WS": true, "Quotation_Mark": true, "QMark": true, "Radical": true, "Regional_Indicator": true, "RI": true, "Sentence_Terminal": true, "STerm": true, "Soft_Dotted": true, "SD": true, "Terminal_Punctuation": true, "Term": true, "Unified_Ideograph": true, "UIdeo": true, "Uppercase": true, "Upper": true, "Variation_Selector": true, "VS": true, "White_Space": true, "space": true, "XID_Continue": true, "XIDC": true, "XID_Start": true, "XIDS": true,
}

// Table 68: Binary Unicode properties of strings
// https://tc39.es/ecma262/#table-binary-unicode-properties-of-strings
var binaryUnicodePropertiesOfStrings = map[string]bool{
	"Basic_Emoji": true, "Emoji_Keycap_Sequence": true, "RGI_Emoji_Modifier_Sequence": true, "RGI_Emoji_Flag_Sequence": true, "RGI_Emoji_Tag_Sequence": true, "RGI_Emoji_ZWJ_Sequence": true, "RGI_Emoji": true,
}

// Unicode 15.1 - General_Category values
var generalCategoryValues = map[string]bool{
	"C": true, "Other": true, "Cc": true, "Control": true, "cntrl": true, "Cf": true, "Format": true, "Cn": true, "Unassigned": true, "Co": true, "Private_Use": true, "Cs": true, "Surrogate": true, "L": true, "Letter": true, "LC": true, "Cased_Letter": true, "Ll": true, "Lowercase_Letter": true, "Lm": true, "Modifier_Letter": true, "Lo": true, "Other_Letter": true, "Lt": true, "Titlecase_Letter": true, "Lu": true, "Uppercase_Letter": true, "M": true, "Mark": true, "Combining_Mark": true, "Mc": true, "Spacing_Mark": true, "Me": true, "Enclosing_Mark": true, "Mn": true, "Nonspacing_Mark": true, "N": true, "Number": true, "Nd": true, "Decimal_Number": true, "digit": true, "Nl": true, "Letter_Number": true, "No": true, "Other_Number": true, "P": true, "Punctuation": true, "punct": true, "Pc": true, "Connector_Punctuation": true, "Pd": true, "Dash_Punctuation": true, "Pe": true, "Close_Punctuation": true, "Pf": true, "Final_Punctuation": true, "Pi": true, "Initial_Punctuation": true, "Po": true, "Other_Punctuation": true, "Ps": true, "Open_Punctuation": true, "S": true, "Symbol": true, "Sc": true, "Currency_Symbol": true, "Sk": true, "Modifier_Symbol": true, "Sm": true, "Math_Symbol": true, "So": true, "Other_Symbol": true, "Z": true, "Separator": true, "Zl": true, "Line_Separator": true, "Zp": true, "Paragraph_Separator": true, "Zs": true, "Space_Separator": true,
}

// Unicode 15.1 - Script values
var scriptValues = map[string]bool{
	"Adlm": true, "Adlam": true, "Aghb": true, "Caucasian_Albanian": true, "Ahom": true, "Arab": true, "Arabic": true, "Armi": true, "Imperial_Aramaic": true, "Armn": true, "Armenian": true, "Avst": true, "Avestan": true, "Bali": true, "Balinese": true, "Bamu": true, "Bamum": true, "Bass": true, "Bassa_Vah": true, "Batk": true, "Batak": true, "Beng": true, "Bengali": true, "Bhks": true, "Bhaiksuki": true, "Bopo": true, "Bopomofo": true, "Brah": true, "Brahmi": true, "Brai": true, "Braille": true, "Bugi": true, "Buginese": true, "Buhd": true, "Buhid": true, "Cakm": true, "Chakma": true, "Cans": true, "Canadian_Aboriginal": true, "Cari": true, "Carian": true, "Cham": true, "Cher": true, "Cherokee": true, "Chrs": true, "Chorasmian": true, "Copt": true, "Coptic": true, "Qaac": true, "Cpmn": true, "Cypro_Minoan": true, "Cprt": true, "Cypriot": true, "Cyrl": true, "Cyrillic": true, "Deva": true, "Devanagari": true, "Diak": true, "Dives_Akuru": true, "Dogr": true, "Dogra": true, "Dsrt": true, "Deseret": true, "Dupl": true, "Duployan": true, "Egyp": true, "Egyptian_Hieroglyphs": true, "Elba": true, "Elbasan": true, "Elym": true, "Elymaic": true, "Ethi": true, "Ethiopic": true, "Geor": true, "Georgian": true, "Glag": true, "Glagolitic": true, "Gong": true, "Gunjala_Gondi": true, "Gonm": true, "Masaram_Gondi": true, "Goth": true, "Gothic": true, "Gran": true, "Grantha": true, "Grek": true, "Greek": true, "Gujr": true, "Gujarati": true, "Guru": true, "Gurmukhi": true, "Hang": true, "Hangul": true, "Hani": true, "Han": true, "Hano": true, "Hanunoo": true, "Hatr": true, "Hatran": true, "Hebr": true, "Hebrew": true, "Hira": true, "Hiragana": true, "Hluw": true, "Anatolian_Hieroglyphs": true, "Hmng": true, "Pahawh_Hmong": true, "Hmnp": true, "Nyiakeng_Puachue_Hmong": true, "Hrkt": true, "Katakana_Or_Hiragana": true, "Hung": true, "Old_Hungarian": true, "Ital": true, "Old_Italic": true, "Java": true, "Javanese": true, "Kali": true, "Kayah_Li": true, "Kana": true, "Katakana": true, "Kawi": true, "Khar": true, "Kharoshthi": true, "Khmr": true, "Khmer": true, "Khoj": true, "Khojki": true, "Kits": true, "Khitan_Small_Script": true, "Knda": true, "Kannada": true, "Kthi": true, "Kaithi": true, "Lana": true, "Tai_Tham": true, "Laoo": true, "Lao": true, "Latn": true, "Latin": true, "Lepc": true, "Lepcha": true, "Limb": true, "Limbu": true, "Lina": true, "Linear_A": true, "Linb": true, "Linear_B": true, "Lisu": true, "Lyci": true, "Lycian": true, "Lydi": true, "Lydian": true, "Mahj": true, "Mahajani": true, "Maka": true, "Makasar": true, "Mand": true, "Mandaic": true, "Mani": true, "Manichaean": true, "Marc": true, "Marchen": true, "Medf": true, "Medefaidrin": true, "Mend": true, "Mende_Kikakui": true, "Merc": true, "Meroitic_Cursive": true, "Mero": true, "Meroitic_Hieroglyphs": true, "Mlym": true, "Malayalam": true, "Modi": true, "Mong": true, "Mongolian": true, "Mroo": true, "Mro": true, "Mtei": true, "Meetei_Mayek": true, "Mult": true, "Multani": true, "Mymr": true, "Myanmar": true, "Nagm": true, "Nag_Mundari": true, "Nand": true, "Nandinagari": true, "Narb": true, "Old_North_Arabian": true, "Nbat": true, "Nabataean": true, "Newa": true, "Nkoo": true, "Nko": true, "Nshu": true, "Nushu": true, "Ogam": true, "Ogham": true, "Olck": true, "Ol_Chiki": true, "Orkh": true, "Old_Turkic": true, "Orya": true, "Oriya": true, "Osge": true, "Osage": true, "Osma": true, "Osmanya": true, "Ougr": true, "Old_Uyghur": true, "Palm": true, "Palmyrene": true, "Pauc": true, "Pau_Cin_Hau": true, "Perm": true, "Old_Permic": true, "Phag": true, "Phags_Pa": true, "Phli": true, "Inscriptional_Pahlavi": true, "Phlp": true, "Psalter_Pahlavi": true, "Phnx": true, "Phoenician": true, "Plrd": true, "Miao": true, "Prti": true, "Inscriptional_Parthian": true, "Rjng": true, "Rejang": true, "Rohg": true, "Hanifi_Rohingya": true, "Runr": true, "Runic": true, "Samr": true, "Samaritan": true, "Sarb": true, "Old_South_Arabian": true, "Saur": true, "Saurashtra": true, "Sgnw": true, "SignWriting": true, "Shaw": true, "Shavian": true, "Shrd": true, "Sharada": true, "Sidd": true, "Siddham": true, "Sind": true, "Khudawadi": true, "Sinh": true, "Sinhala": true, "Sogd": true, "Sogdian": true, "Sogo": true, "Old_Sogdian": true, "Sora": true, "Sora_Sompeng": true, "Soyo": true, "Soyombo": true, "Sund": true, "Sundanese": true, "Sylo": true, "Syloti_Nagri": true, "Syrc": true, "Syriac": true, "Tagb": true, "Tagbanwa": true, "Takr": true, "Takri": true, "Tale": true, "Tai_Le": true, "Talu": true, "New_Tai_Lue": true, "Taml": true, "Tamil": true, "Tang": true, "Tangut": true, "Tavt": true, "Tai_Viet": true, "Telu": true, "Telugu": true, "Tfng": true, "Tifinagh": true, "Tglg": true, "Tagalog": true, "Thaa": true, "Thaana": true, "Thai": true, "Tibt": true, "Tibetan": true, "Tirh": true, "Tirhuta": true, "Tnsa": true, "Tangsa": true, "Toto": true, "Ugar": true, "Ugaritic": true, "Vaii": true, "Vai": true, "Vith": true, "Vithkuqi": true, "Wara": true, "Warang_Citi": true, "Wcho": true, "Wancho": true, "Xpeo": true, "Old_Persian": true, "Xsux": true, "Cuneiform": true, "Yezi": true, "Yezidi": true, "Yiii": true, "Yi": true, "Zanb": true, "Zanabazar_Square": true, "Zinh": true, "Inherited": true, "Qaai": true, "Zyyy": true, "Common": true, "Zzzz": true, "Unknown": true,
}

// Map of non-binary property names to their canonical names
var nonBinaryUnicodePropertyNames = map[string]string{
	"General_Category":  "General_Category",
	"gc":                "General_Category",
	"Script":            "Script",
	"sc":                "Script",
	"Script_Extensions": "Script_Extensions",
	"scx":               "Script_Extensions",
}

func (v *RegExpValidator) isValidUnicodePropertyName(name string) bool {
	return generalCategoryValues[name] || binaryUnicodeProperties[name]
}

func (v *RegExpValidator) isValidNonBinaryUnicodePropertyName(name string) bool {
	_, ok := nonBinaryUnicodePropertyNames[name]
	return ok
}

func (v *RegExpValidator) isValidUnicodeProperty(name, value string) bool {
	// Get canonical name
	canonicalName := nonBinaryUnicodePropertyNames[name]
	if canonicalName == "General_Category" {
		return generalCategoryValues[value]
	}
	if canonicalName == "Script" || canonicalName == "Script_Extensions" {
		return scriptValues[value]
	}
	return false
}

func (v *RegExpValidator) scanIdentifier(ch rune) {
	start := v.pos
	if ch != 0 && (IsIdentifierStart(ch) || ch == '_' || ch == '$') {
		v.pos++
		for v.pos < v.end {
			ch = v.charAtOffset(0)
			if IsIdentifierPart(ch) || ch == '_' || ch == '$' {
				v.pos++
			} else {
				break
			}
		}
	}
	v.tokenValue = v.text[start:v.pos]
}

func (v *RegExpValidator) scanCharacterEscape(atomEscape bool) string {
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

func (v *RegExpValidator) scanEscapeSequence(atomEscape bool) string {
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
		if v.pos >= v.end || !(v.charAtOffset(0) >= '0' && v.charAtOffset(0) <= '9') {
			return "\x00"
		}
		// This is an octal escape starting with \0
		// falls through to handle as octal
		if v.charAtOffset(0) >= '0' && v.charAtOffset(0) <= '7' {
			v.pos++
		}
		fallthrough

	case '1', '2', '3':
		// Can be up to 3 octal digits
		if v.pos < v.end && v.charAtOffset(0) >= '0' && v.charAtOffset(0) <= '7' {
			v.pos++
		}
		fallthrough

	case '4', '5', '6', '7':
		// Can be 1 or 2 octal digits (already consumed one above for 1-3)
		if v.pos < v.end && v.charAtOffset(0) >= '0' && v.charAtOffset(0) <= '7' {
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

	case 'b', 't', 'n', 'v', 'f', 'r':
		// Standard escape sequences
		return string(ch)

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
		// Parse hex value and return the character
		code := 0
		for i := hexStart; i < v.pos; i++ {
			digit := v.text[i]
			if digit >= '0' && digit <= '9' {
				code = code*16 + int(digit-'0')
			} else if digit >= 'a' && digit <= 'f' {
				code = code*16 + int(digit-'a'+10)
			} else if digit >= 'A' && digit <= 'F' {
				code = code*16 + int(digit-'A'+10)
			}
		}
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
			// Parse hex value
			code := 0
			for i := hexStart; i < v.pos-1; i++ { // -1 to skip closing brace
				digit := v.text[i]
				if digit >= '0' && digit <= '9' {
					code = code*16 + int(digit-'0')
				} else if digit >= 'a' && digit <= 'f' {
					code = code*16 + int(digit-'a'+10)
				} else if digit >= 'A' && digit <= 'F' {
					code = code*16 + int(digit-'A'+10)
				}
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
			// Parse hex value
			code := 0
			for i := hexStart; i < v.pos; i++ {
				digit := v.text[i]
				if digit >= '0' && digit <= '9' {
					code = code*16 + int(digit-'0')
				} else if digit >= 'a' && digit <= 'f' {
					code = code*16 + int(digit-'a'+10)
				} else if digit >= 'A' && digit <= 'F' {
					code = code*16 + int(digit-'A'+10)
				}
			}
			// For surrogates, we need to preserve the actual value since string(rune(surrogate))
			// converts to 0xFFFD. We encode the surrogate as UTF-16BE bytes.
			var escapedValueString string
			if code >= 0xD800 && code <= 0xDFFF {
				// Surrogate - encode as 2-byte sequence (UTF-16BE)
				escapedValueString = string([]byte{byte(code >> 8), byte(code & 0xFF)})
			} else {
				escapedValueString = string(rune(code))
			}
			// In Unicode mode, check for surrogate pairs
			if v.anyUnicodeMode && code >= 0xD800 && code <= 0xDBFF &&
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
					nextCode := 0
					for i := nextStart + 2; i < nextPos; i++ {
						digit := v.text[i]
						if digit >= '0' && digit <= '9' {
							nextCode = nextCode*16 + int(digit-'0')
						} else if digit >= 'a' && digit <= 'f' {
							nextCode = nextCode*16 + int(digit-'a'+10)
						} else if digit >= 'A' && digit <= 'F' {
							nextCode = nextCode*16 + int(digit-'A'+10)
						}
					}
					// Check if it's a low surrogate
					if nextCode >= 0xDC00 && nextCode <= 0xDFFF {
						// Combine surrogates into a single code point
						// Formula: 0x10000 + ((high - 0xD800) << 10) + (low - 0xDC00)
						highSurrogate := code
						lowSurrogate := nextCode
						combinedCodePoint := 0x10000 + ((highSurrogate - 0xD800) << 10) + (lowSurrogate - 0xDC00)
						v.pos = nextPos
						return string(rune(combinedCodePoint))
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
		if v.anyUnicodeMode || (v.anyUnicodeModeOrNonAnnexB && IsIdentifierPart(ch)) {
			v.error(diagnostics.This_character_cannot_be_escaped_in_a_regular_expression, start, v.pos-start)
		}
		return string(ch)
	}
}

// codePointAt returns the code point value at the start of the string
// If the string starts with a high surrogate followed by a low surrogate, they are combined
// Surrogates from escape sequences are encoded as 2-byte UTF-16BE sequences
func codePointAt(s string) rune {
	if len(s) == 0 {
		return 0
	}
	// Check if this is a UTF-16BE encoded surrogate (2 bytes, first byte in surrogate range)
	if len(s) == 2 {
		firstByte := uint16(s[0])
		if firstByte >= 0xD8 && firstByte <= 0xDF {
			// This is a surrogate encoded as UTF-16BE
			return rune((uint16(s[0]) << 8) | uint16(s[1]))
		}
	}
	first, size := utf8.DecodeRuneInString(s)
	// Check if it's a high surrogate (0xD800-0xDBFF)
	if first >= 0xD800 && first <= 0xDBFF && len(s) > size {
		second, _ := utf8.DecodeRuneInString(s[size:])
		// Check if it's a low surrogate (0xDC00-0xDFFF)
		if second >= 0xDC00 && second <= 0xDFFF {
			// Combine surrogates to get the code point
			// CodePoint = 0x10000 + ((high & 0x3FF) << 10) | (low & 0x3FF)
			return 0x10000 + ((first & 0x3FF) << 10) | (second & 0x3FF)
		}
	}
	return first
}

// charSize returns the number of UTF-16 code units needed to represent a code point
// This matches TypeScript's charSize function
func charSize(ch rune) int {
	if ch >= 0x10000 {
		// Code points >= 0x10000 require surrogate pairs in UTF-16 (2 code units)
		return 2
	}
	if ch == 0 {
		return 0
	}
	return 1
}

// utf16ComparisonValue returns the UTF-16 code unit value to use for range comparison
// In non-Unicode mode, JavaScript compares the low surrogate for code points >= 0x10000
func utf16ComparisonValue(s string) rune {
	if len(s) == 0 {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(s)
	// For code points >= 0x10000, return the low surrogate value
	// This matches JavaScript's behavior in non-Unicode mode where surrogate pairs
	// are treated as separate code units
	if r >= 0x10000 {
		return 0xDC00 + ((r - 0x10000) & 0x3FF)
	}
	return r
}

// stringCharSize returns the expected string length for a code point in our representation
// Surrogate pairs are represented as 2-rune strings, each rune encoded as UTF-8
func stringCharSize(ch rune) int {
	if ch >= 0x10000 {
		// High surrogate (0xD800-0xDBFF): 3 bytes in UTF-8
		// Low surrogate (0xDC00-0xDFFF): 3 bytes in UTF-8
		// Total: 6 bytes
		return 6
	}
	if ch >= 0x800 {
		return 3
	}
	if ch >= 0x80 {
		return 2
	}
	if ch == 0 {
		return 0
	}
	return 1
}

// utf16Length returns the UTF-16 length of a string, matching JavaScript's string.length
// This counts UTF-16 code units, where surrogate pairs count as 2 units
// Handles both UTF-8 encoded strings and special 2-byte UTF-16BE surrogate encodings
func utf16Length(s string) int {
	// Check if this is a 2-byte UTF-16BE surrogate encoding
	// These are used to preserve surrogate values in patterns like \uD835
	if len(s) == 2 {
		firstByte := uint16(s[0])
		if firstByte >= 0xD8 && firstByte <= 0xDF {
			// This is a surrogate encoded as UTF-16BE
			return 1
		}
	}

	// Otherwise, count UTF-16 code units from UTF-8 runes
	length := 0
	for _, r := range s {
		if r >= 0x10000 {
			length += 2 // Surrogate pair
		} else {
			length += 1
		}
	}
	return length
}

func (v *RegExpValidator) scanGroupName(isReference bool) {
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
// >= U+10000 are returned as individual surrogates (matching TypeScript scanner.ts line 3536).
func (v *RegExpValidator) scanSourceCharacter() string {
	// Check if we have a pending low surrogate from splitting a previous character
	if !v.anyUnicodeMode && v.pendingLowSurrogate != 0 {
		low := v.pendingLowSurrogate
		size := v.pendingLowSurrogate >> 16 // High 16 bits store the UTF-8 size
		v.pendingLowSurrogate = 0
		// Now advance v.pos to consume the UTF-8 character
		v.pos += int(size)
		// Return the low surrogate encoded as UTF-16BE 2-byte sequence
		return string([]byte{byte(low >> 8), byte(low & 0xFF)})
	}

	// Decode the next UTF-8 character
	r, s := utf8.DecodeRuneInString(v.text[v.pos:])
	if r == utf8.RuneError {
		v.pos++
		return v.text[v.pos-1 : v.pos]
	}

	if v.anyUnicodeMode {
		// In Unicode mode, consume the full character
		v.pos += s
		return v.text[v.pos-s : v.pos]
	}

	// In non-Unicode mode, JavaScript treats characters as UTF-16 code units.
	// For code points >= 0x10000, they are surrogate pairs, and we need to
	// return them one UTF-16 code unit at a time (like TypeScript does).
	if r >= 0x10000 {
		// This character requires a surrogate pair in UTF-16.
		// Return the high surrogate now, and save the low surrogate for next time.
		high := 0xD800 + ((r-0x10000)>>10)&0x3FF
		low := 0xDC00 + ((r - 0x10000) & 0x3FF)
		// Store the low surrogate in lower 16 bits, and UTF-8 size in upper 16 bits
		v.pendingLowSurrogate = low | (rune(s) << 16)
		// Return the high surrogate encoded as UTF-16BE 2-byte sequence
		return string([]byte{byte(high >> 8), byte(high & 0xFF)})
	}

	// For code points < 0x10000, consume and return the character normally
	v.pos += s
	return v.text[v.pos-s : v.pos]
}

// ClassRanges ::= ClassAtom ('-' ClassAtom)?
// Scans character class content like [a-z] or [^0-9].
// TypeScript reference: scanner.ts line 2990
func (v *RegExpValidator) scanClassRanges() {
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
				minCodePoint := codePointAt(atom)
				maxCodePoint := codePointAt(rangeEnd)

				// Get the expected sizes (in UTF-16 code units)
				minExpectedSize := charSize(minCodePoint)
				maxExpectedSize := charSize(maxCodePoint)

				// Check if both are "complete" characters
				// TypeScript checks: minCharacter.length === charSize(minCharacterValue)
				// where .length is the UTF-16 length in JavaScript strings.
				// We need to calculate the UTF-16 length of our strings.
				minUTF16Length := utf16Length(atom)
				maxUTF16Length := utf16Length(rangeEnd)

				// A character is complete if its UTF-16 length matches the expected size
				// TypeScript checks: minCharacter.length === charSize(minCharacterValue)
				// where .length is the UTF-16 length in JavaScript strings.
				minIsComplete := minUTF16Length == minExpectedSize
				maxIsComplete := maxUTF16Length == maxExpectedSize

				if minIsComplete && maxIsComplete && minCodePoint > maxCodePoint {
					// TypeScript compares code points directly. In non-Unicode mode,
					// literal characters >= 0x10000 are scanned as individual surrogates
					// by scanSourceCharacter(), so the code points being compared are
					// already the surrogate values (0xD800-0xDFFF).
					// Escape sequences like \u{1D608} return the full character, so the
					// code points are the actual values (>= 0x10000).
					v.error(diagnostics.Range_out_of_order_in_character_class, atomStart, v.pos-atomStart)
				}
			}
		}
	}
}

func (v *RegExpValidator) isClassContentExit(ch rune) bool {
	return ch == ']' || ch == 0 || v.pos >= v.end
}

// ClassAtom ::=
//
//	| SourceCharacter but not one of '\' or ']'
//	| '\' ClassEscape
//
// TypeScript reference: scanner.ts line 3406
func (v *RegExpValidator) scanClassAtom() string {
	if v.charAtOffset(0) == '\\' {
		v.pos++
		return v.scanClassEscape()
	}
	return v.scanSourceCharacter()
}

// ClassEscape ::=
//
//	| 'b'
//	| '-'
//	| CharacterClassEscape
//	| CharacterEscape
//
// TypeScript reference: scanner.ts line 3406
func (v *RegExpValidator) scanClassEscape() string {
	if v.scanCharacterClassEscape() {
		return ""
	}
	return v.scanCharacterEscape(false)
}

type classSetExpressionType int

const (
	classSetExpressionNone classSetExpressionType = iota
	classSetExpressionSubtraction
	classSetExpressionIntersection
)

func (v *RegExpValidator) scanClassSetExpression() {
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
				// Don't report TS1518 for the second operand of a range
				// expressionMayContainStrings tracking is still needed
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

func (v *RegExpValidator) scanClassSetSubExpression(expressionType classSetExpressionType) {
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
func (v *RegExpValidator) scanClassSetOperand() string {
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
func (v *RegExpValidator) scanClassStringDisjunctionContents() {
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
func (v *RegExpValidator) scanClassSetCharacter() string {
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
