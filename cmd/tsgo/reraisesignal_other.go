//go:build !unix

package main

import "os"

// reRaiseSignal is a no-op on platforms that cannot re-deliver a termination
// signal to the current process (notably Windows, where os.Process.Signal does
// not implement Interrupt). It always returns 0 so the caller falls back to
// exiting with a numeric status.
func reRaiseSignal(sig os.Signal) int {
	return 0
}
