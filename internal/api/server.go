package api

import (
	"context"
	"fmt"
	"io"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

// StdioServerOptions configures the STDIO-based API server.
type StdioServerOptions struct {
	In                 io.ReadCloser
	Out                io.WriteCloser
	Err                io.Writer
	Cwd                string
	DefaultLibraryPath string
	// FS is the filesystem to use. If nil, the OS filesystem is used.
	FS vfs.FS
	// Callbacks specifies which filesystem operations should be delegated
	// to the client (e.g., "readFile", "fileExists"). Empty means no callbacks.
	Callbacks []string
	// Async enables JSON-RPC protocol with async connection handling.
	// When false (default), uses MessagePack protocol with sync connection.
	Async bool
}

// StdioServer runs an API session over STDIO using MessagePack protocol.
// This is the entry point for the synchronous STDIO-based API used by
// native TypeScript tooling integration.
type StdioServer struct {
	projectSession *project.Session
	session        *Session
	callbackFS     *CallbackFS
	options        *StdioServerOptions
}

// NewStdioServer creates a new STDIO-based API server.
func NewStdioServer(options *StdioServerOptions) *StdioServer {
	if options.Cwd == "" {
		panic("StdioServerOptions.Cwd is required")
	}

	baseFS := options.FS
	if baseFS == nil {
		baseFS = bundled.WrapFS(osvfs.FS())
	}

	// Wrap the base FS with CallbackFS if callbacks are requested
	var callbackFS *CallbackFS
	var fs vfs.FS = baseFS
	if len(options.Callbacks) > 0 {
		callbackFS = NewCallbackFS(baseFS, options.Callbacks)
		fs = callbackFS
	}

	projectSession := project.NewSession(&project.SessionInit{
		BackgroundCtx: context.Background(),
		Logger:        nil, // TODO: Add logging support
		FS:            fs,
		Options: &project.SessionOptions{
			CurrentDirectory:   options.Cwd,
			DefaultLibraryPath: options.DefaultLibraryPath,
			PositionEncoding:   lsproto.PositionEncodingKindUTF8,
			LoggingEnabled:     false,
		},
	})

	session := NewSession(projectSession, &SessionOptions{
		UseBinaryResponses: !options.Async, // Only msgpack uses binary responses
		SyncRequests:       !options.Async, // Only msgpack protocol requires synchronous request handling
	})

	return &StdioServer{
		projectSession: projectSession,
		session:        session,
		callbackFS:     callbackFS,
		options:        options,
	}
}

// Run starts the server and blocks until the connection closes.
func (s *StdioServer) Run(ctx context.Context) error {
	transport := NewStdioTransport(s.options.In, s.options.Out)
	defer transport.Close()
	defer s.session.Close()

	// Accept connection from transport
	rwc, err := transport.Accept()
	if err != nil {
		return fmt.Errorf("failed to accept connection: %w", err)
	}

	// Create protocol and connection based on async mode
	var conn Conn
	if s.options.Async {
		protocol := NewJSONRPCProtocol(rwc)
		conn = NewAsyncConnWithProtocol(rwc, protocol, s.session)
	} else {
		protocol := NewMessagePackProtocol(rwc)
		conn = NewSyncConn(rwc, protocol, s.session)
	}

	// If callbacks are enabled, set the connection on the FS
	if s.callbackFS != nil {
		s.callbackFS.SetConnection(ctx, conn)
	}

	return conn.Run(ctx)
}

// Close closes the server and releases resources.
func (s *StdioServer) Close() {
	s.session.Close()
	s.projectSession.Close()
}
