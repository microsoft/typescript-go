package transpile_test

import (
	"context"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/transpile"
	"gotest.tools/v3/assert"
)

func TestTranspileModule(t *testing.T) {
	t.Parallel()

	t.Run("commonJS", func(t *testing.T) {
		t.Parallel()
		output := transpile.TranspileModule(t.Context(), "export const x = 0;", transpile.Options{
			CompilerOptions: &core.CompilerOptions{
				Module:  core.ModuleKindCommonJS,
				NewLine: core.NewLineKindLF,
			},
			ReportDiagnostics: true,
		})
		assert.Equal(t, len(output.Diagnostics), 0)
		assert.Equal(t, output.OutputText, "\"use strict\";\nObject.defineProperty(exports, \"__esModule\", { value: true });\nexports.x = void 0;\nexports.x = 0;\n")
	})

	t.Run("diagnostics", func(t *testing.T) {
		t.Parallel()
		withoutDiagnostics := transpile.TranspileModule(t.Context(), "a b", transpile.Options{})
		assert.Assert(t, withoutDiagnostics.Diagnostics == nil)

		withDiagnostics := transpile.TranspileModule(t.Context(), "a b", transpile.Options{ReportDiagnostics: true})
		assert.Assert(t, len(withDiagnostics.Diagnostics) > 0)

		semanticError := transpile.TranspileModule(t.Context(), "const x: string = 0;", transpile.Options{ReportDiagnostics: true})
		assert.Equal(t, len(semanticError.Diagnostics), 0)
	})

	t.Run("tsxDefaultFileName", func(t *testing.T) {
		t.Parallel()
		output := transpile.TranspileModule(t.Context(), "const x = <div />;", transpile.Options{
			CompilerOptions:   &core.CompilerOptions{Jsx: core.JsxEmitReact},
			ReportDiagnostics: true,
		})
		assert.Equal(t, len(output.Diagnostics), 0)
		assert.Assert(t, strings.Contains(output.OutputText, `React.createElement("div", null)`))
	})

	t.Run("sourceMap", func(t *testing.T) {
		t.Parallel()
		output := transpile.TranspileModule(t.Context(), "export const x = 0;", transpile.Options{
			CompilerOptions: &core.CompilerOptions{SourceMap: core.TSTrue},
		})
		assert.Assert(t, strings.Contains(output.SourceMapText, `"version":3`))
	})

	t.Run("overridesTranspileOptions", func(t *testing.T) {
		t.Parallel()
		output := transpile.TranspileModule(t.Context(), "export const x = 0;", transpile.Options{
			CompilerOptions: &core.CompilerOptions{
				NoEmit:        core.TSTrue,
				NoEmitOnError: core.TSTrue,
				Declaration:   core.TSTrue,
			},
		})
		assert.Assert(t, output.OutputText != "")
		assert.Assert(t, !strings.Contains(output.OutputText, "declare"))
	})

	t.Run("cancelled", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		assert.Assert(t, transpile.TranspileModule(ctx, "export const x = 0;", transpile.Options{}) == nil)
	})
}

func TestTranspileDeclaration(t *testing.T) {
	t.Parallel()

	t.Run("annotatedFunction", func(t *testing.T) {
		t.Parallel()
		output := transpile.TranspileDeclaration(t.Context(), "export function f(x: number): number { return x + 1; }", transpile.Options{
			CompilerOptions:   &core.CompilerOptions{NewLine: core.NewLineKindLF},
			ReportDiagnostics: true,
		})
		assert.Equal(t, len(output.Diagnostics), 0)
		assert.Equal(t, output.OutputText, "export declare function f(x: number): number;\n")
	})

	t.Run("emitDiagnosticsAlwaysReported", func(t *testing.T) {
		t.Parallel()
		output := transpile.TranspileDeclaration(t.Context(), "export function f(x) { return x + 1; }", transpile.Options{})
		assert.Assert(t, len(output.Diagnostics) > 0)
		assert.Assert(t, output.OutputText != "")
	})

	t.Run("barebonesLibSupportsArrayInference", func(t *testing.T) {
		t.Parallel()
		output := transpile.TranspileDeclaration(t.Context(), "export const values = [1, 2, 3];", transpile.Options{
			CompilerOptions: &core.CompilerOptions{NewLine: core.NewLineKindLF},
		})
		assert.Equal(t, output.OutputText, "export declare const values: number[];\n")
	})

	t.Run("declarationMap", func(t *testing.T) {
		t.Parallel()
		output := transpile.TranspileDeclaration(t.Context(), "export const x: number = 0;", transpile.Options{
			CompilerOptions: &core.CompilerOptions{
				DeclarationMap: core.TSTrue,
				NewLine:        core.NewLineKindLF,
			},
		})
		assert.Assert(t, strings.Contains(output.OutputText, "//# sourceMappingURL=module.d.ts.map"))
		assert.Assert(t, strings.Contains(output.SourceMapText, `"version":3`))
	})
}
