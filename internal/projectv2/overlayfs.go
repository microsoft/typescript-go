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
	Kind() core.ScriptKind
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

func (f *diskFile) Kind() core.ScriptKind {
	return core.GetScriptKindFromFileName(f.fileName)
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

func (o *overlay) Kind() core.ScriptKind {
	return o.kind
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

	// Reduced collection of changes that occurred on a single file
	type fileEvents struct {
		openChange   *FileChange
		closeChange  *FileChange
		watchChanged bool
		changes      []*FileChange
		saved        bool
		created      bool
		deleted      bool
	}

	fileEventMap := make(map[lsproto.DocumentUri]*fileEvents)

	for _, change := range changes {
		uri := change.URI
		events, exists := fileEventMap[uri]
		if exists {
			if events.openChange != nil {
				panic("should see no changes after open")
			}
		} else {
			events = &fileEvents{}
			fileEventMap[uri] = events
		}

		switch change.Kind {
		case FileChangeKindOpen:
			events.openChange = &change
			events.closeChange = nil
			events.watchChanged = false
			events.changes = nil
			events.saved = false
			events.created = false
			events.deleted = false
		case FileChangeKindClose:
			events.closeChange = &change
			events.changes = nil
			events.saved = false
			events.watchChanged = false
		case FileChangeKindChange:
			if events.closeChange != nil {
				panic("should see no changes after close")
			}
			events.changes = append(events.changes, &change)
			events.saved = false
			events.watchChanged = false
		case FileChangeKindSave:
			events.saved = true
		case FileChangeKindWatchCreate:
			if events.deleted {
				// Delete followed by create becomes a change
				events.deleted = false
				events.watchChanged = true
			} else {
				events.created = true
			}
		case FileChangeKindWatchChange:
			if !events.created {
				events.watchChanged = true
			}
		case FileChangeKindWatchDelete:
			events.watchChanged = false
			events.saved = false
			// Delete after create cancels out
			if events.created {
				events.created = false
			} else {
				events.deleted = true
			}
		}
	}

	// Process deduplicated events per file
	for uri, events := range fileEventMap {
		path := uri.Path(fs.fs.UseCaseSensitiveFileNames())
		o := newOverlays[path]

		if events.openChange != nil {
			if result.Opened != "" {
				panic("can only process one file open event at a time")
			}
			result.Opened = uri
			newOverlays[path] = newOverlay(
				ls.DocumentURIToFileName(uri),
				events.openChange.Content,
				events.openChange.Version,
				ls.LanguageKindToScriptKind(events.openChange.LanguageKind),
			)
			continue
		}

		if events.closeChange != nil {
			if result.Closed == nil {
				result.Closed = make(map[lsproto.DocumentUri][sha256.Size]byte)
			}
			result.Closed[uri] = events.closeChange.Hash
			delete(newOverlays, path)
		}

		if events.watchChanged {
			if o == nil {
				result.Changed.Add(uri)
			} else if o != nil && o.MatchesDiskText() {
				o = newOverlay(o.FileName(), o.Content(), o.Version(), o.kind)
				o.matchesDiskText = false
				newOverlays[path] = o
			}
		}

		if len(events.changes) > 0 {
			result.Changed.Add(uri)
			if o == nil {
				panic("overlay not found for changed file: " + uri)
			}
			for _, change := range events.changes {
				converters := ls.NewConverters(fs.positionEncoding, func(fileName string) *ls.LineMap {
					return o.LineMap()
				})
				for _, textChange := range change.Changes {
					if partialChange := textChange.TextDocumentContentChangePartial; partialChange != nil {
						newContent := converters.FromLSPTextChange(o, partialChange).ApplyTo(o.content)
						o = newOverlay(o.fileName, newContent, change.Version, o.kind)
					} else if wholeChange := textChange.TextDocumentContentChangeWholeDocument; wholeChange != nil {
						o = newOverlay(o.fileName, wholeChange.Text, change.Version, o.kind)
					}
				}
				o.version = change.Version
				o.hash = sha256.Sum256([]byte(o.content))
				o.matchesDiskText = false
				newOverlays[path] = o
			}
		}

		if events.saved {
			result.Saved.Add(uri)
			if o == nil {
				panic("overlay not found for saved file: " + uri)
			}
			o = newOverlay(o.FileName(), o.Content(), o.Version(), o.kind)
			o.matchesDiskText = true
			newOverlays[path] = o
		}

		if events.created && o == nil {
			result.Created.Add(uri)
		}

		if events.deleted && o == nil {
			result.Deleted.Add(uri)
		}
	}

	fs.overlays = newOverlays
	return result
}
