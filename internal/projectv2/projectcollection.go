package projectv2

import (
	"context"
	"fmt"
	"maps"
	"strings"
	"sync"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
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

type projectCollection struct {
	configuredProjects map[tspath.Path]*Project
	inferredProject    *Project
}

func (c *projectCollection) clone() *projectCollection {
	return &projectCollection{
		configuredProjects: maps.Clone(c.configuredProjects),
		inferredProject:    c.inferredProject,
	}
}

type projectCollectionBuilder struct {
	ctx                       context.Context
	snapshot                  *Snapshot
	configFileRegistryBuilder *configFileRegistryBuilder
	base                      *projectCollection
	changes                   snapshotChange
	dirty                     collections.SyncMap[tspath.Path, *Project]
}

type projectCollectionBuilderEntry struct {
	b       *projectCollectionBuilder
	project *Project
	dirty   bool
}

func newProjectCollectionBuilder(
	ctx context.Context,
	newSnapshot *Snapshot,
	oldProjectCollection *projectCollection,
	oldConfigFileRegistry *configFileRegistry,
	changes snapshotChange,
) *projectCollectionBuilder {
	return &projectCollectionBuilder{
		ctx:                       ctx,
		snapshot:                  newSnapshot,
		base:                      oldProjectCollection,
		configFileRegistryBuilder: newConfigFileRegistryBuilder(newSnapshot, oldConfigFileRegistry),
		changes:                   changes,
	}
}

func (b *projectCollectionBuilder) finalize() (*projectCollection, *configFileRegistry) {
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

func (b *projectCollectionBuilder) updateProject(entry *projectCollectionBuilderEntry) {
	if entry.dirty {
		// !!! right now, the only kind of project update is program loading,
		// so we can just assume that if the project is in the dirty map,
		// it's already been updated. This assumption probably won't hold
		// as this logic gets more fleshed out.
		if entry.project.Program == nil {
			entry.project, _ = entry.project.Clone(b.ctx, b.changes, b.snapshot)
		}
		return
	}
	if project, result := entry.project.Clone(b.ctx, b.changes, b.snapshot); result.changed {
		_, loaded := b.dirty.LoadOrStore(project.configFilePath, project)
		if loaded {
			// I don't think we get into a state where multiple goroutines try to update
			// the same project at the same time; ensure this is the case
			panic("unexpected concurrent project update")
		}
	}
}

func (b *projectCollectionBuilder) computeConfigFileName(fileName string, skipSearchInDirectoryOfFile bool) string {
	searchPath := tspath.GetDirectoryPath(fileName)
	result, _ := tspath.ForEachAncestorDirectory(searchPath, func(directory string) (result string, stop bool) {
		tsconfigPath := tspath.CombinePaths(directory, "tsconfig.json")
		if !skipSearchInDirectoryOfFile && b.snapshot.compilerFS.FileExists(tsconfigPath) {
			return tsconfigPath, true
		}
		jsconfigPath := tspath.CombinePaths(directory, "jsconfig.json")
		if !skipSearchInDirectoryOfFile && b.snapshot.compilerFS.FileExists(jsconfigPath) {
			return jsconfigPath, true
		}
		if strings.HasSuffix(directory, "/node_modules") {
			return "", true
		}
		skipSearchInDirectoryOfFile = false
		return "", false
	})
	b.snapshot.Logf("computeConfigFileName:: File: %s:: Result: %s", fileName, result)
	return result
}

func (b *projectCollectionBuilder) getConfigFileNameForFile(fileName string, path tspath.Path, loadKind projectLoadKind) string {
	if project.IsDynamicFileName(fileName) {
		return ""
	}

	// configName, ok := f.configFileForOpenFiles[path]
	// if ok {
	// 	return configName
	// }

	if loadKind == projectLoadKindFind {
		return ""
	}

	configName := b.computeConfigFileName(fileName, false)

	// if f.IsOpenFile(ls.FileNameToDocumentURI(fileName)) {
	// 	f.configFileForOpenFiles[path] = configName
	// }
	return configName
}

func (b *projectCollectionBuilder) getAncestorConfigFileName(fileName string, path tspath.Path, configFileName string, loadKind projectLoadKind) string {
	if project.IsDynamicFileName(fileName) {
		return ""
	}

	// if ancestorConfigMap, ok := f.configFilesAncestorForOpenFiles[path]; ok {
	// 	if ancestorConfigName, found := ancestorConfigMap[configFileName]; found {
	// 		return ancestorConfigName
	// 	}
	// }

	if loadKind == projectLoadKindFind {
		return ""
	}

	// Look for config in parent folders of config file
	result := b.computeConfigFileName(configFileName, true)

	// if f.IsOpenFile(ls.FileNameToDocumentURI(fileName)) {
	// 	ancestorConfigMap, ok := f.configFilesAncestorForOpenFiles[path]
	// 	if !ok {
	// 		ancestorConfigMap = make(map[string]string)
	// 		f.configFilesAncestorForOpenFiles[path] = ancestorConfigMap
	// 	}
	// 	ancestorConfigMap[configFileName] = result
	// }
	return result
}

func (b *projectCollectionBuilder) findOrAcquireConfig(
	// info *ScriptInfo,
	configFileName string,
	configFilePath tspath.Path,
	loadKind projectLoadKind,
) *tsoptions.ParsedCommandLine {
	switch loadKind {
	case projectLoadKindFind:
		// !!! is this right?
		return b.snapshot.configFileRegistry.getConfig(configFilePath)
	case projectLoadKindCreate:
		return b.configFileRegistryBuilder.acquireConfig(configFileName, configFilePath, nil)
	default:
		panic(fmt.Sprintf("unknown project load kind: %d", loadKind))
	}
}

func (b *projectCollectionBuilder) findOrCreateProject(
	configFileName string,
	configFilePath tspath.Path,
	loadKind projectLoadKind,
) *projectCollectionBuilderEntry {
	if loadKind == projectLoadKindFind {
		return &projectCollectionBuilderEntry{
			b:       b,
			project: b.base.configuredProjects[configFilePath],
			dirty:   false,
		}
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
		b.updateProject(entry)
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
			childConfig := b.findOrAcquireConfig(childConfigFileName, childConfigFilePath, loadKind)
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
	if ancestorConfigName := b.getAncestorConfigFileName(fileName, path, configFileName, loadKind); ancestorConfigName != "" {
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
	config := b.findOrAcquireConfig(configFileName, configFilePath, loadKind)
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
	if configFileName := b.getConfigFileNameForFile(fileName, path, loadKind); configFileName != "" {
		var result openScriptInfoProjectResult
		b.tryFindDefaultConfiguredProjectForScriptInfo(fileName, path, configFileName, loadKind, &result)
		if result.project == nil && result.fallbackDefault != nil {
			result.setProject(result.fallbackDefault)
		}
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

func (b *projectCollectionBuilder) findDefaultConfiguredProject(fileName string, path tspath.Path) *Project {
	if b.snapshot.IsOpenFile(path) {
		result := b.tryFindDefaultConfiguredProjectForOpenScriptInfo(fileName, path, projectLoadKindFind)
		if result != nil && result.project != nil /* !!! && !result.project.deferredClose */ {
			return result.project.project
		}
	}
	return nil
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
