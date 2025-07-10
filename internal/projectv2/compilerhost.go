package projectv2

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

var _ compiler.CompilerHost = (*compilerHost)(nil)

type compilerHost struct {
	configFilePath   tspath.Path
	currentDirectory string
	sessionOptions   *SessionOptions

	overlayFS          *overlayFS
	compilerFS         *compilerFS
	configFileRegistry *ConfigFileRegistry

	project *Project
	builder *projectCollectionBuilder
}

func newCompilerHost(
	currentDirectory string,
	project *Project,
	builder *projectCollectionBuilder,
) *compilerHost {
	return &compilerHost{
		configFilePath:   project.configFilePath,
		currentDirectory: currentDirectory,
		sessionOptions:   builder.sessionOptions,

		overlayFS:  builder.fs,
		compilerFS: &compilerFS{overlayFS: builder.fs},

		project: project,
		builder: builder,
	}
}

func (c *compilerHost) freeze(configFileRegistry *ConfigFileRegistry) {
	c.configFileRegistry = configFileRegistry
	c.builder = nil
	c.project = nil
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
	return c.compilerFS
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
		return c.builder.configFileRegistryBuilder.acquireConfigForProject(fileName, path, c.project)
	}
}

// GetSourceFile implements compiler.CompilerHost. GetSourceFile increments
// the ref count of source files it acquires in the parseCache. There should
// be a corresponding release for each call made.
func (c *compilerHost) GetSourceFile(opts ast.SourceFileParseOptions) *ast.SourceFile {
	c.ensureAlive()
	if fh := c.overlayFS.getFile(opts.FileName); fh != nil {
		projectSet := &collections.SyncSet[tspath.Path]{}
		// !!!
		// projectSet, _ = c.builder.fileAssociations.LoadOrStore(fh.URI().Path(c.FS().UseCaseSensitiveFileNames()), projectSet)
		projectSet.Add(c.project.configFilePath)
		return c.builder.parseCache.acquireDocument(fh, opts, c.getScriptKind(opts.FileName))
	}
	return nil
}

// NewLine implements compiler.CompilerHost.
func (c *compilerHost) NewLine() string {
	return c.sessionOptions.NewLine
}

// Trace implements compiler.CompilerHost.
func (c *compilerHost) Trace(msg string) {
	panic("unimplemented")
}

func (c *compilerHost) getScriptKind(fileName string) core.ScriptKind {
	// Customizing script kind per file extension is a common plugin / LS host customization case
	// which can probably be replaced with static info in the future
	return core.GetScriptKindFromFileName(fileName)
}

var _ vfs.FS = (*compilerFS)(nil)

type compilerFS struct {
	overlayFS *overlayFS
}

// DirectoryExists implements vfs.FS.
func (fs *compilerFS) DirectoryExists(path string) bool {
	return fs.overlayFS.fs.DirectoryExists(path)
}

// FileExists implements vfs.FS.
func (fs *compilerFS) FileExists(path string) bool {
	if fh := fs.overlayFS.getFile(path); fh != nil {
		return true
	}
	return fs.overlayFS.fs.FileExists(path)
}

// GetAccessibleEntries implements vfs.FS.
func (fs *compilerFS) GetAccessibleEntries(path string) vfs.Entries {
	return fs.overlayFS.fs.GetAccessibleEntries(path)
}

// ReadFile implements vfs.FS.
func (fs *compilerFS) ReadFile(path string) (contents string, ok bool) {
	if fh := fs.overlayFS.getFile(path); fh != nil {
		return fh.Content(), true
	}
	return "", false
}

// Realpath implements vfs.FS.
func (fs *compilerFS) Realpath(path string) string {
	return fs.overlayFS.fs.Realpath(path)
}

// Stat implements vfs.FS.
func (fs *compilerFS) Stat(path string) vfs.FileInfo {
	return fs.overlayFS.fs.Stat(path)
}

// UseCaseSensitiveFileNames implements vfs.FS.
func (fs *compilerFS) UseCaseSensitiveFileNames() bool {
	return fs.overlayFS.fs.UseCaseSensitiveFileNames()
}

// WalkDir implements vfs.FS.
func (fs *compilerFS) WalkDir(root string, walkFn vfs.WalkDirFunc) error {
	panic("unimplemented")
}

// WriteFile implements vfs.FS.
func (fs *compilerFS) WriteFile(path string, data string, writeByteOrderMark bool) error {
	panic("unimplemented")
}

// Remove implements vfs.FS.
func (fs *compilerFS) Remove(path string) error {
	panic("unimplemented")
}
