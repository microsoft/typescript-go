package projectv2

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type Kind int

const (
	KindInferred Kind = iota
	KindConfigured
)

var _ compiler.CompilerHost = (*Project)(nil)

type Project struct {
	Name string
	Kind Kind

	CommandLine   *tsoptions.ParsedCommandLine
	Program       *compiler.Program
	rootFileNames collections.OrderedMap[tspath.Path, string] // values are file names
	snapshot      *Snapshot

	currentDirectory string
}

func NewConfiguredProject(
	configFileName string,
	configFilePath tspath.Path,
	snapshot *Snapshot,
) *Project {
	return NewProject(configFileName, KindConfigured, tspath.GetDirectoryPath(configFileName), snapshot)
}

func NewProject(
	name string,
	kind Kind,
	currentDirectory string,
	snapshot *Snapshot,
) *Project {
	return &Project{
		Name:             name,
		Kind:             kind,
		snapshot:         snapshot,
		currentDirectory: currentDirectory,
	}
}

// DefaultLibraryPath implements compiler.CompilerHost.
func (p *Project) DefaultLibraryPath() string {
	return p.snapshot.sessionOptions.DefaultLibraryPath
}

// FS implements compiler.CompilerHost.
func (p *Project) FS() vfs.FS {
	return p.snapshot.compilerFS
}

// GetCurrentDirectory implements compiler.CompilerHost.
func (p *Project) GetCurrentDirectory() string {
	return p.currentDirectory
}

// GetResolvedProjectReference implements compiler.CompilerHost.
func (p *Project) GetResolvedProjectReference(fileName string, path tspath.Path) *tsoptions.ParsedCommandLine {
	panic("unimplemented")
}

// GetSourceFile implements compiler.CompilerHost. GetSourceFile increments
// the ref count of source files it acquires in the parseCache. There should
// be a corresponding release for each call made.
func (p *Project) GetSourceFile(opts ast.SourceFileParseOptions) *ast.SourceFile {
	if fh := p.snapshot.GetFile(ls.FileNameToDocumentURI(opts.FileName)); fh != nil {
		return p.snapshot.parseCache.acquireDocument(fh, opts, p.getScriptKind(opts.FileName))
	}
	return nil
}

// NewLine implements compiler.CompilerHost.
func (p *Project) NewLine() string {
	return p.snapshot.sessionOptions.NewLine
}

// Trace implements compiler.CompilerHost.
func (p *Project) Trace(msg string) {
	panic("unimplemented")
}

func (p *Project) getScriptKind(fileName string) core.ScriptKind {
	// Customizing script kind per file extension is a common plugin / LS host customization case
	// which can probably be replaced with static info in the future
	return core.GetScriptKindFromFileName(fileName)
}

func (p *Project) containsFile(path tspath.Path) bool {
	if p.isRoot(path) {
		return true
	}
	return p.Program != nil && p.Program.GetSourceFileByPath(path) != nil
}

func (p *Project) isRoot(path tspath.Path) bool {
	return p.rootFileNames.Has(path)
}
