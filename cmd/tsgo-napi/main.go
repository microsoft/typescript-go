package main

import "C"

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/microsoft/typescript-go/internal/api"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/json"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/napi"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

// sessionState holds the active session and its context.
type sessionState struct {
	session        *api.Session
	projectSession *project.Session
	ctx            context.Context
	cancel         context.CancelFunc
}

var (
	currentSession *sessionState
	sessionMu      sync.Mutex
)

func init() {
	napi.RegisterInit(initModule)
}

func initModule(env napi.Env, exports napi.Value) (napi.Value, error) {
	createSessionFn, err := env.CreateFunction("createSession", createSession)
	if err != nil {
		return napi.Value{}, err
	}
	if err2 := env.SetNamedProperty(exports, "createSession", createSessionFn); err2 != nil {
		return napi.Value{}, err2
	}

	requestFn, err := env.CreateFunction("request", request)
	if err != nil {
		return napi.Value{}, err
	}
	if err2 := env.SetNamedProperty(exports, "request", requestFn); err2 != nil {
		return napi.Value{}, err2
	}

	requestBinaryFn, err := env.CreateFunction("requestBinary", requestBinary)
	if err != nil {
		return napi.Value{}, err
	}
	if err2 := env.SetNamedProperty(exports, "requestBinary", requestBinaryFn); err2 != nil {
		return napi.Value{}, err2
	}

	closeFn, err := env.CreateFunction("close", closeSession)
	if err != nil {
		return napi.Value{}, err
	}
	if err2 := env.SetNamedProperty(exports, "close", closeFn); err2 != nil {
		return napi.Value{}, err2
	}

	return exports, nil
}

// createSession(cwd: string, defaultLibraryPath?: string, callbacks?: string): void
func createSession(env napi.Env, args []napi.Value) napi.Value {
	if len(args) < 1 {
		_ = env.ThrowError("createSession requires at least 1 argument (cwd)")
		undef, _ := env.GetUndefinedValue()
		return undef
	}

	cwd, err := env.StringValueToString(args[0])
	if err != nil {
		_ = env.ThrowError("createSession: cwd must be a string")
		undef, _ := env.GetUndefinedValue()
		return undef
	}

	// Optional: default library path (required for noembed builds where
	// libs are not compiled into the binary).
	var externalLibPath string
	if len(args) > 1 {
		s, err := env.StringValueToString(args[1])
		if err == nil {
			externalLibPath = s
		}
	}

	var callbacksList []string
	if len(args) > 2 {
		callbacksStr, err := env.StringValueToString(args[2])
		if err == nil && callbacksStr != "" {
			callbacksList = strings.Split(callbacksStr, ",")
		}
	}

	sessionMu.Lock()
	defer sessionMu.Unlock()

	// Clean up existing session
	if currentSession != nil {
		currentSession.session.Close()
		currentSession.cancel()
		currentSession = nil
	}

	// Determine the default library path. In embed mode, bundled.LibPath()
	// returns a virtual "bundled:///" URI and the libs are served from the
	// Go embed FS. In noembed mode the libs must exist on disk; use the
	// path provided by the JS caller so it doesn't depend on os.Executable().
	var defaultLibraryPath string
	if bundled.Embedded {
		defaultLibraryPath = bundled.LibPath()
	} else if externalLibPath != "" {
		defaultLibraryPath = externalLibPath
	} else {
		// Fallback: try bundled.LibPath() which will look next to the
		// executable. This will panic if the libs aren't there, but it
		// preserves the existing behavior for callers that don't pass a path.
		defaultLibraryPath = bundled.LibPath()
	}

	var fs vfs.FS = bundled.WrapFS(osvfs.FS())

	// Wrap the base FS with callbackFS if callbacks are requested
	var callbackFS *napiCallbackFS
	if len(callbacksList) > 0 {
		callbackFS = newNapiCallbackFS(fs, callbacksList)
		fs = callbackFS
	}
	// Suppress unused variable warning
	_ = callbackFS

	ctx, cancel := context.WithCancel(context.Background())

	projectSession := project.NewSession(&project.SessionInit{
		BackgroundCtx: ctx,
		Logger:        nil,
		FS:            fs,
		Options: &project.SessionOptions{
			CurrentDirectory:   cwd,
			DefaultLibraryPath: defaultLibraryPath,
			PositionEncoding:   lsproto.PositionEncodingKindUTF8,
			LoggingEnabled:     false,
		},
	})

	session := api.NewSession(projectSession, &api.SessionOptions{
		UseBinaryResponses: true, // NAPI can handle binary directly
	})

	currentSession = &sessionState{
		session:        session,
		projectSession: projectSession,
		ctx:            ctx,
		cancel:         cancel,
	}

	undef, _ := env.GetUndefinedValue()
	return undef
}

// request(method: string, payload: string): string
func request(env napi.Env, args []napi.Value) napi.Value {
	undef, _ := env.GetUndefinedValue()

	if len(args) < 2 {
		_ = env.ThrowError("request requires 2 arguments (method, payload)")
		return undef
	}

	method, err := env.StringValueToString(args[0])
	if err != nil {
		_ = env.ThrowError("request: method must be a string")
		return undef
	}

	payload, err := env.StringValueToString(args[1])
	if err != nil {
		_ = env.ThrowError("request: payload must be a string")
		return undef
	}

	sessionMu.Lock()
	state := currentSession
	sessionMu.Unlock()

	if state == nil {
		_ = env.ThrowError("request: no active session (call createSession first)")
		return undef
	}

	result, callErr := handleRequest(state, method, []byte(payload))
	if callErr != nil {
		_ = env.ThrowError(callErr.Error())
		return undef
	}

	resultStr, err := env.StringToStringValue(string(result))
	if err != nil {
		_ = env.ThrowError(fmt.Sprintf("request: failed to create string result: %v", err))
		return undef
	}

	return resultStr
}

// requestBinary(method: string, payload: Uint8Array): Buffer
func requestBinary(env napi.Env, args []napi.Value) napi.Value {
	undef, _ := env.GetUndefinedValue()

	if len(args) < 2 {
		_ = env.ThrowError("requestBinary requires 2 arguments (method, payload)")
		return undef
	}

	method, err := env.StringValueToString(args[0])
	if err != nil {
		_ = env.ThrowError("requestBinary: method must be a string")
		return undef
	}

	payload, err := env.BufferToBytes(args[1])
	if err != nil {
		_ = env.ThrowError(fmt.Sprintf("requestBinary: payload must be a Buffer: %v", err))
		return undef
	}

	sessionMu.Lock()
	state := currentSession
	sessionMu.Unlock()

	if state == nil {
		_ = env.ThrowError("requestBinary: no active session (call createSession first)")
		return undef
	}

	result, callErr := handleRequest(state, method, payload)
	if callErr != nil {
		_ = env.ThrowError(callErr.Error())
		return undef
	}

	buf, err := env.BytesToBuffer(result)
	if err != nil {
		_ = env.ThrowError(fmt.Sprintf("requestBinary: failed to create buffer: %v", err))
		return undef
	}

	return buf
}

// handleRequest dispatches a request to the session handler.
func handleRequest(state *sessionState, method string, payload []byte) (result []byte, err error) {
	// Recover from panics in the handler
	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			err = fmt.Errorf("panic: %v\n%s", r, stack)
		}
	}()

	handlerResult, handlerErr := state.session.HandleRequest(state.ctx, method, json.Value(payload))
	if handlerErr != nil {
		return nil, handlerErr
	}

	// Check if result is raw binary
	if raw, ok := handlerResult.(api.RawBinary); ok {
		return []byte(raw), nil
	}

	// Marshal to JSON
	resultBytes, marshalErr := json.Marshal(handlerResult)
	if marshalErr != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", marshalErr)
	}

	return resultBytes, nil
}

// close(): void
func closeSession(env napi.Env, args []napi.Value) napi.Value {
	_ = args
	sessionMu.Lock()
	defer sessionMu.Unlock()

	if currentSession != nil {
		currentSession.session.Close()
		currentSession.cancel()
		currentSession = nil
	}

	undef, _ := env.GetUndefinedValue()
	return undef
}

// napiCallbackFS wraps a base filesystem and delegates certain operations
// to JS callbacks via NAPI. This is a simplified version that doesn't
// support callbacks yet - it just uses the base FS.
type napiCallbackFS struct {
	base             vfs.FS
	enabledCallbacks map[string]bool
}

func newNapiCallbackFS(base vfs.FS, callbacks []string) *napiCallbackFS {
	enabled := make(map[string]bool, len(callbacks))
	for _, cb := range callbacks {
		enabled[cb] = true
	}
	return &napiCallbackFS{
		base:             base,
		enabledCallbacks: enabled,
	}
}

func (fs *napiCallbackFS) UseCaseSensitiveFileNames() bool {
	return fs.base.UseCaseSensitiveFileNames()
}

func (fs *napiCallbackFS) ReadFile(path string) (string, bool) {
	return fs.base.ReadFile(path)
}

func (fs *napiCallbackFS) FileExists(path string) bool {
	return fs.base.FileExists(path)
}

func (fs *napiCallbackFS) DirectoryExists(path string) bool {
	return fs.base.DirectoryExists(path)
}

func (fs *napiCallbackFS) GetAccessibleEntries(path string) vfs.Entries {
	return fs.base.GetAccessibleEntries(path)
}

func (fs *napiCallbackFS) Realpath(path string) string {
	return fs.base.Realpath(path)
}

func (fs *napiCallbackFS) WalkDir(root string, walkFn vfs.WalkDirFunc) error {
	return fs.base.WalkDir(root, walkFn)
}

func (fs *napiCallbackFS) Stat(path string) vfs.FileInfo {
	return fs.base.Stat(path)
}

func (fs *napiCallbackFS) Remove(path string) error {
	return fs.base.Remove(path)
}

func (fs *napiCallbackFS) WriteFile(path string, data string) error {
	return fs.base.WriteFile(path, data)
}

func (fs *napiCallbackFS) Chtimes(path string, aTime time.Time, mTime time.Time) error {
	return fs.base.Chtimes(path, aTime, mTime)
}

func main() {}
