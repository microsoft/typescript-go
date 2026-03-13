package lsp

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"gotest.tools/v3/assert"
)

func TestRenameFilesToWatchedFileEvents(t *testing.T) {
	t.Parallel()

	t.Run("simple rename", func(t *testing.T) {
		t.Parallel()

		events := renameFilesToWatchedFileEvents([]*lsproto.FileRename{
			{
				OldUri: string(lsconv.FileNameToDocumentURI("/home/projects/foo.ts")),
				NewUri: string(lsconv.FileNameToDocumentURI("/home/projects/bar.ts")),
			},
		})

		assert.Equal(t, len(events), 2)
		eventsByURI := map[lsproto.DocumentUri]lsproto.FileChangeType{}
		for _, event := range events {
			eventsByURI[event.Uri] = event.Type
		}
		assert.Equal(t, eventsByURI[lsconv.FileNameToDocumentURI("/home/projects/foo.ts")], lsproto.FileChangeTypeDeleted)
		assert.Equal(t, eventsByURI[lsconv.FileNameToDocumentURI("/home/projects/bar.ts")], lsproto.FileChangeTypeCreated)
	})

	t.Run("swap rename becomes change events", func(t *testing.T) {
		t.Parallel()

		events := renameFilesToWatchedFileEvents([]*lsproto.FileRename{
			{
				OldUri: string(lsconv.FileNameToDocumentURI("/home/projects/a.ts")),
				NewUri: string(lsconv.FileNameToDocumentURI("/home/projects/b.ts")),
			},
			{
				OldUri: string(lsconv.FileNameToDocumentURI("/home/projects/b.ts")),
				NewUri: string(lsconv.FileNameToDocumentURI("/home/projects/a.ts")),
			},
		})

		assert.Equal(t, len(events), 2)
		for _, event := range events {
			assert.Equal(t, event.Type, lsproto.FileChangeTypeChanged)
		}
	})
}
