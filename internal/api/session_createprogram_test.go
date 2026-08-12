package api

import (
	"context"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"gotest.tools/v3/assert"
)

func TestCreateProgram(t *testing.T) {
	t.Parallel()

	const fileName = "/home/projects/p/index.ts"
	projectSession, _ := projecttestutil.Setup(map[string]any{
		fileName: `export const value: string = 1;`,
	})
	defer projectSession.Close()

	session := NewSession(projectSession, nil)
	defer session.Close()
	ctx := context.Background()

	// Create a basic program snapshot without replacing the session's latest snapshot.
	baseResponse, err := session.handleUpdateSnapshot(ctx, &UpdateSnapshotParams{})
	assert.NilError(t, err)

	response, err := session.handleCreateProgram(ctx, &CreateProgramParams{
		RootFiles: []DocumentIdentifier{{FileName: fileName}},
		Options: core.CompilerOptions{
			NoLib:  core.TSTrue,
			Strict: core.TSTrue,
		},
	})
	assert.NilError(t, err)
	assert.Assert(t, response.Snapshot != baseResponse.Snapshot)
	assert.Equal(t, session.latestSnapshot, baseResponse.Snapshot)
	assert.Assert(t, response.Project != nil)
	assert.DeepEqual(t, response.Project.RootFiles, []string{fileName})
	assert.Equal(t, response.Project.CompilerOptions.Strict, core.TSTrue)

	// The program snapshot contains one synthetic project and supports program queries.
	snapshot, err := session.getSnapshotData(response.Snapshot)
	assert.NilError(t, err)
	assert.Equal(t, len(snapshot.snapshot.ProjectCollection.Projects()), 1)

	diagnostics, err := session.handleGetSemanticDiagnostics(ctx, &GetDiagnosticsParams{
		Snapshot: response.Snapshot,
		Project:  response.Project.Id,
		File:     &DocumentIdentifier{FileName: fileName},
	})
	assert.NilError(t, err)
	assert.Assert(t, len(diagnostics) > 0)

	// Releasing the program snapshot disposes it without affecting the base snapshot.
	_, err = session.handleRelease(ctx, &ReleaseParams{Snapshot: response.Snapshot})
	assert.NilError(t, err)
	_, err = session.getSnapshotData(response.Snapshot)
	assert.ErrorContains(t, err, "not found")
	_, err = session.getSnapshotData(baseResponse.Snapshot)
	assert.NilError(t, err)
}
