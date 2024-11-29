package gcilint

import (
	"go/ast"
	"go/token"
	"slices"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("emptycase", New)
}

type emptyCasePlugin struct{}

func New(settings any) (register.LinterPlugin, error) {
	return &emptyCasePlugin{}, nil
}

func (f *emptyCasePlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name: "emptycase",
			Doc:  "finds empty switch/select cases",
			Run:  f.run,
		},
	}, nil
}

func (f *emptyCasePlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

func (f *emptyCasePlugin) run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.SwitchStmt:
				checkCases(pass, file, n.Body)
			case *ast.SelectStmt:
				checkCases(pass, file, n.Body)
			}
			return true
		})
	}

	return nil, nil
}

func checkCases(pass *analysis.Pass, file *ast.File, clause *ast.BlockStmt) {
	endOfBlock := clause.End()

	for i, stmt := range clause.List {
		nextCasePos := endOfBlock
		if j := i + 1; j < len(clause.List) {
			nextCasePos = clause.List[j].Pos()
		}
		checkCaseStatement(pass, file, stmt, nextCasePos)
	}
}

func checkCaseStatement(pass *analysis.Pass, file *ast.File, stmt ast.Stmt, nextCasePos token.Pos) {
	var body []ast.Stmt
	colon := token.NoPos

	switch stmt := stmt.(type) {
	case *ast.CaseClause:
		body = stmt.Body
		colon = stmt.Colon
	case *ast.CommClause:
		body = stmt.Body
		colon = stmt.Colon
	default:
		return
	}

	if len(body) == 1 {
		// Also error on a case statement containing a single empty block.
		block, ok := body[0].(*ast.BlockStmt)
		if !ok || len(block.List) != 0 {
			return
		}
	} else if len(body) != 0 {
		return
	}

	if _, found := slices.BinarySearchFunc(file.Comments, posRange{colon, nextCasePos}, posRangeCmp); found {
		return
	}

	pass.Report(analysis.Diagnostic{
		Pos:     stmt.Pos(),
		End:     colon,
		Message: "this case block is empty and will do nothing",
	})
}

type posRange struct {
	start, end token.Pos
}

func posRangeCmp(c *ast.CommentGroup, target posRange) int {
	if c.End() < target.start {
		return -1
	}
	if c.Pos() >= target.end {
		return 1
	}
	return 0
}
