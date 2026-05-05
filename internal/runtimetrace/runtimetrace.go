// Package runtimetrace wires Go's runtime/trace and flight recorder into a
// process based on environment variables. It is intended for use by long- and
// short-lived tsgo entrypoints; callers Start a Session at process startup and
// must Stop it before exit so traces are flushed.
//
// Environment variables:
//
//	TS_GO_RUNTIME_TRACE
//	    If set to a non-empty path, write a Go execution trace covering the
//	    whole process lifetime to that path. Inspect with `go tool trace`.
//
//	TS_GO_RUNTIME_TRACE_FLIGHT
//	    If set to a non-empty path (file or directory), start a flight
//	    recorder and write a snapshot of the most recent execution-trace
//	    window on shutdown (and on SIGUSR1 on unix). When the value is an
//	    existing directory the snapshot is written into it with a generated
//	    filename.
//
//	TS_GO_RUNTIME_TRACE_FLIGHT_MIN_AGE
//	    Optional time.ParseDuration value used as MinAge for the flight
//	    recorder.
//
//	TS_GO_RUNTIME_TRACE_FLIGHT_MAX_BYTES
//	    Optional unsigned integer used as MaxBytes for the flight recorder.
//
//	TS_GO_RUNTIME_TRACE_DETAIL
//	    If set to a truthy value (1, true, yes, on), instrumentation may emit
//	    additional log events that include potentially sensitive information
//	    such as file paths or identifier names. Defaults to off so that
//	    captured traces are safe to share. Use LogUnsafe* to
//	    gate such events.
//
//nolint:depguard
package runtimetrace

import (
	"fmt"
	"io"
	"os"
	"runtime/trace"
	"strconv"
	"sync"
	"sync/atomic"
)

const (
	envRuntimeTrace              = "TS_GO_RUNTIME_TRACE"
	envRuntimeTraceFlight        = "TS_GO_RUNTIME_TRACE_FLIGHT"
	envRuntimeTraceFlightMinAge  = "TS_GO_RUNTIME_TRACE_FLIGHT_MIN_AGE"
	envRuntimeTraceFlightMaxByte = "TS_GO_RUNTIME_TRACE_FLIGHT_MAX_BYTES"
	envRuntimeTraceDetail        = "TS_GO_RUNTIME_TRACE_DETAIL"
)

// Session represents an active runtime tracing setup (regular trace and/or
// flight recorder) configured from the environment.
type Session struct {
	logWriter io.Writer

	traceFile *os.File
	tracePath string

	flight     *FlightRecorder
	flightPath string

	snapshotOnce sync.Once
	stopUnix     func()
}

// Start inspects the environment and, if requested, starts runtime tracing.
// The returned Session must be Stop()ped before the process exits to ensure
// traces are flushed. Errors are logged to logWriter but not fatal.
func Start(logWriter io.Writer) *Session {
	s := &Session{logWriter: logWriter}

	switch os.Getenv(envRuntimeTraceDetail) {
	case "", "0", "false", "no", "off":
		// detail logging disabled
	default:
		unsafeLogging.Store(true)
	}

	if path := os.Getenv(envRuntimeTrace); path != "" {
		if err := s.startTrace(path); err != nil {
			fmt.Fprintf(logWriter, "runtimetrace: failed to start runtime trace: %v\n", err)
		}
	}

	if path := os.Getenv(envRuntimeTraceFlight); path != "" {
		if err := s.startFlight(path); err != nil {
			fmt.Fprintf(logWriter, "runtimetrace: failed to start flight recorder: %v\n", err)
		} else {
			s.stopUnix = installFlightSignalHandler(s)
		}
	}

	return s
}

func (s *Session) startTrace(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	if err := trace.Start(f); err != nil {
		f.Close()
		return err
	}
	s.traceFile = f
	s.tracePath = path
	fmt.Fprintf(s.logWriter, "runtimetrace: writing runtime trace to %s\n", path)
	return nil
}

func (s *Session) startFlight(path string) error {
	cfg := FlightRecorderConfig{}
	if err := cfg.SetMinAgeString(os.Getenv(envRuntimeTraceFlightMinAge)); err != nil {
		return fmt.Errorf("%s: %w", envRuntimeTraceFlightMinAge, err)
	}
	if v := os.Getenv(envRuntimeTraceFlightMaxByte); v != "" {
		n, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return fmt.Errorf("parse %s=%q: %w", envRuntimeTraceFlightMaxByte, v, err)
		}
		cfg.MaxBytes = n
	}
	fr := &FlightRecorder{}
	if err := fr.Start(cfg); err != nil {
		return err
	}
	s.flight = fr
	s.flightPath = path
	fmt.Fprintf(s.logWriter, "runtimetrace: flight recorder armed; snapshot on shutdown -> %s\n", path)
	return nil
}

// snapshotFlight writes the current flight-recorder window to the configured
// path. Only the first call writes a file. If the configured path is an
// existing directory, a generated filename is used inside it.
func (s *Session) snapshotFlight() {
	if s == nil || s.flight == nil {
		return
	}
	s.snapshotOnce.Do(func() {
		if !s.flight.Enabled() {
			return
		}
		if info, err := os.Stat(s.flightPath); err == nil && info.IsDir() {
			snapshot, err := s.flight.Snapshot(s.flightPath)
			if err != nil {
				fmt.Fprintf(s.logWriter, "runtimetrace: failed to write flight snapshot: %v\n", err)
				return
			}
			fmt.Fprintf(s.logWriter, "runtimetrace: wrote flight recorder snapshot to %s\n", snapshot)
			return
		}
		if err := s.flight.SnapshotTo(s.flightPath); err != nil {
			fmt.Fprintf(s.logWriter, "runtimetrace: failed to write flight snapshot: %v\n", err)
			return
		}
		fmt.Fprintf(s.logWriter, "runtimetrace: wrote flight recorder snapshot to %s\n", s.flightPath)
	})
}

// Stop flushes any active runtime trace and writes a final flight-recorder
// snapshot if one has not yet been written. Safe to call on a nil receiver.
func (s *Session) Stop() {
	if s == nil {
		return
	}
	if s.stopUnix != nil {
		s.stopUnix()
		s.stopUnix = nil
	}
	if s.flight != nil {
		s.snapshotFlight()
		_ = s.flight.Stop()
		s.flight = nil
	}
	if s.traceFile != nil {
		trace.Stop()
		if err := s.traceFile.Close(); err != nil {
			fmt.Fprintf(s.logWriter, "runtimetrace: failed to close runtime trace %s: %v\n", s.tracePath, err)
		} else {
			fmt.Fprintf(s.logWriter, "runtimetrace: wrote runtime trace to %s\n", s.tracePath)
		}
		s.traceFile = nil
	}
}

// flightSnapshotInFlight guards against concurrent SIGUSR1 deliveries racing
// with shutdown. Used by the unix signal handler.
var flightSnapshotInFlight atomic.Bool
