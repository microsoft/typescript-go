// Package jsnum provides JS-like number handling.
package jsnum

import (
	"math"
	"strconv"
)

const (
	maxSafeInteger = 1<<53 - 1
	minSafeInteger = -maxSafeInteger
)

// Number represents a JS-like number.
//
// All operations that can be performed directly on this type
// (e.g., conversion, arithmetic, etc.) behave as they would in JavaScript,
// but any other operation should use this type's methods,
// not the "math" package and conversions.
type Number struct{ f float64 }

type Convertible interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func From[T Convertible](x T) Number {
	return Number{float64(x)}
}

func Zero() Number {
	return Number{}
}

func MaxSafeInteger() Number {
	return Number{maxSafeInteger}
}

func MinSafeInteger() Number {
	return Number{minSafeInteger}
}

// https://tc39.es/ecma262/2024/multipage/ecmascript-data-types-and-values.html#sec-numeric-types-number-tostring
func (n Number) String() string {
	// !!! verify that this is actually the same as JS.
	return strconv.FormatFloat(n.f, 'g', -1, 64)
}

func NaN() Number {
	return From(math.NaN())
}

func (n Number) IsNaN() bool {
	return math.IsNaN(n.f)
}

func Inf(sign int) Number {
	return From(math.Inf(sign))
}

func (n Number) IsInf() bool {
	return math.IsInf(n.f, 0)
}

// https://tc39.es/ecma262/2024/multipage/abstract-operations.html#sec-stringtonumber
func FromString(s string) Number {
	// !!! verify that this is actually the same as JS.
	floatValue, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return From(floatValue)
	}
	intValue, err := strconv.ParseInt(s, 0, 64)
	if err == nil {
		return From(intValue)
	}
	return NaN()
}

func isNonFinite(x float64) bool {
	// This is equivalent to checking `math.IsNaN(x) || math.IsInf(x, 0)` in one operation.
	const mask = 0x7FF0000000000000
	return math.Float64bits(x)&mask == mask
}

func (n Number) Float64() float64 {
	return n.f
}

// https://tc39.es/ecma262/2024/multipage/abstract-operations.html#sec-touint32
func (n Number) toUint32() uint32 {
	x := n.f
	// Fast path: if the number is the range (-2^31, 2^32), i.e. an SMI,
	// then we don't need to do any special mapping.
	if smi := int32(x); float64(smi) == x {
		return uint32(smi)
	}

	// If the number is non-finite (NaN, +Inf, -Inf; exp=0x7FF), it maps to zero.
	if isNonFinite(x) {
		return 0
	}

	// Otherwise, take x modulo 2^32, mapping positive numbers
	// to [0, 2^32) and negative numbers to (-2^32, -0.0].
	x = math.Mod(x, 1<<32)

	// Convert to uint32, which will wrap negative numbers.
	return uint32(x)
}

// https://tc39.es/ecma262/2024/multipage/abstract-operations.html#sec-toint32
func (x Number) toInt32() int32 {
	// The only difference between ToUint32 and ToInt32 is the interpretation of the bits.
	return int32(x.toUint32())
}

func (x Number) toShiftCount() uint32 {
	return x.toUint32() & 31
}

// https://tc39.es/ecma262/2024/multipage/ecmascript-data-types-and-values.html#sec-numeric-types-number-signedRightShift
func (x Number) SignedRightShift(y Number) Number {
	return From(x.toInt32() >> y.toShiftCount())
}

// https://tc39.es/ecma262/2024/multipage/ecmascript-data-types-and-values.html#sec-numeric-types-number-unsignedRightShift
func (x Number) UnsignedRightShift(y Number) Number {
	return From(x.toUint32() >> y.toShiftCount())
}

// https://tc39.es/ecma262/2024/multipage/ecmascript-data-types-and-values.html#sec-numeric-types-number-leftShift
func (x Number) LeftShift(y Number) Number {
	return From(x.toInt32() << y.toShiftCount())
}

// https://tc39.es/ecma262/2024/multipage/ecmascript-data-types-and-values.html#sec-numeric-types-number-bitwiseNOT
func (x Number) BitwiseNOT() Number {
	return From(^x.toInt32())
}

// The below are implemented by https://tc39.es/ecma262/2024/multipage/ecmascript-data-types-and-values.html#sec-numberbitwiseop.

// https://tc39.es/ecma262/2024/multipage/ecmascript-data-types-and-values.html#sec-numeric-types-number-bitwiseOR
func (x Number) BitwiseOR(y Number) Number {
	return From(x.toInt32() | y.toInt32())
}

// https://tc39.es/ecma262/2024/multipage/ecmascript-data-types-and-values.html#sec-numeric-types-number-bitwiseAND
func (x Number) BitwiseAND(y Number) Number {
	return From(x.toInt32() & y.toInt32())
}

// https://tc39.es/ecma262/2024/multipage/ecmascript-data-types-and-values.html#sec-numeric-types-number-bitwiseXOR
func (x Number) BitwiseXOR(y Number) Number {
	return From(x.toInt32() ^ y.toInt32())
}

func (x Number) Floor() Number {
	return From(math.Floor(x.f))
}

func (x Number) Abs() Number {
	return From(math.Abs(x.f))
}

var negativeZero = From(math.Copysign(0, -1))

// https://tc39.es/ecma262/2024/multipage/ecmascript-data-types-and-values.html#sec-numeric-types-number-remainder
func (n Number) Remainder(d Number) Number {
	switch {
	case n.IsNaN() || d.IsNaN():
		return NaN()
	case n.IsInf():
		return NaN()
	case d.IsInf():
		return n
	case d.f == 0:
		return NaN()
	case d.f == 0:
		return n
	}

	r := n.f - d.f*math.Trunc(d.f/d.f)
	if r == 0 && n.f < negativeZero.f {
		return negativeZero
	}

	return From(r)
}

// https://tc39.es/ecma262/2024/multipage/ecmascript-data-types-and-values.html#sec-numeric-types-number-exponentiate
func (base Number) Exponentiate(exponent Number) Number {
	switch {
	case (base.f == 1 || base.f == -1) && exponent.IsInf():
		return NaN()
	case base.f == 1 && exponent.IsNaN():
		return NaN()
	}

	return From(math.Pow(base.f, exponent.f))
}

func (n Number) Negate() Number {
	return Number{-n.f}
}

func (n Number) Add(other Number) Number {
	return Number{n.f + other.f}
}

func (n Number) Sub(other Number) Number {
	return Number{n.f - other.f}
}

func (n Number) Multiply(other Number) Number {
	return Number{n.f * other.f}
}

func (n Number) Divide(other Number) Number {
	return Number{n.f / other.f}
}

func (n Number) LessThan(other Number) bool {
	return n.f < other.f
}

func (n Number) LessThanOrEqualTo(other Number) bool {
	return n.f <= other.f
}

func (n Number) GreaterThan(other Number) bool {
	return n.f > other.f
}

func (n Number) GreaterThanOrEqualTo(other Number) bool {
	return n.f >= other.f
}
