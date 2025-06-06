package compiler

import (
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
)

type fileLoaderWorkerTask interface {
	comparable
	FileName() string
	start(loader *fileLoader)
}

type fileLoaderWorker[K fileLoaderWorkerTask, V any] struct {
	wg              core.WorkGroup
	tasksByFileName collections.SyncMap[string, K]
	getSubTasks     func(t K) []K
}

func (w *fileLoaderWorker[K, V]) runAndWait(loader *fileLoader, tasks []K) {
	w.start(loader, tasks)
	w.wg.RunAndWait()
}

func (w *fileLoaderWorker[K, V]) start(loader *fileLoader, tasks []K) {
	if len(tasks) > 0 {
		for i, task := range tasks {
			loadedTask, loaded := w.tasksByFileName.LoadOrStore(task.FileName(), task)
			if loaded {
				// dedup tasks to ensure correct file order, regardless of which task would be started first
				tasks[i] = loadedTask
			} else {
				w.wg.Queue(func() {
					task.start(loader)
					subTasks := w.getSubTasks(task)
					w.start(loader, subTasks)
				})
			}
		}
	}
}

func (w *fileLoaderWorker[K, V]) collect(tasks []K, iterate func(K, []V) V) []V {
	return w.collectWorker(tasks, iterate, core.Set[K]{})
}

func (w *fileLoaderWorker[K, V]) collectWorker(tasks []K, iterate func(K, []V) V, seen core.Set[K]) []V {
	var results []V
	for _, task := range tasks {
		// ensure we only walk each task once
		if seen.Has(task) {
			continue
		}
		seen.Add(task)
		var subResults []V
		if subTasks := w.getSubTasks(task); len(subTasks) > 0 {
			subResults = w.collectWorker(subTasks, iterate, seen)
		}
		results = append(results, iterate(task, subResults))
	}
	return results
}
