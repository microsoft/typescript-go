package projectv2

import (
	"sync"

	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type overlayFS struct {
	fs vfs.FS

	mu       sync.Mutex
	overlays map[lsproto.DocumentUri]*overlay
}

func (fs *overlayFS) ReadFile(uri lsproto.DocumentUri) fileHandle {
	fs.mu.Lock()
	overlay, ok := fs.overlays[uri]
	fs.mu.Unlock()
	if ok {
		return overlay
	}

	content, ok := fs.fs.ReadFile(string(uri))
	if !ok {
		return nil
	}
	return &diskFile{uri: uri, content: content}
}
