//go:build !noassert

package debug

import (
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
