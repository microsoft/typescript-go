//go:build !unix

package runtimetrace

// installFlightSignalHandler is a no-op on platforms without SIGUSR1.
// The flight-recorder snapshot will still be written on Session.Stop.
func installFlightSignalHandler(s *Session) func() {
	return nil
}
