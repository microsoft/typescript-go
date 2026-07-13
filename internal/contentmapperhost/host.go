// Package contentmapperhost drives the configured content mappers: it spawns each mapper's package as a
// child process and talks to it over a JSON-RPC connection (reusing internal/ipc), turning foreign file
// content into TypeScript. Processes are consolidated by mapper identity, so many projects that use the
// same mapper version share a single process.
package contentmapperhost

import (
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/contentmapper"
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
	// CompilerOptions names the compiler options (e.g. "target", "jsx") whose values the mapper wants to
	// receive on each transform. Options not listed are not sent.
	CompilerOptions []string `json:"compilerOptions,omitempty"`
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
	// Mappings is the span map's tuple-array JSON (see spanmap.Marshal). Empty means identity mapping.
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
// tests substitute an in-memory connection. It returns the connection, a closer that tears it down, and
// the compiler option names the mapper declared it depends on.
type dialFunc func(ctx context.Context, mapper *contentmapper.Mapper) (ipc.Conn, io.Closer, []string, error)

// host manages one child process per mapper identity. It is the production implementation of Host.
type host struct {
	ctx    context.Context
	cancel context.CancelFunc
	dial   dialFunc

	mu    sync.Mutex
	conns map[string]*mapperConn
}

type mapperConn struct {
	conn   ipc.Conn
	closer io.Closer
	// optionKeys are the compiler option names the mapper declared in initialize.
	optionKeys []string
	// err, when non-nil, records that this mapper failed to start; it is cached so we do not repeatedly
	// try (and fail) to spawn a broken mapper.
	err error
}

var _ Host = (*host)(nil)

// Spawner starts a content mapper's process, returning its stdio as an io.ReadWriteCloser (Read is the
// process's stdout, Write is its stdin) whose Close tears the process down. This seam keeps os/exec out
// of this package: production hosts spawn a real process, tests supply an in-process pipe.
type Spawner interface {
	Spawn(command []string, dir string) (io.ReadWriteCloser, error)
}

// New creates a Host that spawns each mapper's process via the given spawner and drives it over a
// JSON-RPC connection. The context bounds the lifetime of the spawned processes and their connections: it
// is the CLI's signal context or the LSP's request/session context, and cancelling it (or calling Close)
// tears every mapper down.
func New(ctx context.Context, spawner Spawner) Host {
	return newWithDial(ctx, func(ctx context.Context, mapper *contentmapper.Mapper) (ipc.Conn, io.Closer, []string, error) {
		if len(mapper.Exec) == 0 {
			return nil, nil, nil, fmt.Errorf("content mapper %q declares no command to run", mapper.Package)
		}
		rwc, err := spawner.Spawn(mapper.Exec, mapper.PackageDirectory)
		if err != nil {
			return nil, nil, nil, err
		}
		conn := ipc.NewAsyncConn(rwc, rejectHandler{})
		go func() { _ = conn.Run(ctx) }()
		optionKeys, err := handshake(ctx, conn)
		if err != nil {
			_ = rwc.Close()
			return nil, nil, nil, fmt.Errorf("content mapper %q failed to initialize: %w", mapper.Package, err)
		}
		return conn, rwc, optionKeys, nil
	})
}

func newWithDial(ctx context.Context, dial dialFunc) *host {
	ctx, cancel := context.WithCancel(ctx)
	return &host{ctx: ctx, cancel: cancel, dial: dial, conns: make(map[string]*mapperConn)}
}

// Transform sends the file's content to the mapper's process and decodes the transformed result. The
// mapper receives the subset of the project's compiler options it declared in initialize (an empty
// object if it declared none) and the project's tsconfig file name.
func (h *host) Transform(mapper *contentmapper.Mapper, request Request) (Result, error) {
	conn, optionKeys, err := h.connFor(mapper)
	if err != nil {
		return Result{}, err
	}
	options, err := marshalDeclaredOptions(request.CompilerOptions, optionKeys)
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

// Close shuts down every mapper process.
func (h *host) Close() error {
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
// sharing an identity share a single process; the dial is serialized so an identity is spawned once. It
// also returns the compiler option names that mapper declared it depends on.
func (h *host) connFor(mapper *contentmapper.Mapper) (ipc.Conn, []string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	key := mapper.Identity()
	if mc, ok := h.conns[key]; ok {
		return mc.conn, mc.optionKeys, mc.err
	}
	conn, closer, optionKeys, err := h.dial(h.ctx, mapper)
	h.conns[key] = &mapperConn{conn: conn, closer: closer, optionKeys: optionKeys, err: err}
	return conn, optionKeys, err
}

func handshake(ctx context.Context, conn ipc.Conn) ([]string, error) {
	raw, err := conn.Call(ctx, MethodInitialize, InitializeParams{ProtocolVersion: ProtocolVersion})
	if err != nil {
		return nil, err
	}
	var res InitializeResult
	if err := json.Unmarshal(raw, &res); err != nil {
		return nil, err
	}
	if res.ProtocolVersion != ProtocolVersion {
		return nil, fmt.Errorf("unsupported protocol version %d (expected %d)", res.ProtocolVersion, ProtocolVersion)
	}
	return res.CompilerOptions, nil
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
	if len(res.Mappings) > 0 {
		mappings, err := spanmap.Unmarshal(res.Mappings)
		if err != nil {
			return Result{}, err
		}
		result.Mappings = mappings
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

// compilerOptionFields maps each CompilerOptions option name (its json tag) to its struct field index.
var compilerOptionFields = sync.OnceValue(func() map[string]int {
	t := reflect.TypeFor[core.CompilerOptions]()
	fields := make(map[string]int, t.NumField())
	for i := range t.NumField() {
		name, _, _ := strings.Cut(t.Field(i).Tag.Get("json"), ",")
		if name != "" && name != "-" {
			fields[name] = i
		}
	}
	return fields
})

// marshalDeclaredOptions marshals just the named compiler options, in the given order, skipping any that
// are unset. Marshaling only the requested fields avoids serializing the whole CompilerOptions when a
// mapper depends on few options (or none).
func marshalDeclaredOptions(options *core.CompilerOptions, names []string) (*collections.OrderedMap[string, json.Value], error) {
	out := collections.NewOrderedMapWithSizeHint[string, json.Value](len(names))
	if options == nil || len(names) == 0 {
		return out, nil
	}
	fields := compilerOptionFields()
	v := reflect.ValueOf(options).Elem()
	for _, name := range names {
		i, ok := fields[name]
		if !ok {
			continue
		}
		field := v.Field(i)
		if field.IsZero() {
			continue
		}
		raw, err := json.Marshal(field.Interface())
		if err != nil {
			return nil, err
		}
		out.Set(name, json.Value(raw))
	}
	return out, nil
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
