package project

import (
	"slices"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type scriptInfo struct {
	fileName   string
	path       tspath.Path
	realpath   tspath.Path
	isDynamic  bool
	scriptKind core.ScriptKind
	text       string
	version    int

	isOpen                bool
	pendingReloadFromDisk bool
	matchesDiskText       bool
	deferredDelete        bool

	containingProjects []*Project
}

func newScriptInfo(fileName string, path tspath.Path, scriptKind core.ScriptKind) *scriptInfo {
	return &scriptInfo{
		fileName:   fileName,
		path:       path,
		scriptKind: scriptKind,
	}
}

func (s *scriptInfo) open(newText string) {
	s.isOpen = true
	s.pendingReloadFromDisk = false
	if newText != s.text {
		s.setText(newText)
		s.matchesDiskText = false
		s.markContainingProjectsAsDirty()
	}
}

func (s *scriptInfo) close(fileExists bool) {
	s.isOpen = false
	if fileExists && !s.pendingReloadFromDisk && !s.matchesDiskText {
		s.pendingReloadFromDisk = true
		s.markContainingProjectsAsDirty()
	}
}

func (s *scriptInfo) setText(newText string) {
	s.text = newText
	s.version++
}

func (s *scriptInfo) markContainingProjectsAsDirty() {
	for _, project := range s.containingProjects {
		project.markFileAsDirty(s.path)
	}
}

// attachToProject attaches the script info to the project if it's not already attached
// and returns true if the script info was newly attached.
func (s *scriptInfo) attachToProject(project *Project) bool {
	if !s.isAttached(project) {
		s.containingProjects = append(s.containingProjects, project)
		if project.compilerOptions.PreserveSymlinks != core.TSTrue {
			s.ensureRealpath(project.FS())
		}
		project.onFileAddedOrRemoved(s.isSymlink())
		return true
	}
	return false
}

func (s *scriptInfo) isAttached(project *Project) bool {
	return slices.Contains(s.containingProjects, project)
}

func (s *scriptInfo) isSymlink() bool {
	// !!!
	return false
}

func (s *scriptInfo) isOrphan() bool {
	if s.deferredDelete {
		return true
	}
	for _, project := range s.containingProjects {
		if !project.isOrphan() {
			return false
		}
	}
	return true
}

func (s *scriptInfo) editContent(change ls.TextChange) {
	s.setText(change.ApplyTo(s.text))
	s.markContainingProjectsAsDirty()
}

func (s *scriptInfo) ensureRealpath(fs vfs.FS) {
	if s.realpath == "" {
		if len(s.containingProjects) == 0 {
			panic("scriptInfo must be attached to a project before calling ensureRealpath")
		}
		realpath := fs.Realpath(string(s.path))
		project := s.containingProjects[0]
		s.realpath = project.toPath(realpath)
		if s.realpath != s.path {
			project.projectService.recordSymlink(s)
		}
	}
}

func (s *scriptInfo) getRealpathIfDifferent() (tspath.Path, bool) {
	if s.realpath != "" && s.realpath != s.path {
		return s.realpath, true
	}
	return "", false
}

func (s *scriptInfo) detachAllProjects() {
	for _, project := range s.containingProjects {
		// !!!
		// if (isConfiguredProject(p)) {
		// 	p.getCachedDirectoryStructureHost().addOrDeleteFile(this.fileName, this.path, FileWatcherEventKind.Deleted);
		// }
		isRoot := project.isRoot(s)
		project.removeFile(s, false /*fileExists*/, false /*detachFromProject*/)
		project.onFileAddedOrRemoved(s.isSymlink())
		if isRoot && project.kind != ProjectKindInferred {
			project.addMissingRootFile(s.fileName, s.path)
		}
	}
	s.containingProjects = nil
}

func (s *scriptInfo) detachFromProject(project *Project) {
	if index := slices.Index(s.containingProjects, project); index != -1 {
		s.containingProjects[index].onFileAddedOrRemoved(s.isSymlink())
		s.containingProjects = slices.Delete(s.containingProjects, index, index+1)
	}
}

func (s *scriptInfo) delayReloadNonMixedContentFile() {
	// !!!
	s.pendingReloadFromDisk = true
	s.markContainingProjectsAsDirty()
}
