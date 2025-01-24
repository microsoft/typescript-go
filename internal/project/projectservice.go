package project

import (
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type ProjectService struct {
	host ProjecServicetHost

	documentRegistry *documentRegistry
	scriptInfos      map[tspath.Path]*scriptInfo
	// Contains all the deleted script info's version information so that
	// it does not reset when creating script info again
	filenameToScriptInfoVersion map[tspath.Path]int
}

func NewProjectService(host ProjecServicetHost) *ProjectService {
	return &ProjectService{
		host: host,
		documentRegistry: NewDocumentRegistry(tspath.ComparePathsOptions{
			UseCaseSensitiveFileNames: host.FS().UseCaseSensitiveFileNames(),
			CurrentDirectory:          host.GetCurrentDirectory(),
		}),
		scriptInfos:                 make(map[tspath.Path]*scriptInfo),
		filenameToScriptInfoVersion: make(map[tspath.Path]int),
	}
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
