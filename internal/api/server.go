package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

// ServerOptions configures the API server.
type ServerOptions struct {
	Cwd                string
	DefaultLibraryPath string
	// Transport specifies the transport mechanism.
	// Supported values:
	//   ""        or "stdio" — use In/Out (stdin/stdout)
	//   "pipe=<path>"       — Unix domain socket or Windows named pipe
	//   "fifo=<prefix>"     — two POSIX FIFOs at <prefix>.in / .out (Unix only)
	Transport string
	// Callbacks specifies which filesystem operations should be delegated
	// to the client (e.g., "readFile", "fileExists"). Empty means no callbacks.
	Callbacks []string
	// Async enables JSON-RPC protocol with async connection handling.
	// When false (default), uses MessagePack protocol with sync connection.
	Async bool
}

// Server runs an API session over STDIO using MessagePack protocol.
// This is the entry point for the synchronous API used by
// native TypeScript tooling integration.
type Server struct {
	options *ServerOptions
}

// NewServer creates a new API server.
func NewServer(options *ServerOptions) *Server {
	if options.Cwd == "" {
		panic("ServerOptions.Cwd is required")
	}

	return &Server{
		options: options,
	}
}

// Run starts the server and blocks until the connection closes.
func (s *Server) Run(ctx context.Context) error {
	transport, err := s.createTransport()
	if err != nil {
		return fmt.Errorf("failed to create transport: %w", err)
	}
	defer transport.Close()

	fs := bundled.WrapFS(osvfs.FS())

	// Wrap the base FS with callbackFS if callbacks are requested
	var callbackFS *callbackFS
	if len(s.options.Callbacks) > 0 {
		callbackFS = newCallbackFS(fs, s.options.Callbacks)
		fs = callbackFS
	}

	projectSession := project.NewSession(&project.SessionInit{
		BackgroundCtx: ctx,
		Logger:        nil, // TODO: Add logging support
		FS:            fs,
		Options: &project.SessionOptions{
			CurrentDirectory:   s.options.Cwd,
			DefaultLibraryPath: s.options.DefaultLibraryPath,
			PositionEncoding:   lsproto.PositionEncodingKindUTF8,
			LoggingEnabled:     false,
		},
	})

	session := NewSession(projectSession, &SessionOptions{
		UseBinaryResponses: !s.options.Async, // Only msgpack uses binary responses
	})
	defer session.Close()

	// Accept connection from transport
	rwc, err := transport.Accept()
	if err != nil {
		return fmt.Errorf("failed to accept connection: %w", err)
	}

	// Create protocol and connection based on async mode
	var conn Conn
	if s.options.Async {
		protocol := NewJSONRPCProtocol(rwc)
		conn = NewAsyncConnWithProtocol(rwc, protocol, session)
	} else {
		protocol := NewMessagePackProtocol(rwc)
		conn = NewSyncConn(rwc, protocol, session)
	}

	// If callbacks are enabled, set the connection on the FS
	if callbackFS != nil {
		callbackFS.SetConnection(ctx, conn)
	}

	return conn.Run(ctx)
}

func (s *Server) createTransport() (Transport, error) {
	spec := s.options.Transport
	switch {
	case spec == "" || spec == "stdio":
		return newStdioTransport(), nil
	case strings.HasPrefix(spec, "pipe="):
		path := strings.TrimPrefix(spec, "pipe=")
		return NewPipeTransport(path)
	case strings.HasPrefix(spec, "fifo="):
		prefix := strings.TrimPrefix(spec, "fifo=")
		return newFIFOTransport(prefix)
	default:
		return nil, fmt.Errorf("unknown transport: %q", spec)
	}
}
