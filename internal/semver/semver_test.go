package semver

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

// Versions

func TestTryParseSemver(t *testing.T) {
	var tests = []struct {
		in  string
		out Version
	}{
		{"1.2.3-pre.4+build.5", Version{major: 1, minor: 2, patch: 3, prerelease: []string{"pre", "4"}, build: []string{"build", "5"}}},
		{"1.2.3-pre.4", Version{major: 1, minor: 2, patch: 3, prerelease: []string{"pre", "4"}}},
		{"1.2.3+build.4", Version{major: 1, minor: 2, patch: 3, build: []string{"build", "4"}}},
		{"1.2.3", Version{major: 1, minor: 2, patch: 3}},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			v, err := TryParseVersion(test.in)
			assert.NilError(t, err)
			assertVersion(t, v, test.out)
		})
	}
}

func TestVersionString(t *testing.T) {
	var tests = []struct {
		in  Version
		out string
	}{
		{Version{major: 1, minor: 2, patch: 3, prerelease: []string{"pre", "4"}, build: []string{"build", "5"}}, "1.2.3-pre.4+build.5"},
		{Version{major: 1, minor: 2, patch: 3, prerelease: []string{"pre", "4"}, build: []string{"build"}}, "1.2.3-pre.4+build"},
		{Version{major: 1, minor: 2, patch: 3, build: []string{"build"}}, "1.2.3+build"},
		{Version{major: 1, minor: 2, patch: 3, prerelease: []string{"pre", "4"}}, "1.2.3-pre.4"},
		{Version{major: 1, minor: 2, patch: 3, build: []string{"build", "4"}}, "1.2.3+build.4"},
		{Version{major: 1, minor: 2, patch: 3}, "1.2.3"},
	}

	for _, test := range tests {
		t.Run(test.out, func(t *testing.T) {
			assert.Equal(t, test.in.String(), test.out)
		})
	}
}

func TestVersionCompare(t *testing.T) {
	var tests = []struct {
		v1, v2 string
		want   int
	}{
		// https://semver.org/#spec-item-11
		// > Precedence is determined by the first difference when comparing each of these
		// > identifiers from left to right as follows: Major, minor, and patch versions are
		// > always compared numerically.
		{"1.0.0", "2.0.0", comparisonLessThan},
		{"1.0.0", "1.1.0", comparisonLessThan},
		{"1.0.0", "1.0.1", comparisonLessThan},
		{"2.0.0", "1.0.0", comparisonGreaterThan},
		{"1.1.0", "1.0.0", comparisonGreaterThan},
		{"1.0.1", "1.0.0", comparisonGreaterThan},
		{"1.0.0", "1.0.0", comparisonEqualTo},

		// https://semver.org/#spec-item-11
		// > When major, minor, and patch are equal, a pre-release version has lower
		// > precedence than a normal version.
		{"1.0.0", "1.0.0-pre", comparisonGreaterThan},
		{"1.0.1-pre", "1.0.0", comparisonGreaterThan},
		{"1.0.0-pre", "1.0.0", comparisonLessThan},

		// https://semver.org/#spec-item-11
		// > identifiers consisting of only digits are compared numerically
		{"1.0.0-0", "1.0.0-1", comparisonLessThan},
		{"1.0.0-1", "1.0.0-0", comparisonGreaterThan},
		{"1.0.0-2", "1.0.0-10", comparisonLessThan},
		{"1.0.0-10", "1.0.0-2", comparisonGreaterThan},
		{"1.0.0-0", "1.0.0-0", comparisonEqualTo},

		// https://semver.org/#spec-item-11
		// > identifiers with letters or hyphens are compared lexically in ASCII sort order.
		{"1.0.0-a", "1.0.0-b", comparisonLessThan},
		{"1.0.0-a-2", "1.0.0-a-10", comparisonGreaterThan},
		{"1.0.0-b", "1.0.0-a", comparisonGreaterThan},
		{"1.0.0-a", "1.0.0-a", comparisonEqualTo},
		{"1.0.0-A", "1.0.0-a", comparisonLessThan},

		// https://semver.org/#spec-item-11
		// > Numeric identifiers always have lower precedence than non-numeric identifiers.
		{"1.0.0-0", "1.0.0-alpha", comparisonLessThan},
		{"1.0.0-alpha", "1.0.0-0", comparisonGreaterThan},
		{"1.0.0-0", "1.0.0-0", comparisonEqualTo},
		{"1.0.0-alpha", "1.0.0-alpha", comparisonEqualTo},

		// https://semver.org/#spec-item-11
		// > A larger set of pre-release fields has a higher precedence than a smaller set, if all
		// > of the preceding identifiers are equal.
		{"1.0.0-alpha", "1.0.0-alpha.0", comparisonLessThan},
		{"1.0.0-alpha.0", "1.0.0-alpha", comparisonGreaterThan},

		// https://semver.org/#spec-item-11
		// > Precedence for two pre-release versions with the same major, minor, and patch version
		// > MUST be determined by comparing each dot separated identifier from left to right until
		// > a difference is found [...]
		{"1.0.0-a.0.b.1", "1.0.0-a.0.b.2", comparisonLessThan},
		{"1.0.0-a.0.b.1", "1.0.0-b.0.a.1", comparisonLessThan},
		{"1.0.0-a.0.b.2", "1.0.0-a.0.b.1", comparisonGreaterThan},
		{"1.0.0-b.0.a.1", "1.0.0-a.0.b.1", comparisonGreaterThan},

		// https://semver.org/#spec-item-11
		// > Build metadata does not figure into precedence
		{"1.0.0+build", "1.0.0", comparisonEqualTo},
		{"1.0.0+build.stuff", "1.0.0", comparisonEqualTo},
		{"1.0.0", "1.0.0+build", comparisonEqualTo},
		{"1.0.0+build", "1.0.0+stuff", comparisonEqualTo},
	}

	for _, test := range tests {
		t.Run(test.v1+" <=> "+test.v2, func(t *testing.T) {
			v1, err1 := TryParseVersion(test.v1)
			assert.NilError(t, err1, test.v1)
			v2, err2 := TryParseVersion(test.v2)
			assert.NilError(t, err2, test.v2)
			assert.Equal(t, v1.Compare(&v2), test.want)
		})
	}
}

func assertVersion(t *testing.T, a, b Version) {
	assert.Equal(t, a.major, b.major)
	assert.Equal(t, a.minor, b.minor)
	assert.Equal(t, a.patch, b.patch)
	assert.DeepEqual(t, a.prerelease, b.prerelease)
	assert.DeepEqual(t, a.build, b.build)
}

// Version Ranges

func TestWildcardsHaveSameString(t *testing.T) {
	majorWildcardStrings := []string{
		"",
		"*",
		"*.*",
		"*.*.*",
		"x",
		"x.x",
		"x.x.x",
		"X",
		"X.X",
		"X.X.X",
	}

	minorWildcardStrings := []string{
		"1",
		"1.*",
		"1.*.*",
		"1.x",
		"1.x.x",
		"1.X",
		"1.X.X",
	}

	patchWildcardStrings := []string{
		"1.2",
		"1.2.*",
		"1.2.x",
		"1.2.X",
	}

	assertAllVersionRangesHaveIdenticalStrings(t, "majorWildcardStrings", majorWildcardStrings)
	assertAllVersionRangesHaveIdenticalStrings(t, "minorWildcardStrings", minorWildcardStrings)
	assertAllVersionRangesHaveIdenticalStrings(t, "patchWildcardStrings", patchWildcardStrings)
}

func assertAllVersionRangesHaveIdenticalStrings(t *testing.T, name string, strs []string) {
	t.Run(name, func(t *testing.T) {
		for _, s1 := range strs {
			for _, s2 := range strs {
				t.Run(s1+" == "+s2, func(t *testing.T) {
					v1, ok := TryParseVersionRange(s1)
					assert.Assert(t, ok)
					v2, ok := TryParseVersionRange(s2)
					assert.Assert(t, ok)
					assert.DeepEqual(t, v1.String(), v2.String())
				})
			}
		}
	})
}

type testGoodBad struct {
	good []string
	bad  []string
}

func TestVersionRanges(t *testing.T) {
	assertRangesGoodBad(t, "1", testGoodBad{
		good: []string{"1.0.0", "1.9.9", "1.0.0-pre", "1.0.0+build"},
		bad:  []string{"0.0.0", "2.0.0", "0.0.0-pre", "0.0.0+build"},
	})
	assertRangesGoodBad(t, "1.2", testGoodBad{
		good: []string{"1.2.0", "1.2.9", "1.2.0-pre", "1.2.0+build"},
		bad:  []string{"1.1.0", "1.3.0", "1.1.0-pre", "1.1.0+build"},
	})

	assertRangesGoodBad(t, "1.2.3", testGoodBad{
		good: []string{"1.2.3", "1.2.3+build"},
		bad:  []string{"1.2.2", "1.2.4", "1.2.2-pre", "1.2.2+build", "1.2.3-pre"},
	})

	assertRangesGoodBad(t, "1.2.3-pre", testGoodBad{
		good: []string{"1.2.3-pre", "1.2.3-pre+build.stuff"},
		bad:  []string{"1.2.3", "1.2.3-pre.0", "1.2.3-pre.9", "1.2.3-pre.0+build", "1.2.3-pre.9+build", "1.2.3+build", "1.2.4"},
	})

	assertRangesGoodBad(t, "<3.8.0", testGoodBad{
		good: []string{"3.6", "3.7"},
		bad:  []string{"3.8", "3.9", "4.0"},
	})

	assertRangesGoodBad(t, "<=3.8.0", testGoodBad{
		good: []string{"3.6", "3.7", "3.8"},
		bad:  []string{"3.9", "4.0"},
	})
	assertRangesGoodBad(t, ">3.8.0", testGoodBad{
		good: []string{"3.9", "4.0"},
		bad:  []string{"3.6", "3.7", "3.8"},
	})
	assertRangesGoodBad(t, ">=3.8.0", testGoodBad{
		good: []string{"3.8", "3.9", "4.0"},
		bad:  []string{"3.6", "3.7"},
	})

	assertRangesGoodBad(t, "<3.8.0-0", testGoodBad{
		good: []string{"3.6", "3.7"},
		bad:  []string{"3.8", "3.9", "4.0"},
	})

	assertRangesGoodBad(t, "<=3.8.0-0", testGoodBad{
		good: []string{"3.6", "3.7"},
		bad:  []string{"3.8", "3.9", "4.0"},
	})
}

func TestComparatorsOfVersionRanges(t *testing.T) {
	comparatorsTests := []testForRangeOnVersion{
		// empty (matches everything)
		{"", "2.0.0", true},
		{"", "2.0.0-0", true},
		{"", "1.1.0", true},
		{"", "1.1.0-0", true},
		{"", "1.0.1", true},
		{"", "1.0.1-0", true},
		{"", "1.0.0", true},
		{"", "1.0.0-0", true},
		{"", "0.0.0", true},
		{"", "0.0.0-0", true},

		// wildcard major (matches everything)
		{"*", "2.0.0", true},
		{"*", "2.0.0-0", true},
		{"*", "1.1.0", true},
		{"*", "1.1.0-0", true},
		{"*", "1.0.1", true},
		{"*", "1.0.1-0", true},
		{"*", "1.0.0", true},
		{"*", "1.0.0-0", true},
		{"*", "0.0.0", true},
		{"*", "0.0.0-0", true},

		// wildcard minor
		{"1", "2.0.0", false},
		{"1", "2.0.0-0", false},
		{"1", "1.1.0", true},
		{"1", "1.1.0-0", true},
		{"1", "1.0.1", true},
		{"1", "1.0.1-0", true},
		{"1", "1.0.0", true},
		{"1", "1.0.0-0", true},
		{"1", "0.0.0", false},
		{"1", "0.0.0-0", false},

		// wildcard patch
		{"1.1", "2.0.0", false},
		{"1.1", "2.0.0-0", false},
		{"1.1", "1.1.0", true},
		{"1.1", "1.1.0-0", true},
		{"1.1", "1.0.1", false},
		{"1.1", "1.0.1-0", false},
		{"1.1", "1.0.0", false},
		{"1.1", "1.0.0-0", false},
		{"1.1", "0.0.0", false},
		{"1.1", "0.0.0-0", false},
		{"1.0", "2.0.0", false},
		{"1.0", "2.0.0-0", false},
		{"1.0", "1.1.0", false},
		{"1.0", "1.1.0-0", false},
		{"1.0", "1.0.1", true},
		{"1.0", "1.0.1-0", true},
		{"1.0", "1.0.0", true},
		{"1.0", "1.0.0-0", true},
		{"1.0", "0.0.0", false},
		{"1.0", "0.0.0-0", false},

		// exact
		{"1.1.0", "2.0.0", false},
		{"1.1.0", "2.0.0-0", false},
		{"1.1.0", "1.1.0", true},
		{"1.1.0", "1.1.0-0", false},
		{"1.1.0", "1.0.1", false},
		{"1.1.0", "1.0.1-0", false},
		{"1.1.0", "1.0.0-0", false},
		{"1.1.0", "1.0.0", false},
		{"1.1.0", "0.0.0", false},
		{"1.1.0", "0.0.0-0", false},
		{"1.1.0-0", "2.0.0", false},
		{"1.1.0-0", "2.0.0-0", false},
		{"1.1.0-0", "1.1.0", false},
		{"1.1.0-0", "1.1.0-0", true},
		{"1.1.0-0", "1.0.1", false},
		{"1.1.0-0", "1.0.1-0", false},
		{"1.1.0-0", "1.0.0-0", false},
		{"1.1.0-0", "1.0.0", false},
		{"1.1.0-0", "0.0.0", false},
		{"1.1.0-0", "0.0.0-0", false},
		{"1.0.1", "2.0.0", false},
		{"1.0.1", "2.0.0-0", false},
		{"1.0.1", "1.1.0", false},
		{"1.0.1", "1.1.0-0", false},
		{"1.0.1", "1.0.1", true},
		{"1.0.1", "1.0.1-0", false},
		{"1.0.1", "1.0.0-0", false},
		{"1.0.1", "1.0.0", false},
		{"1.0.1", "0.0.0", false},
		{"1.0.1", "0.0.0-0", false},
		{"1.0.1-0", "2.0.0", false},
		{"1.0.1-0", "2.0.0-0", false},
		{"1.0.1-0", "1.1.0", false},
		{"1.0.1-0", "1.1.0-0", false},
		{"1.0.1-0", "1.0.1", false},
		{"1.0.1-0", "1.0.1-0", true},
		{"1.0.1-0", "1.0.0-0", false},
		{"1.0.1-0", "1.0.0", false},
		{"1.0.1-0", "0.0.0", false},
		{"1.0.1-0", "0.0.0-0", false},
		{"1.0.0", "2.0.0", false},
		{"1.0.0", "2.0.0-0", false},
		{"1.0.0", "1.1.0", false},
		{"1.0.0", "1.1.0-0", false},
		{"1.0.0", "1.0.1", false},
		{"1.0.0", "1.0.1-0", false},
		{"1.0.0", "1.0.0-0", false},
		{"1.0.0", "1.0.0", true},
		{"1.0.0", "0.0.0", false},
		{"1.0.0", "0.0.0-0", false},
		{"1.0.0-0", "2.0.0", false},
		{"1.0.0-0", "2.0.0-0", false},
		{"1.0.0-0", "1.1.0", false},
		{"1.0.0-0", "1.1.0-0", false},
		{"1.0.0-0", "1.0.1", false},
		{"1.0.0-0", "1.0.1-0", false},
		{"1.0.0-0", "1.0.0", false},
		{"1.0.0-0", "1.0.0-0", true},

		// = wildcard major (matches everything)
		{"=*", "2.0.0", true},
		{"=*", "2.0.0-0", true},
		{"=*", "1.1.0", true},
		{"=*", "1.1.0-0", true},
		{"=*", "1.0.1", true},
		{"=*", "1.0.1-0", true},
		{"=*", "1.0.0", true},
		{"=*", "1.0.0-0", true},
		{"=*", "0.0.0", true},
		{"=*", "0.0.0-0", true},

		// = wildcard minor
		{"=1", "2.0.0", false},
		{"=1", "2.0.0-0", false},
		{"=1", "1.1.0", true},
		{"=1", "1.1.0-0", true},
		{"=1", "1.0.1", true},
		{"=1", "1.0.1-0", true},
		{"=1", "1.0.0", true},
		{"=1", "1.0.0-0", true},
		{"=1", "0.0.0", false},
		{"=1", "0.0.0-0", false},

		// = wildcard patch
		{"=1.1", "2.0.0", false},
		{"=1.1", "2.0.0-0", false},
		{"=1.1", "1.1.0", true},
		{"=1.1", "1.1.0-0", true},
		{"=1.1", "1.0.1", false},
		{"=1.1", "1.0.1-0", false},
		{"=1.1", "1.0.0", false},
		{"=1.1", "1.0.0-0", false},
		{"=1.1", "0.0.0", false},
		{"=1.1", "0.0.0-0", false},
		{"=1.0", "2.0.0", false},
		{"=1.0", "2.0.0-0", false},
		{"=1.0", "1.1.0", false},
		{"=1.0", "1.1.0-0", false},
		{"=1.0", "1.0.1", true},
		{"=1.0", "1.0.1-0", true},
		{"=1.0", "1.0.0", true},
		{"=1.0", "1.0.0-0", true},
		{"=1.0", "0.0.0", false},
		{"=1.0", "0.0.0-0", false},

		// = exact
		{"=1.1.0", "2.0.0", false},
		{"=1.1.0", "2.0.0-0", false},
		{"=1.1.0", "1.1.0", true},
		{"=1.1.0", "1.1.0-0", false},
		{"=1.1.0", "1.0.1", false},
		{"=1.1.0", "1.0.1-0", false},
		{"=1.1.0", "1.0.0-0", false},
		{"=1.1.0", "1.0.0", false},
		{"=1.1.0", "0.0.0", false},
		{"=1.1.0", "0.0.0-0", false},
		{"=1.1.0-0", "2.0.0", false},
		{"=1.1.0-0", "2.0.0-0", false},
		{"=1.1.0-0", "1.1.0", false},
		{"=1.1.0-0", "1.1.0-0", true},
		{"=1.1.0-0", "1.0.1", false},
		{"=1.1.0-0", "1.0.1-0", false},
		{"=1.1.0-0", "1.0.0-0", false},
		{"=1.1.0-0", "1.0.0", false},
		{"=1.1.0-0", "0.0.0", false},
		{"=1.1.0-0", "0.0.0-0", false},
		{"=1.0.1", "2.0.0", false},
		{"=1.0.1", "2.0.0-0", false},
		{"=1.0.1", "1.1.0", false},
		{"=1.0.1", "1.1.0-0", false},
		{"=1.0.1", "1.0.1", true},
		{"=1.0.1", "1.0.1-0", false},
		{"=1.0.1", "1.0.0-0", false},
		{"=1.0.1", "1.0.0", false},
		{"=1.0.1", "0.0.0", false},
		{"=1.0.1", "0.0.0-0", false},
		{"=1.0.1-0", "2.0.0", false},
		{"=1.0.1-0", "2.0.0-0", false},
		{"=1.0.1-0", "1.1.0", false},
		{"=1.0.1-0", "1.1.0-0", false},
		{"=1.0.1-0", "1.0.1", false},
		{"=1.0.1-0", "1.0.1-0", true},
		{"=1.0.1-0", "1.0.0-0", false},
		{"=1.0.1-0", "1.0.0", false},
		{"=1.0.1-0", "0.0.0", false},
		{"=1.0.1-0", "0.0.0-0", false},
		{"=1.0.0", "2.0.0", false},
		{"=1.0.0", "2.0.0-0", false},
		{"=1.0.0", "1.1.0", false},
		{"=1.0.0", "1.1.0-0", false},
		{"=1.0.0", "1.0.1", false},
		{"=1.0.0", "1.0.1-0", false},
		{"=1.0.0", "1.0.0-0", false},
		{"=1.0.0", "1.0.0", true},
		{"=1.0.0", "0.0.0", false},
		{"=1.0.0", "0.0.0-0", false},
		{"=1.0.0-0", "2.0.0", false},
		{"=1.0.0-0", "2.0.0-0", false},
		{"=1.0.0-0", "1.1.0", false},
		{"=1.0.0-0", "1.1.0-0", false},
		{"=1.0.0-0", "1.0.1", false},
		{"=1.0.0-0", "1.0.1-0", false},
		{"=1.0.0-0", "1.0.0", false},
		{"=1.0.0-0", "1.0.0-0", true},

		// > wildcard major (matches nothing)
		{">*", "2.0.0", false},
		{">*", "2.0.0-0", false},
		{">*", "1.1.0", false},
		{">*", "1.1.0-0", false},
		{">*", "1.0.1", false},
		{">*", "1.0.1-0", false},
		{">*", "1.0.0", false},
		{">*", "1.0.0-0", false},
		{">*", "0.0.0", false},
		{">*", "0.0.0-0", false},

		// > wildcard minor
		{">1", "2.0.0", true},
		{">1", "2.0.0-0", true},
		{">1", "1.1.0", false},
		{">1", "1.1.0-0", false},
		{">1", "1.0.1", false},
		{">1", "1.0.1-0", false},
		{">1", "1.0.0", false},
		{">1", "1.0.0-0", false},
		{">1", "0.0.0", false},
		{">1", "0.0.0-0", false},

		// > wildcard patch
		{">1.1", "2.0.0", true},
		{">1.1", "2.0.0-0", true},
		{">1.1", "1.1.0", false},
		{">1.1", "1.1.0-0", false},
		{">1.1", "1.0.1", false},
		{">1.1", "1.0.1-0", false},
		{">1.1", "1.0.0", false},
		{">1.1", "1.0.0-0", false},
		{">1.1", "0.0.0", false},
		{">1.1", "0.0.0-0", false},
		{">1.0", "2.0.0", true},
		{">1.0", "2.0.0-0", true},
		{">1.0", "1.1.0", true},
		{">1.0", "1.1.0-0", true},
		{">1.0", "1.0.1", false},
		{">1.0", "1.0.1-0", false},
		{">1.0", "1.0.0", false},
		{">1.0", "1.0.0-0", false},
		{">1.0", "0.0.0", false},
		{">1.0", "0.0.0-0", false},

		// > exact
		{">1.1.0", "2.0.0", true},
		{">1.1.0", "2.0.0-0", true},
		{">1.1.0", "1.1.0", false},
		{">1.1.0", "1.1.0-0", false},
		{">1.1.0", "1.0.1", false},
		{">1.1.0", "1.0.1-0", false},
		{">1.1.0", "1.0.0", false},
		{">1.1.0", "1.0.0-0", false},
		{">1.1.0", "0.0.0", false},
		{">1.1.0", "0.0.0-0", false},
		{">1.1.0-0", "2.0.0", true},
		{">1.1.0-0", "2.0.0-0", true},
		{">1.1.0-0", "1.1.0", true},
		{">1.1.0-0", "1.1.0-0", false},
		{">1.1.0-0", "1.0.1", false},
		{">1.1.0-0", "1.0.1-0", false},
		{">1.1.0-0", "1.0.0", false},
		{">1.1.0-0", "1.0.0-0", false},
		{">1.1.0-0", "0.0.0", false},
		{">1.1.0-0", "0.0.0-0", false},
		{">1.0.1", "2.0.0", true},
		{">1.0.1", "2.0.0-0", true},
		{">1.0.1", "1.1.0", true},
		{">1.0.1", "1.1.0-0", true},
		{">1.0.1", "1.0.1", false},
		{">1.0.1", "1.0.1-0", false},
		{">1.0.1", "1.0.0", false},
		{">1.0.1", "1.0.0-0", false},
		{">1.0.1", "0.0.0", false},
		{">1.0.1", "0.0.0-0", false},
		{">1.0.1-0", "2.0.0", true},
		{">1.0.1-0", "2.0.0-0", true},
		{">1.0.1-0", "1.1.0", true},
		{">1.0.1-0", "1.1.0-0", true},
		{">1.0.1-0", "1.0.1", true},
		{">1.0.1-0", "1.0.1-0", false},
		{">1.0.1-0", "1.0.0", false},
		{">1.0.1-0", "1.0.0-0", false},
		{">1.0.1-0", "0.0.0", false},
		{">1.0.1-0", "0.0.0-0", false},
		{">1.0.0", "2.0.0", true},
		{">1.0.0", "2.0.0-0", true},
		{">1.0.0", "1.1.0", true},
		{">1.0.0", "1.1.0-0", true},
		{">1.0.0", "1.0.1", true},
		{">1.0.0", "1.0.1-0", true},
		{">1.0.0", "1.0.0", false},
		{">1.0.0", "1.0.0-0", false},
		{">1.0.0", "0.0.0", false},
		{">1.0.0", "0.0.0-0", false},
		{">1.0.0-0", "2.0.0", true},
		{">1.0.0-0", "2.0.0-0", true},
		{">1.0.0-0", "1.1.0", true},
		{">1.0.0-0", "1.1.0-0", true},
		{">1.0.0-0", "1.0.1", true},
		{">1.0.0-0", "1.0.1-0", true},
		{">1.0.0-0", "1.0.0", true},
		{">1.0.0-0", "1.0.0-0", false},
		{">1.0.0-0", "0.0.0", false},
		{">1.0.0-0", "0.0.0-0", false},

		// >= wildcard major (matches everything)
		{">=*", "2.0.0", true},
		{">=*", "2.0.0-0", true},
		{">=*", "1.1.0", true},
		{">=*", "1.1.0-0", true},
		{">=*", "1.0.1", true},
		{">=*", "1.0.1-0", true},
		{">=*", "1.0.0", true},
		{">=*", "1.0.0-0", true},
		{">=*", "0.0.0", true},
		{">=*", "0.0.0-0", true},

		// >= wildcard minor
		{">=1", "2.0.0", true},
		{">=1", "2.0.0-0", true},
		{">=1", "1.1.0", true},
		{">=1", "1.1.0-0", true},
		{">=1", "1.0.1", true},
		{">=1", "1.0.1-0", true},
		{">=1", "1.0.0", true},
		{">=1", "1.0.0-0", true},
		{">=1", "0.0.0", false},
		{">=1", "0.0.0-0", false},

		// >= wildcard patch
		{">=1.1", "2.0.0", true},
		{">=1.1", "2.0.0-0", true},
		{">=1.1", "1.1.0", true},
		{">=1.1", "1.1.0-0", true},
		{">=1.1", "1.0.1", false},
		{">=1.1", "1.0.1-0", false},
		{">=1.1", "1.0.0", false},
		{">=1.1", "1.0.0-0", false},
		{">=1.1", "0.0.0", false},
		{">=1.1", "0.0.0-0", false},
		{">=1.0", "2.0.0", true},
		{">=1.0", "2.0.0-0", true},
		{">=1.0", "1.1.0", true},
		{">=1.0", "1.1.0-0", true},
		{">=1.0", "1.0.1", true},
		{">=1.0", "1.0.1-0", true},
		{">=1.0", "1.0.0", true},
		{">=1.0", "1.0.0-0", true},
		{">=1.0", "0.0.0", false},
		{">=1.0", "0.0.0-0", false},

		// >= exact
		{">=1.1.0", "2.0.0", true},
		{">=1.1.0", "2.0.0-0", true},
		{">=1.1.0", "1.1.0", true},
		{">=1.1.0", "1.1.0-0", false},
		{">=1.1.0", "1.0.1", false},
		{">=1.1.0", "1.0.1-0", false},
		{">=1.1.0", "1.0.0", false},
		{">=1.1.0", "1.0.0-0", false},
		{">=1.1.0", "0.0.0", false},
		{">=1.1.0", "0.0.0-0", false},
		{">=1.1.0-0", "2.0.0", true},
		{">=1.1.0-0", "2.0.0-0", true},
		{">=1.1.0-0", "1.1.0", true},
		{">=1.1.0-0", "1.1.0-0", true},
		{">=1.1.0-0", "1.0.1", false},
		{">=1.1.0-0", "1.0.1-0", false},
		{">=1.1.0-0", "1.0.0", false},
		{">=1.1.0-0", "1.0.0-0", false},
		{">=1.1.0-0", "0.0.0", false},
		{">=1.1.0-0", "0.0.0-0", false},
		{">=1.0.1", "2.0.0", true},
		{">=1.0.1", "2.0.0-0", true},
		{">=1.0.1", "1.1.0", true},
		{">=1.0.1", "1.1.0-0", true},
		{">=1.0.1", "1.0.1", true},
		{">=1.0.1", "1.0.1-0", false},
		{">=1.0.1", "1.0.0", false},
		{">=1.0.1", "1.0.0-0", false},
		{">=1.0.1", "0.0.0", false},
		{">=1.0.1", "0.0.0-0", false},
		{">=1.0.1-0", "2.0.0", true},
		{">=1.0.1-0", "2.0.0-0", true},
		{">=1.0.1-0", "1.1.0", true},
		{">=1.0.1-0", "1.1.0-0", true},
		{">=1.0.1-0", "1.0.1", true},
		{">=1.0.1-0", "1.0.1-0", true},
		{">=1.0.1-0", "1.0.0", false},
		{">=1.0.1-0", "1.0.0-0", false},
		{">=1.0.1-0", "0.0.0", false},
		{">=1.0.1-0", "0.0.0-0", false},
		{">=1.0.0", "2.0.0", true},
		{">=1.0.0", "2.0.0-0", true},
		{">=1.0.0", "1.1.0", true},
		{">=1.0.0", "1.1.0-0", true},
		{">=1.0.0", "1.0.1", true},
		{">=1.0.0", "1.0.1-0", true},
		{">=1.0.0", "1.0.0", true},
		{">=1.0.0", "1.0.0-0", false},
		{">=1.0.0", "0.0.0", false},
		{">=1.0.0", "0.0.0-0", false},
		{">=1.0.0-0", "2.0.0", true},
		{">=1.0.0-0", "2.0.0-0", true},
		{">=1.0.0-0", "1.1.0", true},
		{">=1.0.0-0", "1.1.0-0", true},
		{">=1.0.0-0", "1.0.1", true},
		{">=1.0.0-0", "1.0.1-0", true},
		{">=1.0.0-0", "1.0.0", true},
		{">=1.0.0-0", "1.0.0-0", true},
		{">=1.0.0-0", "0.0.0", false},
		{">=1.0.0-0", "0.0.0-0", false},

		// < wildcard major (matches nothing)
		{"<*", "2.0.0", false},
		{"<*", "2.0.0-0", false},
		{"<*", "1.1.0", false},
		{"<*", "1.1.0-0", false},
		{"<*", "1.0.1", false},
		{"<*", "1.0.1-0", false},
		{"<*", "1.0.0", false},
		{"<*", "1.0.0-0", false},
		{"<*", "0.0.0", false},
		{"<*", "0.0.0-0", false},

		// < wildcard minor
		{"<1", "2.0.0", false},
		{"<1", "2.0.0-0", false},
		{"<1", "1.1.0", false},
		{"<1", "1.1.0-0", false},
		{"<1", "1.0.1", false},
		{"<1", "1.0.1-0", false},
		{"<1", "1.0.0", false},
		{"<1", "1.0.0-0", false},
		{"<1", "0.0.0", true},
		{"<1", "0.0.0-0", true},

		// < wildcard patch
		{"<1.1", "2.0.0", false},
		{"<1.1", "2.0.0-0", false},
		{"<1.1", "1.1.0", false},
		{"<1.1", "1.1.0-0", false},
		{"<1.1", "1.0.1", true},
		{"<1.1", "1.0.1-0", true},
		{"<1.1", "1.0.0", true},
		{"<1.1", "1.0.0-0", true},
		{"<1.1", "0.0.0", true},
		{"<1.1", "0.0.0-0", true},
		{"<1.0", "2.0.0", false},
		{"<1.0", "2.0.0-0", false},
		{"<1.0", "1.1.0", false},
		{"<1.0", "1.1.0-0", false},
		{"<1.0", "1.0.1", false},
		{"<1.0", "1.0.1-0", false},
		{"<1.0", "1.0.0", false},
		{"<1.0", "1.0.0-0", false},
		{"<1.0", "0.0.0", true},
		{"<1.0", "0.0.0-0", true},

		// < exact
		{"<1.1.0", "2.0.0", false},
		{"<1.1.0", "2.0.0-0", false},
		{"<1.1.0", "1.1.0", false},
		{"<1.1.0", "1.1.0-0", true},
		{"<1.1.0", "1.0.1", true},
		{"<1.1.0", "1.0.1-0", true},
		{"<1.1.0", "1.0.0", true},
		{"<1.1.0", "1.0.0-0", true},
		{"<1.1.0", "0.0.0", true},
		{"<1.1.0", "0.0.0-0", true},
		{"<1.1.0-0", "2.0.0", false},
		{"<1.1.0-0", "2.0.0-0", false},
		{"<1.1.0-0", "1.1.0", false},
		{"<1.1.0-0", "1.1.0-0", false},
		{"<1.1.0-0", "1.0.1", true},
		{"<1.1.0-0", "1.0.1-0", true},
		{"<1.1.0-0", "1.0.0", true},
		{"<1.1.0-0", "1.0.0-0", true},
		{"<1.1.0-0", "0.0.0", true},
		{"<1.1.0-0", "0.0.0-0", true},
		{"<1.0.1", "2.0.0", false},
		{"<1.0.1", "2.0.0-0", false},
		{"<1.0.1", "1.1.0", false},
		{"<1.0.1", "1.1.0-0", false},
		{"<1.0.1", "1.0.1", false},
		{"<1.0.1", "1.0.1-0", true},
		{"<1.0.1", "1.0.0", true},
		{"<1.0.1", "1.0.0-0", true},
		{"<1.0.1", "0.0.0", true},
		{"<1.0.1", "0.0.0-0", true},
		{"<1.0.1-0", "2.0.0", false},
		{"<1.0.1-0", "2.0.0-0", false},
		{"<1.0.1-0", "1.1.0", false},
		{"<1.0.1-0", "1.1.0-0", false},
		{"<1.0.1-0", "1.0.1", false},
		{"<1.0.1-0", "1.0.1-0", false},
		{"<1.0.1-0", "1.0.0", true},
		{"<1.0.1-0", "1.0.0-0", true},
		{"<1.0.1-0", "0.0.0", true},
		{"<1.0.1-0", "0.0.0-0", true},
		{"<1.0.0", "2.0.0", false},
		{"<1.0.0", "2.0.0-0", false},
		{"<1.0.0", "1.1.0", false},
		{"<1.0.0", "1.1.0-0", false},
		{"<1.0.0", "1.0.1", false},
		{"<1.0.0", "1.0.1-0", false},
		{"<1.0.0", "1.0.0", false},
		{"<1.0.0", "1.0.0-0", true},
		{"<1.0.0", "0.0.0", true},
		{"<1.0.0", "0.0.0-0", true},
		{"<1.0.0-0", "2.0.0", false},
		{"<1.0.0-0", "2.0.0-0", false},
		{"<1.0.0-0", "1.1.0", false},
		{"<1.0.0-0", "1.1.0-0", false},
		{"<1.0.0-0", "1.0.1", false},
		{"<1.0.0-0", "1.0.1-0", false},
		{"<1.0.0-0", "1.0.0", false},
		{"<1.0.0-0", "1.0.0-0", false},
		{"<1.0.0-0", "0.0.0", true},
		{"<1.0.0-0", "0.0.0-0", true},

		// <= wildcard major (matches everything)
		{"<=*", "2.0.0", true},
		{"<=*", "2.0.0-0", true},
		{"<=*", "1.1.0", true},
		{"<=*", "1.1.0-0", true},
		{"<=*", "1.0.1", true},
		{"<=*", "1.0.1-0", true},
		{"<=*", "1.0.0", true},
		{"<=*", "1.0.0-0", true},
		{"<=*", "0.0.0", true},
		{"<=*", "0.0.0-0", true},

		// <= wildcard minor
		{"<=1", "2.0.0", false},
		{"<=1", "2.0.0-0", false},
		{"<=1", "1.1.0", true},
		{"<=1", "1.1.0-0", true},
		{"<=1", "1.0.1", true},
		{"<=1", "1.0.1-0", true},
		{"<=1", "1.0.0", true},
		{"<=1", "1.0.0-0", true},
		{"<=1", "0.0.0", true},
		{"<=1", "0.0.0-0", true},

		// <= wildcard patch
		{"<=1.1", "2.0.0", false},
		{"<=1.1", "2.0.0-0", false},
		{"<=1.1", "1.1.0", true},
		{"<=1.1", "1.1.0-0", true},
		{"<=1.1", "1.0.1", true},
		{"<=1.1", "1.0.1-0", true},
		{"<=1.1", "1.0.0", true},
		{"<=1.1", "1.0.0-0", true},
		{"<=1.1", "0.0.0", true},
		{"<=1.1", "0.0.0-0", true},
		{"<=1.0", "2.0.0", false},
		{"<=1.0", "2.0.0-0", false},
		{"<=1.0", "1.1.0", false},
		{"<=1.0", "1.1.0-0", false},
		{"<=1.0", "1.0.1", true},
		{"<=1.0", "1.0.1-0", true},
		{"<=1.0", "1.0.0", true},
		{"<=1.0", "1.0.0-0", true},
		{"<=1.0", "0.0.0", true},
		{"<=1.0", "0.0.0-0", true},

		// <= exact
		{"<=1.1.0", "2.0.0", false},
		{"<=1.1.0", "2.0.0-0", false},
		{"<=1.1.0", "1.1.0", true},
		{"<=1.1.0", "1.1.0-0", true},
		{"<=1.1.0", "1.0.1", true},
		{"<=1.1.0", "1.0.1-0", true},
		{"<=1.1.0", "1.0.0", true},
		{"<=1.1.0", "1.0.0-0", true},
		{"<=1.1.0", "0.0.0", true},
		{"<=1.1.0", "0.0.0-0", true},
		{"<=1.1.0-0", "2.0.0", false},
		{"<=1.1.0-0", "2.0.0-0", false},
		{"<=1.1.0-0", "1.1.0", false},
		{"<=1.1.0-0", "1.1.0-0", true},
		{"<=1.1.0-0", "1.0.1", true},
		{"<=1.1.0-0", "1.0.1-0", true},
		{"<=1.1.0-0", "1.0.0", true},
		{"<=1.1.0-0", "1.0.0-0", true},
		{"<=1.1.0-0", "0.0.0", true},
		{"<=1.1.0-0", "0.0.0-0", true},
		{"<=1.0.1", "2.0.0", false},
		{"<=1.0.1", "2.0.0-0", false},
		{"<=1.0.1", "1.1.0", false},
		{"<=1.0.1", "1.1.0-0", false},
		{"<=1.0.1", "1.0.1", true},
		{"<=1.0.1", "1.0.1-0", true},
		{"<=1.0.1", "1.0.0", true},
		{"<=1.0.1", "1.0.0-0", true},
		{"<=1.0.1", "0.0.0", true},
		{"<=1.0.1", "0.0.0-0", true},
		{"<=1.0.1-0", "2.0.0", false},
		{"<=1.0.1-0", "2.0.0-0", false},
		{"<=1.0.1-0", "1.1.0", false},
		{"<=1.0.1-0", "1.1.0-0", false},
		{"<=1.0.1-0", "1.0.1", false},
		{"<=1.0.1-0", "1.0.1-0", true},
		{"<=1.0.1-0", "1.0.0", true},
		{"<=1.0.1-0", "1.0.0-0", true},
		{"<=1.0.1-0", "0.0.0", true},
		{"<=1.0.1-0", "0.0.0-0", true},
		{"<=1.0.0", "2.0.0", false},
		{"<=1.0.0", "2.0.0-0", false},
		{"<=1.0.0", "1.1.0", false},
		{"<=1.0.0", "1.1.0-0", false},
		{"<=1.0.0", "1.0.1", false},
		{"<=1.0.0", "1.0.1-0", false},
		{"<=1.0.0", "1.0.0", true},
		{"<=1.0.0", "1.0.0-0", true},
		{"<=1.0.0", "0.0.0", true},
		{"<=1.0.0", "0.0.0-0", true},
		{"<=1.0.0-0", "2.0.0", false},
		{"<=1.0.0-0", "2.0.0-0", false},
		{"<=1.0.0-0", "1.1.0", false},
		{"<=1.0.0-0", "1.1.0-0", false},
		{"<=1.0.0-0", "1.0.1", false},
		{"<=1.0.0-0", "1.0.1-0", false},
		{"<=1.0.0-0", "1.0.0", false},
		{"<=1.0.0-0", "1.0.0-0", true},
		{"<=1.0.0-0", "0.0.0", true},
		{"<=1.0.0-0", "0.0.0-0", true},

		// https://github.com/microsoft/TypeScript/issues/50909
		{">4.8", "4.9.0-beta", true},
		{">=4.9", "4.9.0-beta", true},
		{"<4.9", "4.9.0-beta", false},
		{"<=4.8", "4.9.0-beta", false},
	}
	for _, test := range comparatorsTests {
		assertRangeTest(t, "comparators", test.rangeText, test.versionText, test.expected)
	}
}

func TestConjunctionsOfVersionRanges(t *testing.T) {
	conjunctionTests := []testForRangeOnVersion{
		{">1.0.0 <2.0.0", "1.0.1", true},
		{">1.0.0 <2.0.0", "2.0.0", false},
		{">1.0.0 <2.0.0", "1.0.0", false},
		{">1 >2", "3.0.0", true},
	}
	for _, test := range conjunctionTests {
		assertRangeTest(t, "conjunctions", test.rangeText, test.versionText, test.expected)
	}
}

func TestDisjunctionsOfVersionRanges(t *testing.T) {
	disjunctionTests := []testForRangeOnVersion{
		{">1.0.0 || <1.0.0", "1.0.1", true},
		{">1.0.0 || <1.0.0", "0.0.1", true},
		{">1.0.0 || <1.0.0", "1.0.0", false},
		{">1.0.0 || <1.0.0", "0.0.0", true},
		{">=1.0.0 <2.0.0 || >=3.0.0 <4.0.0", "1.0.0", true},
		{">=1.0.0 <2.0.0 || >=3.0.0 <4.0.0", "2.0.0", false},
		{">=1.0.0 <2.0.0 || >=3.0.0 <4.0.0", "3.0.0", true},
	}
	for _, test := range disjunctionTests {
		assertRangeTest(t, "disjunctions", test.rangeText, test.versionText, test.expected)
	}
}

func TestHyphensOfVersionRanges(t *testing.T) {
	hyphenTests := []testForRangeOnVersion{
		{"1.0.0 - 2.0.0", "1.0.0", true},
		{"1.0.0 - 2.0.0", "1.0.1", true},
		{"1.0.0 - 2.0.0", "2.0.0", true},
		{"1.0.0 - 2.0.0", "2.0.1", false},
		{"1.0.0 - 2.0.0", "0.9.9", false},
		{"1.0.0 - 2.0.0", "3.0.0", false},
	}
	for _, test := range hyphenTests {
		assertRangeTest(t, "hyphens", test.rangeText, test.versionText, test.expected)
	}
}

func TestTildesOfVersionRanges(t *testing.T) {
	tildeTests := []testForRangeOnVersion{
		{"~0", "0.0.0", true},
		{"~0", "0.1.0", true},
		{"~0", "0.1.2", true},
		{"~0", "0.1.9", true},
		{"~0", "1.0.0", false},
		{"~0.1", "0.1.0", true},
		{"~0.1", "0.1.2", true},
		{"~0.1", "0.1.9", true},
		{"~0.1", "0.2.0", false},
		{"~0.1.2", "0.1.2", true},
		{"~0.1.2", "0.1.9", true},
		{"~0.1.2", "0.2.0", false},
		{"~1.0.0", "1.0.0", true},
		{"~1.0.0", "1.0.1", true},
		{"~1", "1.0.0", true},
		{"~1", "1.2.0", true},
		{"~1", "1.2.3", true},
		{"~1", "1.2.0", true},
		{"~1", "1.2.3", true},
		{"~1", "0.0.0", false},
		{"~1", "2.0.0", false},
		{"~1.2", "1.2.0", true},
		{"~1.2", "1.2.3", true},
		{"~1.2", "1.1.0", false},
		{"~1.2", "1.3.0", false},
		{"~1.2.3", "1.2.3", true},
		{"~1.2.3", "1.2.9", true},
		{"~1.2.3", "1.1.0", false},
		{"~1.2.3", "1.3.0", false},
	}
	for _, test := range tildeTests {
		assertRangeTest(t, "tilde", test.rangeText, test.versionText, test.expected)
	}
}

func TestCaretsOfVersionRanges(t *testing.T) {
	caretTests := []testForRangeOnVersion{
		{"^0", "0.0.0", true},
		{"^0", "0.1.0", true},
		{"^0", "0.9.0", true},
		{"^0", "0.1.2", true},
		{"^0", "0.1.9", true},
		{"^0", "1.0.0", false},
		{"^0.1", "0.1.0", true},
		{"^0.1", "0.1.2", true},
		{"^0.1", "0.1.9", true},
		{"^0.1.2", "0.1.2", true},
		{"^0.1.2", "0.1.9", true},
		{"^0.1.2", "0.0.0", false},
		{"^0.1.2", "0.2.0", false},
		{"^0.1.2", "1.0.0", false},
		{"^1", "1.0.0", true},
		{"^1", "1.2.0", true},
		{"^1", "1.2.3", true},
		{"^1", "1.9.0", true},
		{"^1", "0.0.0", false},
		{"^1", "2.0.0", false},
		{"^1.2", "1.2.0", true},
		{"^1.2", "1.2.3", true},
		{"^1.2", "1.9.0", true},
		{"^1.2", "1.1.0", false},
		{"^1.2", "2.0.0", false},
		{"^1.2.3", "1.2.3", true},
		{"^1.2.3", "1.9.0", true},
		{"^1.2.3", "1.2.2", false},
		{"^1.2.3", "2.0.0", false},
	}
	for _, test := range caretTests {
		assertRangeTest(t, "caret", test.rangeText, test.versionText, test.expected)
	}
}

type testForRangeOnVersion struct {
	rangeText   string
	versionText string
	expected    bool
}

func assertRangesGoodBad(t *testing.T, versionRangeString string, tests testGoodBad) {
	t.Run(versionRangeString, func(t *testing.T) {
		versionRange, ok := TryParseVersionRange(versionRangeString)
		assert.Assert(t, ok)
		for _, good := range tests.good {
			v, ok := TryParseVersion(good)
			assert.Assert(t, ok)
			assert.Assert(t, versionRange.Test(&v), "%s should be matched by range %s", good, versionRangeString)
		}

		for _, bad := range tests.bad {
			v, ok := TryParseVersion(bad)
			assert.Assert(t, ok)
			assert.Assert(t, !versionRange.Test(&v), "%s should not be matched by range %s", bad, versionRangeString)
		}
	})
}

func assertRangeTest(t *testing.T, name string, rangeText string, versionText string, inRange bool) {
	testName := fmt.Sprintf("%s (version %s in range %s) == %t", name, versionText, rangeText, inRange)
	t.Run(testName, func(t *testing.T) {
		versionRange, ok := TryParseVersionRange(rangeText)
		assert.Assert(t, ok)
		version, err := TryParseVersion(versionText)
		assert.NilError(t, err)
		assert.Equal(t, versionRange.Test(&version), inRange)
	})
}
