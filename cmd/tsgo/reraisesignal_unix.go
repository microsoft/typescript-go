//go:build unix

package main

import (
	"os"
	"os/signal"
	"syscall"
)

// reRaiseSignal resets the default disposition for sig and re-delivers it to
// this process, so we terminate via the signal itself (yielding the
// conventional 128+signo exit code and the terminal reset an unhandled signal
// produces). It returns the number of the signal that was re-raised, or 0 if
// sig is not a signal this platform can re-deliver, in which case the caller
// should fall back to exiting with a numeric status.
func reRaiseSignal(sig os.Signal) int {
	s, ok := sig.(syscall.Signal)
	if !ok {
		return 0
	}
	signal.Reset(s)
	if proc, err := os.FindProcess(os.Getpid()); err == nil {
		_ = proc.Signal(s)
	}
	return int(s)
}
