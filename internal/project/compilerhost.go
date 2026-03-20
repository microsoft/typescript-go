package project

import (
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/project/logging"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

var _ compiler.CompilerHost = (*compilerHost)(nil)

type compilerHost struct {
	configFilePath   tspath.Path
	currentDirectory string
	sessionOptions   *SessionOptions

	sourceFS           *sourceFS
	configFileRegistry *ConfigFileRegistry
	sourceFilesByPath  *collections.SyncMap[tspath.Path, *loadedSourceFile]

	project *Project
	builder *ProjectCollectionBuilder
	logger  *logging.LogTree
}

type loadedSourceFile struct {
	once sync.Once
	file *ast.SourceFile
}

func newCompilerHost(
	currentDirectory string,
	project *Project,
	builder *ProjectCollectionBuilder,
	logger *logging.LogTree,
) *compilerHost {
	return &compilerHost{
		configFilePath:   project.configFilePath,
		currentDirectory: currentDirectory,
		sessionOptions:   builder.sessionOptions,

		sourceFS: newSourceFS(true, builder.fs, builder.toPath),
		sourceFilesByPath: &collections.SyncMap[tspath.Path, *loadedSourceFile]{},

		project: project,
		builder: builder,
		logger:  logger,
	}
}

// freeze clears references to mutable state to make the compilerHost safe for use
// after the snapshot has been finalized. See the usage in snapshot.go for more details.
func (c *compilerHost) freeze(snapshotFS *SnapshotFS, configFileRegistry *ConfigFileRegistry) {
	if c.builder == nil {
		panic("freeze can only be called once")
	}
	c.sourceFS.source = snapshotFS
	c.sourceFS.DisableTracking()
	c.configFileRegistry = configFileRegistry
	c.sourceFilesByPath = nil
	c.builder = nil
	c.project = nil
	c.logger = nil
}

func (c *compilerHost) ensureAlive() {
	if c.builder == nil || c.project == nil {
		panic("method must not be called after snapshot initialization")
	}
}

// DefaultLibraryPath implements compiler.CompilerHost.
func (c *compilerHost) DefaultLibraryPath() string {
	return c.sessionOptions.DefaultLibraryPath
}

// FS implements compiler.CompilerHost.
func (c *compilerHost) FS() vfs.FS {
	return c.sourceFS
}

// GetCurrentDirectory implements compiler.CompilerHost.
func (c *compilerHost) GetCurrentDirectory() string {
	return c.currentDirectory
}

// GetResolvedProjectReference implements compiler.CompilerHost.
func (c *compilerHost) GetResolvedProjectReference(fileName string, path tspath.Path) *tsoptions.ParsedCommandLine {
	if c.builder == nil {
		return c.configFileRegistry.GetConfig(path)
	} else {
		// acquireConfigForProject will bypass sourceFS, so track the file here.
		c.sourceFS.Track(fileName)
		return c.builder.configFileRegistryBuilder.acquireConfigForProject(fileName, path, c.project, c.logger)
	}
}

// GetSourceFile implements compiler.CompilerHost. Files are cached in parseCache
// and acquired immediately for the in-progress program.
//
// The compiler can ask for the same tspath.Path multiple times while building a
// single program (for example, via different filename casings on a case-insensitive
// filesystem, or when UpdateProgram probes a dirty file before falling back to a
// full rebuild). We only want to acquire parse-cache ownership once per path,
// because snapshot disposal releases source files once per final program file.
func (c *compilerHost) GetSourceFile(opts ast.SourceFileParseOptions) *ast.SourceFile {
	c.ensureAlive()
	entry, _ := c.sourceFilesByPath.LoadOrStore(opts.Path, &loadedSourceFile{})
	entry.once.Do(func() {
		if fh := c.sourceFS.GetFileByPath(opts.FileName, opts.Path); fh != nil {
			key := NewParseCacheKey(opts, fh.Hash(), fh.Kind())
			entry.file = c.builder.parseCache.Acquire(key, fh)
		}
	})
	return entry.file
}

// Trace implements compiler.CompilerHost.
func (c *compilerHost) Trace(msg *diagnostics.Message, args ...any) {
	panic("unimplemented")
}
