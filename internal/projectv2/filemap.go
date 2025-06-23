package projectv2

import (
	"crypto/sha256"
	"maps"

	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
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

type fileMap struct {
	files    map[lsproto.DocumentUri]*diskFile
	overlays map[lsproto.DocumentUri]*overlay
	missing  map[lsproto.DocumentUri]struct{}
}

func newFileMap() *fileMap {
	return &fileMap{
		files:    make(map[lsproto.DocumentUri]*diskFile),
		overlays: make(map[lsproto.DocumentUri]*overlay),
		missing:  make(map[lsproto.DocumentUri]struct{}),
	}
}

func (m *fileMap) clone() *fileMap {
	return &fileMap{
		files:    maps.Clone(m.files),
		overlays: maps.Clone(m.overlays),
		missing:  maps.Clone(m.missing),
	}
}

// Get returns the file handle for the given URI, if it exists.
// The second return value indicates whether the key was known to the map.
// The return value is (nil, true) if the file is known to be missing.
func (m *fileMap) Get(uri lsproto.DocumentUri) (fileHandle, bool) {
	if f, ok := m.overlays[uri]; ok {
		return f, true
	}
	if f, ok := m.files[uri]; ok {
		return f, true
	}
	if _, ok := m.missing[uri]; ok {
		return nil, true
	}
	return nil, false
}

func (m *fileMap) Set(uri lsproto.DocumentUri, f fileHandle) {
	if o, ok := f.(*overlay); ok {
		m.overlays[uri] = o
		delete(m.files, uri)
		delete(m.missing, uri)
	} else if d, ok := f.(*diskFile); ok {
		m.files[uri] = d
		delete(m.overlays, uri)
		delete(m.missing, uri)
	} else if f == nil {
		m.missing[uri] = struct{}{}
	} else {
		panic("unexpected file handle type")
	}
}
