package customlint

import (
	"cmp"
	"go/ast"
	"go/token"
	"go/types"
	"slices"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/cfg"
)

var shadowAnalyzer = &analysis.Analyzer{
	Name:     "shadow",
	Doc:      "check for possible unintended shadowing of variables",
	URL:      "https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/shadow",
	Requires: []*analysis.Analyzer{inspect.Analyzer, ctrlflow.Analyzer},
	Run: func(pass *analysis.Pass) (any, error) {
		return (&shadowPass{
			pass:    pass,
			inspect: pass.ResultOf[inspect.Analyzer].(*inspector.Inspector),
			cfgs:    pass.ResultOf[ctrlflow.Analyzer].(*ctrlflow.CFGs),
		}).run()
	},
}

type shadowPass struct {
	pass    *analysis.Pass
	inspect *inspector.Inspector
	cfgs    *ctrlflow.CFGs

	objectDefs map[types.Object]*ast.Ident
	objectUses map[types.Object][]*ast.Ident
	scopes     map[*types.Scope]ast.Node

	funcDecl *ast.FuncDecl
	funcLits []*ast.FuncLit
}

func (s *shadowPass) run() (any, error) {
	s.objectDefs = make(map[types.Object]*ast.Ident)
	for id, obj := range s.pass.TypesInfo.Defs {
		if obj != nil {
			s.objectDefs[obj] = id
		}
	}

	s.objectUses = make(map[types.Object][]*ast.Ident)
	for id, obj := range s.pass.TypesInfo.Uses {
		if obj != nil {
			s.objectUses[obj] = append(s.objectUses[obj], id)
		}
	}
	for _, uses := range s.objectUses {
		slices.SortFunc(uses, comparePos)
	}

	s.scopes = make(map[*types.Scope]ast.Node, len(s.pass.TypesInfo.Scopes))
	for id, scope := range s.pass.TypesInfo.Scopes {
		s.scopes[scope] = id
	}

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
		(*ast.AssignStmt)(nil),
		(*ast.GenDecl)(nil),
	}
	s.inspect.Nodes(nodeFilter, func(n ast.Node, push bool) (proceed bool) {
		if !push {
			switch n.(type) {
			case *ast.FuncDecl:
				s.funcDecl = nil
			case *ast.FuncLit:
				s.funcLits = s.funcLits[:len(s.funcLits)-1]
			}
			return true
		}

		// var c *cfg.CFG
		switch n := n.(type) {
		case *ast.AssignStmt:
			s.handleAssignment(n)
			return true
		case *ast.GenDecl:
			s.handleAssignment(n)
			return true

		case *ast.FuncDecl:
			s.funcDecl = n
			// c = s.cfgs.FuncDecl(n)
		case *ast.FuncLit:
			s.funcLits = append(s.funcLits, n)
			// c = s.cfgs.FuncLit(n)
		}

		// fmt.Println(n, c.Format(s.pass.Fset))

		return true
	})

	return nil, nil
}

func (s *shadowPass) handleAssignment(n ast.Node) {
	var idents []*ast.Ident

	switch n := n.(type) {
	case *ast.AssignStmt:
		if n.Tok != token.DEFINE {
			return
		}
		for _, expr := range n.Lhs {
			ident, ok := expr.(*ast.Ident)
			if !ok {
				continue
			}
			idents = append(idents, ident)
		}
	case *ast.GenDecl:
		if n.Tok != token.VAR {
			return
		}
		for _, spec := range n.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, ident := range valueSpec.Names {
				idents = append(idents, ident)
			}
		}
	}

	for _, ident := range idents {
		if ident.Name == "_" {
			// Can't shadow the blank identifier.
			continue
		}
		obj := s.pass.TypesInfo.Defs[ident]
		if obj == nil {
			continue
		}
		// obj.Parent.Parent is the surrounding scope. If we can find another declaration
		// starting from there, we have a shadowed identifier.
		_, shadowed := obj.Parent().Parent().LookupParent(obj.Name(), obj.Pos())
		if shadowed == nil {
			continue
		}
		shadowedScope := shadowed.Parent()
		// Don't complain if it's shadowing a universe-declared identifier; that's fine.
		if shadowedScope == types.Universe {
			continue
		}
		// Ignore shadowing a type name, which can never result in a logic error.
		if isTypeName(obj) || isTypeName(shadowed) {
			continue
		}
		// Don't complain if the types differ: that implies the programmer really wants two different things.
		if !types.Identical(obj.Type(), shadowed.Type()) {
			continue
		}
		// Always ban shadowing something outside a function (package-lock declarations).
		// TODO(jakebailey): fix this; catches if statements etc
		// if _, ok := s.scopes[shadowedScope].(*ast.FuncType); !ok {
		// 	s.report(ident, shadowed)
		// 	continue
		// }

		uses := s.objectUses[obj]
		var lastUse *ast.Ident
		if len(uses) > 0 {
			lastUse = uses[len(uses)-1]
		}
		if lastUse == nil {
			// Unused variable.
			continue
		}

		shadowUses := s.objectUses[shadowed]
		idx, _ := slices.BinarySearchFunc(shadowUses, lastUse, comparePos)
		if idx == len(shadowUses) {
			// Shadowed variable is never used after shadowing.
			continue
		}

		nextShadowUse := shadowUses[idx]

		fnScope := s.currentFunctionScope()

		if nextShadowUse.End() > fnScope.End() {
			// line := s.pass.Fset.Position(nextShadowUse.Pos()).Line
			// fmt.Println("next shadow use of", ident.Name, "on line", line, "is outside the current function")
			// continue
			// TODO(jakebailey): either these should always report because
			// their meaning is unclear, or we need to walk up to find the
			// function literal in the shadowed scope and then do CFG reachability
			// from there.
			s.report(ident, shadowed)
			continue
		}

		if s.positionIsReachable(ident, nextShadowUse.Pos()) {
			s.report(ident, shadowed)
		}
	}
}

func (s *shadowPass) report(ident *ast.Ident, shadowed types.Object) {
	line := s.pass.Fset.Position(shadowed.Pos()).Line
	s.pass.ReportRangef(ident, "declaration of %q shadows declaration at line %d", ident.Name, line)
}

func (s *shadowPass) currentFunctionScope() *types.Scope {
	if len(s.funcLits) > 0 {
		last := s.funcLits[len(s.funcLits)-1]
		return s.pass.TypesInfo.Scopes[last.Type]
	}
	return s.pass.TypesInfo.Scopes[s.funcDecl.Type]
}

func (s *shadowPass) currentCFG() *cfg.CFG {
	if len(s.funcLits) > 0 {
		last := s.funcLits[len(s.funcLits)-1]
		return s.cfgs.FuncLit(last)
	}
	return s.cfgs.FuncDecl(s.funcDecl)
}

func (s *shadowPass) positionIsReachable(ident *ast.Ident, pos token.Pos) bool {
	c := s.currentCFG()

	var start *cfg.Block
	for _, b := range c.Blocks {
		if posInBlock(b, ident.Pos()) {
			start = b
			break
		}
	}
	if start == nil {
		return false // TODO: error
	}

	seen := make(map[*cfg.Block]struct{})
	var posReachable func(b *cfg.Block) (found bool)
	posReachable = func(b *cfg.Block) (found bool) {
		if _, ok := seen[b]; ok {
			return false
		}
		seen[b] = struct{}{}

		if posInBlock(b, pos) {
			return true
		}

		return slices.ContainsFunc(b.Succs, posReachable)
	}

	return posReachable(start)
}

func posInBlock(b *cfg.Block, pos token.Pos) bool {
	if len(b.Nodes) == 0 {
		return false
	}

	first := b.Nodes[0]
	last := b.Nodes[len(b.Nodes)-1]

	return first.Pos() <= pos && pos <= last.End()
}

func comparePos[T ast.Node](a, b T) int {
	return cmp.Compare(a.Pos(), b.Pos())
}

func isTypeName(obj types.Object) bool {
	_, ok := obj.(*types.TypeName)
	return ok
}
