package semver

import (
	"math"
	"strconv"
)

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
