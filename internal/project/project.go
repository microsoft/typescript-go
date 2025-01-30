package project

import (
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

var _ ls.Host = (*Project)(nil)

type ProjectKind int

const (
	ProjectKindInferred ProjectKind = iota
	ProjectKindConfigured
	ProjectKindAutoImportProvider
	ProjectKindAuxiliary
)

type Project struct {
	projectService *ProjectService
	kind           ProjectKind

	dirty                     bool
	version                   int
	hasAddedOrRemovedFiles    bool
	hasAddedOrRemovedSymlinks bool
	deferredClose             bool

	configFileName string
	configFilePath tspath.Path
	// rootFileNames was a map from Path to { NormalizedPath, ScriptInfo? } in the original code.
	// But the ProjectService owns script infos, so it's not clear why there was an extra pointer.
	rootFileNames   *collections.OrderedMap[tspath.Path, string]
	compilerOptions *core.CompilerOptions
	program         *compiler.Program
}

func NewConfiguredProject(configFileName string, configFilePath tspath.Path, projectService *ProjectService) *Project {
	return &Project{
		projectService: projectService,
		kind:           ProjectKindConfigured,
		configFileName: configFileName,
		configFilePath: configFilePath,
	}
}

// FS implements LanguageServiceHost.
func (p *Project) FS() vfs.FS {
	return p.projectService.host.FS()
}

// GetCompilerOptions implements LanguageServiceHost.
func (p *Project) GetCompilerOptions() *core.CompilerOptions {
	return p.compilerOptions
}

// GetCurrentDirectory implements LanguageServiceHost.
func (p *Project) GetCurrentDirectory() string {
	return p.projectService.host.GetCurrentDirectory()
}

// GetProjectVersion implements LanguageServiceHost.
func (p *Project) GetProjectVersion() int {
	return p.version
}

// GetRootFileNames implements LanguageServiceHost.
func (p *Project) GetRootFileNames() []string {
	return slices.Collect(p.rootFileNames.Values())
}

// GetSourceFile implements LanguageServiceHost.
func (p *Project) GetSourceFile(fileName string, languageVersion core.ScriptTarget) *ast.SourceFile {
	scriptKind := p.getScriptKind(fileName)
	if scriptInfo := p.getOrCreateScriptInfoAndAttachToProject(fileName, scriptKind); scriptInfo != nil {
		oldSourceFile := p.program.GetSourceFileByPath(scriptInfo.path)
		return p.projectService.documentRegistry.AcquireDocument(scriptInfo, p.GetCompilerOptions(), oldSourceFile, p.program.GetCompilerOptions())
	}
	return nil
}

// NewLine implements LanguageServiceHost.
func (p *Project) NewLine() string {
	return p.projectService.host.NewLine()
}

// Trace implements LanguageServiceHost.
func (p *Project) Trace(msg string) {
	p.projectService.host.Trace(msg)
}

func (p *Project) getOrCreateScriptInfoAndAttachToProject(fileName string, scriptKind core.ScriptKind) *scriptInfo {
	if scriptInfo := p.projectService.getOrCreateScriptInfoNotOpenedByClient(fileName, p.toPath(fileName), scriptKind); scriptInfo != nil {
		scriptInfo.attachToProject(p)
		return scriptInfo
	}
	return nil
}

func (p *Project) getScriptKind(fileName string) core.ScriptKind {
	// Customizing script kind per file extension is a common plugin / LS host customization case
	// which can probably be replaced with static info in the future
	return core.GetScriptKindFromFileName(fileName)
}

func (p *Project) markFileAsDirty(path tspath.Path) {
	p.markAsDirty()
}

func (p *Project) markAsDirty() {
	p.dirty = true
	p.version++
}

func (p *Project) updateIfDirty() {
	// !!! p.invalidateResolutionsOfFailedLookupLocations()
	if p.dirty {
		p.updateGraph()
	}
}

func (p *Project) onFileAddedOrRemoved(isSymlink bool) {
	p.hasAddedOrRemovedFiles = true
	if isSymlink {
		p.hasAddedOrRemovedSymlinks = true
	}
}

func (p *Project) updateGraph() {
}

func (p *Project) isOrphan() bool {
	switch p.kind {
	case ProjectKindInferred:
		return p.rootFileNames.Size() == 0
	case ProjectKindConfigured:
		return p.deferredClose
	default:
		panic("unhandled project kind")
	}
}

func (p *Project) toPath(fileName string) tspath.Path {
	return tspath.ToPath(fileName, p.GetCurrentDirectory(), p.FS().UseCaseSensitiveFileNames())
}

func (p *Project) isRoot(info *scriptInfo) bool {
	return p.rootFileNames.Has(info.path)
}

func (p *Project) removeFile(info *scriptInfo, fileExists bool, detachFromProject bool) {
	if p.isRoot(info) {
		p.rootFileNames.Delete(info.path)
	}
	// !!!
	// if (fileExists) {
	// 	// If file is present, just remove the resolutions for the file
	// 	this.resolutionCache.removeResolutionsOfFile(info.path);
	// } else {
	// 	this.resolutionCache.invalidateResolutionOfFile(info.path);
	// }
	// this.cachedUnresolvedImportsPerFile.delete(info.path);
	if detachFromProject {
		info.detachFromProject(p)
	}
	p.markAsDirty()
}

func (p *Project) addMissingRootFile(fileName string, path tspath.Path) {
	p.rootFileNames.Set(path, fileName)
	p.markAsDirty()
}

func (p *Project) clearSourceMapperCache() {
	// !!!
}
