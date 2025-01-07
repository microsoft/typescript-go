package jsnum

import (
	"fmt"
	"math"
	"testing"

	"gotest.tools/v3/assert"
)

func assertEqualNumber(t *testing.T, got, want Number) {
	t.Helper()
	if want.IsNaN() {
		assert.Assert(t, got.IsNaN())
	} else {
		assert.Equal(t, got, want)
	}
}

var toInt32Tests = []struct {
	name  string
	input float64
	want  int32
	bench bool
}{
	{"0.0", 0, 0, true},
	{"-0.0", negativeZero.f, 0, false},
	{"NaN", NaN().f, 0, true},
	{"+Inf", Inf(1).f, 0, true},
	{"-Inf", Inf(-1).f, 0, true},
	{"MaxInt32", math.MaxInt32, math.MaxInt32, false},
	{"MaxInt32+1", float64(int64(math.MaxInt32) + 1), math.MinInt32, true},
	{"MinInt32", float64(math.MinInt32), math.MinInt32, false},
	{"MinInt32-1", float64(int64(math.MinInt32) - 1), math.MaxInt32, true},
	{"MIN_SAFE_INTEGER", MinSafeInteger().f, 1, false},
	{"MIN_SAFE_INTEGER-1", MinSafeInteger().f - 1, 0, false},
	{"MIN_SAFE_INTEGER+1", MinSafeInteger().f + 1, 2, false},
	{"MAX_SAFE_INTEGER", MaxSafeInteger().f, -1, true},
	{"MAX_SAFE_INTEGER-1", MaxSafeInteger().f - 1, -2, true},
	{"MAX_SAFE_INTEGER+1", MaxSafeInteger().f + 1, 0, true},
	{"-8589934590", -8589934590, 2, false},
	{"0xDEADBEEF", 0xDEADBEEF, -559038737, true},
	{"4294967808", 4294967808, 512, false},
	{"-0.4", -0.4, 0, false},
	{"SmallestNonzeroFloat64", math.SmallestNonzeroFloat64, 0, false},
	{"-SmallestNonzeroFloat64", -math.SmallestNonzeroFloat64, 0, false},
	{"MaxFloat64", math.MaxFloat64, 0, false},
	{"-MaxFloat64", -math.MaxFloat64, 0, false},
	{"Largest subnormal number", float64(math.Float64frombits(0x000FFFFFFFFFFFFF)), 0, false},
	{"Smallest positive normal number", float64(math.Float64frombits(0x0010000000000000)), 0, false},
	{"Largest normal number", math.MaxFloat64, 0, false},
	{"-Largest normal number", -math.MaxFloat64, 0, false},
	{"1.0", 1.0, 1, false},
	{"-1.0", -1.0, -1, false},
	{"1e308", 1e308, 0, false},
	{"-1e308", -1e308, 0, false},
	{"math.Pi", math.Pi, 3, false},
	{"-math.Pi", -math.Pi, -3, false},
	{"math.E", math.E, 2, false},
	{"-math.E", -math.E, -2, false},
	{"0.5", 0.5, 0, false},
	{"-0.5", -0.5, 0, false},
	{"0.49999999999999994", 0.49999999999999994, 0, false},
	{"-0.49999999999999994", -0.49999999999999994, 0, false},
	{"0.5000000000000001", 0.5000000000000001, 0, false},
	{"-0.5000000000000001", -0.5000000000000001, 0, false},
	{"2^31 + 0.5", 2147483648.5, -2147483648, false},
	{"-2^31 - 0.5", -2147483648.5, -2147483648, false},
	{"2^40", 1099511627776, 0, false},
	{"-2^40", -1099511627776, 0, false},
	{"TypeFlagsNarrowable", 536624127, 536624127, true},
}

func TestToInt32(t *testing.T) {
	t.Parallel()

	for _, test := range toInt32Tests {
		input := From(test.input)

		t.Run(fmt.Sprintf("%s (%v)", test.name, input), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, input.toInt32(), test.want)
		})
	}
}

var sink int32

func BenchmarkToInt32(b *testing.B) {
	for _, test := range toInt32Tests {
		if !test.bench {
			continue
		}

		input := From(test.input)

		b.Run(fmt.Sprintf("%s (%v)", test.name, input), func(b *testing.B) {
			for range b.N {
				sink = input.toInt32()
			}
		})
	}
}

func TestBitwiseNOT(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input float64
		want  Number
	}{
		{-2147483649, From(2147483647).BitwiseNOT()},
		{-4294967296, From(0).BitwiseNOT()},
		{-2147483648, From(-2147483648).BitwiseNOT()},
		{-4294967296, From(0).BitwiseNOT()},
	}

	for _, test := range tests {
		input := From(test.input)

		t.Run(input.String(), func(t *testing.T) {
			t.Parallel()
			assertEqualNumber(t, input.BitwiseNOT(), test.want)
		})
	}
}

func TestBitwiseAND(t *testing.T) {
	t.Parallel()

	tests := []struct {
		x, y, want float64
	}{
		{0, 0, 0},
		{0, 1, 0},
		{1, 0, 0},
		{1, 1, 1},
	}

	for _, test := range tests {
		x := From(test.x)
		y := From(test.y)
		want := From(test.want)

		t.Run(fmt.Sprintf("%v & %v", x, y), func(t *testing.T) {
			t.Parallel()
			assertEqualNumber(t, x.BitwiseAND(y), want)
		})
	}
}

func TestBitwiseOR(t *testing.T) {
	t.Parallel()

	tests := []struct {
		x, y, want float64
	}{
		{0, 0, 0},
		{0, 1, 1},
		{1, 0, 1},
		{1, 1, 1},
	}

	for _, test := range tests {
		x := From(test.x)
		y := From(test.y)
		want := From(test.want)

		t.Run(fmt.Sprintf("%v | %v", x, y), func(t *testing.T) {
			t.Parallel()
			assertEqualNumber(t, x.BitwiseOR(y), want)
		})
	}
}

func TestBitwiseXOR(t *testing.T) {
	t.Parallel()

	tests := []struct {
		x, y, want float64
	}{
		{0, 0, 0},
		{0, 1, 1},
		{1, 0, 1},
		{1, 1, 0},
	}

	for _, test := range tests {
		x := From(test.x)
		y := From(test.y)
		want := From(test.want)

		t.Run(fmt.Sprintf("%v ^ %v", x, y), func(t *testing.T) {
			t.Parallel()
			assertEqualNumber(t, x.BitwiseXOR(y), want)
		})
	}
}

func TestSignedRightShift(t *testing.T) {
	t.Parallel()

	tests := []struct {
		x, y, want float64
	}{
		{1, 0, 1},
		{1, 1, 0},
		{1, 2, 0},
		{1, 31, 0},
		{1, 32, 1},
		{-4, 0, -4},
		{-4, 1, -2},
		{-4, 2, -1},
		{-4, 3, -1},
		{-4, 4, -1},
		{-4, 31, -1},
		{-4, 32, -4},
		{-4, 33, -2},
	}

	for _, test := range tests {
		x := From(test.x)
		y := From(test.y)
		want := From(test.want)

		t.Run(fmt.Sprintf("%v >> %v", x, y), func(t *testing.T) {
			t.Parallel()
			assertEqualNumber(t, x.SignedRightShift(y), want)
		})
	}
}

func TestUnsignedRightShift(t *testing.T) {
	t.Parallel()

	tests := []struct {
		x, y, want float64
	}{
		{1, 0, 1},
		{1, 1, 0},
		{1, 2, 0},
		{1, 31, 0},
		{1, 32, 1},
		{-4, 0, 4294967292},
		{-4, 1, 2147483646},
		{-4, 2, 1073741823},
		{-4, 3, 536870911},
		{-4, 4, 268435455},
		{-4, 31, 1},
		{-4, 32, 4294967292},
		{-4, 33, 2147483646},
	}

	for _, test := range tests {
		x := From(test.x)
		y := From(test.y)
		want := From(test.want)

		t.Run(fmt.Sprintf("%v >>> %v", x, y), func(t *testing.T) {
			t.Parallel()
			assertEqualNumber(t, x.UnsignedRightShift(y), want)
		})
	}
}

func TestLeftShift(t *testing.T) {
	t.Parallel()

	tests := []struct {
		x, y, want float64
	}{
		{1, 0, 1},
		{1, 1, 2},
		{1, 2, 4},
		{1, 31, -2147483648},
		{1, 32, 1},
		{-4, 0, -4},
		{-4, 1, -8},
		{-4, 2, -16},
		{-4, 3, -32},
		{-4, 31, 0},
		{-4, 32, -4},
	}

	for _, test := range tests {
		x := From(test.x)
		y := From(test.y)
		want := From(test.want)

		t.Run(fmt.Sprintf("%v << %v", x, y), func(t *testing.T) {
			t.Parallel()
			assertEqualNumber(t, x.LeftShift(y), want)
		})
	}
}

func TestRemainder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		x, y, want float64
	}{
		{NaN().f, 1, NaN().f},
		{1, NaN().f, NaN().f},
		{Inf(1).f, 1, NaN().f},
		{Inf(-1).f, 1, NaN().f},
		{123, Inf(1).f, 123},
		{123, Inf(-1).f, 123},
		{123, 0, NaN().f},
		{123, negativeZero.f, NaN().f},
		{0, 123, 0},
		{negativeZero.f, 123, negativeZero.f},
	}

	for _, test := range tests {
		x := From(test.x)
		y := From(test.y)
		want := From(test.want)

		t.Run(fmt.Sprintf("%v %% %v", x, y), func(t *testing.T) {
			t.Parallel()
			assertEqualNumber(t, x.Remainder(y), want)
		})
	}
}

func TestExponentiate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		x, y, want float64
	}{
		{2, 3, 8},
		{Inf(1).f, 3, Inf(1).f},
		{Inf(1).f, -5, 0},
		{Inf(-1).f, 3, Inf(-1).f},
		{Inf(-1).f, 4, Inf(1).f},
		{Inf(-1).f, -3, negativeZero.f},
		{Inf(-1).f, -4, 0},
		{0, 3, 0},
		{0, -10, Inf(1).f},
		{negativeZero.f, 3, negativeZero.f},
		{negativeZero.f, 4, 0},
		{negativeZero.f, -3, Inf(-1).f},
		{negativeZero.f, -4, Inf(1).f},
		{3, Inf(1).f, Inf(1).f},
		{-3, Inf(1).f, Inf(1).f},
		{3, Inf(-1).f, 0},
		{-3, Inf(-1).f, 0},
		{NaN().f, 3, NaN().f},
		{1, Inf(1).f, NaN().f},
		{1, Inf(-1).f, NaN().f},
		{-1, Inf(1).f, NaN().f},
		{-1, Inf(-1).f, NaN().f},
		{1, NaN().f, NaN().f},
	}

	for _, test := range tests {
		x := From(test.x)
		y := From(test.y)
		want := From(test.want)

		t.Run(fmt.Sprintf("%v ** %v", x, y), func(t *testing.T) {
			t.Parallel()
			assertEqualNumber(t, x.Exponentiate(y), want)
		})
	}
}
