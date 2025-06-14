package compiler

import (
	"math"
	"sync"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type fileLoaderWorkerTask[T any] interface {
	comparable
	FileName() string
	run(loader *fileLoader)
	getSubTasks() []T
	shouldIncreaseDepth() bool
}

type fileLoaderWorker[K fileLoaderWorkerTask[K]] struct {
	wg              core.WorkGroup
	tasksByFileName collections.SyncMap[string, *queuedTask[K]]
	maxDepth        int
}

type queuedTask[K fileLoaderWorkerTask[K]] struct {
	task        K
	mu          sync.Mutex
	run         bool
	lowestDepth int
}

func (w *fileLoaderWorker[K]) runAndWait(loader *fileLoader, tasks []K) {
	w.start(loader, tasks, 0)
	w.wg.RunAndWait()
}

func (w *fileLoaderWorker[K]) start(loader *fileLoader, tasks []K, depth int) {
	for i, task := range tasks {
		newTask := &queuedTask[K]{task: task, lowestDepth: math.MaxInt}
		loadedTask, loaded := w.tasksByFileName.LoadOrStore(task.FileName(), newTask)
		task = loadedTask.task
		if loaded {
			tasks[i] = task
		}

		nextDepth := depth
		if task.shouldIncreaseDepth() {
			nextDepth++
		}

		if nextDepth > w.maxDepth {
			continue
		}

		w.wg.Queue(func() {
			loadedTask.mu.Lock()
			defer loadedTask.mu.Unlock()

			if !loadedTask.run {
				task.run(loader)
				loadedTask.run = true
			}

			if nextDepth < loadedTask.lowestDepth {
				loadedTask.lowestDepth = nextDepth
				subTasks := task.getSubTasks()
				w.start(loader, subTasks, nextDepth)
			}
		})
	}
}

func (w *fileLoaderWorker[K]) collect(loader *fileLoader, tasks []K, iterate func(K, []tspath.Path)) []tspath.Path {
	return w.collectWorker(loader, tasks, iterate, collections.Set[K]{})
}

func (w *fileLoaderWorker[K]) collectWorker(loader *fileLoader, tasks []K, iterate func(K, []tspath.Path), seen collections.Set[K]) []tspath.Path {
	var results []tspath.Path
	for _, task := range tasks {
		// ensure we only walk each task once
		if seen.Has(task) {
			continue
		}
		seen.Add(task)
		// TODO(jakebailey): skip unrun tasks?
		var subResults []tspath.Path
		if subTasks := task.getSubTasks(); len(subTasks) > 0 {
			subResults = w.collectWorker(loader, subTasks, iterate, seen)
		}
		iterate(task, subResults)
		results = append(results, loader.toPath(task.FileName()))
	}
	return results
}
