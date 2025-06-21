package project

import (
	"slices"

	"github.com/microsoft/typescript-go/internal/tspath"
)

type Workspace struct {
	root                string
	folders             []string
	comparePathsOptions tspath.ComparePathsOptions
}

func (w *Workspace) GetProjectRootPath(fileName string) string {
	// If dynamic then use root otherwise try to find the best match in folders
	if !isDynamicFileName(fileName) {
		var bestMatch string
		for _, folder := range w.folders {
			if tspath.ContainsPath(folder, fileName, w.comparePathsOptions) &&
				(bestMatch == "" || len(bestMatch) <= len(folder)) {
				bestMatch = folder
			}
		}

		if bestMatch != "" {
			return bestMatch
		}

		if w.root != "" && tspath.ContainsPath(w.root, fileName, w.comparePathsOptions) {
			return w.root
		}

		return ""
	}
	return w.root
}

func (w *Workspace) SetRoot(root string) {
	w.root = root
}

func (w *Workspace) AddFolder(folder string) {
	w.folders = append(w.folders, folder)
}

func (w *Workspace) RemoveFolder(folder string) {
	w.folders = slices.Delete(w.folders, slices.Index(w.folders, folder), slices.Index(w.folders, folder)+1)
}
