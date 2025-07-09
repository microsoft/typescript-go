package projectv2

import (
	"context"
	"fmt"
	"maps"

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
	_ = b.findOrLoadDefaultConfiguredProjectAndLoadAncestorsForOpenFile(fileName, path, projectLoadKindCreate)
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
		return b.tryFindDefaultConfiguredProjectForOpenScriptInfo(fileName, path, projectLoadKindFind)
	}
	return nil
}

func (b *projectCollectionBuilder) findOrLoadDefaultConfiguredProjectAndLoadAncestorsForOpenFile(
	fileName string,
	path tspath.Path,
	loadKind projectLoadKind,
) *projectCollectionBuilderEntry {
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

type searchNode struct {
	configFileName string
	loadKind       projectLoadKind
}

func (b *projectCollectionBuilder) findOrLoadDefaultConfiguredProjectWorker(
	fileName string,
	path tspath.Path,
	configFileName string,
	loadKind projectLoadKind,
	visited *collections.SyncSet[searchNode],
	fallback *searchNode,
) *projectCollectionBuilderEntry {
	var configs collections.SyncMap[tspath.Path, *tsoptions.ParsedCommandLine]
	if visited == nil {
		visited = &collections.SyncSet[searchNode]{}
	}

	search := core.BreadthFirstSearchParallel(
		searchNode{configFileName: configFileName, loadKind: loadKind},
		func(node searchNode) []searchNode {
			if config, ok := configs.Load(b.toPath(node.configFileName)); ok && len(config.ProjectReferences()) > 0 {
				referenceLoadKind := node.loadKind
				if config.CompilerOptions().DisableReferencedProjectLoad.IsTrue() {
					referenceLoadKind = projectLoadKindFind
				}
				return core.Map(config.ResolvedProjectReferencePaths(), func(configFileName string) searchNode {
					return searchNode{configFileName: configFileName, loadKind: referenceLoadKind}
				})
			}
			return nil
		},
		func(node searchNode) (isResult bool, stop bool) {
			if node.loadKind == projectLoadKindFind && visited.Has(searchNode{configFileName: node.configFileName, loadKind: projectLoadKindCreate}) {
				// We're being asked to find when we've already been asked to create, so we can skip this node.
				// The create search node will have returned the same result we'd find here. (Note that if we
				// cared about the returned search path being determinstic, we would need to figure out whether
				// to return true or false here, but since we only care about the destination node, we can
				// just return false.)
				return false, false
			}
			configFilePath := b.toPath(node.configFileName)
			config := b.configFileRegistryBuilder.findOrAcquireConfigForOpenFile(node.configFileName, configFilePath, path, node.loadKind)
			if config == nil {
				return false, false
			}
			configs.Store(configFilePath, config)
			if len(config.FileNames()) == 0 {
				// Likely a solution tsconfig.json - the search will fan out to its references.
				return false, false
			}

			if config.CompilerOptions().Composite == core.TSTrue {
				// For composite projects, we can get an early negative result.
				// !!! what about declaration files in node_modules? wouldn't it be better to
				//     check project inclusion if the project is already loaded?
				if !config.MatchesFileName(fileName) {
					return false, false
				}
			}

			project := b.findOrCreateProject(node.configFileName, configFilePath, node.loadKind)
			if node.loadKind == projectLoadKindCreate {
				// Ensure project is up to date before checking for file inclusion
				project.updateProgram()
			}

			if project.project.containsFile(path) {
				return true, !project.project.IsSourceFromProjectReference(path)
			}

			return false, false
		},
		visited,
	)

	if search.Stopped {
		project, _ := b.getConfiguredProject(b.toPath(search.Path[0].configFileName))
		return project
	}
	if len(search.Path) > 0 {
		// If we found a project that contains the file, but it is a source from
		// a project reference, record it as a fallback.
		fallback = &search.Path[0]
	}

	// Look for tsconfig.json files higher up the directory tree and do the same. This handles
	// the common case where a higher-level "solution" tsconfig.json contains all projects in a
	// workspace.
	if config, ok := configs.Load(b.toPath(configFileName)); ok && config.CompilerOptions().DisableSolutionSearching.IsTrue() {
		if fallback != nil {
			project, _ := b.getConfiguredProject(b.toPath(fallback.configFileName))
			return project
		}
	}
	if ancestorConfigName := b.configFileRegistryBuilder.getAncestorConfigFileName(fileName, path, configFileName, loadKind); ancestorConfigName != "" {
		return b.findOrLoadDefaultConfiguredProjectWorker(fileName, path, ancestorConfigName, loadKind, visited, fallback)
	}
	if fallback != nil {
		project, _ := b.getConfiguredProject(b.toPath(fallback.configFileName))
		return project
	}
	return nil
}

func (b *projectCollectionBuilder) tryFindDefaultConfiguredProjectForOpenScriptInfo(
	fileName string,
	path tspath.Path,
	loadKind projectLoadKind,
) *projectCollectionBuilderEntry {
	if key, ok := b.fileDefaultProjects[path]; ok {
		if key == "" {
			// The file belongs to the inferred project
			return nil
		}
		entry, _ := b.getConfiguredProject(key)
		return entry
	}
	if configFileName := b.configFileRegistryBuilder.getConfigFileNameForFile(fileName, path, loadKind); configFileName != "" {
		project := b.findOrLoadDefaultConfiguredProjectWorker(fileName, path, configFileName, loadKind, nil, nil)
		if b.fileDefaultProjects == nil {
			b.fileDefaultProjects = make(map[tspath.Path]tspath.Path)
		}
		b.fileDefaultProjects[path] = project.project.configFilePath
		return project
	}
	return nil
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
