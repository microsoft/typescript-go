package projectv2

import "github.com/microsoft/typescript-go/internal/lsp/lsproto"

var _ fileHandle = (*overlay)(nil)

type overlay struct {
	uri             lsproto.DocumentUri
	version         int32
	content         string
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

// MatchesDiskText implements fileHandle.
func (o *overlay) MatchesDiskText() bool {
	return o.matchesDiskText
}
