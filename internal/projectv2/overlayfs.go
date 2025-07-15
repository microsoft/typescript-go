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

type FileHandle interface {
	FileName() string
	Version() int32
	Hash() [sha256.Size]byte
	Content() string
	MatchesDiskText() bool
	IsOverlay() bool
	LineMap() *ls.LineMap
}

type fileBase struct {
	fileName string
	content  string
	hash     [sha256.Size]byte

	lineMapOnce sync.Once
	lineMap     *ls.LineMap
}

func (f *fileBase) FileName() string {
	return f.fileName
}

func (f *fileBase) Hash() [sha256.Size]byte {
	return f.hash
}

func (f *fileBase) Content() string {
	return f.content
}

func (f *fileBase) LineMap() *ls.LineMap {
	f.lineMapOnce.Do(func() {
		f.lineMap = ls.ComputeLineStarts(f.content)
	})
	return f.lineMap
}

type diskFile struct {
	fileBase
}

func newDiskFile(fileName string, content string) *diskFile {
	return &diskFile{
		fileBase: fileBase{
			fileName: fileName,
			content:  content,
			hash:     sha256.Sum256([]byte(content)),
		},
	}
}

var _ FileHandle = (*diskFile)(nil)

func (f *diskFile) Version() int32 {
	return 0
}

func (f *diskFile) MatchesDiskText() bool {
	return true
}

func (f *diskFile) IsOverlay() bool {
	return false
}

var _ FileHandle = (*overlay)(nil)

type overlay struct {
	fileBase
	version         int32
	kind            core.ScriptKind
	matchesDiskText bool
}

func newOverlay(fileName string, content string, version int32, kind core.ScriptKind) *overlay {
	return &overlay{
		fileBase: fileBase{
			fileName: fileName,
			content:  content,
			hash:     sha256.Sum256([]byte(content)),
		},
		version: version,
		kind:    kind,
	}
}

func (o *overlay) Version() int32 {
	return o.version
}

func (o *overlay) Text() string {
	return o.content
}

// MatchesDiskText may return false negatives, but never false positives.
func (o *overlay) MatchesDiskText() bool {
	return o.matchesDiskText
}

func (o *overlay) IsOverlay() bool {
	return true
}

type overlayFS struct {
	toPath           func(string) tspath.Path
	fs               vfs.FS
	positionEncoding lsproto.PositionEncodingKind

	mu       sync.Mutex
	overlays map[tspath.Path]*overlay
}

func newOverlayFS(fs vfs.FS, overlays map[tspath.Path]*overlay, positionEncoding lsproto.PositionEncodingKind, toPath func(string) tspath.Path) *overlayFS {
	return &overlayFS{
		fs:               fs,
		positionEncoding: positionEncoding,
		overlays:         overlays,
		toPath:           toPath,
	}
}

func (fs *overlayFS) getFile(fileName string) FileHandle {
	fs.mu.Lock()
	overlays := fs.overlays
	fs.mu.Unlock()

	path := fs.toPath(fileName)
	if overlay, ok := overlays[path]; ok {
		return overlay
	}

	content, ok := fs.fs.ReadFile(fileName)
	if !ok {
		return nil
	}
	return newDiskFile(fileName, content)
}

func (fs *overlayFS) processChanges(changes []FileChange) FileChangeSummary {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	var result FileChangeSummary
	newOverlays := maps.Clone(fs.overlays)
	for _, change := range changes {
		path := change.URI.Path(fs.fs.UseCaseSensitiveFileNames())
		switch change.Kind {
		case FileChangeKindOpen:
			result.Opened.Add(change.URI)
			newOverlays[path] = newOverlay(
				ls.DocumentURIToFileName(change.URI),
				change.Content,
				change.Version,
				ls.LanguageKindToScriptKind(change.LanguageKind),
			)
		case FileChangeKindChange:
			result.Changed.Add(change.URI)
			o, ok := newOverlays[path]
			if !ok {
				panic("overlay not found for change")
			}
			converters := ls.NewConverters(fs.positionEncoding, func(fileName string) *ls.LineMap {
				return o.LineMap()
			})
			for _, textChange := range change.Changes {
				if partialChange := textChange.TextDocumentContentChangePartial; partialChange != nil {
					newContent := converters.FromLSPTextChange(o, partialChange).ApplyTo(o.content)
					o = newOverlay(o.fileName, newContent, o.version, o.kind)
				} else if wholeChange := textChange.TextDocumentContentChangeWholeDocument; wholeChange != nil {
					o = newOverlay(o.fileName, wholeChange.Text, o.version, o.kind)
				}
			}
			o.version = change.Version
			o.hash = sha256.Sum256([]byte(o.content))
			// Assume the overlay does not match disk text after a change. This field
			// is allowed to be a false negative.
			o.matchesDiskText = false
			newOverlays[path] = o
		case FileChangeKindSave:
			result.Saved.Add(change.URI)
			o, ok := newOverlays[path]
			if !ok {
				panic("overlay not found for save")
			}
			o = newOverlay(o.FileName(), o.Content(), o.Version(), o.kind)
			o.matchesDiskText = true
			newOverlays[path] = o
		case FileChangeKindClose:
			// Remove the overlay for the closed file.
			if result.Closed == nil {
				result.Closed = make(map[lsproto.DocumentUri][sha256.Size]byte)
			}
			result.Closed[change.URI] = change.Hash
			delete(newOverlays, path)
		case FileChangeKindWatchCreate:
			result.Created.Add(change.URI)
		case FileChangeKindWatchChange:
			if o, ok := newOverlays[path]; ok {
				if o.matchesDiskText {
					// Assume the overlay does not match disk text after a change.
					newOverlays[path] = newOverlay(o.FileName(), o.Content(), o.Version(), o.kind)
				}
			} else {
				// Only count this as a change if the file is closed.
				result.Changed.Add(change.URI)
			}
		case FileChangeKindWatchDelete:
			if o, ok := newOverlays[path]; ok {
				if o.matchesDiskText {
					newOverlays[path] = newOverlay(o.FileName(), o.Content(), o.Version(), o.kind)
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
