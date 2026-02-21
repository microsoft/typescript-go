package scanner

import (
	"strconv"
	"unicode/utf8"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/stringutil"
)

type RegularExpressionFlags int32

const (
	RegularExpressionFlagsNone       RegularExpressionFlags = 0
	RegularExpressionFlagsHasIndices RegularExpressionFlags = 1 << 0 // d
	RegularExpressionFlagsGlobal     RegularExpressionFlags = 1 << 1 // g
	RegularExpressionFlagsIgnoreCase RegularExpressionFlags = 1 << 2 // i
	RegularExpressionFlagsMultiline  RegularExpressionFlags = 1 << 3 // m
	RegularExpressionFlagsDotAll     RegularExpressionFlags = 1 << 4 // s
	RegularExpressionFlagsUnicode    RegularExpressionFlags = 1 << 5 // u
	RegularExpressionFlagsUnicodeSets RegularExpressionFlags = 1 << 6 // v
	RegularExpressionFlagsSticky     RegularExpressionFlags = 1 << 7 // y
	RegularExpressionFlagsUnicodeMode RegularExpressionFlags = RegularExpressionFlagsUnicode | RegularExpressionFlagsUnicodeSets
	RegularExpressionFlagsModifiers  RegularExpressionFlags = RegularExpressionFlagsIgnoreCase | RegularExpressionFlagsMultiline | RegularExpressionFlagsDotAll
)

var charToRegExpFlag = map[rune]RegularExpressionFlags{
	'd': RegularExpressionFlagsHasIndices,
	'g': RegularExpressionFlagsGlobal,
	'i': RegularExpressionFlagsIgnoreCase,
	'm': RegularExpressionFlagsMultiline,
	's': RegularExpressionFlagsDotAll,
	'u': RegularExpressionFlagsUnicode,
	'v': RegularExpressionFlagsUnicodeSets,
	'y': RegularExpressionFlagsSticky,
}

var regExpFlagToFirstAvailableLanguageVersion = map[RegularExpressionFlags]core.ScriptTarget{
	RegularExpressionFlagsHasIndices:  core.ScriptTargetES2022,
	RegularExpressionFlagsGlobal:      core.ScriptTargetES5, // ES3 doesn't exist in Go; ES5 is the lowest
	RegularExpressionFlagsIgnoreCase:  core.ScriptTargetES5,
	RegularExpressionFlagsMultiline:   core.ScriptTargetES5,
	RegularExpressionFlagsDotAll:      core.ScriptTargetES2018,
	RegularExpressionFlagsUnicode:     core.ScriptTargetES2015,
	RegularExpressionFlagsUnicodeSets: core.ScriptTargetESNext,
	RegularExpressionFlagsSticky:      core.ScriptTargetES2015,
}

func CharacterToRegularExpressionFlag(ch rune) (RegularExpressionFlags, bool) {
	flag, ok := charToRegExpFlag[ch]
	return flag, ok
}

type classSetExpressionType int

const (
	classSetExpressionTypeUnknown classSetExpressionType = iota
	classSetExpressionTypeClassUnion
	classSetExpressionTypeClassIntersection
	classSetExpressionTypeClassSubtraction
)

type groupNameReference struct {
	pos  int
	end  int
	name string
}

type decimalEscapeValue struct {
	pos   int
	end   int
	value int
}

type regExpParser struct {
	scanner        *Scanner
	end            int
	regExpFlags    RegularExpressionFlags
	isUnterminated bool
	unicodeMode    bool
	unicodeSetsMode bool

	mayContainStrings       bool
	numberOfCapturingGroups int
	groupSpecifiers         map[string]bool
	groupNameReferences     []groupNameReference
	decimalEscapes          []decimalEscapeValue
	namedCapturingGroups    []map[string]bool
}

func (p *regExpParser) pos() int {
	return p.scanner.pos
}

func (p *regExpParser) setPos(v int) {
	p.scanner.pos = v
}

func (p *regExpParser) incPos(n int) {
	p.scanner.pos += n
}

func (p *regExpParser) charAt(pos int) rune {
	if pos < p.end {
		return rune(p.scanner.text[pos])
	}
	return -1
}

func (p *regExpParser) error(msg *diagnostics.Message, pos int, length int, args ...any) {
	p.scanner.errorAt(msg, pos, length, args...)
}

func (p *regExpParser) text() string {
	return p.scanner.text
}

// Disjunction ::= Alternative ('|' Alternative)*
func (p *regExpParser) scanDisjunction(isInGroup bool) {
	for {
		p.namedCapturingGroups = append(p.namedCapturingGroups, make(map[string]bool))
		p.scanAlternative(isInGroup)
		p.namedCapturingGroups = p.namedCapturingGroups[:len(p.namedCapturingGroups)-1]
		if p.charAt(p.pos()) != '|' {
			return
		}
		p.incPos(1)
	}
}

// Alternative ::= Term*
// Term ::=
//     | Assertion
//     | Atom Quantifier?
// Assertion ::=
//     | '^'
//     | '$'
//     | '\b'
//     | '\B'
//     | '(?=' Disjunction ')'
//     | '(?!' Disjunction ')'
//     | '(?<=' Disjunction ')'
//     | '(?<!' Disjunction ')'
// Quantifier ::= QuantifierPrefix '?'?
// QuantifierPrefix ::=
//     | '*'
//     | '+'
//     | '?'
//     | '{' DecimalDigits (',' DecimalDigits?)? '}'
// Atom ::=
//     | PatternCharacter
//     | '.'
//     | '\' AtomEscape
//     | CharacterClass
//     | '(?<' RegExpIdentifierName '>' Disjunction ')'
//     | '(?' RegularExpressionFlags ('-' RegularExpressionFlags)? ':' Disjunction ')'
// CharacterClass ::= unicodeMode
//     ? '[' ClassRanges ']'
//     : '[' ClassSetExpression ']'
func (p *regExpParser) scanAlternative(isInGroup bool) {
	isPreviousTermQuantifiable := false
	for p.pos() < p.end {
		start := p.pos()
		ch := p.charAt(p.pos())
		switch ch {
		case '^', '$':
			p.incPos(1)
			isPreviousTermQuantifiable = false
		case '\\':
			p.incPos(1)
			switch p.charAt(p.pos()) {
			case 'b', 'B':
				p.incPos(1)
				isPreviousTermQuantifiable = false
			default:
				p.scanAtomEscape()
				isPreviousTermQuantifiable = true
			}
		case '(':
			p.incPos(1)
			if p.charAt(p.pos()) == '?' {
				p.incPos(1)
				switch p.charAt(p.pos()) {
				case '=', '!':
					p.incPos(1)
					isPreviousTermQuantifiable = false
				case '<':
					groupNameStart := p.pos()
					p.incPos(1)
					switch p.charAt(p.pos()) {
					case '=', '!':
						p.incPos(1)
						isPreviousTermQuantifiable = false
					default:
						p.scanGroupName(false /*isReference*/)
						p.scanExpectedChar('>')
						if p.scanner.languageVersion() < core.ScriptTargetES2018 {
							p.error(diagnostics.Named_capturing_groups_are_only_available_when_targeting_ES2018_or_later, groupNameStart, p.pos()-groupNameStart)
						}
						p.numberOfCapturingGroups++
						isPreviousTermQuantifiable = true
					}
				default:
					flagsStart := p.pos()
					setFlags := p.scanPatternModifiers(RegularExpressionFlagsNone)
					if p.charAt(p.pos()) == '-' {
						p.incPos(1)
						p.scanPatternModifiers(setFlags)
						if p.pos() == flagsStart+1 {
							p.error(diagnostics.Subpattern_flags_must_be_present_when_there_is_a_minus_sign, flagsStart, p.pos()-flagsStart)
						}
					}
					p.scanExpectedChar(':')
					isPreviousTermQuantifiable = true
				}
			} else {
				p.numberOfCapturingGroups++
				isPreviousTermQuantifiable = true
			}
			p.scanDisjunction(true /*isInGroup*/)
			p.scanExpectedChar(')')
		case '{':
			p.incPos(1)
			digitsStart := p.pos()
			p.scanDigits()
			min := p.scanner.tokenValue
			if p.charAt(p.pos()) == ',' {
				p.incPos(1)
				p.scanDigits()
				max := p.scanner.tokenValue
				if min == "" {
					if max != "" || p.charAt(p.pos()) == '}' {
						p.error(diagnostics.Incomplete_quantifier_Digit_expected, digitsStart, 0)
					} else {
						if p.unicodeMode {
							p.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, start, 1, string(ch))
						}
						isPreviousTermQuantifiable = true
						continue
					}
				}
				if max != "" {
					minVal, _ := strconv.Atoi(min)
					maxVal, _ := strconv.Atoi(max)
					if minVal > maxVal {
						p.error(diagnostics.Numbers_out_of_order_in_quantifier, digitsStart, p.pos()-digitsStart)
					}
				}
			} else if min == "" {
				if p.unicodeMode {
					p.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, start, 1, string(ch))
				}
				isPreviousTermQuantifiable = true
				continue
			}
			p.scanExpectedChar('}')
			p.incPos(-1)
			// falls through to quantifier handling
			fallthrough
		case '*', '+', '?':
			p.incPos(1)
			if p.charAt(p.pos()) == '?' {
				p.incPos(1)
			}
			if !isPreviousTermQuantifiable {
				p.error(diagnostics.There_is_nothing_available_for_repetition, start, p.pos()-start)
			}
			isPreviousTermQuantifiable = false
		case '.':
			p.incPos(1)
			isPreviousTermQuantifiable = true
		case '[':
			p.incPos(1)
			if p.unicodeSetsMode {
				p.scanClassSetExpression()
			} else {
				p.scanClassRanges()
			}
			p.scanExpectedChar(']')
			isPreviousTermQuantifiable = true
		case ')':
			if isInGroup {
				return
			}
			fallthrough
		case ']', '}':
			if p.isUnterminated && !isInGroup {
				return
			}
			if p.unicodeMode || ch == ')' {
				p.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, p.pos(), 1, string(ch))
			}
			p.incPos(1)
			isPreviousTermQuantifiable = true
		case '/', '|':
			return
		default:
			p.scanSourceCharacter()
			isPreviousTermQuantifiable = true
		}
	}
}

func (p *regExpParser) scanPatternModifiers(currFlags RegularExpressionFlags) RegularExpressionFlags {
	for p.pos() < p.end {
		ch := p.charAt(p.pos())
		if !IsIdentifierPart(ch) {
			break
		}
		flag, ok := CharacterToRegularExpressionFlag(ch)
		if !ok {
			p.error(diagnostics.Unknown_regular_expression_flag, p.pos(), 1)
		} else if currFlags&flag != 0 {
			p.error(diagnostics.Duplicate_regular_expression_flag, p.pos(), 1)
		} else if flag&RegularExpressionFlagsModifiers == 0 {
			p.error(diagnostics.This_regular_expression_flag_cannot_be_toggled_within_a_subpattern, p.pos(), 1)
		} else {
			currFlags |= flag
			availableFrom := regExpFlagToFirstAvailableLanguageVersion[flag]
			if p.scanner.languageVersion() < availableFrom {
				p.error(diagnostics.This_regular_expression_flag_is_only_available_when_targeting_0_or_later, p.pos(), 1, GetNameOfScriptTarget(availableFrom))
			}
		}
		p.incPos(1)
	}
	return currFlags
}

// AtomEscape ::=
//     | DecimalEscape
//     | CharacterClassEscape
//     | CharacterEscape
//     | 'k<' RegExpIdentifierName '>'
func (p *regExpParser) scanAtomEscape() {
	switch p.charAt(p.pos()) {
	case 'k':
		p.incPos(1)
		if p.charAt(p.pos()) == '<' {
			p.incPos(1)
			p.scanGroupName(true /*isReference*/)
			p.scanExpectedChar('>')
		} else if p.unicodeMode {
			p.error(diagnostics.X_k_must_be_followed_by_a_capturing_group_name_enclosed_in_angle_brackets, p.pos()-2, 2)
		}
	case 'q':
		if p.unicodeSetsMode {
			p.incPos(1)
			p.error(diagnostics.X_q_is_only_available_inside_character_class, p.pos()-2, 2)
			return
		}
		fallthrough
	default:
		if !p.scanCharacterClassEscape() && !p.scanDecimalEscape() {
			p.scanCharacterEscape(true /*atomEscape*/)
		}
	}
}

// DecimalEscape ::= [1-9] [0-9]*
func (p *regExpParser) scanDecimalEscape() bool {
	ch := p.charAt(p.pos())
	if ch >= '1' && ch <= '9' {
		start := p.pos()
		p.scanDigits()
		val, _ := strconv.Atoi(p.scanner.tokenValue)
		p.decimalEscapes = append(p.decimalEscapes, decimalEscapeValue{pos: start, end: p.pos(), value: val})
		return true
	}
	return false
}

// CharacterEscape ::=
//     | `c` ControlLetter
//     | IdentityEscape
//     | (Other sequences handled by `scanEscapeSequence`)
// IdentityEscape ::=
//     | '^' | '$' | '/' | '\' | '.' | '*' | '+' | '?' | '(' | ')' | '[' | ']' | '{' | '}' | '|'
//     | [~UnicodeMode] (any other non-identifier characters)
func (p *regExpParser) scanCharacterEscape(atomEscape bool) string {
	ch := p.charAt(p.pos())
	switch ch {
	case 'c':
		p.incPos(1)
		ch = p.charAt(p.pos())
		if stringutil.IsASCIILetter(ch) {
			p.incPos(1)
			return string(rune(ch & 0x1f))
		}
		if p.unicodeMode {
			p.error(diagnostics.X_c_must_be_followed_by_an_ASCII_letter, p.pos()-2, 2)
		}
		return string(ch)
	case '^', '$', '/', '\\', '.', '*', '+', '?', '(', ')', '[', ']', '{', '}', '|':
		p.incPos(1)
		return string(ch)
	default:
		if p.pos() >= p.end {
			p.error(diagnostics.Undetermined_character_escape, p.pos()-1, 1)
			return "\\"
		}
		p.incPos(-1) // back up to include the backslash for scanEscapeSequence
		flags := EscapeSequenceScanningFlagsRegularExpression | EscapeSequenceScanningFlagsAnnexB
		if p.unicodeMode {
			flags |= EscapeSequenceScanningFlagsReportErrors | EscapeSequenceScanningFlagsAnyUnicodeMode
		}
		if atomEscape {
			flags |= EscapeSequenceScanningFlagsAtomEscape
		}
		return p.scanner.scanEscapeSequence(flags)
	}
}

func (p *regExpParser) scanGroupName(isReference bool) {
	p.scanner.tokenStart = p.pos()
	p.scanner.scanIdentifier(0)
	if p.pos() == p.scanner.tokenStart {
		p.error(diagnostics.Expected_a_capturing_group_name, p.pos(), 0)
	} else if isReference {
		p.groupNameReferences = append(p.groupNameReferences, groupNameReference{pos: p.scanner.tokenStart, end: p.pos(), name: p.scanner.tokenValue})
	} else if p.namedCapturingGroupsContains(p.scanner.tokenValue) {
		p.error(diagnostics.Named_capturing_groups_with_the_same_name_must_be_mutually_exclusive_to_each_other, p.scanner.tokenStart, p.pos()-p.scanner.tokenStart)
	} else {
		if len(p.namedCapturingGroups) > 0 {
			p.namedCapturingGroups[len(p.namedCapturingGroups)-1][p.scanner.tokenValue] = true
		}
		p.groupSpecifiers[p.scanner.tokenValue] = true
	}
}

func (p *regExpParser) namedCapturingGroupsContains(name string) bool {
	for _, group := range p.namedCapturingGroups {
		if group[name] {
			return true
		}
	}
	return false
}

func (p *regExpParser) isClassContentExit(ch rune) bool {
	return ch == ']' || p.pos() >= p.end
}

// ClassRanges ::= '^'? (ClassAtom ('-' ClassAtom)?)*
func (p *regExpParser) scanClassRanges() {
	if p.charAt(p.pos()) == '^' {
		p.incPos(1)
	}
	for p.pos() < p.end {
		ch := p.charAt(p.pos())
		if p.isClassContentExit(ch) {
			return
		}
		minStart := p.pos()
		minCharacter := p.scanClassAtom()
		if p.charAt(p.pos()) == '-' {
			p.incPos(1)
			ch = p.charAt(p.pos())
			if p.isClassContentExit(ch) {
				return
			}
			if minCharacter == "" {
				p.error(diagnostics.A_character_class_range_must_not_be_bounded_by_another_character_class, minStart, p.pos()-1-minStart)
			}
			maxStart := p.pos()
			maxCharacter := p.scanClassAtom()
			if maxCharacter == "" {
				p.error(diagnostics.A_character_class_range_must_not_be_bounded_by_another_character_class, maxStart, p.pos()-maxStart)
				continue
			}
			if minCharacter == "" {
				continue
			}
			minCharacterValue, minSize := utf8.DecodeRuneInString(minCharacter)
			maxCharacterValue, maxSize := utf8.DecodeRuneInString(maxCharacter)
			if len(minCharacter) == minSize && len(maxCharacter) == maxSize && minCharacterValue > maxCharacterValue {
				p.error(diagnostics.Range_out_of_order_in_character_class, minStart, p.pos()-minStart)
			}
		}
	}
}

// ClassSetExpression ::= '^'? (ClassUnion | ClassIntersection | ClassSubtraction)
// ClassUnion ::= (ClassSetRange | ClassSetOperand)*
// ClassIntersection ::= ClassSetOperand ('&&' ClassSetOperand)+
// ClassSubtraction ::= ClassSetOperand ('--' ClassSetOperand)+
// ClassSetRange ::= ClassSetCharacter '-' ClassSetCharacter
func (p *regExpParser) scanClassSetExpression() {
	isCharacterComplement := false
	if p.charAt(p.pos()) == '^' {
		p.incPos(1)
		isCharacterComplement = true
	}
	expressionMayContainStrings := false
	ch := p.charAt(p.pos())
	if p.isClassContentExit(ch) {
		return
	}
	start := p.pos()
	var operand string
	twoChars := ""
	if p.pos()+1 < p.end {
		twoChars = p.text()[p.pos() : p.pos()+2]
	}
	switch twoChars {
	case "--", "&&":
		p.error(diagnostics.Expected_a_class_set_operand, p.pos(), 0)
		p.mayContainStrings = false
	default:
		operand = p.scanClassSetOperand()
	}
	switch p.charAt(p.pos()) {
	case '-':
		if p.pos()+1 < p.end && p.charAt(p.pos()+1) == '-' {
			if isCharacterComplement && p.mayContainStrings {
				p.error(diagnostics.Anything_that_would_possibly_match_more_than_a_single_character_is_invalid_inside_a_negated_character_class, start, p.pos()-start)
			}
			expressionMayContainStrings = p.mayContainStrings
			p.scanClassSetSubExpression(classSetExpressionTypeClassSubtraction)
			p.mayContainStrings = !isCharacterComplement && expressionMayContainStrings
			return
		}
	case '&':
		if p.pos()+1 < p.end && p.charAt(p.pos()+1) == '&' {
			p.scanClassSetSubExpression(classSetExpressionTypeClassIntersection)
			if isCharacterComplement && p.mayContainStrings {
				p.error(diagnostics.Anything_that_would_possibly_match_more_than_a_single_character_is_invalid_inside_a_negated_character_class, start, p.pos()-start)
			}
			expressionMayContainStrings = p.mayContainStrings
			p.mayContainStrings = !isCharacterComplement && expressionMayContainStrings
			return
		} else {
			p.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, p.pos(), 1, string(ch))
		}
	default:
		if isCharacterComplement && p.mayContainStrings {
			p.error(diagnostics.Anything_that_would_possibly_match_more_than_a_single_character_is_invalid_inside_a_negated_character_class, start, p.pos()-start)
		}
		expressionMayContainStrings = p.mayContainStrings
	}
	for p.pos() < p.end {
		ch = p.charAt(p.pos())
		switch ch {
		case '-':
			p.incPos(1)
			ch = p.charAt(p.pos())
			if p.isClassContentExit(ch) {
				p.mayContainStrings = !isCharacterComplement && expressionMayContainStrings
				return
			}
			if ch == '-' {
				p.incPos(1)
				p.error(diagnostics.Operators_must_not_be_mixed_within_a_character_class_Wrap_it_in_a_nested_class_instead, p.pos()-2, 2)
				start = p.pos() - 2
				operand = p.text()[start:p.pos()]
				continue
			} else {
				if operand == "" {
					p.error(diagnostics.A_character_class_range_must_not_be_bounded_by_another_character_class, start, p.pos()-1-start)
				}
				secondStart := p.pos()
				secondOperand := p.scanClassSetOperand()
				if isCharacterComplement && p.mayContainStrings {
					p.error(diagnostics.Anything_that_would_possibly_match_more_than_a_single_character_is_invalid_inside_a_negated_character_class, secondStart, p.pos()-secondStart)
				}
				expressionMayContainStrings = expressionMayContainStrings || p.mayContainStrings
				if secondOperand == "" {
					p.error(diagnostics.A_character_class_range_must_not_be_bounded_by_another_character_class, secondStart, p.pos()-secondStart)
				} else if operand != "" {
					minCharacterValue, minSize := utf8.DecodeRuneInString(operand)
					maxCharacterValue, maxSize := utf8.DecodeRuneInString(secondOperand)
					if len(operand) == minSize && len(secondOperand) == maxSize && minCharacterValue > maxCharacterValue {
						p.error(diagnostics.Range_out_of_order_in_character_class, start, p.pos()-start)
					}
				}
			}
		case '&':
			start = p.pos()
			p.incPos(1)
			if p.charAt(p.pos()) == '&' {
				p.incPos(1)
				p.error(diagnostics.Operators_must_not_be_mixed_within_a_character_class_Wrap_it_in_a_nested_class_instead, p.pos()-2, 2)
				if p.charAt(p.pos()) == '&' {
					p.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, p.pos(), 1, string(ch))
					p.incPos(1)
				}
			} else {
				p.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, p.pos()-1, 1, string(ch))
			}
			operand = p.text()[start:p.pos()]
			continue
		}
		if p.isClassContentExit(p.charAt(p.pos())) {
			break
		}
		start = p.pos()
		twoChars = ""
		if p.pos()+1 < p.end {
			twoChars = p.text()[p.pos() : p.pos()+2]
		}
		switch twoChars {
		case "--", "&&":
			p.error(diagnostics.Operators_must_not_be_mixed_within_a_character_class_Wrap_it_in_a_nested_class_instead, p.pos(), 2)
			p.incPos(2)
			operand = p.text()[start:p.pos()]
		default:
			operand = p.scanClassSetOperand()
		}
	}
	p.mayContainStrings = !isCharacterComplement && expressionMayContainStrings
}

func (p *regExpParser) scanClassSetSubExpression(expressionType classSetExpressionType) {
	expressionMayContainStrings := p.mayContainStrings
	for p.pos() < p.end {
		ch := p.charAt(p.pos())
		if p.isClassContentExit(ch) {
			break
		}
		switch ch {
		case '-':
			p.incPos(1)
			if p.charAt(p.pos()) == '-' {
				p.incPos(1)
				if expressionType != classSetExpressionTypeClassSubtraction {
					p.error(diagnostics.Operators_must_not_be_mixed_within_a_character_class_Wrap_it_in_a_nested_class_instead, p.pos()-2, 2)
				}
			} else {
				p.error(diagnostics.Operators_must_not_be_mixed_within_a_character_class_Wrap_it_in_a_nested_class_instead, p.pos()-1, 1)
			}
		case '&':
			p.incPos(1)
			if p.charAt(p.pos()) == '&' {
				p.incPos(1)
				if expressionType != classSetExpressionTypeClassIntersection {
					p.error(diagnostics.Operators_must_not_be_mixed_within_a_character_class_Wrap_it_in_a_nested_class_instead, p.pos()-2, 2)
				}
				if p.charAt(p.pos()) == '&' {
					p.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, p.pos(), 1, string(ch))
					p.incPos(1)
				}
			} else {
				p.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, p.pos()-1, 1, string(ch))
			}
		default:
			switch expressionType {
			case classSetExpressionTypeClassSubtraction:
				p.error(diagnostics.X_0_expected, p.pos(), 0, "--")
			case classSetExpressionTypeClassIntersection:
				p.error(diagnostics.X_0_expected, p.pos(), 0, "&&")
			}
		}
		ch = p.charAt(p.pos())
		if p.isClassContentExit(ch) {
			p.error(diagnostics.Expected_a_class_set_operand, p.pos(), 0)
			break
		}
		p.scanClassSetOperand()
		if expressionType == classSetExpressionTypeClassIntersection {
			expressionMayContainStrings = expressionMayContainStrings && p.mayContainStrings
		}
	}
	p.mayContainStrings = expressionMayContainStrings
}

// ClassSetOperand ::=
//     | '[' ClassSetExpression ']'
//     | '\' CharacterClassEscape
//     | '\q{' ClassStringDisjunctionContents '}'
//     | ClassSetCharacter
func (p *regExpParser) scanClassSetOperand() string {
	p.mayContainStrings = false
	switch p.charAt(p.pos()) {
	case '[':
		p.incPos(1)
		p.scanClassSetExpression()
		p.scanExpectedChar(']')
		return ""
	case '\\':
		p.incPos(1)
		if p.scanCharacterClassEscape() {
			return ""
		} else if p.charAt(p.pos()) == 'q' {
			p.incPos(1)
			if p.charAt(p.pos()) == '{' {
				p.incPos(1)
				p.scanClassStringDisjunctionContents()
				p.scanExpectedChar('}')
				return ""
			} else {
				p.error(diagnostics.X_q_must_be_followed_by_string_alternatives_enclosed_in_braces, p.pos()-2, 2)
				return "q"
			}
		}
		p.incPos(-1)
		fallthrough
	default:
		return p.scanClassSetCharacter()
	}
}

// ClassStringDisjunctionContents ::= ClassSetCharacter* ('|' ClassSetCharacter*)*
func (p *regExpParser) scanClassStringDisjunctionContents() {
	characterCount := 0
	for p.pos() < p.end {
		ch := p.charAt(p.pos())
		switch ch {
		case '}':
			if characterCount != 1 {
				p.mayContainStrings = true
			}
			return
		case '|':
			if characterCount != 1 {
				p.mayContainStrings = true
			}
			p.incPos(1)
			characterCount = 0
		default:
			p.scanClassSetCharacter()
			characterCount++
		}
	}
}

// ClassSetCharacter ::=
//     | SourceCharacter -- ClassSetSyntaxCharacter -- ClassSetReservedDoublePunctuator
//     | '\' (CharacterEscape | ClassSetReservedPunctuator | 'b')
func (p *regExpParser) scanClassSetCharacter() string {
	ch := p.charAt(p.pos())
	if ch == '\\' {
		p.incPos(1)
		innerCh := p.charAt(p.pos())
		switch innerCh {
		case 'b':
			p.incPos(1)
			return "\b"
		case '&', '-', '!', '#', '%', ',', ':', ';', '<', '=', '>', '@', '`', '~':
			p.incPos(1)
			return string(innerCh)
		default:
			return p.scanCharacterEscape(false /*atomEscape*/)
		}
	} else if p.pos()+1 < p.end && ch == rune(p.scanner.text[p.pos()+1]) {
		switch ch {
		case '&', '!', '#', '%', '*', '+', ',', '.', ':', ';', '<', '=', '>', '?', '@', '`', '~':
			p.error(diagnostics.A_character_class_must_not_contain_a_reserved_double_punctuator_Did_you_mean_to_escape_it_with_backslash, p.pos(), 2)
			p.incPos(2)
			return p.text()[p.pos()-2 : p.pos()]
		}
	}
	switch ch {
	case '/', '(', ')', '[', ']', '{', '}', '-', '|':
		p.error(diagnostics.Unexpected_0_Did_you_mean_to_escape_it_with_backslash, p.pos(), 1, string(ch))
		p.incPos(1)
		return string(ch)
	}
	return p.scanSourceCharacter()
}

// ClassAtom ::=
//     | SourceCharacter but not one of '\' or ']'
//     | '\' ClassEscape
// ClassEscape ::=
//     | 'b'
//     | '-'
//     | CharacterClassEscape
//     | CharacterEscape
func (p *regExpParser) scanClassAtom() string {
	if p.charAt(p.pos()) == '\\' {
		p.incPos(1)
		ch := p.charAt(p.pos())
		switch ch {
		case 'b':
			p.incPos(1)
			return "\b"
		case '-':
			p.incPos(1)
			return string(ch)
		default:
			if p.scanCharacterClassEscape() {
				return ""
			}
			return p.scanCharacterEscape(false /*atomEscape*/)
		}
	} else {
		return p.scanSourceCharacter()
	}
}

// CharacterClassEscape ::=
//     | 'd' | 'D' | 's' | 'S' | 'w' | 'W'
//     | [+UnicodeMode] ('P' | 'p') '{' UnicodePropertyValueExpression '}'
func (p *regExpParser) scanCharacterClassEscape() bool {
	isCharacterComplement := false
	start := p.pos() - 1
	ch := p.charAt(p.pos())
	switch ch {
	case 'd', 'D', 's', 'S', 'w', 'W':
		p.incPos(1)
		return true
	case 'P':
		isCharacterComplement = true
		fallthrough
	case 'p':
		p.incPos(1)
		if p.charAt(p.pos()) == '{' {
			p.incPos(1)
			propertyNameOrValueStart := p.pos()
			propertyNameOrValue := p.scanWordCharacters()
			if p.charAt(p.pos()) == '=' {
				propertyName := nonBinaryUnicodeProperties[propertyNameOrValue]
				if p.pos() == propertyNameOrValueStart {
					p.error(diagnostics.Expected_a_Unicode_property_name, p.pos(), 0)
				} else if propertyName == "" {
					p.error(diagnostics.Unknown_Unicode_property_name, propertyNameOrValueStart, p.pos()-propertyNameOrValueStart)
					suggestion := p.getSpellingSuggestionForUnicodePropertyName(propertyNameOrValue)
					if suggestion != "" {
						p.error(diagnostics.Did_you_mean_0, propertyNameOrValueStart, p.pos()-propertyNameOrValueStart, suggestion)
					}
				}
				p.incPos(1)
				propertyValueStart := p.pos()
				propertyValue := p.scanWordCharacters()
				if p.pos() == propertyValueStart {
					p.error(diagnostics.Expected_a_Unicode_property_value, p.pos(), 0)
				} else if propertyName != "" {
					values := valuesOfNonBinaryUnicodeProperties[propertyName]
					if values != nil && !values[propertyValue] {
						p.error(diagnostics.Unknown_Unicode_property_value, propertyValueStart, p.pos()-propertyValueStart)
						suggestion := p.getSpellingSuggestionForUnicodePropertyValue(propertyName, propertyValue)
						if suggestion != "" {
							p.error(diagnostics.Did_you_mean_0, propertyValueStart, p.pos()-propertyValueStart, suggestion)
						}
					}
				}
			} else {
				if p.pos() == propertyNameOrValueStart {
					p.error(diagnostics.Expected_a_Unicode_property_name_or_value, p.pos(), 0)
				} else if binaryUnicodePropertiesOfStrings[propertyNameOrValue] {
					if !p.unicodeSetsMode {
						p.error(diagnostics.Any_Unicode_property_that_would_possibly_match_more_than_a_single_character_is_only_available_when_the_Unicode_Sets_v_flag_is_set, propertyNameOrValueStart, p.pos()-propertyNameOrValueStart)
					} else if isCharacterComplement {
						p.error(diagnostics.Anything_that_would_possibly_match_more_than_a_single_character_is_invalid_inside_a_negated_character_class, propertyNameOrValueStart, p.pos()-propertyNameOrValueStart)
					} else {
						p.mayContainStrings = true
					}
				} else if !valuesOfNonBinaryUnicodeProperties["General_Category"][propertyNameOrValue] && !binaryUnicodeProperties[propertyNameOrValue] {
					p.error(diagnostics.Unknown_Unicode_property_name_or_value, propertyNameOrValueStart, p.pos()-propertyNameOrValueStart)
					suggestion := p.getSpellingSuggestionForUnicodePropertyNameOrValue(propertyNameOrValue)
					if suggestion != "" {
						p.error(diagnostics.Did_you_mean_0, propertyNameOrValueStart, p.pos()-propertyNameOrValueStart, suggestion)
					}
				}
			}
			p.scanExpectedChar('}')
			if !p.unicodeMode {
				p.error(diagnostics.Unicode_property_value_expressions_are_only_available_when_the_Unicode_u_flag_or_the_Unicode_Sets_v_flag_is_set, start, p.pos()-start)
			}
		} else if p.unicodeMode {
			p.error(diagnostics.X_0_must_be_followed_by_a_Unicode_property_value_expression_enclosed_in_braces, p.pos()-2, 2, string(ch))
		}
		return true
	}
	return false
}

func (p *regExpParser) getSpellingSuggestionForUnicodePropertyName(name string) string {
	candidates := make([]string, 0, len(nonBinaryUnicodeProperties))
	for k := range nonBinaryUnicodeProperties {
		candidates = append(candidates, k)
	}
	return core.GetSpellingSuggestion(name, candidates, func(s string) string { return s })
}

func (p *regExpParser) getSpellingSuggestionForUnicodePropertyValue(propertyName string, value string) string {
	values := valuesOfNonBinaryUnicodeProperties[propertyName]
	if values == nil {
		return ""
	}
	candidates := make([]string, 0, len(values))
	for k := range values {
		candidates = append(candidates, k)
	}
	return core.GetSpellingSuggestion(value, candidates, func(s string) string { return s })
}

func (p *regExpParser) getSpellingSuggestionForUnicodePropertyNameOrValue(name string) string {
	var candidates []string
	for k := range valuesOfNonBinaryUnicodeProperties["General_Category"] {
		candidates = append(candidates, k)
	}
	for k := range binaryUnicodeProperties {
		candidates = append(candidates, k)
	}
	for k := range binaryUnicodePropertiesOfStrings {
		candidates = append(candidates, k)
	}
	return core.GetSpellingSuggestion(name, candidates, func(s string) string { return s })
}

func (p *regExpParser) scanWordCharacters() string {
	start := p.pos()
	for p.pos() < p.end {
		ch := p.charAt(p.pos())
		if !isWordCharacter(ch) {
			break
		}
		p.incPos(1)
	}
	return p.text()[start:p.pos()]
}

func (p *regExpParser) scanSourceCharacter() string {
	if p.unicodeMode {
		ch, size := utf8.DecodeRuneInString(p.text()[p.pos():])
		if size == 0 || ch == utf8.RuneError {
			return ""
		}
		p.incPos(size)
		return string(ch)
	}
	if p.pos() < p.end {
		ch := p.text()[p.pos()]
		p.incPos(1)
		return string(ch)
	}
	return ""
}

func (p *regExpParser) scanExpectedChar(ch rune) {
	if p.charAt(p.pos()) == ch {
		p.incPos(1)
	} else {
		p.error(diagnostics.X_0_expected, p.pos(), 0, string(ch))
	}
}

func (p *regExpParser) scanDigits() {
	start := p.pos()
	for p.pos() < p.end && stringutil.IsDigit(p.charAt(p.pos())) {
		p.incPos(1)
	}
	p.scanner.tokenValue = p.text()[start:p.pos()]
}

func (p *regExpParser) run() {
	p.scanDisjunction(false /*isInGroup*/)

	for _, reference := range p.groupNameReferences {
		if !p.groupSpecifiers[reference.name] {
			p.error(diagnostics.There_is_no_capturing_group_named_0_in_this_regular_expression, reference.pos, reference.end-reference.pos, reference.name)
			if len(p.groupSpecifiers) > 0 {
				specifiers := make([]string, 0, len(p.groupSpecifiers))
				for k := range p.groupSpecifiers {
					specifiers = append(specifiers, k)
				}
				suggestion := core.GetSpellingSuggestion(reference.name, specifiers, func(s string) string { return s })
				if suggestion != "" {
					p.error(diagnostics.Did_you_mean_0, reference.pos, reference.end-reference.pos, suggestion)
				}
			}
		}
	}
	for _, escape := range p.decimalEscapes {
		if escape.value > p.numberOfCapturingGroups {
			if p.numberOfCapturingGroups > 0 {
				p.error(diagnostics.This_backreference_refers_to_a_group_that_does_not_exist_There_are_only_0_capturing_groups_in_this_regular_expression, escape.pos, escape.end-escape.pos, p.numberOfCapturingGroups)
			} else {
				p.error(diagnostics.This_backreference_refers_to_a_group_that_does_not_exist_There_are_no_capturing_groups_in_this_regular_expression, escape.pos, escape.end-escape.pos)
			}
		}
	}
}

func GetNameOfScriptTarget(target core.ScriptTarget) string {
	switch target {
	case core.ScriptTargetES5:
		return "es5"
	case core.ScriptTargetES2015:
		return "es2015"
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
	case core.ScriptTargetESNext:
		return "esnext"
	default:
		return ""
	}
}
