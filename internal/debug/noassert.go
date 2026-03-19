//go:build noassert

package debug

func Assert[T comparable](value T, message ...any)                                      {}
func CheckNonZero[T comparable](value T, message ...any) T                              { return value }
func AssertEach[TElem any](value []TElem, test func(TElem) bool, message ...any)        {}
func CheckEach[TElem any](value []TElem, test func(TElem) bool, message ...any) []TElem { return value }
