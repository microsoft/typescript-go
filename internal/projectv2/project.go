package projectv2

import (
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type Kind int

const (
	KindInferred Kind = iota
	KindConfigured
)

type PendingReload int

const (
	PendingReloadNone PendingReload = iota
	PendingReloadFileNames
	PendingReloadFull
)

var _ compiler.CompilerHost = (*Project)(nil)
var _ ls.Host = (*Project)(nil)

type Project struct {
	Name           string
	Kind           Kind
	configFileName string
	configFilePath tspath.Path

	CommandLine     *tsoptions.ParsedCommandLine
	Program         *compiler.Program
	LanguageService *ls.LanguageService
	checkerPool     *project.CheckerPool
	rootFileNames   *collections.OrderedMap[tspath.Path, string] // values are file names
	snapshot        *Snapshot

	currentDirectory string
}

func NewConfiguredProject(
	configFileName string,
	configFilePath tspath.Path,
	snapshot *Snapshot,
) *Project {
	p := NewProject(configFileName, KindConfigured, tspath.GetDirectoryPath(configFileName), snapshot)
	p.configFileName = configFileName
	p.configFilePath = configFilePath
	return p
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
		rootFileNames:    &collections.OrderedMap[tspath.Path, string]{},
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

// GetLineMap implements ls.Host.
func (p *Project) GetLineMap(fileName string) *ls.LineMap {
	// !!! cache
	return ls.ComputeLineStarts(p.snapshot.GetFile(ls.FileNameToDocumentURI(fileName)).Content())
}

// GetPositionEncoding implements ls.Host.
func (p *Project) GetPositionEncoding() lsproto.PositionEncodingKind {
	return p.snapshot.sessionOptions.PositionEncoding
}

// GetProgram implements ls.Host.
func (p *Project) GetProgram() *compiler.Program {
	return p.Program
}

func (p *Project) GetRootFileNames() []string {
	return slices.Collect(p.rootFileNames.Values())
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

func (p *Project) IsSourceFromProjectReference(path tspath.Path) bool {
	return p.Program != nil && p.Program.IsSourceFromProjectReference(path)
}

func (p *Project) Clone(newSnapshot *Snapshot) *Project {
	return &Project{
		Name:             p.Name,
		Kind:             p.Kind,
		CommandLine:      p.CommandLine,
		rootFileNames:    p.rootFileNames,
		currentDirectory: p.currentDirectory,
		snapshot:         newSnapshot,
	}
}
