//go:build noassert

package debug

import (
	"fmt"
)

func Assert[T comparable](value T, message ...any)                                                {}
func AssertZero[T comparable](value T, message ...any)                                            {}
func AssertEqual(a fmt.Stringer, b fmt.Stringer, msg ...any)                                      {}
func AssertLessThan(a int, b int, message ...any)                                                 {}
func AssertLessThanOrEqual(a int, b int, message ...any)                                          {}
func AssertGreaterThan(a int, b int, message ...any)                                              {}
func AssertGreaterThanOrEqual(a int, b int, message ...any)                                       {}
func CheckNonZero[T comparable](value T, message ...any) T                                        { return value }
func AssertEach[TElem any](value []TElem, test func(TElem) bool, message ...any)                  {}
func CheckEach[TElem any](value []TElem, test func(TElem) bool, message ...any) []TElem           { return value }
func AssertEachNode[TElem comparable](nodes []TElem, test func(elem TElem) bool, message ...any)  {}
func AssertNode[TElem comparable](node TElem, test func(elem TElem) bool, message ...any)         {}
func AssertNotNode[TElem comparable](node TElem, test func(elem TElem) bool, message ...any)      {}
func AssertOptionalNode[TElem comparable](node TElem, test func(elem TElem) bool, message ...any) {}
func AssertOptionalToken[TElem NodeLike](node TElem, kind int16, message ...any) {
}
func AssertMissingNode[TElem comparable](node TElem, message ...any) {}
