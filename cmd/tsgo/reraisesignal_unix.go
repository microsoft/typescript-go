//go:build unix

package main

import (
	"os"
	"os/signal"
	"syscall"
)

// reRaiseSignal restores the default disposition for sig and re-delivers it to this
// process, terminating us via the signal itself. It returns only if that fails.
func reRaiseSignal(sig os.Signal) {
	s, ok := sig.(syscall.Signal)
	if !ok {
		return
	}
	signal.Reset(s)
	if proc, err := os.FindProcess(os.Getpid()); err == nil {
		_ = proc.Signal(s)
	}
}
