package customlint

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var rangeSymbolTableAnalyzer = &analysis.Analyzer{
	Name: "rangesymboltable",
	Doc:  "finds range statements over ast.SymbolTable, which has nondeterministic iteration order",
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	Run: func(pass *analysis.Pass) (any, error) {
		return (&rangeSymbolTablePass{pass: pass}).run()
	},
}

type rangeSymbolTablePass struct {
	pass *analysis.Pass
}

func (r *rangeSymbolTablePass) run() (any, error) {
	in := r.pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	for c := range in.Root().Preorder((*ast.RangeStmt)(nil)) {
		rangeStmt := c.Node().(*ast.RangeStmt)
		if r.isSymbolTable(r.pass.TypesInfo.TypeOf(rangeStmt.X)) {
			r.pass.Report(analysis.Diagnostic{
				Pos:     rangeStmt.Pos(),
				End:     rangeStmt.X.End(),
				Message: "range over ast.SymbolTable has nondeterministic iteration order",
			})
		}
	}

	return nil, nil
}

// symbolTablePkgPaths contains the real package path and a fake one used for testing,
// since Go's internal package restriction prevents testdata from importing the real path.
var symbolTablePkgPaths = []string{
	"github.com/microsoft/typescript-go/internal/ast",
	"testdata/fakeast",
}

func (r *rangeSymbolTablePass) isSymbolTable(t types.Type) bool {
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	if obj.Name() != "SymbolTable" || obj.Pkg() == nil {
		return false
	}
	for _, path := range symbolTablePkgPaths {
		if obj.Pkg().Path() == path {
			return true
		}
	}
	return false
}
