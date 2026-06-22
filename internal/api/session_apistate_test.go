package api

import (
	"context"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"github.com/microsoft/typescript-go/internal/tspath"
	"gotest.tools/v3/assert"
)

// TestSessionTracksAndReleasesAPIRefs verifies that an API session holds at most
// one ref per opened project/file (opens are idempotent) and releases exactly
// those refs when the session is closed, so it never leaks or over-releases refs
// in the underlying (potentially shared) project session.
func TestSessionTracksAndReleasesAPIRefs(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	t.Run("project opens are idempotent and released on close", func(t *testing.T) {
		t.Parallel()
		const configFileName = "/home/projects/p/tsconfig.json"
		files := map[string]any{
			configFileName:                  `{ "compilerOptions": { "strict": true } }`,
			"/home/projects/p/src/index.ts": `export const x = 1;`,
		}
		projectSession, _ := projecttestutil.Setup(files)
		defer projectSession.Close()
		session := NewSession(projectSession, nil)

		_, err := session.handleUpdateSnapshot(context.Background(), &UpdateSnapshotParams{
			OpenProjects: []DocumentIdentifier{{FileName: configFileName}},
		})
		assert.NilError(t, err)
		assert.Equal(t, session.openProjects.Len(), 1)

		// Opening the same project again must not take an additional ref.
		_, err = session.handleUpdateSnapshot(context.Background(), &UpdateSnapshotParams{
			OpenProjects: []DocumentIdentifier{{FileName: configFileName}},
		})
		assert.NilError(t, err)
		assert.Equal(t, session.openProjects.Len(), 1)

		assert.Assert(t, projectSession.Snapshot().ProjectCollection.ConfiguredProject(tspath.Path(configFileName)) != nil)

		// Closing the session releases the single API ref, so the project is no
		// longer kept loaded.
		session.Close()
		assert.Equal(t, session.openProjects.Len(), 0)
		assert.Assert(t, projectSession.Snapshot().ProjectCollection.ConfiguredProject(tspath.Path(configFileName)) == nil)
	})

	t.Run("explicit close releases the project ref", func(t *testing.T) {
		t.Parallel()
		const configFileName = "/home/projects/p/tsconfig.json"
		files := map[string]any{
			configFileName:                  `{ "compilerOptions": { "strict": true } }`,
			"/home/projects/p/src/index.ts": `export const x = 1;`,
		}
		projectSession, _ := projecttestutil.Setup(files)
		defer projectSession.Close()
		session := NewSession(projectSession, nil)
		defer session.Close()

		_, err := session.handleUpdateSnapshot(context.Background(), &UpdateSnapshotParams{
			OpenProjects: []DocumentIdentifier{{FileName: configFileName}},
		})
		assert.NilError(t, err)
		assert.Equal(t, session.openProjects.Len(), 1)
		assert.Assert(t, projectSession.Snapshot().ProjectCollection.ConfiguredProject(tspath.Path(configFileName)) != nil)

		// Closing a project we hold releases the ref and unloads the project.
		_, err = session.handleUpdateSnapshot(context.Background(), &UpdateSnapshotParams{
			CloseProjects: []DocumentIdentifier{{FileName: configFileName}},
		})
		assert.NilError(t, err)
		assert.Equal(t, session.openProjects.Len(), 0)
		assert.Assert(t, projectSession.Snapshot().ProjectCollection.ConfiguredProject(tspath.Path(configFileName)) == nil)

		// Closing a project we don't hold is a no-op (never over-releases).
		_, err = session.handleUpdateSnapshot(context.Background(), &UpdateSnapshotParams{
			CloseProjects: []DocumentIdentifier{{FileName: configFileName}},
		})
		assert.NilError(t, err)
		assert.Equal(t, session.openProjects.Len(), 0)
	})

	t.Run("file opens are idempotent and released on close", func(t *testing.T) {
		t.Parallel()
		const fileName = "/home/projects/p/src/index.ts"
		files := map[string]any{
			"/home/projects/p/tsconfig.json": `{ "compilerOptions": { "strict": true } }`,
			fileName:                         `export const x = 1;`,
		}
		projectSession, _ := projecttestutil.Setup(files)
		defer projectSession.Close()
		session := NewSession(projectSession, nil)

		_, err := session.handleUpdateSnapshot(context.Background(), &UpdateSnapshotParams{
			OpenFiles: []DocumentIdentifier{{FileName: fileName}},
		})
		assert.NilError(t, err)
		assert.Equal(t, session.openFiles.Len(), 1)

		// Re-opening the same file must not take an additional ref.
		_, err = session.handleUpdateSnapshot(context.Background(), &UpdateSnapshotParams{
			OpenFiles: []DocumentIdentifier{{FileName: fileName}},
		})
		assert.NilError(t, err)
		assert.Equal(t, session.openFiles.Len(), 1)

		// The file should resolve to the configured project via ancestor search.
		assert.Assert(t, projectSession.Snapshot().ProjectCollection.ConfiguredProject(tspath.Path("/home/projects/p/tsconfig.json")) != nil)

		// Closing a file we don't hold is a no-op (never over-releases).
		_, err = session.handleUpdateSnapshot(context.Background(), &UpdateSnapshotParams{
			CloseFiles: []DocumentIdentifier{{FileName: "/home/projects/p/other.ts"}},
		})
		assert.NilError(t, err)
		assert.Equal(t, session.openFiles.Len(), 1)

		// Explicitly closing the held file releases the ref.
		_, err = session.handleUpdateSnapshot(context.Background(), &UpdateSnapshotParams{
			CloseFiles: []DocumentIdentifier{{FileName: fileName}},
		})
		assert.NilError(t, err)
		assert.Equal(t, session.openFiles.Len(), 0)

		session.Close()
		assert.Equal(t, session.openFiles.Len(), 0)
	})
}
