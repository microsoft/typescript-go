//nolint:depguard
package runtimetrace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/trace"
	"sync"
	"time"
)

// FlightRecorderConfig configures a flight recorder.
type FlightRecorderConfig struct {
	// MinAge is the lower bound on the age of an event in the flight
	// recorder's window. If zero, the runtime default is used.
	MinAge time.Duration
	// MaxBytes is the upper bound on the size of the window in bytes.
	// If zero, the runtime default is used.
	MaxBytes uint64
}

// SetMinAgeString parses s as a Go duration and stores it in MinAge.
// An empty string is treated as "leave the default".
func (c *FlightRecorderConfig) SetMinAgeString(s string) error {
	if s == "" {
		return nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("parse minAge %q: %w", s, err)
	}
	c.MinAge = d
	return nil
}

// FlightRecorder manages a Go execution-trace flight recorder. The zero value is
// ready to use. At most one flight recorder may be active at any given time
// in a single process (a runtime/trace constraint).
type FlightRecorder struct {
	mu sync.Mutex
	fr *trace.FlightRecorder
}

// Start activates the flight recorder. Returns an error if it is already
// running or if the runtime cannot start it.
func (r *FlightRecorder) Start(cfg FlightRecorderConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.fr != nil {
		return errors.New("flight recorder already running")
	}
	fr := trace.NewFlightRecorder(trace.FlightRecorderConfig{
		MinAge:   cfg.MinAge,
		MaxBytes: cfg.MaxBytes,
	})
	if err := fr.Start(); err != nil {
		return fmt.Errorf("failed to start flight recorder: %w", err)
	}
	r.fr = fr
	return nil
}

// Snapshot writes the current flight-recorder window to a file in dir and
// returns the path. The recorder remains active afterwards.
func (r *FlightRecorder) Snapshot(dir string) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create snapshot directory: %w", err)
	}
	path := filepath.Join(dir, fmt.Sprintf("%d-%d-flightrecorder.trace", os.Getpid(), time.Now().UnixMilli()))
	if err := r.SnapshotTo(path); err != nil {
		return "", err
	}
	return path, nil
}

// SnapshotTo writes the current flight-recorder window to the given file
// path, creating the parent directory if needed. The recorder remains active
// afterwards.
func (r *FlightRecorder) SnapshotTo(path string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.fr == nil {
		return errors.New("flight recorder not running")
	}
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create snapshot directory: %w", err)
		}
	}
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create flight recorder snapshot file: %w", err)
	}
	defer out.Close()
	if _, err := r.fr.WriteTo(out); err != nil {
		os.Remove(path)
		return fmt.Errorf("failed to write flight recorder snapshot: %w", err)
	}
	return nil
}

// Stop stops the flight recorder. Safe to call when not running.
func (r *FlightRecorder) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.fr == nil {
		return nil
	}
	r.fr.Stop()
	r.fr = nil
	return nil
}

// Enabled reports whether the flight recorder is currently running.
func (r *FlightRecorder) Enabled() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.fr != nil
}
