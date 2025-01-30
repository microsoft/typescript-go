package project

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type projectLoadKind int

const (
	projectLoadKindFind projectLoadKind = iota
	projectLoadKindCreateReplay
	projectLoadKindCreate
	projectLoadKindReload
)

type assignProjectResult struct {
	configFileName string
	retainProjects map[*Project]projectLoadKind
	// configFileErrors []*ast.Diagnostic
}

type openFileArguments struct {
	FileName        string
	Content         string
	ScriptKind      core.ScriptKind
	HasMixedContent bool
	ProjectRootPath string
}

type changeFileArguments struct {
	FileName string
	Changes  []ls.TextChange
}

type ProjectService struct {
	host                ProjecServicetHost
	comparePathsOptions tspath.ComparePathsOptions

	configuredProjects map[tspath.Path]*Project

	documentRegistry *documentRegistry
	scriptInfos      map[tspath.Path]*scriptInfo
	openFiles        map[tspath.Path]string // values are projectRootPath, if provided
	// Contains all the deleted script info's version information so that
	// it does not reset when creating script info again
	filenameToScriptInfoVersion map[tspath.Path]int
	realpathToScriptInfos       map[tspath.Path]map[*scriptInfo]struct{}
}

func NewProjectService(host ProjecServicetHost) *ProjectService {
	return &ProjectService{
		host: host,
		comparePathsOptions: tspath.ComparePathsOptions{
			UseCaseSensitiveFileNames: host.FS().UseCaseSensitiveFileNames(),
			CurrentDirectory:          host.GetCurrentDirectory(),
		},
		documentRegistry: NewDocumentRegistry(tspath.ComparePathsOptions{
			UseCaseSensitiveFileNames: host.FS().UseCaseSensitiveFileNames(),
			CurrentDirectory:          host.GetCurrentDirectory(),
		}),
		scriptInfos:                 make(map[tspath.Path]*scriptInfo),
		filenameToScriptInfoVersion: make(map[tspath.Path]int),
		realpathToScriptInfos:       make(map[tspath.Path]map[*scriptInfo]struct{}),
	}
}

func (s *ProjectService) OpenClientFile(fileName string, fileContent string, scriptKind core.ScriptKind, projectRootPath string) {
	path := tspath.ToPath(fileName, s.host.GetCurrentDirectory(), s.host.FS().UseCaseSensitiveFileNames())
	existing := s.getScriptInfo(path)
	info := s.getOrCreateOpenScriptInfo(fileName, path, fileContent, scriptKind, projectRootPath)
	if existing == nil && info != nil && !info.isDynamic {
		// !!!
		// s.tryInvokeWildcardDirectories(info)
	}
	result := s.assignProjectToOpenedScriptInfo(info)
	s.cleanupProjectsAndScriptInfos(result.retainProjects, []tspath.Path{info.path})
}

func (s *ProjectService) ApplyChangesInOpenFiles(
	openFiles []openFileArguments,
	changedFiles []changeFileArguments,
	closedFiles []string,
) {
	var assignOrphanScriptInfoToInferredProject bool
	existingOpenScriptInfos := make([]*scriptInfo, 0, len(openFiles))
	openScriptInfos := make([]*scriptInfo, 0, len(openFiles))
	openScriptInfoPaths := make([]tspath.Path, 0, len(openFiles))

	for _, openFile := range openFiles {
		openFilePath := tspath.ToPath(openFile.FileName, s.host.GetCurrentDirectory(), s.host.FS().UseCaseSensitiveFileNames())
		existingOpenScriptInfos = append(existingOpenScriptInfos, s.getScriptInfo(openFilePath))
		openScriptInfos = append(openScriptInfos, s.getOrCreateOpenScriptInfo(openFile.FileName, openFilePath, openFile.Content, openFile.ScriptKind, openFile.ProjectRootPath))
		openScriptInfoPaths = append(openScriptInfoPaths, openFilePath)
	}

	for _, changedFile := range changedFiles {
		info := s.getScriptInfo(tspath.ToPath(changedFile.FileName, s.host.GetCurrentDirectory(), s.host.FS().UseCaseSensitiveFileNames()))
		if info == nil {
			panic("scriptInfo for changed file not found")
		}
		s.applyChangesToFile(info, changedFile.Changes)
	}

	for _, closedFile := range closedFiles {
		closedFilePath := tspath.ToPath(closedFile, s.host.GetCurrentDirectory(), s.host.FS().UseCaseSensitiveFileNames())
		assignOrphanScriptInfoToInferredProject = s.closeClientFile(closedFilePath, true /*skipAssignOrphanScriptInfosToInferredProject*/) || assignOrphanScriptInfoToInferredProject
	}

	retainedProjects := make(map[*Project]projectLoadKind)
	for i, existing := range existingOpenScriptInfos {
		if existing == nil && openScriptInfos[i] != nil && !openScriptInfos[i].isDynamic {
			// !!!
			// s.tryInvokeWildcardDirectories(openScriptInfos[i])
		}
	}
	for _, info := range openScriptInfos {
		for project, loadKind := range s.assignProjectToOpenedScriptInfo(info).retainProjects {
			retainedProjects[project] = loadKind
		}
	}

	if assignOrphanScriptInfoToInferredProject {
		// !!!
		// s.assignOrphanScriptInfoToInferredProject()
	}

	if len(openScriptInfos) > 0 {
		s.cleanupProjectsAndScriptInfos(retainedProjects, openScriptInfoPaths)
	}
}

func (s *ProjectService) applyChangesToFile(info *scriptInfo, changes []ls.TextChange) {
	for _, change := range changes {
		info.editContent(change)
	}
}

func (s *ProjectService) closeClientFile(path tspath.Path, skipAssignOrphanScriptInfosToInferredProject bool) bool {
	if info := s.getScriptInfo(path); info != nil {
		return s.closeOpenFile(info, skipAssignOrphanScriptInfosToInferredProject)
	}
	return false
}

func (s *ProjectService) closeOpenFile(info *scriptInfo, skipAssignOrphanScriptInfosToInferredProject bool) bool {
	fileExists := !info.isDynamic && s.host.FS().FileExists(info.fileName)
	info.close(fileExists)
	// s.stopWatchingConfigFilesForScriptInfo(info)

	var ensureProjectsForOpenFiles bool
	// !!! collect all projects that should be removed

	delete(s.openFiles, info.path)

	if !skipAssignOrphanScriptInfosToInferredProject && ensureProjectsForOpenFiles {
		// !!!
		// s.assignOrphanScriptInfoToInferredProject()
	}

	// Cleanup script infos that arent part of any project (eg. those could be closed script infos not referenced by any project)
	// is postponed to next file open so that if file from same project is opened,
	// we wont end up creating same script infos

	// If the current info is being just closed - add the watcher file to track changes
	// But if file was deleted, handle that part
	if fileExists {
		// s.watchClosedScriptInfo(info)
	} else {
		// s.handleDeletedFile(info /*deferredDelete*/, false)
	}
	return ensureProjectsForOpenFiles
}

func (s *ProjectService) handleDeletedFile(info *scriptInfo, deferredDelete bool) {
	if info.isOpen {
		panic("cannot delete an open file")
	}

	s.delayUpdateProjectGraphs(info.containingProjects, false /*clearSourceMapperCache*/)
	// !!!
	// s.handleSourceMapProjects(info)
	info.detachAllProjects()
	if deferredDelete {
		info.delayReloadNonMixedContentFile()
		info.deferredDelete = true
	} else {
		s.deleteScriptInfo(info)
	}
}

func (s *ProjectService) deleteScriptInfo(info *scriptInfo) {
	if info.isOpen {
		panic("cannot delete an open file")
	}
	delete(s.scriptInfos, info.path)
	s.filenameToScriptInfoVersion[info.path] = info.version
	// !!!
	// s.stopWatchingScriptInfo(info)
	if realpath, ok := info.getRealpathIfDifferent(); ok {
		delete(s.realpathToScriptInfos[realpath], info)
	}
	// !!! closeSourceMapFileWatcher
}

func (s *ProjectService) recordSymlink(info *scriptInfo) {
	if scriptInfos, ok := s.realpathToScriptInfos[info.realpath]; ok {
		scriptInfos[info] = struct{}{}
	} else {
		scriptInfos = make(map[*scriptInfo]struct{})
		scriptInfos[info] = struct{}{}
		s.realpathToScriptInfos[info.realpath] = scriptInfos
	}
}

func (s *ProjectService) delayUpdateProjectGraphs(projects []*Project, clearSourceMapperCache bool) {
	for _, project := range projects {
		if clearSourceMapperCache {
			project.clearSourceMapperCache()
		}
		s.delayUpdateProjectGraph(project)
	}
}

func (s *ProjectService) delayUpdateProjectGraph(project *Project) {
	if project.deferredClose {
		return
	}
	project.markAsDirty()
	if project.kind == ProjectKindAutoImportProvider || project.kind == ProjectKindAuxiliary {
		return
	}
	// !!! throttle
	project.updateIfDirty()
}

func (s *ProjectService) getScriptInfo(path tspath.Path) *scriptInfo {
	if info, ok := s.scriptInfos[path]; ok && !info.deferredDelete {
		return info
	}
	return nil
}

func (s *ProjectService) getOrCreateScriptInfoNotOpenedByClient(fileName string, path tspath.Path, scriptKind core.ScriptKind) *scriptInfo {
	if tspath.IsRootedDiskPath(fileName) /* !!! || isDynamicFileName(fileName) */ {
		return s.getOrCreateScriptInfoWorker(fileName, path, scriptKind, false /*openedByClient*/, "" /*fileContent*/, false /*deferredDeleteOk*/)
	}
	// !!!
	// This is non rooted path with different current directory than project service current directory
	// Only paths recognized are open relative file paths
	// const info = this.openFilesWithNonRootedDiskPath.get(this.toCanonicalFileName(fileName))
	// if info {
	// 	return info
	// }

	// This means triple slash references wont be resolved in dynamic and unsaved files
	// which is intentional since we dont know what it means to be relative to non disk files
	return nil
}

func (s *ProjectService) getOrCreateOpenScriptInfo(fileName string, path tspath.Path, fileContent string, scriptKind core.ScriptKind, projectRootPath string) *scriptInfo {
	info := s.getOrCreateScriptInfoWorker(fileName, path, scriptKind, true /*openedByClient*/, fileContent, true /*deferredDeleteOk*/)
	s.openFiles[info.path] = projectRootPath
	return info
}

func (s *ProjectService) getOrCreateScriptInfoWorker(fileName string, path tspath.Path, scriptKind core.ScriptKind, openedByClient bool, fileContent string, deferredDeleteOk bool) *scriptInfo {
	info, ok := s.scriptInfos[path]
	if ok {
		if info.deferredDelete {
			if !openedByClient && !s.host.FS().FileExists(fileName) {
				// If the file is not opened by client and the file does not exist on the disk, return
				return core.IfElse(deferredDeleteOk, info, nil)
			}
			info.deferredDelete = false
		}
	} else if !openedByClient && !s.host.FS().FileExists(fileName) {
		return nil
	} else {
		info = newScriptInfo(fileName, path, scriptKind)
		if prevVersion, ok := s.filenameToScriptInfoVersion[path]; ok {
			info.version = prevVersion
			delete(s.filenameToScriptInfoVersion, path)
		}
		s.scriptInfos[path] = info
		// !!!
		// if !openedByClient {
		// 	this.watchClosedScriptInfo(info)
		// } else if !isRootedDiskPath(fileName) && (!isDynamic || this.currentDirectory != currentDirectory) {
		// 	// File that is opened by user but isn't rooted disk path
		// 	this.openFilesWithNonRootedDiskPath.set(this.toCanonicalFileName(fileName), info)
		// }
	}

	if openedByClient {
		// Opening closed script info
		// either it was created just now, or was part of projects but was closed
		// !!!
		// s.stopWatchingScriptInfo(info)
		info.open(fileContent)
	}
	return info
}

func (s *ProjectService) configFileExists(configFilename string) bool {
	// !!! convoluted cache goes here
	return s.host.FS().FileExists(configFilename)
}

func (s *ProjectService) getConfigFileNameForFile(info *scriptInfo, findFromCacheOnly bool) string {
	// !!!
	// const fromCache = this.getConfigFileNameForFileFromCache(info, findFromCacheOnly);
	// if (fromCache !== undefined) return fromCache || undefined;
	// if (findFromCacheOnly) return undefined;
	//
	// !!!
	// good grief, this is convoluted. I'm skipping so much stuff right now
	projectRootPath := s.openFiles[info.path]
	if info.isDynamic {
		return ""
	}

	searchPath := tspath.GetDirectoryPath(info.fileName)
	fileName, _ := tspath.ForEachAncestorDirectory(searchPath, func(directory string) (result string, stop bool) {
		tsconfigPath := tspath.CombinePaths(directory, "tsconfig.json")
		if s.configFileExists(tsconfigPath) {
			return tsconfigPath, true
		}
		if strings.HasSuffix(directory, "/node_modules") {
			return "", true
		}
		if projectRootPath != "" && !tspath.ContainsPath(projectRootPath, directory, s.comparePathsOptions) {
			return "", true
		}
		return "", false
	})
	return fileName
}

func (s *ProjectService) findConfiguredProjectByName(configFilePath tspath.Path, includeDeferredClosedProjects bool) *Project {
	if result, ok := s.configuredProjects[configFilePath]; ok {
		if includeDeferredClosedProjects || !result.deferredClose {
			return result
		}
	}
	return nil
}

func (s *ProjectService) createConfiguredProject(configFileName string, configFilePath tspath.Path) *Project {
	// !!! config file existence cache stuff omitted
	project := NewConfiguredProject(configFileName, configFilePath, s)
	s.configuredProjects[configFilePath] = project
	// !!!
	// s.createConfigFileWatcherForParsedConfig(configFileName, configFilePath, project)
	return project
}

func (s *ProjectService) findCreateOrReloadConfiguredProject(configFileName string, projectLoadKind projectLoadKind, includeDeferredClosedProjects bool) *Project {
	// !!! many such things omitted
	configFilePath := tspath.ToPath(configFileName, s.host.GetCurrentDirectory(), s.host.FS().UseCaseSensitiveFileNames())
	project := s.findConfiguredProjectByName(configFilePath, includeDeferredClosedProjects)
	switch projectLoadKind {
	case projectLoadKindFind, projectLoadKindCreateReplay:
		return project
	case projectLoadKindCreate:
		if project == nil {
			project = s.createConfiguredProject(configFileName, configFilePath)
		}
	case projectLoadKindReload:
		if project == nil {
			project = s.createConfiguredProject(configFileName, configFilePath)
		}
	default:
		panic("unhandled projectLoadKind")
	}
	return project
}

func (s *ProjectService) tryFindDefaultConfiguredProjectForOpenScriptInfo(info *scriptInfo, projectLoadKind projectLoadKind, includeDeferredClosedProjects bool) *Project {
	findConfigFromCacheOnly := projectLoadKind == projectLoadKindFind || projectLoadKind == projectLoadKindCreateReplay
	if configFileName := s.getConfigFileNameForFile(info, findConfigFromCacheOnly); configFileName != "" {
		// !!! Maybe this recently added "optimized" stuff can be simplified?
		// const optimizedKind = toConfiguredProjectLoadOptimized(kind);
		return s.findCreateOrReloadConfiguredProject(configFileName, projectLoadKind, includeDeferredClosedProjects)
	}
	return nil
}

func (s *ProjectService) tryFindDefaultConfiguredProjectAndLoadAncestorsForOpenScriptInfo(info *scriptInfo, projectLoadKind projectLoadKind) *Project {
	includeDeferredClosedProjects := projectLoadKind == projectLoadKindFind
	result := s.tryFindDefaultConfiguredProjectForOpenScriptInfo(info, projectLoadKind, includeDeferredClosedProjects)
	// !!! I don't even know what an ancestor project is
	return result
}

func (s *ProjectService) assignProjectToOpenedScriptInfo(info *scriptInfo) assignProjectResult {
	var result assignProjectResult
	if project := s.tryFindDefaultConfiguredProjectAndLoadAncestorsForOpenScriptInfo(info, projectLoadKindCreate); project != nil {
		result.configFileName = project.configFileName
		// result.configFileErrors = project.getAllProjectErrors()
	}
	for _, project := range info.containingProjects {
		project.updateIfDirty()
	}
	if info.isOrphan() {
		// !!!
		// more new "optimized" stuff
		// s.assignOrphanScriptInfoToInferredProject(info)
	}
	return result
}

func (s *ProjectService) cleanupProjectsAndScriptInfos(toRetainConfiguredProjects map[*Project]projectLoadKind, openFilesWithRetainedConfiguredProject []tspath.Path) {
	// !!!
}
