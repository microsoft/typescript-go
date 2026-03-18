package customlint

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var bitclearAnalyzer = &analysis.Analyzer{
	Name: "bitclear",
	Doc:  "finds `x &= ^y` and suggests `x &^= y` instead",
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	Run: func(pass *analysis.Pass) (any, error) {
		return (&bitclearPass{pass: pass}).run()
	},
}

type bitclearPass struct {
	pass *analysis.Pass
}

func (b *bitclearPass) run() (any, error) {
	in := b.pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	for c := range in.Root().Preorder((*ast.AssignStmt)(nil)) {
		stmt := c.Node().(*ast.AssignStmt)
		if stmt.Tok != token.AND_ASSIGN {
			continue
		}
		if len(stmt.Rhs) != 1 {
			continue
		}

		rhs := stmt.Rhs[0]
		unary, ok := rhs.(*ast.UnaryExpr)
		if !ok || unary.Op != token.XOR {
			continue
		}

		// Found `x &= ^expr`. Report and suggest `x &^= expr`.
		// If the operand is parenthesized (e.g. `^(A | B)`), strip the parens
		// since they are unnecessary in `x &^= A | B`.

		// Build text edits:
		// 1. Replace `&=` with `&^=`
		// 2. Remove `^` (and parens if present)
		edits := []analysis.TextEdit{
			{Pos: stmt.TokPos, End: stmt.TokPos + token.Pos(len("&=")), NewText: []byte("&^=")},
		}

		if paren, ok := unary.X.(*ast.ParenExpr); ok {
			// Remove `^(` and trailing `)`
			edits = append(edits,
				analysis.TextEdit{Pos: unary.Pos(), End: paren.X.Pos(), NewText: nil},
				analysis.TextEdit{Pos: paren.X.End(), End: paren.End(), NewText: nil},
			)
		} else {
			// Remove `^`
			edits = append(edits,
				analysis.TextEdit{Pos: unary.Pos(), End: unary.X.Pos(), NewText: nil},
			)
		}

		b.pass.Report(analysis.Diagnostic{
			Pos:     stmt.TokPos,
			End:     unary.X.End(),
			Message: "use `&^=` instead of `&= ^`",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message:   "Replace with `&^=`",
					TextEdits: edits,
				},
			},
		})
	}

	return nil, nil
}
