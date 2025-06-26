package projectv2

import (
	"fmt"
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

// type defaultProjectFinder struct {
// 	snapshot                        *Snapshot
// 	configFileForOpenFiles          map[tspath.Path]string            // default config project for open files
// 	configFilesAncestorForOpenFiles map[tspath.Path]map[string]string // ancestor config file for open files
// }

func (s *Snapshot) computeConfigFileName(fileName string, skipSearchInDirectoryOfFile bool) string {
	searchPath := tspath.GetDirectoryPath(fileName)
	result, _ := tspath.ForEachAncestorDirectory(searchPath, func(directory string) (result string, stop bool) {
		tsconfigPath := tspath.CombinePaths(directory, "tsconfig.json")
		if !skipSearchInDirectoryOfFile && s.compilerFS.FileExists(tsconfigPath) {
			return tsconfigPath, true
		}
		jsconfigPath := tspath.CombinePaths(directory, "jsconfig.json")
		if !skipSearchInDirectoryOfFile && s.compilerFS.FileExists(jsconfigPath) {
			return jsconfigPath, true
		}
		if strings.HasSuffix(directory, "/node_modules") {
			return "", true
		}
		skipSearchInDirectoryOfFile = false
		return "", false
	})
	s.Logf("computeConfigFileName:: File: %s:: Result: %s", fileName, result)
	return result
}

func (s *Snapshot) getConfigFileNameForFile(fileName string, path tspath.Path, loadKind projectLoadKind) string {
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

	configName := s.computeConfigFileName(fileName, false)

	// if f.IsOpenFile(ls.FileNameToDocumentURI(fileName)) {
	// 	f.configFileForOpenFiles[path] = configName
	// }
	return configName
}

func (s *Snapshot) getAncestorConfigFileName(fileName string, path tspath.Path, configFileName string, loadKind projectLoadKind) string {
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
	result := s.computeConfigFileName(configFileName, true)

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

func (s *Snapshot) findOrAcquireConfig(
	// info *ScriptInfo,
	configFileName string,
	configFilePath tspath.Path,
	loadKind projectLoadKind,
) *tsoptions.ParsedCommandLine {
	switch loadKind {
	case projectLoadKindFind:
		return s.configFileRegistry.getConfig(configFilePath)
	case projectLoadKindCreate:
		return s.configFileRegistry.acquireConfig(configFileName, configFilePath, nil)
	default:
		panic(fmt.Sprintf("unknown project load kind: %d", loadKind))
	}
}

func (s *Snapshot) findOrCreateProject(
	configFileName string,
	configFilePath tspath.Path,
	loadKind projectLoadKind,
) *Project {
	project := s.configuredProjects[configFilePath]
	if project == nil {
		if loadKind == projectLoadKindFind {
			return nil
		}
		project = NewConfiguredProject(configFileName, configFilePath, s)
	}
	return project
}

func (s *Snapshot) isDefaultConfigForScript(
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
	project := s.findOrCreateProject(configFileName, configFilePath, loadKind)
	return s.isDefaultProject(scriptFileName, scriptPath, project, loadKind, result)
}

func (s *Snapshot) isDefaultProject(
	fileName string,
	path tspath.Path,
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
	if project.containsFile(path) {
		if !project.IsSourceFromProjectReference(path) {
			result.setProject(project)
			return true
		} else if !result.hasFallbackDefault() {
			// Use this project as default if no other project is found
			result.setFallbackDefault(project)
		}
	}
	return false
}

func (s *Snapshot) tryFindDefaultConfiguredProjectFromReferences(
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
	s.tryFindDefaultConfiguredProjectFromReferencesWorker(fileName, path, config, loadKind, result, wg)
	wg.RunAndWait()
	return result.isDone()
}

func (s *Snapshot) tryFindDefaultConfiguredProjectFromReferencesWorker(
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
			childConfigFilePath := s.toPath(childConfigFileName)
			childConfig := s.findOrAcquireConfig(childConfigFileName, childConfigFilePath, loadKind)
			if childConfig == nil || s.isDefaultConfigForScript(fileName, path, childConfigFileName, childConfigFilePath, childConfig, loadKind, result) {
				return
			}
			// Search in references if we cant find default project in current config
			s.tryFindDefaultConfiguredProjectFromReferencesWorker(fileName, path, childConfig, loadKind, result, wg)
		})
	}
}

func (s *Snapshot) tryFindDefaultConfiguredProjectFromAncestor(
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
	if ancestorConfigName := s.getAncestorConfigFileName(fileName, path, configFileName, loadKind); ancestorConfigName != "" {
		return s.tryFindDefaultConfiguredProjectForScriptInfo(fileName, path, ancestorConfigName, loadKind, result)
	}
	return false
}

func (s *Snapshot) tryFindDefaultConfiguredProjectForScriptInfo(
	fileName string,
	path tspath.Path,
	configFileName string,
	loadKind projectLoadKind,
	result *openScriptInfoProjectResult,
) bool {
	// Lookup from parsedConfig if available
	configFilePath := s.toPath(configFileName)
	config := s.findOrAcquireConfig(configFileName, configFilePath, loadKind)
	if config != nil {
		if config.CompilerOptions().Composite == core.TSTrue {
			if s.isDefaultConfigForScript(fileName, path, configFileName, configFilePath, config, loadKind, result) {
				return true
			}
		} else if len(config.FileNames()) > 0 {
			project := s.findOrCreateProject(configFileName, configFilePath, loadKind)
			if s.isDefaultProject(fileName, path, project, loadKind, result) {
				return true
			}
		}
		// Lookup in references
		if s.tryFindDefaultConfiguredProjectFromReferences(fileName, path, config, loadKind, result) {
			return true
		}
	}
	// Lookup in ancestor projects
	if s.tryFindDefaultConfiguredProjectFromAncestor(fileName, path, configFileName, config, loadKind, result) {
		return true
	}
	return false
}

func (s *Snapshot) tryFindDefaultConfiguredProjectForOpenScriptInfo(
	fileName string,
	path tspath.Path,
	loadKind projectLoadKind,
) *openScriptInfoProjectResult {
	if configFileName := s.getConfigFileNameForFile(fileName, path, loadKind); configFileName != "" {
		var result openScriptInfoProjectResult
		s.tryFindDefaultConfiguredProjectForScriptInfo(fileName, path, configFileName, loadKind, &result)
		if result.project == nil && result.fallbackDefault != nil {
			result.setProject(result.fallbackDefault)
		}
		return &result
	}
	return nil
}

func (s *Snapshot) tryFindDefaultConfiguredProjectAndLoadAncestorsForOpenScriptInfo(
	fileName string,
	path tspath.Path,
	loadKind projectLoadKind,
) *openScriptInfoProjectResult {
	result := s.tryFindDefaultConfiguredProjectForOpenScriptInfo(fileName, path, loadKind)
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

func (s *Snapshot) findDefaultConfiguredProject(fileName string, path tspath.Path) *Project {
	if s.IsOpenFile(path) {
		result := s.tryFindDefaultConfiguredProjectForOpenScriptInfo(fileName, path, projectLoadKindFind)
		if result != nil && result.project != nil /* !!! && !result.project.deferredClose */ {
			return result.project
		}
	}
	return nil
}

type openScriptInfoProjectResult struct {
	projectMu         sync.RWMutex
	project           *Project
	fallbackDefaultMu sync.RWMutex
	fallbackDefault   *Project // use this if we cant find actual project
	seenProjects      collections.SyncMap[*Project, projectLoadKind]
	seenConfigs       collections.SyncMap[tspath.Path, projectLoadKind]
}

func (r *openScriptInfoProjectResult) addSeenProject(project *Project, loadKind projectLoadKind) bool {
	if kind, loaded := r.seenProjects.LoadOrStore(project, loadKind); loaded {
		if kind >= loadKind {
			return false
		}
		r.seenProjects.Store(project, loadKind)
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

func (r *openScriptInfoProjectResult) setProject(project *Project) {
	r.projectMu.Lock()
	defer r.projectMu.Unlock()
	if r.project == nil {
		r.project = project
	}
}

func (r *openScriptInfoProjectResult) hasFallbackDefault() bool {
	r.fallbackDefaultMu.RLock()
	defer r.fallbackDefaultMu.RUnlock()
	return r.fallbackDefault != nil
}

func (r *openScriptInfoProjectResult) setFallbackDefault(project *Project) {
	r.fallbackDefaultMu.Lock()
	defer r.fallbackDefaultMu.Unlock()
	if r.fallbackDefault == nil {
		r.fallbackDefault = project
	}
}
