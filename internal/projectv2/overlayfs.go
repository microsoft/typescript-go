package projectv2

import (
	"crypto/sha256"
	"sync"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type overlays struct {
	mu sync.Mutex
	m  map[lsproto.DocumentUri]*overlay
}

type overlayFS struct {
	fs vfs.FS
	*overlays
}

func newOverlayFSFromOverlays(fs vfs.FS, overlays *overlays) *overlayFS {
	return &overlayFS{
		fs:       fs,
		overlays: overlays,
	}
}

func (fs *overlayFS) getFile(uri lsproto.DocumentUri) fileHandle {
	fs.overlays.mu.Lock()
	overlay, ok := fs.overlays.m[uri]
	fs.overlays.mu.Unlock()
	if ok {
		return overlay
	}

	content, ok := fs.fs.ReadFile(string(uri))
	if !ok {
		return nil
	}
	return &diskFile{uri: uri, content: content, hash: sha256.Sum256([]byte(content))}
}

func (fs *overlayFS) replaceOverlay(uri lsproto.DocumentUri, content string, version int32, kind core.ScriptKind) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fs.m[uri] = &overlay{
		uri:     uri,
		content: content,
		hash:    sha256.Sum256([]byte(content)),
		version: version,
		kind:    kind,
	}
}
