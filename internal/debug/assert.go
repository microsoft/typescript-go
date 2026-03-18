//go:build !noassert

package debug

import (
	"fmt"
)

func isZero[T comparable](value T) bool {
	var zero T
	return value == zero
}

func Assert[T comparable](value T, message ...any) {
	if isZero(value) {
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
}

func AssertZero[T comparable](value T, message ...any) {
	if !isZero(value) {
		var msg string
		if len(message) > 0 {
			msg = "Expected zero value: " + fmt.Sprint(message...)
		} else {
			msg = "Expected zero value."
		}
		Fail(msg)
	}
}

func AssertEqual(a fmt.Stringer, b fmt.Stringer, message ...any) {
	if a != b {
		var msg string
		if len(message) == 0 {
			msg = ""
		} else {
			msg = fmt.Sprint(message...)
		}
		Fail(fmt.Sprintf("Expected %s == %s. %s", a.String(), b.String(), msg))
	}
}

func AssertLessThan(a int, b int, message ...any) {
	if a >= b {
		var msg string
		if len(message) == 0 {
			msg = ""
		} else {
			msg = fmt.Sprint(message...)
		}
		Fail(fmt.Sprintf("Expected %d < %d. %s", a, b, msg))
	}
}

func AssertLessThanOrEqual(a int, b int, message ...any) {
	if a > b {
		var msg string
		if len(message) == 0 {
			msg = ""
		} else {
			msg = fmt.Sprint(message...)
		}
		Fail(fmt.Sprintf("Expected %d <= %d. %s", a, b, msg))
	}
}

func AssertGreaterThan(a int, b int, message ...any) {
	if a <= b {
		var msg string
		if len(message) == 0 {
			msg = ""
		} else {
			msg = fmt.Sprint(message...)
		}
		Fail(fmt.Sprintf("Expected %d > %d. %s", a, b, msg))
	}
}

func AssertGreaterThanOrEqual(a int, b int, message ...any) {
	if a < b {
		var msg string
		if len(message) == 0 {
			msg = ""
		} else {
			msg = fmt.Sprint(message...)
		}
		Fail(fmt.Sprintf("Expected %d >= %d. %s", a, b, msg))
	}
}

func CheckNonZero[T comparable](value T, message ...any) T {
	Assert(value, message...)
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

func AssertEachNode[TElem comparable](nodes []TElem, test func(elem TElem) bool, message ...any) {
	if len(message) == 0 {
		message = unexpectedNode
	}
	for _, elem := range nodes {
		AssertNode(elem, test, message...)
	}
}

func AssertNode[TElem comparable](node TElem, test func(elem TElem) bool, message ...any) {
	if len(message) == 0 {
		message = unexpectedNode
	}
	Assert(node, message...)
	if test != nil {
		Assert(test(node), message...)
	}
}

func AssertNotNode[TElem comparable](node TElem, test func(elem TElem) bool, message ...any) {
	if isZero(node) {
		return
	}
	if test == nil {
		return
	}
	if len(message) == 0 {
		message = unexpectedNode
	}
	Assert(!test(node), message...)
}

func AssertOptionalNode[TElem comparable](node TElem, test func(elem TElem) bool, message ...any) {
	if isZero(node) {
		return
	}
	if test == nil {
		return
	}
	if len(message) == 0 {
		message = unexpectedNode
	}
	Assert(test(node), message...)
}

func AssertOptionalToken[TElem NodeLike](node TElem, kind int16, message ...any) {
	if isZero(node) {
		return
	}
	if len(message) == 0 {
		message = unexpectedNode
	}
	Assert(node.KindValue() == kind, message...)
}

func AssertMissingNode[TElem comparable](node TElem, message ...any) {
	if len(message) == 0 {
		message = unexpectedNode
	}
	Assert(isZero(node), message...)
}
