package jsnum

import (
	"fmt"
	"math"
	"testing"

	"gotest.tools/v3/assert"
)

var toInt32Tests = []struct {
	name  string
	input Number
	want  int32
	bench bool
}{
	{"0.0", 0, 0, true},
	{"-0.0", Number(negativeZero), 0, false},
	{"NaN", NaN(), 0, true},
	{"+Inf", Inf(1), 0, true},
	{"-Inf", Inf(-1), 0, true},
	{"MaxInt32", Number(math.MaxInt32), math.MaxInt32, false},
	{"MaxInt32+1", Number(int64(math.MaxInt32) + 1), math.MinInt32, true},
	{"MinInt32", Number(math.MinInt32), math.MinInt32, false},
	{"MinInt32-1", Number(int64(math.MinInt32) - 1), math.MaxInt32, true},
	{"MIN_SAFE_INTEGER", MinSafeInteger, 1, false},
	{"MIN_SAFE_INTEGER-1", MinSafeInteger - 1, 0, false},
	{"MIN_SAFE_INTEGER+1", MinSafeInteger + 1, 2, false},
	{"MAX_SAFE_INTEGER", MaxSafeInteger, -1, true},
	{"MAX_SAFE_INTEGER-1", MaxSafeInteger - 1, -2, true},
	{"MAX_SAFE_INTEGER+1", MaxSafeInteger + 1, 0, true},
	{"-8589934590", -8589934590, 2, false},
	{"0xDEADBEEF", 0xDEADBEEF, -559038737, true},
	{"4294967808", 4294967808, 512, false},
	{"-0.4", -0.4, 0, false},
	{"SmallestNonzeroFloat64", math.SmallestNonzeroFloat64, 0, false},
	{"-SmallestNonzeroFloat64", -math.SmallestNonzeroFloat64, 0, false},
	{"MaxFloat64", math.MaxFloat64, 0, false},
	{"-MaxFloat64", -math.MaxFloat64, 0, false},
	{"Largest subnormal number", Number(math.Float64frombits(0x000FFFFFFFFFFFFF)), 0, false},
	{"Smallest positive normal number", Number(math.Float64frombits(0x0010000000000000)), 0, false},
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
		t.Run(fmt.Sprintf("%s (%v)", test.name, test.input), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, toInt32(test.input), test.want)
		})
	}
}

var sink int32

func BenchmarkToInt32(b *testing.B) {
	for _, test := range toInt32Tests {
		if !test.bench {
			continue
		}

		b.Run(fmt.Sprintf("%s (%v)", test.name, test.input), func(b *testing.B) {
			for range b.N {
				sink = toInt32(test.input)
			}
		})
	}
}

func TestBitwiseNOT(t *testing.T) {
	t.Parallel()

	assert.Equal(t, Number(-2147483649).BitwiseNOT(), Number(2147483647).BitwiseNOT())
	assert.Equal(t, Number(-4294967296).BitwiseNOT(), Number(0).BitwiseNOT())
	assert.Equal(t, Number(-2147483648).BitwiseNOT(), Number(-2147483648).BitwiseNOT())
	assert.Equal(t, Number(-4294967296).BitwiseNOT(), Number(0).BitwiseNOT())
}

func TestBitwiseAND(t *testing.T) {
	t.Parallel()

	assert.Equal(t, Number(1).BitwiseAND(0), Number(0.0))
	assert.Equal(t, Number(1).BitwiseAND(1), Number(1.0))
}

func TestBitwiseOR(t *testing.T) {
	t.Parallel()

	assert.Equal(t, Number(1).BitwiseOR(0), Number(1.0))
	assert.Equal(t, Number(1).BitwiseOR(1), Number(1.0))
}

func TestBitwiseXOR(t *testing.T) {
	t.Parallel()

	assert.Equal(t, Number(1).BitwiseXOR(0), Number(1.0))
	assert.Equal(t, Number(1).BitwiseXOR(1), Number(0.0))
}

func TestSignedRightShift(t *testing.T) {
	t.Parallel()

	assert.Equal(t, Number(1).SignedRightShift(0), Number(1.0))
	assert.Equal(t, Number(1).SignedRightShift(1), Number(0.0))
	assert.Equal(t, Number(1).SignedRightShift(2), Number(0.0))
	assert.Equal(t, Number(1).SignedRightShift(31), Number(0.0))
	assert.Equal(t, Number(1).SignedRightShift(32), Number(1.0))

	assert.Equal(t, Number(-4).SignedRightShift(0), Number(-4.0))
	assert.Equal(t, Number(-4).SignedRightShift(1), Number(-2.0))
	assert.Equal(t, Number(-4).SignedRightShift(2), Number(-1.0))
	assert.Equal(t, Number(-4).SignedRightShift(3), Number(-1.0))
	assert.Equal(t, Number(-4).SignedRightShift(4), Number(-1.0))
	assert.Equal(t, Number(-4).SignedRightShift(31), Number(-1.0))
	assert.Equal(t, Number(-4).SignedRightShift(32), Number(-4.0))
	assert.Equal(t, Number(-4).SignedRightShift(33), Number(-2.0))
}

func TestUnsignedRightShift(t *testing.T) {
	t.Parallel()

	assert.Equal(t, Number(1).UnsignedRightShift(0), Number(1.0))
	assert.Equal(t, Number(1).UnsignedRightShift(1), Number(0.0))
	assert.Equal(t, Number(1).UnsignedRightShift(2), Number(0.0))
	assert.Equal(t, Number(1).UnsignedRightShift(31), Number(0.0))
	assert.Equal(t, Number(1).UnsignedRightShift(32), Number(1.0))

	assert.Equal(t, Number(-4).UnsignedRightShift(0), Number(4294967292.0))
	assert.Equal(t, Number(-4).UnsignedRightShift(1), Number(2147483646.0))
	assert.Equal(t, Number(-4).UnsignedRightShift(2), Number(1073741823.0))
	assert.Equal(t, Number(-4).UnsignedRightShift(3), Number(536870911.0))
	assert.Equal(t, Number(-4).UnsignedRightShift(4), Number(268435455.0))
	assert.Equal(t, Number(-4).UnsignedRightShift(31), Number(1.0))
	assert.Equal(t, Number(-4).UnsignedRightShift(32), Number(4294967292.0))
	assert.Equal(t, Number(-4).UnsignedRightShift(33), Number(2147483646.0))
}

func TestLeftShift(t *testing.T) {
	t.Parallel()

	assert.Equal(t, Number(1).LeftShift(0), Number(1.0))
	assert.Equal(t, Number(1).LeftShift(1), Number(2.0))
	assert.Equal(t, Number(1).LeftShift(2), Number(4.0))
	assert.Equal(t, Number(1).LeftShift(31), Number(-2147483648.0))
	assert.Equal(t, Number(1).LeftShift(32), Number(1.0))

	assert.Equal(t, Number(-4).LeftShift(0), Number(-4.0))
	assert.Equal(t, Number(-4).LeftShift(1), Number(-8.0))
	assert.Equal(t, Number(-4).LeftShift(2), Number(-16.0))
	assert.Equal(t, Number(-4).LeftShift(3), Number(-32.0))
	assert.Equal(t, Number(-4).LeftShift(31), Number(0.0))
	assert.Equal(t, Number(-4).LeftShift(32), Number(-4.0))
}

func TestRemainder(t *testing.T) {
	t.Parallel()

	assert.Assert(t, NaN().Remainder(1).IsNaN())
	assert.Assert(t, Number(1).Remainder(NaN()).IsNaN())

	assert.Assert(t, Inf(1).Remainder(1).IsNaN())
	assert.Assert(t, Inf(-1).Remainder(1).IsNaN())

	assert.Equal(t, Number(123).Remainder(Inf(1)), Number(123.0))
	assert.Equal(t, Number(123).Remainder(Inf(-1)), Number(123.0))

	assert.Assert(t, Number(123).Remainder(0).IsNaN())
	assert.Assert(t, Number(123).Remainder(negativeZero).IsNaN())

	assert.Equal(t, Number(0).Remainder(123), Number(0.0))
	assert.Equal(t, negativeZero.Remainder(123), negativeZero)
}

func TestExponentiate(t *testing.T) {
	t.Parallel()

	assert.Equal(t, Number(2).Exponentiate(3), Number(8.0))

	assert.Equal(t, Inf(1).Exponentiate(3), Inf(1))
	assert.Equal(t, Inf(1).Exponentiate(-5), Number(0.0))

	assert.Equal(t, Inf(-1).Exponentiate(3), Inf(-1))
	assert.Equal(t, Inf(-1).Exponentiate(4), Inf(1))
	assert.Equal(t, Inf(-1).Exponentiate(-3), negativeZero)
	assert.Equal(t, Inf(-1).Exponentiate(-4), Number(0.0))

	assert.Equal(t, Number(0).Exponentiate(3), Number(0.0))
	assert.Equal(t, Number(0).Exponentiate(-10), Inf(1))

	assert.Equal(t, negativeZero.Exponentiate(3), negativeZero)
	assert.Equal(t, negativeZero.Exponentiate(4), Number(0.0))
	assert.Equal(t, negativeZero.Exponentiate(-3), Inf(-1))
	assert.Equal(t, negativeZero.Exponentiate(-4), Inf(1))

	assert.Equal(t, Number(3).Exponentiate(Inf(1)), Inf(1))
	assert.Equal(t, Number(-3).Exponentiate(Inf(1)), Number(math.Inf(1)))

	assert.Equal(t, Number(3).Exponentiate(Inf(-1)), Number(0.0))
	assert.Equal(t, Number(-3).Exponentiate(Inf(-1)), Number(0.0))

	assert.Assert(t, NaN().Exponentiate(3).IsNaN())
	assert.Assert(t, Number(1).Exponentiate(Inf(1)).IsNaN())
	assert.Assert(t, Number(1).Exponentiate(Inf(-1)).IsNaN())
	assert.Assert(t, Number(-1).Exponentiate(Inf(1)).IsNaN())
	assert.Assert(t, Number(-1).Exponentiate(Inf(-1)).IsNaN())
	assert.Assert(t, Number(1).Exponentiate(NaN()).IsNaN())
}
