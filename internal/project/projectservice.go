package project

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/core"
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
	retainProjects []*Project
	// configFileErrors []*ast.Diagnostic
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
	}
}

func (s *ProjectService) OpenClientFile(fileName string, fileContent string, scriptKind core.ScriptKind, projectRootPath string) {
	path := tspath.ToPath(fileName, s.host.GetCurrentDirectory(), s.host.FS().UseCaseSensitiveFileNames())
	existing := s.getScriptInfo(path)
	info := s.getOrCreateOpenScriptInfo(fileName, fileContent, scriptKind, projectRootPath)
	if existing == nil && info != nil && !info.isDynamic {
		// !!!
		// s.tryInvokeWildcardDirectories(info)
	}
	result := s.assignProjectToOpenedScriptInfo(info)
	s.cleanupProjectsAndScriptInfos(result.retainProjects, []tspath.Path{info.path})
}

func (s *ProjectService) getScriptInfo(path tspath.Path) *scriptInfo {
	if info, ok := s.scriptInfos[path]; ok && !info.deferredDelete {
		return info
	}
	return nil
}

func (s *ProjectService) getOrCreateScriptInfoNotOpenedByClient(fileName string, scriptKind core.ScriptKind) *scriptInfo {
	if tspath.IsRootedDiskPath(fileName) /* !!! || isDynamicFileName(fileName) */ {
		return s.getOrCreateScriptInfoWorker(fileName, scriptKind, false /*openedByClient*/, "" /*fileContent*/, false /*deferredDeleteOk*/)
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

func (s *ProjectService) getOrCreateOpenScriptInfo(fileName string, fileContent string, scriptKind core.ScriptKind, projectRootPath string) *scriptInfo {
	info := s.getOrCreateScriptInfoWorker(fileName, scriptKind, true /*openedByClient*/, fileContent, true /*deferredDeleteOk*/)
	s.openFiles[info.path] = projectRootPath
	return info
}

func (s *ProjectService) getOrCreateScriptInfoWorker(fileName string, scriptKind core.ScriptKind, openedByClient bool, fileContent string, deferredDeleteOk bool) *scriptInfo {
	path := tspath.ToPath(fileName, s.host.GetCurrentDirectory(), s.host.FS().UseCaseSensitiveFileNames())
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

func (s *ProjectService) cleanupProjectsAndScriptInfos(toRetainConfiguredProjects []*Project, openFilesWithRetainedConfiguredProject []tspath.Path) {
	// !!!
}
