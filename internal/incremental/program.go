package incremental

import (
	"context"

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

func (p *Program) GetProgram() *compiler.Program {
	if p.program == nil {
		panic("GetProgram should not be called without program")
	}
	return p.program
}

func (p *Program) Options() *core.CompilerOptions {
	return p.state.options
}

func (p *Program) GetSourceFiles() []*ast.SourceFile {
	if p.program == nil {
		panic("GetSourceFiles should not be called without program")
	}
	return p.program.GetSourceFiles()
}

func (p *Program) GetConfigFileParsingDiagnostics() []*ast.Diagnostic {
	if p.program == nil {
		panic("GetConfigFileParsingDiagnostics should not be called without program")
	}
	return p.program.GetConfigFileParsingDiagnostics()
}

func (p *Program) GetSyntacticDiagnostics(ctx context.Context, file *ast.SourceFile) []*ast.Diagnostic {
	if p.program == nil {
		panic("GetSyntacticDiagnostics should not be called without program")
	}
	return p.program.GetSyntacticDiagnostics(ctx, file)
}

func (p *Program) GetBindDiagnostics(ctx context.Context, file *ast.SourceFile) []*ast.Diagnostic {
	if p.program == nil {
		panic("GetBindDiagnostics should not be called without program")
	}
	return p.program.GetBindDiagnostics(ctx, file)
}

func (p *Program) GetOptionsDiagnostics(ctx context.Context) []*ast.Diagnostic {
	if p.program == nil {
		panic("GetOptionsDiagnostics should not be called without program")
	}
	return p.program.GetOptionsDiagnostics(ctx)
}

func (p *Program) GetGlobalDiagnostics(ctx context.Context) []*ast.Diagnostic {
	if p.program == nil {
		panic("GetGlobalDiagnostics should not be called without program")
	}
	return p.program.GetGlobalDiagnostics(ctx)
}

func (p *Program) GetSemanticDiagnostics(ctx context.Context, file *ast.SourceFile) []*ast.Diagnostic {
	if p.program == nil {
		panic("GetSemanticDiagnostics should not be called without program")
	}
	return p.state.getSemanticDiagnostics(ctx, p.program, file)
}

func (p *Program) GetDeclarationDiagnostics(ctx context.Context, file *ast.SourceFile) []*ast.Diagnostic {
	if p.program == nil {
		panic("GetDeclarationDiagnostics should not be called without program")
	}
	return p.state.getDeclarationDiagnostics(ctx, p.program, file)
}

func (p *Program) Emit(ctx context.Context, options compiler.EmitOptions) *compiler.EmitResult {
	if p.program == nil {
		panic("Emit should not be called without program")
	}
	return p.state.emit(ctx, p.program, options)
}
