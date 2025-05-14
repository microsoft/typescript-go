package project

import (
	"fmt"
	"slices"

	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

const (
	fileGlobPattern          = "*.{js,jsx,mjs,cjs,ts,tsx,mts,cts,json}"
	recursiveFileGlobPattern = "**/*.{js,jsx,mjs,cjs,ts,tsx,mts,cts,json}"
)

type watchedFiles[T any] struct {
	p         *Project
	getGlobs  func(data T) []string
	watchKind lsproto.WatchKind

	data      T
	globs     []string
	watcherID WatcherHandle
	watchType string
}

func newWatchedFiles[T any](p *Project, watchKind lsproto.WatchKind, getGlobs func(data T) []string, watchType string) *watchedFiles[T] {
	return &watchedFiles[T]{
		p:         p,
		watchKind: watchKind,
		getGlobs:  getGlobs,
		watchType: watchType,
	}
}

func (w *watchedFiles[T]) update(newData T) {
	if updated, err := w.updateWorker(newData); err != nil {
		w.p.Log(fmt.Sprintf("Failed to update %s watch: %v\n%s", w.watchType, err, formatFileList(w.globs, "\t", hr)))
	} else if updated {
		w.p.Logf("%s watches updated %s:\n%s", w.watchType, w.watcherID, formatFileList(w.globs, "\t", hr))
	}
}

func (w *watchedFiles[T]) updateWorker(newData T) (updated bool, err error) {
	newGlobs := w.getGlobs(newData)
	w.data = newData
	if slices.Equal(w.globs, newGlobs) {
		return false, nil
	}

	w.globs = newGlobs
	if w.watcherID != "" {
		if err = w.p.host.Client().UnwatchFiles(w.watcherID); err != nil {
			return false, err
		}
	}

	w.watcherID = ""
	if len(newGlobs) == 0 {
		return true, nil
	}

	watchers := make([]*lsproto.FileSystemWatcher, 0, len(newGlobs))
	for _, glob := range newGlobs {
		watchers = append(watchers, &lsproto.FileSystemWatcher{
			GlobPattern: lsproto.PatternOrRelativePattern{
				Pattern: &glob,
			},
			Kind: &w.watchKind,
		})
	}
	watcherID, err := w.p.host.Client().WatchFiles(watchers)
	if err != nil {
		return false, err
	}
	w.watcherID = watcherID
	return true, nil
}
