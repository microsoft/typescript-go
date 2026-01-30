package api

import (
	"context"
	"sync"
	"time"

	"github.com/go-json-experiment/json"
	"github.com/microsoft/typescript-go/internal/vfs"
)

// CallbackFS wraps a base filesystem and can delegate certain operations
// to the client via RPC callbacks. This allows the API client to provide
// a virtual filesystem (e.g., in-memory files for testing).
type CallbackFS struct {
	base vfs.FS

	mu               sync.RWMutex
	conn             Conn
	ctx              context.Context
	enabledCallbacks map[string]bool
}

// NewCallbackFS creates a new CallbackFS wrapping the given base filesystem.
func NewCallbackFS(base vfs.FS) *CallbackFS {
	return &CallbackFS{
		base:             base,
		enabledCallbacks: make(map[string]bool),
	}
}

// Configure enables the specified callbacks and sets the connection for RPC calls.
// This should be called when the client sends the "configure" message.
func (fs *CallbackFS) Configure(ctx context.Context, conn Conn, callbacks []string) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	fs.ctx = ctx
	fs.conn = conn
	fs.enabledCallbacks = make(map[string]bool)
	for _, cb := range callbacks {
		fs.enabledCallbacks[cb] = true
	}
}

// isEnabled returns true if the named callback is enabled.
func (fs *CallbackFS) isEnabled(name string) bool {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	return fs.enabledCallbacks[name]
}

// call invokes a callback on the client and returns the result.
func (fs *CallbackFS) call(name string, arg any) ([]byte, error) {
	fs.mu.RLock()
	conn := fs.conn
	ctx := fs.ctx
	fs.mu.RUnlock()

	if conn == nil {
		return nil, nil
	}

	result, err := conn.Call(ctx, name, arg)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Callback names
const (
	cbReadFile             = "readFile"
	cbFileExists           = "fileExists"
	cbDirectoryExists      = "directoryExists"
	cbGetAccessibleEntries = "getAccessibleEntries"
	cbRealpath             = "realpath"
)

// UseCaseSensitiveFileNames implements vfs.FS.
func (fs *CallbackFS) UseCaseSensitiveFileNames() bool {
	return fs.base.UseCaseSensitiveFileNames()
}

// ReadFile implements vfs.FS.
func (fs *CallbackFS) ReadFile(path string) (contents string, ok bool) {
	if fs.isEnabled(cbReadFile) {
		result, err := fs.call(cbReadFile, path)
		if err != nil {
			return "", false
		}
		if result == nil || len(result) == 0 || string(result) == "null" {
			return "", false
		}
		var content string
		if err := json.Unmarshal(result, &content); err != nil {
			return "", false
		}
		return content, true
	}
	return fs.base.ReadFile(path)
}

// FileExists implements vfs.FS.
func (fs *CallbackFS) FileExists(path string) bool {
	if fs.isEnabled(cbFileExists) {
		result, err := fs.call(cbFileExists, path)
		if err != nil {
			return false
		}
		if result == nil || len(result) == 0 {
			return false
		}
		var exists bool
		if err := json.Unmarshal(result, &exists); err != nil {
			return false
		}
		return exists
	}
	return fs.base.FileExists(path)
}

// DirectoryExists implements vfs.FS.
func (fs *CallbackFS) DirectoryExists(path string) bool {
	if fs.isEnabled(cbDirectoryExists) {
		result, err := fs.call(cbDirectoryExists, path)
		if err != nil {
			return false
		}
		if result == nil || len(result) == 0 {
			return false
		}
		var exists bool
		if err := json.Unmarshal(result, &exists); err != nil {
			return false
		}
		return exists
	}
	return fs.base.DirectoryExists(path)
}

// GetAccessibleEntries implements vfs.FS.
func (fs *CallbackFS) GetAccessibleEntries(path string) vfs.Entries {
	if fs.isEnabled(cbGetAccessibleEntries) {
		result, err := fs.call(cbGetAccessibleEntries, path)
		if err != nil {
			return vfs.Entries{}
		}
		if result == nil || len(result) == 0 || string(result) == "null" {
			return vfs.Entries{}
		}
		var entries struct {
			Files       []string `json:"files"`
			Directories []string `json:"directories"`
		}
		if err := json.Unmarshal(result, &entries); err != nil {
			return vfs.Entries{}
		}
		return vfs.Entries{
			Files:       entries.Files,
			Directories: entries.Directories,
		}
	}
	return fs.base.GetAccessibleEntries(path)
}

// Realpath implements vfs.FS.
func (fs *CallbackFS) Realpath(path string) string {
	if fs.isEnabled(cbRealpath) {
		result, err := fs.call(cbRealpath, path)
		if err != nil {
			return path
		}
		if result == nil || len(result) == 0 || string(result) == "null" {
			return path
		}
		var realpath string
		if err := json.Unmarshal(result, &realpath); err != nil {
			return path
		}
		return realpath
	}
	return fs.base.Realpath(path)
}

// WriteFile implements vfs.FS - always delegates to base (no callback support).
func (fs *CallbackFS) WriteFile(path string, data string, writeByteOrderMark bool) error {
	return fs.base.WriteFile(path, data, writeByteOrderMark)
}

// Remove implements vfs.FS - always delegates to base (no callback support).
func (fs *CallbackFS) Remove(path string) error {
	return fs.base.Remove(path)
}

// Chtimes implements vfs.FS - always delegates to base (no callback support).
func (fs *CallbackFS) Chtimes(path string, aTime time.Time, mTime time.Time) error {
	return fs.base.Chtimes(path, aTime, mTime)
}

// Stat implements vfs.FS - always delegates to base (no callback support).
func (fs *CallbackFS) Stat(path string) vfs.FileInfo {
	return fs.base.Stat(path)
}

// WalkDir implements vfs.FS - always delegates to base (no callback support).
func (fs *CallbackFS) WalkDir(root string, walkFn vfs.WalkDirFunc) error {
	return fs.base.WalkDir(root, walkFn)
}
