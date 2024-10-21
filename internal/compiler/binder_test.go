package compiler

import (
	"os"
	"runtime"
	"testing"
)

func BenchmarkBind(b *testing.B) {
	fileName := "../../_submodules/TypeScript/src/compiler/checker.ts"
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		b.Error(err)
	}
	sourceText := string(bytes)

	sourceFiles := make([]*SourceFile, b.N)
	for i := 0; i < b.N; i++ {
		sourceFiles[i] = ParseSourceFile(fileName, sourceText, ScriptTargetESNext)
	}

	compilerOptions := &CompilerOptions{Target: ScriptTargetESNext, ModuleKind: ModuleKindNodeNext}

	// The above parses do a lot of work; ensure GC is settled before we start collecting pefrormance data.
	runtime.GC()
	runtime.GC()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bindSourceFile(sourceFiles[i], compilerOptions)
	}
}
