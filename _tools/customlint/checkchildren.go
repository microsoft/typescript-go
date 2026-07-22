package customlint

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// checkChildrenAnalyzer flags `return` statements in checker methods that are
// dispatched from `checkSourceElementWorker` / `checkExpressionWorker` and that
// exit the method before its child nodes have been checked.
//
// Methods like `checkReturnStatement` contain grammar/error `return`s that fire
// before the child expression is ever checked (e.g. via `checkExpression`).
// When such a return is taken the child is never checked, so the set of
// diagnostics produced depends on whether the child happens to be checked
// elsewhere first. Requiring children to always be checked - even on error
// paths - keeps diagnostics stable regardless of traversal order.
//
// A return is flagged when a child-checking call is reachable *after* it: that
// is what makes the return "early". Checking children up front - even
// conditionally, e.g.
//
//	var exprType *Type
//	if node.Expression() != nil {
//		exprType = c.checkExpressionCached(node.Expression())
//	}
//	if c.grammarError(node) {
//		return
//	}
//
// is therefore accepted, because the child has already been given its chance to
// be checked before the return.
var checkChildrenAnalyzer = &analysis.Analyzer{
	Name: "checkchildren",
	Doc:  "finds early returns in checker dispatch methods that skip checking child nodes",
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	Run: func(pass *analysis.Pass) (any, error) {
		return (&checkChildrenPass{pass: pass}).run()
	},
}

// checkWorkerFuncs are the dispatch methods whose switch cases enumerate the
// per-node-kind checker methods we want to analyze.
var checkWorkerFuncs = map[string]bool{
	"checkSourceElementWorker": true,
	"checkExpressionWorker":    true,
}

type checkChildrenPass struct {
	pass *analysis.Pass
}

func (p *checkChildrenPass) run() (any, error) {
	in := p.pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Collect all methods in the package, keyed by name, and the set of method
	// names dispatched from the worker functions.
	methodsByName := make(map[string]*ast.FuncDecl)
	handlerNames := make(map[string]struct{})

	for cursor := range in.Root().Preorder((*ast.FuncDecl)(nil)) {
		fd := cursor.Node().(*ast.FuncDecl)
		recv, ok := receiverIdent(fd)
		if !ok {
			continue
		}
		methodsByName[fd.Name.Name] = fd
		if checkWorkerFuncs[fd.Name.Name] {
			collectDispatchedMethods(fd, recv, handlerNames)
		}
	}

	for name := range handlerNames {
		fd := methodsByName[name]
		if fd == nil || fd.Body == nil {
			continue
		}
		p.analyzeMethod(fd)
	}

	return nil, nil
}

// collectDispatchedMethods records every `recv.<method>(...)` call appearing in
// the body of a worker function.
func collectDispatchedMethods(fd *ast.FuncDecl, recv string, out map[string]struct{}) {
	ast.Inspect(fd.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if name, ok := receiverCallName(call, recv); ok {
			out[name] = struct{}{}
		}
		return true
	})
}

func (p *checkChildrenPass) analyzeMethod(fd *ast.FuncDecl) {
	recv, ok := receiverIdent(fd)
	if !ok {
		return
	}

	// Collect the child-checking calls and the explicit returns that belong to
	// this method. Nested function literals are skipped: they have their own
	// control flow and returns, which are not this method's concern.
	var checkCalls []*ast.CallExpr
	var returns []*ast.ReturnStmt
	ast.Inspect(fd.Body, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.ReturnStmt:
			returns = append(returns, n)
		case *ast.CallExpr:
			if name, isCall := receiverCallName(n, recv); isCall && isChildCheckName(name) {
				checkCalls = append(checkCalls, n)
			}
		}
		return true
	})

	// If the method never checks any children, there is nothing to enforce:
	// either the node kind is a leaf, or checking is delegated elsewhere.
	if len(checkCalls) == 0 {
		return
	}

	parent := buildParentMap(fd)
	for _, ret := range returns {
		for _, call := range reachableChecksAfter(ret, checkCalls, parent) {
			// A return that only fires when the child is absent (a `nil` guard on
			// the very expression the later call would check) is fine: there is
			// nothing to check on that path.
			if child := firstArg(call); child != nil && guardedNil(ret, child, parent) {
				continue
			}
			p.pass.Report(analysis.Diagnostic{
				Pos:     ret.Pos(),
				End:     ret.End(),
				Message: fd.Name.Name + " returns before checking its child nodes; check children (e.g. via checkExpression, checkSourceElement, or resolveCall) on all paths so diagnostics are stable regardless of traversal order",
			})
			break
		}
	}
}

// reachableChecksAfter returns the child-checking calls that would execute after
// `start` if it fell through instead of returning. Walking up the ancestor
// chain, only statements that genuinely follow `start` in execution order are
// considered; mutually exclusive branches (the other arm of an if, sibling
// switch cases) are not, while loop bodies are revisited in full.
func reachableChecksAfter(start ast.Node, checkCalls []*ast.CallExpr, parent map[ast.Node]ast.Node) []*ast.CallExpr {
	var reachable []*ast.CallExpr
	child := start
	for {
		switch p := parent[child].(type) {
		case nil:
			return reachable
		case *ast.FuncDecl, *ast.FuncLit:
			// Reached the enclosing function without finding a later check.
			return reachable
		case *ast.BlockStmt:
			// Sibling switch/select clauses are alternatives, not successors.
			if !isSwitchClause(child) {
				reachable = appendChecksAfter(reachable, p.List, child, checkCalls)
			}
			child = p
		case *ast.CaseClause:
			reachable = appendChecksAfter(reachable, p.Body, child, checkCalls)
			child = p
		case *ast.CommClause:
			reachable = appendChecksAfter(reachable, p.Body, child, checkCalls)
			child = p
		case *ast.ForStmt:
			// Falling through the loop body re-runs the whole loop.
			if p.Body == child {
				reachable = appendChecksIn(reachable, p.Body, checkCalls)
				if p.Post != nil {
					reachable = appendChecksIn(reachable, p.Post, checkCalls)
				}
				if p.Cond != nil {
					reachable = appendChecksIn(reachable, p.Cond, checkCalls)
				}
			}
			child = p
		case *ast.RangeStmt:
			if p.Body == child {
				reachable = appendChecksIn(reachable, p.Body, checkCalls)
			}
			child = p
		default:
			child = p
		}
	}
}

// appendChecksAfter appends the child-checks contained in the statements that
// follow `child` in `list`.
func appendChecksAfter[T ast.Stmt](reachable []*ast.CallExpr, list []T, child ast.Node, checkCalls []*ast.CallExpr) []*ast.CallExpr {
	idx := -1
	for i, s := range list {
		if ast.Node(s) == child {
			idx = i
			break
		}
	}
	if idx < 0 {
		return reachable
	}
	for _, s := range list[idx+1:] {
		reachable = appendChecksIn(reachable, s, checkCalls)
	}
	return reachable
}

// appendChecksIn appends the child-checks whose call lies within the source
// range of n.
func appendChecksIn(reachable []*ast.CallExpr, n ast.Node, checkCalls []*ast.CallExpr) []*ast.CallExpr {
	for _, call := range checkCalls {
		if n.Pos() <= call.Pos() && call.Pos() < n.End() {
			reachable = append(reachable, call)
		}
	}
	return reachable
}

func isSwitchClause(n ast.Node) bool {
	switch n.(type) {
	case *ast.CaseClause, *ast.CommClause:
		return true
	}
	return false
}

// firstArg returns the first argument of a call, or nil if it has none.
func firstArg(call *ast.CallExpr) ast.Expr {
	if len(call.Args) == 0 {
		return nil
	}
	return call.Args[0]
}

// guardedNil reports whether every path from an enclosing `if` guard to `ret`
// establishes that `childExpr` is nil, i.e. the return only fires when the
// child that a later call would check is absent.
func guardedNil(ret ast.Node, childExpr ast.Expr, parent map[ast.Node]ast.Node) bool {
	child := ret
	for {
		switch p := parent[child].(type) {
		case nil, *ast.FuncDecl, *ast.FuncLit:
			return false
		case *ast.IfStmt:
			if p.Body == child && condImpliesNil(p.Cond, childExpr, true) {
				return true
			}
			if p.Else == child && condImpliesNil(p.Cond, childExpr, false) {
				return true
			}
			child = p
		default:
			child = p
		}
	}
}

// condImpliesNil reports whether `cond` evaluating to `condIsTrue` implies that
// `childExpr` is nil, i.e. `childExpr == nil` (when true) or `childExpr != nil`
// (when false).
func condImpliesNil(cond ast.Expr, childExpr ast.Expr, condIsTrue bool) bool {
	bin, ok := cond.(*ast.BinaryExpr)
	if !ok {
		return false
	}
	wantOp := token.EQL
	if !condIsTrue {
		wantOp = token.NEQ
	}
	if bin.Op != wantOp {
		return false
	}
	return (isNilIdent(bin.Y) && equalExpr(bin.X, childExpr)) ||
		(isNilIdent(bin.X) && equalExpr(bin.Y, childExpr))
}

func isNilIdent(e ast.Expr) bool {
	id, ok := e.(*ast.Ident)
	return ok && id.Name == "nil"
}

// equalExpr reports whether two expressions are structurally identical for the
// simple forms used to identify child nodes (identifiers, selectors, and calls
// such as `node.Expression()`).
func equalExpr(a, b ast.Expr) bool {
	switch a := a.(type) {
	case *ast.Ident:
		b, ok := b.(*ast.Ident)
		return ok && a.Name == b.Name
	case *ast.SelectorExpr:
		b, ok := b.(*ast.SelectorExpr)
		return ok && a.Sel.Name == b.Sel.Name && equalExpr(a.X, b.X)
	case *ast.CallExpr:
		b, ok := b.(*ast.CallExpr)
		if !ok || len(a.Args) != len(b.Args) || !equalExpr(a.Fun, b.Fun) {
			return false
		}
		for i := range a.Args {
			if !equalExpr(a.Args[i], b.Args[i]) {
				return false
			}
		}
		return true
	case *ast.ParenExpr:
		b, ok := b.(*ast.ParenExpr)
		return ok && equalExpr(a.X, b.X)
	case *ast.IndexExpr:
		b, ok := b.(*ast.IndexExpr)
		return ok && equalExpr(a.X, b.X) && equalExpr(a.Index, b.Index)
	case *ast.BasicLit:
		b, ok := b.(*ast.BasicLit)
		return ok && a.Kind == b.Kind && a.Value == b.Value
	default:
		return false
	}
}

// buildParentMap maps each node in fd to its parent node.
func buildParentMap(fd *ast.FuncDecl) map[ast.Node]ast.Node {
	parent := make(map[ast.Node]ast.Node)
	stack := make([]ast.Node, 0, 32)
	ast.Inspect(fd, func(n ast.Node) bool {
		if n == nil {
			stack = stack[:len(stack)-1]
			return true
		}
		if len(stack) > 0 {
			parent[n] = stack[len(stack)-1]
		}
		stack = append(stack, n)
		return true
	})
	return parent
}

// isChildCheckName reports whether a method name recursively checks a child node.
func isChildCheckName(name string) bool {
	return strings.HasPrefix(name, "checkExpression") ||
		strings.HasPrefix(name, "checkSourceElement") ||
		name == "resolveCall"
}

// receiverCallName returns the method name of a `recv.<method>(...)` call.
func receiverCallName(call *ast.CallExpr, recv string) (string, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}
	x, ok := sel.X.(*ast.Ident)
	if !ok || x.Name != recv {
		return "", false
	}
	return sel.Sel.Name, true
}

// receiverIdent returns the name of a method's single receiver variable.
func receiverIdent(fd *ast.FuncDecl) (string, bool) {
	if fd.Recv == nil || len(fd.Recv.List) != 1 {
		return "", false
	}
	names := fd.Recv.List[0].Names
	if len(names) != 1 {
		return "", false
	}
	return names[0].Name, true
}
