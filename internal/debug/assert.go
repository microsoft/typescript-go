package debug

import (
	"cmp"
	"fmt"
)

// Functions in this file use "mid-stack inlining" to ensure the fast path
// (assertion passing) is inlined at call sites, leaving only the slow path
// (assertion failure) as an out-of-line call.
// See https://dave.cheney.net/2020/05/02/mid-stack-inlining-in-go

func Assert[T comparable](value T, message ...any) {
	var zero T
	if value != zero {
		return
	}
	assertSlow(value, message...)
}

func assertSlow[T comparable](value T, message ...any) {
	var prefix string
	if _, ok := any(value).(bool); ok {
		prefix = "False expression"
	} else {
		prefix = "Expected non-zero value"
	}
	var msg string
	if len(message) > 0 {
		msg = prefix + ": " + fmt.Sprint(message...)
	} else {
		msg = prefix + "."
	}
	Fail(msg)
}

func AssertZero[T comparable](value T, message ...any) {
	var zero T
	if value == zero {
		return
	}
	assertZeroSlow(message...)
}

func assertZeroSlow(message ...any) {
	var msg string
	if len(message) > 0 {
		msg = "Expected zero value: " + fmt.Sprint(message...)
	} else {
		msg = "Expected zero value."
	}
	Fail(msg)
}

func AssertEqual[T comparable](a T, b T, message ...any) {
	if a == b {
		return
	}
	assertEqualSlow(a, b, message...)
}

func assertEqualSlow[T comparable](a T, b T, message ...any) {
	var msg string
	if len(message) > 0 {
		msg = fmt.Sprint(message...)
	}
	Fail(fmt.Sprintf("Expected %v == %v. %s", a, b, msg))
}

func AssertLessThan[T cmp.Ordered](a T, b T, message ...any) {
	if a < b {
		return
	}
	assertLessThanSlow(a, b, message...)
}

func assertLessThanSlow[T cmp.Ordered](a T, b T, message ...any) {
	var msg string
	if len(message) > 0 {
		msg = fmt.Sprint(message...)
	}
	Fail(fmt.Sprintf("Expected %v < %v. %s", a, b, msg))
}

func AssertLessThanOrEqual[T cmp.Ordered](a T, b T, message ...any) {
	if a <= b {
		return
	}
	assertLessThanOrEqualSlow(a, b, message...)
}

func assertLessThanOrEqualSlow[T cmp.Ordered](a T, b T, message ...any) {
	var msg string
	if len(message) > 0 {
		msg = fmt.Sprint(message...)
	}
	Fail(fmt.Sprintf("Expected %v <= %v. %s", a, b, msg))
}

func AssertGreaterThan[T cmp.Ordered](a T, b T, message ...any) {
	if a > b {
		return
	}
	assertGreaterThanSlow(a, b, message...)
}

func assertGreaterThanSlow[T cmp.Ordered](a T, b T, message ...any) {
	var msg string
	if len(message) > 0 {
		msg = fmt.Sprint(message...)
	}
	Fail(fmt.Sprintf("Expected %v > %v. %s", a, b, msg))
}

func AssertGreaterThanOrEqual[T cmp.Ordered](a T, b T, message ...any) {
	if a >= b {
		return
	}
	assertGreaterThanOrEqualSlow(a, b, message...)
}

func assertGreaterThanOrEqualSlow[T cmp.Ordered](a T, b T, message ...any) {
	var msg string
	if len(message) > 0 {
		msg = fmt.Sprint(message...)
	}
	Fail(fmt.Sprintf("Expected %v >= %v. %s", a, b, msg))
}

func CheckNonZero[T comparable](value T, message ...any) T {
	var zero T
	if value == zero {
		checkNonZeroSlow(message...)
	}
	return value
}

func checkNonZeroSlow(message ...any) {
	var msg string
	if len(message) > 0 {
		msg = "Expected non-zero value: " + fmt.Sprint(message...)
	} else {
		msg = "Expected non-zero value."
	}
	Fail(msg)
}

func AssertEach[TElem any](value []TElem, test func(TElem) bool, message ...any) {
	for _, elem := range value {
		Assert(test(elem), message...)
	}
}

func CheckEach[TElem any](value []TElem, test func(TElem) bool, message ...any) []TElem {
	AssertEach(value, test, message...)
	return value
}
