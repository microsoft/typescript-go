package project

import (
	"context"
	"iter"
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
)

type CheckerPool struct {
	maxCheckers int
	program     *compiler.Program

	mu                  sync.Mutex
	cond                *sync.Cond
	createCheckersOnce  sync.Once
	checkers            []*checker.Checker
	inUse               map[*checker.Checker]bool
	fileAssociations    map[*ast.SourceFile]*checker.Checker
	requestAssociations map[string]*checker.Checker
}

var _ compiler.CheckerPool = (*CheckerPool)(nil)

func newCheckerPool(maxCheckers int, program *compiler.Program) *CheckerPool {
	pool := &CheckerPool{
		program:             program,
		maxCheckers:         maxCheckers,
		checkers:            make([]*checker.Checker, 0, maxCheckers),
		inUse:               make(map[*checker.Checker]bool),
		requestAssociations: make(map[string]*checker.Checker),
	}

	pool.cond = sync.NewCond(&pool.mu)
	return pool
}

func (p *CheckerPool) GetCheckerForFile(ctx context.Context, file *ast.SourceFile) (*checker.Checker, func()) {
	p.mu.Lock()
	defer p.mu.Unlock()

	requestID := core.GetRequestID(ctx)
	if requestID != "" {
		if checker, ok := p.requestAssociations[requestID]; ok {
			if inUse := p.inUse[checker]; !inUse {
				p.inUse[checker] = true
				return checker, p.createRelease(requestID, checker)
			}
			return checker, noop
		}
	}

	if p.fileAssociations == nil {
		p.fileAssociations = make(map[*ast.SourceFile]*checker.Checker)
	}

	if checker, ok := p.fileAssociations[file]; ok {
		if inUse := p.inUse[checker]; !inUse {
			p.inUse[checker] = true
			if requestID != "" {
				p.requestAssociations[requestID] = checker
			}
			return checker, p.createRelease(requestID, checker)
		}
	}

	checker, release := p.getCheckerLocked(requestID)
	p.fileAssociations[file] = checker
	return checker, release
}

func (p *CheckerPool) GetChecker(ctx context.Context) (*checker.Checker, func()) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.getCheckerLocked(core.GetRequestID(ctx))
}

func (p *CheckerPool) Files(checker *checker.Checker) iter.Seq[*ast.SourceFile] {
	panic("unimplemented")
}

func (p *CheckerPool) GetAllCheckers(ctx context.Context) ([]*checker.Checker, func()) {
	requestID := core.GetRequestID(ctx)
	if requestID == "" {
		panic("cannot call GetAllCheckers on a project.checkerPool without a request ID")
	}

	// A request can only access one checker
	c, release := p.GetChecker(ctx)
	return []*checker.Checker{c}, release
}

func (p *CheckerPool) getCheckerLocked(requestID string) (*checker.Checker, func()) {
	if checker := p.getImmediatelyAvailableChecker(); checker != nil {
		p.inUse[checker] = true
		if requestID != "" {
			p.requestAssociations[requestID] = checker
		}
		return checker, p.createRelease(requestID, checker)
	}

	if len(p.checkers) < p.maxCheckers {
		checker := p.createCheckerLocked()
		p.inUse[checker] = true
		if requestID != "" {
			p.requestAssociations[requestID] = checker
		}
		return checker, p.createRelease(requestID, checker)
	}

	checker := p.waitForAvailableChecker()
	p.inUse[checker] = true
	if requestID != "" {
		p.requestAssociations[requestID] = checker
	}
	return checker, p.createRelease(requestID, checker)
}

func (p *CheckerPool) getImmediatelyAvailableChecker() *checker.Checker {
	if len(p.checkers) == 0 {
		return nil
	}

	for _, checker := range p.checkers {
		if inUse := p.inUse[checker]; !inUse {
			return checker
		}
	}

	return nil
}

func (p *CheckerPool) waitForAvailableChecker() *checker.Checker {
	for {
		p.cond.Wait()
		checker := p.getImmediatelyAvailableChecker()
		if checker != nil {
			return checker
		}
	}
}

func (p *CheckerPool) createRelease(requestId string, checker *checker.Checker) func() {
	return func() {
		p.mu.Lock()
		defer p.mu.Unlock()

		p.inUse[checker] = false
		p.cond.Signal()
	}
}

func (p *CheckerPool) createCheckerLocked() *checker.Checker {
	checker := checker.NewChecker(p.program)
	p.checkers = append(p.checkers, checker)
	return checker
}

func noop() {}
