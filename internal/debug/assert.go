//go:build !noassert

package debug

import (
	"fmt"
)

// Functions in this file use "mid-stack inlining" to ensure the fast path
// (assertion passing) is inlined at call sites, leaving only the slow path
// (assertion failure) as an out-of-line call.
// See https://dave.cheney.net/2020/05/02/mid-stack-inlining-in-go

func Assert(value bool, message ...any) {
	if value {
		return
	}
	assertSlow(message...)
}

func assertSlow(message ...any) {
	var msg string
	if len(message) > 0 {
		msg = "False expression: " + fmt.Sprint(message...)
	} else {
		msg = "False expression."
	}
	Fail(msg)
}
