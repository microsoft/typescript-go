package compiler

import (
	"context"
	"slices"
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

// Process small contiguous blocks of files on the same checker before rotating
// to the next checker. Adjacent files tend to share generic instantiations and
// symbol/type links; assigning blocks to the least-loaded checker by text size
// preserves that locality while keeping checker work balanced.
const maxCheckerAssociationBlockSize = 32

func getCheckerAssociationBlockSize(fileCount int, checkerCount int) int {
	const targetBlocksPerChecker = 4
	if checkerCount <= 1 {
		return maxCheckerAssociationBlockSize
	}
	return min(max(fileCount/(checkerCount*targetBlocksPerChecker), 1), maxCheckerAssociationBlockSize)
}

func getCheckerAssociationsForFileWeights(fileWeights []int, checkerCount int) []int {
	if len(fileWeights) == 0 {
		return nil
	}
	blockSize := getCheckerAssociationBlockSize(len(fileWeights), checkerCount)
	associations := make([]int, len(fileWeights))
	checkerWeights := make([]int, checkerCount)
	for blockStart := 0; blockStart < len(fileWeights); blockStart += blockSize {
		checkerIndex := 0
		for i, weight := range checkerWeights[1:] {
			if weight < checkerWeights[checkerIndex] {
				checkerIndex = i + 1
			}
		}
		blockEnd := min(blockStart+blockSize, len(fileWeights))
		for i := blockStart; i < blockEnd; i++ {
			associations[i] = checkerIndex
			checkerWeights[checkerIndex] += fileWeights[i]
		}
	}
	return associations
}

func refineCheckerAssociationsByGraph(associations []int, fileWeights []int, adjacentFiles [][]int, checkerCount int) {
	if len(associations) == 0 || checkerCount <= 1 {
		return
	}
	checkerWeights := make([]int, checkerCount)
	totalWeight := 0
	maxFileWeight := 0
	for i, checkerIndex := range associations {
		checkerWeights[checkerIndex] += fileWeights[i]
		totalWeight += fileWeights[i]
		maxFileWeight = max(maxFileWeight, fileWeights[i])
	}
	averageCheckerWeight := (totalWeight + checkerCount - 1) / checkerCount
	maxCheckerWeight := max(maxFileWeight, averageCheckerWeight+averageCheckerWeight/50)
	neighborCounts := make([]int, checkerCount)
	for range 2 {
		moved := false
		for fileIndex, currentChecker := range associations {
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
			fileWeights[i] = len(file.Text()) + 3*file.NodeCount + 90*file.SymbolCount
		}
		associations := getCheckerAssociationsForFileWeights(fileWeights, checkerCount)
		adjacentFiles := p.getImportAdjacency()
		refineCheckerAssociationsByGraph(associations, fileWeights, adjacentFiles, checkerCount)
		p.fileAssociations = make(map[*ast.SourceFile]*checker.Checker, len(p.program.files))
		for i, file := range p.program.files {
			p.fileAssociations[file] = p.checkers[associations[i]]
		}
	})
}

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
