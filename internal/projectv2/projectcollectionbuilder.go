package projectv2

import (
	"context"
	"maps"
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type projectLoadKind int

const (
	// Project is not created or updated, only looked up in cache
	projectLoadKindFind projectLoadKind = iota
	// Project is created and then its graph is updated
	projectLoadKindCreate
)

type projectCollectionBuilder struct {
	ctx                       context.Context
	snapshot                  *Snapshot
	configFileRegistryBuilder *configFileRegistryBuilder
	base                      *ProjectCollection
	dirty                     collections.SyncMap[tspath.Path, *Project]
	fileDefaultProjects       map[tspath.Path]tspath.Path
}

func newProjectCollectionBuilder(
	ctx context.Context,
	newSnapshot *Snapshot,
	oldProjectCollection *ProjectCollection,
	oldConfigFileRegistry *ConfigFileRegistry,
) *projectCollectionBuilder {
	return &projectCollectionBuilder{
		ctx:                       ctx,
		snapshot:                  newSnapshot,
		base:                      oldProjectCollection,
		configFileRegistryBuilder: newConfigFileRegistryBuilder(newSnapshot, oldConfigFileRegistry),
	}
}

func (b *projectCollectionBuilder) finalize() (*ProjectCollection, *ConfigFileRegistry) {
	var changed bool
	newProjectCollection := b.base
	b.dirty.Range(func(path tspath.Path, project *Project) bool {
		if !changed {
			newProjectCollection = newProjectCollection.clone()
			if newProjectCollection.configuredProjects == nil {
				newProjectCollection.configuredProjects = make(map[tspath.Path]*Project)
			}
			changed = true
		}
		newProjectCollection.configuredProjects[path] = project
		return true
	})
	if !changed && !maps.Equal(b.fileDefaultProjects, b.base.fileDefaultProjects) {
		newProjectCollection = newProjectCollection.clone()
		newProjectCollection.fileDefaultProjects = b.fileDefaultProjects
	} else if changed {
		newProjectCollection.fileDefaultProjects = b.fileDefaultProjects
	}
	return newProjectCollection, b.configFileRegistryBuilder.finalize()
}

func (b *projectCollectionBuilder) loadOrStoreNewEntry(
	fileName string,
	path tspath.Path,
) (*projectCollectionBuilderEntry, bool) {
	// Check for existence in the base registry first so that all SyncMap
	// access is atomic. We're trying to avoid the scenario where we
	//   1. try to load from the dirty map but find nothing,
	//   2. try to load from the base registry but find nothing, then
	//   3. have to do a subsequent Store in the dirty map for the new entry.
	if prev, ok := b.base.configuredProjects[path]; ok {
		if dirty, ok := b.dirty.Load(path); ok {
			return &projectCollectionBuilderEntry{
				b:       b,
				project: dirty,
				dirty:   true,
			}, true
		}
		return &projectCollectionBuilderEntry{
			b:       b,
			project: prev,
			dirty:   false,
		}, true
	} else {
		entry, loaded := b.dirty.LoadOrStore(path, NewConfiguredProject(fileName, path, b.snapshot))
		return &projectCollectionBuilderEntry{
			b:       b,
			project: entry,
			dirty:   true,
		}, loaded
	}
}

func (b *projectCollectionBuilder) load(path tspath.Path) (*projectCollectionBuilderEntry, bool) {
	if entry, ok := b.dirty.Load(path); ok {
		return &projectCollectionBuilderEntry{
			b:       b,
			project: entry,
			dirty:   true,
		}, true
	}
	if entry, ok := b.base.configuredProjects[path]; ok {
		return &projectCollectionBuilderEntry{
			b:       b,
			project: entry,
			dirty:   false,
		}, true
	}
	return nil, false
}

func (b *projectCollectionBuilder) forEachProject(fn func(entry *projectCollectionBuilderEntry) bool) {
	seenDirty := make(map[tspath.Path]struct{})
	b.dirty.Range(func(path tspath.Path, project *Project) bool {
		entry := &projectCollectionBuilderEntry{
			b:       b,
			project: project,
			dirty:   true,
		}
		seenDirty[path] = struct{}{}
		return fn(entry)
	})
	for path, project := range b.base.configuredProjects {
		if _, ok := seenDirty[path]; !ok {
			entry := &projectCollectionBuilderEntry{
				b:       b,
				project: project,
				dirty:   false,
			}
			if !fn(entry) {
				return
			}
		}
	}
}

func (b *projectCollectionBuilder) markFilesChanged(uris []lsproto.DocumentUri) {
	paths := core.Map(uris, func(uri lsproto.DocumentUri) tspath.Path {
		return uri.Path(b.snapshot.compilerFS.UseCaseSensitiveFileNames())
	})
	b.forEachProject(func(entry *projectCollectionBuilderEntry) bool {
		for _, path := range paths {
			entry.markFileChanged(path)
		}
		return true
	})
}

func (b *projectCollectionBuilder) ensureDefaultProjectForFile(fileName string, path tspath.Path) {
	// See if we can find a default configured project for this file without doing
	// any additional loading.
	if result := b.findDefaultConfiguredProject(fileName, path); result != nil {
		result.updateProgram()
		return
	}

	// Make sure all projects we know about are up to date...
	b.forEachProject(func(entry *projectCollectionBuilderEntry) bool {
		entry.updateProgram()
		return true
	})

	// ...and then try to find the default configured project for this file again.
	if result := b.findDefaultConfiguredProject(fileName, path); result != nil {
		return
	}

	// If we still can't find a default project, create an inferred project for this file.
	// !!!
}

func (b *projectCollectionBuilder) findOrCreateProject(
	configFileName string,
	configFilePath tspath.Path,
	loadKind projectLoadKind,
) *projectCollectionBuilderEntry {
	if loadKind == projectLoadKindFind {
		entry, _ := b.load(configFilePath)
		return entry
	}
	entry, _ := b.loadOrStoreNewEntry(configFileName, configFilePath)
	return entry
}

func (b *projectCollectionBuilder) isDefaultConfigForScript(
	scriptFileName string,
	scriptPath tspath.Path,
	configFileName string,
	configFilePath tspath.Path,
	config *tsoptions.ParsedCommandLine,
	loadKind projectLoadKind,
	result *openScriptInfoProjectResult,
) bool {
	// This currently happens only when finding project for open script info first time file is opened
	// Set seen based on project if present of for config file if its not yet created
	if !result.addSeenConfig(configFilePath, loadKind) {
		return false
	}

	// If the file is listed in root files, then only we can use this project as default project
	if !config.MatchesFileName(scriptFileName) {
		return false
	}

	// Ensure the project is uptodate and created since the file may belong to this project
	project := b.findOrCreateProject(configFileName, configFilePath, loadKind)
	return b.isDefaultProject(scriptFileName, scriptPath, project, loadKind, result)
}

func (b *projectCollectionBuilder) isDefaultProject(
	fileName string,
	path tspath.Path,
	entry *projectCollectionBuilderEntry,
	loadKind projectLoadKind,
	result *openScriptInfoProjectResult,
) bool {
	if entry == nil {
		return false
	}

	// Skip already looked up projects
	if !result.addSeenProject(entry.project, loadKind) {
		return false
	}
	// Make sure project is upto date when in create mode
	if loadKind == projectLoadKindCreate {
		entry.updateProgram()
	}
	// If script info belongs to this project, use this as default config project
	if entry.project.containsFile(path) {
		if !entry.project.IsSourceFromProjectReference(path) {
			result.setProject(entry)
			return true
		} else if !result.hasFallbackDefault() {
			// Use this project as default if no other project is found
			result.setFallbackDefault(entry)
		}
	}
	return false
}

func (b *projectCollectionBuilder) tryFindDefaultConfiguredProjectFromReferences(
	fileName string,
	path tspath.Path,
	config *tsoptions.ParsedCommandLine,
	loadKind projectLoadKind,
	result *openScriptInfoProjectResult,
) bool {
	if len(config.ProjectReferences()) == 0 {
		return false
	}
	wg := core.NewWorkGroup(false)
	b.tryFindDefaultConfiguredProjectFromReferencesWorker(fileName, path, config, loadKind, result, wg)
	wg.RunAndWait()
	return result.isDone()
}

func (b *projectCollectionBuilder) tryFindDefaultConfiguredProjectFromReferencesWorker(
	fileName string,
	path tspath.Path,
	config *tsoptions.ParsedCommandLine,
	loadKind projectLoadKind,
	result *openScriptInfoProjectResult,
	wg core.WorkGroup,
) {
	if config.CompilerOptions().DisableReferencedProjectLoad.IsTrue() {
		loadKind = projectLoadKindFind
	}
	for _, childConfigFileName := range config.ResolvedProjectReferencePaths() {
		wg.Queue(func() {
			childConfigFilePath := b.snapshot.toPath(childConfigFileName)
			childConfig := b.configFileRegistryBuilder.findOrAcquireConfigForOpenFile(childConfigFileName, childConfigFilePath, path, loadKind)
			if childConfig == nil || b.isDefaultConfigForScript(fileName, path, childConfigFileName, childConfigFilePath, childConfig, loadKind, result) {
				return
			}
			// Search in references if we cant find default project in current config
			b.tryFindDefaultConfiguredProjectFromReferencesWorker(fileName, path, childConfig, loadKind, result, wg)
		})
	}
}

func (b *projectCollectionBuilder) tryFindDefaultConfiguredProjectFromAncestor(
	fileName string,
	path tspath.Path,
	configFileName string,
	config *tsoptions.ParsedCommandLine,
	loadKind projectLoadKind,
	result *openScriptInfoProjectResult,
) bool {
	if config != nil && config.CompilerOptions().DisableSolutionSearching.IsTrue() {
		return false
	}
	if ancestorConfigName := b.configFileRegistryBuilder.getAncestorConfigFileName(fileName, path, configFileName, loadKind); ancestorConfigName != "" {
		return b.tryFindDefaultConfiguredProjectForScriptInfo(fileName, path, ancestorConfigName, loadKind, result)
	}
	return false
}

func (b *projectCollectionBuilder) tryFindDefaultConfiguredProjectForScriptInfo(
	fileName string,
	path tspath.Path,
	configFileName string,
	loadKind projectLoadKind,
	result *openScriptInfoProjectResult,
) bool {
	// Lookup from parsedConfig if available
	configFilePath := b.snapshot.toPath(configFileName)
	config := b.configFileRegistryBuilder.findOrAcquireConfigForOpenFile(configFileName, configFilePath, path, loadKind)
	if config != nil {
		if config.CompilerOptions().Composite == core.TSTrue {
			if b.isDefaultConfigForScript(fileName, path, configFileName, configFilePath, config, loadKind, result) {
				return true
			}
		} else if len(config.FileNames()) > 0 {
			project := b.findOrCreateProject(configFileName, configFilePath, loadKind)
			if b.isDefaultProject(fileName, path, project, loadKind, result) {
				return true
			}
		}
		// Lookup in references
		if b.tryFindDefaultConfiguredProjectFromReferences(fileName, path, config, loadKind, result) {
			return true
		}
	}
	// Lookup in ancestor projects
	if b.tryFindDefaultConfiguredProjectFromAncestor(fileName, path, configFileName, config, loadKind, result) {
		return true
	}
	return false
}

func (b *projectCollectionBuilder) tryFindDefaultConfiguredProjectForOpenScriptInfo(
	fileName string,
	path tspath.Path,
	loadKind projectLoadKind,
) *openScriptInfoProjectResult {
	if key, ok := b.fileDefaultProjects[path]; ok {
		entry, _ := b.load(key)
		return &openScriptInfoProjectResult{
			project: entry,
		}
	}
	if configFileName := b.configFileRegistryBuilder.getConfigFileNameForFile(fileName, path, loadKind); configFileName != "" {
		var result openScriptInfoProjectResult
		b.tryFindDefaultConfiguredProjectForScriptInfo(fileName, path, configFileName, loadKind, &result)
		if result.project == nil && result.fallbackDefault != nil {
			result.setProject(result.fallbackDefault)
		}
		if b.fileDefaultProjects == nil {
			b.fileDefaultProjects = make(map[tspath.Path]tspath.Path)
		}
		b.fileDefaultProjects[path] = result.project.project.configFilePath
		return &result
	}
	return nil
}

func (b *projectCollectionBuilder) tryFindDefaultConfiguredProjectAndLoadAncestorsForOpenScriptInfo(
	fileName string,
	path tspath.Path,
	loadKind projectLoadKind,
) *openScriptInfoProjectResult {
	result := b.tryFindDefaultConfiguredProjectForOpenScriptInfo(fileName, path, loadKind)
	if result != nil && result.project != nil {
		// !!! sheetal todo this later
		// // Create ancestor tree for findAllRefs (dont load them right away)
		// forEachAncestorProjectLoad(
		// 	info,
		// 	tsconfigProject!,
		// 	ancestor => {
		// 		seenProjects.set(ancestor.project, kind);
		// 	},
		// 	kind,
		// 	`Creating project possibly referencing default composite project ${defaultProject.getProjectName()} of open file ${info.fileName}`,
		// 	allowDeferredClosed,
		// 	reloadedProjects,
		// 	/*searchOnlyPotentialSolution*/ true,
		// 	delayReloadedConfiguredProjects,
		// );
	}
	return result
}

func (b *projectCollectionBuilder) findDefaultConfiguredProject(fileName string, path tspath.Path) *projectCollectionBuilderEntry {
	if b.snapshot.IsOpenFile(path) {
		result := b.tryFindDefaultConfiguredProjectForOpenScriptInfo(fileName, path, projectLoadKindFind)
		if result != nil && result.project != nil /* !!! && !result.project.deferredClose */ {
			return result.project
		}
	}
	return nil
}

type projectCollectionBuilderEntry struct {
	b       *projectCollectionBuilder
	project *Project
	dirty   bool
}

func (e *projectCollectionBuilderEntry) updateProgram() {
	loadProgram := e.project.dirty
	commandLine := e.b.configFileRegistryBuilder.acquireConfigForProject(e.project.configFileName, e.project.configFilePath, e.project)
	if e.project.CommandLine != commandLine {
		e.ensureProjectCloned()
		e.project.CommandLine = commandLine
		loadProgram = true
	}

	if loadProgram {
		oldProgram := e.project.Program
		e.ensureProjectCloned()
		e.project.CommandLine = commandLine
		var programCloned bool
		var newProgram *compiler.Program
		if e.project.dirtyFilePath != "" {
			newProgram, programCloned = e.project.Program.UpdateProgram(e.project.dirtyFilePath, e.project)
			if !programCloned {
				// !!! wait until accepting snapshot to release documents!
				// !!! make this less janky
				// UpdateProgram called GetSourceFile (acquiring the document) but was unable to use it directly,
				// so it called NewProgram which acquired it a second time. We need to decrement the ref count
				// for the first acquisition.
				e.b.snapshot.parseCache.releaseDocument(newProgram.GetSourceFileByPath(e.project.dirtyFilePath))
			}
		} else {
			newProgram = compiler.NewProgram(
				compiler.ProgramOptions{
					Host:                        e.project,
					Config:                      e.project.CommandLine,
					UseSourceOfProjectReference: true,
					TypingsLocation:             e.project.snapshot.sessionOptions.TypingsLocation,
					JSDocParsingMode:            ast.JSDocParsingModeParseAll,
					CreateCheckerPool: func(program *compiler.Program) compiler.CheckerPool {
						e.project.checkerPool = project.NewCheckerPool(4, program, e.b.snapshot.Log)
						return e.project.checkerPool
					},
				},
			)
		}

		if !programCloned && oldProgram != nil {
			for _, file := range oldProgram.GetSourceFiles() {
				// !!! wait until accepting snapshot to release documents!
				e.b.snapshot.parseCache.releaseDocument(file)
			}
		}

		e.project.Program = newProgram
		// !!! unthread context
		e.project.LanguageService = ls.NewLanguageService(e.b.ctx, e.project)
		e.project.dirty = false
		e.project.dirtyFilePath = ""
	}
}

func (e *projectCollectionBuilderEntry) markFileChanged(path tspath.Path) {
	if e.project.containsFile(path) {
		e.ensureProjectCloned()
		if !e.project.dirty {
			e.project.dirty = true
			e.project.dirtyFilePath = path
		} else if e.project.dirtyFilePath != path {
			e.project.dirtyFilePath = ""
		}
	}
}

func (e *projectCollectionBuilderEntry) ensureProjectCloned() {
	if !e.dirty {
		e.project = e.project.Clone(e.b.snapshot)
		e.dirty = true
		e.b.dirty.Store(e.project.configFilePath, e.project)
	}
}

type openScriptInfoProjectResult struct {
	projectMu         sync.RWMutex
	project           *projectCollectionBuilderEntry // use this if we found actual project
	fallbackDefaultMu sync.RWMutex
	fallbackDefault   *projectCollectionBuilderEntry // use this if we cant find actual project
	seenProjects      collections.SyncMap[tspath.Path, projectLoadKind]
	seenConfigs       collections.SyncMap[tspath.Path, projectLoadKind]
}

func (r *openScriptInfoProjectResult) addSeenProject(project *Project, loadKind projectLoadKind) bool {
	if kind, loaded := r.seenProjects.LoadOrStore(project.configFilePath, loadKind); loaded {
		if kind >= loadKind {
			return false
		}
		r.seenProjects.Store(project.configFilePath, loadKind)
	}
	return true
}

func (r *openScriptInfoProjectResult) addSeenConfig(configPath tspath.Path, loadKind projectLoadKind) bool {
	if kind, loaded := r.seenConfigs.LoadOrStore(configPath, loadKind); loaded {
		if kind >= loadKind {
			return false
		}
		r.seenConfigs.Store(configPath, loadKind)
	}
	return true
}

func (r *openScriptInfoProjectResult) isDone() bool {
	r.projectMu.RLock()
	defer r.projectMu.RUnlock()
	return r.project != nil
}

func (r *openScriptInfoProjectResult) setProject(entry *projectCollectionBuilderEntry) {
	r.projectMu.Lock()
	defer r.projectMu.Unlock()
	if r.project == nil {
		r.project = entry
	}
}

func (r *openScriptInfoProjectResult) hasFallbackDefault() bool {
	r.fallbackDefaultMu.RLock()
	defer r.fallbackDefaultMu.RUnlock()
	return r.fallbackDefault != nil
}

func (r *openScriptInfoProjectResult) setFallbackDefault(entry *projectCollectionBuilderEntry) {
	r.fallbackDefaultMu.Lock()
	defer r.fallbackDefaultMu.Unlock()
	if r.fallbackDefault == nil {
		r.fallbackDefault = entry
	}
}
