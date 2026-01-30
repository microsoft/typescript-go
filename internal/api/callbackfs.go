package api

import (
	"context"
	"fmt"
	"time"

	"github.com/go-json-experiment/json"
	"github.com/microsoft/typescript-go/internal/vfs"
)

// CallbackFS wraps a base filesystem and delegates certain operations
// to the client via RPC callbacks. This allows the API client to provide
// a virtual filesystem (e.g., in-memory files for testing).
//
// The callbacks to enable are specified at construction time via the
// --callbacks CLI flag. The connection is set via SetConnection after
// the transport connection is established.
type CallbackFS struct {
	base             vfs.FS
	enabledCallbacks map[string]bool

	// conn and ctx are set after connection is established
	conn Conn
	ctx  context.Context
}

// Callback names that can be enabled
const (
	CbReadFile             = "readFile"
	CbFileExists           = "fileExists"
	CbDirectoryExists      = "directoryExists"
	CbGetAccessibleEntries = "getAccessibleEntries"
	CbRealpath             = "realpath"
)

// NewCallbackFS creates a new CallbackFS wrapping the given base filesystem.
// The callbacks slice specifies which filesystem operations should be delegated
// to the client (e.g., "readFile", "fileExists").
func NewCallbackFS(base vfs.FS, callbacks []string) *CallbackFS {
	enabled := make(map[string]bool, len(callbacks))
	for _, cb := range callbacks {
		enabled[cb] = true
	}
	return &CallbackFS{
		base:             base,
		enabledCallbacks: enabled,
	}
}

// SetConnection sets the RPC connection for callbacks.
// This must be called after the transport connection is established
// but before any filesystem operations that need callbacks.
func (fs *CallbackFS) SetConnection(ctx context.Context, conn Conn) {
	fs.ctx = ctx
	fs.conn = conn
}

// isEnabled returns true if the named callback is enabled.
func (fs *CallbackFS) isEnabled(name string) bool {
	return fs.enabledCallbacks[name]
}

// call invokes a callback on the client and returns the result.
func (fs *CallbackFS) call(name string, arg any) ([]byte, error) {
	if fs.conn == nil {
		return nil, fmt.Errorf("CallbackFS: %s called before connection set", name)
	}

	result, err := fs.conn.Call(fs.ctx, name, arg)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UseCaseSensitiveFileNames implements vfs.FS.
func (fs *CallbackFS) UseCaseSensitiveFileNames() bool {
	return fs.base.UseCaseSensitiveFileNames()
}

// ReadFile implements vfs.FS.
func (fs *CallbackFS) ReadFile(path string) (contents string, ok bool) {
	if fs.isEnabled(CbReadFile) {
		result, err := fs.call(CbReadFile, path)
		if err != nil {
			panic(err)
		}
		if string(result) == "null" {
			return "", false
		}
		if len(result) > 0 {
			var content string
			if err := json.Unmarshal(result, &content); err != nil {
				panic(err)
			}
			return content, true
		}
	}
	return fs.base.ReadFile(path)
}

// FileExists implements vfs.FS.
func (fs *CallbackFS) FileExists(path string) bool {
	if fs.isEnabled(CbFileExists) {
		result, err := fs.call(CbFileExists, path)
		if err != nil {
			panic(err)
		}
		if len(result) > 0 {
			return string(result) == "true"
		}
	}
	return fs.base.FileExists(path)
}

// DirectoryExists implements vfs.FS.
func (fs *CallbackFS) DirectoryExists(path string) bool {
	if fs.isEnabled(CbDirectoryExists) {
		result, err := fs.call(CbDirectoryExists, path)
		if err != nil {
			panic(err)
		}
		if len(result) > 0 {
			return string(result) == "true"
		}
	}
	return fs.base.DirectoryExists(path)
}

// GetAccessibleEntries implements vfs.FS.
func (fs *CallbackFS) GetAccessibleEntries(path string) vfs.Entries {
	if fs.isEnabled(CbGetAccessibleEntries) {
		result, err := fs.call(CbGetAccessibleEntries, path)
		if err != nil {
			panic(err)
		}
		if len(result) > 0 {
			var rawEntries *struct {
				Files       []string `json:"files"`
				Directories []string `json:"directories"`
			}
			if err := json.Unmarshal(result, &rawEntries); err != nil {
				panic(err)
			}
			if rawEntries != nil {
				return vfs.Entries{
					Files:       rawEntries.Files,
					Directories: rawEntries.Directories,
				}
			}
		}
	}
	return fs.base.GetAccessibleEntries(path)
}

// Realpath implements vfs.FS.
func (fs *CallbackFS) Realpath(path string) string {
	if fs.isEnabled(CbRealpath) {
		result, err := fs.call(CbRealpath, path)
		if err != nil {
			panic(err)
		}
		if len(result) > 0 {
			var realpath string
			if err := json.Unmarshal(result, &realpath); err != nil {
				panic(err)
			}
			return realpath
		}
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
