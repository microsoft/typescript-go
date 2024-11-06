package semver

import (
	"cmp"
	"fmt"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// > A normal version number MUST take the form X.Y.Z where X, Y, and Z are non-negative
// > integers, and MUST NOT contain leading zeroes. X is the major version, Y is the minor
// > version, and Z is the patch version. Each element MUST increase numerically.
//
// NOTE: We differ here in that we allow X and X.Y, with missing parts having the default
// value of `0`.
var versionRegexp = regexp.MustCompile("^(0|[1-9]\\d*)(?:\\.(0|[1-9]\\d*)(?:\\.(0|[1-9]\\d*)(?:-([a-zA-Z0-9-.]+))?(?:\\+([a-zA-Z0-9-.]+))?)?)?$")

// https://semver.org/#spec-item-9
// > A pre-release version MAY be denoted by appending a hyphen and a series of dot separated
// > identifiers immediately following the patch version. Identifiers MUST comprise only ASCII
// > alphanumerics and hyphen [0-9A-Za-z-]. Identifiers MUST NOT be empty. Numeric identifiers
// > MUST NOT include leading zeroes.
var prereleaseRegexp = regexp.MustCompile("^(?:0|[1-9]\\d*|[a-zA-Z-][a-zA-Z0-9-]*)(?:\\.(?:0|[1-9]\\d*|[a-zA-Z-][a-zA-Z0-9-]*))*$")
var prereleasePartRegexp = regexp.MustCompile("^(?:0|[1-9]\\d*|[a-zA-Z-][a-zA-Z0-9-]*)$")

// https://semver.org/#spec-item-10
// > Build metadata MAY be denoted by appending a plus sign and a series of dot separated
// > identifiers immediately following the patch or pre-release version. Identifiers MUST
// > comprise only ASCII alphanumerics and hyphen [0-9A-Za-z-]. Identifiers MUST NOT be empty.
var buildRegExp = regexp.MustCompile("^[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$")
var buildPartRegExp = regexp.MustCompile("^[a-zA-Z0-9-]+$")

// https://semver.org/#spec-item-9
// > Numeric identifiers MUST NOT include leading zeroes.
var numericIdentifierRegExp = regexp.MustCompile("^(?:0|[1-9]\\d*)$")

type Version struct {
	// !!!
	// These should probably all be float64
	major      uint32
	minor      uint32
	patch      uint32
	prerelease []string
	build      []string
}

var _zeroPrerelease = []string{"0"}

var versionZero = Version{
	prerelease: _zeroPrerelease,
}

func incrementMajor(v Version) Version {
	return Version{
		major: v.major + 1,
	}
}

func incrementMinor(v Version) Version {
	return Version{
		major: v.major,
		minor: v.minor + 1,
	}
}

func incrementPatch(v Version) Version {
	return Version{
		major: v.major,
		minor: v.minor,
		patch: v.patch + 1,
	}
}

const comparisonLessThan = -1
const comparisonEqualTo = 0
const comparisonGreaterThan = 1

func (a *Version) Compare(b *Version) int {
	// https://semver.org/#spec-item-11
	// > Precedence is determined by the first difference when comparing each of these
	// > identifiers from left to right as follows: Major, minor, and patch versions are
	// > always compared numerically.
	//
	// https://semver.org/#spec-item-11
	// > Precedence for two pre-release versions with the same major, minor, and patch version
	// > MUST be determined by comparing each dot separated identifier from left to right until
	// > a difference is found [...]
	//
	// https://semver.org/#spec-item-11
	// > Build metadata does not figure into precedence
	switch {
	case a == b:
		return comparisonEqualTo
	case a == nil:
		return comparisonLessThan
	case b == nil:
		return comparisonGreaterThan
	}

	r := cmp.Compare(a.major, b.major)
	if r != 0 {
		return r
	}

	r = cmp.Compare(a.minor, b.minor)
	if r != 0 {
		return r
	}

	r = cmp.Compare(a.patch, b.patch)
	if r != 0 {
		return r
	}

	return comparePreReleaseIdentifiers(a.prerelease, b.prerelease)
}

func comparePreReleaseIdentifiers(left, right []string) int {
	// https://semver.org/#spec-item-11
	// > When major, minor, and patch are equal, a pre-release version has lower precedence
	// > than a normal version.
	if len(left) == 0 {
		if len(right) == 0 {
			return comparisonEqualTo
		}
		return comparisonGreaterThan
	} else if len(right) == 0 {
		return comparisonLessThan
	}

	// https://semver.org/#spec-item-11
	// > Precedence for two pre-release versions with the same major, minor, and patch version
	// > MUST be determined by comparing each dot separated identifier from left to right until
	// > a difference is found [...]
	return slices.CompareFunc(left, right, comparePreReleaseIdentifier)
}

func comparePreReleaseIdentifier(left, right string) int {
	// https://semver.org/#spec-item-11
	// > Precedence for two pre-release versions with the same major, minor, and patch version
	// > MUST be determined by comparing each dot separated identifier from left to right until
	// > a difference is found [...]
	r := strings.Compare(left, right)
	if r == 0 {
		return r
	}

	leftIsNumeric := numericIdentifierRegExp.MatchString(left)
	rightIsNumeric := numericIdentifierRegExp.MatchString(right)

	if leftIsNumeric || rightIsNumeric {
		// https://semver.org/#spec-item-11
		// > Numeric identifiers always have lower precedence than non-numeric identifiers.
		if !rightIsNumeric {
			return comparisonLessThan
		}
		if !leftIsNumeric {
			return comparisonGreaterThan
		}

		// https://semver.org/#spec-item-11
		// > identifiers consisting of only digits are compared numerically
		leftAsNumber := stringToNumber(left)
		rightAsNumber := stringToNumber(right)
		return cmp.Compare(leftAsNumber, rightAsNumber)
	}

	// https://semver.org/#spec-item-11
	// > identifiers with letters or hyphens are compared lexically in ASCII sort order.
	return strings.Compare(left, right)
}

func (v *Version) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d.%d.%d", v.major, v.minor, v.patch)
	if len(v.prerelease) > 0 {
		fmt.Fprintf(&sb, "-%s", strings.Join(v.prerelease, "."))
	}
	if len(v.build) > 0 {
		fmt.Fprintf(&sb, "+%s", strings.Join(v.build, "."))
	}
	return sb.String()
}

type SemverParseError struct {
	origInput string
}

func (e *SemverParseError) Error() string {
	return "Could not parse version string from " + e.origInput
}

func TryParseVersion(text string) (Version, error) {
	var result Version

	match := versionRegexp.FindStringSubmatch(text)
	if match == nil {
		return result, &SemverParseError{origInput: text}
	}

	majorStr := match[1]
	minorStr := match[2]
	patchStr := match[3]
	prereleaseStr := match[4]
	buildStr := match[5]

	result.major = getUintComponent(majorStr)

	if minorStr != "" {
		result.minor = getUintComponent(minorStr)
	}

	if patchStr != "" {
		result.patch = getUintComponent(patchStr)
	}

	if prereleaseStr != "" {
		if !prereleaseRegexp.MatchString(prereleaseStr) {
			return result, &SemverParseError{origInput: text}
		}

		result.prerelease = strings.Split(prereleaseStr, ".")
	}
	if buildStr != "" {
		if !buildRegExp.MatchString(buildStr) {
			return result, &SemverParseError{origInput: text}
		}

		result.build = strings.Split(buildStr, ".")
	}

	return result, nil
}

func stringToNumber(s string) float64 {
	// !!! Copied from the core compiler.
	// !!! This function should behave identically to the expression `+s` in JS
	// This includes parsing binary, octal, and hex numeric strings
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return math.NaN()
	}
	return value
}

func getUintComponent(text string) uint32 {
	r, err := strconv.ParseUint(text, 10, 32)
	if err != nil {
		// !!!
		// Should we actually just panic here?
		panic(err.Error())
	}
	return uint32(r)
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

// https://github.com/npm/node-semver#range-grammar
//
// range-set    ::= range ( logical-or range ) *
// range        ::= hyphen | simple ( ' ' simple ) * | ”
// logical-or   ::= ( ' ' ) * '||' ( ' ' ) *
var logicalOrRegExp = regexp.MustCompile("\\|\\|")
var whitespaceRegExp = regexp.MustCompile("\\s+")

// https://github.com/npm/node-semver#range-grammar
//
// partial      ::= xr ( '.' xr ( '.' xr qualifier ? )? )?
// xr           ::= 'x' | 'X' | '*' | nr
// nr           ::= '0' | ['1'-'9'] ( ['0'-'9'] ) *
// qualifier    ::= ( '-' pre )? ( '+' build )?
// pre          ::= parts
// build        ::= parts
// parts        ::= part ( '.' part ) *
// part         ::= nr | [-0-9A-Za-z]+
var partialRegExp = regexp.MustCompile("^([xX*0]|[1-9]\\d*)(?:\\.([xX*0]|[1-9]\\d*)(?:\\.([xX*0]|[1-9]\\d*)(?:-([a-zA-Z0-9-.]+))?(?:\\+([a-zA-Z0-9-.]+))?)?)?$")

// https://github.com/npm/node-semver#range-grammar
//
// hyphen       ::= partial ' - ' partial
var hyphenRegExp = regexp.MustCompile("^\\s*([a-zA-Z0-9-+.*]+)\\s+-\\s+([a-zA-Z0-9-+.*]+)\\s*$")

// https://github.com/npm/node-semver#range-grammar
//
// simple       ::= primitive | partial | tilde | caret
// primitive    ::= ( '<' | '>' | '>=' | '<=' | '=' ) partial
// tilde        ::= '~' partial
// caret        ::= '^' partial
var rangeRegExp = regexp.MustCompile("^([~^<>=]|<=|>=)?\\s*([a-zA-Z0-9-+.*]+)$")

func (v *VersionRange) String() string {
	var sb strings.Builder
	formatDisjunction(&sb, v.alternatives)
	return sb.String()
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
	ranges := logicalOrRegExp.Split(text, -1)
	for _, r := range ranges {
		// !!!
		// Slight deviation here.
		// Original impementation doesn't split *before* dismissing empty space.
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}

		var comparators []versionComparator

		if hyphenMatch := hyphenRegExp.FindStringSubmatch(r); hyphenMatch != nil {
			if parsedComparators, ok := parseHyphen(hyphenMatch[1], hyphenMatch[2]); ok {
				comparators = append(comparators, parsedComparators...)
			} else {
				return nil, false
			}
		} else {
			for _, simple := range whitespaceRegExp.Split(r, -1) {
				match := rangeRegExp.FindStringSubmatch(strings.TrimSpace(simple))
				if match == nil {
					return nil, false
				}

				if parsedComparators, ok := parseComparator(match[1], match[2]); ok {
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
			operand = incrementMajor(operand)
			operator = rangeLessThan
		case isWildcard(rightResult.patchStr):
			// `...-MAJOR.MINOR.*` gives us `... <MAJOR.(MINOR+1).0`
			operand = incrementMinor(operand)
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

func parsePartial(text string) (partialVersion, bool) {
	match := partialRegExp.FindStringSubmatch(text)
	if match == nil {
		return partialVersion{}, false
	}

	majorStr := match[1]
	minorStr := match[2]
	patchStr := match[3]
	prereleaseStr := match[4]
	buildStr := match[5]

	if minorStr == "" {
		minorStr = "*"
	}
	if patchStr == "" {
		patchStr = "*"
	}

	var majorNumeric, minorNumeric, patchNumeric uint32

	if isWildcard(majorStr) {
		majorNumeric = 0
		minorNumeric = 0
		patchNumeric = 0
	} else {
		majorNumeric = getUintComponent(majorStr)

		if isWildcard(minorStr) {
			minorNumeric = 0
			patchNumeric = 0
		} else {
			minorNumeric = getUintComponent(minorStr)

			if isWildcard(patchStr) {
				patchNumeric = 0
			} else {
				patchNumeric = getUintComponent(patchStr)
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
				secondVersion = incrementMajor(result.version)
			} else {
				secondVersion = incrementMinor(result.version)
			}

			second := versionComparator{rangeLessThan, secondVersion}
			comparatorsResult = []versionComparator{first, second}

		case "^":
			first := versionComparator{rangeGreaterThanEqual, result.version}

			var secondVersion Version
			if result.version.major > 0 || isWildcard(result.minorStr) {
				secondVersion = incrementMajor(result.version)
			} else if result.version.minor > 0 || isWildcard(result.patchStr) {
				secondVersion = incrementMinor(result.version)
			} else {
				secondVersion = incrementPatch(result.version)
			}
			second := versionComparator{rangeLessThan, secondVersion}
			comparatorsResult = []versionComparator{first, second}

		case "<", ">=":
			version := result.version
			if isWildcard(result.minorStr) || isWildcard(result.patchStr) {
				version.prerelease = _zeroPrerelease
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

				version = incrementMajor(version)
				version.prerelease = _zeroPrerelease
			} else if isWildcard(result.patchStr) {
				if operator == rangeLessThanEqual {
					operator = rangeLessThan
				} else {
					operator = rangeGreaterThanEqual
				}

				version = incrementMinor(version)
				version.prerelease = _zeroPrerelease
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
				firstVersion.prerelease = _zeroPrerelease

				var secondVersion Version
				if isWildcard(result.minorStr) {
					secondVersion = incrementMajor(originalVersion)
				} else {
					secondVersion = incrementMinor(originalVersion)
				}
				secondVersion.prerelease = _zeroPrerelease

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

func formatDisjunction(sb *strings.Builder, alternatives [][]versionComparator) {
	for i, alternative := range alternatives {
		if i > 0 {
			sb.WriteString(" || ")
		}
		formatAlternative(sb, alternative)
	}

	if sb.Len() == 0 {
		sb.WriteByte('*')
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
