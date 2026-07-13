package project

import (
	"context"
	"maps"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// APIUpdate creates a new snapshot incorporating the given file changes and the
// supplied API open/close request. The apiRequest may open or close projects and
// files; opens are tracked in the snapshot (ref-counted) so they persist across
// future updates, and closes release a previously taken ref. Even an empty
// apiRequest ensures all API-opened projects and files are kept up to date.
// Returns a ref'd snapshot (which the caller must Deref when done) and any error
// encountered while applying the request, e.g. failing to load a project to open.
func (s *Session) APIUpdate(ctx context.Context, apiFileChanges FileChangeSummary, apiRequest *APISnapshotRequest) (*Snapshot, error) {
	s.snapshotUpdateMu.Lock()
	defer s.snapshotUpdateMu.Unlock()
	s.cancelScheduledSnapshotUpdate()

	fileChanges, overlays, ataChanges, _ := s.flushChanges(ctx)
	mergeFileChangeSummary(&fileChanges, apiFileChanges)

	newSnapshot := s.updateSnapshotRef(ctx, overlays, SnapshotChange{
		apiRequest:  apiRequest,
		fileChanges: fileChanges,
		ataChanges:  ataChanges,
	})
	return newSnapshot, newSnapshot.apiError
}

// APIUpdateTemporary creates a snapshot that layers a temporary in-memory content
// override for a single file on top of the session's current snapshot. Unlike
// APIUpdate, it does not mutate any session state: the session's current snapshot,
// overlays, and pending file changes are left untouched, and the new snapshot is
// not adopted as the session's current snapshot, so it is never observed by other
// requests and is discarded once the caller releases it. This mirrors tsserver's
// runWithTemporaryFileUpdate, which temporarily edits a file's content, runs a
// query against the resulting program, then reverts.
//
// The returned snapshot carries a single reference (the clone ref); the caller
// must call snapshot.Deref(s) when done.
func (s *Session) APIUpdateTemporary(ctx context.Context, uri lsproto.DocumentUri, newText string) *Snapshot {
	s.snapshotUpdateMu.Lock()
	defer s.snapshotUpdateMu.Unlock()

	baseSnapshot := s.snapshot
	path := uri.Path(baseSnapshot.UseCaseSensitiveFileNames())

	// Build the new overlay set from the base snapshot's overlays plus a temporary
	// override for the target file.
	overlays := maps.Clone(baseSnapshot.fs.overlays)
	if overlays == nil {
		overlays = make(map[tspath.Path]*Overlay)
	}
	version := int32(0)
	scriptKind := core.GetScriptKindFromFileName(uri.FileName())
	if existing, ok := overlays[path]; ok {
		version = existing.Version() + 1
		scriptKind = existing.Kind()
	}
	overlays[path] = newOverlay(uri.FileName(), newText, version, scriptKind)

	var fileChanges FileChangeSummary
	fileChanges.Changed.Add(uri)

	newSnapshot := baseSnapshot.Clone(ctx, SnapshotChange{
		fileChanges: fileChanges,
		ResourceRequest: ResourceRequest{
			Documents: []lsproto.DocumentUri{uri},
		},
	}, overlays, s)
	return newSnapshot
}
