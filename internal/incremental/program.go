package incremental

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/outputpaths"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type Program struct {
	snapshot *snapshot
	program  *compiler.Program
}

var _ compiler.AnyProgram = (*Program)(nil)

func NewProgram(program *compiler.Program, oldProgram *Program) *Program {
	return &Program{
		snapshot: newSnapshotForProgram(program, oldProgram),
		program:  program,
	}
}

func (h *Program) panicIfNoProgram(method string) {
	if h.program == nil {
		panic(fmt.Sprintf("%s should not be called without program", method))
	}
}

func (h *Program) GetProgram() *compiler.Program {
	h.panicIfNoProgram("GetProgram")
	return h.program
}

// Options implements compiler.AnyProgram interface.
func (h *Program) Options() *core.CompilerOptions {
	return h.snapshot.options
}

// GetSourceFiles implements compiler.AnyProgram interface.
func (h *Program) GetSourceFiles() []*ast.SourceFile {
	h.panicIfNoProgram("GetSourceFiles")
	return h.program.GetSourceFiles()
}

// GetConfigFileParsingDiagnostics implements compiler.AnyProgram interface.
func (h *Program) GetConfigFileParsingDiagnostics() []*ast.Diagnostic {
	h.panicIfNoProgram("GetConfigFileParsingDiagnostics")
	return h.program.GetConfigFileParsingDiagnostics()
}

// GetSyntacticDiagnostics implements compiler.AnyProgram interface.
func (h *Program) GetSyntacticDiagnostics(ctx context.Context, file *ast.SourceFile) []*ast.Diagnostic {
	h.panicIfNoProgram("GetSyntacticDiagnostics")
	return h.program.GetSyntacticDiagnostics(ctx, file)
}

// GetBindDiagnostics implements compiler.AnyProgram interface.
func (h *Program) GetBindDiagnostics(ctx context.Context, file *ast.SourceFile) []*ast.Diagnostic {
	h.panicIfNoProgram("GetBindDiagnostics")
	return h.program.GetBindDiagnostics(ctx, file)
}

// GetOptionsDiagnostics implements compiler.AnyProgram interface.
func (h *Program) GetOptionsDiagnostics(ctx context.Context) []*ast.Diagnostic {
	h.panicIfNoProgram("GetOptionsDiagnostics")
	return h.program.GetOptionsDiagnostics(ctx)
}

func (h *Program) GetProgramDiagnostics() []*ast.Diagnostic {
	h.panicIfNoProgram("GetProgramDiagnostics")
	return h.program.GetProgramDiagnostics()
}

func (h *Program) GetGlobalDiagnostics(ctx context.Context) []*ast.Diagnostic {
	h.panicIfNoProgram("GetGlobalDiagnostics")
	return h.program.GetGlobalDiagnostics(ctx)
}

// GetSemanticDiagnostics implements compiler.AnyProgram interface.
func (h *Program) GetSemanticDiagnostics(ctx context.Context, file *ast.SourceFile) []*ast.Diagnostic {
	h.panicIfNoProgram("GetSemanticDiagnostics")
	if h.snapshot.options.NoCheck.IsTrue() {
		return nil
	}

	// Ensure all the diagnsotics are cached
	h.collectSemanticDiagnosticsOfAffectedFiles(ctx, file)
	if ctx.Err() != nil {
		return nil
	}

	// Return result from cache
	if file != nil {
		cachedDiagnostics, ok := h.snapshot.semanticDiagnosticsPerFile[file.Path()]
		if !ok {
			panic("After handling all the affected files, there shouldnt be more changes")
		}
		return compiler.FilterNoEmitSemanticDiagnostics(cachedDiagnostics.getDiagnostics(h.program, file), h.snapshot.options)
	}

	var diagnostics []*ast.Diagnostic
	for _, file := range h.program.GetSourceFiles() {
		cachedDiagnostics, ok := h.snapshot.semanticDiagnosticsPerFile[file.Path()]
		if !ok {
			panic("After handling all the affected files, there shouldnt be more changes")
		}
		diagnostics = append(diagnostics, compiler.FilterNoEmitSemanticDiagnostics(cachedDiagnostics.getDiagnostics(h.program, file), h.snapshot.options)...)
	}
	return diagnostics
}

// GetDeclarationDiagnostics implements compiler.AnyProgram interface.
func (h *Program) GetDeclarationDiagnostics(ctx context.Context, file *ast.SourceFile) []*ast.Diagnostic {
	h.panicIfNoProgram("GetDeclarationDiagnostics")
	result := emitFiles(ctx, h, compiler.EmitOptions{
		TargetSourceFile: file,
	}, true)
	if result != nil {
		return result.Diagnostics
	}
	return nil
}

// GetModeForUsageLocation implements compiler.AnyProgram interface.
func (h *Program) Emit(ctx context.Context, options compiler.EmitOptions) *compiler.EmitResult {
	h.panicIfNoProgram("Emit")

	var result *compiler.EmitResult
	if h.snapshot.options.NoEmit.IsTrue() {
		result = &compiler.EmitResult{EmitSkipped: true}
	} else {
		result = compiler.HandleNoEmitOnError(ctx, h, options.TargetSourceFile)
		if ctx.Err() != nil {
			return nil
		}
	}
	if result != nil {
		if options.TargetSourceFile != nil {
			return result
		}

		// Emit buildInfo and combine result
		buildInfoResult := h.emitBuildInfo(ctx, options)
		if buildInfoResult != nil && buildInfoResult.EmittedFiles != nil {
			result.Diagnostics = append(result.Diagnostics, buildInfoResult.Diagnostics...)
			result.EmittedFiles = append(result.EmittedFiles, buildInfoResult.EmittedFiles...)
		}
		return result
	}
	return emitFiles(ctx, h, options, false)
}

// Handle affected files and cache the semantic diagnostics for all of them or the file asked for
func (h *Program) collectSemanticDiagnosticsOfAffectedFiles(ctx context.Context, file *ast.SourceFile) {
	// Get all affected files
	collectAllAffectedFiles(ctx, h)
	if ctx.Err() != nil {
		return
	}

	if len(h.snapshot.semanticDiagnosticsPerFile) == len(h.program.GetSourceFiles()) {
		// If we have all the files,
		return
	}

	var affectedFiles []*ast.SourceFile
	if file != nil {
		_, ok := h.snapshot.semanticDiagnosticsPerFile[file.Path()]
		if ok {
			return
		}
		affectedFiles = []*ast.SourceFile{file}
	} else {
		for _, file := range h.program.GetSourceFiles() {
			if _, ok := h.snapshot.semanticDiagnosticsPerFile[file.Path()]; !ok {
				affectedFiles = append(affectedFiles, file)
			}
		}
	}

	// Get their diagnostics and cache them
	diagnosticsPerFile := h.program.GetSemanticDiagnosticsNoFilter(ctx, affectedFiles)
	// commit changes if no err
	if ctx.Err() != nil {
		return
	}

	// Commit changes to snapshot
	for file, diagnostics := range diagnosticsPerFile {
		h.snapshot.semanticDiagnosticsPerFile[file.Path()] = &diagnosticsOrBuildInfoDiagnosticsWithFileName{diagnostics: diagnostics}
	}
	if len(h.snapshot.semanticDiagnosticsPerFile) == len(h.program.GetSourceFiles()) && h.snapshot.checkPending && !h.snapshot.options.NoCheck.IsTrue() {
		h.snapshot.checkPending = false
	}
	h.snapshot.buildInfoEmitPending = true
}

func (h *Program) emitBuildInfo(ctx context.Context, options compiler.EmitOptions) *compiler.EmitResult {
	buildInfoFileName := outputpaths.GetBuildInfoFileName(h.snapshot.options, tspath.ComparePathsOptions{
		CurrentDirectory:          h.program.GetCurrentDirectory(),
		UseCaseSensitiveFileNames: h.program.UseCaseSensitiveFileNames(),
	})
	if buildInfoFileName == "" {
		return nil
	}

	hasErrors := h.ensureHasErrorsForState(ctx, h.program)
	if !h.snapshot.buildInfoEmitPending && h.snapshot.hasErrors == hasErrors {
		return nil
	}
	h.snapshot.hasErrors = hasErrors
	h.snapshot.buildInfoEmitPending = true
	if ctx.Err() != nil {
		return nil
	}
	buildInfo := snapshotToBuildInfo(h.snapshot, h.program, buildInfoFileName)
	text, err := json.Marshal(buildInfo)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal build info: %v", err))
	}
	if options.WriteFile != nil {
		err = options.WriteFile(buildInfoFileName, string(text), false, &compiler.WriteFileData{
			BuildInfo: &buildInfo,
		})
	} else {
		err = h.program.Host().FS().WriteFile(buildInfoFileName, string(text), false)
	}
	if err != nil {
		return &compiler.EmitResult{
			EmitSkipped: true,
			Diagnostics: []*ast.Diagnostic{
				ast.NewCompilerDiagnostic(diagnostics.Could_not_write_file_0_Colon_1, buildInfoFileName, err.Error()),
			},
		}
	}
	h.snapshot.buildInfoEmitPending = false

	var emittedFiles []string
	if h.snapshot.options.ListEmittedFiles.IsTrue() {
		emittedFiles = []string{buildInfoFileName}
	}
	return &compiler.EmitResult{
		EmitSkipped:  false,
		EmittedFiles: emittedFiles,
	}
}

func (h *Program) ensureHasErrorsForState(ctx context.Context, program *compiler.Program) core.Tristate {
	if h.snapshot.hasErrors != core.TSUnknown {
		return h.snapshot.hasErrors
	}

	// Check semantic and emit diagnostics first as we dont need to ask program about it
	if slices.ContainsFunc(program.GetSourceFiles(), func(file *ast.SourceFile) bool {
		semanticDiagnostics := h.snapshot.semanticDiagnosticsPerFile[file.Path()]
		if semanticDiagnostics == nil {
			// Missing semantic diagnostics in cache will be encoded in incremental buildInfo
			return h.snapshot.options.IsIncremental()
		}
		if len(semanticDiagnostics.diagnostics) > 0 || len(semanticDiagnostics.buildInfoDiagnostics) > 0 {
			// cached semantic diagnostics will be encoded in buildInfo
			return true
		}
		if _, ok := h.snapshot.emitDiagnosticsPerFile[file.Path()]; ok {
			// emit diagnostics will be encoded in buildInfo;
			return true
		}
		return false
	}) {
		// Because semantic diagnostics are recorded in buildInfo, we dont need to encode hasErrors in incremental buildInfo
		// But encode as errors in non incremental buildInfo
		return core.IfElse(h.snapshot.options.IsIncremental(), core.TSFalse, core.TSTrue)
	}
	if len(program.GetConfigFileParsingDiagnostics()) > 0 ||
		len(program.GetSyntacticDiagnostics(ctx, nil)) > 0 ||
		len(program.GetBindDiagnostics(ctx, nil)) > 0 ||
		len(program.GetOptionsDiagnostics(ctx)) > 0 {
		return core.TSTrue
	} else {
		return core.TSFalse
	}
}
