package lsp

import (
	"context"
	"sync"
	"time"

	"github.com/microsoft/typescript-go/internal/project"
)

type SelfManagedProjectInfo struct {
	Project          *project.Project
	LastAccessed     int64
	FilesLastUpdated map[string]int64
}

var (
	mu                sync.Mutex
	projectInfoByName = make(map[string]*SelfManagedProjectInfo)
)

func GetAllSelfManagedProjects(s *Server, ctx context.Context) []*project.Project {
	cleanupStaleProjects(s, ctx)

	mu.Lock()
	defer mu.Unlock()

	var projects []*project.Project
	for _, info := range projectInfoByName {
		projects = append(projects, info.Project)
	}
	return projects
}

func IsSelfManagedProject(projectFileName string) bool {
	mu.Lock()
	defer mu.Unlock()
	return projectInfoByName[projectFileName] != nil
}

func GetOrCreateSelfManagedProjectForFile(s *Server, projectFileName string, file string, ctx context.Context) *project.Project {
	defer cleanupStaleProjects(s, ctx)
	return getProjectForFileImpl(s, projectFileName, file, ctx)
}

func getProjectForFileImpl(s *Server, projectFileName string, file string, ctx context.Context) *project.Project {
	nowMs := time.Now().UnixMilli()
	fileMod := getFileModTimeMs(s, file)

	mu.Lock()
	defer mu.Unlock()

	if info := projectInfoByName[projectFileName]; info != nil {
		prevMod, hadPrev := info.FilesLastUpdated[file]
		if !hadPrev || fileMod == prevMod {
			info.LastAccessed = nowMs
			return info.Project
		}

		closeProject(s, projectFileName, ctx)
		delete(projectInfoByName, projectFileName)
	}

	newProject := createNewSelfManagedProject(s, projectFileName, file, ctx)
	if newProject == nil {
		return nil
	}

	newInfo := &SelfManagedProjectInfo{
		Project:          newProject,
		LastAccessed:     nowMs,
		FilesLastUpdated: make(map[string]int64),
	}
	newInfo.FilesLastUpdated[file] = fileMod
	projectInfoByName[projectFileName] = newInfo

	return newProject
}

func cleanupStaleProjects(s *Server, ctx context.Context) {
	nowMs := time.Now().UnixMilli()
	const ttl = int64(5 * time.Minute / time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	for name, info := range projectInfoByName {
		if nowMs-info.LastAccessed > ttl {
			closeProject(s, name, ctx)
			delete(projectInfoByName, name)
		}
	}
}

func getFileModTimeMs(s *Server, file string) int64 {
	if !s.fs.FileExists(file) {
		return 0
	}
	fi := s.fs.Stat(file)
	if fi == nil {
		return 0
	}
	return fi.ModTime().UnixMilli()
}

func createNewSelfManagedProject(s *Server, projectFileName string, file string, ctx context.Context) *project.Project {
	var p *project.Project
	var err error

	if p, err = s.session.OpenProject(ctx, projectFileName); err == nil && p != nil {
		if p.GetProgram() != nil && p.GetProgram().GetSourceFile(file) != nil {
			return p
		}
	}

	return nil
}

func closeProject(s *Server, projectFileName string, ctx context.Context) {
	err := s.session.CloseProject(ctx, projectFileName)
	if err != nil {
		s.logger.Log("SelfManagedProjects:: Error closing project " + projectFileName + ": " + err.Error())
	}
}
