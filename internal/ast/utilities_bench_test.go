package ast_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/testutil/fixtures"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

func collectFunctionBodies(node *ast.Node, out *[]*ast.Node) {
	node.ForEachChild(func(child *ast.Node) bool {
		if ast.IsFunctionLikeDeclaration(child) {
			if body := child.Body(); body != nil && body.Kind == ast.KindBlock {
				*out = append(*out, body)
			}
		}
		collectFunctionBodies(child, out)
		return false
	})
}

func returnStatementVisitor(*ast.Node) bool { return false }

var forEachReturnStatementSink bool

func BenchmarkForEachReturnStatement(b *testing.B) {
	for _, f := range fixtures.BenchFixtures {
		b.Run(f.Name(), func(b *testing.B) {
			f.SkipIfNotExist(b)

			fileName := tspath.GetNormalizedAbsolutePath(f.Path(), "/")
			path := tspath.ToPath(fileName, "/", osvfs.FS().UseCaseSensitiveFileNames())
			sourceText := f.ReadFile(b)
			scriptKind := core.GetScriptKindFromFileName(fileName)

			sf := parser.ParseSourceFile(ast.SourceFileParseOptions{
				FileName: fileName,
				Path:     path,
			}, sourceText, scriptKind)

			var bodies []*ast.Node
			collectFunctionBodies(sf.AsNode(), &bodies)
			if len(bodies) == 0 {
				b.Skip("no function bodies")
			}

			b.ReportAllocs()
			for b.Loop() {
				for _, body := range bodies {
					forEachReturnStatementSink = ast.ForEachReturnStatement(body, returnStatementVisitor)
				}
			}
		})
	}
}
