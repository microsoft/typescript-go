package compiler

import (
	"context"
	"slices"
	"time"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
)

type AnyProgram interface {
	Options() *core.CompilerOptions
	GetSourceFiles() []*ast.SourceFile
	GetConfigFileParsingDiagnostics() []*ast.Diagnostic
	GetSyntacticDiagnostics(ctx context.Context, file *ast.SourceFile) []*ast.Diagnostic
	GetBindDiagnostics(ctx context.Context, file *ast.SourceFile) []*ast.Diagnostic
	GetOptionsDiagnostics(ctx context.Context) []*ast.Diagnostic
	GetGlobalDiagnostics(ctx context.Context) []*ast.Diagnostic
	GetSemanticDiagnostics(ctx context.Context, file *ast.SourceFile) []*ast.Diagnostic
	GetDeclarationDiagnostics(ctx context.Context, file *ast.SourceFile) []*ast.Diagnostic
	Emit(ctx context.Context, options EmitOptions) *EmitResult
}

func HandleNoEmitOptions(ctx context.Context, program AnyProgram, file *ast.SourceFile) *EmitResult {
	options := program.Options()
	if options.NoEmit.IsTrue() {
		return &EmitResult{
			EmitSkipped: true,
		}
	}

	if !options.NoEmitOnError.IsTrue() {
		return nil // No emit on error is not set, so we can proceed with emitting
	}

	diagnostics := GetDiagnosticsOfAnyProgram(ctx, program, file, true, func(name string, start bool, nameStart time.Time) time.Time { return time.Time{} })
	if len(diagnostics) == 0 {
		return nil // No diagnostics, so we can proceed with emitting
	}
	return &EmitResult{
		Diagnostics: diagnostics,
		EmitSkipped: true,
	}
}

func GetDiagnosticsOfAnyProgram(
	ctx context.Context,
	program AnyProgram,
	file *ast.SourceFile,
	skipNoEmitCheckForDtsDiagnostics bool,
	recordTime func(name string, start bool, nameStart time.Time) time.Time,
) []*ast.Diagnostic {
	allDiagnostics := slices.Clip(program.GetConfigFileParsingDiagnostics())
	configFileParsingDiagnosticsLength := len(allDiagnostics)

	allDiagnostics = append(allDiagnostics, program.GetSyntacticDiagnostics(ctx, file)...)

	if len(allDiagnostics) == configFileParsingDiagnosticsLength {
		// Options diagnostics include global diagnostics (even though we collect them separately),
		// and global diagnostics create checkers, which then bind all of the files. Do this binding
		// early so we can track the time.
		bindStart := recordTime("bind", true, time.Time{})
		_ = program.GetBindDiagnostics(ctx, file)
		recordTime("bind", false, bindStart)

		allDiagnostics = append(allDiagnostics, program.GetOptionsDiagnostics(ctx)...)

		if program.Options().ListFilesOnly.IsFalseOrUnknown() {
			allDiagnostics = append(allDiagnostics, program.GetGlobalDiagnostics(ctx)...)

			if len(allDiagnostics) == configFileParsingDiagnosticsLength {
				// !!! add program diagnostics here instead of merging with the semantic diagnostics for better api usage with with incremental and
				checkStart := recordTime("check", true, time.Time{})
				allDiagnostics = append(allDiagnostics, program.GetSemanticDiagnostics(ctx, file)...)
				recordTime("check", false, checkStart)
			}

			if (skipNoEmitCheckForDtsDiagnostics || program.Options().NoEmit.IsTrue()) && program.Options().GetEmitDeclarations() && len(allDiagnostics) == configFileParsingDiagnosticsLength {
				allDiagnostics = append(allDiagnostics, program.GetDeclarationDiagnostics(ctx, file)...)
			}
		}
	}
	return allDiagnostics
}

func CombineEmitResults(results []*EmitResult) *EmitResult {
	result := &EmitResult{}
	for _, emitter := range results {
		if emitter.EmitSkipped {
			result.EmitSkipped = true
		}
		result.Diagnostics = append(result.Diagnostics, emitter.Diagnostics...)
		if emitter.EmittedFiles != nil {
			result.EmittedFiles = append(result.EmittedFiles, emitter.EmittedFiles...)
		}
		if emitter.SourceMaps != nil {
			result.SourceMaps = append(result.SourceMaps, emitter.SourceMaps...)
		}
	}
	return result
}
