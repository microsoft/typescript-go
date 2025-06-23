package projectv2

import (
	"crypto/sha256"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

var _ fileHandle = (*overlay)(nil)

type overlay struct {
	uri             lsproto.DocumentUri
	version         int32
	content         string
	hash            [sha256.Size]byte
	kind            core.ScriptKind
	matchesDiskText bool
}

// Content implements fileHandle.
func (o *overlay) Content() string {
	return o.content
}

// URI implements fileHandle.
func (o *overlay) URI() lsproto.DocumentUri {
	return o.uri
}

// Version implements fileHandle.
func (o *overlay) Version() int32 {
	return o.version
}

func (o *overlay) Hash() [sha256.Size]byte {
	return o.hash
}

// MatchesDiskText implements fileHandle. May return false positives but never false negatives.
func (o *overlay) MatchesDiskText() bool {
	return o.matchesDiskText
}
