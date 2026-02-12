//go:build !windows

package api

import (
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"syscall"
)

// newPipeListener creates a Unix domain socket listener.
func newPipeListener(path string) (net.Listener, error) {
	// Remove any existing socket file
	_ = os.Remove(path) //nolint:forbidigo
	return net.Listen("unix", path)
}

// GeneratePipePath returns a platform-appropriate pipe path for the given name.
func GeneratePipePath(name string) string {
	//nolint:forbidigo
	return path.Join(os.TempDir(), name)
}

// newServerTransport creates a FIFO transport for the API server on Unix.
// It creates two FIFOs at prefix.in and prefix.out.
func newServerTransport(prefix string) (Transport, error) {
	inPath := prefix + ".in"
	outPath := prefix + ".out"

	// Create the FIFOs. The parent process will open them.
	if err := syscall.Mkfifo(inPath, 0o600); err != nil {
		return nil, fmt.Errorf("failed to create FIFO %s: %w", inPath, err)
	}
	if err := syscall.Mkfifo(outPath, 0o600); err != nil {
		_ = os.Remove(inPath) //nolint:forbidigo
		return nil, fmt.Errorf("failed to create FIFO %s: %w", outPath, err)
	}

	return &FIFOTransport{prefix: prefix}, nil
}

// FIFOTransport uses two POSIX FIFOs for communication.
// The server creates FIFOs at prefix.in and prefix.out; the client opens them.
type FIFOTransport struct {
	prefix string
	used   bool
}

// Accept opens the FIFOs and returns a combined ReadWriteCloser.
// The open order (.out write first, .in read second) must match the
// parent's open order on opposite ends to avoid deadlock.
func (t *FIFOTransport) Accept() (io.ReadWriteCloser, error) {
	if t.used {
		return nil, io.EOF
	}
	t.used = true

	// Open write end first — blocks until parent opens .out for reading
	outFile, err := os.OpenFile(t.prefix+".out", os.O_WRONLY, 0) //nolint:forbidigo
	if err != nil {
		return nil, fmt.Errorf("failed to open FIFO %s.out for writing: %w", t.prefix, err)
	}

	// Open read end — blocks until parent opens .in for writing
	inFile, err := os.OpenFile(t.prefix+".in", os.O_RDONLY, 0) //nolint:forbidigo
	if err != nil {
		outFile.Close() //nolint:forbidigo
		return nil, fmt.Errorf("failed to open FIFO %s.in for reading: %w", t.prefix, err)
	}

	return &fifoConn{reader: inFile, writer: outFile}, nil
}

// Close removes the FIFO files.
func (t *FIFOTransport) Close() error {
	_ = os.Remove(t.prefix + ".in")  //nolint:forbidigo
	_ = os.Remove(t.prefix + ".out") //nolint:forbidigo
	return nil
}

type fifoConn struct {
	reader *os.File //nolint:forbidigo
	writer *os.File //nolint:forbidigo
}

func (c *fifoConn) Read(p []byte) (int, error) {
	return c.reader.Read(p) //nolint:forbidigo
}

func (c *fifoConn) Write(p []byte) (int, error) {
	return c.writer.Write(p) //nolint:forbidigo
}

func (c *fifoConn) Close() error {
	err1 := c.reader.Close() //nolint:forbidigo
	err2 := c.writer.Close() //nolint:forbidigo
	if err1 != nil {
		return err1
	}
	return err2
}
