package api

import (
	"context"
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

	// Wrap the base FS with CallbackFS to support client-provided virtual filesystems
	callbackFS := NewCallbackFS(baseFS)

	projectSession := project.NewSession(&project.SessionInit{
		BackgroundCtx: context.Background(),
		Logger:        nil, // TODO: Add logging support
		FS:            callbackFS,
		Options: &project.SessionOptions{
			CurrentDirectory:   options.Cwd,
			DefaultLibraryPath: options.DefaultLibraryPath,
			PositionEncoding:   lsproto.PositionEncodingKindUTF8,
			LoggingEnabled:     false,
		},
	})

	session := NewSession(projectSession, &SessionOptions{
		CallbackFS:         callbackFS,
		UseBinaryResponses: true,
		SyncRequests:       true, // msgpack protocol requires synchronous request handling
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

	return s.session.RunWithProtocol(ctx, transport, func(rw io.ReadWriter) Protocol {
		return NewMessagePackProtocol(rw)
	})
}

// Close closes the server and releases resources.
func (s *StdioServer) Close() {
	s.session.Close()
	s.projectSession.Close()
}
