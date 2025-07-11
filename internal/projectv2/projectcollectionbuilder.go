package projectv2

import (
	"context"
	"crypto/sha256"
	"fmt"
	"maps"

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
	logger              *logCollector

	ctx                                context.Context
	fs                                 *overlayFS
	base                               *ProjectCollection
	compilerOptionsForInferredProjects *core.CompilerOptions
	configFileRegistryBuilder          *configFileRegistryBuilder

	fileDefaultProjects map[tspath.Path]tspath.Path
	configuredProjects  *dirty.SyncMap[tspath.Path, *Project]
	inferredProject     *dirty.Box[*Project]
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
	logger *logCollector,
) *projectCollectionBuilder {
	if logger != nil {
		logger = logger.Fork("projectCollectionBuilder", "")
	}
	return &projectCollectionBuilder{
		ctx:                                ctx,
		fs:                                 fs,
		compilerOptionsForInferredProjects: compilerOptionsForInferredProjects,
		sessionOptions:                     sessionOptions,
		parseCache:                         parseCache,
		extendedConfigCache:                extendedConfigCache,
		logger:                             logger,
		base:                               oldProjectCollection,
		configFileRegistryBuilder:          newConfigFileRegistryBuilder(fs, oldConfigFileRegistry, extendedConfigCache, sessionOptions),
		configuredProjects:                 dirty.NewSyncMap(oldProjectCollection.configuredProjects, nil),
		inferredProject:                    dirty.NewBox(oldProjectCollection.inferredProject),
	}
}

func (b *projectCollectionBuilder) Finalize() (*ProjectCollection, *ConfigFileRegistry) {
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

	return newProjectCollection, b.configFileRegistryBuilder.finalize()
}

func (b *projectCollectionBuilder) forEachProject(fn func(entry dirty.Value[*Project]) bool) {
	var keepGoing bool
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

func (b *projectCollectionBuilder) DidCloseFile(uri lsproto.DocumentUri, hash [sha256.Size]byte) {
	fileName := uri.FileName()
	path := b.toPath(fileName)
	fh := b.fs.getFile(fileName)
	if fh != nil && fh.Hash() != hash {
		b.forEachProject(func(entry dirty.Value[*Project]) bool {
			b.markFileChanged(path)
			return true
		})
	}
	b.configFileRegistryBuilder.DidCloseFile(path)
	if fh == nil {
		// !!! handleDeletedFile
	}
}

func (b *projectCollectionBuilder) DidOpenFile(uri lsproto.DocumentUri) {
	if b.logger != nil {
		b.logger.Logf("DidOpenFile: %s", uri)
	}
	fileName := uri.FileName()
	path := b.toPath(fileName)
	var toRemoveProjects collections.Set[tspath.Path]
	b.ensureConfiguredProjectAndAncestorsForOpenFile(fileName, path)
	b.forEachProject(func(entry dirty.Value[*Project]) bool {
		toRemoveProjects.Add(entry.Value().configFilePath)
		b.updateProgram(entry)
		return true
	})
	if b.findDefaultProject(fileName, path) == nil {
		b.addFileToInferredProject(fileName, path)
	}

	for _, overlay := range b.fs.overlays {
		if toRemoveProjects.Len() == 0 {
			break
		}
		if p := b.findDefaultProject(overlay.FileName(), b.toPath(overlay.FileName())); p != nil {
			toRemoveProjects.Delete(p.Value().configFilePath)
		}
	}

	for projectPath := range toRemoveProjects.Keys() {
		b.deleteProject(projectPath)
	}
	b.configFileRegistryBuilder.Cleanup()
}

func (b *projectCollectionBuilder) DidChangeFiles(uris []lsproto.DocumentUri) {
	for _, uri := range uris {
		b.markFileChanged(uri.Path(b.fs.fs.UseCaseSensitiveFileNames()))
	}
}

func (b *projectCollectionBuilder) DidRequestFile(uri lsproto.DocumentUri) {
	// See if we can find a default project for this file without doing
	// any additional loading.
	fileName := uri.FileName()
	path := b.toPath(fileName)
	if result := b.findDefaultProject(fileName, path); result != nil {
		b.updateProgram(result)
		return
	}

	// Make sure all projects we know about are up to date...
	var hasChanges bool
	b.configuredProjects.Range(func(entry *dirty.SyncMapEntry[tspath.Path, *Project]) bool {
		hasChanges = b.updateProgram(entry) || hasChanges
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
			b.updateInferredProject(inferredProjectFiles)
		}
	}

	// ...and then try to find the default configured project for this file again.
	if b.findDefaultProject(fileName, path) == nil {
		panic(fmt.Sprintf("no project found for file %s", fileName))
	}
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
	if b.isOpenFile(path) {
		return b.findOrCreateDefaultConfiguredProjectForOpenScriptInfo(fileName, path, projectLoadKindFind)
	}
	return nil
}

func (b *projectCollectionBuilder) ensureConfiguredProjectAndAncestorsForOpenFile(fileName string, path tspath.Path) {
	result := b.findOrCreateDefaultConfiguredProjectForOpenScriptInfo(fileName, path, projectLoadKindCreate)
	if result != nil && result.Value() != nil {
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
}

type searchNode struct {
	configFileName string
	loadKind       projectLoadKind
}

func (b *projectCollectionBuilder) findOrCreateDefaultConfiguredProjectWorker(
	fileName string,
	path tspath.Path,
	configFileName string,
	loadKind projectLoadKind,
	visited *collections.SyncSet[searchNode],
	fallback *searchNode,
) *dirty.SyncMapEntry[tspath.Path, *Project] {
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
				b.updateProgram(project)
			}

			if project.Value().containsFile(path) {
				return true, !project.Value().IsSourceFromProjectReference(path)
			}

			return false, false
		},
		visited,
	)

	if search.Stopped {
		project, _ := b.configuredProjects.Load(b.toPath(search.Path[0].configFileName))
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
			project, _ := b.configuredProjects.Load(b.toPath(fallback.configFileName))
			return project
		}
	}
	if ancestorConfigName := b.configFileRegistryBuilder.getAncestorConfigFileName(fileName, path, configFileName, loadKind); ancestorConfigName != "" {
		return b.findOrCreateDefaultConfiguredProjectWorker(fileName, path, ancestorConfigName, loadKind, visited, fallback)
	}
	if fallback != nil {
		project, _ := b.configuredProjects.Load(b.toPath(fallback.configFileName))
		return project
	}
	return nil
}

func (b *projectCollectionBuilder) findOrCreateDefaultConfiguredProjectForOpenScriptInfo(
	fileName string,
	path tspath.Path,
	loadKind projectLoadKind,
) *dirty.SyncMapEntry[tspath.Path, *Project] {
	if key, ok := b.fileDefaultProjects[path]; ok {
		if key == inferredProjectName {
			// The file belongs to the inferred project
			return nil
		}
		entry, _ := b.configuredProjects.Load(key)
		return entry
	}
	if configFileName := b.configFileRegistryBuilder.getConfigFileNameForFile(fileName, path, loadKind); configFileName != "" {
		project := b.findOrCreateDefaultConfiguredProjectWorker(fileName, path, configFileName, loadKind, nil, nil)
		if b.fileDefaultProjects == nil {
			b.fileDefaultProjects = make(map[tspath.Path]tspath.Path)
		}
		b.fileDefaultProjects[path] = project.Value().configFilePath
		return project
	}
	return nil
}

func (b *projectCollectionBuilder) findOrCreateProject(
	configFileName string,
	configFilePath tspath.Path,
	loadKind projectLoadKind,
) *dirty.SyncMapEntry[tspath.Path, *Project] {
	if loadKind == projectLoadKindFind {
		entry, _ := b.configuredProjects.Load(configFilePath)
		return entry
	}
	entry, _ := b.configuredProjects.LoadOrStore(configFilePath, NewConfiguredProject(configFileName, configFilePath, b))
	return entry
}

func (b *projectCollectionBuilder) toPath(fileName string) tspath.Path {
	return tspath.ToPath(fileName, b.sessionOptions.CurrentDirectory, b.fs.fs.UseCaseSensitiveFileNames())
}

func (b *projectCollectionBuilder) isOpenFile(path tspath.Path) bool {
	_, ok := b.fs.overlays[path]
	return ok
}

func (b *projectCollectionBuilder) updateInferredProject(rootFileNames []string) bool {
	if b.inferredProject.Value() == nil && len(rootFileNames) > 0 {
		b.inferredProject.Set(NewInferredProject(b.sessionOptions.CurrentDirectory, b.compilerOptionsForInferredProjects, rootFileNames, b))
	} else if b.inferredProject.Value() != nil && len(rootFileNames) == 0 {
		b.inferredProject.Delete()
		return true
	} else {
		newCommandLine := tsoptions.NewParsedCommandLine(b.compilerOptionsForInferredProjects, rootFileNames, tspath.ComparePathsOptions{
			UseCaseSensitiveFileNames: b.fs.fs.UseCaseSensitiveFileNames(),
			CurrentDirectory:          b.sessionOptions.CurrentDirectory,
		})
		changed := b.inferredProject.ChangeIf(
			func(p *Project) bool {
				return !maps.Equal(p.CommandLine.FileNamesByPath(), newCommandLine.FileNamesByPath())
			},
			func(p *Project) {
				p.CommandLine = newCommandLine
				p.dirty = true
				p.dirtyFilePath = ""
			},
		)
		if !changed {
			return false
		}
	}
	return b.updateProgram(b.inferredProject)
}

func (b *projectCollectionBuilder) addFileToInferredProject(fileName string, path tspath.Path) bool {
	if b.inferredProject.Value() == nil {
		return b.updateInferredProject([]string{fileName})
	}
	return b.updateInferredProject(append(b.inferredProject.Value().CommandLine.FileNames(), fileName))
}

func (b *projectCollectionBuilder) updateProgram(entry dirty.Value[*Project]) bool {
	var newCommandLine *tsoptions.ParsedCommandLine
	return entry.ChangeIf(
		func(project *Project) bool {
			if project.Kind == KindConfigured {
				commandLine := b.configFileRegistryBuilder.acquireConfigForProject(project.configFileName, project.configFilePath, project)
				if project.CommandLine != commandLine {
					newCommandLine = commandLine
					return true
				}
			}
			return project.dirty
		},
		func(project *Project) {
			if newCommandLine != nil {
				project.CommandLine = newCommandLine
			}
			project.host = newCompilerHost(project.currentDirectory, project, b)
			newProgram, checkerPool := project.CreateProgram()
			project.Program = newProgram
			project.checkerPool = checkerPool
			// !!! unthread context
			project.LanguageService = ls.NewLanguageService(b.ctx, project)
			project.dirty = false
			project.dirtyFilePath = ""
		},
	)
}

func (b *projectCollectionBuilder) markFileChanged(path tspath.Path) {
	b.forEachProject(func(entry dirty.Value[*Project]) bool {
		entry.ChangeIf(
			func(p *Project) bool { return p.containsFile(path) },
			func(p *Project) {
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

func (b *projectCollectionBuilder) deleteProject(path tspath.Path) {
	if project, ok := b.configuredProjects.Load(path); ok {
		if program := project.Value().Program; program != nil {
			program.ForEachResolvedProjectReference(func(path tspath.Path, config *tsoptions.ParsedCommandLine) {
				b.configFileRegistryBuilder.releaseConfigForProject(path, path)
			})
		}
		if project.Value().Kind == KindConfigured {
			b.configFileRegistryBuilder.releaseConfigForProject(path, path)
		}
		project.Delete()
	}
}
