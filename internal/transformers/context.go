package transformers

import (
	"context"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/printer"
)

type transformContextKey int

const (
	compilerOptionsKeyValue transformContextKey = iota
	emitContextKeyValue
)

func WithCompilerOptions(ctx context.Context, options *core.CompilerOptions) context.Context {
	return context.WithValue(ctx, compilerOptionsKeyValue, options)
}

func GetCompilerOptionsFromContext(ctx context.Context) *core.CompilerOptions {
	return ctx.Value(compilerOptionsKeyValue).(*core.CompilerOptions)
}

func WithEmitContext(ctx context.Context, emitContext *printer.EmitContext) context.Context {
	return context.WithValue(ctx, emitContextKeyValue, emitContext)
}

func GetEmitContextFromContext(ctx context.Context) *printer.EmitContext {
	return ctx.Value(emitContextKeyValue).(*printer.EmitContext)
}
