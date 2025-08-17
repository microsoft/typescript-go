package project

import (
	"fmt"
	"strings"
	"sync"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type defaultProjectFinder struct {
	service                         *Service
	configFileForOpenFiles          map[tspath.Path]string            // default config project for open files
	configFilesAncestorForOpenFiles map[tspath.Path]map[string]string // ancestor config file for open files
}

func (f *defaultProjectFinder) computeConfigFileName(fileName string, info *ScriptInfo, skipSearchInDirectoryOfFile bool) string {
	projectRootPath := f.service.openFiles[info.path]
	searchPath := tspath.GetDirectoryPath(fileName)
	result, _ := tspath.ForEachAncestorDirectory(searchPath, func(directory string) (result string, stop bool) {
		tsconfigPath := tspath.CombinePaths(directory, "tsconfig.json")
		if !skipSearchInDirectoryOfFile && f.service.FS().FileExists(tsconfigPath) {
			return tsconfigPath, true
		}
		jsconfigPath := tspath.CombinePaths(directory, "jsconfig.json")
		if !skipSearchInDirectoryOfFile && f.service.FS().FileExists(jsconfigPath) {
			return jsconfigPath, true
		}
		if strings.HasSuffix(directory, "/node_modules") {
			return "", true
		}
		if projectRootPath != "" && !tspath.ContainsPath(projectRootPath, directory, f.service.comparePathsOptions) {
			return "", true
		}
		skipSearchInDirectoryOfFile = false
		return "", false
	})
	f.service.logf("getConfigFileNameForFile:: File: %s ProjectRootPath: %s:: Result: %s", fileName, projectRootPath, result)
	return result
}

func (f *defaultProjectFinder) getConfigFileNameForFile(info *ScriptInfo, loadKind projectLoadKind) string {
	if info.isDynamic {
		return ""
	}

	configName, ok := f.configFileForOpenFiles[info.path]
	if ok {
		return configName
	}

	if loadKind == projectLoadKindFind {
		return ""
	}

	fileName := f.computeConfigFileName(info.fileName, info, false)

	if _, ok := f.service.openFiles[info.path]; ok {
		f.configFileForOpenFiles[info.path] = fileName
	}
	return fileName
}

func (f *defaultProjectFinder) getAncestorConfigFileName(info *ScriptInfo, configFileName string, loadKind projectLoadKind) string {
	if info.isDynamic {
		return ""
	}

	ancestorConfigMap, ok := f.configFilesAncestorForOpenFiles[info.path]
	if ok {
		ancestorConfigName, found := ancestorConfigMap[configFileName]
		if found {
			return ancestorConfigName
		}
	}

	if loadKind == projectLoadKindFind {
		return ""
	}

	// Look for config in parent folders of config file
	fileName := f.computeConfigFileName(configFileName, info, true)

	if _, ok := f.service.openFiles[info.path]; ok {
		ancestorConfigMap, ok := f.configFilesAncestorForOpenFiles[info.path]
		if !ok {
			ancestorConfigMap = make(map[string]string)
			f.configFilesAncestorForOpenFiles[info.path] = ancestorConfigMap
		}
		ancestorConfigMap[configFileName] = fileName
	}
	return fileName
}

func (f *defaultProjectFinder) findOrAcquireConfig(
	info *ScriptInfo,
	configFileName string,
	configFilePath tspath.Path,
	loadKind projectLoadKind,
) *tsoptions.ParsedCommandLine {
	switch loadKind {
	case projectLoadKindFind:
		return f.service.configFileRegistry.getConfig(configFilePath)
	case projectLoadKindCreate:
		return f.service.configFileRegistry.acquireConfig(configFileName, configFilePath, nil, info)
	default:
		panic(fmt.Sprintf("unknown project load kind: %d", loadKind))
	}
}

func (f *defaultProjectFinder) findOrCreateProject(
	configFileName string,
	configFilePath tspath.Path,
	loadKind projectLoadKind,
) *Project {
	project := f.service.ConfiguredProject(configFilePath)
	if project == nil {
		if loadKind == projectLoadKindFind {
			return nil
		}
		project = f.service.createConfiguredProject(configFileName, configFilePath)
	}
	return project
}

func (f *defaultProjectFinder) isDefaultConfigForScriptInfo(
	info *ScriptInfo,
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
	if !config.MatchesFileName(info.fileName) {
		return false
	}

	// Ensure the project is uptodate and created since the file may belong to this project
	project := f.findOrCreateProject(configFileName, configFilePath, loadKind)
	return f.isDefaultProject(info, project, loadKind, result)
}

func (f *defaultProjectFinder) isDefaultProject(
	info *ScriptInfo,
	project *Project,
	loadKind projectLoadKind,
	result *openScriptInfoProjectResult,
) bool {
	if project == nil {
		return false
	}

	// Skip already looked up projects
	if !result.addSeenProject(project, loadKind) {
		return false
	}
	// Make sure project is upto date when in create mode
	if loadKind == projectLoadKindCreate {
		project.updateGraph()
	}
	// If script info belongs to this project, use this as default config project
	if project.containsScriptInfo(info) {
		if !project.isSourceFromProjectReference(info) {
			result.setProject(project)
			return true
		} else if !result.hasFallbackDefault() {
			// Use this project as default if no other project is found
			result.setFallbackDefault(project)
		}
	}
	return false
}

func (f *defaultProjectFinder) tryFindDefaultConfiguredProjectFromReferences(
	info *ScriptInfo,
	config *tsoptions.ParsedCommandLine,
	loadKind projectLoadKind,
	result *openScriptInfoProjectResult,
) bool {
	if len(config.ProjectReferences()) == 0 {
		return false
	}

	type configItem struct {
		fileName string
		path     tspath.Path
		config   *tsoptions.ParsedCommandLine
		loadKind projectLoadKind
	}

	type visitKey struct {
		path     tspath.Path
		loadKind projectLoadKind
	}

	var queue []configItem
	visited := make(map[visitKey]bool)

	initialLoadKind := loadKind
	if config.CompilerOptions().DisableReferencedProjectLoad.IsTrue() {
		initialLoadKind = projectLoadKindFind
	}

	for _, childConfigFileName := range config.ResolvedProjectReferencePaths() {
		childConfigFilePath := f.service.toPath(childConfigFileName)
		key := visitKey{path: childConfigFilePath, loadKind: initialLoadKind}
		if !visited[key] {
			visited[key] = true
			childConfig := f.findOrAcquireConfig(info, childConfigFileName, childConfigFilePath, initialLoadKind)
			if childConfig != nil {
				queue = append(queue, configItem{
					fileName: childConfigFileName,
					path:     childConfigFilePath,
					config:   childConfig,
					loadKind: initialLoadKind,
				})
			}
		}
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if f.isDefaultConfigForScriptInfo(info, current.fileName, current.path, current.config, current.loadKind, result) {
			return true
		}

		if result.isDone() {
			return true
		}

		childLoadKind := current.loadKind
		if current.config.CompilerOptions().DisableReferencedProjectLoad.IsTrue() {
			childLoadKind = projectLoadKindFind
		}

		for _, childConfigFileName := range current.config.ResolvedProjectReferencePaths() {
			childConfigFilePath := f.service.toPath(childConfigFileName)
			key := visitKey{path: childConfigFilePath, loadKind: childLoadKind}

			if !visited[key] {
				visited[key] = true
				childConfig := f.findOrAcquireConfig(info, childConfigFileName, childConfigFilePath, childLoadKind)
				if childConfig != nil {
					queue = append(queue, configItem{
						fileName: childConfigFileName,
						path:     childConfigFilePath,
						config:   childConfig,
						loadKind: childLoadKind,
					})
				}
			}
		}
	}

	return result.isDone()
}

func (f *defaultProjectFinder) tryFindDefaultConfiguredProjectFromAncestor(
	info *ScriptInfo,
	configFileName string,
	config *tsoptions.ParsedCommandLine,
	loadKind projectLoadKind,
	result *openScriptInfoProjectResult,
) bool {
	if config != nil && config.CompilerOptions().DisableSolutionSearching.IsTrue() {
		return false
	}
	if ancestorConfigName := f.getAncestorConfigFileName(info, configFileName, loadKind); ancestorConfigName != "" {
		return f.tryFindDefaultConfiguredProjectForScriptInfo(info, ancestorConfigName, loadKind, result)
	}
	return false
}

func (f *defaultProjectFinder) tryFindDefaultConfiguredProjectForScriptInfo(
	info *ScriptInfo,
	configFileName string,
	loadKind projectLoadKind,
	result *openScriptInfoProjectResult,
) bool {
	// Lookup from parsedConfig if available
	configFilePath := f.service.toPath(configFileName)
	config := f.findOrAcquireConfig(info, configFileName, configFilePath, loadKind)
	if config != nil {
		if config.CompilerOptions().Composite == core.TSTrue {
			if f.isDefaultConfigForScriptInfo(info, configFileName, configFilePath, config, loadKind, result) {
				return true
			}
		} else if len(config.FileNames()) > 0 {
			project := f.findOrCreateProject(configFileName, configFilePath, loadKind)
			if f.isDefaultProject(info, project, loadKind, result) {
				return true
			}
		}
		// Lookup in references
		if f.tryFindDefaultConfiguredProjectFromReferences(info, config, loadKind, result) {
			return true
		}
	}
	// Lookup in ancestor projects
	if f.tryFindDefaultConfiguredProjectFromAncestor(info, configFileName, config, loadKind, result) {
		return true
	}
	return false
}

func (f *defaultProjectFinder) tryFindDefaultConfiguredProjectForOpenScriptInfo(
	info *ScriptInfo,
	loadKind projectLoadKind,
) *openScriptInfoProjectResult {
	if configFileName := f.getConfigFileNameForFile(info, loadKind); configFileName != "" {
		var result openScriptInfoProjectResult
		f.tryFindDefaultConfiguredProjectForScriptInfo(info, configFileName, loadKind, &result)
		if result.project == nil && result.fallbackDefault != nil {
			result.setProject(result.fallbackDefault)
		}
		return &result
	}
	return nil
}

func (f *defaultProjectFinder) tryFindDefaultConfiguredProjectAndLoadAncestorsForOpenScriptInfo(
	info *ScriptInfo,
	projectLoadKind projectLoadKind,
) *openScriptInfoProjectResult {
	result := f.tryFindDefaultConfiguredProjectForOpenScriptInfo(info, projectLoadKind)
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

func (f *defaultProjectFinder) findDefaultConfiguredProject(scriptInfo *ScriptInfo) *Project {
	if f.service.isOpenFile(scriptInfo) {
		result := f.tryFindDefaultConfiguredProjectForOpenScriptInfo(scriptInfo, projectLoadKindFind)
		if result != nil && result.project != nil && !result.project.deferredClose {
			return result.project
		}
	}
	return nil
}

type openScriptInfoProjectResult struct {
	mu              sync.Mutex
	project         *Project
	fallbackDefault *Project // use this if we cant find actual project
	seenProjects    map[*Project]projectLoadKind
	seenConfigs     map[tspath.Path]projectLoadKind
}

func (r *openScriptInfoProjectResult) addSeenProject(project *Project, loadKind projectLoadKind) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.seenProjects == nil {
		r.seenProjects = make(map[*Project]projectLoadKind)
	}

	if kind, exists := r.seenProjects[project]; exists {
		if kind >= loadKind {
			return false
		}
	}
	r.seenProjects[project] = loadKind
	return true
}

func (r *openScriptInfoProjectResult) addSeenConfig(configPath tspath.Path, loadKind projectLoadKind) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.seenConfigs == nil {
		r.seenConfigs = make(map[tspath.Path]projectLoadKind)
	}

	if kind, exists := r.seenConfigs[configPath]; exists {
		if kind >= loadKind {
			return false
		}
	}
	r.seenConfigs[configPath] = loadKind
	return true
}

func (r *openScriptInfoProjectResult) isDone() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.project != nil
}

func (r *openScriptInfoProjectResult) setProject(project *Project) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.project == nil {
		r.project = project
	}
}

func (r *openScriptInfoProjectResult) hasFallbackDefault() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.fallbackDefault != nil
}

func (r *openScriptInfoProjectResult) setFallbackDefault(project *Project) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.fallbackDefault == nil {
		r.fallbackDefault = project
	}
}
