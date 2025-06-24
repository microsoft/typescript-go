package projectv2

import (
	"crypto/sha256"
	"maps"
	"sync"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type fileHandle interface {
	URI() lsproto.DocumentUri
	Version() int32
	Hash() [sha256.Size]byte
	Content() string
	MatchesDiskText() bool
}

type diskFile struct {
	uri     lsproto.DocumentUri
	content string
	hash    [sha256.Size]byte
}

var _ fileHandle = (*diskFile)(nil)

func (f *diskFile) URI() lsproto.DocumentUri {
	return f.uri
}

func (f *diskFile) Version() int32 {
	return 0
}

func (f *diskFile) Hash() [sha256.Size]byte {
	return f.hash
}

func (f *diskFile) Content() string {
	return f.content
}

func (f *diskFile) MatchesDiskText() bool {
	return true
}

var _ fileHandle = (*overlay)(nil)

type overlay struct {
	uri             lsproto.DocumentUri
	version         int32
	content         string
	hash            [sha256.Size]byte
	kind            core.ScriptKind
	matchesDiskText bool
}

func (o *overlay) Content() string {
	return o.content
}

func (o *overlay) URI() lsproto.DocumentUri {
	return o.uri
}

func (o *overlay) Version() int32 {
	return o.version
}

func (o *overlay) Hash() [sha256.Size]byte {
	return o.hash
}

func (o *overlay) FileName() string {
	return ls.DocumentURIToFileName(o.uri)
}

func (o *overlay) Text() string {
	return o.content
}

// MatchesDiskText may return false negatives, but never false positives.
func (o *overlay) MatchesDiskText() bool {
	return o.matchesDiskText
}

type overlayFS struct {
	fs vfs.FS

	mu       sync.Mutex
	overlays map[lsproto.DocumentUri]*overlay
}

func newOverlayFS(fs vfs.FS, overlays map[lsproto.DocumentUri]*overlay) *overlayFS {
	return &overlayFS{
		fs:       fs,
		overlays: overlays,
	}
}

func (fs *overlayFS) getFile(uri lsproto.DocumentUri) fileHandle {
	fs.mu.Lock()
	overlays := fs.overlays
	fs.mu.Unlock()

	if overlay, ok := overlays[uri]; ok {
		return overlay
	}

	content, ok := fs.fs.ReadFile(string(uri))
	if !ok {
		return nil
	}
	return &diskFile{uri: uri, content: content, hash: sha256.Sum256([]byte(content))}
}

func (fs *overlayFS) updateOverlays(changes []FileChange, converters *ls.Converters) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	newOverlays := maps.Clone(fs.overlays)
	for _, change := range changes {
		switch change.Kind {
		case FileChangeKindOpen:
			newOverlays[change.URI] = &overlay{
				uri:     change.URI,
				content: change.Content,
				hash:    sha256.Sum256([]byte(change.Content)),
				version: change.Version,
				kind:    ls.LanguageKindToScriptKind(change.LanguageKind),
			}
		case FileChangeKindChange:
			o, ok := newOverlays[change.URI]
			if !ok {
				panic("overlay not found for change")
			}
			for _, textChange := range change.Changes {
				if partialChange := textChange.TextDocumentContentChangePartial; partialChange != nil {
					newContent := converters.FromLSPTextChange(o, partialChange).ApplyTo(o.content)
					o = &overlay{uri: o.uri, content: newContent} // need intermediate structs to pass back into FromLSPTextChange
				} else if wholeChange := textChange.TextDocumentContentChangeWholeDocument; wholeChange != nil {
					o = &overlay{uri: o.uri, content: wholeChange.Text}
				}
			}
			o.version = change.Version
			o.hash = sha256.Sum256([]byte(o.content))
			// Assume the overlay does not match disk text after a change. This field
			// is allowed to be a false negative.
			o.matchesDiskText = false
		case FileChangeKindSave:
			o, ok := newOverlays[change.URI]
			if !ok {
				panic("overlay not found for save")
			}
			newOverlays[change.URI] = &overlay{
				uri:             o.uri,
				content:         o.Content(),
				hash:            o.Hash(),
				version:         o.Version(),
				matchesDiskText: true,
			}
		case FileChangeKindClose:
			// Remove the overlay for the closed file.
			delete(newOverlays, change.URI)
		case FileChangeKindWatchAdd, FileChangeKindWatchChange, FileChangeKindWatchDelete:
			// !!! set matchesDiskText?
		}
	}

	fs.overlays = newOverlays
}
