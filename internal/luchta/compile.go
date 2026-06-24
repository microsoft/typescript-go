package luchta

import (
	"context"
	"io"
	"path/filepath"
	"sort"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
	"github.com/microsoft/typescript-go/internal/pnp"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
	"github.com/microsoft/typescript-go/internal/vfs/pnpvfs"
)

// CompileResult holds the outcome of compiling one package.
type CompileResult struct {
	ExitCode int
	Inputs   []string
	Outputs  []string
	// Diagnostics are the raw compiler diagnostics, rendered into a SARIF
	// report by the worker (see DiagnosticsToSARIF).
	Diagnostics []*ast.Diagnostic
	// InternalError is set for an operational failure that is not a code-level
	// diagnostic (e.g. failing to clean stale outputs). The worker logs it to
	// stderr rather than placing it in the SARIF report.
	InternalError string
}

var tsconfigCandidates = []string{"tsconfig.build.json", "tsconfig.json"}

// compilerFS builds the vfs for compiling under cwd, enabling Yarn PnP when a
// .pnp.cjs manifest exists at or above cwd. Returns the FS and the PnP API (nil
// when not a PnP workspace).
func compilerFS(cwd string) (vfs.FS, *pnp.PnpApi) {
	fsys := bundled.WrapFS(osvfs.FS())
	pnpApi := pnp.InitPnpApi(fsys, cwd) // nil when there is no .pnp.cjs above cwd
	if pnpApi != nil {
		fsys = pnpvfs.From(fsys)
	}
	return fsys, pnpApi
}

// CompilePackage replicates the project's tsc worker: for each candidate tsconfig
// that exists, parse it, clean stale outputs, build a Program, collect diagnostics,
// and emit. Returns aggregated inputs/outputs and a non-zero ExitCode on any error.
func CompilePackage(ctx context.Context, cwd string) CompileResult {
	fsys, pnpApi := compilerFS(cwd)
	sys := newRunSystem(cwd, fsys, bundled.LibPath(), io.Discard, pnpApi)
	extendedConfigCache := &tsc.ExtendedConfigCache{}

	inputs := collections.NewSetWithSizeHint[string](4)
	var outputs []string
	var allDiags []*ast.Diagnostic
	internalErr := ""
	exitCode := 0

	// Report every tsconfig the build *could* parse as a literal input, even
	// when it is currently absent. luchta records a missing literal as an
	// "absent" sentinel, so adding the file later flips absent->present and
	// busts the cache. Reporting only the existing candidate would let a newly
	// added tsconfig.build.json (which this worker would then parse) go
	// undetected and the cache would incorrectly skip the rebuild.
	for _, name := range tsconfigCandidates {
		inputs.Add(name)
	}

	for _, name := range tsconfigCandidates {
		configPath := filepath.Join(cwd, name)
		if !fsys.FileExists(configPath) {
			continue
		}

		parsed, errs := tsoptions.GetParsedCommandLineOfConfigFile(
			configPath, &core.CompilerOptions{}, nil, sys, extendedConfigCache,
		)
		if len(errs) > 0 {
			allDiags = append(allDiags, errs...)
			exitCode = 1
			break
		}
		for _, p := range includePatterns(parsed) {
			inputs.Add(p)
		}

		opts := parsed.CompilerOptions()
		if err := CleanOutputs(cwd, outDirOf(opts), rootDirOf(opts), opts.NoEmit.IsTrue()); err != nil {
			internalErr = "clean outputs failed: " + err.Error()
			exitCode = 1
			break
		}

		host := compiler.NewCachedFSCompilerHost(cwd, fsys, sys.DefaultLibraryPath(), extendedConfigCache, pnpApi, nil)
		program := compiler.NewProgram(compiler.ProgramOptions{Config: parsed, Host: host})

		diags := collectAllDiagnostics(ctx, program)

		emitResult := program.Emit(ctx, compiler.EmitOptions{})
		outputs = append(outputs, emitResult.EmittedFiles...)
		diags = append(diags, emitResult.Diagnostics...)

		if len(diags) > 0 {
			allDiags = append(allDiags, diags...)
			exitCode = 1
			break
		}
	}

	relativized := RelativizeOutputs(cwd, outputs)
	sort.Strings(relativized)
	return CompileResult{
		ExitCode:      exitCode,
		Inputs:        sortedKeys(inputs),
		Outputs:       relativized,
		Diagnostics:   allDiags,
		InternalError: internalErr,
	}
}

// collectAllDiagnostics gathers pre-emit diagnostics exactly the way the tsc CLI
// does (via compiler.GetDiagnosticsOfAnyProgram), so skipLibCheck and the
// "don't report semantic errors when there are syntactic errors" behavior are
// honored. Binder/duplicate-identifier errors surface through the semantic pass,
// which filters declaration files under skipLibCheck; appending raw
// GetBindDiagnostics(nil) here (as this previously did) wrongly reports errors
// inside library .d.ts files — e.g. playwright-core's intentionally duplicated
// `ElectronType` — even when skipLibCheck is set. Declaration-emit diagnostics
// are collected separately from the subsequent program.Emit call.
func collectAllDiagnostics(ctx context.Context, p *compiler.Program) []*ast.Diagnostic {
	return compiler.GetDiagnosticsOfAnyProgram(
		ctx,
		p,
		nil,   // file == nil: gather across all (non-skipped) files
		false, // skipNoEmitCheckForDtsDiagnostics
		p.GetBindDiagnostics,
		p.GetSemanticDiagnostics,
	)
}

func outDirOf(o *core.CompilerOptions) string {
	if o != nil && o.OutDir != "" {
		return o.OutDir
	}
	return "dist/types"
}

func rootDirOf(o *core.CompilerOptions) string {
	if o != nil && o.RootDir != "" {
		return o.RootDir
	}
	return "src"
}

// includePatterns extracts the tsconfig "include" globs from the raw config,
// defaulting to ["src/**"] (mirrors the JS worker's input reporting).
func includePatterns(parsed *tsoptions.ParsedCommandLine) []string {
	def := []string{"src/**"}
	if parsed == nil || parsed.Raw == nil {
		return def
	}
	var inc any
	switch raw := parsed.Raw.(type) {
	case map[string]any:
		inc = raw["include"]
	case *collections.OrderedMap[string, any]:
		if v, ok := raw.Get("include"); ok {
			inc = v
		}
	}
	arr, ok := inc.([]any)
	if !ok || len(arr) == 0 {
		return def
	}
	out := make([]string, 0, len(arr))
	for _, v := range arr {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	if len(out) == 0 {
		return def
	}
	return out
}

func sortedKeys(s *collections.Set[string]) []string {
	keys := make([]string, 0, s.Len())
	for k := range s.Keys() {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
