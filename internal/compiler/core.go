package compiler

import "strings"

func EquateStringsCaseInsensitive(a, b string) bool {
	// !!!
	// return a == b || strings.ToUpper(a) == strings.ToUpper(b)
	return strings.EqualFold(a, b)
}

func EquateStringsCaseSensitive(a, b string) bool {
	return a == b
}

func GetStringEqualityComparer(ignoreCase bool) func(a, b string) bool {
	if ignoreCase {
		return EquateStringsCaseInsensitive
	}
	return EquateStringsCaseSensitive
}

type Comparison = int32

const (
	ComparisonLessThan    Comparison = -1
	ComparisonEqual       Comparison = 0
	ComparisonGreaterThan Comparison = 1
)

func CompareStringsCaseInsensitive(a string, b string) Comparison {
	if a == b {
		return ComparisonEqual
	}
	return int32(strings.Compare(strings.ToUpper(a), strings.ToUpper(b)))
}

func CompareStringsCaseSensitive(a string, b string) Comparison {
	return int32(strings.Compare(a, b))
}

func GetStringComparer(ignoreCase bool) func(a, b string) Comparison {
	if ignoreCase {
		return CompareStringsCaseInsensitive
	}
	return CompareStringsCaseSensitive
}
