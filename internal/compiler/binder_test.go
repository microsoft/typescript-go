package compiler

import (
	"runtime"
	"testing"
)

func BenchmarkBind(b *testing.B) {
	srcCompilerCheckerTS.SkipIfNotExist(b)

	fileName := srcCompilerCheckerTS.Path()
	sourceText := srcCompilerCheckerTS.ReadFile(b)

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
