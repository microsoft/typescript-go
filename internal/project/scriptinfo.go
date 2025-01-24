package project

import (
	"slices"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type scriptInfo struct {
	fileName   string
	path       tspath.Path
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
			// !!!
			// s.ensureRealPath()
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
