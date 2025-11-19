package autoimport

import (
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/binder"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type failedAmbientModuleLookupSource struct {
	mu          sync.Mutex
	fileName    string
	packageName string
}

type resolver struct {
	toPath         func(fileName string) tspath.Path
	host           RegistryCloneHost
	moduleResolver *module.Resolver

	rootFiles []*ast.SourceFile

	resolvedModules                          collections.SyncMap[tspath.Path, *collections.SyncMap[module.ModeAwareCacheKey, *module.ResolvedModule]]
	possibleFailedAmbientModuleLookupTargets collections.SyncSet[string]
	possibleFailedAmbientModuleLookupSources collections.SyncMap[tspath.Path, *failedAmbientModuleLookupSource]
}

func newResolver(rootFileNames []string, host RegistryCloneHost, moduleResolver *module.Resolver, toPath func(fileName string) tspath.Path) *resolver {
	r := &resolver{
		toPath:         toPath,
		host:           host,
		moduleResolver: moduleResolver,
		rootFiles:      make([]*ast.SourceFile, 0, len(rootFileNames)),
	}
	for _, fileName := range rootFileNames {
		// !!! if we don't end up storing files in the ParseCache, this would be repeated
		r.rootFiles = append(r.rootFiles, r.GetSourceFile(fileName))
	}
	return r
}

// BindSourceFiles implements checker.Program.
func (r *resolver) BindSourceFiles() {
	// We will bind as we parse
}

// SourceFiles implements checker.Program.
func (r *resolver) SourceFiles() []*ast.SourceFile {
	return r.rootFiles
}

// Options implements checker.Program.
func (r *resolver) Options() *core.CompilerOptions {
	return &core.CompilerOptions{
		NoCheck: core.TSTrue,
	}
}

// GetCurrentDirectory implements checker.Program.
func (r *resolver) GetCurrentDirectory() string {
	return r.host.GetCurrentDirectory()
}

// UseCaseSensitiveFileNames implements checker.Program.
func (r *resolver) UseCaseSensitiveFileNames() bool {
	return r.host.FS().UseCaseSensitiveFileNames()
}

// GetSourceFile implements checker.Program.
func (r *resolver) GetSourceFile(fileName string) *ast.SourceFile {
	// !!! local cache
	file := r.host.GetSourceFile(fileName, r.toPath(fileName))
	binder.BindSourceFile(file)
	return file
}

// GetDefaultResolutionModeForFile implements checker.Program.
func (r *resolver) GetDefaultResolutionModeForFile(file ast.HasFileName) core.ResolutionMode {
	// !!!
	return core.ModuleKindESNext
}

// GetEmitModuleFormatOfFile implements checker.Program.
func (r *resolver) GetEmitModuleFormatOfFile(sourceFile ast.HasFileName) core.ModuleKind {
	return core.ModuleKindESNext
}

// GetEmitSyntaxForUsageLocation implements checker.Program.
func (r *resolver) GetEmitSyntaxForUsageLocation(sourceFile ast.HasFileName, usageLocation *ast.StringLiteralLike) core.ResolutionMode {
	return core.ModuleKindESNext
}

// GetImpliedNodeFormatForEmit implements checker.Program.
func (r *resolver) GetImpliedNodeFormatForEmit(sourceFile ast.HasFileName) core.ModuleKind {
	return core.ModuleKindESNext
}

// GetModeForUsageLocation implements checker.Program.
func (r *resolver) GetModeForUsageLocation(file ast.HasFileName, moduleSpecifier *ast.StringLiteralLike) core.ResolutionMode {
	return core.ModuleKindESNext
}

// GetResolvedModule implements checker.Program.
func (r *resolver) GetResolvedModule(currentSourceFile ast.HasFileName, moduleReference string, mode core.ResolutionMode) *module.ResolvedModule {
	cache, _ := r.resolvedModules.LoadOrStore(currentSourceFile.Path(), &collections.SyncMap[module.ModeAwareCacheKey, *module.ResolvedModule]{})
	if resolved, ok := cache.Load(module.ModeAwareCacheKey{Name: moduleReference, Mode: mode}); ok {
		return resolved
	}
	resolved, _ := r.moduleResolver.ResolveModuleName(moduleReference, currentSourceFile.FileName(), mode, nil)
	resolved, _ = cache.LoadOrStore(module.ModeAwareCacheKey{Name: moduleReference, Mode: mode}, resolved)
	// !!! failed lookup locations
	// !!! also successful lookup locations, for that matter, need to cause invalidation
	if !resolved.IsResolved() && !tspath.PathIsRelative(moduleReference) {
		r.possibleFailedAmbientModuleLookupTargets.Add(moduleReference)
		r.possibleFailedAmbientModuleLookupSources.LoadOrStore(currentSourceFile.Path(), &failedAmbientModuleLookupSource{
			fileName: currentSourceFile.FileName(),
		})
	}
	return resolved
}

// GetSourceFileForResolvedModule implements checker.Program.
func (r *resolver) GetSourceFileForResolvedModule(fileName string) *ast.SourceFile {
	return r.GetSourceFile(fileName)
}

// GetResolvedModules implements checker.Program.
func (r *resolver) GetResolvedModules() map[tspath.Path]module.ModeAwareCache[*module.ResolvedModule] {
	// only used when producing diagnostics, which hopefully the checker won't do
	return nil
}

// ---

// GetSourceFileMetaData implements checker.Program.
func (r *resolver) GetSourceFileMetaData(path tspath.Path) ast.SourceFileMetaData {
	panic("unimplemented")
}

// CommonSourceDirectory implements checker.Program.
func (r *resolver) CommonSourceDirectory() string {
	panic("unimplemented")
}

// FileExists implements checker.Program.
func (r *resolver) FileExists(fileName string) bool {
	panic("unimplemented")
}

// GetGlobalTypingsCacheLocation implements checker.Program.
func (r *resolver) GetGlobalTypingsCacheLocation() string {
	panic("unimplemented")
}

// GetImportHelpersImportSpecifier implements checker.Program.
func (r *resolver) GetImportHelpersImportSpecifier(path tspath.Path) *ast.Node {
	panic("unimplemented")
}

// GetJSXRuntimeImportSpecifier implements checker.Program.
func (r *resolver) GetJSXRuntimeImportSpecifier(path tspath.Path) (moduleReference string, specifier *ast.Node) {
	panic("unimplemented")
}

// GetNearestAncestorDirectoryWithPackageJson implements checker.Program.
func (r *resolver) GetNearestAncestorDirectoryWithPackageJson(dirname string) string {
	panic("unimplemented")
}

// GetPackageJsonInfo implements checker.Program.
func (r *resolver) GetPackageJsonInfo(pkgJsonPath string) modulespecifiers.PackageJsonInfo {
	panic("unimplemented")
}

// GetProjectReferenceFromOutputDts implements checker.Program.
func (r *resolver) GetProjectReferenceFromOutputDts(path tspath.Path) *tsoptions.SourceOutputAndProjectReference {
	panic("unimplemented")
}

// GetProjectReferenceFromSource implements checker.Program.
func (r *resolver) GetProjectReferenceFromSource(path tspath.Path) *tsoptions.SourceOutputAndProjectReference {
	panic("unimplemented")
}

// GetRedirectForResolution implements checker.Program.
func (r *resolver) GetRedirectForResolution(file ast.HasFileName) *tsoptions.ParsedCommandLine {
	panic("unimplemented")
}

// GetRedirectTargets implements checker.Program.
func (r *resolver) GetRedirectTargets(path tspath.Path) []string {
	panic("unimplemented")
}

// GetResolvedModuleFromModuleSpecifier implements checker.Program.
func (r *resolver) GetResolvedModuleFromModuleSpecifier(file ast.HasFileName, moduleSpecifier *ast.StringLiteralLike) *module.ResolvedModule {
	panic("unimplemented")
}

// GetSourceOfProjectReferenceIfOutputIncluded implements checker.Program.
func (r *resolver) GetSourceOfProjectReferenceIfOutputIncluded(file ast.HasFileName) string {
	panic("unimplemented")
}

// IsSourceFileDefaultLibrary implements checker.Program.
func (r *resolver) IsSourceFileDefaultLibrary(path tspath.Path) bool {
	panic("unimplemented")
}

// IsSourceFromProjectReference implements checker.Program.
func (r *resolver) IsSourceFromProjectReference(path tspath.Path) bool {
	panic("unimplemented")
}

// SourceFileMayBeEmitted implements checker.Program.
func (r *resolver) SourceFileMayBeEmitted(sourceFile *ast.SourceFile, forceDtsEmit bool) bool {
	panic("unimplemented")
}

var _ checker.Program = (*resolver)(nil)
