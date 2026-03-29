package inliners

import (
	"strings"
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/debug"
	"github.com/microsoft/typescript-go/internal/jsnum"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/scanner"
	"github.com/microsoft/typescript-go/internal/transformers"
)

type ConstEnumInliningTransformer struct {
	transformers.Transformer
	compilerOptions   *core.CompilerOptions
	currentSourceFile *ast.SourceFile
	emitResolver      printer.EmitResolver
}

var constEnumPool = sync.Pool{New: func() any { return &ConstEnumInliningTransformer{} }}

func getConstEnumInliningTransformer() *ConstEnumInliningTransformer {
	return constEnumPool.Get().(*ConstEnumInliningTransformer)
}

func putConstEnumInliningTransformer(tx *ConstEnumInliningTransformer) {
	dispose, visit := tx.SaveState()
	*tx = ConstEnumInliningTransformer{}
	tx.RestoreState(dispose, visit)
	constEnumPool.Put(tx)
}

func NewConstEnumInliningTransformer(opt *transformers.TransformOptions) *transformers.Transformer {
	compilerOptions := opt.CompilerOptions
	emitContext := opt.Context
	if compilerOptions.GetIsolatedModules() {
		debug.Fail("const enums are not inlined under isolated modules")
	}
	tx := getConstEnumInliningTransformer()
	tx.compilerOptions = compilerOptions
	tx.emitResolver = opt.EmitResolver
	if tx.GetDispose() == nil {
		tx.SetDispose(func() { putConstEnumInliningTransformer(tx) })
	}
	return tx.NewTransformer(tx.visit, emitContext)
}

func (tx *ConstEnumInliningTransformer) visit(node *ast.Node) *ast.Node {
	switch node.Kind {
	case ast.KindPropertyAccessExpression, ast.KindElementAccessExpression:
		{
			parse := tx.EmitContext().ParseNode(node)
			if parse == nil {
				return tx.Visitor().VisitEachChild(node)
			}
			value := tx.emitResolver.GetConstantValue(parse)
			if value != nil {
				var replacement *ast.Node
				switch v := value.(type) {
				case jsnum.Number:
					if v.IsInf() {
						if v.Abs() == v {
							replacement = tx.Factory().NewIdentifier("Infinity")
						} else {
							replacement = tx.Factory().NewPrefixUnaryExpression(ast.KindMinusToken, tx.Factory().NewIdentifier("Infinity"))
						}
					} else if v.IsNaN() {
						replacement = tx.Factory().NewIdentifier("NaN")
					} else if v.Abs() == v {
						replacement = tx.Factory().NewNumericLiteral(v.String(), ast.TokenFlagsNone)
					} else {
						replacement = tx.Factory().NewPrefixUnaryExpression(ast.KindMinusToken, tx.Factory().NewNumericLiteral(v.Abs().String(), ast.TokenFlagsNone))
					}
				case string:
					replacement = tx.Factory().NewStringLiteral(v, ast.TokenFlagsNone)
				case jsnum.PseudoBigInt: // technically not supported by strada, and issues a checker error, handled here for completeness
					if v == (jsnum.PseudoBigInt{}) {
						replacement = tx.Factory().NewBigIntLiteral("0", ast.TokenFlagsNone)
					} else if !v.Negative {
						replacement = tx.Factory().NewBigIntLiteral(v.Base10Value, ast.TokenFlagsNone)
					} else {
						replacement = tx.Factory().NewPrefixUnaryExpression(ast.KindMinusToken, tx.Factory().NewBigIntLiteral(v.Base10Value, ast.TokenFlagsNone))
					}
				}

				if tx.compilerOptions.RemoveComments.IsFalseOrUnknown() {
					original := tx.EmitContext().MostOriginal(node)
					if original != nil && !ast.NodeIsSynthesized(original) {
						originalText := scanner.GetTextOfNode(original)
						escapedText := safeMultiLineComment(originalText)
						tx.EmitContext().AddSyntheticTrailingComment(replacement, ast.KindMultiLineCommentTrivia, escapedText, false)
					}
				}
				return replacement
			}
			return tx.Visitor().VisitEachChild(node)
		}
	}
	return tx.Visitor().VisitEachChild(node)
}

func safeMultiLineComment(text string) string {
	var b strings.Builder
	b.Grow(len(text) + 2)
	b.WriteByte(' ')
	for {
		i := strings.Index(text, "*/")
		if i < 0 {
			break
		}
		b.WriteString(text[:i])
		b.WriteString("*_/")
		text = text[i+2:]
	}
	b.WriteString(text)
	b.WriteByte(' ')
	return b.String()
}
