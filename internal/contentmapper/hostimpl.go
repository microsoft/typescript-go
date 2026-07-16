package contentmapper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/ipc"
	"github.com/microsoft/typescript-go/internal/json"
	"github.com/microsoft/typescript-go/internal/spanmap"
)

// ProtocolVersion is the content mapper protocol version this host speaks.
const ProtocolVersion = 1

const (
	MethodInitialize = "initialize"
	MethodTransform  = "transform"
)

// InitializeParams is the parameter object for the initialize request.
type InitializeParams struct {
	ProtocolVersion int `json:"protocolVersion"`
}

// InitializeResult is the mapper's response to the initialize request.
type InitializeResult struct {
	ProtocolVersion int `json:"protocolVersion"`
}

// TransformParams is the parameter object for the transform request.
type TransformParams struct {
	FileName string `json:"fileName"`
	Content  string `json:"content"`
	// ConfigFileName is the tsconfig file name of the project the file is being transformed for.
	ConfigFileName string `json:"configFileName"`
	// CompilerOptions holds the values of the options the mapper declared in initialize, keyed by option
	// name and ordered by the mapper's declaration. It is an empty object when the mapper declared none.
	CompilerOptions *collections.OrderedMap[string, json.Value] `json:"compilerOptions"`
}

// TransformResult is the mapper's response to a transform request.
type TransformResult struct {
	Text        string          `json:"text"`
	ScriptKind  core.ScriptKind `json:"scriptKind,omitempty"`
	Diagnostics []Diagnostic    `json:"diagnostics,omitempty"`
	// Mappings is the span map's tuple-array JSON (see spanmap.Marshal). Absent or empty means the output
	// is fully synthesized (no part corresponds to the original).
	Mappings json.Value `json:"mappings,omitempty"`
}

// Diagnostic is an error reported by a mapper.
type Diagnostic struct {
	MessageText string `json:"messageText"`
	Start       int    `json:"start"`
	Length      int    `json:"length"`
	Code        int32  `json:"code,omitempty"`
	Source      string `json:"source,omitempty"`
}

// dialFunc establishes a running connection to a mapper. In production it spawns the mapper's process;
// tests substitute an in-memory connection. It returns the connection and a closer that tears it down.
type dialFunc func(ctx context.Context, mapper *Mapper) (ipc.Conn, io.Closer, error)

// host manages one child process per mapper identity. It is the production implementation of Host.
type host struct {
	ctx    context.Context
	cancel context.CancelFunc
	stop   func() bool
	dial   dialFunc

	mu    sync.Mutex
	conns map[string]*mapperConn
}

type mapperConn struct {
	conn   ipc.Conn
	closer io.Closer
	// err, when non-nil, records that this mapper failed to start; it is cached so we do not repeatedly
	// try (and fail) to spawn a broken mapper.
	err error
}

var _ Host = (*host)(nil)

// Spawner starts a child process, returning its stdio as an io.ReadWriteCloser (Read is the
// process's stdout, Write is its stdin) whose Close tears the process down. This seam keeps os/exec out
// of this package: production hosts spawn a real process, tests supply an in-process pipe.
type Spawner interface {
	Spawn(command []string, dir string) (io.ReadWriteCloser, error)
}

// SpawnerFunc adapts a spawn function to the Spawner interface.
type SpawnerFunc func(command []string, dir string) (io.ReadWriteCloser, error)

func (f SpawnerFunc) Spawn(command []string, dir string) (io.ReadWriteCloser, error) {
	return f(command, dir)
}

// NewHost creates a Host that spawns each mapper's process via the given spawner and drives it over a
// JSON-RPC connection. The host's lifetime is bound to ctx: cancelling it (e.g. the CLI's signal context
// on SIGINT, or a build/watch session ending) tears every mapper process down, so owners of a session
// context need not close the host explicitly. Close does the same synchronously.
func NewHost(ctx context.Context, spawner Spawner) Host {
	return newWithDial(ctx, func(ctx context.Context, mapper *Mapper) (ipc.Conn, io.Closer, error) {
		if len(mapper.Exec) == 0 {
			return nil, nil, fmt.Errorf("content mapper %q declares no command to run", mapper.Package)
		}
		rwc, err := spawner.Spawn(mapper.Exec, mapper.PackageDirectory)
		if err != nil {
			return nil, nil, err
		}
		conn := ipc.NewAsyncConn(rwc, rejectHandler{})
		go func() { _ = conn.Run(ctx) }()
		if err := handshake(ctx, conn); err != nil {
			_ = rwc.Close()
			return nil, nil, fmt.Errorf("content mapper %q failed to initialize: %w", mapper.Package, err)
		}
		return conn, rwc, nil
	})
}

func newWithDial(ctx context.Context, dial dialFunc) *host {
	hostCtx, cancel := context.WithCancel(ctx)
	h := &host{ctx: hostCtx, cancel: cancel, dial: dial, conns: make(map[string]*mapperConn)}
	h.stop = context.AfterFunc(ctx, func() { _ = h.Close() })
	return h
}

// Transform sends the file's content to the mapper's process and decodes the transformed result. The
// mapper receives the subset of the project's compiler options it declared in initialize (an empty
// object if it declared none) and the project's tsconfig file name.
func (h *host) Transform(mapper *Mapper, request Request) (Result, error) {
	conn, err := h.connFor(mapper)
	if err != nil {
		return Result{}, err
	}
	options, err := mapper.MarshalDeclaredOptions(request.CompilerOptions)
	if err != nil {
		return Result{}, err
	}
	raw, err := conn.Call(h.ctx, MethodTransform, TransformParams{
		FileName:        request.FileName,
		Content:         request.Content,
		ConfigFileName:  request.ConfigFileName,
		CompilerOptions: options,
	})
	if err != nil {
		return Result{}, err
	}
	return decodeTransformResult(raw)
}

// Close shuts down every mapper process. It is safe to call more than once and is invoked automatically
// when the context passed to New is cancelled.
func (h *host) Close() error {
	h.stop()
	h.cancel()
	h.mu.Lock()
	defer h.mu.Unlock()
	var errs []error
	for _, mc := range h.conns {
		if mc.closer != nil {
			if err := mc.closer.Close(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	h.conns = nil
	return errors.Join(errs...)
}

// connFor returns the connection for a mapper's identity, spawning its process on first use. Mappers
// sharing an identity share a single process; the dial is serialized so an identity is spawned once.
func (h *host) connFor(mapper *Mapper) (ipc.Conn, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.conns == nil {
		return nil, errors.New("content mapper host is closed")
	}
	key := mapper.Identity()
	if mc, ok := h.conns[key]; ok {
		return mc.conn, mc.err
	}
	conn, closer, err := h.dial(h.ctx, mapper)
	h.conns[key] = &mapperConn{conn: conn, closer: closer, err: err}
	return conn, err
}

func handshake(ctx context.Context, conn ipc.Conn) error {
	raw, err := conn.Call(ctx, MethodInitialize, InitializeParams{ProtocolVersion: ProtocolVersion})
	if err != nil {
		return err
	}
	var res InitializeResult
	if err := json.Unmarshal(raw, &res); err != nil {
		return err
	}
	if res.ProtocolVersion != ProtocolVersion {
		return fmt.Errorf("unsupported protocol version %d (expected %d)", res.ProtocolVersion, ProtocolVersion)
	}
	return nil
}

func decodeTransformResult(raw json.Value) (Result, error) {
	var res TransformResult
	if err := json.Unmarshal(raw, &res); err != nil {
		return Result{}, err
	}
	// Any script kind the mapper does not produce a valid, non-Unknown value for defaults to a .ts file.
	scriptKind := core.ScriptKindTS
	switch res.ScriptKind {
	case core.ScriptKindJS, core.ScriptKindJSX, core.ScriptKindTS, core.ScriptKindTSX, core.ScriptKindJSON:
		scriptKind = res.ScriptKind
	}
	result := Result{
		Text:       res.Text,
		ScriptKind: scriptKind,
	}
	// A successful transform always carries a span map. Absent or empty mappings describe fully
	// synthesized output (no segment corresponds to the original), so decode to an empty map rather than
	// nil, which would mean "not content-mapped".
	if len(res.Mappings) > 0 {
		mappings, err := spanmap.Unmarshal(res.Mappings)
		if err != nil {
			return Result{}, err
		}
		result.Mappings = mappings
	} else {
		result.Mappings = spanmap.New(nil)
	}
	for _, d := range res.Diagnostics {
		result.Diagnostics = append(result.Diagnostics, ast.NewExternalDiagnostic(
			nil,
			core.NewTextRange(d.Start, d.Start+d.Length),
			d.Source,
			diagnostics.CategoryError,
			d.Code,
			d.MessageText,
		))
	}
	return result, nil
}

// rejectHandler rejects any request initiated by the mapper. The content mapper protocol is currently
// parent-driven only; a request from the child is a protocol violation.
type rejectHandler struct{}

func (rejectHandler) HandleRequest(ctx context.Context, method string, params json.Value) (any, error) {
	return nil, fmt.Errorf("content mapper sent an unexpected request: %s", method)
}

func (rejectHandler) HandleNotification(ctx context.Context, method string, params json.Value) error {
	return nil
}
