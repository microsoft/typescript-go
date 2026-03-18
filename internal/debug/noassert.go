//go:build noassert

package debug

import (
	"cmp"
)

func Assert[T comparable](value T, message ...any)                                      {}
func AssertZero[T comparable](value T, message ...any)                                  {}
func AssertEqual[T comparable](a T, b T, msg ...any)                                    {}
func AssertLessThan[T cmp.Ordered](a T, b T, message ...any)                            {}
func AssertLessThanOrEqual[T cmp.Ordered](a T, b T, message ...any)                     {}
func AssertGreaterThan[T cmp.Ordered](a T, b T, message ...any)                         {}
func AssertGreaterThanOrEqual[T cmp.Ordered](a T, b T, message ...any)                  {}
func CheckNonZero[T comparable](value T, message ...any) T                              { return value }
func AssertEach[TElem any](value []TElem, test func(TElem) bool, message ...any)        {}
func CheckEach[TElem any](value []TElem, test func(TElem) bool, message ...any) []TElem { return value }
