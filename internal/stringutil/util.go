package stringutil

import "strings"

func EquateCaseInsensitive(a, b string) bool {
	// !!!
	// return a == b || strings.ToUpper(a) == strings.ToUpper(b)
	return strings.EqualFold(a, b)
}

func EquateCaseSensitive(a, b string) bool {
	return a == b
}
