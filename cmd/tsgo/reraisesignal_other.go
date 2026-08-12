//go:build !unix

package main

import "os"

// reRaiseSignal is a no-op here: these platforms cannot re-deliver a termination
// signal to the current process (on Windows, os.Process.Signal rejects Interrupt).
func reRaiseSignal(sig os.Signal) {
}
