package projectv2

import (
	"crypto/sha256"
	"maps"
	"sync"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tspath"
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
	overlays map[tspath.Path]*overlay
}

func newOverlayFS(fs vfs.FS, overlays map[tspath.Path]*overlay) *overlayFS {
	return &overlayFS{
		fs:       fs,
		overlays: overlays,
	}
}

func (fs *overlayFS) getFile(uri lsproto.DocumentUri) fileHandle {
	fs.mu.Lock()
	overlays := fs.overlays
	fs.mu.Unlock()

	fileName := ls.DocumentURIToFileName(uri)
	path := tspath.ToPath(fileName, "", fs.fs.UseCaseSensitiveFileNames())
	if overlay, ok := overlays[path]; ok {
		return overlay
	}

	content, ok := fs.fs.ReadFile(fileName)
	if !ok {
		return nil
	}
	return &diskFile{uri: uri, content: content, hash: sha256.Sum256([]byte(content))}
}

func (fs *overlayFS) processChanges(changes []FileChange, converters *ls.Converters) FileChangeSummary {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	var result FileChangeSummary
	newOverlays := maps.Clone(fs.overlays)
	for _, change := range changes {
		path := change.URI.Path(fs.fs.UseCaseSensitiveFileNames())
		switch change.Kind {
		case FileChangeKindOpen:
			result.Opened.Add(change.URI)
			newOverlays[path] = &overlay{
				uri:     change.URI,
				content: change.Content,
				hash:    sha256.Sum256([]byte(change.Content)),
				version: change.Version,
				kind:    ls.LanguageKindToScriptKind(change.LanguageKind),
			}
		case FileChangeKindChange:
			result.Changed.Add(change.URI)
			o, ok := newOverlays[path]
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
			result.Saved.Add(change.URI)
			o, ok := newOverlays[path]
			if !ok {
				panic("overlay not found for save")
			}
			newOverlays[path] = &overlay{
				uri:             o.URI(),
				content:         o.Content(),
				hash:            o.Hash(),
				version:         o.Version(),
				matchesDiskText: true,
			}
		case FileChangeKindClose:
			// Remove the overlay for the closed file.
			result.Closed.Add(change.URI)
			delete(newOverlays, path)
		case FileChangeKindWatchCreate:
			result.Created.Add(change.URI)
		case FileChangeKindWatchChange:
			if o, ok := newOverlays[path]; ok {
				if o.matchesDiskText {
					// Assume the overlay does not match disk text after a change.
					newOverlays[path] = &overlay{
						uri:             o.URI(),
						content:         o.Content(),
						hash:            o.Hash(),
						version:         o.Version(),
						matchesDiskText: false,
					}
				}
			} else {
				// Only count this as a change if the file is closed.
				result.Changed.Add(change.URI)
			}
		case FileChangeKindWatchDelete:
			if o, ok := newOverlays[path]; ok {
				if o.matchesDiskText {
					newOverlays[path] = &overlay{
						uri:             o.URI(),
						content:         o.Content(),
						hash:            o.Hash(),
						version:         o.Version(),
						matchesDiskText: false,
					}
				}
			} else {
				// Only count this as a deletion if the file is closed.
				result.Deleted.Add(change.URI)
			}
		default:
			panic("unhandled file change kind")
		}
	}

	fs.overlays = newOverlays
	return result
}
