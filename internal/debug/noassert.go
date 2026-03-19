//go:build noassert

package debug

func Assert(value bool, message ...any)                                          {}
func AssertEach[TElem any](value []TElem, test func(TElem) bool, message ...any) {}
