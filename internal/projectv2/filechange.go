package projectv2

import "github.com/microsoft/typescript-go/internal/lsp/lsproto"

type FileChangeKind int

const (
	FileChangeKindOpen FileChangeKind = iota
	FileChangeKindClose
	FileChangeKindChange
	FileChangeKindSave
	FileChangeKindWatchAdd
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
