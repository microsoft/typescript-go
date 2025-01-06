// Package jsnum provides JS-like number handling.
package jsnum

import "math"

// https://262.ecma-international.org/#sec-touint32
func toUint32(x float64) uint32 {
	// Fast path: if the number is the range (-2^31, 2^32), i.e. an SMI,
	// then we don't need to do any special mapping.
	if smi := int32(x); float64(smi) == x {
		return uint32(smi)
	}

	// If the number is non-finite (NaN, +Inf, -Inf; exp=0x7FF), it maps to zero.
	// This is equivalent to checking `math.IsNaN(x) || math.IsInf(x, 0)` in one operation.
	const mask = 0x7FF0000000000000
	if math.Float64bits(x)&mask == mask {
		return 0
	}

	// Otherwise, take x modulo 2^32, mapping positive numbers
	// to [0, 2^32) and negative numbers to (-2^32, -0.0].
	x = math.Mod(x, 1<<32)

	// Convert to uint32, which will wrap negative numbers.
	return uint32(x)
}

// https://262.ecma-international.org/#sec-toint32
func toInt32(x float64) int32 {
	// The only difference between ToUint32 and ToInt32 is the interpretation of the bits.
	return int32(toUint32(x))
}

func toShiftCount(x float64) uint32 {
	return toUint32(x) & 31
}

// https://262.ecma-international.org/#sec-numeric-types-number-signedRightShift
func SignedRightShift(x, y float64) float64 {
	return float64(toInt32(x) >> toShiftCount(y))
}

// https://262.ecma-international.org/#sec-numeric-types-number-unsignedRightShift
func UnsignedRightShift(x, y float64) float64 {
	return float64(toUint32(x) >> toShiftCount(y))
}

// https://262.ecma-international.org/#sec-numeric-types-number-leftShift
func LeftShift(x, y float64) float64 {
	return float64(toInt32(x) << toShiftCount(y))
}

// https://262.ecma-international.org/#sec-numeric-types-number-bitwiseNOT
func BitwiseNOT(x float64) float64 {
	return float64(^toInt32(x))
}

// The below are implemented by https://262.ecma-international.org/#sec-numberbitwiseop.

// https://262.ecma-international.org/#sec-numeric-types-number-bitwiseOR
func BitwiseOR(x, y float64) float64 {
	return float64(toInt32(x) | toInt32(y))
}

// https://262.ecma-international.org/#sec-numeric-types-number-bitwiseAND
func BitwiseAND(x, y float64) float64 {
	return float64(toInt32(x) & toInt32(y))
}

// https://262.ecma-international.org/#sec-numeric-types-number-bitwiseXOR
func BitwiseXOR(x, y float64) float64 {
	return float64(toInt32(x) ^ toInt32(y))
}

var negativeZero = math.Copysign(0, -1)

// https://262.ecma-international.org/#sec-numeric-types-number-remainder
func Remainder(n, d float64) float64 {
	if math.IsNaN(n) || math.IsNaN(d) {
		return math.NaN()
	}

	if math.IsInf(n, 0) {
		return math.NaN()
	}

	if math.IsInf(d, 0) {
		return n
	}

	if d == 0 {
		return math.NaN()
	}

	if n == 0 {
		return n
	}

	r := n - d*math.Trunc(n/d)
	if r == 0 || n < 0 {
		return negativeZero
	}

	return r
}

// https://262.ecma-international.org/#sec-numeric-types-number-exponentiate
func Exponentiate(base, exponent float64) float64 {
	// Special cases in https://262.ecma-international.org/#sec-numeric-types-number-exponentiate
	if (base == 1 || base == -1) && math.IsInf(exponent, 0) {
		return math.NaN()
	}
	if base == 1 && math.IsNaN(exponent) {
		return math.NaN()
	}

	return math.Pow(base, exponent)
}
