package projectv2

import (
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type StateChangeKind int

const (
	StateChangeKindFile StateChangeKind = iota
	StateChangeKindProgramLoad
)

type StateChange struct {
}

type Snapshot struct {
	id uint64

	// Session options are immutable for the server lifetime,
	// so can be a pointer.
	sessionOptions *SessionOptions

	fs                 *overlayFS
	files              *fileMap
	configuredProjects map[tspath.Path]*Project
}

// ReadFile is stable over the lifetime of the snapshot. It first consults its
// own cache (which includes keys for missing files), and only delegates to the
// file system if the key is not known to the cache.
func (s *Snapshot) ReadFile(uri lsproto.DocumentUri) (string, bool) {
	if f, ok := s.files.Get(uri); ok {
		if f == nil {
			return "", false
		}
		return f.Content(), true
	}
	fh := s.fs.ReadFile(uri)
	s.files.Set(uri, fh)
	if fh != nil {
		return fh.Content(), true
	}
	return "", false
}
