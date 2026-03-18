//go:build noassert

package debug

import (
	"cmp"
	"fmt"
)

func Assert[T comparable](value T, message ...any)                                              {}
func AssertZero[T comparable](value T, message ...any)                                          {}
func AssertEqual(a fmt.Stringer, b fmt.Stringer, msg ...any)                                    {}
func AssertLessThan[T cmp.Ordered](a T, b T, message ...any)                                    {}
func AssertLessThanOrEqual[T cmp.Ordered](a T, b T, message ...any)                             {}
func AssertGreaterThan[T cmp.Ordered](a T, b T, message ...any)                                 {}
func AssertGreaterThanOrEqual[T cmp.Ordered](a T, b T, message ...any)                          {}
func CheckNonZero[T comparable](value T, message ...any) T                                      { return value }
func AssertEach[TElem any](value []TElem, test func(TElem) bool, message ...any)                {}
func CheckEach[TElem any](value []TElem, test func(TElem) bool, message ...any) []TElem         { return value }
func AssertEachNode[TElem NodeLike](nodes []TElem, test func(elem TElem) bool, message ...any)  {}
func AssertNode[TElem NodeLike](node TElem, test func(elem TElem) bool, message ...any)         {}
func AssertNotNode[TElem NodeLike](node TElem, test func(elem TElem) bool, message ...any)      {}
func AssertOptionalNode[TElem NodeLike](node TElem, test func(elem TElem) bool, message ...any) {}
func AssertOptionalToken[TElem NodeLike](node TElem, kind int16, message ...any) {
}
func AssertMissingNode[TElem NodeLike](node TElem, message ...any) {}
