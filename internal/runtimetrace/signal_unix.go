//go:build unix

package runtimetrace

import (
	"os"
	"os/signal"
	"syscall"
)

// installFlightSignalHandler arranges for SIGUSR1 to trigger a flight-recorder
// snapshot. The returned function tears the handler down.
func installFlightSignalHandler(s *Session) func() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGUSR1)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ch:
				if flightSnapshotInFlight.Swap(true) {
					continue
				}
				s.snapshotFlight()
				flightSnapshotInFlight.Store(false)
				// snapshotFlight is sync.Once-guarded, so subsequent
				// signals are no-ops; keep the loop running to drain
				// the channel cleanly.
			}
		}
	}()
	return func() {
		signal.Stop(ch)
		close(done)
	}
}
