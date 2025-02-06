package customlint

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var shadow = &analysis.Analyzer{
	Name: "shadow",
	Doc:  "bans shadowing",
	Run:  runShadow,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func runShadow(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
		(*ast.AssignStmt)(nil),
		(*ast.CaseClause)(nil),
		(*ast.IfStmt)(nil),
		(*ast.ForStmt)(nil),
		(*ast.BlockStmt)(nil),
	}

	// The inspect package doesn't tell us up front which file is being used,
	// so keep track of it as part of the traversal. The file is the first node
	// so will be set before any other nodes are visited.
	var file *ast.File
	scopes := []map[string]bool{}

	inspect.Nodes(nodeFilter, func(n ast.Node, push bool) bool {
		if !push {
			switch n.(type) {
			case *ast.FuncDecl, *ast.CaseClause, *ast.FuncLit, *ast.File, *ast.BlockStmt, *ast.IfStmt, *ast.ForStmt:
				scopes = scopes[:len(scopes)-1]
			}
			return false
		}
		switch n := n.(type) {
		case *ast.File:
			file = n
			scopes = append(scopes, make(map[string]bool))
		case *ast.FuncDecl, *ast.CaseClause, *ast.FuncLit, *ast.BlockStmt, *ast.IfStmt, *ast.ForStmt:
			scopes = append(scopes, make(map[string]bool))
		case *ast.AssignStmt:
			if n.Tok == token.DEFINE {
				checkShadowingAssignment(pass, file, n, scopes)
			}
		}
		return true
	})

	return nil, nil
}

func checkShadowingAssignment(pass *analysis.Pass, file *ast.File, stmt *ast.AssignStmt, scopes []map[string]bool) {
	for _, name := range stmt.Lhs {
		switch name := name.(type) {
		case *ast.Ident:
			if name.Name != "_" && name.Name != "err" && name.Name != "ok" && name.Name != "ch" && name.Name != "size" && resolve(name.Name, scopes) {
				pass.Report(analysis.Diagnostic{
					Pos:     stmt.Pos(),
					End:     stmt.End(),
					Message: fmt.Sprintf("this declaration shadows a previous declaration of %s", name.Name),
				})
			}
			scopes[len(scopes)-1][name.Name] = true
		default:
			continue
		}
	}
}

func resolve(name string, scopes []map[string]bool) bool {
	for _, scope := range scopes {
		if _, ok := scope[name]; ok {
			return true
		}
	}
	return false
}
