package scanner

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/stringutil"
)

// RegularExpressionFlags mirrors TypeScript's RegularExpressionFlags enum
type RegularExpressionFlags int32

const (
	RegularExpressionFlagsNone           RegularExpressionFlags = 0
	RegularExpressionFlagsHasIndices     RegularExpressionFlags = 1 << 0 // d
	RegularExpressionFlagsGlobal         RegularExpressionFlags = 1 << 1 // g
	RegularExpressionFlagsIgnoreCase     RegularExpressionFlags = 1 << 2 // i
	RegularExpressionFlagsMultiline      RegularExpressionFlags = 1 << 3 // m
	RegularExpressionFlagsDotAll         RegularExpressionFlags = 1 << 4 // s
	RegularExpressionFlagsUnicode        RegularExpressionFlags = 1 << 5 // u
	RegularExpressionFlagsUnicodeSets    RegularExpressionFlags = 1 << 6 // v
	RegularExpressionFlagsSticky         RegularExpressionFlags = 1 << 7 // y
	RegularExpressionFlagsAnyUnicodeMode RegularExpressionFlags = RegularExpressionFlagsUnicode | RegularExpressionFlagsUnicodeSets
	RegularExpressionFlagsModifiers      RegularExpressionFlags = RegularExpressionFlagsIgnoreCase | RegularExpressionFlagsMultiline | RegularExpressionFlagsDotAll
)

func characterCodeToRegularExpressionFlag(ch rune) RegularExpressionFlags {
	switch ch {
	case 'd':
		return RegularExpressionFlagsHasIndices
	case 'g':
		return RegularExpressionFlagsGlobal
	case 'i':
		return RegularExpressionFlagsIgnoreCase
	case 'm':
		return RegularExpressionFlagsMultiline
	case 's':
		return RegularExpressionFlagsDotAll
	case 'u':
		return RegularExpressionFlagsUnicode
	case 'v':
		return RegularExpressionFlagsUnicodeSets
	case 'y':
		return RegularExpressionFlagsSticky
	}
	return RegularExpressionFlagsNone
}

var regExpFlagToFirstAvailableLanguageVersion = map[RegularExpressionFlags]core.ScriptTarget{
	RegularExpressionFlagsHasIndices:  core.ScriptTargetES2022,
	RegularExpressionFlagsDotAll:      core.ScriptTargetES2018,
	RegularExpressionFlagsUnicode:     core.ScriptTargetES2015,
	RegularExpressionFlagsUnicodeSets: core.ScriptTargetES2024,
	RegularExpressionFlagsSticky:      core.ScriptTargetES2015,
}

func (s *Scanner) checkRegularExpressionFlagAvailability(flag RegularExpressionFlags, size int) {
	availableFrom, ok := regExpFlagToFirstAvailableLanguageVersion[flag]
	if ok && s.languageVersion < availableFrom {
		s.errorAt(diagnostics.This_regular_expression_flag_is_only_available_when_targeting_0_or_later, s.pos, size, scriptTargetName(availableFrom))
	}
}

func scriptTargetName(target core.ScriptTarget) string {
	switch target {
	case core.ScriptTargetES5:
		return "es5"
	case core.ScriptTargetES2015:
		return "es6"
	case core.ScriptTargetES2016:
		return "es2016"
	case core.ScriptTargetES2017:
		return "es2017"
	case core.ScriptTargetES2018:
		return "es2018"
	case core.ScriptTargetES2019:
		return "es2019"
	case core.ScriptTargetES2020:
		return "es2020"
	case core.ScriptTargetES2021:
		return "es2021"
	case core.ScriptTargetES2022:
		return "es2022"
	case core.ScriptTargetES2023:
		return "es2023"
	case core.ScriptTargetES2024:
		return "es2024"
	case core.ScriptTargetES2025:
		return "es2025"
	case core.ScriptTargetESNext:
		return "esnext"
	default:
		return strings.ToLower(target.String())
	}
}

type classSetExpressionType int

const (
	classSetExpressionTypeUnknown classSetExpressionType = iota
	classSetExpressionTypeClassUnion
	classSetExpressionTypeClassIntersection
	classSetExpressionTypeClassSubtraction
)

type groupNameRef struct {
	pos  int
	end  int
	name string
}

type decimalEscapeRef struct {
	pos   int
	end   int
	value int
}

type regExpWorker struct {
	s   *Scanner
	end int // end of the regex body text

	// Regex flags / grammar parameters
	unicodeSetsMode           bool
	anyUnicodeMode            bool
	anyUnicodeModeOrNonAnnexB bool
	annexB                    bool
	namedCaptureGroups        bool

	// State
	mayContainStrings       bool
	numberOfCapturingGroups int
	groupSpecifiers         map[string]bool
	groupNameReferences     []groupNameRef
	decimalEscapes          []decimalEscapeRef

	namedCapturingGroupsScopeStack []map[string]bool
	topNamedCapturingGroupsScope   map[string]bool
}

func newRegExpWorker(s *Scanner, end int, flags RegularExpressionFlags, annexB bool, namedCaptureGroups bool) *regExpWorker {
	unicodeSetsMode := flags&RegularExpressionFlagsUnicodeSets != 0
	anyUnicodeMode := flags&RegularExpressionFlagsAnyUnicodeMode != 0
	return &regExpWorker{
		s:                         s,
		end:                       end,
		unicodeSetsMode:           unicodeSetsMode,
		anyUnicodeMode:            anyUnicodeMode,
		anyUnicodeModeOrNonAnnexB: anyUnicodeMode || !annexB,
		annexB:                    annexB,
		namedCaptureGroups:        namedCaptureGroups,
	}
}

// charCodeChecked returns the rune at pos, or -1 if pos is out of bounds
func (w *regExpWorker) charCodeChecked(pos int) rune {
	if pos >= 0 && pos < w.end {
		r, _ := utf8.DecodeRuneInString(w.s.text[pos:])
		return r
	}
	return -1
}

// charCodeUnchecked returns the rune at pos without bounds checking
func (w *regExpWorker) charCodeUnchecked(pos int) rune {
	r, _ := utf8.DecodeRuneInString(w.s.text[pos:])
	return r
}

// charSize returns the UTF-8 byte size of a rune
func (w *regExpWorker) charSize(ch rune) int {
	return utf8.RuneLen(ch)
}

func (w *regExpWorker) error(msg *diagnostics.Message, pos int, length int, args ...any) {
	w.s.errorAt(msg, pos, length, args...)
}

func (w *regExpWorker) getText(start, end int) string {
	if start >= len(w.s.text) || end > len(w.s.text) || start < 0 {
		return ""
	}
	return w.s.text[start:end]
}

func (w *regExpWorker) run() {
	w.scanDisjunction(false)

	// Validate group name references
	for _, reference := range w.groupNameReferences {
		if w.groupSpecifiers == nil || !w.groupSpecifiers[reference.name] {
			w.error(diagnostics.There_is_no_capturing_group_named_0_in_this_regular_expression, reference.pos, reference.end-reference.pos, reference.name)
			if w.groupSpecifiers != nil {
				candidates := make([]string, 0, len(w.groupSpecifiers))
				for k := range w.groupSpecifiers {
					candidates = append(candidates, k)
				}
				suggestion := core.GetSpellingSuggestion(reference.name, candidates, func(s string) string { return s })
				if suggestion != "" {
					w.error(diagnostics.Did_you_mean_0, reference.pos, reference.end-reference.pos, suggestion)
				}
			}
		}
	}

	// Validate decimal escapes
	for _, escape := range w.decimalEscapes {
		if escape.value > w.numberOfCapturingGroups {
			if w.numberOfCapturingGroups > 0 {
				w.error(diagnostics.This_backreference_refers_to_a_group_that_does_not_exist_There_are_only_0_capturing_groups_in_this_regular_expression, escape.pos, escape.end-escape.pos, w.numberOfCapturingGroups)
			} else {
				w.error(diagnostics.This_backreference_refers_to_a_group_that_does_not_exist_There_are_no_capturing_groups_in_this_regular_expression, escape.pos, escape.end-escape.pos)
			}
		}
	}
}

// Disjunction ::= Alternative ('|' Alternative)*
func (w *regExpWorker) scanDisjunction(isInGroup bool) {
	for {
		w.namedCapturingGroupsScopeStack = append(w.namedCapturingGroupsScopeStack, w.topNamedCapturingGroupsScope)
		w.topNamedCapturingGroupsScope = nil
		w.scanAlternative(isInGroup)
		w.topNamedCapturingGroupsScope = w.namedCapturingGroupsScopeStack[len(w.namedCapturingGroupsScopeStack)-1]
		w.namedCapturingGroupsScopeStack = w.namedCapturingGroupsScopeStack[:len(w.namedCapturingGroupsScopeStack)-1]
		if w.charCodeChecked(w.s.pos) != '|' {
			return
		}
		w.s.pos++
	}
}

// Alternative ::= Term*
func (w *regExpWorker) scanAlternative(isInGroup bool) {
	isPreviousTermQuantifiable := false
	for {
		start := w.s.pos
		ch := w.charCodeChecked(w.s.pos)
		switch ch {
		case -1: // EOF
			return
		case '^', '$':
			w.s.pos++
			isPreviousTermQuantifiable = false
		case '\\':
			w.s.pos++
			switch w.charCodeChecked(w.s.pos) {
			case 'b', 'B':
				w.s.pos++
				isPreviousTermQuantifiable = false
			default:
				w.scanAtomEscape()
				isPreviousTermQuantifiable = true
			}
		case '(':
			w.s.pos++
			if w.charCodeChecked(w.s.pos) == '?' {
				w.s.pos++
				switch w.charCodeChecked(w.s.pos) {
				case '=', '!':
					w.s.pos++
					// In Annex B, (?=Disjunction) and (?!Disjunction) are quantifiable
					isPreviousTermQuantifiable = !w.anyUnicodeModeOrNonAnnexB
				case '<':
					groupNameStart := w.s.pos
					w.s.pos++
					switch w.charCodeChecked(w.s.pos) {
					case '=', '!':
						w.s.pos++
						isPreviousTermQuantifiable = false
					default:
						w.scanGroupName(false)
						w.scanExpectedChar('>')
						if w.s.languageVersion < core.ScriptTargetES2018 {
							w.error(diagnostics.Named_capturing_groups_are_only_available_when_targeting_ES2018_or_later, groupNameStart, w.s.pos-groupNameStart)
						}
						w.numberOfCapturingGroups++
						isPreviousTermQuantifiable = true
					}
				default:
					modStart := w.s.pos
					setFlags := w.scanPatternModifiers(RegularExpressionFlagsNone)
					if w.charCodeChecked(w.s.pos) == '-' {
						w.s.pos++
						w.scanPatternModifiers(setFlags)
						if w.s.pos == modStart+1 {
							w.error(diagnostics.Subpattern_flags_must_be_present_when_there_is_a_minus_sign, modStart, w.s.pos-modStart)
						}
					}
					w.scanExpectedChar(':')
					isPreviousTermQuantifiable = true
				}
			} else {
				w.numberOfCapturingGroups++
				isPreviousTermQuantifiable = true
			}
			w.scanDisjunction(true)
			w.scanExpectedChar(')')
		case '{':
			w.s.pos++
			digitsStart := w.s.pos
			min, _ := w.s.scanDigits()
			if !w.anyUnicodeModeOrNonAnnexB && min == "" {
				isPreviousTermQuantifiable = true
				break
			}
			if w.charCodeChecked(w.s.pos) == ',' {
				w.s.pos++
				max, _ := w.s.scanDigits()
				if min == "" {
					if max != "" || w.charCodeChecked(w.s.pos) == '}' {
						w.error(diagnostics.Incomplete_quantifier_Digit_expected, digitsStart, 0)
					} else {
						w.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, start, 1, string(ch))
						isPreviousTermQuantifiable = true
						break
					}
				} else if max != "" {
					minVal, _ := strconv.Atoi(min)
					maxVal, _ := strconv.Atoi(max)
					if minVal > maxVal && (w.anyUnicodeModeOrNonAnnexB || w.charCodeChecked(w.s.pos) == '}') {
						w.error(diagnostics.Numbers_out_of_order_in_quantifier, digitsStart, w.s.pos-digitsStart)
					}
				}
			} else if min == "" {
				if w.anyUnicodeModeOrNonAnnexB {
					w.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, start, 1, string(ch))
				}
				isPreviousTermQuantifiable = true
				break
			}
			if w.charCodeChecked(w.s.pos) != '}' {
				if w.anyUnicodeModeOrNonAnnexB {
					w.error(diagnostics.X_0_expected, w.s.pos, 0, "}")
					w.s.pos--
				} else {
					isPreviousTermQuantifiable = true
					break
				}
			}
			// falls through to quantifier suffix
			fallthrough
		case '*', '+', '?':
			w.s.pos++
			if w.charCodeChecked(w.s.pos) == '?' {
				// Non-greedy
				w.s.pos++
			}
			if !isPreviousTermQuantifiable {
				w.error(diagnostics.There_is_nothing_available_for_repetition, start, w.s.pos-start)
			}
			isPreviousTermQuantifiable = false
		case '.':
			w.s.pos++
			isPreviousTermQuantifiable = true
		case '[':
			w.s.pos++
			if w.unicodeSetsMode {
				w.scanClassSetExpression()
			} else {
				w.scanClassRanges()
			}
			w.scanExpectedChar(']')
			isPreviousTermQuantifiable = true
		case ')':
			if isInGroup {
				return
			}
			// falls through
			fallthrough
		case ']', '}':
			if w.anyUnicodeModeOrNonAnnexB || ch == ')' {
				w.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, w.s.pos, 1, string(ch))
			}
			w.s.pos++
			isPreviousTermQuantifiable = true
		case '/', '|':
			return
		default:
			w.scanSourceCharacter()
			isPreviousTermQuantifiable = true
		}
	}
}

func (w *regExpWorker) scanPatternModifiers(currFlags RegularExpressionFlags) RegularExpressionFlags {
	for {
		ch := w.charCodeChecked(w.s.pos)
		if ch == -1 || !IsIdentifierPart(ch) {
			break
		}
		size := w.charSize(ch)
		flag := characterCodeToRegularExpressionFlag(ch)
		if flag == RegularExpressionFlagsNone {
			w.error(diagnostics.Unknown_regular_expression_flag, w.s.pos, size)
		} else if currFlags&flag != 0 {
			w.error(diagnostics.Duplicate_regular_expression_flag, w.s.pos, size)
		} else if flag&RegularExpressionFlagsModifiers == 0 {
			w.error(diagnostics.This_regular_expression_flag_cannot_be_toggled_within_a_subpattern, w.s.pos, size)
		} else {
			currFlags |= flag
			w.s.checkRegularExpressionFlagAvailability(flag, size)
		}
		w.s.pos += size
	}
	return currFlags
}

// AtomEscape ::=
//
//	| DecimalEscape
//	| CharacterClassEscape
//	| CharacterEscape
//	| 'k<' RegExpIdentifierName '>'
func (w *regExpWorker) scanAtomEscape() {
	switch w.charCodeChecked(w.s.pos) {
	case 'k':
		w.s.pos++
		if w.charCodeChecked(w.s.pos) == '<' {
			w.s.pos++
			w.scanGroupName(true)
			w.scanExpectedChar('>')
		} else if w.anyUnicodeModeOrNonAnnexB || w.namedCaptureGroups {
			w.error(diagnostics.X_k_must_be_followed_by_a_capturing_group_name_enclosed_in_angle_brackets, w.s.pos-2, 2)
		}
	case 'q':
		if w.unicodeSetsMode {
			w.s.pos++
			w.error(diagnostics.X_q_is_only_available_inside_character_class, w.s.pos-2, 2)
			return
		}
		// falls through
		fallthrough
	default:
		if !w.scanCharacterClassEscape() && !w.scanDecimalEscape() {
			w.scanCharacterEscape(true)
		}
	}
}

// DecimalEscape ::= [1-9] [0-9]*
func (w *regExpWorker) scanDecimalEscape() bool {
	ch := w.charCodeChecked(w.s.pos)
	if ch >= '1' && ch <= '9' {
		start := w.s.pos
		digits, _ := w.s.scanDigits()
		val, _ := strconv.Atoi(digits)
		w.decimalEscapes = append(w.decimalEscapes, decimalEscapeRef{pos: start, end: w.s.pos, value: val})
		return true
	}
	return false
}

// CharacterEscape ::=
//
//	| `c` ControlLetter
//	| IdentityEscape
//	| (Other sequences handled by `scanEscapeSequence`)
func (w *regExpWorker) scanCharacterEscape(atomEscape bool) string {
	ch := w.charCodeChecked(w.s.pos)
	switch ch {
	case -1: // EOF
		w.error(diagnostics.Undetermined_character_escape, w.s.pos-1, 1)
		return "\\"
	case 'c':
		w.s.pos++
		ch = w.charCodeChecked(w.s.pos)
		if stringutil.IsASCIILetter(ch) {
			w.s.pos++
			return string(rune(ch & 0x1f))
		}
		if w.anyUnicodeModeOrNonAnnexB {
			w.error(diagnostics.X_c_must_be_followed_by_an_ASCII_letter, w.s.pos-2, 2)
		} else if atomEscape {
			// Annex B treats
			//   ExtendedAtom : `\` [lookahead = `c`]
			// as the single character `\` when `c` isn't followed by a valid control character
			w.s.pos--
			return "\\"
		}
		return string(ch)
	case '^', '$', '/', '\\', '.', '*', '+', '?', '(', ')', '[', ']', '{', '}', '|':
		w.s.pos++
		return string(ch)
	default:
		w.s.pos--
		var flags EscapeSequenceScanningFlags
		flags = EscapeSequenceScanningFlagsRegularExpression
		if w.annexB {
			flags |= EscapeSequenceScanningFlagsAnnexB
		}
		if w.anyUnicodeMode {
			flags |= EscapeSequenceScanningFlagsAnyUnicodeMode
		}
		if atomEscape {
			flags |= EscapeSequenceScanningFlagsAtomEscape
		}
		return w.s.scanEscapeSequence(flags)
	}
}

func (w *regExpWorker) scanGroupName(isReference bool) {
	w.s.tokenStart = w.s.pos
	w.s.scanIdentifier(0)
	if w.s.pos == w.s.tokenStart {
		w.error(diagnostics.Expected_a_capturing_group_name, w.s.pos, 0)
	} else if isReference {
		w.groupNameReferences = append(w.groupNameReferences, groupNameRef{pos: w.s.tokenStart, end: w.s.pos, name: w.s.tokenValue})
	} else if w.topNamedCapturingGroupsScope != nil && w.topNamedCapturingGroupsScope[w.s.tokenValue] || w.scopeStackHas(w.s.tokenValue) {
		w.error(diagnostics.Named_capturing_groups_with_the_same_name_must_be_mutually_exclusive_to_each_other, w.s.tokenStart, w.s.pos-w.s.tokenStart)
	} else {
		if w.topNamedCapturingGroupsScope == nil {
			w.topNamedCapturingGroupsScope = make(map[string]bool)
		}
		w.topNamedCapturingGroupsScope[w.s.tokenValue] = true
		if w.groupSpecifiers == nil {
			w.groupSpecifiers = make(map[string]bool)
		}
		w.groupSpecifiers[w.s.tokenValue] = true
	}
}

func (w *regExpWorker) scopeStackHas(name string) bool {
	for _, scope := range w.namedCapturingGroupsScopeStack {
		if scope != nil && scope[name] {
			return true
		}
	}
	return false
}

func (w *regExpWorker) isClassContentExit(ch rune) bool {
	return ch == ']' || ch == -1 || w.s.pos >= w.end
}

// ClassRanges ::= '^'? (ClassAtom ('-' ClassAtom)?)*
func (w *regExpWorker) scanClassRanges() {
	if w.charCodeChecked(w.s.pos) == '^' {
		// character complement
		w.s.pos++
	}
	for {
		ch := w.charCodeChecked(w.s.pos)
		if w.isClassContentExit(ch) {
			return
		}
		minStart := w.s.pos
		minCharacter := w.scanClassAtom()
		if w.charCodeChecked(w.s.pos) == '-' {
			w.s.pos++
			ch = w.charCodeChecked(w.s.pos)
			if w.isClassContentExit(ch) {
				return
			}
			if minCharacter == "" && w.anyUnicodeModeOrNonAnnexB {
				w.error(diagnostics.A_character_class_range_must_not_be_bounded_by_another_character_class, minStart, w.s.pos-1-minStart)
			}
			maxStart := w.s.pos
			maxCharacter := w.scanClassAtom()
			if maxCharacter == "" && w.anyUnicodeModeOrNonAnnexB {
				w.error(diagnostics.A_character_class_range_must_not_be_bounded_by_another_character_class, maxStart, w.s.pos-maxStart)
				continue
			}
			if minCharacter == "" {
				continue
			}
			minCharacterValue, _ := utf8.DecodeRuneInString(minCharacter)
			maxCharacterValue, _ := utf8.DecodeRuneInString(maxCharacter)
			if len(minCharacter) == w.charSize(minCharacterValue) &&
				len(maxCharacter) == w.charSize(maxCharacterValue) &&
				minCharacterValue > maxCharacterValue {
				w.error(diagnostics.Range_out_of_order_in_character_class, minStart, w.s.pos-minStart)
			}
		}
	}
}

// ClassSetExpression ::= '^'? (ClassUnion | ClassIntersection | ClassSubtraction)
func (w *regExpWorker) scanClassSetExpression() {
	isCharacterComplement := false
	if w.charCodeChecked(w.s.pos) == '^' {
		w.s.pos++
		isCharacterComplement = true
	}
	expressionMayContainStrings := false
	ch := w.charCodeChecked(w.s.pos)
	if w.isClassContentExit(ch) {
		return
	}
	start := w.s.pos
	var operand string
	twoChars := w.getText(w.s.pos, w.s.pos+2)
	if twoChars == "--" || twoChars == "&&" {
		w.error(diagnostics.Expected_a_class_set_operand, w.s.pos, 0)
		w.mayContainStrings = false
	} else {
		operand = w.scanClassSetOperand()
	}
	switch w.charCodeChecked(w.s.pos) {
	case '-':
		if w.charCodeChecked(w.s.pos+1) == '-' {
			if isCharacterComplement && w.mayContainStrings {
				w.error(diagnostics.Anything_that_would_possibly_match_more_than_a_single_character_is_invalid_inside_a_negated_character_class, start, w.s.pos-start)
			}
			expressionMayContainStrings = w.mayContainStrings
			w.scanClassSetSubExpression(classSetExpressionTypeClassSubtraction)
			if !isCharacterComplement {
				w.mayContainStrings = expressionMayContainStrings
			} else {
				w.mayContainStrings = false
			}
			return
		}
	case '&':
		if w.charCodeChecked(w.s.pos+1) == '&' {
			w.scanClassSetSubExpression(classSetExpressionTypeClassIntersection)
			if isCharacterComplement && w.mayContainStrings {
				w.error(diagnostics.Anything_that_would_possibly_match_more_than_a_single_character_is_invalid_inside_a_negated_character_class, start, w.s.pos-start)
			}
			expressionMayContainStrings = w.mayContainStrings
			if !isCharacterComplement {
				w.mayContainStrings = expressionMayContainStrings
			} else {
				w.mayContainStrings = false
			}
			return
		}
		w.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, w.s.pos, 1, string(ch))
	default:
		if isCharacterComplement && w.mayContainStrings {
			w.error(diagnostics.Anything_that_would_possibly_match_more_than_a_single_character_is_invalid_inside_a_negated_character_class, start, w.s.pos-start)
		}
		expressionMayContainStrings = w.mayContainStrings
	}
	for {
		ch = w.charCodeChecked(w.s.pos)
		if ch == -1 {
			break
		}
		switch ch {
		case '-':
			w.s.pos++
			ch = w.charCodeChecked(w.s.pos)
			if w.isClassContentExit(ch) {
				if !isCharacterComplement {
					w.mayContainStrings = expressionMayContainStrings
				} else {
					w.mayContainStrings = false
				}
				return
			}
			if ch == '-' {
				w.s.pos++
				w.error(diagnostics.Operators_must_not_be_mixed_within_a_character_class_Wrap_it_in_a_nested_class_instead, w.s.pos-2, 2)
				start = w.s.pos - 2
				operand = w.getText(start, w.s.pos)
				continue
			}
			if operand == "" {
				w.error(diagnostics.A_character_class_range_must_not_be_bounded_by_another_character_class, start, w.s.pos-1-start)
			}
			secondStart := w.s.pos
			secondOperand := w.scanClassSetOperand()
			if isCharacterComplement && w.mayContainStrings {
				w.error(diagnostics.Anything_that_would_possibly_match_more_than_a_single_character_is_invalid_inside_a_negated_character_class, secondStart, w.s.pos-secondStart)
			}
			expressionMayContainStrings = expressionMayContainStrings || w.mayContainStrings
			if secondOperand == "" {
				w.error(diagnostics.A_character_class_range_must_not_be_bounded_by_another_character_class, secondStart, w.s.pos-secondStart)
				break
			}
			if operand == "" {
				break
			}
			minCharacterValue, _ := utf8.DecodeRuneInString(operand)
			maxCharacterValue, _ := utf8.DecodeRuneInString(secondOperand)
			if len(operand) == w.charSize(minCharacterValue) &&
				len(secondOperand) == w.charSize(maxCharacterValue) &&
				minCharacterValue > maxCharacterValue {
				w.error(diagnostics.Range_out_of_order_in_character_class, start, w.s.pos-start)
			}
		case '&':
			start = w.s.pos
			w.s.pos++
			if w.charCodeChecked(w.s.pos) == '&' {
				w.s.pos++
				w.error(diagnostics.Operators_must_not_be_mixed_within_a_character_class_Wrap_it_in_a_nested_class_instead, w.s.pos-2, 2)
				if w.charCodeChecked(w.s.pos) == '&' {
					w.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, w.s.pos, 1, string(ch))
					w.s.pos++
				}
			} else {
				w.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, w.s.pos-1, 1, string(ch))
			}
			operand = w.getText(start, w.s.pos)
			continue
		}
		if w.isClassContentExit(w.charCodeChecked(w.s.pos)) {
			break
		}
		start = w.s.pos
		twoChars = w.getText(w.s.pos, w.s.pos+2)
		if twoChars == "--" || twoChars == "&&" {
			w.error(diagnostics.Operators_must_not_be_mixed_within_a_character_class_Wrap_it_in_a_nested_class_instead, w.s.pos, 2)
			w.s.pos += 2
			operand = w.getText(start, w.s.pos)
		} else {
			operand = w.scanClassSetOperand()
		}
	}
	if !isCharacterComplement {
		w.mayContainStrings = expressionMayContainStrings
	} else {
		w.mayContainStrings = false
	}
}

func (w *regExpWorker) scanClassSetSubExpression(expressionType classSetExpressionType) {
	expressionMayContainStrings := w.mayContainStrings
	for {
		ch := w.charCodeChecked(w.s.pos)
		if w.isClassContentExit(ch) {
			break
		}
		// Provide user-friendly diagnostic messages
		switch ch {
		case '-':
			w.s.pos++
			if w.charCodeChecked(w.s.pos) == '-' {
				w.s.pos++
				if expressionType != classSetExpressionTypeClassSubtraction {
					w.error(diagnostics.Operators_must_not_be_mixed_within_a_character_class_Wrap_it_in_a_nested_class_instead, w.s.pos-2, 2)
				}
			} else {
				w.error(diagnostics.Operators_must_not_be_mixed_within_a_character_class_Wrap_it_in_a_nested_class_instead, w.s.pos-1, 1)
			}
		case '&':
			w.s.pos++
			if w.charCodeChecked(w.s.pos) == '&' {
				w.s.pos++
				if expressionType != classSetExpressionTypeClassIntersection {
					w.error(diagnostics.Operators_must_not_be_mixed_within_a_character_class_Wrap_it_in_a_nested_class_instead, w.s.pos-2, 2)
				}
				if w.charCodeChecked(w.s.pos) == '&' {
					w.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, w.s.pos, 1, string(ch))
					w.s.pos++
				}
			} else {
				w.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, w.s.pos-1, 1, string(ch))
			}
		default:
			switch expressionType {
			case classSetExpressionTypeClassSubtraction:
				w.error(diagnostics.X_0_expected, w.s.pos, 0, "--")
			case classSetExpressionTypeClassIntersection:
				w.error(diagnostics.X_0_expected, w.s.pos, 0, "&&")
			}
		}
		ch = w.charCodeChecked(w.s.pos)
		if w.isClassContentExit(ch) {
			w.error(diagnostics.Expected_a_class_set_operand, w.s.pos, 0)
			break
		}
		w.scanClassSetOperand()
		// Used only if expressionType is Intersection
		expressionMayContainStrings = expressionMayContainStrings && w.mayContainStrings
	}
	w.mayContainStrings = expressionMayContainStrings
}

// ClassSetOperand ::=
//
//	| '[' ClassSetExpression ']'
//	| '\' CharacterClassEscape
//	| '\q{' ClassStringDisjunctionContents '}'
//	| ClassSetCharacter
func (w *regExpWorker) scanClassSetOperand() string {
	w.mayContainStrings = false
	switch w.charCodeChecked(w.s.pos) {
	case -1: // EOF
		return ""
	case '[':
		w.s.pos++
		w.scanClassSetExpression()
		w.scanExpectedChar(']')
		return ""
	case '\\':
		w.s.pos++
		if w.scanCharacterClassEscape() {
			return ""
		}
		if w.charCodeChecked(w.s.pos) == 'q' {
			w.s.pos++
			if w.charCodeChecked(w.s.pos) == '{' {
				w.s.pos++
				w.scanClassStringDisjunctionContents()
				w.scanExpectedChar('}')
				return ""
			}
			w.error(diagnostics.X_q_must_be_followed_by_string_alternatives_enclosed_in_braces, w.s.pos-2, 2)
			return "q"
		}
		w.s.pos--
		// falls through
		fallthrough
	default:
		return w.scanClassSetCharacter()
	}
}

// ClassStringDisjunctionContents ::= ClassSetCharacter* ('|' ClassSetCharacter*)*
func (w *regExpWorker) scanClassStringDisjunctionContents() {
	characterCount := 0
	for {
		ch := w.charCodeChecked(w.s.pos)
		switch ch {
		case -1: // EOF
			return
		case '}':
			if characterCount != 1 {
				w.mayContainStrings = true
			}
			return
		case '|':
			if characterCount != 1 {
				w.mayContainStrings = true
			}
			w.s.pos++
			characterCount = 0
		default:
			w.scanClassSetCharacter()
			characterCount++
		}
	}
}

// ClassSetCharacter ::=
//
//	| SourceCharacter -- ClassSetSyntaxCharacter -- ClassSetReservedDoublePunctuator
//	| '\' (CharacterEscape | ClassSetReservedPunctuator | 'b')
func (w *regExpWorker) scanClassSetCharacter() string {
	ch := w.charCodeChecked(w.s.pos)
	if ch == -1 {
		return ""
	}
	if ch == '\\' {
		w.s.pos++
		innerCh := w.charCodeChecked(w.s.pos)
		switch innerCh {
		case 'b':
			w.s.pos++
			return "\b"
		case '&', '-', '!', '#', '%', ',', ':', ';', '<', '=', '>', '@', '`', '~':
			w.s.pos++
			return string(innerCh)
		default:
			return w.scanCharacterEscape(false)
		}
	}
	// Check for reserved double punctuators
	if w.s.pos+1 < w.end {
		nextCh := w.charCodeChecked(w.s.pos + 1)
		if ch == nextCh {
			switch ch {
			case '&', '!', '#', '%', '*', '+', ',', '.', ':', ';', '<', '=', '>', '?', '@', '`', '~':
				w.error(diagnostics.A_character_class_must_not_contain_a_reserved_double_punctuator_Did_you_mean_to_escape_it_with_backslash, w.s.pos, 2)
				w.s.pos += 2
				return w.getText(w.s.pos-2, w.s.pos)
			}
		}
	}
	switch ch {
	case '/', '(', ')', '[', ']', '{', '}', '-', '|':
		w.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, w.s.pos, 1, string(ch))
		w.s.pos++
		return string(ch)
	}
	return w.scanSourceCharacter()
}

// ClassAtom ::=
//
//	| SourceCharacter but not one of '\' or ']'
//	| '\' ClassEscape
func (w *regExpWorker) scanClassAtom() string {
	if w.charCodeChecked(w.s.pos) == '\\' {
		w.s.pos++
		ch := w.charCodeChecked(w.s.pos)
		switch ch {
		case 'b':
			w.s.pos++
			return "\b"
		case '-':
			w.s.pos++
			return string(ch)
		default:
			if w.scanCharacterClassEscape() {
				return ""
			}
			return w.scanCharacterEscape(false)
		}
	}
	return w.scanSourceCharacter()
}

// CharacterClassEscape ::=
//
//	| 'd' | 'D' | 's' | 'S' | 'w' | 'W'
//	| [+UnicodeMode] ('P' | 'p') '{' UnicodePropertyValueExpression '}'
func (w *regExpWorker) scanCharacterClassEscape() bool {
	isCharacterComplement := false
	start := w.s.pos - 1
	ch := w.charCodeChecked(w.s.pos)
	switch ch {
	case 'd', 'D', 's', 'S', 'w', 'W':
		w.s.pos++
		return true
	case 'P':
		isCharacterComplement = true
		fallthrough
	case 'p':
		w.s.pos++
		if w.charCodeChecked(w.s.pos) == '{' {
			w.s.pos++
			propertyNameOrValueStart := w.s.pos
			propertyNameOrValue := w.scanWordCharacters()
			if w.charCodeChecked(w.s.pos) == '=' {
				propertyName, hasName := nonBinaryUnicodeProperties[propertyNameOrValue]
				if w.s.pos == propertyNameOrValueStart {
					w.error(diagnostics.Expected_a_Unicode_property_name, w.s.pos, 0)
				} else if !hasName {
					w.error(diagnostics.Unknown_Unicode_property_name, propertyNameOrValueStart, w.s.pos-propertyNameOrValueStart)
					candidates := make([]string, 0, len(nonBinaryUnicodeProperties))
					for k := range nonBinaryUnicodeProperties {
						candidates = append(candidates, k)
					}
					suggestion := core.GetSpellingSuggestion(propertyNameOrValue, candidates, func(s string) string { return s })
					if suggestion != "" {
						w.error(diagnostics.Did_you_mean_0, propertyNameOrValueStart, w.s.pos-propertyNameOrValueStart, suggestion)
					}
				}
				w.s.pos++
				propertyValueStart := w.s.pos
				propertyValue := w.scanWordCharacters()
				if w.s.pos == propertyValueStart {
					w.error(diagnostics.Expected_a_Unicode_property_value, w.s.pos, 0)
				} else if hasName {
					valuesSet := valuesOfNonBinaryUnicodeProperties[propertyName]
					if valuesSet != nil && !valuesSet[propertyValue] {
						w.error(diagnostics.Unknown_Unicode_property_value, propertyValueStart, w.s.pos-propertyValueStart)
						candidates := make([]string, 0, len(valuesSet))
						for k := range valuesSet {
							candidates = append(candidates, k)
						}
						suggestion := core.GetSpellingSuggestion(propertyValue, candidates, func(s string) string { return s })
						if suggestion != "" {
							w.error(diagnostics.Did_you_mean_0, propertyValueStart, w.s.pos-propertyValueStart, suggestion)
						}
					}
				}
			} else {
				if w.s.pos == propertyNameOrValueStart {
					w.error(diagnostics.Expected_a_Unicode_property_name_or_value, w.s.pos, 0)
				} else if binaryUnicodePropertiesOfStrings[propertyNameOrValue] {
					if !w.unicodeSetsMode {
						w.error(diagnostics.Any_Unicode_property_that_would_possibly_match_more_than_a_single_character_is_only_available_when_the_Unicode_Sets_v_flag_is_set, propertyNameOrValueStart, w.s.pos-propertyNameOrValueStart)
					} else if isCharacterComplement {
						w.error(diagnostics.Anything_that_would_possibly_match_more_than_a_single_character_is_invalid_inside_a_negated_character_class, propertyNameOrValueStart, w.s.pos-propertyNameOrValueStart)
					} else {
						w.mayContainStrings = true
					}
				} else if !valuesOfNonBinaryUnicodeProperties["General_Category"][propertyNameOrValue] && !binaryUnicodeProperties[propertyNameOrValue] {
					w.error(diagnostics.Unknown_Unicode_property_name_or_value, propertyNameOrValueStart, w.s.pos-propertyNameOrValueStart)
					allCandidates := make([]string, 0)
					for k := range valuesOfNonBinaryUnicodeProperties["General_Category"] {
						allCandidates = append(allCandidates, k)
					}
					for k := range binaryUnicodeProperties {
						allCandidates = append(allCandidates, k)
					}
					for k := range binaryUnicodePropertiesOfStrings {
						allCandidates = append(allCandidates, k)
					}
					suggestion := core.GetSpellingSuggestion(propertyNameOrValue, allCandidates, func(s string) string { return s })
					if suggestion != "" {
						w.error(diagnostics.Did_you_mean_0, propertyNameOrValueStart, w.s.pos-propertyNameOrValueStart, suggestion)
					}
				}
			}
			w.scanExpectedChar('}')
			if !w.anyUnicodeMode {
				w.error(diagnostics.Unicode_property_value_expressions_are_only_available_when_the_Unicode_u_flag_or_the_Unicode_Sets_v_flag_is_set, start, w.s.pos-start)
			}
		} else if w.anyUnicodeModeOrNonAnnexB {
			w.error(diagnostics.X_0_must_be_followed_by_a_Unicode_property_value_expression_enclosed_in_braces, w.s.pos-2, 2, string(ch))
		} else {
			w.s.pos--
			return false
		}
		return true
	}
	return false
}

func (w *regExpWorker) scanWordCharacters() string {
	var value strings.Builder
	for {
		ch := w.charCodeChecked(w.s.pos)
		if ch == -1 || !isWordCharacter(ch) {
			break
		}
		value.WriteRune(ch)
		w.s.pos++
	}
	return value.String()
}

func (w *regExpWorker) scanSourceCharacter() string {
	if w.anyUnicodeMode {
		ch := w.charCodeChecked(w.s.pos)
		size := w.charSize(ch)
		w.s.pos += size
		if size > 0 {
			return w.getText(w.s.pos-size, w.s.pos)
		}
		return ""
	}
	if w.s.pos < w.end {
		ch := w.s.text[w.s.pos]
		w.s.pos++
		return string(ch)
	}
	return ""
}

func (w *regExpWorker) scanExpectedChar(ch rune) {
	if w.charCodeChecked(w.s.pos) == ch {
		w.s.pos++
	} else {
		w.error(diagnostics.X_0_expected, w.s.pos, 0, string(ch))
	}
}

// Unicode property data tables

// Table 66: Non-binary Unicode property aliases and their canonical property names
// https://tc39.es/ecma262/#table-nonbinary-unicode-properties
var nonBinaryUnicodeProperties = map[string]string{
	"General_Category":  "General_Category",
	"gc":                "General_Category",
	"Script":            "Script",
	"sc":                "Script",
	"Script_Extensions": "Script_Extensions",
	"scx":               "Script_Extensions",
}

// Table 67: Binary Unicode property aliases and their canonical property names
// https://tc39.es/ecma262/#table-binary-unicode-properties
var binaryUnicodeProperties = map[string]bool{
	"ASCII": true, "ASCII_Hex_Digit": true, "AHex": true, "Alphabetic": true, "Alpha": true,
	"Any": true, "Assigned": true, "Bidi_Control": true, "Bidi_C": true, "Bidi_Mirrored": true,
	"Bidi_M": true, "Case_Ignorable": true, "CI": true, "Cased": true,
	"Changes_When_Casefolded": true, "CWCF": true, "Changes_When_Casemapped": true, "CWCM": true,
	"Changes_When_Lowercased": true, "CWL": true, "Changes_When_NFKC_Casefolded": true, "CWKCF": true,
	"Changes_When_Titlecased": true, "CWT": true, "Changes_When_Uppercased": true, "CWU": true,
	"Dash": true, "Default_Ignorable_Code_Point": true, "DI": true, "Deprecated": true, "Dep": true,
	"Diacritic": true, "Dia": true, "Emoji": true, "Emoji_Component": true, "EComp": true,
	"Emoji_Modifier": true, "EMod": true, "Emoji_Modifier_Base": true, "EBase": true,
	"Emoji_Presentation": true, "EPres": true, "Extended_Pictographic": true, "ExtPict": true,
	"Extender": true, "Ext": true, "Grapheme_Base": true, "Gr_Base": true, "Grapheme_Extend": true,
	"Gr_Ext": true, "Hex_Digit": true, "Hex": true, "IDS_Binary_Operator": true, "IDSB": true,
	"IDS_Trinary_Operator": true, "IDST": true, "ID_Continue": true, "IDC": true, "ID_Start": true,
	"IDS": true, "Ideographic": true, "Ideo": true, "Join_Control": true, "Join_C": true,
	"Logical_Order_Exception": true, "LOE": true, "Lowercase": true, "Lower": true, "Math": true,
	"Noncharacter_Code_Point": true, "NChar": true, "Pattern_Syntax": true, "Pat_Syn": true,
	"Pattern_White_Space": true, "Pat_WS": true, "Quotation_Mark": true, "QMark": true,
	"Radical": true, "Regional_Indicator": true, "RI": true, "Sentence_Terminal": true, "STerm": true,
	"Soft_Dotted": true, "SD": true, "Terminal_Punctuation": true, "Term": true,
	"Unified_Ideograph": true, "UIdeo": true, "Uppercase": true, "Upper": true,
	"Variation_Selector": true, "VS": true, "White_Space": true, "space": true,
	"XID_Continue": true, "XIDC": true, "XID_Start": true, "XIDS": true,
}

// Table 68: Binary Unicode properties of strings
// https://tc39.es/ecma262/#table-binary-unicode-properties-of-strings
var binaryUnicodePropertiesOfStrings = map[string]bool{
	"Basic_Emoji":                 true,
	"Emoji_Keycap_Sequence":       true,
	"RGI_Emoji_Modifier_Sequence": true,
	"RGI_Emoji_Flag_Sequence":     true,
	"RGI_Emoji_Tag_Sequence":      true,
	"RGI_Emoji_ZWJ_Sequence":      true,
	"RGI_Emoji":                   true,
}

// Unicode 15.1
// Values of non-binary Unicode properties
var valuesOfNonBinaryUnicodeProperties = map[string]map[string]bool{
	"General_Category": {
		"C": true, "Other": true, "Cc": true, "Control": true, "cntrl": true, "Cf": true, "Format": true,
		"Cn": true, "Unassigned": true, "Co": true, "Private_Use": true, "Cs": true, "Surrogate": true,
		"L": true, "Letter": true, "LC": true, "Cased_Letter": true, "Ll": true, "Lowercase_Letter": true,
		"Lm": true, "Modifier_Letter": true, "Lo": true, "Other_Letter": true, "Lt": true,
		"Titlecase_Letter": true, "Lu": true, "Uppercase_Letter": true, "M": true, "Mark": true,
		"Combining_Mark": true, "Mc": true, "Spacing_Mark": true, "Me": true, "Enclosing_Mark": true,
		"Mn": true, "Nonspacing_Mark": true, "N": true, "Number": true, "Nd": true, "Decimal_Number": true,
		"digit": true, "Nl": true, "Letter_Number": true, "No": true, "Other_Number": true,
		"P": true, "Punctuation": true, "punct": true, "Pc": true, "Connector_Punctuation": true,
		"Pd": true, "Dash_Punctuation": true, "Pe": true, "Close_Punctuation": true,
		"Pf": true, "Final_Punctuation": true, "Pi": true, "Initial_Punctuation": true,
		"Po": true, "Other_Punctuation": true, "Ps": true, "Open_Punctuation": true,
		"S": true, "Symbol": true, "Sc": true, "Currency_Symbol": true, "Sk": true,
		"Modifier_Symbol": true, "Sm": true, "Math_Symbol": true, "So": true, "Other_Symbol": true,
		"Z": true, "Separator": true, "Zl": true, "Line_Separator": true, "Zp": true,
		"Paragraph_Separator": true, "Zs": true, "Space_Separator": true,
	},
	"Script": {
		"Adlm": true, "Adlam": true, "Aghb": true, "Caucasian_Albanian": true, "Ahom": true,
		"Arab": true, "Arabic": true, "Armi": true, "Imperial_Aramaic": true, "Armn": true,
		"Armenian": true, "Avst": true, "Avestan": true, "Bali": true, "Balinese": true,
		"Bamu": true, "Bamum": true, "Bass": true, "Bassa_Vah": true, "Batk": true, "Batak": true,
		"Beng": true, "Bengali": true, "Bhks": true, "Bhaiksuki": true, "Bopo": true,
		"Bopomofo": true, "Brah": true, "Brahmi": true, "Brai": true, "Braille": true,
		"Bugi": true, "Buginese": true, "Buhd": true, "Buhid": true, "Cakm": true, "Chakma": true,
		"Cans": true, "Canadian_Aboriginal": true, "Cari": true, "Carian": true, "Cham": true,
		"Cher": true, "Cherokee": true, "Chrs": true, "Chorasmian": true, "Copt": true,
		"Coptic": true, "Qaac": true, "Cpmn": true, "Cypro_Minoan": true, "Cprt": true,
		"Cypriot": true, "Cyrl": true, "Cyrillic": true, "Deva": true, "Devanagari": true,
		"Diak": true, "Dives_Akuru": true, "Dogr": true, "Dogra": true, "Dsrt": true,
		"Deseret": true, "Dupl": true, "Duployan": true, "Egyp": true, "Egyptian_Hieroglyphs": true,
		"Elba": true, "Elbasan": true, "Elym": true, "Elymaic": true, "Ethi": true,
		"Ethiopic": true, "Geor": true, "Georgian": true, "Glag": true, "Glagolitic": true,
		"Gong": true, "Gunjala_Gondi": true, "Gonm": true, "Masaram_Gondi": true, "Goth": true,
		"Gothic": true, "Gran": true, "Grantha": true, "Grek": true, "Greek": true,
		"Gujr": true, "Gujarati": true, "Guru": true, "Gurmukhi": true, "Hang": true,
		"Hangul": true, "Hani": true, "Han": true, "Hano": true, "Hanunoo": true, "Hatr": true,
		"Hatran": true, "Hebr": true, "Hebrew": true, "Hira": true, "Hiragana": true,
		"Hluw": true, "Anatolian_Hieroglyphs": true, "Hmng": true, "Pahawh_Hmong": true,
		"Hmnp": true, "Nyiakeng_Puachue_Hmong": true, "Hrkt": true, "Katakana_Or_Hiragana": true,
		"Hung": true, "Old_Hungarian": true, "Ital": true, "Old_Italic": true, "Java": true,
		"Javanese": true, "Kali": true, "Kayah_Li": true, "Kana": true, "Katakana": true,
		"Kawi": true, "Khar": true, "Kharoshthi": true, "Khmr": true, "Khmer": true,
		"Khoj": true, "Khojki": true, "Kits": true, "Khitan_Small_Script": true, "Knda": true,
		"Kannada": true, "Kthi": true, "Kaithi": true, "Lana": true, "Tai_Tham": true,
		"Laoo": true, "Lao": true, "Latn": true, "Latin": true, "Lepc": true, "Lepcha": true,
		"Limb": true, "Limbu": true, "Lina": true, "Linear_A": true, "Linb": true,
		"Linear_B": true, "Lisu": true, "Lyci": true, "Lycian": true, "Lydi": true,
		"Lydian": true, "Mahj": true, "Mahajani": true, "Maka": true, "Makasar": true,
		"Mand": true, "Mandaic": true, "Mani": true, "Manichaean": true, "Marc": true,
		"Marchen": true, "Medf": true, "Medefaidrin": true, "Mend": true, "Mende_Kikakui": true,
		"Merc": true, "Meroitic_Cursive": true, "Mero": true, "Meroitic_Hieroglyphs": true,
		"Mlym": true, "Malayalam": true, "Modi": true, "Mong": true, "Mongolian": true,
		"Mroo": true, "Mro": true, "Mtei": true, "Meetei_Mayek": true, "Mult": true,
		"Multani": true, "Mymr": true, "Myanmar": true, "Nagm": true, "Nag_Mundari": true,
		"Nand": true, "Nandinagari": true, "Narb": true, "Old_North_Arabian": true, "Nbat": true,
		"Nabataean": true, "Newa": true, "Nkoo": true, "Nko": true, "Nshu": true, "Nushu": true,
		"Ogam": true, "Ogham": true, "Olck": true, "Ol_Chiki": true, "Orkh": true,
		"Old_Turkic": true, "Orya": true, "Oriya": true, "Osge": true, "Osage": true,
		"Osma": true, "Osmanya": true, "Ougr": true, "Old_Uyghur": true, "Palm": true,
		"Palmyrene": true, "Pauc": true, "Pau_Cin_Hau": true, "Perm": true, "Old_Permic": true,
		"Phag": true, "Phags_Pa": true, "Phli": true, "Inscriptional_Pahlavi": true, "Phlp": true,
		"Psalter_Pahlavi": true, "Phnx": true, "Phoenician": true, "Plrd": true, "Miao": true,
		"Prti": true, "Inscriptional_Parthian": true, "Rjng": true, "Rejang": true, "Rohg": true,
		"Hanifi_Rohingya": true, "Runr": true, "Runic": true, "Samr": true, "Samaritan": true,
		"Sarb": true, "Old_South_Arabian": true, "Saur": true, "Saurashtra": true, "Sgnw": true,
		"SignWriting": true, "Shaw": true, "Shavian": true, "Shrd": true, "Sharada": true,
		"Sidd": true, "Siddham": true, "Sind": true, "Khudawadi": true, "Sinh": true,
		"Sinhala": true, "Sogd": true, "Sogdian": true, "Sogo": true, "Old_Sogdian": true,
		"Sora": true, "Sora_Sompeng": true, "Soyo": true, "Soyombo": true, "Sund": true,
		"Sundanese": true, "Sylo": true, "Syloti_Nagri": true, "Syrc": true, "Syriac": true,
		"Tagb": true, "Tagbanwa": true, "Takr": true, "Takri": true, "Tale": true, "Tai_Le": true,
		"Talu": true, "New_Tai_Lue": true, "Taml": true, "Tamil": true, "Tang": true,
		"Tangut": true, "Tavt": true, "Tai_Viet": true, "Telu": true, "Telugu": true,
		"Tfng": true, "Tifinagh": true, "Tglg": true, "Tagalog": true, "Thaa": true,
		"Thaana": true, "Thai": true, "Tibt": true, "Tibetan": true, "Tirh": true,
		"Tirhuta": true, "Tnsa": true, "Tangsa": true, "Toto": true, "Ugar": true,
		"Ugaritic": true, "Vaii": true, "Vai": true, "Vith": true, "Vithkuqi": true,
		"Wara": true, "Warang_Citi": true, "Wcho": true, "Wancho": true, "Xpeo": true,
		"Old_Persian": true, "Xsux": true, "Cuneiform": true, "Yezi": true, "Yezidi": true,
		"Yiii": true, "Yi": true, "Zanb": true, "Zanabazar_Square": true, "Zinh": true,
		"Inherited": true, "Qaai": true, "Zyyy": true, "Common": true, "Zzzz": true, "Unknown": true,
	},
	// Script_Extensions values are the same as Script values
	"Script_Extensions": nil, // initialized in init()
}

func init() {
	valuesOfNonBinaryUnicodeProperties["Script_Extensions"] = valuesOfNonBinaryUnicodeProperties["Script"]
}
