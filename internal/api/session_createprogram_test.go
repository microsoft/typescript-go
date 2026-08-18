package api

import (
	"context"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"github.com/microsoft/typescript-go/internal/tspath"
	"gotest.tools/v3/assert"
)

func TestCreateProgram(t *testing.T) {
	t.Parallel()

	const fileName = "/home/projects/p/index.ts"
	projectSession, sessionUtils := projecttestutil.Setup(map[string]any{
		fileName: `export const value: string = 1;`,
	})
	defer projectSession.Close()

	session := NewSession(projectSession, nil)
	defer session.Close()
	ctx := context.Background()

	// Create a basic program snapshot without replacing the session's latest snapshot.
	baseResponse, err := session.handleUpdateSnapshot(ctx, &UpdateSnapshotParams{})
	assert.NilError(t, err)
	// The valid unsaved LSP overlay is visible to createProgram.
	projectSession.DidOpenFile(
		ctx,
		DocumentIdentifier{FileName: fileName}.ToURI(projectSession.GetCurrentDirectory()),
		1,
		`export const value: string = "valid overlay";`,
		lsproto.LanguageKindTypeScript,
	)

	response, err := session.handleCreateProgram(ctx, &CreateProgramParams{
		RootFiles: []DocumentIdentifier{{FileName: fileName}},
		CreateProgramOptions: CreateProgramOptions{
			CompilerOptions: core.CompilerOptions{
				NoLib:  core.TSTrue,
				Strict: core.TSTrue,
			},
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
	assert.Equal(t, len(diagnostics), 0)

	// Update from the old program using an explicit disk change summary.
	assert.NilError(t, sessionUtils.FS().WriteFile(fileName, `export const value: string = "valid on disk";`))
	updatedResponse, err := session.handleCreateProgram(ctx, &CreateProgramParams{
		RootFiles: []DocumentIdentifier{{FileName: fileName}},
		CreateProgramOptions: CreateProgramOptions{
			CompilerOptions: core.CompilerOptions{
				NoLib:  core.TSTrue,
				Strict: core.TSTrue,
			},
		},
		OldProgram: &CreateProgramOldProgramParams{
			Snapshot: response.Snapshot,
			Project:  response.Project.Id,
		},
		FileChanges: &APIFileChanges{
			Changed: []DocumentIdentifier{{FileName: fileName}},
		},
	})
	assert.NilError(t, err)
	updatedSnapshot, err := session.getSnapshotData(updatedResponse.Snapshot)
	assert.NilError(t, err)
	updatedProject := updatedSnapshot.snapshot.ProjectCollection.InferredProject()
	assert.Assert(t, updatedProject != nil)

	updatedDiagnostics, err := session.handleGetSemanticDiagnostics(ctx, &GetDiagnosticsParams{
		Snapshot: updatedResponse.Snapshot,
		Project:  updatedResponse.Project.Id,
		File:     &DocumentIdentifier{FileName: fileName},
	})
	assert.NilError(t, err)
	assert.Equal(t, len(updatedDiagnostics), 0)

	oldDiagnostics, err := session.handleGetSemanticDiagnostics(ctx, &GetDiagnosticsParams{
		Snapshot: response.Snapshot,
		Project:  response.Project.Id,
		File:     &DocumentIdentifier{FileName: fileName},
	})
	assert.NilError(t, err)
	assert.Equal(t, len(oldDiagnostics), 0)

	// Releasing the program snapshot disposes it without affecting the base snapshot.
	_, err = session.handleRelease(ctx, &ReleaseParams{Snapshot: updatedResponse.Snapshot})
	assert.NilError(t, err)
	_, err = session.handleRelease(ctx, &ReleaseParams{Snapshot: response.Snapshot})
	assert.NilError(t, err)
	_, err = session.getSnapshotData(response.Snapshot)
	assert.ErrorContains(t, err, "not found")
	_, err = session.getSnapshotData(baseResponse.Snapshot)
	assert.NilError(t, err)
}

func TestCreateProgramReusesProgram(t *testing.T) {
	t.Parallel()

	const fileName = "/home/projects/p/index.ts"
	projectSession, sessionUtils := projecttestutil.Setup(map[string]any{
		fileName: `export const value: string = 1;`,
	})
	defer projectSession.Close()

	session := NewSession(projectSession, nil)
	defer session.Close()
	ctx := context.Background()

	// Build the initial program.
	oldResponse, err := session.handleCreateProgram(ctx, &CreateProgramParams{
		RootFiles: []DocumentIdentifier{{FileName: fileName}},
		CreateProgramOptions: CreateProgramOptions{
			CompilerOptions: core.CompilerOptions{
				NoLib:  core.TSTrue,
				Strict: core.TSTrue,
			},
		},
	})
	assert.NilError(t, err)

	// A single named file change should take the Program.UpdateProgram reuse path.
	assert.NilError(t, sessionUtils.FS().WriteFile(fileName, `export const value: string = "valid";`))
	updatedResponse, err := session.handleCreateProgram(ctx, &CreateProgramParams{
		RootFiles: []DocumentIdentifier{{FileName: fileName}},
		CreateProgramOptions: CreateProgramOptions{
			CompilerOptions: core.CompilerOptions{
				NoLib:  core.TSTrue,
				Strict: core.TSTrue,
			},
		},
		OldProgram: &CreateProgramOldProgramParams{
			Snapshot: oldResponse.Snapshot,
			Project:  oldResponse.Project.Id,
		},
		FileChanges: &APIFileChanges{
			Changed: []DocumentIdentifier{{FileName: fileName}},
		},
	})
	assert.NilError(t, err)

	updatedSnapshot, err := session.getSnapshotData(updatedResponse.Snapshot)
	assert.NilError(t, err)
	updatedProject := updatedSnapshot.snapshot.ProjectCollection.InferredProject()
	assert.Assert(t, updatedProject != nil)
	assert.Equal(t, updatedProject.ProgramUpdateKind, project.ProgramUpdateKindCloned)

	// Changing compiler options replaces the command line and intentionally skips reuse.
	changedOptionsResponse, err := session.handleCreateProgram(ctx, &CreateProgramParams{
		RootFiles: []DocumentIdentifier{{FileName: fileName}},
		CreateProgramOptions: CreateProgramOptions{
			CompilerOptions: core.CompilerOptions{
				NoLib:  core.TSTrue,
				Strict: core.TSFalse,
			},
		},
		OldProgram: &CreateProgramOldProgramParams{
			Snapshot: oldResponse.Snapshot,
			Project:  oldResponse.Project.Id,
		},
	})
	assert.NilError(t, err)
	changedOptionsSnapshot, err := session.getSnapshotData(changedOptionsResponse.Snapshot)
	assert.NilError(t, err)
	changedOptionsProject := changedOptionsSnapshot.snapshot.ProjectCollection.InferredProject()
	assert.Assert(t, changedOptionsProject != nil)
	assert.Equal(t, changedOptionsProject.CommandLine.CompilerOptions().Strict, core.TSFalse)
	assert.Equal(t, changedOptionsProject.ProgramUpdateKind, project.ProgramUpdateKindSameFileNames)
}

func TestCreateProgramProjectReferencesAndReuse(t *testing.T) {
	t.Parallel()

	const (
		fileName        = "/home/projects/app/index.ts"
		libConfigName   = "/home/projects/lib/tsconfig.json"
		otherConfigName = "/home/projects/other/tsconfig.json"
	)
	projectSession, sessionUtils := projecttestutil.Setup(map[string]any{
		fileName:                        `export const value: string = 1;`,
		libConfigName:                   `{ "compilerOptions": { "composite": true, "noLib": true }, "files": ["index.ts"] }`,
		"/home/projects/lib/index.ts":   `export const lib = 1;`,
		otherConfigName:                 `{ "compilerOptions": { "composite": true, "noLib": true }, "files": ["index.ts"] }`,
		"/home/projects/other/index.ts": `export const other = 1;`,
	})
	defer projectSession.Close()

	session := NewSession(projectSession, nil)
	defer session.Close()
	ctx := context.Background()
	libReference := &core.ProjectReference{Path: libConfigName, OriginalPath: libConfigName}

	// A config-less createProgram command line still carries and resolves project references.
	oldResponse, err := session.handleCreateProgram(ctx, &CreateProgramParams{
		RootFiles: []DocumentIdentifier{{FileName: fileName}},
		CreateProgramOptions: CreateProgramOptions{
			CompilerOptions:   core.CompilerOptions{NoLib: core.TSTrue, Strict: core.TSTrue},
			ProjectReferences: []*core.ProjectReference{libReference},
		},
	})
	assert.NilError(t, err)
	assert.DeepEqual(t, oldResponse.Project.ParsedCommandLine.ProjectReferences, []*core.ProjectReference{libReference})
	oldSnapshot, err := session.getSnapshotData(oldResponse.Snapshot)
	assert.NilError(t, err)
	resolvedReferences := oldSnapshot.snapshot.ProjectCollection.InferredProject().Program.GetResolvedProjectReferences()
	assert.Equal(t, len(resolvedReferences), 1)
	assert.Equal(t, resolvedReferences[0].ConfigName(), libConfigName)

	// originalPath is display syntax only; matching path/circular values preserve
	// command-line identity, so one changed file can reuse the old program.
	assert.NilError(t, sessionUtils.FS().WriteFile(fileName, `export const value: string = "valid";`))
	equivalentLibReference := &core.ProjectReference{Path: libConfigName, OriginalPath: "../lib"}
	reusedResponse, err := session.handleCreateProgram(ctx, &CreateProgramParams{
		RootFiles: []DocumentIdentifier{{FileName: fileName}},
		CreateProgramOptions: CreateProgramOptions{
			CompilerOptions:   core.CompilerOptions{NoLib: core.TSTrue, Strict: core.TSTrue},
			ProjectReferences: []*core.ProjectReference{equivalentLibReference},
		},
		OldProgram: &CreateProgramOldProgramParams{
			Snapshot: oldResponse.Snapshot,
			Project:  oldResponse.Project.Id,
		},
		FileChanges: &APIFileChanges{Changed: []DocumentIdentifier{{FileName: fileName}}},
	})
	assert.NilError(t, err)
	reusedSnapshot, err := session.getSnapshotData(reusedResponse.Snapshot)
	assert.NilError(t, err)
	assert.Equal(t, reusedSnapshot.snapshot.ProjectCollection.InferredProject().ProgramUpdateKind, project.ProgramUpdateKindCloned)

	// Changing references changes the command line and therefore requires a full program update.
	otherReference := &core.ProjectReference{Path: otherConfigName, OriginalPath: otherConfigName}
	changedResponse, err := session.handleCreateProgram(ctx, &CreateProgramParams{
		RootFiles: []DocumentIdentifier{{FileName: fileName}},
		CreateProgramOptions: CreateProgramOptions{
			CompilerOptions:   core.CompilerOptions{NoLib: core.TSTrue, Strict: core.TSTrue},
			ProjectReferences: []*core.ProjectReference{otherReference},
		},
		OldProgram: &CreateProgramOldProgramParams{
			Snapshot: oldResponse.Snapshot,
			Project:  oldResponse.Project.Id,
		},
	})
	assert.NilError(t, err)
	changedSnapshot, err := session.getSnapshotData(changedResponse.Snapshot)
	assert.NilError(t, err)
	changedProject := changedSnapshot.snapshot.ProjectCollection.InferredProject()
	assert.Equal(t, changedProject.ProgramUpdateKind, project.ProgramUpdateKindSameFileNames)
	assert.DeepEqual(t, changedProject.CommandLine.ProjectReferences(), []*core.ProjectReference{otherReference})
}

func TestCreateProgramFromConfiguredProgramDoesNotRetainOtherProjects(t *testing.T) {
	t.Parallel()

	const (
		configFileName      = "/home/projects/p/tsconfig.json"
		fileName            = "/home/projects/p/index.ts"
		otherConfigFileName = "/home/projects/other/tsconfig.json"
		otherFileName       = "/home/projects/other/index.ts"
	)
	projectSession, sessionUtils := projecttestutil.Setup(map[string]any{
		configFileName:      `{ "compilerOptions": { "noLib": true, "strict": true }, "files": ["index.ts"] }`,
		fileName:            `export const value: string = 1;`,
		otherConfigFileName: `{ "files": ["index.ts"] }`,
		otherFileName:       `export const other = 1;`,
	})
	defer projectSession.Close()

	session := NewSession(projectSession, nil)
	defer session.Close()
	ctx := context.Background()

	// Load two configured projects so the selected old program comes from a non-synthetic, multi-project snapshot.
	baseResponse, err := session.handleUpdateSnapshot(ctx, &UpdateSnapshotParams{
		OpenProjects: []DocumentIdentifier{{FileName: configFileName}, {FileName: otherConfigFileName}},
	})
	assert.NilError(t, err)
	var baseProject *ProjectResponse
	for _, candidate := range baseResponse.Projects {
		if candidate.ConfigFileName == configFileName {
			baseProject = candidate
			break
		}
	}
	assert.Assert(t, baseProject != nil)
	rootFiles := make([]DocumentIdentifier, len(baseProject.RootFiles))
	for i, rootFile := range baseProject.RootFiles {
		rootFiles[i] = DocumentIdentifier{FileName: rootFile}
	}

	// Derive only the selected configured program as a synthetic createProgram project.
	assert.NilError(t, sessionUtils.FS().WriteFile(fileName, `export const value: string = "valid";`))
	updatedResponse, err := session.handleCreateProgram(ctx, &CreateProgramParams{
		RootFiles: rootFiles,
		CreateProgramOptions: CreateProgramOptions{
			CompilerOptions: core.CompilerOptions{
				NoLib:  core.TSTrue,
				Strict: core.TSTrue,
			},
		},
		OldProgram: &CreateProgramOldProgramParams{
			Snapshot: baseResponse.Snapshot,
			Project:  baseProject.Id,
		},
		FileChanges: &APIFileChanges{
			Changed: []DocumentIdentifier{{FileName: fileName}},
		},
	})
	assert.NilError(t, err)

	updatedSnapshot, err := session.getSnapshotData(updatedResponse.Snapshot)
	assert.NilError(t, err)
	// The derived snapshot must not retain configured projects or config state from the unrelated base project.
	assert.Equal(t, len(updatedSnapshot.snapshot.ProjectCollection.Projects()), 1)
	assert.Equal(t, len(updatedSnapshot.snapshot.ProjectCollection.ConfiguredProjects()), 0)
	assert.Assert(t, updatedSnapshot.snapshot.ConfigFileRegistry.GetConfig(tspath.Path(otherConfigFileName)) == nil)
	updatedProject := updatedSnapshot.snapshot.ProjectCollection.InferredProject()
	assert.Assert(t, updatedProject != nil)
	// Configured programs carry ConfigFilePath in their compiler options, while
	// createProgram options intentionally do not. The command lines therefore differ,
	// so this safely rebuilds instead of taking the single-file reuse path.
	assert.Equal(t, updatedProject.ProgramUpdateKind, project.ProgramUpdateKindSameFileNames)
	updatedDiagnostics, err := session.handleGetSemanticDiagnostics(ctx, &GetDiagnosticsParams{
		Snapshot: updatedResponse.Snapshot,
		Project:  updatedResponse.Project.Id,
		File:     &DocumentIdentifier{FileName: fileName},
	})
	assert.NilError(t, err)
	assert.Equal(t, len(updatedDiagnostics), 0)

	// Disposing the derived snapshot must leave the original configured snapshot queryable.
	_, err = session.handleRelease(ctx, &ReleaseParams{Snapshot: updatedResponse.Snapshot})
	assert.NilError(t, err)
	baseDiagnostics, err := session.handleGetSemanticDiagnostics(ctx, &GetDiagnosticsParams{
		Snapshot: baseResponse.Snapshot,
		Project:  baseProject.Id,
		File:     &DocumentIdentifier{FileName: fileName},
	})
	assert.NilError(t, err)
	assert.Equal(t, len(baseDiagnostics), 1)
}
