package projectv2

import (
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

type FileChangeKind int

const (
	FileChangeKindOpen FileChangeKind = iota
	FileChangeKindClose
	FileChangeKindChange
	FileChangeKindSave
	FileChangeKindWatchCreate
	FileChangeKindWatchChange
	FileChangeKindWatchDelete
)

type FileChange struct {
	Kind         FileChangeKind
	URI          lsproto.DocumentUri
	Version      int32                                    // Only set for Open/Change
	Content      string                                   // Only set for Open
	LanguageKind lsproto.LanguageKind                     // Only set for Open
	Changes      []lsproto.TextDocumentContentChangeEvent // Only set for Change
}

type FileChangeSummary struct {
	Opened  collections.Set[lsproto.DocumentUri]
	Closed  collections.Set[lsproto.DocumentUri]
	Changed collections.Set[lsproto.DocumentUri]
	Saved   collections.Set[lsproto.DocumentUri]
	Created collections.Set[lsproto.DocumentUri]
	Deleted collections.Set[lsproto.DocumentUri]
}

func (f FileChangeSummary) IsEmpty() bool {
	return f.Opened.Len() == 0 && f.Closed.Len() == 0 && f.Changed.Len() == 0 && f.Saved.Len() == 0 && f.Created.Len() == 0 && f.Deleted.Len() == 0
}
