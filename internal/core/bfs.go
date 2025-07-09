package core

import (
	"math"
	"sync"
	"sync/atomic"

	"github.com/microsoft/typescript-go/internal/collections"
)

type BreadthFirstSearchResult[N comparable] struct {
	Stopped bool
	Path    []N
}

// BreadthFirstSearchParallel performs a breadth-first search on a graph
// starting from the given node. It processes nodes in parallel and returns the path
// from the first node that satisfies the `visit` function back to the start node.
func BreadthFirstSearchParallel[N comparable](
	start N,
	neighbors func(N) []N,
	visit func(node N) (isResult bool, stop bool),
	visited *collections.SyncSet[N],
) BreadthFirstSearchResult[N] {
	if visited == nil {
		visited = &collections.SyncSet[N]{}
	}

	type job struct {
		node   N
		parent *job
	}

	type result struct {
		stop bool
		job  *job
		next *collections.OrderedMap[N, *job]
	}

	var fallback *job
	// processLevel processes each node at the current level in parallel.
	// It produces either a list of jobs to be processed in the next level,
	// or a result if the visit function returns true for any node.
	processLevel := func(index int, jobs *collections.OrderedMap[N, *job]) result {
		var lowestFallback atomic.Int64
		var lowestGoal atomic.Int64
		var nextJobCount atomic.Int64
		lowestGoal.Store(math.MaxInt64)
		lowestFallback.Store(math.MaxInt64)
		next := make([][]*job, jobs.Size())
		var wg sync.WaitGroup
		i := 0
		for j := range jobs.Values() {
			wg.Add(1)
			go func(i int, j *job) {
				defer wg.Done()
				if int64(i) >= lowestGoal.Load() {
					return // Stop processing if we already found a lower result
				}

				// If we have already visited this node, skip it.
				if !visited.AddIfAbsent(j.node) {
					// Note that if we are here, we already visited this node at a
					// previous *level*, which means `visit` must have returned false,
					// so we don't need to update our result indices. This holds true
					// because we deduplicated jobs before queuing the level.
					return
				}

				isResult, stop := visit(j.node)
				if isResult {
					// We found a result, so we will stop at this level, but an
					// earlier job may still find a true result at a lower index.
					if stop {
						updateMin(&lowestGoal, int64(i))
						return
					}
					if fallback == nil {
						updateMin(&lowestFallback, int64(i))
					}
				}

				if int64(i) >= lowestGoal.Load() {
					// If `visit` is expensive, it's likely that by the time we get here,
					// a different job has already found a lower index result, so we
					// don't even need to collect the next jobs.
					return
				}
				// Add the next level jobs
				neighborNodes := neighbors(j.node)
				if len(neighborNodes) > 0 {
					nextJobCount.Add(int64(len(neighborNodes)))
					next[i] = Map(neighborNodes, func(child N) *job {
						return &job{node: child, parent: j}
					})
				}
			}(i, j)
			i++
		}
		wg.Wait()
		if index := lowestGoal.Load(); index != math.MaxInt64 {
			// If we found a result, return it immediately.
			_, job, _ := jobs.EntryAt(int(index))
			return result{stop: true, job: job}
		}
		if fallback == nil {
			if index := lowestFallback.Load(); index != math.MaxInt64 {
				_, fallback, _ = jobs.EntryAt(int(index))
			}
		}
		nextJobs := collections.NewOrderedMapWithSizeHint[N, *job](int(nextJobCount.Load()))
		for _, jobs := range next {
			for _, j := range jobs {
				if !nextJobs.Has(j.node) {
					// Deduplicate synchronously to avoid messy locks and spawning
					// unnecessary goroutines.
					nextJobs.Set(j.node, j)
				}
			}
		}
		return result{next: nextJobs}
	}

	createPath := func(job *job) []N {
		var path []N
		for job != nil {
			path = append(path, job.node)
			job = job.parent
		}
		return path
	}

	levelIndex := 0
	level := collections.NewOrderedMapFromList([]collections.MapEntry[N, *job]{
		{Key: start, Value: &job{node: start}},
	})
	for level.Size() > 0 {
		result := processLevel(levelIndex, level)
		if result.stop {
			return BreadthFirstSearchResult[N]{Stopped: true, Path: createPath(result.job)}
		} else if result.job != nil && fallback == nil {
			fallback = result.job
		}
		level = result.next
		levelIndex++
	}
	return BreadthFirstSearchResult[N]{Stopped: false, Path: createPath(fallback)}
}

// updateMin updates the atomic integer `a` to the candidate value if it is less than the current value.
func updateMin(a *atomic.Int64, candidate int64) bool {
	for {
		current := a.Load()
		if current < candidate {
			return false
		}
		if a.CompareAndSwap(current, candidate) {
			return true
		}
	}
}
