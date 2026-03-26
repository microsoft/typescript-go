package semver

import (
	"strings"
)

// isVersionRangeChar returns true for characters valid in version range partials: [a-zA-Z0-9-+.*]
func isVersionRangeChar(c byte) bool {
	return isPreOrBuildChar(c) || c == '+' || c == '*'
}

// parseNrOrWildcard parses "0", "[1-9][0-9]*", "x", "X", or "*" from s starting at pos.
func parseNrOrWildcard(s string, pos int) (string, int, bool) {
	if pos >= len(s) {
		return "", pos, false
	}
	c := s[pos]
	if c == 'x' || c == 'X' || c == '*' {
		return string(c), pos + 1, true
	}
	return parseNr(s, pos)
}

// matchPartial parses a partial version string:
//
//	partial ::= xr ('.' xr ('.' xr qualifier?)?)?
func matchPartial(text string) (majorStr, minorStr, patchStr, prereleaseStr, buildStr string, ok bool) {
	pos := 0

	majorStr, pos, ok = parseNrOrWildcard(text, pos)
	if !ok {
		return majorStr, minorStr, patchStr, prereleaseStr, buildStr, ok
	}

	if pos < len(text) && text[pos] == '.' {
		pos++
		minorStr, pos, ok = parseNrOrWildcard(text, pos)
		if !ok {
			return majorStr, minorStr, patchStr, prereleaseStr, buildStr, ok
		}

		if pos < len(text) && text[pos] == '.' {
			pos++
			patchStr, pos, ok = parseNrOrWildcard(text, pos)
			if !ok {
				return majorStr, minorStr, patchStr, prereleaseStr, buildStr, ok
			}

			if pos < len(text) && text[pos] == '-' {
				pos++
				start := pos
				for pos < len(text) && isPreOrBuildChar(text[pos]) {
					pos++
				}
				if pos == start {
					ok = false
					return majorStr, minorStr, patchStr, prereleaseStr, buildStr, ok
				}
				prereleaseStr = text[start:pos]
			}

			if pos < len(text) && text[pos] == '+' {
				pos++
				start := pos
				for pos < len(text) && isPreOrBuildChar(text[pos]) {
					pos++
				}
				if pos == start {
					ok = false
					return majorStr, minorStr, patchStr, prereleaseStr, buildStr, ok
				}
				buildStr = text[start:pos]
			}
		}
	}

	ok = pos == len(text)
	return majorStr, minorStr, patchStr, prereleaseStr, buildStr, ok
}

// matchHyphen tries to match "partial - partial" (hyphen range) in the string.
func matchHyphen(s string) (left, right string, ok bool) {
	s = strings.TrimSpace(s)

	// Parse left: [a-zA-Z0-9-+.*]+
	pos := 0
	for pos < len(s) && isVersionRangeChar(s[pos]) {
		pos++
	}
	if pos == 0 {
		return "", "", false
	}
	left = s[:pos]

	// Require at least one whitespace before '-'
	wsStart := pos
	for pos < len(s) && (s[pos] == ' ' || s[pos] == '\t' || s[pos] == '\n' || s[pos] == '\f' || s[pos] == '\r') {
		pos++
	}
	if pos == wsStart {
		return "", "", false
	}

	// Require '-'
	if pos >= len(s) || s[pos] != '-' {
		return "", "", false
	}
	pos++

	// Require at least one whitespace after '-'
	wsStart = pos
	for pos < len(s) && (s[pos] == ' ' || s[pos] == '\t' || s[pos] == '\n' || s[pos] == '\f' || s[pos] == '\r') {
		pos++
	}
	if pos == wsStart {
		return "", "", false
	}

	// Parse right: [a-zA-Z0-9-+.*]+
	rightStart := pos
	for pos < len(s) && isVersionRangeChar(s[pos]) {
		pos++
	}
	if pos == rightStart {
		return "", "", false
	}
	right = s[rightStart:pos]

	ok = pos == len(s)
	return left, right, ok
}

// matchRange parses an optional operator followed by a version partial.
func matchRange(s string) (operator, version string, ok bool) {
	pos := 0

	// Parse optional operator: ~, ^, <, >, <=, >=, =
	if pos < len(s) {
		switch s[pos] {
		case '~', '^', '=':
			operator = string(s[pos])
			pos++
		case '<', '>':
			if pos+1 < len(s) && s[pos+1] == '=' {
				operator = s[pos : pos+2]
				pos += 2
			} else {
				operator = string(s[pos])
				pos++
			}
		}
	}

	// Skip optional whitespace
	for pos < len(s) && (s[pos] == ' ' || s[pos] == '\t' || s[pos] == '\n' || s[pos] == '\f' || s[pos] == '\r') {
		pos++
	}

	// Parse version: [a-zA-Z0-9-+.*]+
	versionStart := pos
	for pos < len(s) && isVersionRangeChar(s[pos]) {
		pos++
	}
	if pos == versionStart {
		return "", "", false
	}
	version = s[versionStart:pos]

	ok = pos == len(s)
	return operator, version, ok
}

type VersionRange struct {
	alternatives [][]versionComparator
}

type versionComparator struct {
	operator comparatorOperator
	operand  Version
}

type comparatorOperator string

const (
	rangeLessThan         comparatorOperator = "<"
	rangeLessThanEqual    comparatorOperator = "<="
	rangeEqual            comparatorOperator = "="
	rangeGreaterThanEqual comparatorOperator = ">="
	rangeGreaterThan      comparatorOperator = ">"
)

func (v *VersionRange) String() string {
	var sb strings.Builder
	formatDisjunction(&sb, v.alternatives)
	return sb.String()
}

func formatDisjunction(sb *strings.Builder, alternatives [][]versionComparator) {
	origLen := sb.Len()

	for i, alternative := range alternatives {
		if i > 0 {
			sb.WriteString(" || ")
		}
		formatAlternative(sb, alternative)
	}

	if sb.Len() == origLen {
		sb.WriteString("*")
	}
}

func formatAlternative(sb *strings.Builder, comparators []versionComparator) {
	for i, comparator := range comparators {
		if i > 0 {
			sb.WriteByte(' ')
		}
		formatComparator(sb, comparator)
	}
}

func formatComparator(sb *strings.Builder, comparator versionComparator) {
	sb.WriteString(string(comparator.operator))
	sb.WriteString(comparator.operand.String())
}

func (v *VersionRange) Test(version *Version) bool {
	return testDisjunction(v.alternatives, version)
}

func testDisjunction(alternatives [][]versionComparator, version *Version) bool {
	// an empty disjunction is treated as "*" (all versions)
	if len(alternatives) == 0 {
		return true
	}

	for _, alternative := range alternatives {
		if testAlternative(alternative, version) {
			return true
		}
	}

	return false
}

func testAlternative(alternative []versionComparator, version *Version) bool {
	for _, comparator := range alternative {
		if !testComparator(comparator, version) {
			return false
		}
	}
	return true
}

func testComparator(comparator versionComparator, version *Version) bool {
	cmp := version.Compare(&comparator.operand)
	switch comparator.operator {
	case rangeLessThan:
		return cmp < 0
	case rangeLessThanEqual:
		return cmp <= 0
	case rangeEqual:
		return cmp == 0
	case rangeGreaterThanEqual:
		return cmp >= 0
	case rangeGreaterThan:
		return cmp > 0
	default:
		panic("Unexpected operator: " + comparator.operator)
	}
}

func TryParseVersionRange(text string) (VersionRange, bool) {
	alternatives, ok := parseAlternatives(text)
	return VersionRange{alternatives: alternatives}, ok
}

func parseAlternatives(text string) ([][]versionComparator, bool) {
	var alternatives [][]versionComparator

	text = strings.TrimSpace(text)
	for r := range strings.SplitSeq(text, "||") {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}

		var comparators []versionComparator

		if left, right, matched := matchHyphen(r); matched {
			if parsedComparators, ok := parseHyphen(left, right); ok {
				comparators = append(comparators, parsedComparators...)
			} else {
				return nil, false
			}
		} else {
			for simple := range strings.FieldsSeq(r) {
				op, ver, matched := matchRange(strings.TrimSpace(simple))
				if !matched {
					return nil, false
				}

				if parsedComparators, ok := parseComparator(op, ver); ok {
					comparators = append(comparators, parsedComparators...)
				} else {
					return nil, false
				}
			}
		}

		alternatives = append(alternatives, comparators)
	}

	return alternatives, true
}

func parseHyphen(left, right string) ([]versionComparator, bool) {
	leftResult, leftOk := parsePartial(left)
	if !leftOk {
		return nil, false
	}

	rightResult, rightOk := parsePartial(right)
	if !rightOk {
		return nil, false
	}

	var comparators []versionComparator
	if !isWildcard(leftResult.majorStr) {
		// `MAJOR.*.*-...` gives us `>=MAJOR.0.0 ...`
		comparators = append(comparators, versionComparator{
			operator: rangeGreaterThanEqual,
			operand:  leftResult.version,
		})
	}

	if !isWildcard(rightResult.majorStr) {
		var operator comparatorOperator
		operand := rightResult.version

		switch {
		case isWildcard(rightResult.minorStr):
			// `...-MAJOR.*.*` gives us `... <(MAJOR+1).0.0`
			operand = operand.incrementMajor()
			operator = rangeLessThan
		case isWildcard(rightResult.patchStr):
			// `...-MAJOR.MINOR.*` gives us `... <MAJOR.(MINOR+1).0`
			operand = operand.incrementMinor()
			operator = rangeLessThan
		default:
			// `...-MAJOR.MINOR.PATCH` gives us `... <=MAJOR.MINOR.PATCH`
			operator = rangeLessThanEqual
		}

		comparators = append(comparators, versionComparator{
			operator: operator,
			operand:  operand,
		})
	}

	return comparators, true
}

type partialVersion struct {
	version  Version
	majorStr string
	minorStr string
	patchStr string
}

// Produces a "partial" version
func parsePartial(text string) (partialVersion, bool) {
	majorStr, minorStr, patchStr, prereleaseStr, buildStr, ok := matchPartial(text)
	if !ok {
		return partialVersion{}, false
	}

	if minorStr == "" {
		minorStr = "*"
	}
	if patchStr == "" {
		patchStr = "*"
	}

	var majorNumeric, minorNumeric, patchNumeric uint32
	var err error

	if isWildcard(majorStr) {
		majorNumeric = 0
		minorNumeric = 0
		patchNumeric = 0
	} else {
		majorNumeric, err = getUintComponent(majorStr)
		if err != nil {
			return partialVersion{}, false
		}

		if isWildcard(minorStr) {
			minorNumeric = 0
			patchNumeric = 0
		} else {
			minorNumeric, err = getUintComponent(minorStr)
			if err != nil {
				return partialVersion{}, false
			}

			if isWildcard(patchStr) {
				patchNumeric = 0
			} else {
				patchNumeric, err = getUintComponent(patchStr)
				if err != nil {
					return partialVersion{}, false
				}
			}
		}
	}

	var prerelease []string
	if prereleaseStr != "" {
		prerelease = strings.Split(prereleaseStr, ".")
	}

	var build []string
	if buildStr != "" {
		build = strings.Split(buildStr, ".")
	}

	result := partialVersion{
		version: Version{
			major:      majorNumeric,
			minor:      minorNumeric,
			patch:      patchNumeric,
			prerelease: prerelease,
			build:      build,
		},
		majorStr: majorStr,
		minorStr: minorStr,
		patchStr: patchStr,
	}

	return result, true
}

func parseComparator(op string, text string) ([]versionComparator, bool) {
	operator := comparatorOperator(op)

	result, ok := parsePartial(text)
	if !ok {
		return nil, false
	}

	var comparatorsResult []versionComparator

	if !isWildcard(result.majorStr) {
		switch operator {
		case "~":
			first := versionComparator{rangeGreaterThanEqual, result.version}

			var secondVersion Version
			if isWildcard(result.minorStr) {
				secondVersion = result.version.incrementMajor()
			} else {
				secondVersion = result.version.incrementMinor()
			}

			second := versionComparator{rangeLessThan, secondVersion}
			comparatorsResult = []versionComparator{first, second}

		case "^":
			first := versionComparator{rangeGreaterThanEqual, result.version}

			var secondVersion Version
			if result.version.major > 0 || isWildcard(result.minorStr) {
				secondVersion = result.version.incrementMajor()
			} else if result.version.minor > 0 || isWildcard(result.patchStr) {
				secondVersion = result.version.incrementMinor()
			} else {
				secondVersion = result.version.incrementPatch()
			}
			second := versionComparator{rangeLessThan, secondVersion}
			comparatorsResult = []versionComparator{first, second}

		case "<", ">=":
			version := result.version
			if isWildcard(result.minorStr) || isWildcard(result.patchStr) {
				version.prerelease = []string{"0"}
			}
			comparatorsResult = []versionComparator{
				{operator, version},
			}

		case "<=", ">":
			version := result.version
			if isWildcard(result.minorStr) {
				if operator == rangeLessThanEqual {
					operator = rangeLessThan
				} else {
					operator = rangeGreaterThanEqual
				}

				version = version.incrementMajor()
				version.prerelease = []string{"0"}
			} else if isWildcard(result.patchStr) {
				if operator == rangeLessThanEqual {
					operator = rangeLessThan
				} else {
					operator = rangeGreaterThanEqual
				}

				version = version.incrementMinor()
				version.prerelease = []string{"0"}
			}

			comparatorsResult = []versionComparator{
				{operator, version},
			}
		case "=", "":
			// normalize empty string to `=`
			operator = rangeEqual

			if isWildcard(result.minorStr) || isWildcard(result.patchStr) {
				originalVersion := result.version

				firstVersion := originalVersion
				firstVersion.prerelease = []string{"0"}

				var secondVersion Version
				if isWildcard(result.minorStr) {
					secondVersion = originalVersion.incrementMajor()
				} else {
					secondVersion = originalVersion.incrementMinor()
				}
				secondVersion.prerelease = []string{"0"}

				comparatorsResult = []versionComparator{
					{rangeGreaterThanEqual, firstVersion},
					{rangeLessThan, secondVersion},
				}
			} else {
				comparatorsResult = []versionComparator{
					{operator, result.version},
				}
			}
		default:
			panic("Unexpected operator: " + operator)
		}
	} else {
		if operator == "<" || operator == ">" {
			comparatorsResult = []versionComparator{
				// < 0.0.0-0
				{rangeLessThan, versionZero},
			}
		}
	}

	return comparatorsResult, true
}

func isWildcard(text string) bool {
	return text == "*" || text == "x" || text == "X"
}
