package compiler

import (
	"context"
	"slices"
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/core"
)

type CheckerPool interface {
	GetChecker(ctx context.Context) (*checker.Checker, func())
	GetCheckerForFile(ctx context.Context, file *ast.SourceFile) (*checker.Checker, func())
	GetCheckerForFileExclusive(ctx context.Context, file *ast.SourceFile) (*checker.Checker, func())
	GetGlobalDiagnostics() []*ast.Diagnostic
}

type checkerPool struct {
	program *Program

	createCheckersOnce sync.Once
	checkers           []*checker.Checker
	locks              []*sync.Mutex
	fileAssociations   map[*ast.SourceFile]*checker.Checker
}

var _ CheckerPool = (*checkerPool)(nil)

func newCheckerPool(program *Program) *checkerPool {
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
	}

	return pool
}

func (p *checkerPool) GetCheckerForFile(ctx context.Context, file *ast.SourceFile) (*checker.Checker, func()) {
	p.createCheckers()
	return p.fileAssociations[file], noop
}

func (p *checkerPool) GetCheckerForFileExclusive(ctx context.Context, file *ast.SourceFile) (*checker.Checker, func()) {
	p.createCheckers()
	c := p.fileAssociations[file]
	idx := slices.Index(p.checkers, c)
	p.locks[idx].Lock()
	return c, sync.OnceFunc(func() {
		p.locks[idx].Unlock()
	})
}

func (p *checkerPool) GetChecker(ctx context.Context) (*checker.Checker, func()) {
	p.createCheckers()
	checker := p.checkers[0]
	return checker, noop
}

func (p *checkerPool) createCheckers() {
	p.createCheckersOnce.Do(func() {
		checkerCount := len(p.checkers)
		wg := core.NewWorkGroup(p.program.SingleThreaded())
		for i := range checkerCount {
			wg.Queue(func() {
				p.checkers[i], p.locks[i] = checker.NewChecker(p.program)
			})
		}

		wg.RunAndWait()

		p.fileAssociations = make(map[*ast.SourceFile]*checker.Checker, len(p.program.files))
		for i, file := range p.program.files {
			p.fileAssociations[file] = p.checkers[i%checkerCount]
		}
	})
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

// forEachCheckerGroupDo groups the provided files by their associated checker and
// processes each group in parallel, one task per checker. Within each group, files
// are processed sequentially in their original order relative to the input slice.
// The callback receives the checker (held exclusively) along with the file index and file.
func (p *checkerPool) forEachCheckerGroupDo(ctx context.Context, files []*ast.SourceFile, singleThreaded bool, cb func(c *checker.Checker, fileIndex int, file *ast.SourceFile)) {
	p.createCheckers()

	checkerCount := len(p.checkers)
	// Build reverse map from checker pointer to index for efficient grouping.
	checkerIndices := make(map[*checker.Checker]int, checkerCount)
	for i, c := range p.checkers {
		checkerIndices[c] = i
	}

	// Group file indices by their associated checker, preserving relative order.
	groups := make([][]int, checkerCount)
	for i, file := range files {
		c := p.fileAssociations[file]
		idx := checkerIndices[c]
		groups[idx] = append(groups[idx], i)
	}

	// Process each checker's files in parallel, one task per checker.
	wg := core.NewWorkGroup(singleThreaded)
	for checkerIdx := range checkerCount {
		if len(groups[checkerIdx]) == 0 {
			continue
		}
		wg.Queue(func() {
			p.locks[checkerIdx].Lock()
			defer p.locks[checkerIdx].Unlock()
			for _, fileIdx := range groups[checkerIdx] {
				cb(p.checkers[checkerIdx], fileIdx, files[fileIdx])
			}
		})
	}
	wg.RunAndWait()
}

func noop() {}
