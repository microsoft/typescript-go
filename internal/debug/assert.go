//go:build !noassert

// Functions in this file use "mid-stack inlining" to ensure the fast path
// (assertion passing) is inlined at call sites, leaving only the slow path
// (assertion failure) as an out-of-line call.
// See https://dave.cheney.net/2020/05/02/mid-stack-inlining-in-go

package debug

import (
	"cmp"
	"fmt"
)

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
		assertSlow(value, message...)
	}
	return value
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

var unexpectedNode []any = []any{"Unexpected node."}

func AssertEachNode[TElem NodeLike](nodes []TElem, test func(elem TElem) bool, message ...any) {
	if len(message) == 0 {
		message = unexpectedNode
	}
	for _, elem := range nodes {
		AssertNode(elem, test, message...)
	}
}

func AssertNode[TElem NodeLike](node TElem, test func(elem TElem) bool, message ...any) {
	var zero TElem
	if node != zero && (test == nil || test(node)) {
		return
	}
	assertNodeSlow(node, test, message...)
}

func assertNodeSlow[TElem NodeLike](node TElem, test func(elem TElem) bool, message ...any) {
	if len(message) == 0 {
		message = unexpectedNode
	}
	var zero TElem
	if node == zero {
		Fail("False expression: " + fmt.Sprint(message...))
	}
	if test != nil && !test(node) {
		Fail("False expression: " + fmt.Sprint(message...))
	}
}

func AssertNotNode[TElem NodeLike](node TElem, test func(elem TElem) bool, message ...any) {
	var zero TElem
	if node == zero || test == nil || !test(node) {
		return
	}
	assertNodeFail(message...)
}

func AssertOptionalNode[TElem NodeLike](node TElem, test func(elem TElem) bool, message ...any) {
	var zero TElem
	if node == zero || test == nil || test(node) {
		return
	}
	assertNodeFail(message...)
}

func AssertOptionalToken[TElem NodeLike](node TElem, kind int16, message ...any) {
	var zero TElem
	if node == zero || node.KindValue() == kind {
		return
	}
	assertNodeFail(message...)
}

func AssertMissingNode[TElem NodeLike](node TElem, message ...any) {
	var zero TElem
	if node == zero {
		return
	}
	assertNodeFail(message...)
}

func assertNodeFail(message ...any) {
	if len(message) == 0 {
		message = unexpectedNode
	}
	Fail("False expression: " + fmt.Sprint(message...))
}
