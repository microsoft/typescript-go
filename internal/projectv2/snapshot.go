package projectv2

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/projectv2/ata"
	"github.com/microsoft/typescript-go/internal/tspath"
)

var snapshotID atomic.Uint64

type Snapshot struct {
	id       uint64
	parentId uint64
	refCount atomic.Int32

	// Session options are immutable for the server lifetime,
	// so can be a pointer.
	sessionOptions *SessionOptions
	toPath         func(fileName string) tspath.Path

	// Immutable state, cloned between snapshots
	diskFiles                          map[tspath.Path]*diskFile
	ProjectCollection                  *ProjectCollection
	ConfigFileRegistry                 *ConfigFileRegistry
	compilerOptionsForInferredProjects *core.CompilerOptions
	builderLogs                        *logCollector
}

// NewSnapshot
func NewSnapshot(
	diskFiles map[tspath.Path]*diskFile,
	sessionOptions *SessionOptions,
	parseCache *parseCache,
	extendedConfigCache *extendedConfigCache,
	configFileRegistry *ConfigFileRegistry,
	compilerOptionsForInferredProjects *core.CompilerOptions,
	toPath func(fileName string) tspath.Path,
) *Snapshot {
	id := snapshotID.Add(1)
	s := &Snapshot{
		id: id,

		sessionOptions: sessionOptions,
		toPath:         toPath,

		diskFiles:                          diskFiles,
		ConfigFileRegistry:                 configFileRegistry,
		ProjectCollection:                  &ProjectCollection{toPath: toPath},
		compilerOptionsForInferredProjects: compilerOptionsForInferredProjects,
	}

	s.refCount.Store(1)
	return s
}

func (s *Snapshot) GetDefaultProject(uri lsproto.DocumentUri) *Project {
	fileName := ls.DocumentURIToFileName(uri)
	path := s.toPath(fileName)
	return s.ProjectCollection.GetDefaultProject(fileName, path)
}

func (s *Snapshot) ID() uint64 {
	return s.id
}

type SnapshotChange struct {
	// fileChanges are the changes that have occurred since the last snapshot.
	fileChanges FileChangeSummary
	// requestedURIs are URIs that were requested by the client.
	// The new snapshot should ensure projects for these URIs have loaded programs.
	requestedURIs []lsproto.DocumentUri
	// compilerOptionsForInferredProjects is the compiler options to use for inferred projects.
	// It should only be set the value in the next snapshot should be changed. If nil, the
	// value from the previous snapshot will be copied to the new snapshot.
	compilerOptionsForInferredProjects *core.CompilerOptions
	// ataChanges contains ATA-related changes to apply to projects in the new snapshot.
	ataChanges map[tspath.Path]*ATAStateChange
}

// ATAStateChange represents a change to a project's ATA state.
type ATAStateChange struct {
	ProjectID tspath.Path
	// TypingsInfo is the new typings info for the project.
	TypingsInfo *ata.TypingsInfo
	// TypingsFiles is the new list of typing files for the project.
	TypingsFiles []string
	Logs         *logCollector
}

func (s *Snapshot) Clone(ctx context.Context, change SnapshotChange, session *Session) *Snapshot {
	var logger *logCollector
	if session.options.LoggingEnabled {
		logger = NewLogCollector(fmt.Sprintf("Cloning snapshot %d", s.id))
		defer logger.Close()
	}

	start := time.Now()
	fs := newSnapshotFSBuilder(session.fs.fs, session.fs.overlays, s.diskFiles, session.options.PositionEncoding, s.toPath)
	fs.markDirtyFiles(change.fileChanges)

	compilerOptionsForInferredProjects := s.compilerOptionsForInferredProjects
	if change.compilerOptionsForInferredProjects != nil {
		// !!! mark inferred projects as dirty?
		compilerOptionsForInferredProjects = change.compilerOptionsForInferredProjects
	}

	projectCollectionBuilder := newProjectCollectionBuilder(
		ctx,
		fs,
		s.ProjectCollection,
		s.ConfigFileRegistry,
		compilerOptionsForInferredProjects,
		s.sessionOptions,
		session.parseCache,
		session.extendedConfigCache,
	)

	if change.ataChanges != nil {
		projectCollectionBuilder.DidUpdateATAState(change.ataChanges, logger.Fork("DidUpdateATAState"))
	}

	if !change.fileChanges.IsEmpty() {
		projectCollectionBuilder.DidChangeFiles(change.fileChanges, logger.Fork("DidChangeFiles"))
	}

	for _, uri := range change.requestedURIs {
		projectCollectionBuilder.DidRequestFile(uri, logger.Fork("DidRequestFile"))
	}

	projectCollection, configFileRegistry := projectCollectionBuilder.Finalize(logger)
	snapshotFS, _ := fs.Finalize()

	newSnapshot := NewSnapshot(
		snapshotFS.diskFiles,
		s.sessionOptions,
		session.parseCache,
		session.extendedConfigCache,
		nil,
		s.compilerOptionsForInferredProjects,
		s.toPath,
	)

	newSnapshot.parentId = s.id
	newSnapshot.ProjectCollection = projectCollection
	newSnapshot.ConfigFileRegistry = configFileRegistry
	newSnapshot.builderLogs = logger

	for _, project := range newSnapshot.ProjectCollection.Projects() {
		if project.Program != nil {
			project.host.freeze(snapshotFS, newSnapshot.ConfigFileRegistry)
			session.programCounter.Ref(project.Program)
		}
	}
	for path, config := range newSnapshot.ConfigFileRegistry.configs {
		if config.commandLine != nil && config.commandLine.ConfigFile != nil {
			if prevConfig, ok := s.ConfigFileRegistry.configs[path]; ok {
				if prevConfig.commandLine != nil && config.commandLine.ConfigFile == prevConfig.commandLine.ConfigFile {
					for _, file := range prevConfig.commandLine.ExtendedSourceFiles() {
						// Ref count extended configs that were already loaded in the previous snapshot.
						// New/changed ones were handled during config file registry building.
						session.extendedConfigCache.Ref(s.toPath(file))
					}
				}
			}
		}
	}

	logger.Logf("Finished cloning snapshot %d into snapshot %d in %v", s.id, newSnapshot.id, time.Since(start))
	return newSnapshot
}

func (s *Snapshot) Ref() {
	s.refCount.Add(1)
}

func (s *Snapshot) Deref() bool {
	return s.refCount.Add(-1) == 0
}

func (s *Snapshot) dispose(session *Session) {
	for _, project := range s.ProjectCollection.Projects() {
		if project.Program != nil && session.programCounter.Deref(project.Program) {
			for _, file := range project.Program.SourceFiles() {
				session.parseCache.Release(file)
			}
		}
	}
	for _, config := range s.ConfigFileRegistry.configs {
		if config.commandLine != nil {
			for _, file := range config.commandLine.ExtendedSourceFiles() {
				session.extendedConfigCache.Release(session.toPath(file))
			}
		}
	}
}
