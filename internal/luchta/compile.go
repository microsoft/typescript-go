package luchta

import (
	"bytes"
	"context"
	"path/filepath"
	"sort"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
	"github.com/microsoft/typescript-go/internal/locale"
	"github.com/microsoft/typescript-go/internal/pnp"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
	"github.com/microsoft/typescript-go/internal/vfs/pnpvfs"
)

// CompileResult holds the outcome of compiling one package.
type CompileResult struct {
	ExitCode    int
	Inputs      []string
	Outputs     []string
	Diagnostics string
}

var tsconfigCandidates = []string{"tsconfig.build.json", "tsconfig.json"}

// CompilePackage replicates the project's tsc worker: for each candidate tsconfig
// that exists, parse it, clean stale outputs, build a Program, collect diagnostics,
// and emit. Returns aggregated inputs/outputs and a non-zero ExitCode on any error.
func CompilePackage(ctx context.Context, cwd string) CompileResult {
	fsys := bundled.WrapFS(osvfs.FS())
	pnpApi := pnp.InitPnpApi(fsys, cwd) // nil when there is no .pnp.cjs above cwd
	if pnpApi != nil {
		fsys = pnpvfs.From(fsys)
	}
	var diagBuf bytes.Buffer
	sys := newRunSystem(cwd, fsys, bundled.LibPath(), &diagBuf, pnpApi)
	extendedConfigCache := &tsc.ExtendedConfigCache{}

	inputs := collections.NewSetWithSizeHint[string](4)
	var outputs []string
	exitCode := 0

	for _, name := range tsconfigCandidates {
		configPath := filepath.Join(cwd, name)
		if !fsys.FileExists(configPath) {
			continue
		}
		inputs.Add(name)

		parsed, errs := tsoptions.GetParsedCommandLineOfConfigFile(
			configPath, &core.CompilerOptions{}, nil, sys, extendedConfigCache,
		)
		if len(errs) > 0 {
			reportDiagnostics(sys, &diagBuf, parsed, errs)
			exitCode = 1
			break
		}
		for _, p := range includePatterns(parsed) {
			inputs.Add(p)
		}

		opts := parsed.CompilerOptions()
		if err := CleanOutputs(cwd, outDirOf(opts), rootDirOf(opts), opts.NoEmit.IsTrue()); err != nil {
			diagBuf.WriteString("clean outputs failed: " + err.Error() + "\n")
			exitCode = 1
			break
		}

		host := compiler.NewCachedFSCompilerHost(cwd, fsys, sys.DefaultLibraryPath(), extendedConfigCache, pnpApi, nil)
		program := compiler.NewProgram(compiler.ProgramOptions{Config: parsed, Host: host})

		diags := collectAllDiagnostics(ctx, program)

		emitResult := program.Emit(ctx, compiler.EmitOptions{
			WriteFile: func(fileName, text string, data *compiler.WriteFileData) error {
				outputs = append(outputs, fileName)
				return osvfs.FS().WriteFile(fileName, text)
			},
		})
		diags = append(diags, emitResult.Diagnostics...)

		if len(diags) > 0 {
			reportDiagnostics(sys, &diagBuf, parsed, diags)
			exitCode = 1
			break
		}
	}

	if inputs.Len() == 0 {
		inputs.Add("src/**")
	}
	return CompileResult{
		ExitCode:    exitCode,
		Inputs:      sortedKeys(inputs),
		Outputs:     RelativizeOutputs(cwd, outputs),
		Diagnostics: diagBuf.String(),
	}
}

func collectAllDiagnostics(ctx context.Context, p *compiler.Program) []*ast.Diagnostic {
	var d []*ast.Diagnostic
	d = append(d, p.GetConfigFileParsingDiagnostics()...)
	d = append(d, p.GetSyntacticDiagnostics(ctx, nil)...)
	d = append(d, p.GetProgramDiagnostics()...)
	d = append(d, p.GetBindDiagnostics(ctx, nil)...)
	d = append(d, p.GetGlobalDiagnostics(ctx)...)
	d = append(d, p.GetSemanticDiagnostics(ctx, nil)...)
	d = append(d, p.GetDeclarationDiagnostics(ctx, nil)...)
	return d
}

func reportDiagnostics(sys tsc.System, w *bytes.Buffer, parsed *tsoptions.ParsedCommandLine, diags []*ast.Diagnostic) {
	var opts *core.CompilerOptions
	if parsed != nil {
		opts = parsed.CompilerOptions()
	}
	report := tsc.CreateDiagnosticReporter(sys, w, locale.Default, opts)
	for _, d := range diags {
		report(d)
	}
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
