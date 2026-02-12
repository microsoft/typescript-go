package api

import (
	"io"
	"net"
	"os"
	"sync"
)

// Transport is an interface for accepting a single connection from an API client.
type Transport interface {
	// Accept waits for and returns the connection.
	Accept() (io.ReadWriteCloser, error)
	// Close releases transport resources.
	Close() error
}

// pipeTransport accepts a connection on a Unix domain socket or Windows named pipe.
type pipeTransport struct {
	listener net.Listener
	once     sync.Once
	closeErr error
}

// NewPipeTransport creates a new transport listening on the given path.
// On Unix, this creates a Unix domain socket. On Windows, this creates a named pipe.
func NewPipeTransport(path string) (Transport, error) {
	listener, err := newPipeListener(path)
	if err != nil {
		return nil, err
	}
	return &pipeTransport{listener: listener}, nil
}

func (t *pipeTransport) close() error {
	t.once.Do(func() {
		t.closeErr = t.listener.Close()
	})
	return t.closeErr
}

func (t *pipeTransport) Accept() (io.ReadWriteCloser, error) {
	conn, err := t.listener.Accept()
	if closeErr := t.close(); err != nil {
		return nil, err
	} else if closeErr != nil {
		conn.Close()
		return nil, closeErr
	}
	return conn, nil
}

func (t *pipeTransport) Close() error {
	return t.close()
}

// stdioTransport wraps stdin/stdout as a single-connection transport.
type stdioTransport struct {
	used bool
}

func newStdioTransport() *stdioTransport {
	return &stdioTransport{}
}

func (t *stdioTransport) Accept() (io.ReadWriteCloser, error) {
	if t.used {
		return nil, io.EOF
	}
	t.used = true
	return &stdioConn{
		reader: os.Stdin,  //nolint:forbidigo
		writer: os.Stdout, //nolint:forbidigo
	}, nil
}

func (t *stdioTransport) Close() error {
	return nil
}

type stdioConn struct {
	reader *os.File //nolint:forbidigo
	writer *os.File //nolint:forbidigo
}

func (c *stdioConn) Read(p []byte) (int, error) {
	return c.reader.Read(p) //nolint:forbidigo
}

func (c *stdioConn) Write(p []byte) (int, error) {
	return c.writer.Write(p) //nolint:forbidigo
}

func (c *stdioConn) Close() error {
	return c.writer.Close() //nolint:forbidigo
}
