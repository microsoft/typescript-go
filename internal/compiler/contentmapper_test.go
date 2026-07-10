package compiler_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/contentmapper"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/diagnosticwriter"
	"github.com/microsoft/typescript-go/internal/spanmap"
	"github.com/microsoft/typescript-go/internal/testutil/baseline"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

type fakeContentMapperRunner struct {
	transform func(fileName string, content string) (compiler.ContentMapperResult, error)
}

func (r fakeContentMapperRunner) Transform(fileName string, content string) (compiler.ContentMapperResult, error) {
	return r.transform(fileName, content)
}

const vueRawContent = "<template>original</template>"

// synthesizedSpanMap maps an entirely generated transformed text to no original location (a fully
// synthesized file). It satisfies the span map contract: it tiles the transformed text and its (empty)
// original span stays within the original content.
func synthesizedSpanMap(transformed string) *spanmap.SpanMap {
	return spanmap.New([]spanmap.Segment{{
		GenStart:  0,
		GenEnd:    core.TextPos(len(transformed)),
		OrigStart: 0,
		OrigEnd:   0,
		Kind:      spanmap.KindSynthesized,
	}})
}

func newContentMapperProgram(t *testing.T, runner compiler.ContentMapperRunner, files map[string]string, rootFiles []string) *compiler.Program {
	t.Helper()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}
	fs := vfstest.FromMap[any](nil, false /*useCaseSensitiveFileNames*/)
	fs = bundled.WrapFS(fs)
	for name, content := range files {
		_ = fs.WriteFile(name, content)
	}

	config := &tsoptions.ParsedCommandLine{
		ParsedConfig: &core.ParsedOptions{
			FileNames: rootFiles,
			CompilerOptions: &core.CompilerOptions{
				SkipLibCheck:     core.TSTrue,
				Module:           core.ModuleKindESNext,
				ModuleResolution: core.ModuleResolutionKindBundler,
			},
			ContentMappers: []*contentmapper.Mapper{{Definition: contentmapper.Definition{Package: "vue", Extensions: []string{".vue"}}, Manifest: contentmapper.Manifest{Name: "vue-mapper", Version: "1.0.0"}}},
		},
	}
	return compiler.NewProgram(compiler.ProgramOptions{
		Config:              config,
		Host:                compiler.NewCompilerHost("/src", fs, bundled.LibPath(), nil, nil),
		ContentMapperRunner: runner,
		// Load files on the calling goroutine for deterministic diagnostics ordering.
		SingleThreaded: core.TSTrue,
	})
}

func collectContentMapperDiagnostics(program *compiler.Program) []*ast.Diagnostic {
	ctx := context.Background()
	return slices.Concat(
		program.GetSyntacticDiagnostics(ctx, nil),
		program.GetSemanticDiagnostics(ctx, nil),
		program.GetProgramDiagnostics(),
	)
}

func TestContentMapperDiagnostics(t *testing.T) {
	t.Parallel()

	// A mapper reports problems it finds in the file's content (e.g. a syntax error in the embedded
	// script) as result.Diagnostics. A real runner deserializes these from another process: it does
	// not have our diagnostics.Message values, so it builds each one from an already-localized message,
	// a code namespaced by the mapper's own source prefix, and a range into the original content.
	const componentVue = `<template>
  <div>{{ greeting }}</div>
</template>

<script lang="ts">
export const greeting = "hello" oops
</script>
`
	files := map[string]string{
		"/src/app.ts": `import { greeting } from "./Component.vue";
console.log(greeting);`,
		"/src/Component.vue": componentVue,
	}
	runner := fakeContentMapperRunner{
		transform: func(fileName string, content string) (compiler.ContentMapperResult, error) {
			// The mapper turns the <script> block into valid TypeScript but reports the stray token it
			// found, pointing back into the original .vue content.
			start := strings.Index(content, "oops")
			return compiler.ContentMapperResult{
				Text:       "export const greeting = \"hello\";\n",
				ScriptKind: core.ScriptKindTS,
				Mappings:   synthesizedSpanMap("export const greeting = \"hello\";\n"),
				Diagnostics: []*ast.Diagnostic{
					ast.NewExternalDiagnostic(
						nil, // the loader associates the source file
						core.NewTextRange(start, start+len("oops")),
						"vue",
						diagnostics.CategoryError,
						1002,
						"Unexpected token.",
					),
				},
			}, nil
		},
	}
	program := newContentMapperProgram(t, runner, files, []string{"/src/app.ts"})
	baseline.Run(t, "contentMapperDiagnostics.txt", contentMapperBaseline(program, collectContentMapperDiagnostics(program)), baseline.Options{Subfolder: "contentMappers"})
}

func TestContentMapperSpanMapping(t *testing.T) {
	t.Parallel()

	// A <script> block is passed through verbatim: its TypeScript is copied into the transformed text
	// unchanged, so a type error the compiler finds there must be reported at the corresponding line in
	// the original .vue file, not at its (shifted) position in the transformed text. The mapper records
	// a verbatim span mapping so the compiler diagnostic is mapped back exactly.
	const scriptBody = "const greeting: number = \"not a number\";\nexport { greeting };\n"
	componentVue := "<template>\n  <div>{{ greeting }}</div>\n</template>\n\n<script lang=\"ts\">\n" + scriptBody + "</script>\n"
	scriptStart := strings.Index(componentVue, scriptBody)

	files := map[string]string{
		"/src/app.ts": `import { greeting } from "./Component.vue";
console.log(greeting);`,
		"/src/Component.vue": componentVue,
	}
	runner := fakeContentMapperRunner{
		transform: func(fileName string, content string) (compiler.ContentMapperResult, error) {
			return compiler.ContentMapperResult{
				Text:       scriptBody,
				ScriptKind: core.ScriptKindTS,
				Mappings: spanmap.New([]spanmap.Segment{{
					GenStart:  0,
					GenEnd:    core.TextPos(len(scriptBody)),
					OrigStart: core.TextPos(scriptStart),
					OrigEnd:   core.TextPos(scriptStart + len(scriptBody)),
					Kind:      spanmap.KindVerbatim,
				}}),
			}, nil
		},
	}
	program := newContentMapperProgram(t, runner, files, []string{"/src/app.ts"})
	baseline.Run(t, "contentMapperSpanMapping.txt", contentMapperBaseline(program, collectContentMapperDiagnostics(program)), baseline.Options{Subfolder: "contentMappers"})
}

func TestContentMapperGeneratedCode(t *testing.T) {
	t.Parallel()

	// Some transformed code is synthesized by the mapper and has no counterpart in the original file
	// (e.g. a generated runtime call). A compiler error there cannot be pointed at the original, so it is
	// shown against the transformed text with a note that the location is generated.
	const transformed = "export const el = jsxRuntime(Widget);\n"
	files := map[string]string{
		"/src/app.ts":        `import "./Component.vue";`,
		"/src/Component.vue": "<template>\n  <Widget />\n</template>\n",
	}
	runner := fakeContentMapperRunner{
		transform: func(fileName string, content string) (compiler.ContentMapperResult, error) {
			return compiler.ContentMapperResult{
				Text:       transformed,
				ScriptKind: core.ScriptKindTS,
				Mappings:   synthesizedSpanMap(transformed),
			}, nil
		},
	}
	program := newContentMapperProgram(t, runner, files, []string{"/src/app.ts"})
	baseline.Run(t, "contentMapperGeneratedCode.txt", contentMapperBaseline(program, collectContentMapperDiagnostics(program)), baseline.Options{Subfolder: "contentMappers"})
}

func TestContentMapperInvalidMappings(t *testing.T) {
	t.Parallel()

	// Mappings are a required, enforced part of the content mapper contract. Each malformed map is
	// attributed to the mapper with a specific diagnostic instead of surfacing untrustworthy positions.
	const transformed = "export const x = 1;\n"
	const original = "<template>x</template>\n"

	verbatimAll := func(kind spanmap.Kind, origEnd int) *spanmap.SpanMap {
		return spanmap.New([]spanmap.Segment{{
			GenStart: 0, GenEnd: core.TextPos(len(transformed)),
			OrigStart: 0, OrigEnd: core.TextPos(origEnd), Kind: kind,
		}})
	}

	testCases := []struct {
		name     string
		mappings *spanmap.SpanMap
		wantCode int32
	}{
		{"missing", nil, 100027},
		{
			"coverage",
			spanmap.New([]spanmap.Segment{{GenStart: 0, GenEnd: 3, OrigStart: 0, OrigEnd: 0, Kind: spanmap.KindSynthesized}}),
			100028,
		},
		{
			"outOfBounds",
			verbatimAll(spanmap.KindSynthesized, len(original)+50),
			100029,
		},
		{
			// A verbatim segment whose original text differs from the transformed text.
			"verbatimMismatch",
			verbatimAll(spanmap.KindVerbatim, len(transformed)),
			100030,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			files := map[string]string{
				"/src/app.ts":        `import "./Component.vue";`,
				"/src/Component.vue": original,
			}
			runner := fakeContentMapperRunner{
				transform: func(fileName string, content string) (compiler.ContentMapperResult, error) {
					return compiler.ContentMapperResult{
						Text:       transformed,
						ScriptKind: core.ScriptKindTS,
						Mappings:   tc.mappings,
					}, nil
				},
			}
			program := newContentMapperProgram(t, runner, files, []string{"/src/app.ts"})
			diags := collectContentMapperDiagnostics(program)
			found := slices.ContainsFunc(diags, func(d *ast.Diagnostic) bool { return d.Code() == tc.wantCode })
			assert.Assert(t, found, "expected diagnostic TS%d attributing the invalid mapping, got: %v", tc.wantCode, diags)
		})
	}
}

func TestContentMapperFileFailure(t *testing.T) {
	t.Parallel()

	files := map[string]string{
		"/src/app.ts": `import { greeting } from "./Component.vue";
export const bad: number = greeting;`,
		"/src/Component.vue": vueRawContent,
	}
	runner := fakeContentMapperRunner{
		transform: func(fileName string, content string) (compiler.ContentMapperResult, error) {
			return compiler.ContentMapperResult{}, errors.New("broken pipe")
		},
	}
	program := newContentMapperProgram(t, runner, files, []string{"/src/app.ts"})
	baseline.Run(t, "contentMapperFileFailure.txt", contentMapperBaseline(program, collectContentMapperDiagnostics(program)), baseline.Options{Subfolder: "contentMappers"})
}

func TestContentMapperDisabledAfterRepeatedFailures(t *testing.T) {
	t.Parallel()

	// A mapper that keeps failing is disabled after a bounded number of failures: the individual
	// failures are reported, then a single program diagnostic notes the mapper was disabled, and the
	// remaining files it would have handled are silently substituted with empty files.
	files := map[string]string{
		"/src/app.ts": `import "./A.vue";
import "./B.vue";
import "./C.vue";
import "./D.vue";
import "./E.vue";
import "./F.vue";
import "./G.vue";`,
		"/src/A.vue": vueRawContent,
		"/src/B.vue": vueRawContent,
		"/src/C.vue": vueRawContent,
		"/src/D.vue": vueRawContent,
		"/src/E.vue": vueRawContent,
		"/src/F.vue": vueRawContent,
		"/src/G.vue": vueRawContent,
	}
	runner := fakeContentMapperRunner{
		transform: func(fileName string, content string) (compiler.ContentMapperResult, error) {
			return compiler.ContentMapperResult{}, errors.New("mapper protocol version mismatch")
		},
	}
	program := newContentMapperProgram(t, runner, files, []string{"/src/app.ts"})
	baseline.Run(t, "contentMapperDisabledAfterRepeatedFailures.txt", contentMapperBaseline(program, collectContentMapperDiagnostics(program)), baseline.Options{Subfolder: "contentMappers"})
}

// contentMapperBaseline renders the content-mapped source files (original and transformed text, script
// kind, and associated content mapper) followed by the program's diagnostics as the CLI would show them:
// with a source code frame for each diagnostic. Terminal color escapes are stripped so the baseline
// matches tsc output redirected to a file.
func contentMapperBaseline(program *compiler.Program, diagnostics []*ast.Diagnostic) string {
	var b strings.Builder
	formatOpts := &diagnosticwriter.FormattingOptions{
		NewLine: "\n",
		ComparePathsOptions: tspath.ComparePathsOptions{
			CurrentDirectory:          "/src",
			UseCaseSensitiveFileNames: false,
		},
	}
	relative := func(fileName string) string {
		return tspath.ConvertToRelativePath(fileName, formatOpts.ComparePathsOptions)
	}

	b.WriteString("=== Content-mapped files ===\n")
	for _, file := range program.GetSourceFiles() {
		mapper := program.GetContentMapper(file)
		if mapper == nil {
			continue
		}
		fmt.Fprintf(&b, "\n//// [%s] (ScriptKind: %s, ContentMapper: %v)\n", relative(file.FileName()), file.ScriptKind, mapper.Extensions)
		b.WriteString("--- Original ---\n")
		b.WriteString(file.OriginalText())
		b.WriteString("\n--- Transformed ---\n")
		b.WriteString(file.Text())
		if !strings.HasSuffix(file.Text(), "\n") {
			b.WriteString("\n")
		}
	}

	b.WriteString("\n=== Diagnostics ===\n")
	if len(diagnostics) == 0 {
		b.WriteString("\n" + baseline.NoContent + "\n")
	} else {
		b.WriteString("\n")
		var rendered strings.Builder
		diagnosticwriter.FormatDiagnosticsWithColorAndContext(&rendered, diagnosticwriter.ToDiagnostics(diagnosticwriter.WrapASTDiagnostics(diagnostics)), formatOpts)
		b.WriteString(ansiEscape.ReplaceAllString(rendered.String(), ""))
		b.WriteString("\n")
	}
	return b.String()
}

var ansiEscape = regexp.MustCompile("\x1b\\[[0-9;]*m")
