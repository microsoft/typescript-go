package semver

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isAlphanumericOrHyphen(c byte) bool {
	return isDigit(c) || isLetter(c) || c == '-'
}

// isPreOrBuildChar returns true for characters valid in prerelease/build segments: [a-zA-Z0-9-.]
func isPreOrBuildChar(c byte) bool {
	return isAlphanumericOrHyphen(c) || c == '.'
}

// parseNr parses "0" or "[1-9][0-9]*" from s starting at pos.
func parseNr(s string, pos int) (string, int, bool) {
	if pos >= len(s) || !isDigit(s[pos]) {
		return "", pos, false
	}
	start := pos
	for pos < len(s) && isDigit(s[pos]) {
		pos++
	}
	nr := s[start:pos]
	if len(nr) > 1 && nr[0] == '0' {
		return "", start, false
	}
	return nr, pos, true
}

// https://semver.org/#spec-item-9
// > Numeric identifiers MUST NOT include leading zeroes.
func isNumericIdentifier(s string) bool {
	if s == "" {
		panic("isNumericIdentifier called with empty string")
	}
	_, end, ok := parseNr(s, 0)
	return ok && end == len(s)
}

// isValidPrereleasePart validates a single prerelease identifier.
// https://semver.org/#spec-item-9
// > A pre-release version MAY be denoted by appending a hyphen and a series of dot separated
// > identifiers immediately following the patch version. Identifiers MUST comprise only ASCII
// > alphanumerics and hyphen [0-9A-Za-z-]. Identifiers MUST NOT be empty. Numeric identifiers
// > MUST NOT include leading zeroes.
func isValidPrereleasePart(s string) bool {
	if s == "" {
		return false
	}
	if isNumericIdentifier(s) {
		return true
	}
	if !isLetter(s[0]) && s[0] != '-' {
		return false
	}
	for i := 1; i < len(s); i++ {
		if !isAlphanumericOrHyphen(s[i]) {
			panic("isValidPrereleasePart called with non-alphanumeric character: " + s)
		}
	}
	return true
}

func isValidPrerelease(s string) bool {
	if s == "" {
		panic("isValidPrerelease called with empty string")
	}
	for part := range strings.SplitSeq(s, ".") {
		if !isValidPrereleasePart(part) {
			return false
		}
	}
	return true
}

// https://semver.org/#spec-item-10
// > Build metadata MAY be denoted by appending a plus sign and a series of dot separated
// > identifiers immediately following the patch or pre-release version. Identifiers MUST
// > comprise only ASCII alphanumerics and hyphen [0-9A-Za-z-]. Identifiers MUST NOT be empty.
func isValidBuild(s string) bool {
	if s == "" {
		panic("isValidBuild called with empty string")
	}
	for part := range strings.SplitSeq(s, ".") {
		if part == "" {
			return false
		}
		for i := range len(part) {
			if !isAlphanumericOrHyphen(part[i]) {
				panic("isValidBuild called with non-alphanumeric character: " + s)
			}
		}
	}
	return true
}

type Version struct {
	major      uint32
	minor      uint32
	patch      uint32
	prerelease []string
	build      []string
}

var versionZero = Version{
	prerelease: []string{"0"},
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

const (
	comparisonLessThan    = -1
	comparisonEqualTo     = 0
	comparisonGreaterThan = 1
)

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
	compareResult := strings.Compare(left, right)
	if compareResult == 0 {
		return compareResult
	}

	leftIsNumeric := isNumericIdentifier(left)
	rightIsNumeric := isNumericIdentifier(right)

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
		leftAsNumber, leftErr := getUintComponent(left)
		rightAsNumber, rightErr := getUintComponent(right)
		if leftErr != nil || rightErr != nil {
			// This should only happen in the event of an overflow.
			// If so, use the lengths or fall back to string comparison.
			leftLen := len(left)
			rightLen := len(right)
			lenCompare := cmp.Compare(leftLen, rightLen)
			if lenCompare == 0 {
				return compareResult
			} else {
				return lenCompare
			}
		}
		return cmp.Compare(leftAsNumber, rightAsNumber)
	}

	// https://semver.org/#spec-item-11
	// > identifiers with letters or hyphens are compared lexically in ASCII sort order.
	return compareResult
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
	return fmt.Sprintf("Could not parse version string from %q", e.origInput)
}

func TryParseVersion(text string) (Version, error) {
	var result Version
	pos := 0

	// Parse major: 0 | [1-9][0-9]*
	majorStr, newPos, ok := parseNr(text, pos)
	if !ok {
		return result, &SemverParseError{origInput: text}
	}
	pos = newPos

	var minorStr, patchStr, prereleaseStr, buildStr string

	// Optional .minor
	if pos < len(text) && text[pos] == '.' {
		pos++
		minorStr, newPos, ok = parseNr(text, pos)
		if !ok {
			return result, &SemverParseError{origInput: text}
		}
		pos = newPos

		// Optional .patch
		if pos < len(text) && text[pos] == '.' {
			pos++
			patchStr, newPos, ok = parseNr(text, pos)
			if !ok {
				return result, &SemverParseError{origInput: text}
			}
			pos = newPos

			// Optional -prerelease
			if pos < len(text) && text[pos] == '-' {
				pos++
				start := pos
				for pos < len(text) && isPreOrBuildChar(text[pos]) {
					pos++
				}
				if pos == start {
					return result, &SemverParseError{origInput: text}
				}
				prereleaseStr = text[start:pos]
			}

			// Optional +build
			if pos < len(text) && text[pos] == '+' {
				pos++
				start := pos
				for pos < len(text) && isPreOrBuildChar(text[pos]) {
					pos++
				}
				if pos == start {
					return result, &SemverParseError{origInput: text}
				}
				buildStr = text[start:pos]
			}
		}
	}

	if pos != len(text) {
		return result, &SemverParseError{origInput: text}
	}

	var err error

	result.major, err = getUintComponent(majorStr)
	if err != nil {
		return result, err
	}

	if minorStr != "" {
		result.minor, err = getUintComponent(minorStr)
		if err != nil {
			return result, err
		}
	}

	if patchStr != "" {
		result.patch, err = getUintComponent(patchStr)
		if err != nil {
			return result, err
		}
	}

	if prereleaseStr != "" {
		if !isValidPrerelease(prereleaseStr) {
			return result, &SemverParseError{origInput: text}
		}
		result.prerelease = strings.Split(prereleaseStr, ".")
	}

	if buildStr != "" {
		if !isValidBuild(buildStr) {
			return result, &SemverParseError{origInput: text}
		}
		result.build = strings.Split(buildStr, ".")
	}

	return result, nil
}

func MustParse(text string) Version {
	v, err := TryParseVersion(text)
	if err != nil {
		panic(err)
	}
	return v
}

func getUintComponent(text string) (uint32, error) {
	r, err := strconv.ParseUint(text, 10, 32)
	return uint32(r), err
}
