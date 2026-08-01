package compiler

import (
	"context"
	"math"
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

// The checker work proxy cannot predict demand-driven semantic work: checking a
// small root can populate a large amount of dependency state. These safety factors
// were selected from cross-project sweeps on VS Code, TypeScript, MUI docs, and
// XState at 2, 4, and 8 checkers. They are deliberately project-independent.
//
// Source-dominated projects at any checker count balance implementation roots
// directly, using unmodified file weights and a 12x balance penalty. Other projects
// use program order; at four or more checkers they use a 4x implementation-file
// base weight and a 16x balance penalty, while smaller checker pools use the
// published weighted FENNEL penalty without scaling.
const (
	checkerAssociationSourceFileWeightMultiplier   = 4
	checkerAssociationBalancePenaltyMultiplier     = 16
	checkerAssociationPrioritizedSourcePenalty     = 12
	checkerAssociationStrongBalanceMinCheckerCount = 4
)

type checkerAssociationPolicy struct {
	prioritizeSourceFiles      bool
	sourceFileWeightMultiplier int
	balancePenaltyMultiplier   int
}

func getCheckerAssociationPolicy(totalWeight int, declarationWeight int, checkerCount int) checkerAssociationPolicy {
	if shouldPrioritizeSourceFiles(totalWeight, declarationWeight, checkerCount) {
		return checkerAssociationPolicy{
			prioritizeSourceFiles:      true,
			sourceFileWeightMultiplier: 1,
			balancePenaltyMultiplier:   checkerAssociationPrioritizedSourcePenalty,
		}
	}
	if checkerCount >= checkerAssociationStrongBalanceMinCheckerCount {
		return checkerAssociationPolicy{
			sourceFileWeightMultiplier: checkerAssociationSourceFileWeightMultiplier,
			balancePenaltyMultiplier:   checkerAssociationBalancePenaltyMultiplier,
		}
	}
	return checkerAssociationPolicy{
		sourceFileWeightMultiplier: 1,
		balancePenaltyMultiplier:   1,
	}
}

// getCheckerAssociationsInOrder partitions the import graph using a weighted adaptation
// of FENNEL's streaming graph-partitioning objective with gamma = 3/2. Each file
// is placed where it has the most already-placed neighbors, minus the incremental
// convex load penalty. The published alpha = m*sqrt(k)/n^(3/2) becomes
// m*sqrt(k)/W^(3/2), where W is total estimated checker work.
//
// A nil order means stable program order. The preferred maximum checker weight is
// the larger of the largest file and 101% of average. If no checker can accept a
// file under that bound, the file is assigned to the least-loaded checker. Ties
// are deterministic.
func getCheckerAssociationsInOrder(fileWeights []int, adjacentFiles [][]int, fileOrder []int, checkerCount int, penaltyMultiplier int) []int {
	if len(fileWeights) == 0 {
		return nil
	}

	totalWeight := 0
	maxFileWeight := 0
	edgeCount := 0
	for i, weight := range fileWeights {
		totalWeight += weight
		maxFileWeight = max(maxFileWeight, weight)
		edgeCount += len(adjacentFiles[i])
	}

	associations := make([]int, len(fileWeights))
	for i := range associations {
		associations[i] = -1
	}
	checkerWeights := make([]int, checkerCount)
	averageCheckerWeight := (totalWeight + checkerCount - 1) / checkerCount
	maxCheckerWeight := max(maxFileWeight, averageCheckerWeight+averageCheckerWeight/100)
	totalWeightFloat := float64(totalWeight)
	alpha := float64(penaltyMultiplier) * float64(edgeCount/2) * math.Sqrt(float64(checkerCount)) / (totalWeightFloat * math.Sqrt(totalWeightFloat))
	neighborCounts := make([]int, checkerCount)

	for position := range fileWeights {
		fileIndex := position
		if fileOrder != nil {
			fileIndex = fileOrder[position]
		}

		clear(neighborCounts)
		for _, adjacentFile := range adjacentFiles[fileIndex] {
			if checkerIndex := associations[adjacentFile]; checkerIndex >= 0 {
				neighborCounts[checkerIndex]++
			}
		}

		bestChecker := -1
		bestScore := math.Inf(-1)
		for checkerIndex, checkerWeight := range checkerWeights {
			if checkerWeight+fileWeights[fileIndex] > maxCheckerWeight {
				continue
			}
			oldWeight := float64(checkerWeight)
			newWeight := float64(checkerWeight + fileWeights[fileIndex])
			penalty := alpha * (newWeight*math.Sqrt(newWeight) - oldWeight*math.Sqrt(oldWeight))
			score := float64(neighborCounts[checkerIndex]) - penalty
			if score > bestScore || score == bestScore && (bestChecker < 0 || checkerWeight < checkerWeights[bestChecker]) {
				bestChecker = checkerIndex
				bestScore = score
			}
		}
		if bestChecker < 0 {
			bestChecker = 0
			for checkerIndex, checkerWeight := range checkerWeights[1:] {
				if checkerWeight < checkerWeights[bestChecker] {
					bestChecker = checkerIndex + 1
				}
			}
		}
		associations[fileIndex] = bestChecker
		checkerWeights[bestChecker] += fileWeights[fileIndex]
	}
	return associations
}

// getCheckerAssociationOrder places implementation files before declarations and
// orders each group by descending estimated work. Returning nil preserves program
// order without allocating an index array.
func getCheckerAssociationOrder(fileWeights []int, isDeclarationFile []bool, prioritizeSourceFiles bool) []int {
	if !prioritizeSourceFiles {
		return nil
	}
	fileOrder := make([]int, len(fileWeights))
	for i := range fileOrder {
		fileOrder[i] = i
	}
	sort.Slice(fileOrder, func(i, j int) bool {
		left := fileOrder[i]
		right := fileOrder[j]
		if isDeclarationFile[left] != isDeclarationFile[right] {
			return !isDeclarationFile[left]
		}
		if fileWeights[left] != fileWeights[right] {
			return fileWeights[left] > fileWeights[right]
		}
		return left < right
	})
	return fileOrder
}

func getCheckerAssociationBaseWeight(nodeCount int, textLength int) int {
	return max(nodeCount+textLength/checkerAssociationTextWeightDivisor, 1)
}

// shouldPrioritizeSourceFiles reports whether declaration-file base work is at
// most half of one average checker load. Cross-project sweeps found this to be the
// stable boundary where source-first ordering improved root balance without losing
// the declaration locality needed by declaration-heavy projects.
func shouldPrioritizeSourceFiles(totalWeight int, declarationWeight int, checkerCount int) bool {
	return declarationWeight*checkerCount*2 <= totalWeight
}

// getCheckerAssociationWeights combines local estimated work with dependency
// fanout. The import unit is normalized so that total import weight equals total
// base weight for the project, avoiding a project-specific tuning constant.
func getCheckerAssociationWeights(baseWeights []int, importCounts []int) []int {
	totalBaseWeight := 0
	totalImports := 0
	for i, baseWeight := range baseWeights {
		totalBaseWeight += baseWeight
		totalImports += importCounts[i]
	}
	importWeight := 0
	if totalImports > 0 {
		importWeight = max(totalBaseWeight/totalImports, 1)
	}
	fileWeights := make([]int, len(baseWeights))
	for i, baseWeight := range baseWeights {
		fileWeights[i] = baseWeight + importCounts[i]*importWeight
	}
	return fileWeights
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

		associations := make([]int, len(p.program.files))
		if checkerCount > 1 {
			baseWeights := make([]int, len(p.program.files))
			importCounts := make([]int, len(p.program.files))
			isDeclarationFile := make([]bool, len(p.program.files))
			totalBaseWeight := 0
			declarationBaseWeight := 0
			for i, file := range p.program.files {
				baseWeight := getCheckerAssociationBaseWeight(file.NodeCount, len(file.Text()))
				totalBaseWeight += baseWeight
				if file.IsDeclarationFile {
					declarationBaseWeight += baseWeight
				}
				baseWeights[i] = baseWeight
				importCounts[i] = len(file.Imports())
				isDeclarationFile[i] = file.IsDeclarationFile
			}
			// Rebenchmark the vscode, self-compiler, mui-docs, and xstate-main
			// TypeScript-benchmarking scenarios at 2, 4, and 8 checkers before
			// changing this policy or its constants.
			policy := getCheckerAssociationPolicy(totalBaseWeight, declarationBaseWeight, checkerCount)
			if policy.sourceFileWeightMultiplier != 1 {
				// Apply this before import normalization: increasing total base
				// weight also increases the project-normalized cost of every import.
				for i, declaration := range isDeclarationFile {
					if !declaration {
						baseWeights[i] *= policy.sourceFileWeightMultiplier
					}
				}
			}
			fileWeights := getCheckerAssociationWeights(baseWeights, importCounts)
			adjacentFiles := p.getImportAdjacency()
			fileOrder := getCheckerAssociationOrder(fileWeights, isDeclarationFile, policy.prioritizeSourceFiles)
			associations = getCheckerAssociationsInOrder(fileWeights, adjacentFiles, fileOrder, checkerCount, policy.balancePenaltyMultiplier)
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
