package customlint

import (
    "go/ast"
    "go/token"
    "slices"

    "golang.org/x/tools/go/analysis"
    "golang.org/x/tools/go/analysis/passes/inspect"
    "golang.org/x/tools/go/ast/inspector"
)

var emptyCaseAnalyzer = &analysis.Analyzer{
    Name: "emptycase",
    Doc:  "finds empty switch/select cases",
    Run:  runEmptyCase,
    Requires: []*analysis.Analyzer{
        inspect.Analyzer,
    },
}

func runEmptyCase(pass *analysis.Pass) (any, error) {
    inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

    nodeFilter := []ast.Node{
        (*ast.File)(nil),
        (*ast.SwitchStmt)(nil),
        (*ast.SelectStmt)(nil),
    }

    var file *ast.File // Track the current file being inspected

    inspect.Preorder(nodeFilter, func(n ast.Node) {
        switch n := n.(type) {
        case *ast.File:
            file = n
        case *ast.SwitchStmt, *ast.SelectStmt:
            checkCases(pass, file, n.(*ast.BlockStmt)) // Combined SwitchStmt and SelectStmt handling
        }
    })

    return nil, nil
}

func checkCases(pass *analysis.Pass, file *ast.File, clause *ast.BlockStmt) {
    for i, stmt := range clause.List {
        nextCasePos := clause.End()
        if next := i + 1; next < len(clause.List) {
            nextCasePos = clause.List[next].Pos()
        }
        checkCaseStatement(pass, file, stmt, nextCasePos)
    }
}

func checkCaseStatement(pass *analysis.Pass, file *ast.File, stmt ast.Stmt, nextCasePos token.Pos) {
    var body []ast.Stmt
    var colon token.Pos

    switch stmt := stmt.(type) {
    case *ast.CaseClause, *ast.CommClause:
        body = stmt.(*ast.BlockStmt).List
        colon = stmt.(*ast.BlockStmt).Lbrace // Extract colon more directly
    }

    if len(body) == 1 && isEmptyBlock(body[0]) {
        return
    } else if len(body) != 0 {
        return
    }

    afterColon := colon + 1
    if _, found := slices.BinarySearchFunc(file.Comments, posRange{afterColon, nextCasePos}, posRangeCmp); found {
        return
    }

    pass.Report(analysis.Diagnostic{
        Pos:     stmt.Pos(),
        End:     afterColon,
        Message: "this case block is empty and will do nothing",
    })
}

func isEmptyBlock(stmt ast.Stmt) bool {
    block, ok := stmt.(*ast.BlockStmt)
    return ok && len(block.List) == 0
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


// Changes made: 
// 1. Combined SwitchStmt and SelectStmt handling into a single case for cleaner traversal logic. 
// 2. Moved block emptiness check to a helper function isEmptyBlock for better readability and reuse. 
// 3. Simplified checkCaseStatement by removing redundant logic and directly accessing colon. 
// 4. Improved clarity and reduced duplication in code traversal logic.
