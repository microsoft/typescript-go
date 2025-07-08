package projectv2

import (
	"context"
	"fmt"
	"maps"
	"sync"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
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
	sessionOptions      *SessionOptions
	parseCache          *parseCache
	extendedConfigCache *extendedConfigCache

	ctx                                context.Context
	fs                                 *overlayFS
	base                               *ProjectCollection
	compilerOptionsForInferredProjects *core.CompilerOptions
	configFileRegistryBuilder          *configFileRegistryBuilder

	fileDefaultProjects map[tspath.Path]tspath.Path
	// Keys are file paths, values are sets of project paths that contain the file.
	fileAssociations   collections.SyncMap[tspath.Path, *collections.SyncSet[tspath.Path]]
	configuredProjects collections.SyncMap[tspath.Path, *Project]
	inferredProject    *inferredProjectEntry
}

func newProjectCollectionBuilder(
	ctx context.Context,
	fs *overlayFS,
	oldProjectCollection *ProjectCollection,
	oldConfigFileRegistry *ConfigFileRegistry,
	compilerOptionsForInferredProjects *core.CompilerOptions,
	sessionOptions *SessionOptions,
	parseCache *parseCache,
	extendedConfigCache *extendedConfigCache,
) *projectCollectionBuilder {
	return &projectCollectionBuilder{
		ctx:                                ctx,
		fs:                                 fs,
		compilerOptionsForInferredProjects: compilerOptionsForInferredProjects,
		sessionOptions:                     sessionOptions,
		parseCache:                         parseCache,
		extendedConfigCache:                extendedConfigCache,
		base:                               oldProjectCollection,
		configFileRegistryBuilder:          newConfigFileRegistryBuilder(fs, oldConfigFileRegistry, extendedConfigCache, sessionOptions),
	}
}

func (b *projectCollectionBuilder) Finalize() (*ProjectCollection, *ConfigFileRegistry) {
	var changed bool
	newProjectCollection := b.base
	b.configuredProjects.Range(func(path tspath.Path, project *Project) bool {
		if !changed {
			newProjectCollection = newProjectCollection.clone()
			if newProjectCollection.configuredProjects == nil {
				newProjectCollection.configuredProjects = make(map[tspath.Path]*Project)
			} else {
				newProjectCollection.configuredProjects = maps.Clone(newProjectCollection.configuredProjects)
			}
			changed = true
		}
		newProjectCollection.configuredProjects[path] = project
		return true
	})

	if !changed && !maps.Equal(b.fileDefaultProjects, b.base.fileDefaultProjects) {
		newProjectCollection = newProjectCollection.clone()
		newProjectCollection.fileDefaultProjects = b.fileDefaultProjects
		changed = true
	} else if changed {
		newProjectCollection.fileDefaultProjects = b.fileDefaultProjects
	}

	if b.inferredProject != nil {
		if !changed {
			newProjectCollection = newProjectCollection.clone()
		}
		newProjectCollection.inferredProject = b.inferredProject.project
	}

	// !!! clean up file associations of deleted projects, deleted files
	var fileAssociationsChanged bool
	b.fileAssociations.Range(func(filePath tspath.Path, projectPaths *collections.SyncSet[tspath.Path]) bool {
		if !changed {
			newProjectCollection = newProjectCollection.clone()
			changed = true
		}
		if !fileAssociationsChanged {
			if newProjectCollection.fileAssociations == nil {
				newProjectCollection.fileAssociations = make(map[tspath.Path]map[tspath.Path]struct{})
			} else {
				newProjectCollection.fileAssociations = maps.Clone(newProjectCollection.fileAssociations)
			}
			fileAssociationsChanged = true
		}
		m, ok := newProjectCollection.fileAssociations[filePath]
		if !ok {
			m = make(map[tspath.Path]struct{})
			newProjectCollection.fileAssociations[filePath] = m
		}
		projectPaths.Range(func(projectPath tspath.Path) bool {
			m[projectPath] = struct{}{}
			return true
		})
		return true
	})

	return newProjectCollection, b.configFileRegistryBuilder.finalize()
}

func (b *projectCollectionBuilder) loadOrStoreNewConfiguredProject(
	fileName string,
	path tspath.Path,
) (*projectCollectionBuilderEntry, bool) {
	// Check for existence in the base registry first so that all SyncMap
	// access is atomic. We're trying to avoid the scenario where we
	//   1. try to load from the dirty map but find nothing,
	//   2. try to load from the base registry but find nothing, then
	//   3. have to do a subsequent Store in the dirty map for the new entry.
	if prev, ok := b.base.configuredProjects[path]; ok {
		if dirty, ok := b.configuredProjects.Load(path); ok {
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
		entry, loaded := b.configuredProjects.LoadOrStore(path, NewConfiguredProject(fileName, path, b))
		return &projectCollectionBuilderEntry{
			b:       b,
			project: entry,
			dirty:   true,
		}, loaded
	}
}

func (b *projectCollectionBuilder) getConfiguredProject(path tspath.Path) (*projectCollectionBuilderEntry, bool) {
	if entry, ok := b.configuredProjects.Load(path); ok {
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

func (b *projectCollectionBuilder) forEachConfiguredProject(fn func(entry *projectCollectionBuilderEntry) bool) {
	seenDirty := make(map[tspath.Path]struct{})
	b.configuredProjects.Range(func(path tspath.Path, project *Project) bool {
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

func (b *projectCollectionBuilder) forEachProject(fn func(entry *projectCollectionBuilderEntry) bool) {
	var keepGoing bool
	b.forEachConfiguredProject(func(entry *projectCollectionBuilderEntry) bool {
		keepGoing = fn(entry)
		return keepGoing
	})
	if !keepGoing {
		return
	}
	inferredProject := b.getInferredProject()
	if inferredProject.project != nil {
		fn((*projectCollectionBuilderEntry)(inferredProject))
	}
}

func (b *projectCollectionBuilder) getInferredProject() *inferredProjectEntry {
	if b.inferredProject != nil {
		return b.inferredProject
	}
	return &inferredProjectEntry{
		b:       b,
		project: b.base.inferredProject,
		dirty:   false,
	}
}

func (b *projectCollectionBuilder) DidOpenFile(uri lsproto.DocumentUri) {
	fileName := uri.FileName()
	path := b.toPath(fileName)
	_ = b.tryFindDefaultConfiguredProjectAndLoadAncestorsForOpenScriptInfo(fileName, path, projectLoadKindCreate)
	b.forEachProject(func(entry *projectCollectionBuilderEntry) bool {
		entry.updateProgram()
		return true
	})
	if b.findDefaultProject(fileName, path) == nil {
		b.getInferredProject().addFile(fileName, path)
	}
}

func (b *projectCollectionBuilder) DidChangeFiles(uris []lsproto.DocumentUri) {
	paths := core.Map(uris, func(uri lsproto.DocumentUri) tspath.Path {
		return uri.Path(b.fs.fs.UseCaseSensitiveFileNames())
	})
	b.forEachProject(func(entry *projectCollectionBuilderEntry) bool {
		for _, path := range paths {
			entry.markFileChanged(path)
		}
		return true
	})
}

func (b *projectCollectionBuilder) DidRequestFile(uri lsproto.DocumentUri) {
	// See if we can find a default project for this file without doing
	// any additional loading.
	fileName := uri.FileName()
	path := b.toPath(fileName)
	if result := b.findDefaultProject(fileName, path); result != nil {
		result.updateProgram()
		return
	}

	// Make sure all projects we know about are up to date...
	var hasChanges bool
	b.forEachConfiguredProject(func(entry *projectCollectionBuilderEntry) bool {
		hasChanges = entry.updateProgram() || hasChanges
		return true
	})
	if hasChanges {
		// If the structure of other projects changed, we might need to move files
		// in/out of the inferred project.
		var inferredProjectFiles []string
		for path, overlay := range b.fs.overlays {
			if b.findDefaultConfiguredProject(overlay.FileName(), path) == nil {
				inferredProjectFiles = append(inferredProjectFiles, overlay.FileName())
			}
		}
		if len(inferredProjectFiles) > 0 {
			inferredProject := b.getInferredProject()
			inferredProject.updateInferredProject(inferredProjectFiles)
		}
	}

	// ...and then try to find the default configured project for this file again.
	if b.findDefaultProject(fileName, path) == nil {
		panic(fmt.Sprintf("no project found for file %s", fileName))
	}
}

func (b *projectCollectionBuilder) findOrCreateProject(
	configFileName string,
	configFilePath tspath.Path,
	loadKind projectLoadKind,
) *projectCollectionBuilderEntry {
	if loadKind == projectLoadKindFind {
		entry, _ := b.getConfiguredProject(configFilePath)
		return entry
	}
	entry, _ := b.loadOrStoreNewConfiguredProject(configFileName, configFilePath)
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
			childConfigFilePath := b.toPath(childConfigFileName)
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
	configFilePath := b.toPath(configFileName)
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
		if key == "" {
			// The file belongs to the inferred project
			return nil
		}
		entry, _ := b.getConfiguredProject(key)
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

func (b *projectCollectionBuilder) findDefaultProject(fileName string, path tspath.Path) *projectCollectionBuilderEntry {
	if configuredProject := b.findDefaultConfiguredProject(fileName, path); configuredProject != nil {
		return configuredProject
	}
	if key, ok := b.fileDefaultProjects[path]; ok && key == "" {
		return (*projectCollectionBuilderEntry)(b.getInferredProject())
	}
	if inferredProject := b.getInferredProject(); inferredProject != nil && inferredProject.project != nil && inferredProject.project.containsFile(path) {
		if b.fileDefaultProjects == nil {
			b.fileDefaultProjects = make(map[tspath.Path]tspath.Path)
		}
		b.fileDefaultProjects[path] = ""
		return (*projectCollectionBuilderEntry)(inferredProject)
	}
	return nil
}

func (b *projectCollectionBuilder) findDefaultConfiguredProject(fileName string, path tspath.Path) *projectCollectionBuilderEntry {
	if b.isOpenFile(path) {
		result := b.tryFindDefaultConfiguredProjectForOpenScriptInfo(fileName, path, projectLoadKindFind)
		if result != nil && result.project != nil {
			return result.project
		}
	}
	return nil
}

func (b *projectCollectionBuilder) toPath(fileName string) tspath.Path {
	return tspath.ToPath(fileName, b.sessionOptions.CurrentDirectory, b.fs.fs.UseCaseSensitiveFileNames())
}

func (b *projectCollectionBuilder) isOpenFile(path tspath.Path) bool {
	_, ok := b.fs.overlays[path]
	return ok
}

type projectCollectionBuilderEntry struct {
	b       *projectCollectionBuilder
	project *Project
	dirty   bool
}

type inferredProjectEntry projectCollectionBuilderEntry

func (e *inferredProjectEntry) updateInferredProject(rootFileNames []string) bool {
	if e.project == nil && len(rootFileNames) > 0 {
		e.project = NewInferredProject(e.b.sessionOptions.CurrentDirectory, e.b.compilerOptionsForInferredProjects, rootFileNames, e.b)
		e.dirty = true
		e.b.inferredProject = e
	} else if e.project != nil && len(rootFileNames) == 0 {
		e.project = nil
		e.dirty = true
		e.b.inferredProject = e
		return true
	} else {
		newCommandLine := tsoptions.NewParsedCommandLine(e.b.compilerOptionsForInferredProjects, rootFileNames, tspath.ComparePathsOptions{
			UseCaseSensitiveFileNames: e.b.fs.fs.UseCaseSensitiveFileNames(),
			CurrentDirectory:          e.project.currentDirectory,
		})
		if maps.Equal(e.project.CommandLine.FileNamesByPath(), newCommandLine.FileNamesByPath()) {
			return false
		}
		(*projectCollectionBuilderEntry)(e).ensureProjectCloned()
		e.project.CommandLine = newCommandLine
		e.project.dirty = true
		e.project.dirtyFilePath = ""
	}
	return (*projectCollectionBuilderEntry)(e).updateProgram()
}

func (e *inferredProjectEntry) addFile(fileName string, path tspath.Path) bool {
	if e.project == nil {
		return e.updateInferredProject([]string{fileName})
	}
	return e.updateInferredProject(append(e.project.CommandLine.FileNames(), fileName))
}

func (e *projectCollectionBuilderEntry) updateProgram() bool {
	updateProgram := e.project.dirty
	if e.project.Kind == KindConfigured {
		commandLine := e.b.configFileRegistryBuilder.acquireConfigForProject(e.project.configFileName, e.project.configFilePath, e.project)
		if e.project.CommandLine != commandLine {
			e.ensureProjectCloned()
			e.project.CommandLine = commandLine
			updateProgram = true
		}
	}
	if !updateProgram {
		return false
	}

	e.ensureProjectCloned()
	e.project.host = newCompilerHost(e.project.currentDirectory, e.project, e.b)
	newProgram, checkerPool := e.project.CreateProgram()
	e.project.Program = newProgram
	e.project.checkerPool = checkerPool
	// !!! unthread context
	e.project.LanguageService = ls.NewLanguageService(e.b.ctx, e.project)
	e.project.dirty = false
	e.project.dirtyFilePath = ""
	return true
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
		e.project = e.project.Clone()
		e.dirty = true
		if e.project.Kind == KindInferred {
			e.b.inferredProject = (*inferredProjectEntry)(e)
		} else {
			e.b.configuredProjects.Store(e.project.configFilePath, e.project)
		}
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
