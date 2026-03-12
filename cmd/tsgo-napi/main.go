package main

import "C"

import (
	"context"
	"fmt"
	"runtime/debug"
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
	callbackFS     *napiCallbackFS
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

// createSession(cwd: string, defaultLibraryPath?: string, fsCallbacks?: object): void
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
		typ, _ := env.TypeOf(args[1])
		if typ == "string" {
			s, err2 := env.StringValueToString(args[1])
			if err2 == nil {
				externalLibPath = s
			}
		}
	}

	// Optional: FS callbacks object with methods like readFile, fileExists, etc.
	var fsCallbacksObj *napi.Value
	if len(args) > 2 {
		typ, _ := env.TypeOf(args[2])
		if typ == "object" {
			fsCallbacksObj = &args[2]
		}
	}

	sessionMu.Lock()
	defer sessionMu.Unlock()

	// Clean up existing session
	if currentSession != nil {
		currentSession.session.Close()
		currentSession.cancel()
		if currentSession.callbackFS != nil {
			currentSession.callbackFS.release(env)
		}
		currentSession = nil
	}

	// Determine the default library path.
	var defaultLibraryPath string
	if bundled.Embedded {
		defaultLibraryPath = bundled.LibPath()
	} else if externalLibPath != "" {
		defaultLibraryPath = externalLibPath
	} else {
		defaultLibraryPath = bundled.LibPath()
	}

	var fs vfs.FS = bundled.WrapFS(osvfs.FS())

	// Wrap the base FS with callbackFS if a JS FS object is provided
	var callbackFS *napiCallbackFS
	if fsCallbacksObj != nil {
		var err2 error
		callbackFS, err2 = newNapiCallbackFS(env, fs, *fsCallbacksObj)
		if err2 != nil {
			_ = env.ThrowError(fmt.Sprintf("createSession: failed to set up FS callbacks: %v", err2))
			undef, _ := env.GetUndefinedValue()
			return undef
		}
		fs = callbackFS
	}

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
		callbackFS:     callbackFS,
	}

	undef, _ := env.GetUndefinedValue()
	return undef
}

// handleRequestResult is the result of a HandleRequest call from a goroutine.
type handleRequestResult struct {
	result []byte
	err    error
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

	result, callErr := dispatchRequest(env, state, method, []byte(payload))
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

	result, callErr := dispatchRequest(env, state, method, payload)
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

// dispatchRequest runs HandleRequest. If the session has FS callbacks,
// it runs the handler in a goroutine and pumps the callback channel on
// the main thread so that goroutines spawned by the handler can call
// JS functions through the main thread. If there are no callbacks, it
// runs the handler directly on the current (main) thread.
func dispatchRequest(env napi.Env, state *sessionState, method string, payload []byte) ([]byte, error) {
	if state.callbackFS == nil {
		// No callbacks — run directly, no goroutine needed.
		return handleRequest(state, method, payload)
	}

	// Run HandleRequest in a goroutine so the main thread can pump callbacks.
	doneCh := make(chan handleRequestResult, 1)
	cbFS := state.callbackFS

	// Activate the callback channel for this request. Goroutines
	// calling FS methods will send callbackRequests here.
	cbFS.activate()
	defer cbFS.deactivate()

	go func() {
		result, err := handleRequest(state, method, payload)
		doneCh <- handleRequestResult{result: result, err: err}
	}()

	// Pump: service callback requests from goroutines on the main thread.
	for {
		select {
		case result := <-doneCh:
			return result.result, result.err
		case req := <-cbFS.callbackCh:
			// Execute the JS callback on the main thread.
			req.response <- cbFS.executeCallback(env, req)
		}
	}
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
		if currentSession.callbackFS != nil {
			currentSession.callbackFS.release(env)
		}
		currentSession = nil
	}

	undef, _ := env.GetUndefinedValue()
	return undef
}

// ── Callback FS with channel-based main-thread dispatch ─────────

// callbackRequest is sent from a goroutine to the main thread when a
// FS method needs to call a JS callback.
type callbackRequest struct {
	ref      *napi.Ref               // reference to the JS callback function
	arg      string                  // string argument (typically a file path)
	response chan<- callbackResponse // channel to send the result back
}

// callbackResponse is the result of a JS callback execution.
type callbackResponse struct {
	resultType string // "undefined", "null", "boolean", "string", "object"
	strResult  string // valid when resultType is "string"
	boolResult bool   // valid when resultType is "boolean"
	err        error
}

// napiCallbackFS wraps a base filesystem and delegates certain operations
// to JS callbacks. When called from any goroutine, it marshals the call
// to the main Node.js thread via a channel, preserving full Go parallelism.
type napiCallbackFS struct {
	base vfs.FS

	// References to JS callback functions, preventing GC.
	// nil means the callback is not enabled.
	readFileRef             *napi.Ref
	fileExistsRef           *napi.Ref
	directoryExistsRef      *napi.Ref
	getAccessibleEntriesRef *napi.Ref
	realpathRef             *napi.Ref

	// callbackCh receives callback requests from goroutines.
	// The main thread pumps this channel during request dispatch.
	callbackCh chan callbackRequest

	// activeMu protects the active flag.
	activeMu sync.Mutex
	active   bool
}

func newNapiCallbackFS(env napi.Env, base vfs.FS, fsObj napi.Value) (*napiCallbackFS, error) {
	fs := &napiCallbackFS{
		base:       base,
		callbackCh: make(chan callbackRequest),
	}

	// For each known callback name, check if the JS object has that method
	// and create a reference to it.
	fs.readFileRef = extractCallbackRef(env, fsObj, "readFile")
	fs.fileExistsRef = extractCallbackRef(env, fsObj, "fileExists")
	fs.directoryExistsRef = extractCallbackRef(env, fsObj, "directoryExists")
	fs.getAccessibleEntriesRef = extractCallbackRef(env, fsObj, "getAccessibleEntries")
	fs.realpathRef = extractCallbackRef(env, fsObj, "realpath")

	return fs, nil
}

// extractCallbackRef checks if a JS object has a function property with the
// given name. If so, creates and returns a reference to it; otherwise nil.
func extractCallbackRef(env napi.Env, obj napi.Value, name string) *napi.Ref {
	val, err := env.GetNamedProperty(obj, name)
	if err != nil {
		return nil
	}
	typ, err := env.TypeOf(val)
	if err != nil || typ != "function" {
		return nil
	}
	ref, err := env.CreateReference(val)
	if err != nil {
		return nil
	}
	return &ref
}

// release deletes all callback references.
func (fs *napiCallbackFS) release(env napi.Env) {
	for _, ref := range []*napi.Ref{
		fs.readFileRef,
		fs.fileExistsRef,
		fs.directoryExistsRef,
		fs.getAccessibleEntriesRef,
		fs.realpathRef,
	} {
		if ref != nil {
			_ = env.DeleteReference(*ref)
		}
	}
}

// activate enables the callback channel for a request dispatch cycle.
func (fs *napiCallbackFS) activate() {
	fs.activeMu.Lock()
	fs.active = true
	fs.activeMu.Unlock()
}

// deactivate disables the callback channel after a request completes.
func (fs *napiCallbackFS) deactivate() {
	fs.activeMu.Lock()
	fs.active = false
	fs.activeMu.Unlock()
}

// sendCallback sends a callback request to the main thread and blocks
// until the JS callback has been executed and the response is available.
func (fs *napiCallbackFS) sendCallback(ref *napi.Ref, arg string) callbackResponse {
	fs.activeMu.Lock()
	isActive := fs.active
	fs.activeMu.Unlock()

	if !isActive {
		// Not inside a request dispatch cycle — this shouldn't happen in
		// normal operation but return a fallthrough response to be safe.
		return callbackResponse{resultType: "undefined"}
	}

	respCh := make(chan callbackResponse, 1)
	fs.callbackCh <- callbackRequest{
		ref:      ref,
		arg:      arg,
		response: respCh,
	}
	return <-respCh
}

// executeCallback runs a JS callback on the main thread (called by the
// callback pump in dispatchRequest).
func (fs *napiCallbackFS) executeCallback(env napi.Env, req callbackRequest) callbackResponse {
	fn, err := env.GetReferenceValue(*req.ref)
	if err != nil {
		return callbackResponse{err: fmt.Errorf("failed to get callback reference: %w", err)}
	}
	argVal, err := env.StringToStringValue(req.arg)
	if err != nil {
		return callbackResponse{err: fmt.Errorf("failed to create string arg: %w", err)}
	}
	undef, _ := env.GetUndefinedValue()
	result, err := env.CallFunction(undef, fn, []napi.Value{argVal})
	if err != nil {
		return callbackResponse{err: fmt.Errorf("callback failed: %w", err)}
	}

	typ, _ := env.TypeOf(result)
	resp := callbackResponse{resultType: typ}

	switch typ {
	case "string":
		resp.strResult, _ = env.StringValueToString(result)
	case "boolean":
		resp.boolResult, _ = env.BooleanValueToBool(result)
	}

	return resp
}

func (fs *napiCallbackFS) UseCaseSensitiveFileNames() bool {
	return fs.base.UseCaseSensitiveFileNames()
}

// ReadFile implements vfs.FS.
//
// The readFile callback follows the same contract as the JS FileSystem interface:
//   - Return undefined → fall back to real FS
//   - Return null → file not found (no fallback)
//   - Return string → file content
func (fs *napiCallbackFS) ReadFile(path string) (contents string, ok bool) {
	if fs.readFileRef != nil {
		resp := fs.sendCallback(fs.readFileRef, path)
		if resp.err != nil {
			panic(fmt.Sprintf("napiCallbackFS.ReadFile: %v", resp.err))
		}
		switch resp.resultType {
		case "undefined":
			// Fall through to real FS
		case "null":
			return "", false
		case "string":
			return resp.strResult, true
		}
	}
	return fs.base.ReadFile(path)
}

// FileExists implements vfs.FS.
func (fs *napiCallbackFS) FileExists(path string) bool {
	if fs.fileExistsRef != nil {
		resp := fs.sendCallback(fs.fileExistsRef, path)
		if resp.err != nil {
			panic(fmt.Sprintf("napiCallbackFS.FileExists: %v", resp.err))
		}
		if resp.resultType == "boolean" {
			return resp.boolResult
		}
		// undefined/null → fall through
	}
	return fs.base.FileExists(path)
}

// DirectoryExists implements vfs.FS.
func (fs *napiCallbackFS) DirectoryExists(path string) bool {
	if fs.directoryExistsRef != nil {
		resp := fs.sendCallback(fs.directoryExistsRef, path)
		if resp.err != nil {
			panic(fmt.Sprintf("napiCallbackFS.DirectoryExists: %v", resp.err))
		}
		if resp.resultType == "boolean" {
			return resp.boolResult
		}
		// undefined/null → fall through
	}
	return fs.base.DirectoryExists(path)
}

// GetAccessibleEntries implements vfs.FS.
func (fs *napiCallbackFS) GetAccessibleEntries(path string) vfs.Entries {
	if fs.getAccessibleEntriesRef != nil {
		resp := fs.sendCallback(fs.getAccessibleEntriesRef, path)
		if resp.err != nil {
			panic(fmt.Sprintf("napiCallbackFS.GetAccessibleEntries: %v", resp.err))
		}
		if resp.resultType == "string" {
			var rawEntries *struct {
				Files       []string `json:"files"`
				Directories []string `json:"directories"`
			}
			if err := json.Unmarshal([]byte(resp.strResult), &rawEntries); err != nil {
				panic(fmt.Sprintf("napiCallbackFS.GetAccessibleEntries: failed to unmarshal: %v", err))
			}
			if rawEntries != nil {
				return vfs.Entries{
					Files:       rawEntries.Files,
					Directories: rawEntries.Directories,
				}
			}
		}
		// undefined/null → fall through
	}
	return fs.base.GetAccessibleEntries(path)
}

// Realpath implements vfs.FS.
func (fs *napiCallbackFS) Realpath(path string) string {
	if fs.realpathRef != nil {
		resp := fs.sendCallback(fs.realpathRef, path)
		if resp.err != nil {
			panic(fmt.Sprintf("napiCallbackFS.Realpath: %v", resp.err))
		}
		if resp.resultType == "string" {
			return resp.strResult
		}
		// undefined/null → fall through
	}
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
