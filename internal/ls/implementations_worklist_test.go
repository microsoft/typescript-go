package ls

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

// Bug D (withfig): runaway memory in the go-to-implementation worklist.
//
// provideSymbolsAndEntries' implementations worklist (findallreferences.go) dedups the
// worklist only by *node* (seenNodes) and does `implementationEntries = core.Concatenate(...)`
// with NO symbol-level dedup. When an interface member has K implementations, a go-to-
// implementation query runs ~K full program-wide searches (one per implementation node), and
// each search returns all K implementations, so `implementationEntries` accumulates O(K^2)
// references. On a large, deeply-typed file (withfig `wp.ts` + `@withfig/autocomplete-types`)
// this blows up to many GB and the server is OOM-killed ("connection closed prematurely:
// undefined").
//
// The final LSP response dedups by node, so the blow-up is invisible from the response; this
// white-box test inspects the pre-dedup `SymbolsAndEntries` that provideSymbolsAndEntries
// returns. It asserts the accumulated reference count grows *linearly* with K: quadratic
// (current, buggy) => the count roughly quadruples when K doubles; a correct (symbol-deduped)
// implementation stays ~linear (roughly doubles). Pre-fix this test FAILS; post-fix it PASSES.
func TestImplementationsWorklistDoesNotBlowUp(t *testing.T) {
	t.Parallel()

	measure := func(k int) int {
		var b strings.Builder
		b.WriteString("interface I { m(): void; }\n")
		for i := range k {
			fmt.Fprintf(&b, "const a%d: I = { m() {} };\n", i)
		}
		b.WriteString("declare const i: I;\n")
		b.WriteString("i.m();\n")
		content := b.String()

		fs := vfstest.FromMap(map[string]string{
			"/repro.ts":      content,
			"/tsconfig.json": `{ "compilerOptions": {}, "files": ["repro.ts"] }`,
		}, false /*useCaseSensitiveFileNames*/)
		fs = bundled.WrapFS(fs)

		host := compiler.NewCompilerHost("/", fs, bundled.LibPath(), nil, nil)
		parsed, errors := tsoptions.GetParsedCommandLineOfConfigFile("/tsconfig.json", &core.CompilerOptions{}, nil, host, nil)
		assert.Equal(t, len(errors), 0)
		program := compiler.NewProgram(compiler.ProgramOptions{Config: parsed, Host: host})
		program.BindSourceFiles()
		program.GetSemanticDiagnostics(context.Background(), program.GetSourceFile("/repro.ts"))

		sourceFile := program.GetSourceFile("/repro.ts")
		converters := lsconv.NewConverters(lsproto.PositionEncodingKindUTF8, func(_ string) *lsconv.LSPLineMap {
			return lsconv.ComputeLSPLineStarts(content)
		})
		l := &LanguageService{program: program, converters: converters}

		// Position of the `m` property in the final `i.m();`.
		offset := strings.LastIndex(content, "i.m") + len("i.")
		pos := converters.PositionToLineAndCharacter(sourceFile, core.TextPos(offset))

		data, ok := l.provideSymbolsAndEntries(context.Background(), "file:///repro.ts", pos, false /*isRename*/, true /*implementations*/)
		assert.Assert(t, ok)
		total := 0
		for _, se := range data.SymbolsAndEntries {
			total += len(se.references)
		}
		return total
	}

	const k = 40
	small := measure(k)
	large := measure(2 * k)
	t.Logf("accumulated references: K=%d -> %d, K=%d -> %d (ratio %.2f)", k, small, 2*k, large, float64(large)/float64(small))

	// Linear growth ~2x when K doubles; the O(K^2) bug grows ~4x. Fail above 3x.
	assert.Assert(t, large <= small*3,
		"implementations worklist scales superlinearly: K=%d -> %d refs, K=%d -> %d refs (expected ~linear); "+
			"provideSymbolsAndEntries accumulates duplicate entries without symbol-level dedup", k, small, 2*k, large)
}
