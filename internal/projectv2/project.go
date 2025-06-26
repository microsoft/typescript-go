package projectv2

import (
	"context"
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
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

type projectChange struct {
	changedURIs   []tspath.Path
	requestedURIs []struct {
		path           tspath.Path
		defaultProject *Project
	}
}

type projectChangeResult struct {
	changed bool
}

func (p *Project) Clone(ctx context.Context, change projectChange, newSnapshot *Snapshot) (*Project, projectChangeResult) {
	var result projectChangeResult
	var loadProgram bool
	// var pendingReload PendingReload
	for _, file := range change.requestedURIs {
		if file.defaultProject == p {
			loadProgram = true
			break
		}
	}

	var singleChangedFile tspath.Path
	if p.Program != nil || !loadProgram {
		for _, path := range change.changedURIs {
			if p.containsFile(path) {
				loadProgram = true
				if p.Program == nil {
					break
				} else if singleChangedFile == "" {
					singleChangedFile = path
				} else {
					singleChangedFile = ""
					break
				}
			}
		}
	}

	if loadProgram {
		result.changed = true
		newProject := &Project{
			Name:             p.Name,
			Kind:             p.Kind,
			CommandLine:      p.CommandLine,
			rootFileNames:    p.rootFileNames,
			currentDirectory: p.currentDirectory,
			snapshot:         newSnapshot,
		}

		var cloned bool
		var newProgram *compiler.Program
		oldProgram := p.Program
		if singleChangedFile != "" {
			newProgram, cloned = p.Program.UpdateProgram(singleChangedFile, newProject)
			if !cloned {
				// !!! make this less janky
				// UpdateProgram called GetSourceFile (acquiring the document) but was unable to use it directly,
				// so it called NewProgram which acquired it a second time. We need to decrement the ref count
				// for the first acquisition.
				p.snapshot.parseCache.releaseDocument(newProgram.GetSourceFileByPath(singleChangedFile))
			}
		} else {
			newProgram = compiler.NewProgram(
				compiler.ProgramOptions{
					Host:                        newProject,
					Config:                      newProject.CommandLine,
					UseSourceOfProjectReference: true,
					TypingsLocation:             newProject.snapshot.sessionOptions.TypingsLocation,
					JSDocParsingMode:            ast.JSDocParsingModeParseAll,
				},
			)
		}

		if !cloned {
			for _, file := range oldProgram.GetSourceFiles() {
				p.snapshot.parseCache.releaseDocument(file)
			}
		}

		newProject.Program = newProgram
		newProject.LanguageService = ls.NewLanguageService(ctx, newProject)
		return newProject, result
	}

	return p, result
}
