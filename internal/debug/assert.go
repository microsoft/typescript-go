//go:build !noassert

package debug

import (
	"fmt"
)

func Assert(expression bool, message ...any) {
	if !expression {
		var msg string
		if len(message) > 0 {
			msg = "False expression: " + fmt.Sprint(message...)
		} else {
			msg = "False expression."
		}
		Fail(msg)
	}
}

func isZero[T comparable](value T) bool {
	var zero T
	return value == zero
}

func AssertNil[T comparable](value T, message ...any) {
	if !isZero(value) {
		var msg string
		if len(message) > 0 {
			msg = "Nil expression: " + fmt.Sprint(message...)
		} else {
			msg = "Nil expression."
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

func AssertIsDefined[T comparable](value T, message ...any) {
	if isZero(value) {
		var msg string
		if len(message) == 0 {
			msg = ""
		} else {
			msg = fmt.Sprint(message...)
		}
		Fail(msg)
	}
}

func CheckDefined[T comparable](value T, message ...any) T {
	AssertIsDefined(value, message...)
	return value
}

func AssertEachIsDefined[TElem comparable](value []TElem, message ...any) {
	for _, elem := range value {
		AssertIsDefined(elem, message...)
	}
}

func CheckEachIsDefined[TElem comparable](value []TElem, message ...any) []TElem {
	AssertEachIsDefined(value, message...)
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
	AssertIsDefined(node, message...)
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

func AssertOptionalToken[TElem interface {
	comparable
	KindValue() int16
}](node TElem, kind int16, message ...any) {
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
