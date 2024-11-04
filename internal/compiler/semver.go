package compiler

import (
	"cmp"
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
var versionRegexp = makeRegexp("^(0|[1-9]\\d*)(?:\\.(0|[1-9]\\d*)(?:\\.(0|[1-9]\\d*)(?:-([a-zA-Z0-9-.]+))?(?:\\+([a-zA-Z0-9-.]+))?)?)?$")

// https://semver.org/#spec-item-9
// > A pre-release version MAY be denoted by appending a hyphen and a series of dot separated
// > identifiers immediately following the patch version. Identifiers MUST comprise only ASCII
// > alphanumerics and hyphen [0-9A-Za-z-]. Identifiers MUST NOT be empty. Numeric identifiers
// > MUST NOT include leading zeroes.
var prereleaseRegexp = makeRegexp("^(?:0|[1-9]\\d*|[a-zA-Z-][a-zA-Z0-9-]*)(?:\\.(?:0|[1-9]\\d*|[a-zA-Z-][a-zA-Z0-9-]*))*$")
var prereleasePartRegexp = makeRegexp("^(?:0|[1-9]\\d*|[a-zA-Z-][a-zA-Z0-9-]*)$")

// https://semver.org/#spec-item-10
// > Build metadata MAY be denoted by appending a plus sign and a series of dot separated
// > identifiers immediately following the patch or pre-release version. Identifiers MUST
// > comprise only ASCII alphanumerics and hyphen [0-9A-Za-z-]. Identifiers MUST NOT be empty.
var buildRegExp = makeRegexp("^[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$")
var buildPartRegExp = makeRegexp("^[a-zA-Z0-9-]+$")

// https://semver.org/#spec-item-9
// > Numeric identifiers MUST NOT include leading zeroes.
var numericIdentifierRegExp = makeRegexp("^(?:0|[1-9]\\d*)$")

type Version struct {
	// !!!
	// These should probably all be float64
	major      uint32
	minor      uint32
	patch      uint32
	prerelease []string
	build      []string
}

const ComparisonLessThan = -1
const ComparisonEqualTo = 0
const ComparisonGreaterThan = 1

func compareSemverVersion(a, b *Version) int {
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
		return ComparisonEqualTo
	case a == nil:
		return ComparisonLessThan
	case b == nil:
		return ComparisonGreaterThan
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
			return ComparisonEqualTo
		}
		return ComparisonGreaterThan
	} else if len(right) == 0 {
		return ComparisonLessThan
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
			return ComparisonLessThan
		}
		if !leftIsNumeric {
			return ComparisonGreaterThan
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

type SemverParseError struct {
	origInput string
}

func (e *SemverParseError) Error() string {
	return "Could not parse version string from " + e.origInput
}

func parseSemver(text string) (Version, error) {
	var result Version

	match := versionRegexp.FindStringSubmatch(text)
	if match == nil {
		return result, &SemverParseError{origInput: text}
	}
	match = match[1:]

	result.major = getUintComponent(match[0])
	result.minor = getUintComponent(match[1])
	result.patch = getUintComponent(match[2])

	prerelease := match[3]
	if prerelease != "" {
		result.prerelease = strings.Split(match[3], ".")
	}
	build := match[4]
	if build != "" {
		result.build = strings.Split(build, ".")
	}

	return result, nil
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
	alternatives []versionComparator
}

type versionComparator struct {
	operand  Version
	operator comparatorOperator
}

type comparatorOperator int

const (
	rangeLessThan comparatorOperator = iota
	rangeLessThanEqual
	rangeEqual
	rangeGreaterThanEqual
	rangeGreaterThan
)

// https://github.com/npm/node-semver#range-grammar
//
// range-set    ::= range ( logical-or range ) *
// range        ::= hyphen | simple ( ' ' simple ) * | ”
// logical-or   ::= ( ' ' ) * '||' ( ' ' ) *
var logicalOrRegExp = makeRegexp("\\|\\|")
var whitespaceRegExp = makeRegexp("\\s+")

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
var partialRegExp = makeRegexp("^([x*0]|[1-9]\\d*)(?:\\.([x*0]|[1-9]\\d*)(?:\\.([x*0]|[1-9]\\d*)(?:-([a-zA-Z0-9-.]+))?(?:\\+([a-zA-Z0-9-.]+))?)?)?$")
