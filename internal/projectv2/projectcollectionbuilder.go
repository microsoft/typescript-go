package projectv2

import (
	"context"
	"crypto/sha256"
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/dirty"
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
	fs                                 *snapshotFSBuilder
	base                               *ProjectCollection
	compilerOptionsForInferredProjects *core.CompilerOptions
	configFileRegistryBuilder          *configFileRegistryBuilder

	projectsAffectedByConfigChanges map[tspath.Path]struct{}
	filesAffectedByConfigChanges    map[tspath.Path]struct{}
	fileDefaultProjects             map[tspath.Path]tspath.Path
	configuredProjects              *dirty.SyncMap[tspath.Path, *Project]
	inferredProject                 *dirty.Box[*Project]
}

func newProjectCollectionBuilder(
	ctx context.Context,
	fs *snapshotFSBuilder,
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
		projectsAffectedByConfigChanges:    make(map[tspath.Path]struct{}),
		filesAffectedByConfigChanges:       make(map[tspath.Path]struct{}),
		configFileRegistryBuilder:          newConfigFileRegistryBuilder(fs, oldConfigFileRegistry, extendedConfigCache, sessionOptions, nil),
		configuredProjects:                 dirty.NewSyncMap(oldProjectCollection.configuredProjects, nil),
		inferredProject:                    dirty.NewBox(oldProjectCollection.inferredProject),
	}
}

func (b *projectCollectionBuilder) Finalize(logger *logCollector) (*ProjectCollection, *ConfigFileRegistry) {
	b.markProjectsAffectedByConfigChanges(logger)
	var changed bool
	newProjectCollection := b.base
	ensureCloned := func() {
		if !changed {
			newProjectCollection = newProjectCollection.clone()
			changed = true
		}
	}

	if configuredProjects, configuredProjectsChanged := b.configuredProjects.Finalize(); configuredProjectsChanged {
		ensureCloned()
		newProjectCollection.configuredProjects = configuredProjects
	}

	if !changed && !maps.Equal(b.fileDefaultProjects, b.base.fileDefaultProjects) {
		ensureCloned()
		newProjectCollection.fileDefaultProjects = b.fileDefaultProjects
	}

	if newInferredProject, inferredProjectChanged := b.inferredProject.Finalize(); inferredProjectChanged {
		ensureCloned()
		newProjectCollection.inferredProject = newInferredProject
	}

	configFileRegistry := b.configFileRegistryBuilder.Finalize()
	newProjectCollection.configFileRegistry = configFileRegistry
	return newProjectCollection, configFileRegistry
}

func (b *projectCollectionBuilder) forEachProject(fn func(entry dirty.Value[*Project]) bool) {
	keepGoing := true
	b.configuredProjects.Range(func(entry *dirty.SyncMapEntry[tspath.Path, *Project]) bool {
		keepGoing = fn(entry)
		return keepGoing
	})
	if !keepGoing {
		return
	}
	if b.inferredProject.Value() != nil {
		fn(b.inferredProject)
	}
}

func (b *projectCollectionBuilder) DidCloseFile(uri lsproto.DocumentUri, hash [sha256.Size]byte, logger *logCollector) {
	fileName := uri.FileName()
	path := b.toPath(fileName)
	fh := b.fs.GetFileByPath(fileName, path)
	if fh == nil || fh.Hash() != hash {
		b.forEachProject(func(entry dirty.Value[*Project]) bool {
			b.markFileChanged(path, logger)
			return true
		})
	}
	if b.inferredProject.Value() != nil {
		rootFilesMap := b.inferredProject.Value().CommandLine.FileNamesByPath()
		if fileName, ok := rootFilesMap[path]; ok {
			rootFiles := b.inferredProject.Value().CommandLine.FileNames()
			index := slices.Index(rootFiles, fileName)
			newRootFiles := slices.Delete(rootFiles, index, index+1)
			b.updateInferredProject(newRootFiles, logger)
		}
	}
	b.configFileRegistryBuilder.DidCloseFile(path)
}

func (b *projectCollectionBuilder) DidOpenFile(uri lsproto.DocumentUri, logger *logCollector) {
	fileName := uri.FileName()
	path := b.toPath(fileName)
	var toRemoveProjects collections.Set[tspath.Path]
	openFileResult := b.ensureConfiguredProjectAndAncestorsForOpenFile(fileName, path, logger)
	b.configuredProjects.Range(func(entry *dirty.SyncMapEntry[tspath.Path, *Project]) bool {
		toRemoveProjects.Add(entry.Value().configFilePath)
		b.updateProgram(entry, logger)
		return true
	})

	var inferredProjectFiles []string
	for _, overlay := range b.fs.overlays {
		if p := b.findDefaultConfiguredProject(overlay.FileName(), b.toPath(overlay.FileName())); p != nil {
			toRemoveProjects.Delete(p.Value().configFilePath)
		} else {
			inferredProjectFiles = append(inferredProjectFiles, overlay.FileName())
		}
	}

	for projectPath := range toRemoveProjects.Keys() {
		if !openFileResult.retain.Has(projectPath) {
			if p, ok := b.configuredProjects.Load(projectPath); ok {
				b.deleteConfiguredProject(p, logger)
			}
		}
	}
	b.updateInferredProject(inferredProjectFiles, logger)
	b.configFileRegistryBuilder.Cleanup()
}

func (b *projectCollectionBuilder) DidDeleteFiles(uris []lsproto.DocumentUri, logger *logCollector) {
	for _, uri := range uris {
		path := uri.Path(b.fs.fs.UseCaseSensitiveFileNames())
		result := b.configFileRegistryBuilder.DidDeleteFile(path)
		maps.Copy(b.projectsAffectedByConfigChanges, result.affectedProjects)
		maps.Copy(b.filesAffectedByConfigChanges, result.affectedFiles)
		if result.IsEmpty() {
			b.forEachProject(func(entry dirty.Value[*Project]) bool {
				entry.ChangeIf(
					func(p *Project) bool { return (!p.dirty || p.dirtyFilePath != "") && p.containsFile(path) },
					func(p *Project) {
						p.dirty = true
						p.dirtyFilePath = ""
						logger.Logf("Marked project %s as dirty", p.configFileName)
					},
				)
				return true
			})
		} else if logger != nil {
			logChangeFileResult(result, logger)
		}
	}
}

// DidCreateFiles is only called when file watching is enabled.
func (b *projectCollectionBuilder) DidCreateFiles(uris []lsproto.DocumentUri, logger *logCollector) {
	// !!! some way to stop iterating when everything that can be marked has been marked?
	for _, uri := range uris {
		fileName := uri.FileName()
		path := uri.Path(b.fs.fs.UseCaseSensitiveFileNames())
		result := b.configFileRegistryBuilder.DidCreateFile(fileName, path)
		maps.Copy(b.projectsAffectedByConfigChanges, result.affectedProjects)
		maps.Copy(b.filesAffectedByConfigChanges, result.affectedFiles)
		if logger != nil {
			logChangeFileResult(result, logger)
		}
		b.forEachProject(func(entry dirty.Value[*Project]) bool {
			entry.ChangeIf(
				func(p *Project) bool {
					if p.dirty && p.dirtyFilePath == "" {
						return false
					}
					if _, ok := p.failedLookupsWatch.input[path]; ok {
						return true
					}
					if _, ok := p.affectingLocationsWatch.input[path]; ok {
						return true
					}
					return false
				},
				func(p *Project) {
					p.dirty = true
					p.dirtyFilePath = ""
					logger.Logf("Marked project %s as dirty", p.configFileName)
				},
			)
			return true
		})
	}
}

func (b *projectCollectionBuilder) DidChangeFiles(uris []lsproto.DocumentUri, logger *logCollector) {
	for _, uri := range uris {
		path := uri.Path(b.fs.fs.UseCaseSensitiveFileNames())
		result := b.configFileRegistryBuilder.DidChangeFile(path)
		maps.Copy(b.projectsAffectedByConfigChanges, result.affectedProjects)
		maps.Copy(b.filesAffectedByConfigChanges, result.affectedFiles)
		if result.IsEmpty() {
			b.markFileChanged(path, logger)
		} else if logger != nil {
			logChangeFileResult(result, logger)
		}
	}
}

func logChangeFileResult(result changeFileResult, logger *logCollector) {
	if len(result.affectedProjects) > 0 {
		logger.Logf("Config file change affected projects: %v", slices.Collect(maps.Keys(result.affectedProjects)))
	}
	if len(result.affectedFiles) > 0 {
		logger.Logf("Config file change affected config file lookups for %d files", len(result.affectedFiles))
	}
}

func (b *projectCollectionBuilder) DidRequestFile(uri lsproto.DocumentUri, logger *logCollector) {
	startTime := time.Now()
	fileName := uri.FileName()

	hasChanges := b.markProjectsAffectedByConfigChanges(logger)

	// See if we can find a default project without updating a bunch of stuff.
	path := b.toPath(fileName)
	if result := b.findDefaultProject(fileName, path); result != nil {
		hasChanges = b.updateProgram(result, logger) || hasChanges
		if result.Value() != nil {
			return
		}
	}

	// Make sure all projects we know about are up to date...
	b.configuredProjects.Range(func(entry *dirty.SyncMapEntry[tspath.Path, *Project]) bool {
		hasChanges = b.updateProgram(entry, logger) || hasChanges
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
			b.updateInferredProject(inferredProjectFiles, logger)
		}
	}

	// ...and then try to find the default configured project for this file again.
	if b.findDefaultProject(fileName, path) == nil {
		panic(fmt.Sprintf("no project found for file %s", fileName))
	}

	if logger != nil {
		elapsed := time.Since(startTime)
		logger.Log(fmt.Sprintf("Completed file request for %s in %v", fileName, elapsed))
	}
}

func (b *projectCollectionBuilder) DidUpdateATAState(ataChanges map[tspath.Path]*ATAStateChange, logger *logCollector) {
	updateProject := func(project dirty.Value[*Project], ataChange *ATAStateChange) {
		project.ChangeIf(
			func(p *Project) bool {
				if p == nil {
					return false
				}
				return ataChange.TypingsInfo.Equals(p.ComputeTypingsInfo())
			},
			func(p *Project) {
				p.installedTypingsInfo = ataChange.TypingsInfo
				if !slices.Equal(p.typingsFiles, ataChange.TypingsFiles) {
					p.typingsFiles = ataChange.TypingsFiles
					p.dirty = true
					p.dirtyFilePath = ""
				}
			},
		)
	}

	for projectPath, ataChange := range ataChanges {
		ataChange.Logs.WriteLogs(logger.Fork("Typings Installer Logs for " + string(projectPath)))
		if projectPath == inferredProjectName {
			updateProject(b.inferredProject, ataChange)
		} else if project, ok := b.configuredProjects.Load(projectPath); ok {
			updateProject(project, ataChange)
		}

		if logger != nil {
			logger.Log(fmt.Sprintf("Updated ATA state for project %s", projectPath))
		}
	}
}

func (b *projectCollectionBuilder) markProjectsAffectedByConfigChanges(logger *logCollector) bool {
	for projectPath := range b.projectsAffectedByConfigChanges {
		project, ok := b.configuredProjects.Load(projectPath)
		if !ok {
			panic(fmt.Sprintf("project %s affected by config change not found", projectPath))
		}
		project.ChangeIf(
			func(p *Project) bool { return !p.dirty || p.dirtyFilePath != "" },
			func(p *Project) {
				p.dirty = true
				p.dirtyFilePath = ""
			},
		)
	}

	// Recompute default projects for open files that now have different config file presence.
	var hasChanges bool
	for path := range b.filesAffectedByConfigChanges {
		fileName := b.fs.overlays[path].FileName()
		_ = b.ensureConfiguredProjectAndAncestorsForOpenFile(fileName, path, logger)
		hasChanges = true
	}

	b.projectsAffectedByConfigChanges = nil
	b.filesAffectedByConfigChanges = nil
	return hasChanges
}

func (b *projectCollectionBuilder) findDefaultProject(fileName string, path tspath.Path) dirty.Value[*Project] {
	if configuredProject := b.findDefaultConfiguredProject(fileName, path); configuredProject != nil {
		return configuredProject
	}
	if key, ok := b.fileDefaultProjects[path]; ok && key == inferredProjectName {
		return b.inferredProject
	}
	if inferredProject := b.inferredProject.Value(); inferredProject != nil && inferredProject.containsFile(path) {
		if b.fileDefaultProjects == nil {
			b.fileDefaultProjects = make(map[tspath.Path]tspath.Path)
		}
		b.fileDefaultProjects[path] = inferredProjectName
		return b.inferredProject
	}
	return nil
}

func (b *projectCollectionBuilder) findDefaultConfiguredProject(fileName string, path tspath.Path) *dirty.SyncMapEntry[tspath.Path, *Project] {
	// !!! look in fileDefaultProjects first?
	// Sort configured projects so we can use a deterministic "first" as a last resort.
	var configuredProjectPaths []tspath.Path
	configuredProjects := make(map[tspath.Path]*dirty.SyncMapEntry[tspath.Path, *Project])
	b.configuredProjects.Range(func(entry *dirty.SyncMapEntry[tspath.Path, *Project]) bool {
		configuredProjectPaths = append(configuredProjectPaths, entry.Key())
		configuredProjects[entry.Key()] = entry
		return true
	})
	slices.Sort(configuredProjectPaths)

	project, multipleCandidates := findDefaultConfiguredProjectFromProgramInclusion(fileName, path, configuredProjectPaths, func(path tspath.Path) *Project {
		return configuredProjects[path].Value()
	})

	if multipleCandidates {
		if p := b.findOrCreateDefaultConfiguredProjectForOpenScriptInfo(fileName, path, projectLoadKindFind, nil).project; p != nil {
			return p
		}
	}

	return configuredProjects[project]
}

func (b *projectCollectionBuilder) ensureConfiguredProjectAndAncestorsForOpenFile(fileName string, path tspath.Path, logger *logCollector) searchResult {
	result := b.findOrCreateDefaultConfiguredProjectForOpenScriptInfo(fileName, path, projectLoadKindCreate, logger)
	if result.project != nil {
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
	logger         *logCollector
}

type searchResult struct {
	project *dirty.SyncMapEntry[tspath.Path, *Project]
	retain  collections.Set[tspath.Path]
}

func (b *projectCollectionBuilder) findOrCreateDefaultConfiguredProjectWorker(
	fileName string,
	path tspath.Path,
	configFileName string,
	loadKind projectLoadKind,
	visited *collections.SyncSet[searchNode],
	fallback *searchResult,
	logger *logCollector,
) searchResult {
	var configs collections.SyncMap[tspath.Path, *tsoptions.ParsedCommandLine]
	if visited == nil {
		visited = &collections.SyncSet[searchNode]{}
	}

	search := core.BreadthFirstSearchParallelEx(
		searchNode{configFileName: configFileName, loadKind: loadKind, logger: logger},
		func(node searchNode) []searchNode {
			if config, ok := configs.Load(b.toPath(node.configFileName)); ok && len(config.ProjectReferences()) > 0 {
				referenceLoadKind := node.loadKind
				if config.CompilerOptions().DisableReferencedProjectLoad.IsTrue() {
					referenceLoadKind = projectLoadKindFind
				}

				var logger *logCollector
				references := config.ResolvedProjectReferencePaths()
				if len(references) > 0 && node.logger != nil {
					logger = node.logger.Fork(fmt.Sprintf("Searching %d project references of %s", len(references), node.configFileName))
				}
				return core.Map(references, func(configFileName string) searchNode {
					return searchNode{configFileName: configFileName, loadKind: referenceLoadKind, logger: logger.Fork("Searching project reference " + configFileName)}
				})
			}
			return nil
		},
		func(node searchNode) (isResult bool, stop bool) {
			configFilePath := b.toPath(node.configFileName)
			config := b.configFileRegistryBuilder.findOrAcquireConfigForOpenFile(node.configFileName, configFilePath, path, node.loadKind)
			if config == nil {
				return false, false
			}
			configs.Store(configFilePath, config)
			if len(config.FileNames()) == 0 {
				// Likely a solution tsconfig.json - the search will fan out to its references.
				node.logger.Log("Project does not contain file (no root files)")
				return false, false
			}

			if config.CompilerOptions().Composite == core.TSTrue {
				// For composite projects, we can get an early negative result.
				// !!! what about declaration files in node_modules? wouldn't it be better to
				//     check project inclusion if the project is already loaded?
				if !config.MatchesFileName(fileName) {
					node.logger.Log("Project does not contain file (by composite config inclusion)")
					return false, false
				}
			}

			project := b.findOrCreateProject(node.configFileName, configFilePath, node.loadKind, node.logger)
			if node.loadKind == projectLoadKindCreate {
				// Ensure project is up to date before checking for file inclusion
				b.updateProgram(project, node.logger)
			}

			if project.Value().containsFile(path) {
				isDirectInclusion := !project.Value().IsSourceFromProjectReference(path)
				if node.logger != nil {
					node.logger.Logf("Project contains file %s", core.IfElse(isDirectInclusion, "directly", "as a source of a referenced project"))
				}
				return true, isDirectInclusion
			}

			node.logger.Log("Project does not contain file")
			return false, false
		},
		core.BreadthFirstSearchOptions[searchNode]{
			Visited: visited,
			PreprocessLevel: func(level *core.BreadthFirstSearchLevel[searchNode]) {
				level.Range(func(node searchNode) bool {
					if node.loadKind == projectLoadKindFind && level.Has(searchNode{configFileName: node.configFileName, loadKind: projectLoadKindCreate, logger: node.logger}) {
						// Remove find requests when a create request for the same project is already present.
						level.Delete(node)
					}
					return true
				})
			},
		},
	)

	var retain collections.Set[tspath.Path]
	var project *dirty.SyncMapEntry[tspath.Path, *Project]
	if len(search.Path) > 0 {
		project, _ = b.configuredProjects.Load(b.toPath(search.Path[0].configFileName))
		// If we found a project, we retain each project along the BFS path.
		// We don't want to retain everything we visited since BFS can terminate
		// early, and we don't want to retain nondeterministically.
		for _, node := range search.Path {
			retain.Add(b.toPath(node.configFileName))
		}
	}

	if search.Stopped {
		// Found a project that directly contains the file.
		return searchResult{
			project: project,
			retain:  retain,
		}
	}

	if project != nil {
		// If we found a project that contains the file, but it is a source from
		// a project reference, record it as a fallback.
		fallback = &searchResult{
			project: project,
			retain:  retain,
		}
	}

	// Look for tsconfig.json files higher up the directory tree and do the same. This handles
	// the common case where a higher-level "solution" tsconfig.json contains all projects in a
	// workspace.
	if config, ok := configs.Load(b.toPath(configFileName)); ok && config.CompilerOptions().DisableSolutionSearching.IsTrue() {
		if fallback != nil {
			return *fallback
		}
	}
	if ancestorConfigName := b.configFileRegistryBuilder.getAncestorConfigFileName(fileName, path, configFileName, loadKind); ancestorConfigName != "" {
		return b.findOrCreateDefaultConfiguredProjectWorker(
			fileName,
			path,
			ancestorConfigName,
			loadKind,
			visited,
			fallback,
			logger.Fork(fmt.Sprintf("Searching ancestor config file at %s", ancestorConfigName)),
		)
	}
	if fallback != nil {
		return *fallback
	}
	// If we didn't find anything, we can retain everything we visited,
	// since the whole graph must have been traversed (i.e., the set of
	// retained projects is guaranteed to be deterministic).
	visited.Range(func(node searchNode) bool {
		retain.Add(b.toPath(node.configFileName))
		return true
	})
	return searchResult{retain: retain}
}

func (b *projectCollectionBuilder) findOrCreateDefaultConfiguredProjectForOpenScriptInfo(
	fileName string,
	path tspath.Path,
	loadKind projectLoadKind,
	logger *logCollector,
) searchResult {
	if key, ok := b.fileDefaultProjects[path]; ok {
		if key == inferredProjectName {
			// The file belongs to the inferred project
			return searchResult{}
		}
		entry, _ := b.configuredProjects.Load(key)
		return searchResult{project: entry}
	}
	if configFileName := b.configFileRegistryBuilder.getConfigFileNameForFile(fileName, path, loadKind); configFileName != "" {
		startTime := time.Now()
		result := b.findOrCreateDefaultConfiguredProjectWorker(
			fileName,
			path,
			configFileName,
			loadKind,
			nil,
			nil,
			logger.Fork(fmt.Sprintf("Searching for default configured project for %s", fileName)),
		)
		if result.project != nil {
			if b.fileDefaultProjects == nil {
				b.fileDefaultProjects = make(map[tspath.Path]tspath.Path)
			}
			b.fileDefaultProjects[path] = result.project.Value().configFilePath
		}
		if logger != nil {
			elapsed := time.Since(startTime)
			if result.project != nil {
				logger.Log(fmt.Sprintf("Found default configured project for %s: %s (in %v)", fileName, result.project.Value().configFileName, elapsed))
			} else {
				logger.Log(fmt.Sprintf("No default configured project found for %s (searched in %v)", fileName, elapsed))
			}
		}
		return result
	}
	return searchResult{}
}

func (b *projectCollectionBuilder) findOrCreateProject(
	configFileName string,
	configFilePath tspath.Path,
	loadKind projectLoadKind,
	logger *logCollector,
) *dirty.SyncMapEntry[tspath.Path, *Project] {
	if loadKind == projectLoadKindFind {
		entry, _ := b.configuredProjects.Load(configFilePath)
		return entry
	}
	entry, _ := b.configuredProjects.LoadOrStore(configFilePath, NewConfiguredProject(configFileName, configFilePath, b, logger))
	return entry
}

func (b *projectCollectionBuilder) toPath(fileName string) tspath.Path {
	return tspath.ToPath(fileName, b.sessionOptions.CurrentDirectory, b.fs.fs.UseCaseSensitiveFileNames())
}

func (b *projectCollectionBuilder) updateInferredProject(rootFileNames []string, logger *logCollector) bool {
	if len(rootFileNames) == 0 {
		if b.inferredProject.Value() != nil {
			if logger != nil {
				logger.Log("Deleting inferred project")
			}
			b.inferredProject.Delete()
			return true
		}
		return false
	}

	if b.inferredProject.Value() == nil {
		b.inferredProject.Set(NewInferredProject(b.sessionOptions.CurrentDirectory, b.compilerOptionsForInferredProjects, rootFileNames, b, logger))
	} else {
		newCompilerOptions := b.inferredProject.Value().CommandLine.CompilerOptions()
		if b.compilerOptionsForInferredProjects != nil {
			newCompilerOptions = b.compilerOptionsForInferredProjects
		}
		newCommandLine := tsoptions.NewParsedCommandLine(newCompilerOptions, rootFileNames, tspath.ComparePathsOptions{
			UseCaseSensitiveFileNames: b.fs.fs.UseCaseSensitiveFileNames(),
			CurrentDirectory:          b.sessionOptions.CurrentDirectory,
		})
		changed := b.inferredProject.ChangeIf(
			func(p *Project) bool {
				return !maps.Equal(p.CommandLine.FileNamesByPath(), newCommandLine.FileNamesByPath())
			},
			func(p *Project) {
				if logger != nil {
					logger.Log(fmt.Sprintf("Updating inferred project config with %d root files", len(rootFileNames)))
				}
				p.CommandLine = newCommandLine
				p.dirty = true
				p.dirtyFilePath = ""
			},
		)
		if !changed {
			return false
		}
	}
	return b.updateProgram(b.inferredProject, logger)
}

// updateProgram updates the program for the given project entry if necessary. It returns
// a boolean indicating whether the update could have caused any structure-affecting changes.
func (b *projectCollectionBuilder) updateProgram(entry dirty.Value[*Project], logger *logCollector) bool {
	var updateProgram bool
	var filesChanged bool
	configFileName := entry.Value().configFileName
	startTime := time.Now()
	entry.Locked(func(entry dirty.Value[*Project]) {
		if entry.Value().Kind == KindConfigured {
			commandLine := b.configFileRegistryBuilder.acquireConfigForProject(entry.Value().configFileName, entry.Value().configFilePath, entry.Value())
			if entry.Value().CommandLine != commandLine {
				updateProgram = true
				if commandLine == nil {
					b.deleteConfiguredProject(entry, logger)
					filesChanged = true
					return
				}
				entry.Change(func(p *Project) { p.CommandLine = commandLine })
			}
		}
		if !updateProgram {
			updateProgram = entry.Value().dirty
		}
		if updateProgram {
			entry.Change(func(project *Project) {
				project.host = newCompilerHost(project.currentDirectory, project, b)
				result := project.CreateProgram()
				project.Program = result.Program
				project.checkerPool = result.CheckerPool
				project.ProgramUpdateKind = result.UpdateKind
				if result.UpdateKind == ProgramUpdateKindNewFiles {
					filesChanged = true
					if b.sessionOptions.WatchEnabled {
						failedLookupsWatch, affectingLocationsWatch := project.CloneWatchers()
						project.failedLookupsWatch = failedLookupsWatch
						project.affectingLocationsWatch = affectingLocationsWatch
					}
				}
				// !!! unthread context
				project.LanguageService = ls.NewLanguageService(b.ctx, project)
				project.dirty = false
				project.dirtyFilePath = ""
			})
			delete(b.projectsAffectedByConfigChanges, entry.Value().configFilePath)
		}
	})
	if updateProgram && logger != nil {
		elapsed := time.Since(startTime)
		logger.Log(fmt.Sprintf("Program update for %s completed in %v", configFileName, elapsed))
	}
	return filesChanged
}

func (b *projectCollectionBuilder) markFileChanged(path tspath.Path, logger *logCollector) {
	b.forEachProject(func(entry dirty.Value[*Project]) bool {
		entry.ChangeIf(
			func(p *Project) bool { return (!p.dirty || p.dirtyFilePath != path) && p.containsFile(path) },
			func(p *Project) {
				if logger != nil {
					logger.Log(fmt.Sprintf("Marking project %s as dirty due to file change %s", p.configFileName, path))
				}
				if !p.dirty {
					p.dirty = true
					p.dirtyFilePath = path
				} else if p.dirtyFilePath != path {
					p.dirtyFilePath = ""
				}
			})
		return true
	})
}

func (b *projectCollectionBuilder) deleteConfiguredProject(project dirty.Value[*Project], logger *logCollector) {
	projectPath := project.Value().configFilePath
	if logger != nil {
		logger.Log(fmt.Sprintf("Deleting configured project: %s", project.Value().configFileName))
	}
	if program := project.Value().Program; program != nil {
		program.ForEachResolvedProjectReference(func(referencePath tspath.Path, config *tsoptions.ParsedCommandLine) {
			b.configFileRegistryBuilder.releaseConfigForProject(referencePath, projectPath)
		})
	}
	b.configFileRegistryBuilder.releaseConfigForProject(projectPath, projectPath)
	project.Delete()
}
