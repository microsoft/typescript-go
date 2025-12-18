package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync/atomic"

	"github.com/go-json-experiment/json/jsontext"
	"github.com/microsoft/typescript-go/internal/api/encoder"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
)

var sessionIDCounter atomic.Uint64

// Session represents an API session that provides programmatic access
// to TypeScript language services through the LSP server.
// It implements the Handler interface to process incoming API requests.
// The session retains a snapshot until the client explicitly requests an update,
// ensuring consistency across multiple requests.
type Session struct {
	id             string
	projectSession *project.Session
	conn           *Conn
	onClose        func()

	// snapshot is the current snapshot for this session.
	// It is retained until the client requests an update.
	snapshot        *project.Snapshot
	snapshotRelease func()
}

// Ensure Session implements Handler
var _ Handler = (*Session)(nil)

// NewSession creates a new API session with the given project session.
// The onClose callback is called when the session is closed to allow
// cleanup (e.g., removing from a server's session map).
func NewSession(projectSession *project.Session, onClose func()) *Session {
	id := sessionIDCounter.Add(1)
	return &Session{
		id:             formatSessionID(id),
		projectSession: projectSession,
		onClose:        onClose,
	}
}

// ID returns the unique identifier for this session.
func (s *Session) ID() string {
	return s.id
}

// ProjectSession returns the underlying project session.
func (s *Session) ProjectSession() *project.Session {
	return s.projectSession
}

// SetConn sets the connection for this session.
func (s *Session) SetConn(conn *Conn) {
	s.conn = conn
}

// Conn returns the connection for this session.
func (s *Session) Conn() *Conn {
	return s.conn
}

// ensureSnapshot lazily initializes the snapshot if it's nil.
func (s *Session) ensureSnapshot() {
	if s.snapshot == nil {
		s.snapshot, s.snapshotRelease = s.projectSession.Snapshot()
	}
}

// HandleRequest implements Handler.
func (s *Session) HandleRequest(ctx context.Context, method string, params jsontext.Value) (any, error) {
	parsed, err := unmarshalPayload(method, params)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	// Ensure we have a snapshot for request processing
	s.ensureSnapshot()

	switch method {
	case string(MethodGetDefaultProjectForFile):
		return s.handleGetDefaultProjectForFile(ctx, parsed.(*GetDefaultProjectForFileParams))
	case string(MethodGetSourceFile):
		return s.handleGetSourceFile(ctx, parsed.(*GetSourceFileParams))
	case "ping":
		return "pong", nil
	default:
		return nil, fmt.Errorf("unknown method: %s", method)
	}
}

// HandleNotification implements Handler.
func (s *Session) HandleNotification(ctx context.Context, method string, params jsontext.Value) error {
	// TODO: Implement notification handling
	return nil
}

// handleGetDefaultProjectForFile returns the default project for a given file.
func (s *Session) handleGetDefaultProjectForFile(ctx context.Context, params *GetDefaultProjectForFileParams) (*ProjectResponse, error) {
	uri := lsproto.DocumentUri("file://" + params.FileName)

	proj := s.snapshot.GetDefaultProject(uri)
	if proj == nil {
		return nil, fmt.Errorf("%w: no project found for file %s", ErrClientError, params.FileName)
	}

	return NewProjectResponse(proj), nil
}

// handleGetSourceFile returns a source file from a project.
func (s *Session) handleGetSourceFile(ctx context.Context, params *GetSourceFileParams) (*SourceFileResponse, error) {
	projectName := parseProjectHandle(params.Project)
	proj := s.snapshot.ProjectCollection.GetProjectByPath(projectName)
	if proj == nil {
		return nil, fmt.Errorf("%w: project %s not found", ErrClientError, projectName)
	}

	program := proj.GetProgram()
	if program == nil {
		return nil, fmt.Errorf("%w: project has no program", ErrClientError)
	}

	sourceFile := program.GetSourceFile(params.FileName)
	if sourceFile == nil {
		return nil, fmt.Errorf("%w: source file not found: %s", ErrClientError, params.FileName)
	}

	// Encode the source file to binary format
	handle := FileHandle(sourceFile)
	data, err := encoder.EncodeSourceFile(sourceFile, string(handle))
	if err != nil {
		return nil, fmt.Errorf("failed to encode source file: %w", err)
	}

	// Base64 encode for JSON transport
	return &SourceFileResponse{
		Data: base64.StdEncoding.EncodeToString(data),
	}, nil
}

// Close closes the session and triggers the onClose callback.
func (s *Session) Close() {
	if s.snapshotRelease != nil {
		s.snapshotRelease()
		s.snapshotRelease = nil
		s.snapshot = nil
	}
	if s.conn != nil {
		s.conn.Close()
	}
	if s.onClose != nil {
		s.onClose()
	}
}

func formatSessionID(id uint64) string {
	return fmt.Sprintf("api-session-%d", id)
}
