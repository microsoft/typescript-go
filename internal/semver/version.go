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
var versionRegexp = regexp.MustCompile(`^(0|[1-9]\d*)(?:\.(0|[1-9]\d*)(?:\.(0|[1-9]\d*)(?:-([a-zA-Z0-9-.]+))?(?:\+([a-zA-Z0-9-.]+))?)?)?$`)

// https://semver.org/#spec-item-9
// > A pre-release version MAY be denoted by appending a hyphen and a series of dot separated
// > identifiers immediately following the patch version. Identifiers MUST comprise only ASCII
// > alphanumerics and hyphen [0-9A-Za-z-]. Identifiers MUST NOT be empty. Numeric identifiers
// > MUST NOT include leading zeroes.
var prereleaseRegexp = regexp.MustCompile(`^(?:0|[1-9]\d*|[a-zA-Z-][a-zA-Z0-9-]*)(?:\.(?:0|[1-9]\d*|[a-zA-Z-][a-zA-Z0-9-]*))*$`)
var prereleasePartRegexp = regexp.MustCompile(`^(?:0|[1-9]\d*|[a-zA-Z-][a-zA-Z0-9-]*)$`)

// https://semver.org/#spec-item-10
// > Build metadata MAY be denoted by appending a plus sign and a series of dot separated
// > identifiers immediately following the patch or pre-release version. Identifiers MUST
// > comprise only ASCII alphanumerics and hyphen [0-9A-Za-z-]. Identifiers MUST NOT be empty.
var buildRegExp = regexp.MustCompile(`^[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$`)
var buildPartRegExp = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

// https://semver.org/#spec-item-9
// > Numeric identifiers MUST NOT include leading zeroes.
var numericIdentifierRegExp = regexp.MustCompile(`^(?:0|[1-9]\d*)$`)

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

func (v *Version) incrementMajor() Version {
	return Version{
		major: v.major + 1,
	}
}

func (v *Version) incrementMinor() Version {
	return Version{
		major: v.major,
		minor: v.minor + 1,
	}
}

func (v *Version) incrementPatch() Version {
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
