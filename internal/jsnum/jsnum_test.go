package jsnum

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
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
	{"Largest subnormal number", numberFromBits(0x000FFFFFFFFFFFFF), 0, false},
	{"Smallest positive normal number", numberFromBits(0x0010000000000000), 0, false},
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
		t.Run(fmt.Sprintf("%s (%v)", test.name, float64(test.input)), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.input.toInt32(), test.want)
		})
	}
}

var sink int32

func BenchmarkToInt32(b *testing.B) {
	for _, test := range toInt32Tests {
		if !test.bench {
			continue
		}

		b.Run(fmt.Sprintf("%s (%v)", test.name, float64(test.input)), func(b *testing.B) {
			for range b.N {
				sink = test.input.toInt32()
			}
		})
	}
}

func TestBitwiseNOT(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input, want Number
	}{
		{-2147483649, Number(2147483647).BitwiseNOT()},
		{-4294967296, Number(0).BitwiseNOT()},
		{-2147483648, Number(-2147483648).BitwiseNOT()},
		{-4294967296, Number(0).BitwiseNOT()},
	}

	for _, test := range tests {
		t.Run(test.input.String(), func(t *testing.T) {
			t.Parallel()
			assertEqualNumber(t, test.input.BitwiseNOT(), test.want)
		})
	}
}

func TestBitwiseAND(t *testing.T) {
	t.Parallel()

	tests := []struct {
		x, y, want Number
	}{
		{0, 0, 0},
		{0, 1, 0},
		{1, 0, 0},
		{1, 1, 1},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v & %v", test.x, test.y), func(t *testing.T) {
			t.Parallel()
			assertEqualNumber(t, test.x.BitwiseAND(test.y), test.want)
		})
	}
}

func TestBitwiseOR(t *testing.T) {
	t.Parallel()

	tests := []struct {
		x, y, want Number
	}{
		{0, 0, 0},
		{0, 1, 1},
		{1, 0, 1},
		{1, 1, 1},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v | %v", test.x, test.y), func(t *testing.T) {
			t.Parallel()
			assertEqualNumber(t, test.x.BitwiseOR(test.y), test.want)
		})
	}
}

func TestBitwiseXOR(t *testing.T) {
	t.Parallel()

	tests := []struct {
		x, y, want Number
	}{
		{0, 0, 0},
		{0, 1, 1},
		{1, 0, 1},
		{1, 1, 0},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v ^ %v", test.x, test.y), func(t *testing.T) {
			t.Parallel()
			assertEqualNumber(t, test.x.BitwiseXOR(test.y), test.want)
		})
	}
}

func TestSignedRightShift(t *testing.T) {
	t.Parallel()

	tests := []struct {
		x, y, want Number
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
		t.Run(fmt.Sprintf("%v >> %v", test.x, test.y), func(t *testing.T) {
			t.Parallel()
			assertEqualNumber(t, test.x.SignedRightShift(test.y), test.want)
		})
	}
}

func TestUnsignedRightShift(t *testing.T) {
	t.Parallel()

	tests := []struct {
		x, y, want Number
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
		t.Run(fmt.Sprintf("%v >>> %v", test.x, test.y), func(t *testing.T) {
			t.Parallel()
			assertEqualNumber(t, test.x.UnsignedRightShift(test.y), test.want)
		})
	}
}

func TestLeftShift(t *testing.T) {
	t.Parallel()

	tests := []struct {
		x, y, want Number
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
		t.Run(fmt.Sprintf("%v << %v", test.x, test.y), func(t *testing.T) {
			t.Parallel()
			assertEqualNumber(t, test.x.LeftShift(test.y), test.want)
		})
	}
}

func TestRemainder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		x, y, want Number
	}{
		{NaN(), 1, NaN()},
		{1, NaN(), NaN()},
		{Inf(1), 1, NaN()},
		{Inf(-1), 1, NaN()},
		{123, Inf(1), 123},
		{123, Inf(-1), 123},
		{123, 0, NaN()},
		{123, negativeZero, NaN()},
		{0, 123, 0},
		{negativeZero, 123, negativeZero},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v %% %v", test.x, test.y), func(t *testing.T) {
			t.Parallel()
			assertEqualNumber(t, test.x.Remainder(test.y), test.want)
		})
	}
}

func TestExponentiate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		x, y, want Number
	}{
		{2, 3, 8},
		{Inf(1), 3, Inf(1)},
		{Inf(1), -5, 0},
		{Inf(-1), 3, Inf(-1)},
		{Inf(-1), 4, Inf(1)},
		{Inf(-1), -3, negativeZero},
		{Inf(-1), -4, 0},
		{0, 3, 0},
		{0, -10, Inf(1)},
		{negativeZero, 3, negativeZero},
		{negativeZero, 4, 0},
		{negativeZero, -3, Inf(-1)},
		{negativeZero, -4, Inf(1)},
		{3, Inf(1), Inf(1)},
		{-3, Inf(1), Inf(1)},
		{3, Inf(-1), 0},
		{-3, Inf(-1), 0},
		{NaN(), 3, NaN()},
		{1, Inf(1), NaN()},
		{1, Inf(-1), NaN()},
		{-1, Inf(1), NaN()},
		{-1, Inf(-1), NaN()},
		{1, NaN(), NaN()},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v ** %v", test.x, test.y), func(t *testing.T) {
			t.Parallel()
			assertEqualNumber(t, test.x.Exponentiate(test.y), test.want)
		})
	}
}

type stringTest struct {
	number Number
	str    string
}

var stringTests = slices.Concat([]stringTest{
	{NaN(), "NaN"},
	{Inf(1), "Infinity"},
	{Inf(-1), "-Infinity"},
	{0, "0"},
	{negativeZero, "0"},
	{1, "1"},
	{-1, "-1"},
	{0.3, "0.3"},
	{-0.3, "-0.3"},
	{1.5, "1.5"},
	{-1.5, "-1.5"},
	{1e308, "1e+308"},
	{-1e308, "-1e+308"},
	{math.Pi, "3.141592653589793"},
	{-math.Pi, "-3.141592653589793"},
	{MaxSafeInteger, "9007199254740991"},
	{MinSafeInteger, "-9007199254740991"},
	{numberFromBits(0x000FFFFFFFFFFFFF), "2.225073858507201e-308"},
	{numberFromBits(0x0010000000000000), "2.2250738585072014e-308"},
	{1234567.8, "1234567.8"},
	{19686109595169230000, "19686109595169230000"},
	{123.456, "123.456"},
	{-123.456, "-123.456"},
}, ryuTests)

func TestString(t *testing.T) {
	t.Parallel()

	for _, test := range stringTests {
		fInput := float64(test.number)

		t.Run(fmt.Sprintf("%v", fInput), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.number.String(), test.str)
		})
	}
}

var fromStringTests = []stringTest{
	{NaN(), "    NaN"},
	{Inf(1), "Infinity    "},
	{Inf(-1), "    -Infinity"},
	{1, "1."},
	{1, "1.0   "},
	{NaN(), "whoops"},
}

func TestFromString(t *testing.T) {
	t.Parallel()

	for _, test := range stringTests {
		t.Run(test.str, func(t *testing.T) {
			t.Parallel()
			assertEqualNumber(t, FromString(test.str), test.number)
			assertEqualNumber(t, FromString(test.str+" "), test.number)
			assertEqualNumber(t, FromString(" "+test.str), test.number)
		})
	}
}

func getNodeExe(t testing.TB) string {
	t.Helper()

	const exeName = "node"
	exe, err := exec.LookPath(exeName)
	if err != nil {
		t.Skipf("%s not found: %v", exeName, err)
	}
	return exe
}

func TestNodeString(t *testing.T) {
	t.Parallel()

	exe := getNodeExe(t)

	t.Run("stringTests", func(t *testing.T) {
		t.Parallel()

		// These tests should roundtrip both ways.
		stringTestsResults := getRealStringResults(t, exe, stringTests)
		for i, test := range stringTests {
			t.Run(fmt.Sprintf("%v", float64(test.number)), func(t *testing.T) {
				t.Parallel()
				assertEqualNumber(t, stringTestsResults[i].number, test.number)
				assert.Equal(t, stringTestsResults[i].str, test.str)
			})
		}
	})

	t.Run("fromStringTests", func(t *testing.T) {
		t.Parallel()

		// These tests should convert the string to the same number.
		fromStringTestsResults := getRealStringResults(t, exe, fromStringTests)
		for i, test := range fromStringTests {
			t.Run(fmt.Sprintf("fromString %q", test.str), func(t *testing.T) {
				t.Parallel()
				assertEqualNumber(t, fromStringTestsResults[i].number, test.number)
			})
		}
	})
}

func FuzzNodeString(f *testing.F) {
	exe := getNodeExe(f)

	for _, test := range stringTests {
		f.Add(float64(test.number))
	}
	for _, test := range fromStringTests {
		f.Add(float64(test.number))
	}

	f.Fuzz(func(t *testing.T, f float64) {
		n := Number(f)

		results := getRealStringResults(t, exe, []stringTest{{number: n}})
		assert.Equal(t, len(results), 1)
		assert.Equal(t, results[0].str, n.String())
	})
}

func getRealStringResults(t testing.TB, exe string, tests []stringTest) []stringTest {
	t.Helper()
	tmpdir := t.TempDir()

	type data struct {
		Bits [2]uint32 `json:"bits"`
		Str  string    `json:"str"`
	}

	inputData := make([]data, len(tests))
	for i, test := range tests {
		inputData[i] = data{
			Bits: numberToUint32Array(test.number),
			Str:  test.str,
		}
	}

	jsonInput, err := json.Marshal(inputData)
	assert.NilError(t, err)

	jsonInputPath := filepath.Join(tmpdir, "input.json")
	err = os.WriteFile(jsonInputPath, jsonInput, 0o644)
	assert.NilError(t, err)

	script := `
		const fs = require('fs');

		function fromBits(bits) {
			const buffer = new ArrayBuffer(8);
			(new Uint32Array(buffer))[0] = bits[0];
			(new Uint32Array(buffer))[1] = bits[1];
			return new Float64Array(buffer)[0];
		}

		function toBits(number) {
			const buffer = new ArrayBuffer(8);
			(new Float64Array(buffer))[0] = number;
			return [(new Uint32Array(buffer))[0], (new Uint32Array(buffer))[1]];
		}

		const input = JSON.parse(fs.readFileSync(process.argv[2], 'utf8'));

		const output = input.map((input) => ({
			str: ""+fromBits(input.bits),
			bits: toBits(+input.str),	
		}));

		process.stdout.write(JSON.stringify(output));
	`

	scriptPath := filepath.Join(tmpdir, "script.cjs")
	err = os.WriteFile(scriptPath, []byte(script), 0o644)
	assert.NilError(t, err)

	execCmd := exec.Command(exe, scriptPath, jsonInputPath)

	stdout, err := execCmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			t.Fatalf("failed to execute: %v\n%s", err, exitErr.Stderr)
		} else {
			t.Fatalf("failed to execute: %v", err)
		}
	}

	var outputData []data

	err = json.Unmarshal(stdout, &outputData)
	assert.NilError(t, err)

	assert.Equal(t, len(outputData), len(tests))

	output := make([]stringTest, len(tests))
	for i, outputDatum := range outputData {
		output[i] = stringTest{
			number: uint32ArrayToNumber(outputDatum.Bits),
			str:    outputDatum.Str,
		}
	}

	return output
}

func numberFromBits(b uint64) Number {
	return Number(math.Float64frombits(b))
}

func numberToBits(n Number) uint64 {
	return math.Float64bits(float64(n))
}

func numberToUint32Array(n Number) [2]uint32 {
	bits := numberToBits(n)
	return [2]uint32{uint32(bits), uint32(bits >> 32)}
}

func uint32ArrayToNumber(a [2]uint32) Number {
	bits := uint64(a[0]) | uint64(a[1])<<32
	return numberFromBits(bits)
}
