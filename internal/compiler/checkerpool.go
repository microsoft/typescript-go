package compiler

import (
	"context"
	"slices"
	"sort"
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tracing"
)

// CheckerPool is implemented by the project system to provide checkers with
// request-scoped lifetime and reclamation. It returns a checker and a release
// function that must be called when the caller is done with the checker.
// The returned checker must not be accessed concurrently; each acquisition is exclusive.
// If file is non-nil, the pool may use it as an affinity hint to return the same
// checker for the same file across calls.
type CheckerPool interface {
	GetChecker(ctx context.Context, file *ast.SourceFile) (*checker.Checker, func())
}

type checkerPool struct {
	program *Program
	tracing *tracing.Tracing

	createCheckersOnce sync.Once
	checkers           []*checker.Checker
	locks              []*sync.Mutex
	fileAssociations   map[*ast.SourceFile]*checker.Checker
}

var _ CheckerPool = (*checkerPool)(nil)

const checkerAssociationTextWeightDivisor = 100

// getCheckerAssociationsForFileWeights builds the initial mapping from file index
// to checker index using longest-processing-time-first scheduling. Files with
// the largest estimated checker work are assigned first to the least-loaded
// checker, minimizing the slowest checker bucket before graph refinement nudges
// the mapping toward import locality.
func getCheckerAssociationsForFileWeights(fileWeights []int, checkerCount int) []int {
	if len(fileWeights) == 0 {
		return nil
	}
	associations := make([]int, len(fileWeights))
	checkerWeights := make([]int, checkerCount)
	fileIndices := make([]int, len(fileWeights))
	for i := range fileIndices {
		fileIndices[i] = i
	}
	sort.Slice(fileIndices, func(i, j int) bool {
		left := fileIndices[i]
		right := fileIndices[j]
		if fileWeights[left] != fileWeights[right] {
			return fileWeights[left] > fileWeights[right]
		}
		return left < right
	})
	for _, fileIndex := range fileIndices {
		checkerIndex := 0
		for i, weight := range checkerWeights[1:] {
			if weight < checkerWeights[checkerIndex] {
				checkerIndex = i + 1
			}
		}
		associations[fileIndex] = checkerIndex
		checkerWeights[checkerIndex] += fileWeights[fileIndex]
	}
	return associations
}

// refineCheckerAssociationsByGraph nudges the initial file-index-to-checker-index
// mapping toward the import graph. For each file, it counts which checkers own the
// file's import neighbors and moves the file to the checker with the largest net
// neighbor gain, as long as the move stays within a small load-balance cap. This
// favors sharing cached checker state among related files without letting dense
// import clusters undo the weight balancing from the initial pass.
// This is a deliberately small one-vertex local-search refinement, similar in
// spirit to balanced graph-partitioning heuristics like Kernighan-Lin and
// Fiduccia-Mattheyses, but without their heavier gain queues or swap sequences;
// see https://en.wikipedia.org/wiki/Kernighan%E2%80%93Lin_algorithm and
// https://en.wikipedia.org/wiki/Fiduccia%E2%80%93Mattheyses_algorithm.
func refineCheckerAssociationsByGraph(associations []int, fileWeights []int, adjacentFiles [][]int, checkerCount int) {
	if len(associations) == 0 || checkerCount <= 1 {
		return
	}
	checkerWeights := make([]int, checkerCount)
	totalWeight := 0
	maxFileWeight := 0
	// Reconstruct checker loads from the current mapping, and remember the
	// largest single file because any legal cap must be able to fit it.
	for i, checkerIndex := range associations {
		checkerWeights[checkerIndex] += fileWeights[i]
		totalWeight += fileWeights[i]
		maxFileWeight = max(maxFileWeight, fileWeights[i])
	}
	averageCheckerWeight := (totalWeight + checkerCount - 1) / checkerCount
	maxCheckerWeight := max(maxFileWeight, averageCheckerWeight+averageCheckerWeight/50)
	neighborCounts := make([]int, checkerCount)
	// Make a bounded number of greedy passes. Later moves can create better
	// placements for files already visited, but this should remain a cheap,
	// predictable refinement rather than an expensive graph partitioner.
	for range 2 {
		moved := false
		for fileIndex, currentChecker := range associations {
			// Count how many import neighbors of this file are currently assigned
			// to each checker.
			clear(neighborCounts)
			for _, adjacentFile := range adjacentFiles[fileIndex] {
				neighborCounts[associations[adjacentFile]]++
			}
			bestChecker := currentChecker
			bestGain := 0
			for candidate := range checkerCount {
				if candidate == currentChecker || checkerWeights[candidate]+fileWeights[fileIndex] > maxCheckerWeight {
					continue
				}
				gain := neighborCounts[candidate] - neighborCounts[currentChecker]
				// Prefer a checker that gains more colocated import neighbors. If
				// the gain is tied, prefer the lighter checker so equal-quality
				// graph moves still improve balance.
				if gain > bestGain || gain == bestGain && gain > 0 && checkerWeights[candidate] < checkerWeights[bestChecker] {
					bestChecker = candidate
					bestGain = gain
				}
			}
			if bestChecker != currentChecker {
				associations[fileIndex] = bestChecker
				checkerWeights[currentChecker] -= fileWeights[fileIndex]
				checkerWeights[bestChecker] += fileWeights[fileIndex]
				moved = true
			}
		}
		if !moved {
			break
		}
	}
}

func newCheckerPool(program *Program) *checkerPool {
	return newCheckerPoolWithTracing(program, nil)
}

func newCheckerPoolWithTracing(program *Program, tr *tracing.Tracing) *checkerPool {
	checkerCount := 4
	if program.SingleThreaded() {
		checkerCount = 1
	} else if c := program.Options().Checkers; c != nil {
		checkerCount = *c
	}

	checkerCount = max(min(checkerCount, len(program.files), 256), 1)

	pool := &checkerPool{
		program:  program,
		checkers: make([]*checker.Checker, checkerCount),
		locks:    make([]*sync.Mutex, checkerCount),
		tracing:  tr,
	}

	return pool
}

// GetChecker implements CheckerPool. When file is non-nil, returns the checker
// associated with that file; otherwise returns the first checker.
func (p *checkerPool) GetChecker(ctx context.Context, file *ast.SourceFile) (*checker.Checker, func()) {
	if file != nil {
		return p.getCheckerForFileExclusive(ctx, file)
	}
	p.createCheckers()
	c := p.checkers[0]
	p.locks[0].Lock()
	return c, sync.OnceFunc(func() {
		p.locks[0].Unlock()
	})
}

// getCheckerForFileNonExclusive returns the checker for the given file without locking.
// This is only safe when the caller guarantees no concurrent access to the same checker,
// e.g. for read-only operations like obtaining an emit resolver.
func (p *checkerPool) getCheckerForFileNonExclusive(file *ast.SourceFile) (*checker.Checker, func()) {
	p.createCheckers()
	return p.fileAssociations[file], noop
}

func (p *checkerPool) getCheckerForFileExclusive(ctx context.Context, file *ast.SourceFile) (*checker.Checker, func()) {
	p.createCheckers()
	c := p.fileAssociations[file]
	idx := slices.Index(p.checkers, c)
	p.locks[idx].Lock()
	return c, sync.OnceFunc(func() {
		p.locks[idx].Unlock()
	})
}

// getCheckerNonExclusive returns the first checker without locking.
func (p *checkerPool) getCheckerNonExclusive() (*checker.Checker, func()) {
	p.createCheckers()
	return p.checkers[0], noop
}

func (p *checkerPool) createCheckers() {
	p.createCheckersOnce.Do(func() {
		checkerCount := len(p.checkers)
		wg := core.NewWorkGroup(p.program.SingleThreaded())
		for i := range checkerCount {
			wg.Queue(func() {
				var tracer *checker.Tracer
				if p.tracing != nil {
					tracer = checker.NewTracer(p.tracing, i)
				}
				p.checkers[i], p.locks[i] = checker.NewChecker(p.program, tracer)
			})
		}

		wg.RunAndWait()

		fileWeights := make([]int, len(p.program.files))
		for i, file := range p.program.files {
			fileWeights[i] = max(file.NodeCount+len(file.Text())/checkerAssociationTextWeightDivisor, 1)
		}
		// The association algorithm uses p.program.files indices throughout:
		// start with a work-balanced pass, then refine that mapping using import adjacency.
		associations := getCheckerAssociationsForFileWeights(fileWeights, checkerCount)
		if checkerCount > 1 {
			adjacentFiles := p.getImportAdjacency()
			refineCheckerAssociationsByGraph(associations, fileWeights, adjacentFiles, checkerCount)
		}
		p.fileAssociations = make(map[*ast.SourceFile]*checker.Checker, len(p.program.files))
		for i, file := range p.program.files {
			p.fileAssociations[file] = p.checkers[associations[i]]
		}
	})
}

// getImportAdjacency returns an undirected import graph represented by file
// index. A directed import from A to B makes both files adjacent because either
// file can benefit from sharing checker caches with the other.
func (p *checkerPool) getImportAdjacency() [][]int {
	fileIndices := make(map[*ast.SourceFile]int, len(p.program.files))
	for i, file := range p.program.files {
		fileIndices[file] = i
	}
	adjacentFiles := make([][]int, len(p.program.files))
	for fileIndex, file := range p.program.files {
		resolvedModules := p.program.resolvedModules[file.Path()]
		for _, resolved := range resolvedModules {
			if resolved == nil || !resolved.IsResolved() {
				continue
			}
			importedFile := p.program.GetSourceFileForResolvedModule(resolved.ResolvedFileName)
			importedIndex, ok := fileIndices[importedFile]
			if !ok || importedIndex == fileIndex {
				continue
			}
			adjacentFiles[fileIndex] = append(adjacentFiles[fileIndex], importedIndex)
			adjacentFiles[importedIndex] = append(adjacentFiles[importedIndex], fileIndex)
		}
	}
	return adjacentFiles
}

// Runs `cb` for each checker in the pool concurrently, locking and unlocking checker mutexes as it goes,
// making it safe to call `forEachCheckerParallel` from many threads simultaneously.
func (p *checkerPool) forEachCheckerParallel(cb func(idx int, c *checker.Checker)) {
	p.createCheckers()
	wg := core.NewWorkGroup(p.program.SingleThreaded())
	for idx, checker := range p.checkers {
		wg.Queue(func() {
			p.locks[idx].Lock()
			defer p.locks[idx].Unlock()
			cb(idx, checker)
		})
	}
	wg.RunAndWait()
}

func (p *checkerPool) GetGlobalDiagnostics() []*ast.Diagnostic {
	p.createCheckers()
	globalDiagnostics := make([][]*ast.Diagnostic, len(p.checkers))
	p.forEachCheckerParallel(func(idx int, checker *checker.Checker) {
		globalDiagnostics[idx] = checker.GetGlobalDiagnostics()
	})
	return SortAndDeduplicateDiagnostics(slices.Concat(globalDiagnostics...))
}

// forEachCheckerGroupDo runs one task per checker in parallel. Each task iterates
// the provided files, processing only those assigned to its checker. Within each
// checker's set, files are visited in their original order.
func (p *checkerPool) forEachCheckerGroupDo(ctx context.Context, files []*ast.SourceFile, singleThreaded bool, cb func(c *checker.Checker, fileIndex int, file *ast.SourceFile)) {
	p.createCheckers()

	checkerCount := len(p.checkers)
	wg := core.NewWorkGroup(singleThreaded)
	for checkerIdx := range checkerCount {
		wg.Queue(func() {
			p.locks[checkerIdx].Lock()
			defer p.locks[checkerIdx].Unlock()
			for i, file := range files {
				if checker := p.checkers[checkerIdx]; checker == p.fileAssociations[file] {
					cb(checker, i, file)
				}
			}
		})
	}
	wg.RunAndWait()
}

func noop() {}
