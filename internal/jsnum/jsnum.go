// Package jsnum provides JS-like number handling.
package jsnum

import (
	"encoding/json"
	"errors"
	"math"
	"math/big"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/microsoft/typescript-go/internal/stringutil"
)

const (
	MaxSafeInteger Number = 1<<53 - 1
	MinSafeInteger Number = -MaxSafeInteger
)

// Number represents a JS-like number.
//
// All operations that can be performed directly on this type
// (e.g., conversion, arithmetic, etc.) behave as they would in JavaScript,
// but any other operation should use this type's methods,
// not the "math" package and conversions.
type Number float64

func (n Number) String() string {
	switch {
	case n.IsNaN():
		return "NaN"
	case n.IsInf():
		if n < 0 {
			return "-Infinity"
		}
		return "Infinity"
	}

	// Fast path: for safe integers, directly convert to string.
	if MinSafeInteger <= n && n <= MaxSafeInteger {
		if i := int64(n); float64(i) == float64(n) {
			return strconv.FormatInt(i, 10)
		}
	}

	// Otherwise, the Go json package handles this correctly.
	b, _ := json.Marshal(float64(n)) //nolint:errchkjson
	return string(b)
}

// https://tc39.es/ecma262/2024/multipage/abstract-operations.html#sec-stringtonumber
func FromString(s string) Number {
	// Implemeting StringToNumber exactly as written in the spec involves
	// writing a parser, along with the conversion of the parsed AST into the
	// actual value.
	//
	// We've already implemented a number parser in the scanner, but we can't
	// import it here. We also do not have the conversion implemented since we
	// previously just wrote `+literal` and let the runtime handle it.
	//
	// The strategy below is to instead break the number apart and fix it up
	// such that Go's own parsing functionality can handle it. This won't be
	// the fastest method, but it saves us from writing the full parser and
	// conversion logic.

	s = strings.TrimFunc(s, isStrWhiteSpace)

	switch s {
	case "":
		return 0
	case "Infinity", "+Infinity":
		return Inf(1)
	case "-Infinity":
		return Inf(-1)
	}

	for _, r := range s {
		if !isNumberRune(r) {
			return NaN()
		}
	}

	if n, ok := tryParseInt(s); ok {
		return n
	}

	// Cut this off first so we can ensure -0 is returned as -0.
	s, negative := strings.CutPrefix(s, "-")

	if !negative {
		s, _ = strings.CutPrefix(s, "+")
	}

	if first, _ := utf8.DecodeRuneInString(s); !stringutil.IsDigit(first) && first != '.' {
		return NaN()
	}

	f := parseFloatString(s)
	if math.IsNaN(f) {
		return NaN()
	}

	sign := 1.0
	if negative {
		sign = -1.0
	}
	return Number(math.Copysign(f, sign))
}

func isStrWhiteSpace(r rune) bool {
	// This is different than stringutil.IsWhiteSpaceLike.

	// https://tc39.es/ecma262/2024/multipage/ecmascript-language-lexical-grammar.html#prod-LineTerminator
	// https://tc39.es/ecma262/2024/multipage/ecmascript-language-lexical-grammar.html#prod-WhiteSpace

	switch r {
	// LineTerminator
	case '\n', '\r', 0x2028, 0x2029:
		return true
	// WhiteSpace
	case '\t', '\v', '\f', 0xFEFF:
		return true
	}

	// WhiteSpace
	return unicode.Is(unicode.Zs, r)
}

var errUnknownPrefix = errors.New("unknown number prefix")

func tryParseInt(s string) (Number, bool) {
	var i int64
	var err error
	var hasIntResult bool

	if len(s) > 2 {
		prefix, rest := s[:2], s[2:]
		switch prefix {
		case "0b", "0B":
			if !isAllBinaryDigits(rest) {
				return NaN(), true
			}
			i, err = strconv.ParseInt(rest, 2, 64)
			hasIntResult = true
		case "0o", "0O":
			if !isAllOctalDigits(rest) {
				return NaN(), true
			}
			i, err = strconv.ParseInt(rest, 8, 64)
			hasIntResult = true
		case "0x", "0X":
			if !isAllHexDigits(rest) {
				return NaN(), true
			}
			i, err = strconv.ParseInt(rest, 16, 64)
			hasIntResult = true
		}
	}

	if !hasIntResult {
		// StringToNumber does not parse leading zeros as octal.
		s = trimLeadingZeros(s)
		if !isAllDigits(s) {
			return 0, false
		}
		i, err = strconv.ParseInt(s, 10, 64)
		hasIntResult = true
	}

	if hasIntResult && err == nil {
		return Number(i), true
	}

	// Using this to parse large integers.
	bi, ok := new(big.Int).SetString(s, 0)
	if !ok {
		return NaN(), true
	}

	f, _ := bi.Float64()
	return Number(f), true
}

func parseFloatString(s string) float64 {
	var hasDot, hasExp bool

	// <a>
	// <a>.<b>
	// <a>.<b>e<c>
	// <a>e<c>
	var a, b, c, rest string

	a, rest, hasDot = strings.Cut(s, ".")
	if hasDot {
		// <a>.<b>
		// <a>.<b>e<c>
		b, c, hasExp = cutAny(rest, "eE")
	} else {
		// <a>
		// <a>e<c>
		a, c, hasExp = cutAny(s, "eE")
	}

	var sb strings.Builder
	sb.Grow(len(a) + len(b) + len(c) + 3)

	if a == "" {
		if hasDot && b == "" {
			return math.NaN()
		}
		if hasExp && c == "" {
			return math.NaN()
		}
		sb.WriteString("0")
	} else {
		a = trimLeadingZeros(a)
		if !isAllDigits(a) {
			return math.NaN()
		}
		sb.WriteString(a)
	}

	if hasDot {
		sb.WriteString(".")
		if b == "" {
			sb.WriteString("0")
		} else {
			b = trimTrailingZeros(b)
			if !isAllDigits(b) {
				return math.NaN()
			}
			sb.WriteString(b)
		}
	}

	if hasExp {
		sb.WriteString("e")

		c, negative := strings.CutPrefix(c, "-")
		if negative {
			sb.WriteString("-")
		} else {
			c, _ = strings.CutPrefix(c, "+")
		}
		c = trimLeadingZeros(c)
		if !isAllDigits(c) {
			return math.NaN()
		}
		sb.WriteString(c)
	}

	return stringToFloat64(sb.String())
}

func cutAny(s string, cutset string) (before, after string, found bool) {
	if i := strings.IndexAny(s, cutset); i >= 0 {
		before = s[:i]
		afterAndFound := s[i:]
		_, size := utf8.DecodeRuneInString(afterAndFound)
		after = afterAndFound[size:]
		return before, after, true
	}
	return s, "", false
}

func trimLeadingZeros(s string) string {
	for strings.HasPrefix(s, "0") {
		rest, ok := strings.CutPrefix(s, "0")
		if !ok {
			return s
		}
		if rest == "" {
			return "0"
		}
		s = rest
	}
	return s
}

func trimTrailingZeros(s string) string {
	for strings.HasSuffix(s, "0") {
		rest, ok := strings.CutSuffix(s, "0")
		if !ok {
			return s
		}
		if rest == "" {
			return "0"
		}
		s = rest
	}
	return s
}

func stringToFloat64(s string) float64 {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	} else {
		if errors.Is(err, strconv.ErrRange) {
			return f
		}
	}

	i, err := strconv.ParseInt(s, 0, 64)
	if err == nil {
		return float64(i)
	}

	return math.NaN()
}

func isAllDigits(s string) bool {
	for _, r := range s {
		if !stringutil.IsDigit(r) {
			return false
		}
	}
	return true
}

func isAllBinaryDigits(s string) bool {
	for _, r := range s {
		if r != '0' && r != '1' {
			return false
		}
	}
	return true
}

func isAllOctalDigits(s string) bool {
	for _, r := range s {
		if !stringutil.IsOctalDigit(r) {
			return false
		}
	}
	return true
}

func isAllHexDigits(s string) bool {
	for _, r := range s {
		if !stringutil.IsHexDigit(r) {
			return false
		}
	}
	return true
}

func isNumberRune(r rune) bool {
	if stringutil.IsDigit(r) {
		return true
	}

	if 'a' <= r && r <= 'f' {
		return true
	}

	if 'A' <= r && r <= 'F' {
		return true
	}

	switch r {
	case '.', '-', '+', 'x', 'X', 'o', 'O':
		return true
	}

	return false
}

func NaN() Number {
	return Number(math.NaN())
}

func (n Number) IsNaN() bool {
	return math.IsNaN(float64(n))
}

func Inf(sign int) Number {
	return Number(math.Inf(sign))
}

func (n Number) IsInf() bool {
	return math.IsInf(float64(n), 0)
}

func isNonFinite(x float64) bool {
	// This is equivalent to checking `math.IsNaN(x) || math.IsInf(x, 0)` in one operation.
	const mask = 0x7FF0000000000000
	return math.Float64bits(x)&mask == mask
}

// https://tc39.es/ecma262/2024/multipage/abstract-operations.html#sec-touint32
func (n Number) toUint32() uint32 {
	x := float64(n)
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
	return Number(x.toInt32() >> y.toShiftCount())
}

// https://tc39.es/ecma262/2024/multipage/ecmascript-data-types-and-values.html#sec-numeric-types-number-unsignedRightShift
func (x Number) UnsignedRightShift(y Number) Number {
	return Number(x.toUint32() >> y.toShiftCount())
}

// https://tc39.es/ecma262/2024/multipage/ecmascript-data-types-and-values.html#sec-numeric-types-number-leftShift
func (x Number) LeftShift(y Number) Number {
	return Number(x.toInt32() << y.toShiftCount())
}

// https://tc39.es/ecma262/2024/multipage/ecmascript-data-types-and-values.html#sec-numeric-types-number-bitwiseNOT
func (x Number) BitwiseNOT() Number {
	return Number(^x.toInt32())
}

// The below are implemented by https://tc39.es/ecma262/2024/multipage/ecmascript-data-types-and-values.html#sec-numberbitwiseop.

// https://tc39.es/ecma262/2024/multipage/ecmascript-data-types-and-values.html#sec-numeric-types-number-bitwiseOR
func (x Number) BitwiseOR(y Number) Number {
	return Number(x.toInt32() | y.toInt32())
}

// https://tc39.es/ecma262/2024/multipage/ecmascript-data-types-and-values.html#sec-numeric-types-number-bitwiseAND
func (x Number) BitwiseAND(y Number) Number {
	return Number(x.toInt32() & y.toInt32())
}

// https://tc39.es/ecma262/2024/multipage/ecmascript-data-types-and-values.html#sec-numeric-types-number-bitwiseXOR
func (x Number) BitwiseXOR(y Number) Number {
	return Number(x.toInt32() ^ y.toInt32())
}

func (x Number) trunc() Number {
	return Number(math.Trunc(float64(x)))
}

func (x Number) Floor() Number {
	return Number(math.Floor(float64(x)))
}

func (x Number) Abs() Number {
	return Number(math.Abs(float64(x)))
}

var negativeZero = Number(math.Copysign(0, -1))

// https://tc39.es/ecma262/2024/multipage/ecmascript-data-types-and-values.html#sec-numeric-types-number-remainder
func (n Number) Remainder(d Number) Number {
	switch {
	case n.IsNaN() || d.IsNaN():
		return NaN()
	case n.IsInf():
		return NaN()
	case d.IsInf():
		return n
	case d == 0:
		return NaN()
	case n == 0:
		return n
	}

	r := n - d*(n/d).trunc()
	if r == 0 || n < 0 {
		return negativeZero
	}

	return r
}

// https://tc39.es/ecma262/2024/multipage/ecmascript-data-types-and-values.html#sec-numeric-types-number-exponentiate
func (base Number) Exponentiate(exponent Number) Number {
	switch {
	case (base == 1 || base == -1) && exponent.IsInf():
		return NaN()
	case base == 1 && exponent.IsNaN():
		return NaN()
	}

	return Number(math.Pow(float64(base), float64(exponent)))
}
