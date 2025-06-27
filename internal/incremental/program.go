package incremental

import (
	"context"
	"fmt"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
)

type Program struct {
	state   *programState
	program *compiler.Program
}

var _ compiler.AnyProgram = (*Program)(nil)

func NewProgram(program *compiler.Program, oldProgram *Program) *Program {
	return &Program{
		state:   newProgramState(program, oldProgram),
		program: program,
	}
}

func (p *Program) panicIfNoProgram(method string) {
	if p.program == nil {
		panic(fmt.Sprintf("%s should not be called without program", method))
	}
}

func (p *Program) GetProgram() *compiler.Program {
	p.panicIfNoProgram("GetProgram")
	return p.program
}

func (p *Program) Options() *core.CompilerOptions {
	return p.state.options
}

func (p *Program) GetSourceFiles() []*ast.SourceFile {
	p.panicIfNoProgram("GetSourceFiles")
	return p.program.GetSourceFiles()
}

func (p *Program) GetConfigFileParsingDiagnostics() []*ast.Diagnostic {
	p.panicIfNoProgram("GetConfigFileParsingDiagnostics")
	return p.program.GetConfigFileParsingDiagnostics()
}

func (p *Program) GetSyntacticDiagnostics(ctx context.Context, file *ast.SourceFile) []*ast.Diagnostic {
	p.panicIfNoProgram("GetSyntacticDiagnostics")
	return p.program.GetSyntacticDiagnostics(ctx, file)
}

func (p *Program) GetBindDiagnostics(ctx context.Context, file *ast.SourceFile) []*ast.Diagnostic {
	p.panicIfNoProgram("GetBindDiagnostics")
	return p.program.GetBindDiagnostics(ctx, file)
}

func (p *Program) GetOptionsDiagnostics(ctx context.Context) []*ast.Diagnostic {
	p.panicIfNoProgram("GetOptionsDiagnostics")
	return p.program.GetOptionsDiagnostics(ctx)
}

func (p *Program) GetGlobalDiagnostics(ctx context.Context) []*ast.Diagnostic {
	p.panicIfNoProgram("GetGlobalDiagnostics")
	return p.program.GetGlobalDiagnostics(ctx)
}

func (p *Program) GetSemanticDiagnostics(ctx context.Context, file *ast.SourceFile) []*ast.Diagnostic {
	p.panicIfNoProgram("GetSemanticDiagnostics")
	return p.state.getSemanticDiagnostics(ctx, p.program, file)
}

func (p *Program) GetDeclarationDiagnostics(ctx context.Context, file *ast.SourceFile) []*ast.Diagnostic {
	p.panicIfNoProgram("GetDeclarationDiagnostics")
	return p.state.getDeclarationDiagnostics(ctx, p.program, file)
}

func (p *Program) Emit(ctx context.Context, options compiler.EmitOptions) *compiler.EmitResult {
	p.panicIfNoProgram("Emit")
	return p.state.emit(ctx, p.program, options)
}
