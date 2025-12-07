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

type CheckerPool interface {
	GetChecker(ctx context.Context) (*checker.Checker, func())
	GetCheckerForFile(ctx context.Context, file *ast.SourceFile) (*checker.Checker, func())
	GetCheckerForFileExclusive(ctx context.Context, file *ast.SourceFile) (*checker.Checker, func())
}

type checkerPool struct {
	checkerCount int
	program      *Program
	tracing      *tracing.Tracing

	createCheckersOnce sync.Once
	checkers           []*checker.Checker
	locks              []*sync.Mutex
	fileAssociations   map[*ast.SourceFile]*checker.Checker
}

var _ CheckerPool = (*checkerPool)(nil)

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
		program:      program,
		checkerCount: checkerCount,
		checkers:     make([]*checker.Checker, checkerCount),
		locks:        make([]*sync.Mutex, checkerCount),
		tracing:      tr,
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
				var tracer checker.TypeTracer
				if p.tracing != nil {
					tracer = checker.NewTracingTypeTracer(p.tracing.NewTypeTracer(i))
				}
				p.checkers[i], p.locks[i] = checker.NewCheckerWithTracer(p.program, tracer)
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

func noop() {}
